<template>
  <div class="voice-call-page">
    <!-- 通话前 -->
    <template v-if="!inCall">
      <NavBar title="语音通话" />
      <div class="pre-call">
        <div class="contact-avatar">🎧</div>
        <h2 class="contact-name">{{ contactName }}</h2>
        <p class="call-hint">点击拨打发起语音通话</p>
        <button class="call-start-btn" @click="startCall">
          <span>📞</span>
          <span>发起通话</span>
        </button>
      </div>
    </template>

    <!-- 通话中 -->
    <template v-else>
      <div class="calling-screen">
        <div class="call-bg"></div>
        <div class="call-content">
          <div class="caller-avatar-large">🎧</div>
          <h2 class="caller-name">{{ contactName }}</h2>
          <div class="call-status">{{ callStatus }}</div>
          <div v-if="timer" class="call-timer">{{ timer }}</div>

          <!-- 音波动画 -->
          <div v-if="connected" class="sound-waves">
            <div class="wave" v-for="i in 5" :key="i" :style="{ animationDelay: i * 0.15 + 's' }"></div>
          </div>

          <div class="call-controls">
            <button class="ctrl-btn" :class="{ active: muted }" @click="muted = !muted">
              <span>{{ muted ? '🔇' : '🎤' }}</span>
              <span>{{ muted ? '已静音' : '静音' }}</span>
            </button>
            <button class="ctrl-btn" :class="{ active: speakerOn }" @click="speakerOn = !speakerOn">
              <span>{{ speakerOn ? '📢' : '🔈' }}</span>
              <span>扬声器</span>
            </button>
            <button class="ctrl-btn" @click="showKeypad = !showKeypad">
              <span>⌨️</span>
              <span>键盘</span>
            </button>
          </div>

          <button class="hangup-btn" @click="endCall">
            <span>📵</span>
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

const contactName = ref('宝贝')
const inCall = ref(false)
const connected = ref(false)
const callStatus = ref('正在呼叫...')
const timer = ref('')
const muted = ref(false)
const speakerOn = ref(false)
const showKeypad = ref(false)
let timerInterval: ReturnType<typeof setInterval> | null = null
let connectTimeout: ReturnType<typeof setTimeout> | null = null

function startCall() {
  inCall.value = true
  connected.value = false
  callStatus.value = '正在呼叫...'
  timer.value = ''

  connectTimeout = setTimeout(() => {
    connected.value = true
    callStatus.value = '通话中'
    let seconds = 0
    timerInterval = setInterval(() => {
      seconds++
      const m = Math.floor(seconds / 60)
      const s = seconds % 60
      timer.value = `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
    }, 1000)
  }, 2500)
}

function endCall() {
  inCall.value = false
  connected.value = false
  if (timerInterval) { clearInterval(timerInterval); timerInterval = null }
  if (connectTimeout) { clearTimeout(connectTimeout); connectTimeout = null }
}

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
  if (connectTimeout) clearTimeout(connectTimeout)
})
</script>

<style scoped>
.voice-call-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.pre-call {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 32px;
}

.contact-avatar {
  font-size: 72px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.contact-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.call-hint {
  font-size: 14px;
  color: var(--text-tertiary);
  margin: 0;
}

.call-start-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  background: var(--accent-green);
  color: #fff;
  border: none;
  border-radius: 28px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 16px;
  box-shadow: 0 4px 16px rgba(52, 199, 89, 0.3);
}

.calling-screen {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.call-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}

.call-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #fff;
  width: 100%;
  padding: 0 32px;
}

.caller-avatar-large {
  font-size: 64px;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.caller-name {
  font-size: 28px;
  font-weight: 300;
  margin: 0;
}

.call-status {
  font-size: 14px;
  opacity: 0.6;
}

.call-timer {
  font-size: 20px;
  font-weight: 300;
  font-variant-numeric: tabular-nums;
  opacity: 0.8;
}

.sound-waves {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 40px;
  margin: 12px 0;
}

.wave {
  width: 4px;
  background: rgba(52, 199, 89, 0.8);
  border-radius: 2px;
  animation: waveAnim 1s ease-in-out infinite;
}

@keyframes waveAnim {
  0%, 100% { height: 8px; }
  50% { height: 32px; }
}

.call-controls {
  display: flex;
  gap: 20px;
  margin: 24px 0;
}

.ctrl-btn {
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
  min-width: 70px;
}

.ctrl-btn.active {
  background: rgba(255, 255, 255, 0.25);
}

.ctrl-btn span:first-child { font-size: 24px; }
.ctrl-btn span:last-child { font-size: 11px; opacity: 0.7; }

.hangup-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #e74c3c;
  border: none;
  font-size: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16px;
  box-shadow: 0 4px 16px rgba(231, 76, 60, 0.4);
}
</style>
