<template>
  <div class="voice-call-page">
    <!-- 通话前 -->
    <template v-if="!inCall">
      <NavBar title="语音通话" />
      <div class="pre-call">
        <div class="contact-avatar-wrap">
          <img v-if="contactAvatar" :src="contactAvatar" class="contact-avatar-img" />
          <span v-else class="contact-avatar-emoji">♪</span>
        </div>
        <h2 class="contact-name">{{ contactName }}</h2>
        <p class="contact-number" v-if="contactNumber">{{ formatPhone(contactNumber) }}</p>
        <p class="call-hint">{{ characterId ? '点击拨打发起 AI 语音通话' : '点击拨打发起语音通话' }}</p>
        <button class="call-start-btn" @click="startCall">
          <span>☎</span>
          <span>发起通话</span>
        </button>
      </div>
    </template>

    <!-- 通话中 -->
    <template v-else>
      <div class="calling-screen">
        <div class="call-bg"></div>
        <div class="call-content">
          <!-- 顶部信息 -->
          <div class="caller-avatar-large">
            <img v-if="contactAvatar" :src="contactAvatar" class="caller-avatar-img" />
            <span v-else>♪</span>
          </div>
          <h2 class="caller-name">{{ contactName }}</h2>
          <div class="call-status">{{ callStatus }}</div>
          <div v-if="timer" class="call-timer">{{ timer }}</div>

          <!-- 音波动画 -->
          <div v-if="connected" class="sound-waves">
            <div class="wave" v-for="i in 5" :key="i"
              :class="{ active: isAITyping }"
              :style="{ animationDelay: i * 0.15 + 's' }">
            </div>
          </div>

          <!-- AI 对话区域 -->
          <div v-if="connected && characterId" class="call-chat-area" ref="chatRef">
            <div
              v-for="msg in callMessages"
              :key="msg.id"
              class="call-msg"
              :class="{ 'is-me': msg.from === 'me' }"
            >
              <div class="call-msg-bubble" :class="msg.from === 'me' ? 'msg-me' : 'msg-other'">
                {{ msg.text }}
              </div>
            </div>
            <!-- AI 正在说话 -->
            <div v-if="isAITyping" class="call-msg">
              <div class="call-msg-bubble msg-other">
                <div v-if="streamingText">{{ streamingText }}</div>
                <div v-else class="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
              </div>
            </div>
          </div>

          <!-- 通话中输入框（AI模式） -->
          <div v-if="connected && characterId" class="call-input-area">
            <input
              v-model="inputText"
              type="text"
              :placeholder="'说点什么...'"
              class="call-input"
              @keyup.enter="sendCallMessage"
            />
            <button class="call-send-btn" :disabled="!inputText.trim() || isAITyping" @click="sendCallMessage">
              发送
            </button>
          </div>

          <!-- 控制按钮 -->
          <div class="call-controls">
            <button class="ctrl-btn" :class="{ active: muted }" @click="muted = !muted">
              <span>{{ muted ? '◇' : '◉' }}</span>
              <span>{{ muted ? '已静音' : '静音' }}</span>
            </button>
            <button class="ctrl-btn" :class="{ active: speakerOn }" @click="speakerOn = !speakerOn">
              <span>{{ speakerOn ? '◉' : '◎' }}</span>
              <span>扬声器</span>
            </button>
            <button class="ctrl-btn" @click="showKeypad = !showKeypad">
              <span>⌨️</span>
              <span>键盘</span>
            </button>
          </div>

          <button class="hangup-btn" @click="endCall">
            <span>✕</span>
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
const callStatus = ref('正在呼叫...')
const timer = ref('')
const muted = ref(false)
const speakerOn = ref(false)
const showKeypad = ref(false)

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
  callStatus.value = '正在呼叫...'
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
  }, 2500)
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
    callType: 'voice',
    time: timeStr,
    duration,
  })

  // 返回电话页面
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

  // 模拟短暂延迟
  await new Promise(resolve => setTimeout(resolve, 500 + Math.random() * 1000))

  try {
    const history = callMessages.value.map(m => ({ from: m.from, text: m.text }))
    const aiMessages = buildCallMessages(character, history, 'voice')

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
  // 从路由参数获取联系人信息
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
  background: var(--accent-green);
  color: #fff;
  border: none;
  border-radius: 28px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 16px;
  box-shadow: 0 4px 16px rgba(52, 199, 89, 0.3);
}

/* 通话中 */
.calling-screen {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
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
  color: #fff;
  width: 100%;
  height: 100%;
  padding: 24px 16px 16px;
  box-sizing: border-box;
}

.caller-avatar-large {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  overflow: hidden;
  flex-shrink: 0;
}

.caller-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.caller-name {
  font-size: 22px;
  font-weight: 300;
  margin: 8px 0 2px;
}

.call-status {
  font-size: 13px;
  opacity: 0.6;
}

.call-timer {
  font-size: 16px;
  font-weight: 300;
  font-variant-numeric: tabular-nums;
  opacity: 0.8;
  margin-top: 2px;
}

/* 音波 */
.sound-waves {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 30px;
  margin: 8px 0;
  flex-shrink: 0;
}

.wave {
  width: 4px;
  background: rgba(52, 199, 89, 0.5);
  border-radius: 2px;
  animation: waveAnim 1s ease-in-out infinite;
  transition: background 0.3s;
}

.wave.active {
  background: rgba(52, 199, 89, 1);
}

@keyframes waveAnim {
  0%, 100% { height: 6px; }
  50% { height: 24px; }
}

/* 通话中对话区域 */
.call-chat-area {
  flex: 1;
  width: 100%;
  overflow-y: auto;
  padding: 8px 0;
  min-height: 0;
}

.call-msg {
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
}

.call-msg.is-me {
  align-items: flex-end;
}

.call-msg-bubble {
  max-width: 85%;
  padding: 8px 14px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.msg-me {
  background: rgba(0, 122, 255, 0.8);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.msg-other {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  border-bottom-left-radius: 4px;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  align-items: center;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.6);
  animation: typingBounce 1.2s infinite;
}

.typing-indicator span:nth-child(2) { animation-delay: 0.15s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.3s; }

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-4px); opacity: 1; }
}

/* 通话输入框 */
.call-input-area {
  display: flex;
  gap: 8px;
  width: 100%;
  padding: 8px 0;
  flex-shrink: 0;
}

.call-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 14px;
  outline: none;
}

.call-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.call-input:focus {
  border-color: rgba(255, 255, 255, 0.4);
}

.call-send-btn {
  padding: 8px 16px;
  background: rgba(0, 122, 255, 0.8);
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.call-send-btn:disabled {
  opacity: 0.4;
}

/* 控制按钮 */
.call-controls {
  display: flex;
  gap: 16px;
  margin: 8px 0;
  flex-shrink: 0;
}

.ctrl-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 14px;
  padding: 10px 14px;
  cursor: pointer;
  color: #fff;
  min-width: 60px;
}

.ctrl-btn.active {
  background: rgba(255, 255, 255, 0.25);
}

.ctrl-btn span:first-child { font-size: 20px; }
.ctrl-btn span:last-child { font-size: 10px; opacity: 0.7; }

.hangup-btn {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #e74c3c;
  border: none;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 8px;
  box-shadow: 0 4px 16px rgba(231, 76, 60, 0.4);
  flex-shrink: 0;
}
</style>
