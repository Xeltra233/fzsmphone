<template>
  <div class="offline-date-page">
    <NavBar title="线下约会" />

    <div class="page-content">
      <!-- 约会状态 -->
      <div class="date-banner">
        <div class="banner-bg"></div>
        <div class="banner-content">
          <span class="date-emoji">♥</span>
          <h2>计划一场甜蜜约会</h2>
          <p>选择约会类型，生成完美计划</p>
        </div>
      </div>

      <!-- 约会类型选择 -->
      <div class="section">
        <h3>选择约会类型</h3>
        <div class="type-grid">
          <div
            v-for="t in dateTypes"
            :key="t.label"
            class="type-card"
            :class="{ selected: selectedType === t.label }"
            @click="selectedType = t.label"
          >
            <span class="type-icon">{{ t.icon }}</span>
            <span class="type-label">{{ t.label }}</span>
          </div>
        </div>
      </div>

      <!-- 生成按钮 -->
      <button class="generate-btn" @click="generatePlan" :disabled="!selectedType || generating">
        {{ generating ? '生成中...' : '✦ 生成约会计划' }}
      </button>

      <!-- 约会计划 -->
      <div v-if="plan" class="plan-card">
        <div class="plan-header">
          <h3>{{ plan.title }}</h3>
          <button class="refresh-btn" @click="generatePlan">↻</button>
        </div>

        <div class="timeline">
          <div v-for="(step, idx) in plan.steps" :key="idx" class="timeline-item">
            <div class="timeline-dot" :class="{ done: step.done }" @click="step.done = !step.done">
              {{ step.done ? '✓' : idx + 1 }}
            </div>
            <div class="timeline-content">
              <div class="step-time">{{ step.time }}</div>
              <div class="step-title">{{ step.title }}</div>
              <div class="step-desc">{{ step.desc }}</div>
            </div>
          </div>
        </div>

        <div class="plan-footer">
          <div class="plan-stat">
            <span>¤</span>
            <span>预算 ¥{{ plan.budget }}</span>
          </div>
          <div class="plan-stat">
            <span>⏱️</span>
            <span>约 {{ plan.duration }}</span>
          </div>
        </div>
      </div>

      <!-- 历史约会 -->
      <div v-if="pastDates.length > 0" class="section">
        <h3>约会记录</h3>
        <div class="past-list">
          <div v-for="d in pastDates" :key="d.id" class="past-item">
            <span class="past-icon">{{ d.icon }}</span>
            <div class="past-info">
              <span class="past-title">{{ d.title }}</span>
              <span class="past-date">{{ d.date }}</span>
            </div>
            <span class="past-rating">{{ '★'.repeat(d.rating) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSettingsStore } from '@/stores/settings'
import { sendAIRequest } from '@/utils/aiService'

const settingsStore = useSettingsStore()

const selectedType = ref('')
const generating = ref(false)

const dateTypes = [
  { icon: '◈', label: '美食约会' },
  { icon: '▷', label: '电影约会' },
  { icon: '△', label: '户外约会' },
  { icon: '▢', label: '逛街约会' },
  { icon: '▶', label: '游戏约会' },
  { icon: '▥', label: '文艺约会' },
  { icon: '◈', label: '演出约会' },
  { icon: '⌂', label: '宅家约会' },
]

interface PlanStep {
  time: string
  title: string
  desc: string
  done: boolean
}

interface DatePlan {
  title: string
  steps: PlanStep[]
  budget: number
  duration: string
}

const plan = ref<DatePlan | null>(null)

interface PastDate {
  id: number
  icon: string
  title: string
  date: string
  rating: number
}

const pastDates = ref<PastDate[]>([
  { id: 1, icon: '', title: '日料探店', date: '2024-12-20', rating: 5 },
  { id: 2, icon: '', title: '看《你的名字》', date: '2024-12-15', rating: 4 },
  { id: 3, icon: '', title: '西湖骑行', date: '2024-11-28', rating: 5 },
])

const planTemplates: Record<string, DatePlan> = {
  '美食约会': {
    title: '◈ 美食探店之旅',
    steps: [
      { time: '17:00', title: '出发准备', desc: '换上好看的衣服，化个淡妆', done: false },
      { time: '17:30', title: '甜品下午茶', desc: '先去一家网红甜品店打卡', done: false },
      { time: '18:30', title: '正式晚餐', desc: '去预定的餐厅享用晚餐', done: false },
      { time: '20:00', title: '散步消食', desc: '沿着河边散步聊天', done: false },
      { time: '21:00', title: '夜市小吃', desc: '买点小吃边走边吃', done: false },
    ],
    budget: 500,
    duration: '4小时',
  },
  '电影约会': {
    title: '▷ 浪漫观影日',
    steps: [
      { time: '14:00', title: '见面碰头', desc: '在商场入口见面', done: false },
      { time: '14:30', title: '买电影票', desc: '选一部两人都感兴趣的电影', done: false },
      { time: '15:00', title: '观看电影', desc: '一起看电影，分享爆米花', done: false },
      { time: '17:30', title: '讨论剧情', desc: '找个咖啡厅聊聊电影感想', done: false },
      { time: '18:30', title: '晚餐', desc: '就近找一家好评餐厅', done: false },
    ],
    budget: 350,
    duration: '5小时',
  },
  '户外约会': {
    title: '△ 阳光户外日',
    steps: [
      { time: '09:00', title: '早起出发', desc: '带上野餐垫和食物', done: false },
      { time: '10:00', title: '到达公园', desc: '选个好位置铺好野餐垫', done: false },
      { time: '10:30', title: '拍照打卡', desc: '在花丛中拍照留念', done: false },
      { time: '12:00', title: '野餐时光', desc: '享受准备好的美食', done: false },
      { time: '14:00', title: '骑行/散步', desc: '环湖骑行或林间散步', done: false },
    ],
    budget: 200,
    duration: '6小时',
  },
}

async function generatePlan() {
  if (!selectedType.value) return
  generating.value = true

  const s = settingsStore.settings
  let aiPlan: DatePlan | null = null

  if (s.apiKey) {
    try {
      const prompt = `Generate a date plan for type "${selectedType.value}". Return ONLY valid JSON:
{"title":"...","steps":[{"time":"HH:MM","title":"...","desc":"...","done":false},...],"budget":number,"duration":"Xh"}
Generate 4-6 steps. Use Chinese. Be creative and romantic.`

      const response = await sendAIRequest({
        apiKey: s.apiKey,
        apiUrl: settingsStore.getApiUrl(),
        model: s.model,
        messages: [{ role: 'user', content: prompt }],
        temperature: 0.9,
        maxTokens: 500,
        stream: false,
        timeout: s.timeout,
      })

      if (response.content) {
        const jsonMatch = response.content.match(/\{[\s\S]*\}/)
        if (jsonMatch) {
          const parsed = JSON.parse(jsonMatch[0]) as DatePlan
          if (parsed.steps?.length > 0) {
            aiPlan = parsed
            aiPlan.steps = aiPlan.steps.map(step => ({ ...step, done: false }))
          }
        }
      }
    } catch { /* fallback to template */ }
  }

  setTimeout(() => {
    plan.value = aiPlan || planTemplates[selectedType.value] || {
      title: `${selectedType.value}`,
      steps: [
        { time: '15:00', title: 'Start', desc: `Enjoy ${selectedType.value}`, done: false },
        { time: '18:00', title: 'Dinner', desc: 'Find a nice place', done: false },
        { time: '20:00', title: 'Wrap up', desc: 'Head home happy', done: false },
      ],
      budget: 300,
      duration: '5h',
    }
    generating.value = false
  }, aiPlan ? 300 : 1000)
}

const DATES_KEY = 'offline-date-history'

function saveDateHistory() {
  try {
    localStorage.setItem(DATES_KEY, JSON.stringify(pastDates.value))
  } catch { /* ignore */ }
}

function completePlan() {
  if (!plan.value || !selectedType.value) return
  const typeObj = dateTypes.find(t => t.label === selectedType.value)
  const now = new Date()
  const dateStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
  pastDates.value.unshift({
    id: Date.now(),
    icon: typeObj?.icon || '',
    title: plan.value.title,
    date: dateStr,
    rating: Math.floor(Math.random() * 2) + 4,
  })
  if (pastDates.value.length > 20) pastDates.value = pastDates.value.slice(0, 20)
  saveDateHistory()
  plan.value = null
}

onMounted(() => {
  try {
    const saved = localStorage.getItem(DATES_KEY)
    if (saved) pastDates.value = JSON.parse(saved)
  } catch { /* ignore */ }
})
</script>

<style scoped>
.offline-date-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.page-content {
  padding: 0 16px 16px;
}

.date-banner {
  position: relative;
  padding: 28px 16px;
  text-align: center;
  margin-bottom: 16px;
  border-radius: 0 0 20px 20px;
  overflow: hidden;
}

.banner-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #ff6b9d 0%, #ff8e53 100%);
  opacity: 0.15;
}

