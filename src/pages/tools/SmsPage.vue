<template>
  <div class="sms-page">
    <NavBar title="短信" />

    <!-- 搜索栏 -->
    <div v-if="!activeConv" class="search-bar">
      <input
        v-model="searchText"
        type="text"
        placeholder="搜索短信..."
        class="search-input"
      />
    </div>

    <!-- 对话列表 -->
    <div v-if="!activeConv" class="conversations-list">
      <div v-if="filteredConversations.length === 0" class="empty-state">
        <span class="empty-icon">◌</span>
        <p>暂无短信对话</p>
        <p class="empty-hint">请先在角色管理中创建角色</p>
      </div>
      <div
        v-for="conv in filteredConversations"
        :key="conv.id"
        class="conv-item"
        @click="openChat(conv)"
      >
        <div class="conv-avatar" :style="{ background: conv.color }">
          <img v-if="conv.avatar" :src="conv.avatar" class="conv-avatar-img" />
          <span v-else>{{ conv.name.charAt(0) }}</span>
        </div>
        <div class="conv-info">
          <div class="conv-top">
            <span class="conv-name">{{ conv.name }}</span>
            <span class="conv-time">{{ conv.time }}</span>
          </div>
          <div class="conv-bottom">
            <span class="conv-preview">{{ conv.lastMsg || '暂无消息' }}</span>
            <span v-if="conv.unread" class="conv-badge">{{ conv.unread }}</span>
          </div>
        </div>
      </div>

      <!-- 新建短信按钮 -->
      <button class="fab-btn" @click="showNewSms = true">✎</button>
    </div>

    <!-- 聊天详情 -->
    <template v-if="activeConv">
      <div class="chat-header" @click="closeChat">
        <span class="back-arrow">‹</span>
        <div class="chat-header-info">
          <span class="chat-title">{{ activeConv.name }}</span>
          <span v-if="activeConv.characterId" class="chat-subtitle">AI 角色</span>
        </div>
      </div>

      <!-- API 未配置提示 -->
      <div v-if="activeConv.characterId && !settingsStore.hasChatProviderAccess()" class="api-warning" @click="goToSettings">
        △ 未配置 API Key，AI 无法回复，点击前往设置
      </div>

      <div class="messages-area" ref="messagesRef">
        <div v-if="activeConv.messages.length === 0" class="empty-chat">
          <span>◌</span>
          <p>开始和 {{ activeConv.name }} 发短信吧~</p>
        </div>
        <div
          v-for="msg in activeConv.messages"
          :key="msg.id"
          class="sms-message"
          :class="{ 'is-self': msg.from === 'me' }"
        >
          <div class="sms-bubble" :class="msg.from === 'me' ? 'bubble-blue' : 'bubble-gray'">
            {{ msg.text }}
          </div>
          <span class="sms-time">{{ msg.time }}</span>
        </div>

        <!-- AI 正在输入 -->
        <div v-if="isTyping" class="sms-message">
          <div class="sms-bubble bubble-gray typing-bubble">
            <div v-if="streamingText" class="streaming-text">{{ streamingText }}</div>
            <div v-else class="typing-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
          <span class="sms-time">{{ activeConv.name }} 正在输入...</span>
        </div>
      </div>

      <div class="sms-input-area">
        <input
          v-model="smsText"
          type="text"
          placeholder="输入短信..."
          class="sms-input"
          @keyup.enter="sendSms"
        />
        <button class="sms-send" :disabled="!smsText.trim() || isTyping" @click="sendSms">
          发送
        </button>
      </div>
    </template>

    <!-- 新建短信弹窗 -->
    <div v-if="showNewSms" class="overlay" @click.self="showNewSms = false">
      <div class="new-sms-panel">
        <div class="panel-header">
          <h3>新短信</h3>
          <button class="close-btn" @click="showNewSms = false">✕</button>
        </div>

        <!-- 角色选择 -->
        <div v-if="availableCharacters.length > 0" class="char-section">
          <div class="section-label">选择角色</div>
          <div class="char-list">
            <div
              v-for="char in availableCharacters"
              :key="char.id"
              class="char-option"
              @click="createFromCharacter(char)"
            >
              <div class="char-avatar" :style="{ background: getCharColor(char.id) }">
                <img v-if="char.avatar" :src="char.avatar" class="char-avatar-img" />
                <span v-else>{{ (char.name || '?').charAt(0) }}</span>
              </div>
              <span class="char-name">{{ char.name }}</span>
            </div>
          </div>
        </div>

        <div class="divider-line">
          <span>或自定义</span>
        </div>

        <div class="form-group">
          <label>收件人</label>
          <input v-model="newRecipient" placeholder="输入名称" />
        </div>
        <div class="form-group">
          <label>号码</label>
          <input v-model="newNumber" placeholder="输入号码（可选）" />
        </div>
        <div class="form-group">
          <label>内容</label>
          <textarea v-model="newContent" placeholder="输入第一条短信..." rows="3"></textarea>
        </div>
        <button class="send-new-btn" :disabled="!newRecipient" @click="sendNewSms">
          发送
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { useSmsStore } from '@/stores/sms'
import { useSettingsStore } from '@/stores/settings'
import { useCharactersStore } from '@/stores/characters'
import {
  sendAIRequest,
  buildSmsMessages,
  getCharacterById,
} from '@/utils/aiService'
import type { CharacterData } from '@/utils/aiService'

