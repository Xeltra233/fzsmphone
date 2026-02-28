<template>
  <div class="video-call-page">
    <!-- 通话前 -->
    <template v-if="!inCall">
      <NavBar title="视频通话" />
      <div class="pre-call">
        <div class="contact-avatar-wrap">
          <img v-if="contactAvatar" :src="contactAvatar" class="contact-avatar-img" />
          <span v-else class="contact-avatar-emoji">▷</span>
        </div>
        <h2 class="contact-name">{{ contactName }}</h2>
        <p class="contact-number" v-if="contactNumber">{{ formatPhone(contactNumber) }}</p>
        <p class="call-hint">{{ characterId ? '点击发起 AI 视频通话' : '点击发起视频通话' }}</p>
        <button class="call-start-btn" @click="startCall">
          <span>▷</span>
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
            <div class="avatar-circle">
              <img v-if="contactAvatar" :src="contactAvatar" class="avatar-circle-img" />
              <span v-else>{{ contactName.charAt(0) }}</span>
            </div>
            <span v-if="!connected" class="connecting-text">{{ callStatus }}</span>
          </div>
        </div>

        <!-- 自己画面（小窗） -->
        <div class="local-video" :class="{ minimized: localMinimized }" @click="localMinimized = !localMinimized">
          <div class="local-placeholder">
            <span>○</span>
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

        <!-- AI 对话覆盖层 -->
        <div v-if="connected && characterId" class="chat-overlay">
          <div class="chat-scroll" ref="chatRef">
            <div
              v-for="msg in callMessages"
              :key="msg.id"
              class="video-msg"
              :class="{ 'is-me': msg.from === 'me' }"
            >
              <div class="video-msg-bubble" :class="msg.from === 'me' ? 'vmsg-me' : 'vmsg-other'">
                {{ msg.text }}
              </div>
            </div>
            <!-- AI 正在说话 -->
            <div v-if="isAITyping" class="video-msg">
              <div class="video-msg-bubble vmsg-other">
                <div v-if="streamingText">{{ streamingText }}</div>
                <div v-else class="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
              </div>
            </div>
          </div>

          <div class="chat-input-row">
            <input
              v-model="inputText"
              type="text"
              placeholder="说点什么..."
              class="chat-input"
              @keyup.enter="sendCallMessage"
            />
            <button class="chat-send" :disabled="!inputText.trim() || isAITyping" @click="sendCallMessage">
              发送
            </button>
          </div>
        </div>

        <!-- 底部控制栏 -->
        <div class="bottom-controls">
          <button class="ctl-btn" :class="{ active: muted }" @click="muted = !muted">
            <span>{{ muted ? '◇' : '◉' }}</span>
          </button>
          <button class="ctl-btn" :class="{ active: !cameraOn }" @click="cameraOn = !cameraOn">
            <span>{{ cameraOn ? '▣' : '✕' }}</span>
          </button>
          <button class="hangup-btn" @click="endCall">
            <span>✕</span>
          </button>
          <button class="ctl-btn" @click="flipCamera = !flipCamera">
            <span>↻</span>
          </button>
          <button class="ctl-btn" :class="{ active: beautyOn }" @click="beautyOn = !beautyOn">
            <span>✦</span>
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { usePhoneStore } from '@/stores/phone'
import { useSettingsStore } from '@/stores/settings'
import {
  sendAIRequest,
  buildCallMessages,
  getCharacterById,
} from '@/utils/aiService'

interface CallMsg {
  id: string
  from: 'me' | 'other'
  text: string
}

const route = useRoute()
const router = useRouter()
const phoneStore = usePhoneStore()
const settingsStore = useSettingsStore()

const contactName = ref('未知')
const contactNumber = ref('')
const contactAvatar = ref('')
const characterId = ref('')

const inCall = ref(false)
const connected = ref(false)
const callStatus = ref('正在连接...')
const timer = ref('')
const muted = ref(false)
const cameraOn = ref(true)
const flipCamera = ref(false)
const beautyOn = ref(false)
const localMinimized = ref(false)

const callMessages = ref<CallMsg[]>([])
const inputText = ref('')
const isAITyping = ref(false)
const streamingText = ref('')
const chatRef = ref<HTMLElement | null>(null)

let timerInterval: ReturnType<typeof setInterval> | null = null
let connectTimeout: ReturnType<typeof setTimeout> | null = null
let callSeconds = 0
let abortController: AbortController | null = null

function formatPhone(num: string): string {
  if (num.length <= 3) return num
  if (num.length <= 7) return `${num.slice(0, 3)} ${num.slice(3)}`
  return `${num.slice(0, 3)} ${num.slice(3, 7)} ${num.slice(7)}`
}

