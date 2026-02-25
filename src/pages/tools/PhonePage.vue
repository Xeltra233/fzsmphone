<template>
  <div class="phone-page">
    <NavBar title="电话" />

    <!-- Tab 切换 -->
    <div class="tab-bar">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 拨号盘 -->
    <div v-if="activeTab === 'dial'" class="dial-view">
      <div class="number-display">
        <span class="phone-number">{{ displayNumber || '输入号码' }}</span>
        <button v-if="phoneNumber" class="backspace-btn" @click="backspace">⌫</button>
      </div>

      <!-- 匹配的联系人提示 -->
      <div v-if="matchedContact" class="matched-contact">
        <div class="matched-avatar" :style="{ background: matchedContact.color }">
          <img v-if="matchedContact.avatar" :src="matchedContact.avatar" class="matched-avatar-img" />
          <span v-else>{{ matchedContact.name.charAt(0) }}</span>
        </div>
        <span class="matched-name">{{ matchedContact.name }}</span>
        <span class="matched-tag">AI 角色</span>
      </div>

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

      <div class="action-row">
        <button class="action-btn" @click="activeTab = 'contacts'">
          <span>👤</span>
          <span class="action-label">通讯录</span>
        </button>
        <button class="call-btn" @click="makeCall('voice')">
          <span>📞</span>
        </button>
        <button class="action-btn" @click="makeCall('video')">
          <span>📹</span>
          <span class="action-label">视频通话</span>
        </button>
      </div>
    </div>

    <!-- 通讯录 -->
    <div v-if="activeTab === 'contacts'" class="contacts-view">
      <div v-if="phoneStore.contacts.length === 0" class="empty-state">
        <span class="empty-icon">👤</span>
        <p>暂无联系人</p>
        <p class="empty-hint">请先在角色管理中创建角色</p>
      </div>
      <div
        v-for="contact in phoneStore.contacts"
        :key="contact.id"
        class="contact-item"
      >
        <div class="contact-avatar" :style="{ background: contact.color }">
          <img v-if="contact.avatar" :src="contact.avatar" class="contact-avatar-img" />
          <span v-else>{{ contact.name.charAt(0) }}</span>
        </div>
        <div class="contact-info">
          <span class="contact-name">{{ contact.name }}</span>
          <span class="contact-number">{{ formatPhoneDisplay(contact.number) }}</span>
        </div>
        <div class="contact-actions">
          <button class="contact-call-btn" @click="callContact(contact, 'voice')">📞</button>
          <button class="contact-call-btn" @click="callContact(contact, 'video')">📹</button>
        </div>
      </div>
    </div>

    <!-- 通话记录 -->
    <div v-if="activeTab === 'history'" class="history-view">
      <div v-if="phoneStore.callHistory.length === 0" class="empty-state">
        <span class="empty-icon">📱</span>
        <p>暂无通话记录</p>
      </div>
      <div
        v-for="call in phoneStore.callHistory"
        :key="call.id"
        class="history-item"
        @click="redial(call)"
      >
        <div class="call-icon" :class="call.type">
          {{ call.type === 'incoming' ? '📥' : call.type === 'outgoing' ? '📤' : '📵' }}
        </div>
        <div class="call-info">
          <span class="call-name" :class="{ missed: call.type === 'missed' }">
            {{ call.name || call.number }}
          </span>
          <span class="call-meta">
            {{ call.callType === 'video' ? '📹 视频' : '📞 语音' }} · {{ call.time }}
          </span>
        </div>
        <span class="call-duration">{{ call.duration || '--' }}</span>
      </div>
      <button v-if="phoneStore.callHistory.length > 0" class="clear-history-btn" @click="clearHistory">
        清除所有通话记录
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { usePhoneStore } from '@/stores/phone'
import type { Contact } from '@/stores/phone'

const router = useRouter()
const phoneStore = usePhoneStore()

const phoneNumber = ref('')
const activeTab = ref<'dial' | 'contacts' | 'history'>('dial')

const tabs = [
  { key: 'dial' as const, label: '拨号' },
  { key: 'contacts' as const, label: '通讯录' },
  { key: 'history' as const, label: '记录' },
]

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

