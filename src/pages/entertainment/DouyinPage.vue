<template>
  <div class="douyin-page">
    <!-- 顶部导航 -->
    <div class="douyin-header">
      <button class="header-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="#fff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6" />
        </svg>
      </button>
      <div class="header-tabs">
        <button class="htab" :class="{ active: activeTab === 'following' }" @click="activeTab = 'following'">关注</button>
        <button class="htab" :class="{ active: activeTab === 'recommend' }" @click="activeTab = 'recommend'">推荐</button>
      </div>
      <div class="header-right">
        <button class="header-btn" @click="showPromptEditor = true">⚙</button>
        <button class="header-btn" @click="handleGenerate" :disabled="store.generating">
          <svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ spinning: store.generating }">
            <polyline points="23 4 23 10 17 10" />
            <polyline points="1 20 1 14 7 14" />
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 主内容 -->
    <div class="douyin-body" @wheel="handleScroll" @touchstart="touchStart" @touchend="touchEnd">
      <div v-if="store.lastError" class="error-banner" @click="store.lastError = ''">
        △ {{ store.lastError }}
      </div>

      <div v-if="store.generating && store.douyinVideos.length === 0" class="loading-state">
        <div class="loading-spinner"></div>
        <span>AI正在生成内容...</span>
      </div>

      <div v-else-if="store.douyinVideos.length === 0" class="empty-state">
        <div class="empty-icon">▷</div>
        <div class="empty-title">暂无视频</div>
        <div class="empty-sub">点击右上角刷新按钮生成AI内容</div>
        <button class="generate-btn" @click="handleGenerate" :disabled="store.generating">
          ✦ 生成抖音内容
        </button>
      </div>

      <!-- 竖屏全屏视频卡片 -->
      <div v-else class="video-container" :style="{ transform: `translateY(-${currentIndex * 100}%)` }">
        <div v-for="(video, idx) in store.douyinVideos" :key="video.id" class="video-card">
          <!-- 模拟视频背景 -->
          <div class="video-bg" :style="{ background: getVideoImages(video).length ? '#000' : getVideoBg(idx) }">
            <img v-if="getVideoImages(video).length" :src="getVideoImages(video)[0]" class="video-bg-img" alt="" />
            <div class="video-text-overlay">
              <div class="video-desc-main">{{ video.description }}</div>
            </div>
          </div>

          <!-- 右侧操作栏 -->
          <div class="video-sidebar">
            <div class="sidebar-item">
              <div class="sidebar-avatar" :style="{ background: getAvatarColor(video.author) }">
                {{ video.author.charAt(0) }}
              </div>
            </div>
            <div class="sidebar-item">
              <button class="sidebar-btn" :class="{ liked: video.isLiked }" @click="store.toggleDouyinLike(video.id)">
                <svg viewBox="0 0 24 24" width="28" height="28" :fill="video.isLiked ? '#fe2c55' : 'none'" :stroke="video.isLiked ? '#fe2c55' : '#fff'" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
                </svg>
              </button>
              <span class="sidebar-count">{{ video.likes }}</span>
            </div>
            <div class="sidebar-item">
              <button class="sidebar-btn" @click="toggleComments(video)">
                <svg viewBox="0 0 24 24" width="28" height="28" fill="none" stroke="#fff" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                </svg>
              </button>
              <span class="sidebar-count">{{ video.commentCount }}</span>
            </div>
            <div class="sidebar-item">
              <button class="sidebar-btn">
                <svg viewBox="0 0 24 24" width="28" height="28" fill="none" stroke="#fff" stroke-width="2">
                  <path d="M17 1l4 4-4 4" /><path d="M3 11V9a4 4 0 0 1 4-4h14" />
                  <path d="M7 23l-4-4 4-4" /><path d="M21 13v2a4 4 0 0 1-4 4H3" />
                </svg>
              </button>
              <span class="sidebar-count">{{ video.shares }}</span>
            </div>
            <div class="sidebar-item">
              <button class="sidebar-btn delete-btn" @click="store.deleteDouyinVideo(video.id)">
                <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="#fff" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
          </div>

          <!-- 底部信息 -->
          <div class="video-bottom">
            <div class="video-author">@{{ video.author }}</div>
            <div class="video-desc">{{ video.description }}</div>
            <div v-if="video.tags.length" class="video-tags">
              <span v-for="tag in video.tags" :key="tag" class="vtag">#{{ tag }}</span>
            </div>
          </div>

          <!-- 评论弹窗 -->
          <Transition name="slide-up">
            <div v-if="video.showComments" class="comments-panel" @click.stop>
              <div class="comments-header">
                <span>{{ video.comments.length }} 条评论</span>
                <button class="close-comments" @click="video.showComments = false">✕</button>
              </div>
              <div class="comments-list">
                <div v-for="c in video.comments" :key="c.id" class="dy-comment">
                  <div class="dy-comment-avatar" :style="{ background: getAvatarColor(c.author) }">{{ c.author.charAt(0) }}</div>
                  <div class="dy-comment-body">
                    <div class="dy-comment-name">{{ c.author }}
                      <span v-if="c.isAuthorReply" class="author-badge">作者</span>
                    </div>
                    <div class="dy-comment-text">{{ c.content }}</div>
                    <div class="dy-comment-meta">
                      <span>{{ formatRelativeTime(c.timestamp) }}</span>
                      <button class="dy-comment-like" :class="{ liked: c.isLiked }" @click="store.toggleDouyinCommentLike(video.id, c.id)">
                        ♥ {{ c.likes || '' }}
                      </button>
                    </div>
                  </div>
                </div>
                <div v-if="!video.comments.length" class="no-comments">暂无评论</div>
              </div>
            </div>
          </Transition>
        </div>
      </div>

      <!-- 滑动提示 -->
      <div v-if="store.douyinVideos.length > 1 && currentIndex === 0" class="swipe-hint">
        ↑ 上滑看下一个
      </div>
    </div>

    <!-- 提示词编辑器 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="showPromptEditor" class="prompt-overlay" @click.self="showPromptEditor = false">
          <div class="prompt-sheet">
            <div class="sheet-bar">
              <button class="sheet-cancel" @click="showPromptEditor = false">关闭</button>
              <span class="sheet-title">抖音提示词</span>
              <button class="sheet-post" @click="savePrompt">保存</button>
            </div>
            <textarea v-model="promptText" class="sheet-textarea" placeholder="输入抖音生成提示词..."></textarea>
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
import { ref, onMounted } from 'vue'
import { useSocialAIStore } from '@/stores/socialAI'
import { getAvatarColor, formatRelativeTime } from '@/utils/socialParsers'
import type { DouyinVideo } from '@/utils/socialParsers'
import { getPromptTemplate, setPromptTemplate, resetPromptTemplate } from '@/utils/socialPrompts'

