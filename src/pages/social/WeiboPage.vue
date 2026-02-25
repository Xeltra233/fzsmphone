<template>
  <div class="weibo-page">
    <NavBar title="微博" back-to="/">
      <template #right>
        <button class="icon-btn" @click="showCompose = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- Tabs -->
    <div class="weibo-tabs">
      <button class="wtab" :class="{ active: tab === 'hot' }" @click="tab = 'hot'">热门</button>
      <button class="wtab" :class="{ active: tab === 'latest' }" @click="tab = 'latest'">最新</button>
      <button class="wtab" :class="{ active: tab === 'mine' }" @click="tab = 'mine'">我的</button>
    </div>

    <!-- Feed -->
    <div class="weibo-feed">
      <div v-if="posts.length === 0" class="empty-state">
        <div class="empty-emoji">📝</div>
        <div class="empty-title">暂无微博</div>
        <div class="empty-sub">发第一条微博吧</div>
      </div>

      <div v-for="post in posts" :key="post.id" class="weibo-card">
        <!-- Header -->
        <div class="wb-header">
          <div class="wb-avatar">
            {{ post.author?.display_name?.charAt(0) || '👤' }}
          </div>
          <div class="wb-user">
            <span class="wb-name">{{ post.author?.display_name || '匿名用户' }}</span>
            <span class="wb-time">{{ formatTime(post.created_at) }}</span>
          </div>
          <div v-if="post.is_hot" class="hot-badge">🔥 热门</div>
        </div>

        <!-- Content -->
        <p class="wb-content">{{ post.content }}</p>

        <!-- Images -->
        <div v-if="post.images?.length" class="wb-images" :class="`grid-${Math.min(post.images.length, 9)}`">
          <div v-for="(img, i) in post.images.slice(0, 9)" :key="i" class="wb-img">
            <img :src="img" alt="" />
          </div>
        </div>

        <!-- Actions -->
        <div class="wb-actions">
          <button class="wb-action" @click.stop>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 1l4 4-4 4" /><path d="M3 11V9a4 4 0 0 1 4-4h14" />
              <path d="M7 23l-4-4 4-4" /><path d="M21 13v2a4 4 0 0 1-4 4H3" />
            </svg>
            <span>{{ post.reposts || 0 }}</span>
          </button>
          <button class="wb-action" @click.stop>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
            </svg>
            <span>{{ post.comments || 0 }}</span>
          </button>
          <button class="wb-action" @click.stop>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
            </svg>
            <span>{{ post.likes || 0 }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Compose -->
    <Teleport to="body">
      <Transition name="slide-up">
        <div v-if="showCompose" class="compose-overlay" @click.self="showCompose = false">
          <div class="compose-sheet">
            <div class="sheet-bar">
              <button @click="showCompose = false" class="sheet-cancel">取消</button>
              <span class="sheet-title">发微博</span>
              <button class="sheet-post" :disabled="!newContent.trim()" @click="submitWeibo">发布</button>
            </div>
            <textarea
              v-model="newContent"
              class="sheet-textarea"
              placeholder="分享新鲜事..."
              autofocus
            ></textarea>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { api } from '@/api/client'

interface WeiboPost {
  id: number
  content: string
  images: string[]
  likes: number
  reposts: number
  comments: number
  is_hot: boolean
  created_at: string
  author?: { display_name: string; avatar_url: string }
}

const tab = ref('hot')
const posts = ref<WeiboPost[]>([])
const showCompose = ref(false)
const newContent = ref('')

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diffMin = Math.floor((now.getTime() - d.getTime()) / 60000)
  if (diffMin < 1) return '刚刚'
  if (diffMin < 60) return `${diffMin}分钟前`
  const diffH = Math.floor(diffMin / 60)
  if (diffH < 24) return `${diffH}小时前`
  return `${d.getMonth() + 1}-${d.getDate()}`
}

async function fetchPosts() {
  try {
    const res = await api.get<{ data: WeiboPost[] }>('/api/weibo')
    posts.value = res.data || []
  } catch { /* empty */ }
}