const router = useRouter()
const smsStore = useSmsStore()
const settingsStore = useSettingsStore()
const charactersStore = useCharactersStore()

const searchText = ref('')
const smsText = ref('')
const showNewSms = ref(false)
const newRecipient = ref('')
const newNumber = ref('')
const newContent = ref('')
const messagesRef = ref<HTMLElement | null>(null)
const isTyping = ref(false)
const streamingText = ref('')
let abortController: AbortController | null = null

// 当前活跃对话
const activeConv = computed(() => smsStore.currentConversation)

// 过滤后的对话列表
const filteredConversations = computed(() => {
  const list = smsStore.sortedConversations
  if (!searchText.value) return list
  const q = searchText.value.toLowerCase()
  return list.filter(
    (c) => c.name.toLowerCase().includes(q) || c.lastMsg.toLowerCase().includes(q)
  )
})

// 可用角色列表（用于新建短信）
const availableCharacters = computed(() => {
  return charactersStore.charItems.map((item: any) => ({
    id: String(item.id),
    type: 'char',
    name: item.name,
    avatar: item.avatar_url || '',
  }))
})

function getCharColor(id: string): string {
  const colors = ['#ff2d55', '#007aff', '#ff9500', '#5856d6', '#34c759', '#af52de', '#ff6348']
  let hash = 0
  for (let i = 0; i < id.length; i++) {
    hash = ((hash << 5) - hash) + id.charCodeAt(i)
    hash |= 0
  }
  return colors[Math.abs(hash) % colors.length]
}

function openChat(conv: any) {
  smsStore.openConversation(conv.id)
  nextTick(() => scrollToBottom())
}

function closeChat() {
  smsStore.closeConversation()
}

function goToSettings() {
  router.push('/customize')
}

async function sendSms() {
  if (!smsText.value.trim() || isTyping.value) return
  const conv = activeConv.value
  if (!conv) return

  const text = smsText.value.trim()
  smsText.value = ''

  // 发送用户消息
  smsStore.sendMessage(conv.id, text)
  scrollToBottom()

  // 如果关联了角色，触发 AI 回复
  if (conv.characterId) {
    await triggerAIReply(conv.id, conv.characterId)
  }
}

async function triggerAIReply(conversationId: string, characterId: string) {
  const s = settingsStore.settings
  if (!s.apiKey) return

  const character = getCharacterById(characterId)
  if (!character) return

  // 模拟打字延迟（1-3秒）
  isTyping.value = true
  streamingText.value = ''
  await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))

  abortController = new AbortController()

  try {
    // 构建短信消息历史
    const history = smsStore.getMessageHistory(conversationId, 20)
    const aiMessages = buildSmsMessages(
      character,
      history.map(m => ({ from: m.from, text: m.text })),
    )

    const apiUrl = settingsStore.getApiUrl()
    const isStream = s.streamEnabled

    const response = await sendAIRequest({
      apiKey: s.apiKey,
      apiUrl: apiUrl,
      model: s.model,
      messages: aiMessages,
      temperature: s.temperature,
      maxTokens: 200, // 短信回复限制更短
      stream: isStream,
      timeout: s.timeout,
      signal: abortController.signal,
      onChunk: (chunk: string) => {
        streamingText.value += chunk
        scrollToBottom()
      },
    })

    isTyping.value = false
    const content = isStream ? streamingText.value.trim() : response.content
    streamingText.value = ''

    if (!content) {
      smsStore.addAIReply(conversationId, '（消息发送失败）')
      scrollToBottom()
      return
    }

    // 添加 AI 回复
    smsStore.addAIReply(conversationId, content)
    scrollToBottom()
  } catch (err: any) {
    isTyping.value = false
    streamingText.value = ''

    if (err.name === 'AbortError' || err.message?.includes('abort')) {
      return
    }

    smsStore.addAIReply(conversationId, `✕ 回复失败：${err.message}`)
    scrollToBottom()
  } finally {
    abortController = null
  }
}

function createFromCharacter(char: any) {
  const conv = smsStore.createConversation(char.id)
  smsStore.openConversation(conv.id)
  showNewSms.value = false
  nextTick(() => scrollToBottom())
}