function startCall() {
  inCall.value = true
  connected.value = false
  callStatus.value = '正在连接...'
  timer.value = ''
  callMessages.value = []
  callSeconds = 0

  connectTimeout = setTimeout(async () => {
    connected.value = true
    callStatus.value = '通话中'
    callSeconds = 0
    timerInterval = setInterval(() => {
      callSeconds++
      const m = Math.floor(callSeconds / 60)
      const s = callSeconds % 60
      timer.value = `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
    }, 1000)

    // 如果有角色，AI 先打招呼
    if (characterId.value) {
      await triggerAIReply()
    }
  }, 3000)
}

function endCall() {
  inCall.value = false
  connected.value = false
  if (timerInterval) { clearInterval(timerInterval); timerInterval = null }
  if (connectTimeout) { clearTimeout(connectTimeout); connectTimeout = null }
  if (abortController) { abortController.abort(); abortController = null }

  // 保存通话记录
  const duration = callSeconds > 0
    ? `${Math.floor(callSeconds / 60)}:${(callSeconds % 60).toString().padStart(2, '0')}`
    : '0:00'

  const now = new Date()
  const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`

  phoneStore.addCallRecord({
    name: contactName.value,
    number: contactNumber.value,
    characterId: characterId.value,
    avatar: contactAvatar.value,
    type: 'outgoing',
    callType: 'video',
    time: timeStr,
    duration,
  })

  router.back()
}

async function sendCallMessage() {
  if (!inputText.value.trim() || isAITyping.value) return

  const text = inputText.value.trim()
  inputText.value = ''

  callMessages.value.push({
    id: `me-${Date.now()}`,
    from: 'me',
    text,
  })
  scrollChat()

  await triggerAIReply()
}

async function triggerAIReply() {
  if (!characterId.value) return

  const s = settingsStore.settings
  if (!s.apiKey) {
    callMessages.value.push({
      id: `sys-${Date.now()}`,
      from: 'other',
      text: '△ 未配置 API Key，无法进行 AI 通话',
    })
    scrollChat()
    return
  }

  const character = getCharacterById(characterId.value)
  if (!character) return

  isAITyping.value = true
  streamingText.value = ''
  abortController = new AbortController()

  await new Promise(resolve => setTimeout(resolve, 500 + Math.random() * 1000))

  try {
    const history = callMessages.value.map(m => ({ from: m.from, text: m.text }))
    const aiMessages = buildCallMessages(character, history, 'video')

    const apiUrl = settingsStore.getApiUrl()
    const isStream = s.streamEnabled

    const response = await sendAIRequest({
      apiKey: s.apiKey,
      apiUrl: apiUrl,
      model: s.model,
      messages: aiMessages,
      temperature: s.temperature,
      maxTokens: 150,
      stream: isStream,
      timeout: s.timeout,
      signal: abortController.signal,
      onChunk: (chunk: string) => {
        streamingText.value += chunk
        scrollChat()
      },
    })

    isAITyping.value = false
    const content = isStream ? streamingText.value.trim() : response.content
    streamingText.value = ''

    if (content) {
      callMessages.value.push({
        id: `ai-${Date.now()}`,
        from: 'other',
        text: content,
      })
      scrollChat()
    }
  } catch (err: any) {
    isAITyping.value = false
    streamingText.value = ''
    if (err.name !== 'AbortError') {
      callMessages.value.push({
        id: `err-${Date.now()}`,
        from: 'other',
        text: `✕ ${err.message}`,
      })
      scrollChat()
    }
  } finally {
    abortController = null
  }
}

function scrollChat() {
  nextTick(() => {
    if (chatRef.value) {
      chatRef.value.scrollTop = chatRef.value.scrollHeight
    }
  })
}

onMounted(() => {
  const q = route.query
  if (q.name) contactName.value = String(q.name)
  if (q.number) contactNumber.value = String(q.number)
  if (q.characterId) {
    characterId.value = String(q.characterId)
    const char = getCharacterById(characterId.value)
    if (char) {
      contactName.value = char.name || contactName.value
      contactAvatar.value = char.avatar || ''
    }
  }
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
  if (connectTimeout) clearTimeout(connectTimeout)
  if (abortController) abortController.abort()
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
  gap: 12px;
  padding: 32px;
}

.contact-avatar-wrap {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.contact-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.contact-avatar-emoji {
  font-size: 72px;
}

.contact-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.contact-number {
  font-size: 14px;
  color: var(--text-tertiary);
  margin: 0;
  letter-spacing: 1px;
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
  overflow: hidden;
}

.avatar-circle-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

/* AI 对话覆盖层 */
.chat-overlay {
  position: absolute;
  bottom: 90px;
  left: 0;
  right: 0;
  max-height: 45%;
  display: flex;
  flex-direction: column;
  z-index: 8;
  pointer-events: auto;
}

.chat-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 8px 16px;
  min-height: 0;
}

.video-msg {
  display: flex;
  flex-direction: column;
  margin-bottom: 6px;
}

.video-msg.is-me {
  align-items: flex-end;
}

.video-msg-bubble {
  max-width: 80%;
  padding: 7px 12px;
  border-radius: 14px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
  backdrop-filter: blur(12px);
}

.vmsg-me {
  background: rgba(0, 122, 255, 0.75);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.vmsg-other {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border-bottom-left-radius: 4px;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  align-items: center;
}

.typing-indicator span {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.6);
  animation: typingBounce 1.2s infinite;
}

.typing-indicator span:nth-child(2) { animation-delay: 0.15s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.3s; }

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-3px); opacity: 1; }
}

.chat-input-row {
  display: flex;
  gap: 8px;
  padding: 6px 16px;
}

.chat-input {
  flex: 1;
  padding: 8px 14px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  background: rgba(0, 0, 0, 0.4);
  color: #fff;
  font-size: 13px;
  outline: none;
  backdrop-filter: blur(8px);
}

.chat-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.chat-input:focus {
  border-color: rgba(255, 255, 255, 0.4);
}

.chat-send {
  padding: 6px 14px;
  background: rgba(0, 122, 255, 0.7);
  color: #fff;
  border: none;
  border-radius: 18px;
  font-size: 13px;
  cursor: pointer;
  backdrop-filter: blur(8px);
}

.chat-send:disabled {
  opacity: 0.4;
}

/* 底部控制栏 */
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
  z-index: 15;
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
