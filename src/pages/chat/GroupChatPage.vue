<template>
  <div class="group-chat-page">
    <NavBar title="群聊">
      <template #right>
        <button class="members-btn" @click="showMembers = !showMembers">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span class="member-count">{{ members.length }}</span>
        </button>
      </template>
    </NavBar>

    <!-- No API Key warning -->
    <div v-if="!settingsStore.settings.apiKey" class="api-warning" @click="goToSettings">
      <span>未配置 API Key，回复将使用模拟模式</span>
    </div>

    <!-- 消息列表 -->
    <div class="messages-area" ref="messagesRef">
      <div v-if="messages.length === 0" class="empty-chat">
        <svg viewBox="0 0 48 48" width="56" height="56" fill="none" stroke="var(--text-quaternary)" stroke-width="1.5">
          <rect x="4" y="6" width="40" height="32" rx="6"/>
          <path d="M16 26h16M16 20h10"/>
          <path d="M14 38l-4 6v-6"/>
        </svg>
        <p>群聊已创建，开始聊天吧~</p>
      </div>

      <div
        v-for="msg in messages"
        :key="msg.id"
        class="message-row"
        :class="{ 'is-self': msg.role === 'user' }"
      >
        <div v-if="msg.role !== 'user'" class="msg-avatar" :style="{ background: getMemberColor(msg.senderName) }">
          <span>{{ (msg.senderName || '?').charAt(0) }}</span>
        </div>
        <div class="msg-content">
          <span v-if="msg.role !== 'user'" class="msg-sender">{{ msg.senderName }}</span>
          <div class="msg-bubble" :class="msg.role === 'user' ? 'bubble-self' : 'bubble-other'">
            {{ msg.content }}
          </div>
          <span class="msg-time">{{ formatTime(msg.created_at) }}</span>
        </div>
        <div v-if="msg.role === 'user'" class="msg-avatar self-avatar">
          <span>我</span>
        </div>
      </div>

      <!-- Typing indicator -->
      <div v-if="typingMember" class="message-row">
        <div class="msg-avatar" :style="{ background: getMemberColor(typingMember) }">
          <span>{{ typingMember.charAt(0) }}</span>
        </div>
        <div class="msg-content">
          <span class="msg-sender">{{ typingMember }}</span>
          <div class="msg-bubble bubble-other typing-bubble">
            <div v-if="streamingText" class="streaming-text">{{ streamingText }}</div>
            <div v-else class="typing-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区 -->
    <div class="input-area">
      <div class="input-row">
        <input
          v-model="inputText"
          type="text"
          placeholder="发送消息..."
          class="chat-input"
          @keyup.enter="sendMessage"
          :disabled="isReplying"
        />
        <button class="send-btn" :disabled="!inputText.trim() || isReplying" @click="sendMessage">
          <svg viewBox="0 0 24 24" fill="currentColor" width="18" height="18">
            <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 成员面板 -->
    <div v-if="showMembers" class="members-overlay" @click.self="showMembers = false">
      <div class="members-panel">
        <div class="panel-header">
          <h3>群成员 ({{ members.length }})</h3>
          <button class="close-btn" @click="showMembers = false">✕</button>
        </div>
        <div class="members-list">
          <div v-for="member in members" :key="member.id" class="member-item">
            <div class="member-avatar" :style="{ background: member.color }">
              {{ member.name.charAt(0) }}
            </div>
            <div class="member-info">
              <span class="member-name">{{ member.name }}</span>
              <span class="member-role">{{ member.role }}</span>
            </div>
            <span v-if="member.isOwner" class="owner-badge">群主</span>
          </div>
        </div>

        <div class="panel-actions">
          <button class="action-btn add-btn" @click="showAddMember = true">
            + 添加成员
          </button>
        </div>
      </div>
    </div>

    <!-- 添加成员弹窗 -->
    <div v-if="showAddMember" class="members-overlay" @click.self="showAddMember = false">
      <div class="members-panel">
        <div class="panel-header">
          <h3>添加角色到群聊</h3>
          <button class="close-btn" @click="showAddMember = false">✕</button>
        </div>
        <div v-if="availableCharacters.length === 0" class="empty-chars">
          <p>没有更多角色可添加</p>
        </div>
        <div v-else class="members-list">
          <div
            v-for="char in availableCharacters"
            :key="char.id"
            class="member-item clickable"
            @click="addMember(char)"
          >
            <div class="member-avatar">
              {{ char.name.charAt(0) }}
            </div>
            <div class="member-info">
              <span class="member-name">{{ char.name }}</span>
              <span class="member-role">{{ char.description?.slice(0, 30) || '角色' }}</span>
            </div>
            <span class="add-icon">+</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { useChatStore, type Character } from '@/stores/chat'
