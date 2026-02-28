<template>
  <div class="reverse-peek-page">
    <NavBar title="反向偷看" />

    <div class="page-content">
      <!-- 状态卡 -->
      <div class="status-card">
        <div class="status-icon" :class="{ alert: caught }">
          {{ caught ? '◈' : '▢' }}
        </div>
        <h2>{{ caught ? 'TA在偷看你的手机！' : '设置手机陷阱' }}</h2>
        <p>{{ caught ? '已记录偷看行为' : '当TA偷看你手机时自动触发' }}</p>
      </div>

      <!-- 陷阱设置 -->
      <div class="trap-settings">
        <h3>陷阱模式</h3>
        <div class="trap-list">
          <div
            v-for="trap in traps"
            :key="trap.id"
            class="trap-item"
            :class="{ active: trap.enabled }"
            @click="trap.enabled = !trap.enabled"
          >
            <span class="trap-icon">{{ trap.icon }}</span>
            <div class="trap-info">
              <span class="trap-name">{{ trap.name }}</span>
              <span class="trap-desc">{{ trap.desc }}</span>
            </div>
            <div class="trap-toggle" :class="{ on: trap.enabled }">
              <div class="toggle-dot"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 模拟触发 -->
      <div class="simulate-section">
        <button class="simulate-btn" @click="simulateCatch" :disabled="caught">
          {{ caught ? '已触发' : '◎ 模拟偷看触发' }}
        </button>
        <button v-if="caught" class="reset-btn" @click="resetTrap">
          重置陷阱
        </button>
      </div>

      <!-- 触发记录 -->
      <div v-if="catchLogs.length > 0" class="logs-section">
        <h3>偷看记录</h3>
        <div class="log-list">
          <div v-for="log in catchLogs" :key="log.id" class="log-item">
            <span class="log-icon">{{ log.icon }}</span>
            <div class="log-info">
              <span class="log-title">{{ log.title }}</span>
              <span class="log-time">{{ log.time }}</span>
            </div>
            <span class="log-detail">{{ log.detail }}</span>
          </div>
        </div>
      </div>

      <!-- 假消息设置 -->
      <div class="fake-section">
        <h3>假消息诱饵</h3>
        <p class="fake-hint">设置假的通知和聊天记录，当TA偷看时看到的内容</p>
        <div class="fake-list">
          <div v-for="fake in fakeMessages" :key="fake.id" class="fake-item">
            <span class="fake-icon">{{ fake.icon }}</span>
            <div class="fake-info">
              <span class="fake-name">{{ fake.from }}</span>
              <span class="fake-msg">{{ fake.content }}</span>
            </div>
            <button class="edit-fake" @click="editFake(fake)">编辑</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

const caught = ref(false)

interface Trap {
  id: number
  icon: string
  name: string
  desc: string
  enabled: boolean
}

const traps = ref<Trap[]>([
  { id: 1, icon: '▣', name: '自动拍照', desc: '偷看时自动拍一张照片', enabled: true },
  { id: 2, icon: '▲', name: '震动提醒', desc: '你的手机会收到提醒', enabled: true },
  { id: 3, icon: '▤', name: '记录行为', desc: '记录偷看了哪些内容', enabled: true },
  { id: 4, icon: '◈', name: '显示假消息', desc: '显示预设的假消息', enabled: false },
  { id: 5, icon: '■', name: '自动锁屏', desc: '3秒后自动锁定屏幕', enabled: false },
])

interface CatchLog {
  id: number
  icon: string
  title: string
  time: string
  detail: string
}

const catchLogs = ref<CatchLog[]>([])

interface FakeMessage {
  id: number
  icon: string
  from: string
  content: string
}

const fakeMessages = ref<FakeMessage[]>([
  { id: 1, icon: '◌', from: '老板', content: '这个月奖金翻倍！' },
  { id: 2, icon: '▢', from: '快递', content: '您有一个惊喜包裹到了' },
  { id: 3, icon: '✉', from: '银行', content: '您的账户入账 ¥50,000' },
])

function simulateCatch() {
  caught.value = true
  const enabledTraps = traps.value.filter(t => t.enabled)
  enabledTraps.forEach(trap => {
    catchLogs.value.unshift({
      id: Date.now() + trap.id,
      icon: trap.icon,
      title: `${trap.name} 已触发`,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
      detail: trap.desc,
    })
  })
  catchLogs.value.unshift({
    id: Date.now(),
    icon: '△',
    title: '检测到偷看行为',
    time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    detail: '有人正在查看你的手机',
  })
}

