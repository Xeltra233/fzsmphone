<template>
  <div class="xhs-page">
    <!-- 顶部导航 -->
    <div class="xhs-header">
      <button class="header-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
      </button>
      <div class="header-tabs">
        <button class="htab" :class="{ active: activeTab === 'discover' }" @click="activeTab = 'discover'">发现</button>
        <button class="htab" :class="{ active: activeTab === 'me' }" @click="activeTab = 'me'">我</button>
      </div>
      <div class="header-right">
        <button class="header-btn" @click="showPromptEditor = true">⚙</button>
        <button class="header-btn" @click="handleGenerate" :disabled="store.generating">
          <svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ spinning: store.generating }">
            <polyline points="23 4 23 10 17 10" />
            <polyline points="1 20 1 14 7 14" />
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 主内容 -->
    <div class="xhs-body">
      <!-- 发现/瀑布流 -->
      <div v-if="activeTab === 'discover'" class="tab-content">
        <div v-if="store.lastError" class="error-banner" @click="store.lastError = ''">
          △ {{ store.lastError }}
          <span class="error-dismiss">点击关闭</span>
        </div>

        <div v-if="store.generating && store.xhsNotes.length === 0" class="loading-state">
          <div class="loading-spinner"></div>
          <span>AI正在生成笔记...</span>
        </div>

        <div v-else-if="store.xhsNotes.length === 0" class="empty-state">
          <div class="empty-icon">▤</div>
          <div class="empty-title">暂无笔记</div>
          <div class="empty-sub">点击右上角刷新按钮生成AI内容</div>
          <button class="generate-btn" @click="handleGenerate" :disabled="store.generating">
            ✦ 生成小红书内容
          </button>
        </div>

        <!-- 瀑布流双列 -->
        <div v-else class="waterfall">
          <div class="waterfall-col">
            <div v-for="note in leftCol" :key="note.id" class="note-card" @click="openNote(note)">
              <div class="note-cover" :style="{ background: getNoteBg(note) }">
                <div class="cover-text">{{ note.title || note.content.slice(0, 20) }}</div>
              </div>
              <div class="note-body">
                <div class="note-title">{{ note.title }}</div>
                <div class="note-footer">
                  <div class="note-author">
                    <div class="author-dot" :style="{ background: getAvatarColor(note.author) }"></div>
                    <span>{{ note.author }}</span>
                  </div>
                  <button class="note-like" :class="{ liked: note.isLiked }" @click.stop="store.toggleXhsLike(note.id)">
                    ♥ {{ note.likes }}
                  </button>
                </div>
              </div>
              <button class="note-delete" @click.stop="store.deleteXhsNote(note.id)">✕</button>
            </div>
          </div>
          <div class="waterfall-col">
            <div v-for="note in rightCol" :key="note.id" class="note-card" @click="openNote(note)">
              <div class="note-cover" :style="{ background: getNoteBg(note) }">
                <div class="cover-text">{{ note.title || note.content.slice(0, 20) }}</div>
              </div>
              <div class="note-body">
                <div class="note-title">{{ note.title }}</div>
                <div class="note-footer">
                  <div class="note-author">
                    <div class="author-dot" :style="{ background: getAvatarColor(note.author) }"></div>
                    <span>{{ note.author }}</span>
                  </div>
                  <button class="note-like" :class="{ liked: note.isLiked }" @click.stop="store.toggleXhsLike(note.id)">
                    ♥ {{ note.likes }}
                  </button>
                </div>
              </div>
              <button class="note-delete" @click.stop="store.deleteXhsNote(note.id)">✕</button>
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
            <div class="stat-item"><span class="stat-num">{{ store.xhsNotes.length }}</span><span class="stat-label">笔记</span></div>
            <div class="stat-item"><span class="stat-num">{{ totalLikes }}</span><span class="stat-label">获赞</span></div>
            <div class="stat-item"><span class="stat-num">{{ totalCollects }}</span><span class="stat-label">收藏</span></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 笔记详情弹窗 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="selectedNote" class="detail-overlay" @click.self="selectedNote = null">
          <div class="detail-sheet">
            <div class="sheet-bar">
              <button class="sheet-cancel" @click="selectedNote = null">关闭</button>
              <span class="sheet-title">笔记详情</span>
              <span></span>
            </div>
            <div class="detail-body">
              <div class="detail-author-row">
                <div class="detail-avatar" :style="{ background: getAvatarColor(selectedNote.author) }">{{ selectedNote.author.charAt(0) }}</div>
                <span class="detail-author-name">{{ selectedNote.author }}</span>
                <span class="detail-time">{{ formatRelativeTime(selectedNote.timestamp) }}</span>
              </div>
              <div class="detail-title">{{ selectedNote.title }}</div>
              <div class="detail-content">{{ selectedNote.content }}</div>
              <div v-if="selectedNote.tags.length" class="detail-tags">
                <span v-for="tag in selectedNote.tags" :key="tag" class="tag">#{{ tag }}</span>
              </div>
              <div class="detail-actions">
                <button class="action-btn" :class="{ active: selectedNote.isLiked }" @click="store.toggleXhsLike(selectedNote.id)">
                  ♥ {{ selectedNote.likes }}
                </button>
                <button class="action-btn" :class="{ active: selectedNote.isCollected }" @click="store.toggleXhsCollect(selectedNote.id)">
                  ★ {{ selectedNote.collects }}
                </button>
                <span class="action-btn">◌ {{ selectedNote.comments?.length || 0 }}</span>
              </div>

              <!-- 评论 -->
              <div class="detail-comments">
                <div v-for="c in selectedNote.comments" :key="c.id" class="comment-item">
                  <span class="comment-author" :style="{ color: getAvatarColor(c.author) }">{{ c.author }}</span>
                  <span v-if="c.replyTo" class="reply-to"> 回复 <span :style="{ color: getAvatarColor(c.replyTo) }">{{ c.replyTo }}</span></span>
                  <span class="comment-sep">：</span>
                  <span class="comment-text">{{ c.content }}</span>
                </div>
                <div v-if="!selectedNote.comments?.length" class="no-comments">暂无评论</div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 提示词编辑器 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="showPromptEditor" class="detail-overlay" @click.self="showPromptEditor = false">
          <div class="detail-sheet prompt-sheet">
            <div class="sheet-bar">
              <button class="sheet-cancel" @click="showPromptEditor = false">关闭</button>
              <span class="sheet-title">小红书提示词</span>
              <button class="sheet-post" @click="savePrompt">保存</button>
            </div>
            <textarea v-model="promptText" class="sheet-textarea" placeholder="输入小红书生成提示词..."></textarea>
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
import { ref, computed, onMounted } from 'vue'
import { useSocialAIStore } from '@/stores/socialAI'
import { getAvatarColor, formatRelativeTime } from '@/utils/socialParsers'
import type { XhsNote } from '@/utils/socialParsers'
import { getPromptTemplate, setPromptTemplate, resetPromptTemplate } from '@/utils/socialPrompts'