async function submitWeibo() {
  if (!newContent.value.trim()) return
  try {
    await api.post('/api/weibo', { content: newContent.value.trim() })
    newContent.value = ''
    showCompose.value = false
    fetchPosts()
  } catch { /* empty */ }
}

onMounted(fetchPosts)
</script>

<style scoped>
.weibo-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
}

.weibo-tabs {
  display: flex;
  background: var(--bg-primary);
  border-bottom: 0.5px solid var(--separator);
  flex-shrink: 0;
}

.wtab {
  flex: 1;
  padding: 10px 0;
  border: none;
  background: none;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  position: relative;
  -webkit-tap-highlight-color: transparent;
}

.wtab.active {
  color: var(--ios-orange);
  font-weight: 600;
}

.wtab.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 3px;
  background: var(--ios-orange);
  border-radius: 2px;
}

.weibo-feed {
  flex: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.weibo-card {
  background: var(--bg-primary);
  padding: 14px 16px;
  border-bottom: 6px solid var(--bg-secondary);
}

.wb-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.wb-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ff6b6b, #ee5a24);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.wb-user {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.wb-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.wb-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.hot-badge {
  font-size: 12px;
  color: var(--ios-orange);
  font-weight: 600;
}

.wb-content {
  font-size: 15px;
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0 0 10px;
  white-space: pre-wrap;
}

.wb-images {
  display: grid;
  gap: 4px;
  margin-bottom: 10px;
}

.wb-images.grid-1 { grid-template-columns: 1fr; max-width: 240px; }
.wb-images.grid-2 { grid-template-columns: 1fr 1fr; }
.wb-images.grid-3 { grid-template-columns: 1fr 1fr 1fr; }
.wb-images.grid-4 { grid-template-columns: 1fr 1fr; }
.wb-images.grid-5,
.wb-images.grid-6 { grid-template-columns: 1fr 1fr 1fr; }
.wb-images.grid-7,
.wb-images.grid-8,
.wb-images.grid-9 { grid-template-columns: 1fr 1fr 1fr; }

.wb-img {
  aspect-ratio: 1;
  border-radius: 6px;
  overflow: hidden;
  background: var(--fill-tertiary);
}

.wb-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.wb-actions {
  display: flex;
  justify-content: space-around;
  padding-top: 8px;
}

.wb-action {
  display: flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: none;
  color: var(--text-tertiary);
  font-size: 13px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 16px;
  transition: background 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.wb-action:active {
  background: var(--fill-tertiary);
}

.wb-action svg {
  width: 18px;
  height: 18px;
}

/* Empty & Compose reuse similar patterns */
.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; }
.empty-emoji { font-size: 48px; margin-bottom: 12px; }
.empty-title { font-size: 17px; font-weight: 600; color: var(--text-secondary); margin-bottom: 4px; }
.empty-sub { font-size: 14px; color: var(--text-tertiary); }

.icon-btn {
  width: 32px; height: 32px; border: none; background: none;
  color: var(--brand-primary); cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.icon-btn:active { opacity: 0.5; }
.icon-btn svg { width: 22px; height: 22px; }

.compose-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.3);
  z-index: 1000; display: flex; align-items: flex-end;
}

.compose-sheet {
  width: 100%; background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
}

.sheet-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px; border-bottom: 0.5px solid var(--separator);
}

.sheet-title { font-size: 17px; font-weight: 600; color: var(--text-primary); }

.sheet-cancel, .sheet-post {
  border: none; background: none; font-size: 16px; cursor: pointer; padding: 4px 8px;
}
.sheet-cancel { color: var(--text-secondary); }
.sheet-post { color: var(--ios-orange); font-weight: 600; }
.sheet-post:disabled { color: var(--text-quaternary); }

.sheet-textarea {
  width: 100%; min-height: 120px; padding: 16px;
  border: none; background: none; font-size: 15px; line-height: 1.6;
  color: var(--text-primary); outline: none; resize: none; font-family: inherit;
}
.sheet-textarea::placeholder { color: var(--text-quaternary); }

.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s var(--ease-ios); }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); opacity: 0.8; }
</style>
