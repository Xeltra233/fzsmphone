<template>
  <div class="zhihu-page">
    <!-- 顶部导航 -->
    <div class="zhihu-header">
      <button class="header-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
      </button>
      <span class="header-title">知乎</span>
      <div class="header-right">
        <button class="header-btn" @click="showPromptEditor = true">
          <svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3" />
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
          </svg>
        </button>
        <button class="header-btn" @click="handleGenerate" :disabled="store.generating">
          <svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ spinning: store.generating }">
            <polyline points="23 4 23 10 17 10" />
            <polyline points="1 20 1 14 7 14" />
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 底部Tab -->
    <div class="bottom-tabs">
      <button v-for="t in tabs" :key="t.key" class="bottom-tab" :class="{ active: activeTab === t.key }" @click="activeTab = t.key">
        <span class="tab-icon">{{ t.icon }}</span>
        <span class="tab-label">{{ t.label }}</span>
      </button>
    </div>

    <!-- 主内容 -->
    <div class="zhihu-body">
      <!-- 推荐/问题列表 -->
      <div v-if="activeTab === 'recommend'" class="tab-content">
        <div v-if="store.lastError" class="error-banner" @click="store.lastError = ''">
          △ {{ store.lastError }}
          <span class="error-dismiss">点击关闭</span>
        </div>

        <div v-if="store.generating && store.zhihuQuestions.length === 0" class="loading-state">
          <div class="loading-spinner"></div>
          <span>AI正在生成知乎内容...</span>
        </div>

        <div v-else-if="store.zhihuQuestions.length === 0" class="empty-state">
          <div class="empty-icon">❓</div>
          <div class="empty-title">暂无问题</div>
          <div class="empty-sub">点击右上角刷新按钮生成AI内容</div>
          <button class="generate-btn" @click="handleGenerate" :disabled="store.generating">
            ✦ 生成知乎内容
          </button>
        </div>

        <div v-else class="question-list">
          <div v-for="q in store.zhihuQuestions" :key="q.id" class="question-card" @click="openQuestion(q)">
            <div class="q-title">{{ q.title }}</div>
            <div v-if="q.description" class="q-desc">{{ truncateText(q.description, 80) }}</div>
            <div class="q-meta">
              <span class="q-answers">{{ q.answers.length }} 个回答</span>
              <span class="q-views">{{ q.views }} 浏览</span>
              <span class="q-followers">{{ q.followers }} 关注</span>
            </div>
            <!-- 首条回答预览 -->
            <div v-if="q.answers.length > 0" class="q-preview">
              <div class="preview-author">
                <div class="preview-avatar" :style="{ background: getAvatarColor(q.answers[0].author) }">{{ q.answers[0].author.charAt(0) }}</div>
                <span class="preview-name">{{ q.answers[0].author }}</span>
              </div>
              <div class="preview-text">{{ truncateText(q.answers[0].content, 100) }}</div>
              <div v-if="getItemImages(q.answers[0]).length" class="preview-images">
                <img v-for="(img, idx) in getItemImages(q.answers[0])" :key="idx" :src="img" class="preview-gen-img" alt="" />
              </div>
              <div class="preview-stats">
                <span>▲ {{ q.answers[0].upvotes }}</span>
              </div>
            </div>
            <button class="q-delete" @click.stop="store.deleteZhihuQuestion(q.id)">✕</button>
          </div>
        </div>
      </div>

      <!-- 问题详情 -->
      <div v-if="activeTab === 'detail' && selectedQuestion" class="tab-content">
        <button class="back-to-list" @click="activeTab = 'recommend'">← 返回列表</button>
        <div class="detail-header">
          <div class="detail-title">{{ selectedQuestion.title }}</div>
          <div v-if="selectedQuestion.description" class="detail-desc">{{ selectedQuestion.description }}</div>
          <div class="detail-stats">
            <span>{{ selectedQuestion.followers }} 人关注</span>
            <span>{{ selectedQuestion.views }} 次浏览</span>
            <span>{{ selectedQuestion.answers.length }} 个回答</span>
          </div>
        </div>

        <div class="answers-list">
          <div v-for="a in selectedQuestion.answers" :key="a.id" class="answer-card">
            <div class="ans-user-row">
              <div class="ans-avatar" :style="{ background: getAvatarColor(a.author) }">{{ a.author.charAt(0) }}</div>
              <div class="ans-info">
                <span class="ans-name">{{ a.author }}</span>
                <span class="ans-time">{{ formatRelativeTime(a.timestamp) }}</span>
              </div>
            </div>
            <div class="ans-content">{{ a.content }}</div>
            <div v-if="getItemImages(a).length" class="ans-images">
              <div v-for="(img, idx) in getItemImages(a)" :key="idx" class="img-wrapper">
                <img :src="img" class="ans-gen-img" alt="" />
                <button v-if="(a as any).imagePrompt" class="regen-btn" :disabled="store.regeneratingImages.has(`${a.id}-${idx}`)" @click="store.regenerateImage('zhihu', a.id, idx)">
                  <span v-if="store.regeneratingImages.has(`${a.id}-${idx}`)" class="regen-spin"></span>
                  <svg v-else viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10" /><polyline points="1 20 1 14 7 14" /><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" /></svg>
                </button>
              </div>
            </div>
            <div class="ans-actions">
              <button class="vote-btn" :class="{ voted: a.isUpvoted }" @click="store.toggleZhihuVote(selectedQuestion!.id, a.id, 'up')">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" :stroke="a.isUpvoted ? '#0066ff' : 'currentColor'" stroke-width="2">
                  <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3H14z" />
                  <path d="M7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3" />
                </svg>
                {{ a.upvotes }}
              </button>
              <button class="vote-btn down" :class="{ voted: a.isDownvoted }" @click="store.toggleZhihuVote(selectedQuestion!.id, a.id, 'down')">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" :stroke="a.isDownvoted ? '#ff4757' : 'currentColor'" stroke-width="2">
                  <path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3H10z" />
                  <path d="M17 2h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3" />
                </svg>
              </button>
              <button class="ans-comment-btn" @click="toggleAnswerComments(a.id)">
                ◌ {{ a.comments?.length || 0 }}
              </button>
            </div>

            <!-- 评论区 -->
            <div v-if="expandedAnswers[a.id]" class="ans-comments">
              <div v-for="c in a.comments" :key="c.id" class="comment-item">
                <span class="comment-author" :style="{ color: getAvatarColor(c.author) }">{{ c.author }}</span>
                <span class="comment-sep">：</span>
                <span class="comment-text">{{ c.content }}</span>
              </div>
              <div v-if="!a.comments?.length" class="no-comments">暂无评论</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 我的 -->
      <div v-if="activeTab === 'me'" class="tab-content">
        <div class="profile-card">
          <div class="profile-avatar" :style="{ background: getAvatarColor(userName) }">{{ userName.charAt(0) }}</div>
          <div class="profile-name">{{ userName }}</div>
          <div class="profile-stats">
            <div class="stat-item">
              <span class="stat-num">{{ store.zhihuQuestions.length }}</span>
              <span class="stat-label">关注问题</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 提示词编辑器 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="showPromptEditor" class="compose-overlay" @click.self="showPromptEditor = false">
          <div class="compose-sheet prompt-sheet">
            <div class="sheet-bar">
              <button class="sheet-cancel" @click="showPromptEditor = false">关闭</button>
              <span class="sheet-title">知乎提示词</span>
              <button class="sheet-post" @click="savePrompt">保存</button>
            </div>
            <textarea v-model="promptText" class="sheet-textarea prompt-textarea" placeholder="输入知乎生成提示词..."></textarea>
            <div class="prompt-actions">
              <button class="prompt-action-btn" @click="resetPrompt">重置默认</button>
              <button class="prompt-action-btn danger" @click="clearAllData">清空数据</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { getScopedItem } from '@/utils/userScopedStorage'