const displayNumber = computed(() => {
  const num = phoneNumber.value
  if (num.length <= 3) return num
  if (num.length <= 7) return `${num.slice(0, 3)} ${num.slice(3)}`
  return `${num.slice(0, 3)} ${num.slice(3, 7)} ${num.slice(7)}`
})

// 根据输入的号码匹配联系人
const matchedContact = computed(() => {
  if (!phoneNumber.value || phoneNumber.value.length < 3) return null
  return phoneStore.contacts.find(c =>
    c.number.includes(phoneNumber.value) || phoneNumber.value.includes(c.number)
  ) || null
})

function formatPhoneDisplay(num: string): string {
  if (num.length <= 3) return num
  if (num.length <= 7) return `${num.slice(0, 3)} ${num.slice(3)}`
  return `${num.slice(0, 3)} ${num.slice(3, 7)} ${num.slice(7)}`
}

function pressKey(key: string) {
  if (phoneNumber.value.length < 15) {
    phoneNumber.value += key
  }
}

function backspace() {
  phoneNumber.value = phoneNumber.value.slice(0, -1)
}

function makeCall(callType: 'voice' | 'video') {
  if (!phoneNumber.value) return

  // 查找匹配的联系人
  const contact = matchedContact.value
  const characterId = contact?.characterId || ''
  const name = contact?.name || ''

  // 跳转到通话页面
  const route = callType === 'video' ? '/video-call' : '/voice-call'
  router.push({
    path: route,
    query: {
      number: phoneNumber.value,
      name: name,
      characterId: characterId,
    },
  })
}

function callContact(contact: Contact, callType: 'voice' | 'video') {
  const route = callType === 'video' ? '/video-call' : '/voice-call'
  router.push({
    path: route,
    query: {
      number: contact.number,
      name: contact.name,
      characterId: contact.characterId,
    },
  })
}

function redial(call: any) {
  const route = call.callType === 'video' ? '/video-call' : '/voice-call'
  router.push({
    path: route,
    query: {
      number: call.number,
      name: call.name,
      characterId: call.characterId || '',
    },
  })
}

function clearHistory() {
  if (confirm('确定清除所有通话记录？')) {
    phoneStore.clearCallHistory()
  }
}

onMounted(() => {
  phoneStore.loadCallHistory()
  phoneStore.loadContacts()
})
</script>

<style scoped>
.phone-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

/* Tab Bar */
.tab-bar {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  padding: 0 16px;
}

.tab-item {
  flex: 1;
  padding: 10px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-tertiary);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-item.active {
  color: var(--accent-blue);
  border-bottom-color: var(--accent-blue);
}

/* 拨号盘视图 */
.dial-view {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.number-display {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 16px 8px;
  gap: 8px;
  min-height: 50px;
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

/* 匹配联系人提示 */
.matched-contact {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 4px 16px 8px;
}

.matched-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  overflow: hidden;
}

.matched-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.matched-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.matched-tag {
  font-size: 11px;
  color: var(--accent-blue);
  background: rgba(0, 122, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.dial-pad {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  padding: 8px 40px;
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
  padding: 12px 16px 28px;
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

/* 通讯录 */
.contacts-view {
  flex: 1;
  overflow-y: auto;
}

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

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-color);
}

.contact-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  flex-shrink: 0;
  overflow: hidden;
}

.contact-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.contact-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.contact-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.contact-number {
  font-size: 13px;
  color: var(--text-tertiary);
}

.contact-actions {
  display: flex;
  gap: 8px;
}

.contact-call-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--bg-secondary);
  border: none;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
}

.contact-call-btn:active {
  background: var(--bg-tertiary);
}

/* 通话记录 */
.history-view {
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
  transition: background 0.15s;
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
  gap: 2px;
}

.call-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
}

.call-name.missed {
  color: var(--accent-red);
}

.call-meta {
  font-size: 12px;
  color: var(--text-tertiary);
}

.call-duration {
  font-size: 13px;
  color: var(--text-secondary);
}

.clear-history-btn {
  display: block;
  width: calc(100% - 32px);
  margin: 16px auto;
  padding: 12px;
  background: none;
  border: 1px solid var(--accent-red);
  border-radius: 12px;
  color: var(--accent-red);
  font-size: 14px;
  cursor: pointer;
  text-align: center;
}

.clear-history-btn:active {
  background: rgba(255, 59, 48, 0.1);
}
</style>
