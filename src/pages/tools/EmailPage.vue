<template>
  <div class="email-page">
    <NavBar title="邮箱" back>
      <template #right>
        <button class="icon-btn" @click="showCompose = true" title="写邮件">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
      </template>
    </NavBar>

    <!-- Generate Button -->
    <div v-if="store.emailItems.length === 0 && !store.emailLoading" class="empty-state">
      <div class="empty-icon"><svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="10" width="40" height="28" rx="4"/><polyline points="4,14 24,30 44,14"/></svg></div>
      <div class="empty-title">收件箱为空</div>
      <div class="empty-sub">点击下方按钮让AI生成邮件</div>
      <button class="gen-btn" :disabled="store.generating" @click="store.generateEmailContent()">
        {{ store.generating ? '生成中...' : '◆ AI 生成邮件' }}
      </button>
    </div>

    <!-- Loading -->
    <div v-else-if="store.emailLoading" class="loading-wrap">
      <div class="loading-spinner"></div>
      <span>正在生成邮件...</span>
    </div>

    <!-- Inbox List -->
    <div v-else-if="!selectedEmail" class="inbox-list">
      <div
        v-for="mail in store.emailItems"
        :key="mail.id"
        class="mail-item"
        :class="{ unread: !mail.isRead }"
        @click="openEmail(mail)"
      >
        <div class="mail-avatar">{{ mail.senderName[0] }}</div>
        <div class="mail-info">
          <div class="mail-top">
            <span class="mail-sender">{{ mail.senderName }}</span>
            <span class="mail-time">{{ mail.timestamp }}</span>
          </div>
          <div class="mail-subject">{{ mail.subject }}</div>
          <div class="mail-preview">{{ mail.preview }}</div>
        </div>
        <div v-if="mail.isStarred" class="mail-star">★</div>
      </div>
      <button class="gen-btn refresh-btn" :disabled="store.generating" @click="store.generateEmailContent()">
        {{ store.generating ? '生成中...' : '◆ 生成更多邮件' }}
      </button>
    </div>

    <!-- Email Detail -->
    <div v-else class="email-detail">
      <button class="back-link" @click="selectedEmail = null">← 返回收件箱</button>
      <div class="detail-header">
        <h2 class="detail-subject">{{ selectedEmail.subject }}</h2>
        <div class="detail-meta">
          <span class="detail-sender">{{ selectedEmail.senderName }} &lt;{{ selectedEmail.senderEmail }}&gt;</span>
          <span class="detail-time">{{ selectedEmail.timestamp }}</span>
        </div>
      </div>
      <div class="detail-body">{{ selectedEmail.content }}</div>
      <div class="detail-actions">
        <button class="action-btn" @click="toggleStar(selectedEmail)">{{ selectedEmail.isStarred ? '★ 已标星' : '☆ 标星' }}</button>
        <button class="action-btn" @click="deleteEmail(selectedEmail)">删除</button>
      </div>
    </div>

    <!-- Compose Modal -->
    <Teleport to="#phone-overlay">
      <div v-if="showCompose" class="modal-overlay" @click.self="showCompose = false">
        <div class="compose-modal">
          <div class="compose-header">
            <span>新邮件</span>
            <button class="close-btn" @click="showCompose = false">✕</button>
          </div>
          <input v-model="composeTo" class="compose-input" placeholder="收件人" />
          <input v-model="composeSubject" class="compose-input" placeholder="主题" />
          <textarea v-model="composeBody" class="compose-textarea" placeholder="正文..."></textarea>
          <button class="gen-btn" @click="sendCompose">发送</button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'
import type { EmailItem } from '@/utils/socialParsers'

const store = useSocialAIStore()
const selectedEmail = ref<EmailItem | null>(null)
const showCompose = ref(false)
const composeTo = ref('')
const composeSubject = ref('')
const composeBody = ref('')

onMounted(() => { store.loadData('email') })

function openEmail(mail: EmailItem) {
  mail.isRead = true
  store.saveData('email')
  selectedEmail.value = mail
}

function toggleStar(mail: EmailItem) {
  mail.isStarred = !mail.isStarred
  store.saveData('email')
}