import { useSocialAIStore } from '@/stores/socialAI'
import { getAvatarColor, formatRelativeTime, truncateText } from '@/utils/socialParsers'
import type { ZhihuQuestion } from '@/utils/socialParsers'
import { getPromptTemplate, setPromptTemplate, resetPromptTemplate } from '@/utils/socialPrompts'

const store = useSocialAIStore()

const activeTab = ref<'recommend' | 'detail' | 'me'>('recommend')
const tabs = [
  { key: 'recommend' as const, icon: '⌂', label: '推荐' },
  { key: 'me' as const, icon: '○', label: '我的' },
]

const selectedQuestion = ref<ZhihuQuestion | null>(null)
const showPromptEditor = ref(false)
const expandedAnswers = reactive<Record<string, boolean>>({})
const promptText = ref('')

const userName = computed(() => {
  try {
    const currentId = localStorage.getItem('currentPersonaId')
    if (currentId) {
      const charsStr = getScopedItem('characters')
      if (charsStr) {
        const chars = JSON.parse(charsStr)
        const user = chars.find((c: any) => c.id === currentId && c.type === 'user')
        if (user?.name) return user.name
      }
    }
  } catch { /* ignore */ }
  return '我'
})

onMounted(() => {
  store.loadData('zhihu')
  promptText.value = getPromptTemplate('zhihu')
})

async function handleGenerate() {
  await store.generateZhihuContent()
}

function openQuestion(q: ZhihuQuestion) {
  selectedQuestion.value = q
  activeTab.value = 'detail'
}

function toggleAnswerComments(answerId: string) {
  expandedAnswers[answerId] = !expandedAnswers[answerId]
}