function sendNewSms() {
  if (!newRecipient.value) return

  const conv = smsStore.createCustomConversation(
    newRecipient.value,
    newNumber.value || newRecipient.value,
  )

  if (newContent.value.trim()) {
    smsStore.sendMessage(conv.id, newContent.value.trim())
  }

  smsStore.openConversation(conv.id)
  newRecipient.value = ''
  newNumber.value = ''
  newContent.value = ''
  showNewSms.value = false
  nextTick(() => scrollToBottom())
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

onMounted(() => {
  charactersStore.fetchCharacters()
  smsStore.loadConversations()
})
</script>

<style scoped>
.sms-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.search-bar {
  padding: 8px 16px;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  border: none;
  border-radius: 12px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 32px;
  color: var(--text-tertiary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.empty-state p {
  margin: 4px 0;
  font-size: 15px;
}

.empty-hint {
  font-size: 13px !important;
  opacity: 0.7;
}

.conversations-list {
  flex: 1;
  overflow-y: auto;
  position: relative;
}

.conv-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.15s;
}

.conv-item:active {
  background: var(--bg-tertiary);
}

.conv-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  flex-shrink: 0;
  overflow: hidden;
}

.conv-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conv-info {
  flex: 1;
  min-width: 0;
}

.conv-top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.conv-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.conv-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.conv-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.conv-preview {
  font-size: 13px;
  color: var(--text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 220px;
}

.conv-badge {
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: var(--accent-red);
  color: #fff;
  border-radius: 10px;
  font-size: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.fab-btn {
  position: fixed;
  bottom: 80px;
  right: 24px;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: var(--accent-blue);
  border: none;
  font-size: 22px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.3);
  z-index: 10;
}

/* API 警告 */
.api-warning {
  background: linear-gradient(135deg, #ff9500, #ff6348);
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  text-align: center;
  padding: 8px 16px;
  cursor: pointer;
}

/* 聊天详情 */
.chat-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}

.back-arrow {
  font-size: 24px;
  color: var(--accent-blue);
  font-weight: 300;
}

.chat-header-info {
  display: flex;
  flex-direction: column;
}

.chat-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.chat-subtitle {
  font-size: 11px;
  color: var(--text-tertiary);
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: var(--text-tertiary);
}

.empty-chat span {
  font-size: 40px;
  margin-bottom: 8px;
}

.empty-chat p {
  font-size: 14px;
  margin: 0;
}

.sms-message {
  display: flex;
  flex-direction: column;
  margin-bottom: 12px;
}

.sms-message.is-self {
  align-items: flex-end;
}

.sms-bubble {
  max-width: 80%;
  padding: 10px 14px;
  border-radius: 18px;
  font-size: 15px;
  line-height: 1.5;
  word-break: break-word;
}

.bubble-blue {
  background: var(--accent-blue);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.bubble-gray {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-bottom-left-radius: 4px;
}

.sms-time {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 3px;
  padding: 0 4px;
}

/* 打字动画 */
.typing-bubble {
  padding: 14px 18px;
}

.typing-dots {
  display: flex;
  gap: 5px;
  align-items: center;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-tertiary);
  animation: typingBounce 1.2s infinite;
}

.typing-dots span:nth-child(2) { animation-delay: 0.15s; }
.typing-dots span:nth-child(3) { animation-delay: 0.3s; }

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-6px); opacity: 1; }
}

.streaming-text {
  font-size: 15px;
  line-height: 1.5;
}

.sms-input-area {
  display: flex;
  gap: 8px;
  padding: 8px 16px 16px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.sms-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--border-color);
  border-radius: 20px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
}

.sms-input:focus {
  border-color: var(--accent-blue);
}

.sms-send {
  padding: 8px 16px;
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  transition: opacity 0.15s;
}

.sms-send:disabled {
  opacity: 0.4;
}

/* 新建短信 */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 24px;
}

.new-sms-panel {
  width: 100%;
  max-width: 360px;
  max-height: 80vh;
  background: var(--bg-secondary);
  border-radius: 20px;
  padding: 24px;
  overflow-y: auto;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.panel-header h3 {
  margin: 0;
  font-size: 17px;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: var(--text-tertiary);
  cursor: pointer;
}

/* 角色选择 */
.char-section {
  margin-bottom: 16px;
}

.section-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 10px;
  font-weight: 600;
}

.char-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.char-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 8px;
  border-radius: 12px;
  transition: background 0.15s;
  min-width: 64px;
}

.char-option:active {
  background: var(--bg-tertiary);
}

.char-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  overflow: hidden;
}

.char-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.char-name {
  font-size: 12px;
  color: var(--text-primary);
  max-width: 64px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: center;
}

.divider-line {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
  color: var(--text-tertiary);
  font-size: 12px;
}

.divider-line::before,
.divider-line::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

.form-group {
  margin-bottom: 14px;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
  font-family: inherit;
  resize: none;
}

.send-new-btn {
  width: 100%;
  padding: 12px;
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  cursor: pointer;
}

.send-new-btn:disabled {
  opacity: 0.4;
}
</style>