import { useSettingsStore } from '@/stores/settings'
import {
  sendAIRequest,
  buildSystemPrompt,
  getActivePreset,
  getCurrentUserPersona,
  getCharacterById,
} from '@/utils/aiService'
import type { AIMessage, CharacterData } from '@/utils/aiService'

interface GroupMember {
  id: string
  name: string
  role: string
  color: string
  isOwner: boolean
  characterData?: CharacterData | null
}

interface GroupMessage {
  id: number
  role: string
  senderName: string
  content: string
  created_at: string
}

const route = useRoute()
const router = useRouter()
const chatStore = useChatStore()
const settingsStore = useSettingsStore()
const groupId = route.params.groupId as string

const messages = ref<GroupMessage[]>([])
const inputText = ref('')
const showMembers = ref(false)
const showAddMember = ref(false)
const messagesRef = ref<HTMLElement | null>(null)
const isReplying = ref(false)
const typingMember = ref('')
const streamingText = ref('')

const STORAGE_KEY = `group-messages-${groupId}`
const MEMBERS_KEY = `group-members-${groupId}`

const memberColors = [
  '#5856d6', '#af52de', '#ff2d55', '#ff9500',
  '#34c759', '#e74c3c', '#1abc9c', '#f39c12',
  '#007aff', '#5ac8fa', '#ff6b6b', '#6c5ce7',
]

const members = ref<GroupMember[]>([
  { id: '0', name: '我', role: '群主', color: '#007aff', isOwner: true },
])

const availableCharacters = computed(() => {
  const memberIds = new Set(members.value.map((m) => m.id))
  return chatStore.characters.filter((c) => !memberIds.has(String(c.id)))
})

// Load persisted data
onMounted(async () => {
  await chatStore.fetchCharacters()

  // Restore members from conversation data or localStorage
  const conv = chatStore.conversations.find(c => c.id === groupId)
  const savedMembers = loadSavedMembers()

  if (savedMembers && savedMembers.length > 1) {
    // Restore from localStorage, re-hydrate character data
    members.value = savedMembers.map(m => ({
      ...m,
      characterData: m.id !== '0' ? getCharacterById(m.id) : null,
    }))
  } else if (conv?.groupMemberIds && conv.groupMemberIds.length > 0) {
    // First time opening — build members from conversation's groupMemberIds
    for (const charId of conv.groupMemberIds) {
      const char = getCharacterById(charId)
      if (char) {
        addMemberFromCharData(char)
      }
    }
    saveMembers()
  } else if (members.value.length === 1 && chatStore.characters.length > 0) {
    // Fallback: auto-add first 2 characters
    const toAdd = chatStore.characters.slice(0, 2)
    toAdd.forEach((c) => {
      const char = getCharacterById(String(c.id))
      if (char) addMemberFromCharData(char)
    })
    saveMembers()
  }

  // Load persisted messages
  loadMessages()
})

function addMemberFromCharData(char: CharacterData) {
  const colorIdx = members.value.length % memberColors.length
  members.value.push({
    id: String(char.id),
    name: char.name || '未命名',
    role: char.description?.slice(0, 20) || '角色',
    color: memberColors[colorIdx],
    isOwner: false,
    characterData: char,
  })
}

function addMember(char: Character) {
  const charData = getCharacterById(String(char.id))
  addMemberFromCharData(charData || char as any)
  showAddMember.value = false
  saveMembers()

  // Update conversation's groupMemberIds
  const conv = chatStore.conversations.find(c => c.id === groupId)
  if (conv) {
    conv.groupMemberIds = members.value.filter(m => m.id !== '0').map(m => m.id)
    chatStore.saveConversations()
  }
}

function getMemberColor(name: string): string {
  const member = members.value.find(m => m.name === name)
  if (member) return member.color
  // Fallback: hash-based color
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash)
  return memberColors[Math.abs(hash) % memberColors.length]
}

// Persistence helpers
function loadMessages() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) messages.value = JSON.parse(saved)
  } catch { /* ignore */ }
}

function saveMessages() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(messages.value.slice(-200)))
  } catch { /* ignore */ }
}