function savePrompt() {
  setPromptTemplate('zhihu', promptText.value)
  showPromptEditor.value = false
}

function resetPrompt() {
  resetPromptTemplate('zhihu')
  promptText.value = getPromptTemplate('zhihu')
}

function clearAllData() {
  if (confirm('确定要清空所有知乎数据吗？')) {
    store.clearData('zhihu')
    showPromptEditor.value = false
  }
}

function getItemImages(item: any): string[] {
  return Array.isArray(item?.images) ? item.images.filter((s: string) => !!s) : []
}
</script>

<style scoped>
.zhihu-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f6f6f6;
  position: relative;
}

.zhihu-header {
  --status-bar-height: 54px;
  height: calc(var(--status-bar-height) + 44px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--status-bar-height) 14px 0;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.header-btn { width: 34px; height: 34px; border: none; background: none; color: #666; cursor: pointer; display: flex; align-items: center; justify-content: center; border-radius: 50%; }
.header-btn:active { background: #f5f5f5; }
.header-btn:disabled { opacity: 0.4; }
.header-right { display: flex; gap: 4px; }
.header-title { font-size: 18px; font-weight: 700; color: #0066ff; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.bottom-tabs { display: flex; background: #fff; border-top: 1px solid #f0f0f0; flex-shrink: 0; order: 2; position: absolute; bottom: 0; left: 0; right: 0; z-index: 10; padding-bottom: env(safe-area-inset-bottom); }
.bottom-tab { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 6px 0; border: none; background: none; cursor: pointer; color: #999; transition: color 0.2s; }
.bottom-tab.active { color: #0066ff; }
.tab-icon { font-size: 20px; }
.tab-label { font-size: 10px; font-weight: 500; }

.zhihu-body { flex: 1; overflow-y: auto; padding-bottom: 60px; -webkit-overflow-scrolling: touch; }
.zhihu-body::-webkit-scrollbar { display: none; }
.tab-content { min-height: 100%; }

.error-banner { background: #fff3cd; color: #856404; padding: 10px 16px; font-size: 13px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; }
.error-dismiss { font-size: 11px; color: #a08040; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 80px 20px; gap: 12px; color: #999; }
.loading-spinner { width: 32px; height: 32px; border: 3px solid #f0f0f0; border-top-color: #0066ff; border-radius: 50%; animation: spin 0.8s linear infinite; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; gap: 8px; }
.empty-icon { font-size: 44px; margin-bottom: 4px; }
.empty-title { font-size: 16px; font-weight: 600; color: #666; }
.empty-sub { font-size: 13px; color: #aaa; }
.generate-btn { margin-top: 16px; padding: 10px 28px; background: linear-gradient(135deg, #0066ff, #3399ff); color: #fff; border: none; border-radius: 22px; font-size: 15px; font-weight: 600; cursor: pointer; }
.generate-btn:disabled { opacity: 0.5; }

/* 问题列表 */
.question-list { padding: 0; }
.question-card { background: #fff; padding: 16px; border-bottom: 1px solid #f0f0f0; position: relative; cursor: pointer; transition: background 0.15s; }
.question-card:active { background: #fafafa; }
.q-title { font-size: 16px; font-weight: 600; color: #1a1a1a; line-height: 1.5; margin-bottom: 6px; }
.q-desc { font-size: 14px; color: #666; line-height: 1.5; margin-bottom: 8px; }
.q-meta { display: flex; gap: 16px; font-size: 12px; color: #999; margin-bottom: 10px; }
.q-delete { position: absolute; top: 12px; right: 12px; width: 24px; height: 24px; border: none; background: none; color: #ccc; cursor: pointer; font-size: 14px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.q-delete:active { background: #f5f5f5; color: #999; }

/* 回答预览 */
.q-preview { background: #f7f8fa; border-radius: 8px; padding: 12px; }
.preview-author { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.preview-avatar { width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 12px; font-weight: 700; flex-shrink: 0; }
.preview-name { font-size: 13px; font-weight: 600; color: #333; }
.preview-text { font-size: 14px; color: #444; line-height: 1.6; }
.preview-stats { font-size: 12px; color: #999; margin-top: 6px; }

/* 问题详情 */
.back-to-list { border: none; background: none; color: #0066ff; font-size: 14px; padding: 12px 16px; cursor: pointer; display: block; }
.detail-header { background: #fff; padding: 16px; border-bottom: 8px solid #f6f6f6; }
.detail-title { font-size: 20px; font-weight: 700; color: #1a1a1a; line-height: 1.4; margin-bottom: 8px; }
.detail-desc { font-size: 14px; color: #666; line-height: 1.6; margin-bottom: 10px; }
.detail-stats { display: flex; gap: 16px; font-size: 13px; color: #999; }

/* 回答列表 */
.answers-list { padding: 0; }
.answer-card { background: #fff; padding: 16px; border-bottom: 8px solid #f6f6f6; }
.ans-user-row { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
.ans-avatar { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 14px; font-weight: 700; flex-shrink: 0; }
.ans-info { display: flex; flex-direction: column; }
.ans-name { font-size: 15px; font-weight: 600; color: #333; }
.ans-time { font-size: 12px; color: #bbb; }
.ans-content { font-size: 15px; color: #333; line-height: 1.8; word-break: break-word; margin-bottom: 12px; white-space: pre-wrap; }
.ans-images { margin-bottom: 12px; }
.ans-gen-img { width: 100%; border-radius: 8px; margin-bottom: 4px; }
.preview-images { margin-top: 6px; }
.preview-gen-img { width: 100%; border-radius: 6px; margin-bottom: 4px; }

.ans-actions { display: flex; align-items: center; gap: 16px; padding-top: 8px; border-top: 1px solid #f5f5f5; }
.vote-btn { display: flex; align-items: center; gap: 4px; border: 1px solid #e8e8e8; background: #f7f8fa; color: #8590a6; font-size: 13px; cursor: pointer; padding: 6px 14px; border-radius: 16px; transition: all 0.15s; }
.vote-btn.voted { background: #eef3ff; border-color: #0066ff; color: #0066ff; }
.vote-btn.down { padding: 6px 10px; }
.vote-btn.down.voted { background: #fff0f0; border-color: #ff4757; color: #ff4757; }
.ans-comment-btn { border: none; background: none; color: #8590a6; font-size: 13px; cursor: pointer; padding: 6px 10px; }

.ans-comments { background: #f7f8fa; border-radius: 8px; padding: 10px 12px; margin-top: 10px; }
.comment-item { padding: 6px 0; border-bottom: 1px solid #f0f0f0; font-size: 14px; line-height: 1.6; }
.comment-item:last-child { border-bottom: none; }
.comment-author { font-weight: 600; font-size: 13px; }
.comment-sep { color: #999; }
.comment-text { color: #333; }
.no-comments { text-align: center; font-size: 13px; color: #ccc; padding: 8px 0; }

/* 个人 */
.profile-card { background: #fff; padding: 24px; text-align: center; border-bottom: 8px solid #f6f6f6; }
.profile-avatar { width: 60px; height: 60px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 24px; font-weight: 700; margin: 0 auto 10px; }
.profile-name { font-size: 18px; font-weight: 700; color: #333; margin-bottom: 12px; }
.profile-stats { display: flex; justify-content: center; gap: 32px; }
.stat-item { display: flex; flex-direction: column; align-items: center; gap: 2px; }
.stat-num { font-size: 18px; font-weight: 700; color: #333; }
.stat-label { font-size: 12px; color: #999; }

/* 弹窗 */
.compose-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.4); z-index: 100; display: flex; align-items: flex-end; }
.compose-sheet { width: 100%; background: #fff; border-radius: 12px 12px 0 0; max-height: 80%; display: flex; flex-direction: column; }
.prompt-sheet { max-height: 90%; }
.sheet-bar { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #f0f0f0; flex-shrink: 0; }
.sheet-cancel { border: none; background: none; color: #666; font-size: 15px; cursor: pointer; }
.sheet-title { font-size: 16px; font-weight: 600; color: #333; }
.sheet-post { border: none; background: none; color: #0066ff; font-size: 15px; font-weight: 600; cursor: pointer; }
.sheet-textarea { flex: 1; padding: 14px 16px; border: none; resize: none; font-size: 15px; color: #333; line-height: 1.6; outline: none; min-height: 120px; }
.sheet-textarea::placeholder { color: #ccc; }
.prompt-textarea { min-height: 200px; font-family: monospace; font-size: 13px; }
.prompt-actions { display: flex; gap: 8px; padding: 10px 16px; border-top: 1px solid #f0f0f0; }
.prompt-action-btn { flex: 1; padding: 10px; border: 1px solid #e8e8e8; background: #f7f8fa; border-radius: 8px; font-size: 13px; color: #666; cursor: pointer; }
.prompt-action-btn.danger { color: #ff4757; border-color: #ffcccc; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; }
.slide-up-enter-from .compose-sheet, .slide-up-leave-to .compose-sheet { transform: translateY(100%); }
.img-wrapper { position: relative; display: inline-block; width: 100%; }
.regen-btn { position: absolute; top: 6px; right: 6px; width: 28px; height: 28px; border-radius: 50%; background: rgba(0,0,0,0.5); border: none; color: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center; opacity: 0; transition: opacity 0.2s; z-index: 2; }
.img-wrapper:hover .regen-btn { opacity: 1; }
.regen-btn:disabled { cursor: wait; opacity: 1 !important; }
.regen-spin { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: spin 0.8s linear infinite; display: block; }
</style>
