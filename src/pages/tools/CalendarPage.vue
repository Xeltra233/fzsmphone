<template>
  <div class="calendar-page">
    <NavBar title="日历" back />

    <!-- Empty -->
    <div v-if="store.calendarEvents.length === 0 && !store.calendarLoading" class="empty-state">
      <div class="empty-icon"><svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="2"><rect x="6" y="8" width="36" height="36" rx="4"/><line x1="16" y1="4" x2="16" y2="12"/><line x1="32" y1="4" x2="32" y2="12"/><line x1="6" y1="20" x2="42" y2="20"/></svg></div>
      <div class="empty-title">日程为空</div>
      <div class="empty-sub">让AI根据世界观生成日程事件</div>
      <button class="gen-btn" :disabled="store.generating" @click="store.generateCalendarContent()">
        ◆ AI 生成日程
      </button>
    </div>

    <!-- Loading -->
    <div v-else-if="store.calendarLoading" class="loading-wrap">
      <div class="loading-spinner"></div>
      <span>正在生成日程...</span>
    </div>

    <!-- Calendar View -->
    <div v-else class="cal-content">
      <!-- Mini Calendar -->
      <div class="mini-cal">
        <div class="cal-header">
          <button class="cal-nav" @click="changeMonth(-1)">‹</button>
          <span class="cal-month">{{ currentYear }}年{{ currentMonth + 1 }}月</span>
          <button class="cal-nav" @click="changeMonth(1)">›</button>
        </div>
        <div class="cal-weekdays">
          <span v-for="d in weekdays" :key="d">{{ d }}</span>
        </div>
        <div class="cal-grid">
          <div
            v-for="day in calDays"
            :key="day.key"
            class="cal-day"
            :class="{ other: day.otherMonth, today: day.isToday, active: day.dateStr === selectedDate, hasEvent: day.hasEvent }"
            @click="selectedDate = day.dateStr"
          >
            <span>{{ day.day }}</span>
            <span v-if="day.hasEvent" class="event-dot"></span>
          </div>
        </div>
      </div>

      <!-- Events for selected date -->
      <div class="events-section">
        <div class="events-title">{{ selectedDate || '选择日期' }} 的事件</div>
        <div v-if="filteredEvents.length === 0" class="no-events">当日无事件</div>
        <div v-for="evt in filteredEvents" :key="evt.id" class="evt-card" @click="selectedEvent = evt">
          <div class="evt-time">{{ evt.time }}</div>
          <div class="evt-info">
            <div class="evt-title">{{ evt.title }}</div>
            <div class="evt-type">{{ evt.type }}</div>
          </div>
        </div>
      </div>

      <!-- Event Detail -->
      <Teleport to="#phone-overlay">
        <div v-if="selectedEvent" class="modal-overlay" @click.self="selectedEvent = null">
          <div class="event-detail-card">
            <div class="edc-title">{{ selectedEvent.title }}</div>
            <div class="edc-row"><span class="edc-label">▫</span>{{ selectedEvent.date }} {{ selectedEvent.time }}</div>
            <div class="edc-row"><span class="edc-label">◇</span>{{ selectedEvent.type }}</div>
            <div class="edc-row"><span class="edc-label">◌</span>{{ selectedEvent.description }}</div>
            <div v-if="selectedEvent.participants.length" class="edc-row">
              <span class="edc-label">◎</span>
              <span v-for="p in selectedEvent.participants" :key="p" class="participant-tag">{{ p }}</span>
            </div>
            <button class="gen-btn close-detail" @click="selectedEvent = null">关闭</button>
          </div>
        </div>
      </Teleport>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'
import type { CalendarEvent } from '@/utils/socialParsers'

const store = useSocialAIStore()
const selectedDate = ref('')
const selectedEvent = ref<CalendarEvent | null>(null)
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

onMounted(() => {
  store.loadData('calendar')
  const today = new Date()
  selectedDate.value = formatDateStr(today)
})

function formatDateStr(d: Date): string {
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
}

function changeMonth(dir: number) {
  currentMonth.value += dir
  if (currentMonth.value < 0) { currentMonth.value = 11; currentYear.value-- }
  if (currentMonth.value > 11) { currentMonth.value = 0; currentYear.value++ }
}

const eventDates = computed(() => {
  const set = new Set<string>()
  store.calendarEvents.forEach(e => set.add(e.date))
  return set
})