const store = useSocialAIStore()

const activeTab = ref<'following' | 'recommend'>('recommend')
const currentIndex = ref(0)
const showPromptEditor = ref(false)
const promptText = ref('')

const VIDEO_BGS = [
  'linear-gradient(180deg, #1a1a2e, #16213e, #0f3460)',
  'linear-gradient(180deg, #2d1b69, #11052c, #3d0066)',
  'linear-gradient(180deg, #0c0032, #190061, #240090)',
  'linear-gradient(180deg, #000428, #004e92)',
  'linear-gradient(180deg, #200122, #6f0000)',
  'linear-gradient(180deg, #141e30, #243b55)',
  'linear-gradient(180deg, #1e130c, #9a8478)',
  'linear-gradient(180deg, #0f2027, #203a43, #2c5364)',
]

function getVideoBg(idx: number): string {
  return VIDEO_BGS[idx % VIDEO_BGS.length]
}

let touchStartY = 0

function touchStart(e: TouchEvent) {
  touchStartY = e.touches[0].clientY
}

function touchEnd(e: TouchEvent) {
  const dy = touchStartY - e.changedTouches[0].clientY
  if (Math.abs(dy) > 60) {
    if (dy > 0 && currentIndex.value < store.douyinVideos.length - 1) {
      currentIndex.value++
    } else if (dy < 0 && currentIndex.value > 0) {
      currentIndex.value--
    }
  }
}

