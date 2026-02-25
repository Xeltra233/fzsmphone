<template>
  <div class="phone-page">
    <NavBar title="电话" />

    <!-- 号码显示 -->
    <div class="number-display">
      <span class="phone-number">{{ displayNumber || '输入号码' }}</span>
      <button v-if="phoneNumber" class="backspace-btn" @click="backspace">⌫</button>
    </div>

    <!-- 拨号盘 -->
    <div class="dial-pad">
      <div
        v-for="key in dialKeys"
        :key="key.num"
        class="dial-key"
        @click="pressKey(key.num)"
      >
        <span class="key-num">{{ key.num }}</span>
        <span class="key-letters">{{ key.letters }}</span>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="action-row">
      <button class="action-btn" @click="addContact">
        <span>➕</span>
        <span class="action-label">添加联系人</span>
      </button>
      <button class="call-btn" @click="makeCall">
        <span>📞</span>
      </button>
      <button class="action-btn" @click="showRecent = true">
        <span>🕐</span>
        <span class="action-label">通话记录</span>
      </button>
    </div>

    <!-- 通话记录面板 -->
    <div v-if="showRecent" class="overlay" @click.self="showRecent = false">
      <div class="recent-panel">
        <div class="panel-header">
          <h3>通话记录</h3>
          <button class="close-btn" @click="showRecent = false">✕</button>
        </div>
        <div v-if="callHistory.length === 0" class="empty-history">
          <span>📱</span>
          <p>暂无通话记录</p>
        </div>
        <div v-else class="history-list">
          <div
            v-for="call in callHistory"
            :key="call.id"
            class="history-item"
            @click="phoneNumber = call.number; showRecent = false"
          >
            <div class="call-icon" :class="call.type">
              {{ call.type === 'incoming' ? '📥' : call.type === 'outgoing' ? '📤' : '📵' }}
            </div>
            <div class="call-info">
              <span class="call-name">{{ call.name || call.number }}</span>
              <span class="call-time">{{ call.time }}</span>
            </div>
            <span class="call-duration">{{ call.duration }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 通话中界面 -->
    <div v-if="calling" class="calling-overlay">
      <div class="calling-screen">
        <div class="calling-avatar">📞</div>
        <div class="calling-number">{{ displayNumber }}</div>
        <div class="calling-status">{{ callStatus }}</div>
        <div class="calling-timer" v-if="callTimer">{{ callTimer }}</div>
        <div class="calling-actions">
          <button class="calling-action-btn" @click="toggleMute">
            <span>{{ muted ? '🔇' : '🔊' }}</span>
            <span>{{ muted ? '取消静音' : '静音' }}</span>
          </button>
          <button class="calling-action-btn" @click="toggleSpeaker">
            <span>{{ speaker ? '📢' : '🔈' }}</span>
            <span>{{ speaker ? '关闭扬声器' : '扬声器' }}</span>
          </button>
          <button class="calling-action-btn" @click="toggleKeypad">
            <span>⌨️</span>
            <span>键盘</span>
          </button>
        </div>
        <button class="hangup-btn" @click="hangUp">
          <span>📵</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

const phoneNumber = ref('')
const showRecent = ref(false)
const calling = ref(false)
const callStatus = ref('正在拨号...')
const callTimer = ref('')
const muted = ref(false)
const speaker = ref(false)
let timerInterval: ReturnType<typeof setInterval> | null = null

const dialKeys = [
  { num: '1', letters: '' },
  { num: '2', letters: 'ABC' },
  { num: '3', letters: 'DEF' },
  { num: '4', letters: 'GHI' },
  { num: '5', letters: 'JKL' },
  { num: '6', letters: 'MNO' },
  { num: '7', letters: 'PQRS' },
  { num: '8', letters: 'TUV' },
  { num: '9', letters: 'WXYZ' },
  { num: '*', letters: '' },
  { num: '0', letters: '+' },
  { num: '#', letters: '' },
]

interface CallRecord {
  id: number
  name: string
  number: string
  type: 'incoming' | 'outgoing' | 'missed'
  time: string
  duration: string
}

const callHistory = ref<CallRecord[]>([
  { id: 1, name: '宝贝', number: '13800138000', type: 'outgoing', time: '今天 14:30', duration: '5:23' },
  { id: 2, name: '', number: '10086', type: 'incoming', time: '今天 11:15', duration: '2:01' },
  { id: 3, name: '外卖小哥', number: '15912345678', type: 'missed', time: '昨天 18:45', duration: '' },
  { id: 4, name: '妈妈', number: '13912345678', type: 'incoming', time: '昨天 12:00', duration: '15:42' },
])

const displayNumber = computed(() => {
  const num = phoneNumber.value
  if (num.length <= 3) return num
  if (num.length <= 7) return `${num.slice(0, 3)} ${num.slice(3)}`
  return `${num.slice(0, 3)} ${num.slice(3, 7)} ${num.slice(7)}`
})

function pressKey(key: string) {
  if (phoneNumber.value.length < 15) {
    phoneNumber.value += key
  }
}

function backspace() {
  phoneNumber.value = phoneNumber.value.slice(0, -1)
}

function makeCall() {
  if (!phoneNumber.value) return
  calling.value = true
  callStatus.value = '正在拨号...'
  callTimer.value = ''
  muted.value = false
  speaker.value = false

  // 模拟接通
  setTimeout(() => {
    callStatus.value = '通话中'
    let seconds = 0
    timerInterval = setInterval(() => {
      seconds++
      const m = Math.floor(seconds / 60)
      const s = seconds % 60
      callTimer.value = `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
    }, 1000)
  }, 2000)
}

function hangUp() {
  calling.value = false
  if (timerInterval) {
    clearInterval(timerInterval)
    timerInterval = null
  }
  // 添加到通话记录
  callHistory.value.unshift({
    id: Date.now(),
    name: '',
    number: phoneNumber.value,
    type: 'outgoing',
    time: '刚刚',
    duration: callTimer.value || '0:00',
  })
}

function toggleMute() { muted.value = !muted.value }
function toggleSpeaker() { speaker.value = !speaker.value }
function toggleKeypad() { /* 可扩展 */ }

function addContact() {
  if (!phoneNumber.value) return
  console.log('添加联系人:', phoneNumber.value)
}
</script>

<style scoped>
.phone-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.number-display {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px 12px;
  gap: 8px;
  min-height: 60px;
}

.phone-number {
  font-size: 28px;
  font-weight: 300;
  color: var(--text-primary);
  letter-spacing: 2px;
}

.backspace-btn {
  background: none;
  border: none;
  font-size: 22px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px;
}

.dial-pad {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  padding: 12px 40px;
  flex: 1;
  align-content: center;
}

.dial-key {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--bg-secondary);
  cursor: pointer;
  margin: 0 auto;
  transition: background 0.15s;
  user-select: none;
}

.dial-key:active {
  background: var(--bg-tertiary);
}

.key-num {
  font-size: 28px;
  font-weight: 400;
  color: var(--text-primary);
  line-height: 1;
}

.key-letters {
  font-size: 10px;
  color: var(--text-tertiary);
  letter-spacing: 2px;
  margin-top: 2px;
  min-height: 14px;
}

.action-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32px;
  padding: 16px 16px 32px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  cursor: pointer;
}

.action-btn span:first-child {
  font-size: 20px;
}

.action-label {
  font-size: 11px;
  color: var(--text-tertiary);
}

.call-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--accent-green);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(52, 199, 89, 0.3);
  transition: transform 0.15s;
}

.call-btn:active {
  transform: scale(0.92);
}

/* 通话记录 */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.recent-panel {
  width: 100%;
  max-width: 393px;
  max-height: 70vh;
  background: var(--bg-secondary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: var(--text-tertiary);
  cursor: pointer;
}

.empty-history {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px;
  color: var(--text-tertiary);
}

.empty-history span {
  font-size: 40px;
  margin-bottom: 8px;
}

.history-list {
  flex: 1;
  overflow-y: auto;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}

.history-item:active {
  background: var(--bg-tertiary);
}

.call-icon {
  font-size: 20px;
}

.call-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.call-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
}

.call-info .missed .call-name {
  color: var(--accent-red);
}

.call-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.call-duration {
  font-size: 13px;
  color: var(--text-secondary);
}

/* 通话中 */
.calling-overlay {
  position: fixed;
  inset: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.calling-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #fff;
  width: 100%;
  padding: 0 32px;
}

.calling-avatar {
  font-size: 64px;
  margin-bottom: 8px;
}

.calling-number {
  font-size: 28px;
  font-weight: 300;
  letter-spacing: 2px;
}

.calling-status {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.calling-timer {
  font-size: 18px;
  font-weight: 300;
  color: rgba(255, 255, 255, 0.8);
  font-variant-numeric: tabular-nums;
}

.calling-actions {
  display: flex;
  gap: 24px;
  margin: 32px 0;
}

.calling-action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 16px;
  padding: 14px 18px;
  cursor: pointer;
  color: #fff;
}

.calling-action-btn span:first-child {
  font-size: 24px;
}

.calling-action-btn span:last-child {
  font-size: 11px;
  opacity: 0.7;
}

.hangup-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #e74c3c;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  cursor: pointer;
  margin-top: 16px;
  box-shadow: 0 4px 16px rgba(231, 76, 60, 0.4);
}
</style>