const store = useSocialAIStore()

const activeTab = ref<'discover' | 'me'>('discover')
const selectedNote = ref<XhsNote | null>(null)
const showPromptEditor = ref(false)
const promptText = ref('')

const NOTE_COLORS = [
  'linear-gradient(135deg, #ffecd2, #fcb69f)',
  'linear-gradient(135deg, #a1c4fd, #c2e9fb)',
  'linear-gradient(135deg, #d4fc79, #96e6a1)',
  'linear-gradient(135deg, #fbc2eb, #a6c1ee)',
  'linear-gradient(135deg, #fdcbf1, #e6dee9)',
  'linear-gradient(135deg, #f6d365, #fda085)',
  'linear-gradient(135deg, #89f7fe, #66a6ff)',
  'linear-gradient(135deg, #fddb92, #d1fdff)',
  'linear-gradient(135deg, #c1dfc4, #deecdd)',
  'linear-gradient(135deg, #ffeaa7, #dfe6e9)',
]

function getNoteBg(note: XhsNote): string {
  let hash = 0
  for (let i = 0; i < note.id.length; i++) hash = hash * 31 + note.id.charCodeAt(i)
  return NOTE_COLORS[Math.abs(hash) % NOTE_COLORS.length]
}

const leftCol = computed(() => store.xhsNotes.filter((_, i) => i % 2 === 0))
const rightCol = computed(() => store.xhsNotes.filter((_, i) => i % 2 === 1))