function handleScroll(e: WheelEvent) {
  if (e.deltaY > 30 && currentIndex.value < store.douyinVideos.length - 1) {
    currentIndex.value++
  } else if (e.deltaY < -30 && currentIndex.value > 0) {
    currentIndex.value--
  }
}

function toggleComments(video: DouyinVideo) {
  video.showComments = !video.showComments
}

onMounted(() => {
  store.loadData('douyin')
  promptText.value = getPromptTemplate('douyin')
})

async function handleGenerate() {
  await store.generateDouyinContent()
}

function savePrompt() {
  setPromptTemplate('douyin', promptText.value)
  showPromptEditor.value = false
}

function resetPrompt() {
  resetPromptTemplate('douyin')
  promptText.value = getPromptTemplate('douyin')
}

function clearAllData() {
  if (confirm('确定要清空所有抖音数据吗？')) {
    store.clearData('douyin')
    showPromptEditor.value = false
  }
}

function getVideoImages(video: any): string[] {
  return Array.isArray(video?.images) ? video.images.filter((s: string) => !!s) : []
}
</script>

<style scoped>
.douyin-page { height: 100%; display: flex; flex-direction: column; background: #000; position: relative; overflow: hidden; }

.douyin-header {
  position: absolute; top: 0; left: 0; right: 0; z-index: 20;
  --status-bar-height: 54px;
  height: calc(var(--status-bar-height) + 44px);
  display: flex; align-items: center; justify-content: space-between;
  padding: var(--status-bar-height) 10px 0;
}
.header-btn { width: 34px; height: 34px; border: none; background: none; color: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center; border-radius: 50%; font-size: 16px; }
.header-btn:active { background: rgba(255,255,255,0.1); }
.header-btn:disabled { opacity: 0.4; }
.header-right { display: flex; gap: 2px; }
.header-tabs { display: flex; gap: 0; }
.htab { border: none; background: none; font-size: 16px; font-weight: 600; color: rgba(255,255,255,0.6); cursor: pointer; padding: 4px 14px; position: relative; }
.htab.active { color: #fff; }
.htab.active::after { content: ''; position: absolute; bottom: -2px; left: 50%; transform: translateX(-50%); width: 20px; height: 3px; background: #fff; border-radius: 2px; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.douyin-body { flex: 1; position: relative; overflow: hidden; }

.error-banner { position: absolute; top: 100px; left: 10px; right: 10px; background: rgba(255,193,7,0.9); color: #333; padding: 10px 16px; font-size: 13px; cursor: pointer; border-radius: 8px; z-index: 15; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: 12px; color: rgba(255,255,255,0.7); }
.loading-spinner { width: 32px; height: 32px; border: 3px solid rgba(255,255,255,0.2); border-top-color: #fe2c55; border-radius: 50%; animation: spin 0.8s linear infinite; }

.empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: 8px; color: rgba(255,255,255,0.7); }
.empty-icon { font-size: 44px; margin-bottom: 4px; }
.empty-title { font-size: 16px; font-weight: 600; }
.empty-sub { font-size: 13px; color: rgba(255,255,255,0.5); }
.generate-btn { margin-top: 16px; padding: 10px 28px; background: linear-gradient(135deg, #fe2c55, #ff6b81); color: #fff; border: none; border-radius: 22px; font-size: 15px; font-weight: 600; cursor: pointer; }
.generate-btn:disabled { opacity: 0.5; }

/* 视频容器 */
.video-container { width: 100%; height: 100%; transition: transform 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94); }
.video-card { width: 100%; height: 100%; position: relative; overflow: hidden; }
.video-bg { position: absolute; top: 0; left: 0; right: 0; bottom: 0; display: flex; align-items: center; justify-content: center; }
.video-bg-img { position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover; }
.video-text-overlay { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center; width: 70%; z-index: 1; }
.video-desc-main { color: rgba(255,255,255,0.15); font-size: 28px; font-weight: 900; line-height: 1.3; word-break: break-word; }

/* 右侧操作栏 */
.video-sidebar { position: absolute; right: 10px; bottom: 130px; display: flex; flex-direction: column; align-items: center; gap: 18px; z-index: 10; }
.sidebar-item { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.sidebar-avatar { width: 44px; height: 44px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 18px; font-weight: 700; border: 2px solid #fff; }
.sidebar-btn { border: none; background: none; cursor: pointer; padding: 0; }
.sidebar-btn:active { transform: scale(0.9); }
.sidebar-count { font-size: 12px; color: #fff; font-weight: 600; }
.delete-btn { opacity: 0.5; }
.delete-btn:hover { opacity: 1; }

/* 底部信息 */
.video-bottom { position: absolute; bottom: 20px; left: 12px; right: 70px; z-index: 10; }
.video-author { font-size: 16px; font-weight: 700; color: #fff; margin-bottom: 6px; }
.video-desc { font-size: 14px; color: rgba(255,255,255,0.9); line-height: 1.5; margin-bottom: 6px; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
.video-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.vtag { font-size: 13px; color: rgba(255,255,255,0.8); font-weight: 500; }

/* 滑动提示 */
.swipe-hint { position: absolute; bottom: 80px; left: 50%; transform: translateX(-50%); color: rgba(255,255,255,0.4); font-size: 12px; animation: bounce 2s infinite; z-index: 10; }
@keyframes bounce { 0%,100% { transform: translateX(-50%) translateY(0); } 50% { transform: translateX(-50%) translateY(-8px); } }

/* 评论面板 */
.comments-panel { position: absolute; bottom: 0; left: 0; right: 0; height: 60%; background: #fff; border-radius: 12px 12px 0 0; z-index: 30; display: flex; flex-direction: column; }
.comments-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #f0f0f0; font-size: 15px; font-weight: 600; color: #333; }
.close-comments { border: none; background: none; font-size: 18px; color: #999; cursor: pointer; }
.comments-list { flex: 1; overflow-y: auto; padding: 8px 16px; }
.dy-comment { display: flex; gap: 10px; padding: 10px 0; border-bottom: 1px solid #f5f5f5; }
.dy-comment-avatar { width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 13px; font-weight: 700; flex-shrink: 0; }
.dy-comment-body { flex: 1; }
.dy-comment-name { font-size: 13px; color: #999; margin-bottom: 4px; }
.author-badge { background: #fe2c55; color: #fff; font-size: 10px; padding: 1px 5px; border-radius: 3px; margin-left: 4px; }
.dy-comment-text { font-size: 14px; color: #333; line-height: 1.5; }
.dy-comment-meta { display: flex; align-items: center; justify-content: space-between; margin-top: 6px; font-size: 12px; color: #bbb; }
.dy-comment-like { border: none; background: none; font-size: 12px; color: #ccc; cursor: pointer; }
.dy-comment-like.liked { color: #fe2c55; }
.no-comments { text-align: center; font-size: 13px; color: #ccc; padding: 30px 0; }

/* 提示词弹窗 */
.prompt-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.4); z-index: 100; display: flex; align-items: flex-end; }
.prompt-sheet { width: 100%; background: #fff; border-radius: 12px 12px 0 0; max-height: 90%; display: flex; flex-direction: column; }
.sheet-bar { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #f0f0f0; flex-shrink: 0; }
.sheet-cancel { border: none; background: none; color: #666; font-size: 15px; cursor: pointer; }
.sheet-title { font-size: 16px; font-weight: 600; color: #333; }
.sheet-post { border: none; background: none; color: #fe2c55; font-size: 15px; font-weight: 600; cursor: pointer; }
.sheet-textarea { flex: 1; padding: 14px 16px; border: none; resize: none; font-size: 13px; color: #333; line-height: 1.6; outline: none; min-height: 200px; font-family: monospace; }
.sheet-textarea::placeholder { color: #ccc; }
.prompt-actions { display: flex; gap: 8px; padding: 10px 16px; border-top: 1px solid #f0f0f0; }
.prompt-action-btn { flex: 1; padding: 10px; border: 1px solid #e8e8e8; background: #f7f8fa; border-radius: 8px; font-size: 13px; color: #666; cursor: pointer; }
.prompt-action-btn.danger { color: #ff4757; border-color: #ffcccc; }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; }
.slide-up-enter-from .comments-panel, .slide-up-leave-to .comments-panel,
.slide-up-enter-from .prompt-sheet, .slide-up-leave-to .prompt-sheet { transform: translateY(100%); }
</style>