function resetTrap() {
  caught.value = false
}

function editFake(fake: FakeMessage) {
  const newContent = prompt('编辑假消息内容:', fake.content)
  if (newContent !== null) {
    fake.content = newContent
    saveFakeMessages()
  }
}

const TRAP_KEY = 'reverse-peek-traps'
const LOG_KEY = 'reverse-peek-logs'
const FAKE_KEY = 'reverse-peek-fakes'

function saveTraps() {
  try {
    localStorage.setItem(TRAP_KEY, JSON.stringify(traps.value.map(t => ({ id: t.id, enabled: t.enabled }))))
  } catch { /* ignore */ }
}

function saveLogs() {
  try {
    localStorage.setItem(LOG_KEY, JSON.stringify(catchLogs.value))
  } catch { /* ignore */ }
}

function saveFakeMessages() {
  try {
    localStorage.setItem(FAKE_KEY, JSON.stringify(fakeMessages.value))
  } catch { /* ignore */ }
}

watch(traps, saveTraps, { deep: true })
watch(catchLogs, saveLogs, { deep: true })

onMounted(() => {
  try {
    const savedTraps = localStorage.getItem(TRAP_KEY)
    if (savedTraps) {
      const settings = JSON.parse(savedTraps) as { id: number; enabled: boolean }[]
      settings.forEach(s => {
        const trap = traps.value.find(t => t.id === s.id)
        if (trap) trap.enabled = s.enabled
      })
    }
    const savedLogs = localStorage.getItem(LOG_KEY)
    if (savedLogs) catchLogs.value = JSON.parse(savedLogs)
    const savedFakes = localStorage.getItem(FAKE_KEY)
    if (savedFakes) fakeMessages.value = JSON.parse(savedFakes)
  } catch { /* ignore */ }
})
</script>

<style scoped>
.reverse-peek-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.page-content {
  padding: 16px;
}

.status-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  background: var(--bg-secondary);
  border-radius: 20px;
  margin-bottom: 20px;
}

.status-icon {
  font-size: 52px;
  transition: transform 0.3s;
}

.status-icon.alert {
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-8px); }
  75% { transform: translateX(8px); }
}

.status-card h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.status-card p {
  margin: 0;
  font-size: 13px;
  color: var(--text-tertiary);
}

.trap-settings h3,
.logs-section h3,
.fake-section h3 {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 10px;
}

.trap-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
}

.trap-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: var(--bg-secondary);
  border-radius: 14px;
  cursor: pointer;
}

.trap-icon { font-size: 24px; }

.trap-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.trap-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.trap-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

.trap-toggle {
  width: 44px;
  height: 26px;
  border-radius: 13px;
  background: var(--bg-tertiary);
  padding: 2px;
  transition: background 0.2s;
}

.trap-toggle.on {
  background: var(--accent-green);
}

.toggle-dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.2s;
}

.trap-toggle.on .toggle-dot {
  transform: translateX(18px);
}

.simulate-section {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.simulate-btn {
  flex: 1;
  padding: 14px;
  background: var(--accent-red);
  color: #fff;
  border: none;
  border-radius: 14px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
}

.simulate-btn:disabled {
  opacity: 0.5;
}

.reset-btn {
  padding: 14px 20px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: none;
  border-radius: 14px;
  font-size: 14px;
  cursor: pointer;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 20px;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.log-icon { font-size: 20px; }

.log-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.log-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.log-time {
  font-size: 11px;
  color: var(--text-tertiary);
}

.log-detail {
  font-size: 11px;
  color: var(--text-secondary);
  max-width: 100px;
  text-align: right;
}

.fake-hint {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0 0 10px;
}

.fake-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.fake-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.fake-icon { font-size: 20px; }

.fake-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.fake-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.fake-msg {
  font-size: 12px;
  color: var(--text-tertiary);
}

.edit-fake {
  padding: 4px 12px;
  background: var(--bg-tertiary);
  border: none;
  border-radius: 8px;
  font-size: 12px;
  color: var(--accent-blue);
  cursor: pointer;
}
</style>