const userName = computed(() => {
  try {
    const currentId = localStorage.getItem('currentPersonaId')
    if (currentId) {
      const charsStr = localStorage.getItem('characters')
      if (charsStr) {
        const chars = JSON.parse(charsStr)
        const user = chars.find((c: any) => c.id === currentId && c.type === 'user')
        if (user?.name) return user.name
      }
    }
  } catch { /* ignore */ }
  return '我'
})

const totalLikes = computed(() => store.xhsNotes.reduce((s, n) => s + n.likes, 0))
const totalCollects = computed(() => store.xhsNotes.reduce((s, n) => s + n.collects, 0))

onMounted(() => {
  store.loadData('xiaohongshu')
  promptText.value = getPromptTemplate('xiaohongshu')
})

async function handleGenerate() {
  await store.generateXhsContent()
}

function openNote(note: XhsNote) {
  selectedNote.value = note
}

function savePrompt() {
  setPromptTemplate('xiaohongshu', promptText.value)
  showPromptEditor.value = false
}

function resetPrompt() {
  resetPromptTemplate('xiaohongshu')
  promptText.value = getPromptTemplate('xiaohongshu')
}

function clearAllData() {
  if (confirm('确定要清空所有小红书数据吗？')) {
    store.clearData('xiaohongshu')
    showPromptEditor.value = false
  }
}
</script>