function loadSavedMembers(): GroupMember[] | null {
  try {
    const saved = localStorage.getItem(MEMBERS_KEY)
    if (saved) return JSON.parse(saved)
  } catch { /* ignore */ }
  return null
}

function saveMembers() {
  try {
    // Don't save characterData to localStorage
    const toSave = members.value.map(m => ({
      id: m.id, name: m.name, role: m.role,
      color: m.color, isOwner: m.isOwner,
    }))
    localStorage.setItem(MEMBERS_KEY, JSON.stringify(toSave))
  } catch { /* ignore */ }
}

function goToSettings() {
  router.push('/customize')
}

async function sendMessage() {
  if (!inputText.value.trim() || isReplying.value) return
  const text = inputText.value.trim()
  inputText.value = ''

  // Add user message
  const userMsg: GroupMessage = {
    id: Date.now(),
    role: 'user',
    senderName: '我',
    content: text,
    created_at: new Date().toISOString(),
  }
  messages.value.push(userMsg)
  saveMessages()
  scrollToBottom()

  // Update conversation preview
  const conv = chatStore.conversations.find(c => c.id === groupId)
  if (conv) {
    conv.last_message = text.slice(0, 50)
    conv.last_at = userMsg.created_at
    chatStore.saveConversations()
  }

  // Generate replies from group members
  const otherMembers = members.value.filter((m) => m.id !== '0')
  if (otherMembers.length === 0) return

  isReplying.value = true

  // Pick random 1-3 members to reply
  const replyCount = Math.min(otherMembers.length, Math.floor(Math.random() * 3) + 1)
  const shuffled = [...otherMembers].sort(() => Math.random() - 0.5)
  const repliers = shuffled.slice(0, replyCount)

  const hasApiKey = !!settingsStore.settings.apiKey

  for (let i = 0; i < repliers.length; i++) {
    const replier = repliers[i]
    typingMember.value = replier.name
    streamingText.value = ''
    scrollToBottom()

    let replyContent: string

    if (hasApiKey) {
      try {
        replyContent = await getAIReply(replier, text)
      } catch {
        replyContent = generateFallbackReply(replier.name, text)
      }
    } else {
      await delay(600 + Math.random() * 1000)
      replyContent = generateFallbackReply(replier.name, text)
    }

    typingMember.value = ''
    streamingText.value = ''

    const reply: GroupMessage = {
      id: Date.now() + i + 1,
      role: 'assistant',
      senderName: replier.name,
      content: replyContent,
      created_at: new Date().toISOString(),
    }
    messages.value.push(reply)
    saveMessages()
    scrollToBottom()

    // Update conversation preview
    if (conv) {
      conv.last_message = `${replier.name}: ${replyContent.slice(0, 40)}`
      conv.last_at = reply.created_at
      chatStore.saveConversations()
    }

    // Small gap between different members replying
    if (i < repliers.length - 1) {
      await delay(300 + Math.random() * 500)
    }
  }

  isReplying.value = false
}

async function getAIReply(member: GroupMember, userText: string): Promise<string> {
  const s = settingsStore.settings
  const charData = member.characterData || getCharacterById(member.id)
  const userPersona = getCurrentUserPersona()
  const preset = getActivePreset()

  // Build a group-chat-aware system prompt
  const otherNames = members.value
    .filter(m => m.id !== '0' && m.id !== member.id)
    .map(m => m.name)

  let systemPrompt = buildSystemPrompt(charData, userPersona, s.maxLength, [], preset)
  systemPrompt += `\n\n<群聊规则>
你正在一个群聊中，群里还有：${otherNames.join('、')}。
回复要简短自然，像真人群聊一样，每条消息1-3句话。
用中文回复。不需要加角色名前缀。
</群聊规则>`

  // Build recent message history
  const contextLimit = s.contextSize || 20
  const recent = messages.value.slice(-contextLimit)
  const aiMessages: AIMessage[] = [
    { role: 'system', content: systemPrompt },
    ...recent.map(m => ({
      role: (m.role === 'user' ? 'user' : 'assistant') as 'user' | 'assistant',
      content: m.role === 'user' ? m.content : `[${m.senderName}]: ${m.content}`,
    })),
  ]

  const apiUrl = settingsStore.getApiUrl()
  const isStream = s.streamEnabled

  const response = await sendAIRequest({
    apiKey: s.apiKey,
    apiUrl: apiUrl,
    model: s.model,
    messages: aiMessages,
    temperature: Math.min((s.temperature || 0.7) + 0.1, 1.5),
    maxTokens: Math.min(s.maxLength || 1000, 1000), // 群聊保持简短
    stream: isStream,
    timeout: s.timeout,
    onChunk: (chunk: string) => {
      streamingText.value += chunk
      scrollToBottom()
    },
  })

  let content = isStream ? streamingText.value.trim() : response.content
  // Remove any accidental role prefixes
  content = content.replace(/^\[?[\u4e00-\u9fff]+\]?[:：]\s*/g, '').trim()
  return content || generateFallbackReply(member.name, userText)
}

