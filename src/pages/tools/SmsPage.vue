<template>
  <div class="sms-page">
    <NavBar title="短信" />

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input
        v-model="searchText"
        type="text"
        placeholder="搜索短信..."
        class="search-input"
      />
    </div>

    <!-- 对话列表 -->
    <div v-if="!activeChat" class="conversations-list">
      <div
        v-for="conv in filteredConversations"
        :key="conv.id"
        class="conv-item"
        @click="openChat(conv)"
      >
        <div class="conv-avatar" :style="{ background: conv.color }">
          {{ conv.name.charAt(0) }}
        </div>
        <div class="conv-info">
          <div class="conv-top">
            <span class="conv-name">{{ conv.name }}</span>
            <span class="conv-time">{{ conv.time }}</span>
          </div>
          <div class="conv-bottom">
            <span class="conv-preview">{{ conv.lastMsg }}</span>
            <span v-if="conv.unread" class="conv-badge">{{ conv.unread }}</span>
          </div>
        </div>
      </div>

      <!-- 新建短信按钮 -->
      <button class="fab-btn" @click="showNewSms = true">✏️</button>
    </div>

    <!-- 聊天详情 -->
    <template v-if="activeChat">
      <div class="chat-header" @click="activeChat = null">
        <span class="back-arrow">‹</span>
        <span class="chat-title">{{ activeChat.name }}</span>
      </div>

      <div class="messages-area" ref="messagesRef">
        <div
          v-for="msg in activeChat.messages"
          :key="msg.id"
          class="sms-message"
          :class="{ 'is-self': msg.from === 'me' }"
        >
          <div class="sms-bubble" :class="msg.from === 'me' ? 'bubble-blue' : 'bubble-gray'">
            {{ msg.text }}
          </div>
          <span class="sms-time">{{ msg.time }}</span>
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
        <button class="sms-send" :disabled="!smsText.trim()" @click="sendSms">
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
        <div class="form-group">
          <label>收件人</label>
          <input v-model="newRecipient" placeholder="输入号码或姓名" />
        </div>
        <div class="form-group">
          <label>内容</label>
          <textarea v-model="newContent" placeholder="输入短信内容..." rows="4"></textarea>
        </div>
        <button class="send-new-btn" :disabled="!newRecipient || !newContent" @click="sendNewSms">
          发送
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

interface SmsMessage {
  id: number
  from: 'me' | 'other'
  text: string
  time: string
}

interface SmsConversation {
  id: number
  name: string
  number: string
  color: string
  lastMsg: string
  time: string
  unread: number
  messages: SmsMessage[]
}

const searchText = ref('')
const smsText = ref('')
const activeChat = ref<SmsConversation | null>(null)
const showNewSms = ref(false)
const newRecipient = ref('')
const newContent = ref('')
const messagesRef = ref<HTMLElement | null>(null)

const conversations = ref<SmsConversation[]>([
  {
    id: 1, name: '10086', number: '10086', color: '#007aff',
    lastMsg: '您本月话费已出账，金额58.00元', time: '14:30', unread: 1,
    messages: [
      { id: 1, from: 'other', text: '【中国移动】尊敬的客户，您本月话费已出账，金额58.00元，请及时缴费。', time: '14:30' },
      { id: 2, from: 'other', text: '【中国移动】您的流量已使用80%，建议购买流量包。', time: '昨天 10:00' },
    ],
  },
  {
    id: 2, name: '宝贝', number: '13800138000', color: '#ff2d55',
    lastMsg: '到家了吗？', time: '12:15', unread: 0,
    messages: [
      { id: 1, from: 'other', text: '你在干嘛呀', time: '12:00' },
      { id: 2, from: 'me', text: '在忙工作呢~', time: '12:05' },
      { id: 3, from: 'other', text: '到家了吗？', time: '12:15' },
      { id: 4, from: 'me', text: '快了快了，马上到', time: '12:16' },
    ],
  },
  {
    id: 3, name: '快递', number: '95311', color: '#ff9500',
    lastMsg: '您的快递已签收', time: '昨天', unread: 0,
    messages: [
      { id: 1, from: 'other', text: '【顺丰速运】您的包裹已到达，请凭取件码到驿站取件。取件码：6688', time: '昨天 09:30' },
      { id: 2, from: 'other', text: '【顺丰速运】您的快递已签收，感谢使用顺丰。', time: '昨天 15:00' },
    ],
  },
  {
    id: 4, name: '验证码', number: '106xxxx', color: '#5856d6',
    lastMsg: '您的验证码为 886432', time: '前天', unread: 0,
    messages: [
      { id: 1, from: 'other', text: '【某APP】您的验证码为 886432，5分钟内有效，请勿泄露。', time: '前天 20:15' },
    ],
  },
])

const filteredConversations = computed(() => {
  if (!searchText.value) return conversations.value
  const q = searchText.value.toLowerCase()
  return conversations.value.filter(
    (c) => c.name.toLowerCase().includes(q) || c.lastMsg.toLowerCase().includes(q)
  )
})

function openChat(conv: SmsConversation) {
  activeChat.value = conv
  conv.unread = 0
  nextTick(() => scrollToBottom())
}

function sendSms() {
  if (!smsText.value.trim() || !activeChat.value) return
  activeChat.value.messages.push({
    id: Date.now(),
    from: 'me',
    text: smsText.value.trim(),
    time: '刚刚',
  })
  activeChat.value.lastMsg = smsText.value.trim()
  activeChat.value.time = '刚刚'
  smsText.value = ''
  scrollToBottom()
}

function sendNewSms() {
  if (!newRecipient.value || !newContent.value) return
  const newConv: SmsConversation = {
    id: Date.now(),
    name: newRecipient.value,
    number: newRecipient.value,
    color: '#34c759',
    lastMsg: newContent.value,
    time: '刚刚',
    unread: 0,
    messages: [
      { id: 1, from: 'me', text: newContent.value, time: '刚刚' },
    ],
  }
  conversations.value.unshift(newConv)
  newRecipient.value = ''
  newContent.value = ''
  showNewSms.value = false
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesRef.value) {
      messagesRef.value.scrollTop = messagesRef.value.scrollHeight
    }
  })
}
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

.chat-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
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
  max-width: 340px;
  background: var(--bg-secondary);
  border-radius: 20px;
  padding: 24px;
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