<style scoped>
.xhs-page { height: 100%; display: flex; flex-direction: column; background: #fff; position: relative; }

.xhs-header {
  --status-bar-height: 54px;
  height: calc(var(--status-bar-height) + 44px);
  display: flex; align-items: center; justify-content: space-between;
  padding: var(--status-bar-height) 10px 0;
  background: #fff; border-bottom: 1px solid #f0f0f0; flex-shrink: 0;
}
.header-btn { width: 34px; height: 34px; border: none; background: none; color: #666; cursor: pointer; display: flex; align-items: center; justify-content: center; border-radius: 50%; font-size: 16px; }
.header-btn:active { background: #f5f5f5; }
.header-btn:disabled { opacity: 0.4; }
.header-right { display: flex; gap: 2px; }
.header-tabs { display: flex; gap: 0; }
.htab { border: none; background: none; font-size: 16px; font-weight: 600; color: #999; cursor: pointer; padding: 4px 12px; position: relative; }
.htab.active { color: #333; }
.htab.active::after { content: ''; position: absolute; bottom: -2px; left: 50%; transform: translateX(-50%); width: 20px; height: 3px; background: #ff2442; border-radius: 2px; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.xhs-body { flex: 1; overflow-y: auto; -webkit-overflow-scrolling: touch; }
.xhs-body::-webkit-scrollbar { display: none; }
.tab-content { min-height: 100%; }

.error-banner { background: #fff3cd; color: #856404; padding: 10px 16px; font-size: 13px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; }
.error-dismiss { font-size: 11px; color: #a08040; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 80px 20px; gap: 12px; color: #999; }
.loading-spinner { width: 32px; height: 32px; border: 3px solid #f0f0f0; border-top-color: #ff2442; border-radius: 50%; animation: spin 0.8s linear infinite; }

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; gap: 8px; }
.empty-icon { font-size: 44px; margin-bottom: 4px; }
.empty-title { font-size: 16px; font-weight: 600; color: #666; }
.empty-sub { font-size: 13px; color: #aaa; }
.generate-btn { margin-top: 16px; padding: 10px 28px; background: linear-gradient(135deg, #ff2442, #ff6b81); color: #fff; border: none; border-radius: 22px; font-size: 15px; font-weight: 600; cursor: pointer; }
.generate-btn:disabled { opacity: 0.5; }

/* 瀑布流 */
.waterfall { display: flex; gap: 6px; padding: 6px; }
.waterfall-col { flex: 1; display: flex; flex-direction: column; gap: 6px; }
.note-card { background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,0.06); position: relative; cursor: pointer; transition: transform 0.15s; }
.note-card:active { transform: scale(0.98); }
.note-cover { height: 120px; display: flex; align-items: flex-end; padding: 10px; position: relative; }
.cover-text { color: #fff; font-size: 13px; font-weight: 600; text-shadow: 0 1px 4px rgba(0,0,0,0.3); line-height: 1.4; word-break: break-word; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.note-body { padding: 8px 10px 10px; }
.note-title { font-size: 13px; font-weight: 600; color: #333; line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; margin-bottom: 6px; }
.note-footer { display: flex; justify-content: space-between; align-items: center; }
.note-author { display: flex; align-items: center; gap: 4px; font-size: 11px; color: #999; }
.author-dot { width: 14px; height: 14px; border-radius: 50%; flex-shrink: 0; }
.note-like { border: none; background: none; font-size: 12px; color: #999; cursor: pointer; }
.note-like.liked { color: #ff2442; }
.note-delete { position: absolute; top: 6px; right: 6px; width: 20px; height: 20px; border: none; background: rgba(0,0,0,0.3); color: #fff; cursor: pointer; font-size: 11px; border-radius: 50%; display: flex; align-items: center; justify-content: center; opacity: 0; transition: opacity 0.2s; }
.note-card:hover .note-delete { opacity: 1; }

/* 详情弹窗 */
.detail-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.4); z-index: 100; display: flex; align-items: flex-end; }
.detail-sheet { width: 100%; background: #fff; border-radius: 12px 12px 0 0; max-height: 85%; display: flex; flex-direction: column; }
.prompt-sheet { max-height: 90%; }
.sheet-bar { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #f0f0f0; flex-shrink: 0; }
.sheet-cancel { border: none; background: none; color: #666; font-size: 15px; cursor: pointer; }
.sheet-title { font-size: 16px; font-weight: 600; color: #333; }
.sheet-post { border: none; background: none; color: #ff2442; font-size: 15px; font-weight: 600; cursor: pointer; }
.sheet-textarea { flex: 1; padding: 14px 16px; border: none; resize: none; font-size: 13px; color: #333; line-height: 1.6; outline: none; min-height: 200px; font-family: monospace; }
.sheet-textarea::placeholder { color: #ccc; }

.detail-body { flex: 1; overflow-y: auto; padding: 16px; }
.detail-author-row { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }
.detail-avatar { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 14px; font-weight: 700; flex-shrink: 0; }
.detail-author-name { font-size: 15px; font-weight: 600; color: #333; }
.detail-time { font-size: 12px; color: #bbb; margin-left: auto; }
.detail-title { font-size: 18px; font-weight: 700; color: #333; margin-bottom: 10px; line-height: 1.4; }
.detail-content { font-size: 15px; color: #444; line-height: 1.8; white-space: pre-wrap; word-break: break-word; margin-bottom: 12px; }
.detail-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 12px; }
.tag { font-size: 13px; color: #ff2442; }
.detail-actions { display: flex; gap: 16px; padding: 10px 0; border-top: 1px solid #f0f0f0; border-bottom: 1px solid #f0f0f0; margin-bottom: 12px; }
.action-btn { border: none; background: none; font-size: 14px; color: #999; cursor: pointer; }
.action-btn.active { color: #ff2442; }

.detail-comments { padding-top: 4px; }
.comment-item { padding: 8px 0; border-bottom: 1px solid #f0f0f0; font-size: 14px; line-height: 1.6; }
.comment-item:last-child { border-bottom: none; }
.comment-author { font-weight: 600; font-size: 13px; }
.reply-to { font-size: 12px; color: #999; }
.comment-sep { color: #999; }
.comment-text { color: #333; }
.no-comments { text-align: center; font-size: 13px; color: #ccc; padding: 16px 0; }

/* 个人 */
.profile-card { background: #fff; padding: 24px; text-align: center; }
.profile-avatar { width: 60px; height: 60px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 24px; font-weight: 700; margin: 0 auto 10px; }
.profile-name { font-size: 18px; font-weight: 700; color: #333; margin-bottom: 12px; }
.profile-stats { display: flex; justify-content: center; gap: 32px; }
.stat-item { display: flex; flex-direction: column; align-items: center; gap: 2px; }
.stat-num { font-size: 18px; font-weight: 700; color: #333; }
.stat-label { font-size: 12px; color: #999; }

.prompt-actions { display: flex; gap: 8px; padding: 10px 16px; border-top: 1px solid #f0f0f0; }
.prompt-action-btn { flex: 1; padding: 10px; border: 1px solid #e8e8e8; background: #f7f8fa; border-radius: 8px; font-size: 13px; color: #666; cursor: pointer; }
.prompt-action-btn.danger { color: #ff4757; border-color: #ffcccc; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; }
.slide-up-enter-from .detail-sheet, .slide-up-leave-to .detail-sheet { transform: translateY(100%); }
</style>