function generateFallbackReply(_name: string, userMsg: string): string {
  const replies = [
    `嗯嗯，我觉得是这样的~`,
    `${userMsg.slice(0, 10)}...这个话题很有趣！`,
    `哈哈哈，说得对！`,
    `我也这么想的呢`,
    `嗯...让我想想`,
    `有道理！`,
    `真的吗？太棒了！`,
    `我有不同看法哦~`,
    `+1`,
    `了解了~`,
    `确实如此`,
    `哈哈好有意思`,
  ]
  return replies[Math.floor(Math.random() * replies.length)]
}

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}

function formatTime(dateStr: string): string {
  const d = new Date(dateStr)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style scoped>
.group-chat-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.members-btn {
  background: none;
  border: none;
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--accent-blue, #007aff);
  cursor: pointer;
  padding: 4px;
}

.member-count {
  font-size: 13px;
  font-weight: 600;
}

.api-warning {
  padding: 8px 16px;
  background: rgba(255, 149, 0, 0.12);
  color: var(--accent-orange, #ff9500);
  font-size: 13px;
  text-align: center;
  cursor: pointer;
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
}

.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-tertiary);
  gap: 12px;
}

.empty-chat p {
  font-size: 14px;
  margin: 0;
}

.message-row {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  align-items: flex-start;
}

.message-row.is-self {
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.self-avatar {
  background: var(--accent-blue, #007aff);
}

.msg-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.is-self .msg-content {
  align-items: flex-end;
}

.msg-sender {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 2px;
  padding-left: 4px;
}

.msg-bubble {
  padding: 10px 14px;
  border-radius: 18px;
  font-size: 15px;
  line-height: 1.5;
  word-break: break-word;
}

.bubble-other {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-bottom-left-radius: 4px;
}

.bubble-self {
  background: var(--accent-blue, #007aff);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.msg-time {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 2px;
  padding: 0 4px;
}

/* Typing indicator */
.typing-bubble {
  min-width: 60px;
}

.typing-dots {
  display: flex;
  gap: 4px;
  padding: 2px 0;
}

.typing-dots span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-tertiary);
  animation: typingBounce 1.4s ease-in-out infinite;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-4px); opacity: 1; }
}

.streaming-text {
  font-size: 15px;
  line-height: 1.5;
}

/* 输入区 */
.input-area {
  padding: 8px 16px 16px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

.input-row {
  display: flex;
  gap: 8px;
}

.chat-input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid var(--border-color);
  border-radius: 20px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 15px;
  outline: none;
}

.chat-input:focus {
  border-color: var(--accent-blue, #007aff);
}

.chat-input:disabled {
  opacity: 0.6;
}

.send-btn {
  width: 40px;
  height: 40px;
  background: var(--accent-blue, #007aff);
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: opacity 0.2s;
}

.send-btn:disabled {
  opacity: 0.4;
}

/* 成员面板 */
.members-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.members-panel {
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

.members-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 16px;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.member-item:last-child {
  border-bottom: none;
}

.member-item.clickable {
  cursor: pointer;
  border-radius: 8px;
  padding: 10px 8px;
}

.member-item.clickable:active {
  background: var(--bg-tertiary);
}

.member-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-purple, #af52de);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.member-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.member-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.member-role {
  font-size: 12px;
  color: var(--text-tertiary);
}

.owner-badge {
  font-size: 11px;
  color: var(--accent-orange, #ff9500);
  background: rgba(255, 149, 0, 0.12);
  padding: 2px 8px;
  border-radius: 10px;
}

.add-icon {
  font-size: 20px;
  color: var(--accent-blue, #007aff);
  font-weight: 600;
}

.panel-actions {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}

.action-btn {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  cursor: pointer;
}

.add-btn {
  background: var(--accent-blue, #007aff);
  color: #fff;
}

.empty-chars {
  padding: 30px 16px;
  text-align: center;
  color: var(--text-tertiary);
}
</style>