.banner-content {
  position: relative;
}

.date-emoji { font-size: 48px; }

.banner-content h2 {
  margin: 8px 0 4px;
  font-size: 20px;
  color: var(--text-primary);
}

.banner-content p {
  margin: 0;
  font-size: 13px;
  color: var(--text-tertiary);
}

.section {
  margin-bottom: 16px;
}

.section h3 {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 10px;
}

.type-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

.type-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 14px 4px;
  background: var(--bg-secondary);
  border-radius: 14px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}

.type-card.selected {
  border-color: var(--accent-pink);
  background: rgba(255, 107, 157, 0.08);
}

.type-icon { font-size: 26px; }
.type-label { font-size: 11px; color: var(--text-primary); font-weight: 500; }

.generate-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #ff6b9d, #ff8e53);
  color: #fff;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 16px;
}

.generate-btn:disabled { opacity: 0.5; }

.plan-card {
  background: var(--bg-secondary);
  border-radius: 16px;
  padding: 16px;
  margin-bottom: 16px;
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.plan-header h3 {
  margin: 0;
  font-size: 17px;
  color: var(--text-primary);
}

.refresh-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
}

.timeline {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.timeline-item {
  display: flex;
  gap: 12px;
  padding: 8px 0;
}

.timeline-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: var(--text-secondary);
  flex-shrink: 0;
  cursor: pointer;
}

.timeline-dot.done {
  background: var(--accent-green);
  color: #fff;
}

.timeline-content { flex: 1; }

.step-time {
  font-size: 11px;
  color: var(--accent-pink);
  font-weight: 600;
}

.step-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.step-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

.plan-footer {
  display: flex;
  gap: 20px;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.plan-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}

.past-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.past-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.past-icon { font-size: 24px; }

.past-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.past-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.past-date {
  font-size: 11px;
  color: var(--text-tertiary);
}

.past-rating {
  font-size: 12px;
}
</style>
