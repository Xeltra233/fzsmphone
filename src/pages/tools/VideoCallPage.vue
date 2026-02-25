<template>
  <div class="video-call-page">
    <!-- 通话前 -->
    <template v-if="!inCall">
      <NavBar title="视频通话" />
      <div class="pre-call">
        <div class="contact-avatar">📹</div>
        <h2 class="contact-name">{{ contactName }}</h2>
        <p class="call-hint">点击发起视频通话</p>
        <button class="call-start-btn" @click="startCall">
          <span>📹</span>
          <span>发起视频通话</span>
        </button>
      </div>
    </template>

    <!-- 视频通话中 -->
    <template v-else>
      <div class="video-screen">
        <!-- 对方画面（全屏） -->
        <div class="remote-video">
          <div class="video-placeholder">
            <div class="avatar-circle">{{ contactName.charAt(0) }}</div>
            <span v-if="!connected" class="connecting-text">{{ callStatus }}</span>
          </div>
        </div>

        <!-- 自己画面（小窗） -->
        <div class="local-video" :class="{ minimized: localMinimized }" @click="localMinimized = !localMinimized">
          <div class="local-placeholder">
            <span>👤</span>
          </div>
        </div>

        <!-- 顶部信息 -->
        <div class="top-bar">
          <div class="call-info-top">
            <span class="caller-name-top">{{ contactName }}</span>
            <span v-if="timer" class="timer-top">{{ timer }}</span>
            <span v-else class="status-top">{{ callStatus }}</span>
          </div>
        </div>

        <!-- 底部控制栏 -->
        <div class="bottom-controls">
          <button class="ctl-btn" :class="{ active: muted }" @click="muted = !muted">
            <span>{{ muted ? '🔇' : '🎤' }}</span>
          </button>
          <button class="ctl-btn" :class="{ active: !cameraOn }" @click="cameraOn = !cameraOn">
            <span>{{ cameraOn ? '📷' : '🚫' }}</span>
          </button>
          <button class="hangup-btn" @click="endCall">
            <span>📵</span>
          </button>
          <button class="ctl-btn" @click="flipCamera = !flipCamera">
            <span>🔄</span>
          </button>
          <button class="ctl-btn" :class="{ active: beautyOn }" @click="beautyOn = !beautyOn">
            <span>✨</span>
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
const callStatus = ref('正在连接...')
const timer = ref('')
const muted = ref(false)
const cameraOn = ref(true)
const flipCamera = ref(false)
const beautyOn = ref(false)
const localMinimized = ref(false)
let timerInterval: ReturnType<typeof setInterval> | null = null
let connectTimeout: ReturnType<typeof setTimeout> | null = null

function startCall() {
  inCall.value = true
  connected.value = false
  callStatus.value = '正在连接...'
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
  }, 3000)
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
.video-call-page {
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
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 28px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 16px;
  box-shadow: 0 4px 16px rgba(0, 122, 255, 0.3);
}

/* 视频通话界面 */
.video-screen {
  position: relative;
  height: 100%;
  background: #000;
}

.remote-video {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.avatar-circle {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  color: #fff;
  font-weight: 600;
}

.connecting-text {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

.local-video {
  position: absolute;
  top: 60px;
  right: 16px;
  width: 110px;
  height: 150px;
  border-radius: 14px;
  overflow: hidden;
  border: 2px solid rgba(255, 255, 255, 0.3);
  z-index: 10;
  cursor: pointer;
  transition: all 0.3s;
}

.local-video.minimized {
  width: 70px;
  height: 96px;
}

.local-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #2d3436, #636e72);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.top-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: linear-gradient(180deg, rgba(0,0,0,0.6) 0%, transparent 100%);
  z-index: 5;
}

.call-info-top {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.caller-name-top {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.timer-top {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  font-variant-numeric: tabular-nums;
}

.status-top {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.bottom-controls {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 24px 16px 36px;
  background: linear-gradient(0deg, rgba(0,0,0,0.7) 0%, transparent 100%);
  z-index: 5;
}

.ctl-btn {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  border: none;
  font-size: 22px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);
}

.ctl-btn.active {
  background: rgba(255, 255, 255, 0.35);
}

.hangup-btn {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #e74c3c;
  border: none;
  font-size: 26px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(231, 76, 60, 0.4);
}
</style>