function deleteEmail(mail: EmailItem) {
  const idx = store.emailItems.indexOf(mail)
  if (idx >= 0) store.emailItems.splice(idx, 1)
  store.saveData('email')
  selectedEmail.value = null
}

function sendCompose() {
  showCompose.value = false
  composeTo.value = ''
  composeSubject.value = ''
  composeBody.value = ''
}
</script>

<style scoped>
.email-page { height: 100%; display: flex; flex-direction: column; background: var(--bg-primary); }
.empty-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; padding: 40px 20px; }
.empty-icon { width: 56px; height: 56px; color: var(--text-tertiary); }
.empty-icon svg { width: 100%; height: 100%; }
.empty-title { font-size: 18px; font-weight: 700; color: var(--text-primary); }
.empty-sub { font-size: 13px; color: var(--text-tertiary); }
.gen-btn { padding: 10px 24px; border: none; border-radius: 20px; background: var(--brand-primary); color: #fff; font-size: 14px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.gen-btn:disabled { opacity: 0.5; }
.loading-wrap { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; color: var(--text-secondary); }
.loading-spinner { width: 32px; height: 32px; border: 3px solid var(--separator); border-top-color: var(--brand-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.inbox-list { flex: 1; overflow-y: auto; padding: 0 0 80px; }
.mail-item { display: flex; align-items: flex-start; gap: 12px; padding: 14px 16px; border-bottom: 0.5px solid var(--separator); cursor: pointer; }
.mail-item.unread { background: var(--fill-tertiary); }
.mail-avatar { width: 40px; height: 40px; border-radius: 20px; background: var(--brand-primary); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 700; flex-shrink: 0; }
.mail-info { flex: 1; min-width: 0; }
.mail-top { display: flex; justify-content: space-between; align-items: center; }
.mail-sender { font-size: 15px; font-weight: 600; color: var(--text-primary); }
.mail-time { font-size: 12px; color: var(--text-tertiary); }
.mail-subject { font-size: 14px; font-weight: 500; color: var(--text-primary); margin-top: 2px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.mail-preview { font-size: 13px; color: var(--text-tertiary); margin-top: 2px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.mail-star { color: #FFD700; font-size: 16px; }
.refresh-btn { margin: 16px auto; display: block; }

.email-detail { flex: 1; overflow-y: auto; padding: 0 16px 80px; }
.back-link { background: none; border: none; color: var(--brand-primary); font-size: 14px; padding: 12px 0; cursor: pointer; }
.detail-header { border-bottom: 0.5px solid var(--separator); padding-bottom: 12px; }
.detail-subject { font-size: 20px; font-weight: 700; color: var(--text-primary); margin: 0; }
.detail-meta { margin-top: 8px; }
.detail-sender { font-size: 14px; color: var(--text-secondary); display: block; }
.detail-time { font-size: 12px; color: var(--text-tertiary); }
.detail-body { padding: 16px 0; font-size: 15px; color: var(--text-primary); line-height: 1.6; white-space: pre-wrap; }
.detail-actions { display: flex; gap: 12px; padding: 12px 0; }
.action-btn { padding: 8px 16px; border: 1px solid var(--separator); border-radius: 16px; background: var(--bg-secondary); color: var(--text-primary); font-size: 13px; cursor: pointer; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.compose-modal { width: 90%; max-width: 400px; background: var(--bg-primary); border-radius: 16px; padding: 20px; display: flex; flex-direction: column; gap: 10px; }
.compose-header { display: flex; justify-content: space-between; align-items: center; font-size: 16px; font-weight: 700; color: var(--text-primary); }
.close-btn { background: none; border: none; font-size: 18px; color: var(--text-tertiary); cursor: pointer; }
.compose-input { border: 1px solid var(--separator); border-radius: 10px; padding: 10px 12px; font-size: 14px; background: var(--bg-secondary); color: var(--text-primary); outline: none; }
.compose-textarea { border: 1px solid var(--separator); border-radius: 10px; padding: 10px 12px; font-size: 14px; background: var(--bg-secondary); color: var(--text-primary); outline: none; min-height: 120px; resize: none; }
</style>