const calDays = computed(() => {
  const y = currentYear.value, m = currentMonth.value
  const firstDay = new Date(y, m, 1).getDay()
  const daysInMonth = new Date(y, m + 1, 0).getDate()
  const today = formatDateStr(new Date())
  const days: { key: string; day: number; dateStr: string; otherMonth: boolean; isToday: boolean; hasEvent: boolean }[] = []

  // Previous month padding
  const prevDays = new Date(y, m, 0).getDate()
  for (let i = firstDay - 1; i >= 0; i--) {
    const d = prevDays - i
    const ds = `${y}-${String(m).padStart(2,'0')}-${String(d).padStart(2,'0')}`
    days.push({ key: `p${d}`, day: d, dateStr: ds, otherMonth: true, isToday: false, hasEvent: false })
  }

  // Current month
  for (let d = 1; d <= daysInMonth; d++) {
    const ds = `${y}-${String(m+1).padStart(2,'0')}-${String(d).padStart(2,'0')}`
    days.push({ key: `c${d}`, day: d, dateStr: ds, otherMonth: false, isToday: ds === today, hasEvent: eventDates.value.has(ds) })
  }

  // Next month padding
  const remaining = 42 - days.length
  for (let d = 1; d <= remaining; d++) {
    const ds = `${y}-${String(m+2).padStart(2,'0')}-${String(d).padStart(2,'0')}`
    days.push({ key: `n${d}`, day: d, dateStr: ds, otherMonth: true, isToday: false, hasEvent: false })
  }

  return days
})

const filteredEvents = computed(() => {
  return store.calendarEvents.filter(e => e.date === selectedDate.value)
})
</script>

<style scoped>
.calendar-page { height: 100%; display: flex; flex-direction: column; background: var(--bg-primary); }
.empty-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; padding: 40px 20px; }
.empty-icon { width: 56px; height: 56px; color: var(--text-tertiary); }
.empty-icon svg { width: 100%; height: 100%; }
.empty-title { font-size: 18px; font-weight: 700; color: var(--text-primary); }
.empty-sub { font-size: 13px; color: var(--text-tertiary); }
.gen-btn { padding: 10px 24px; border: none; border-radius: 20px; background: var(--brand-primary); color: #fff; font-size: 14px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.gen-btn:disabled { opacity: 0.5; }
.loading-wrap { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; color: var(--text-secondary); }
.loading-spinner { width: 32px; height: 32px; border: 3px solid var(--separator); border-top-color: var(--brand-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.cal-content { flex: 1; overflow-y: auto; padding-bottom: 80px; }
.mini-cal { padding: 8px 16px; }
.cal-header { display: flex; align-items: center; justify-content: space-between; padding: 8px 0; }
.cal-nav { background: none; border: none; font-size: 24px; color: var(--brand-primary); cursor: pointer; padding: 4px 8px; }
.cal-month { font-size: 16px; font-weight: 700; color: var(--text-primary); }
.cal-weekdays { display: grid; grid-template-columns: repeat(7, 1fr); text-align: center; font-size: 12px; color: var(--text-tertiary); padding: 6px 0; }
.cal-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 2px; }
.cal-day { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 38px; border-radius: 10px; font-size: 14px; cursor: pointer; position: relative; color: var(--text-primary); }
.cal-day.other { color: var(--text-quaternary); }
.cal-day.today { font-weight: 700; background: var(--fill-tertiary); }
.cal-day.active { background: var(--brand-primary); color: #fff; }
.event-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--brand-primary); position: absolute; bottom: 3px; }
.cal-day.active .event-dot { background: #fff; }

.events-section { padding: 12px 16px; }
.events-title { font-size: 15px; font-weight: 700; color: var(--text-primary); margin-bottom: 8px; }
.no-events { font-size: 13px; color: var(--text-tertiary); text-align: center; padding: 20px; }
.evt-card { display: flex; align-items: center; gap: 12px; padding: 12px; background: var(--bg-secondary); border-radius: 12px; margin-bottom: 8px; cursor: pointer; }
.evt-time { font-size: 14px; font-weight: 700; color: var(--brand-primary); min-width: 50px; }
.evt-info { flex: 1; }
.evt-title { font-size: 15px; font-weight: 600; color: var(--text-primary); }
.evt-type { font-size: 12px; color: var(--text-tertiary); margin-top: 2px; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.event-detail-card { width: 90%; max-width: 360px; background: var(--bg-primary); border-radius: 16px; padding: 20px; }
.edc-title { font-size: 18px; font-weight: 700; color: var(--text-primary); margin-bottom: 12px; }
.edc-row { font-size: 14px; color: var(--text-secondary); padding: 4px 0; display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.edc-label { font-size: 14px; color: var(--text-tertiary); }
.participant-tag { padding: 2px 10px; background: var(--brand-primary); color: #fff; border-radius: 10px; font-size: 12px; }
.close-detail { margin-top: 16px; width: 100%; }
</style>
