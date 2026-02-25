<template>
  <div class="forum-page">
    <NavBar title="论坛" back-to="/">
      <template #right>
        <button class="icon-btn" @click="showCompose = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- Board Tabs -->
    <div class="board-tabs">
      <button
        v-for="tab in boards"
        :key="tab.id"
        class="tab-btn"
        :class="{ active: currentBoard === tab.id }"
        @click="currentBoard = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Posts List -->
    <div class="posts-list">
      <div v-if="posts.length === 0" class="empty-state">
        <div class="empty-emoji">📋</div>
        <div class="empty-title">暂无帖子</div>
        <div class="empty-subtitle">来发第一个帖子吧</div>
      </div>

      <div
        v-for="post in posts"
        :key="post.id"
        class="post-card"
        @click="$router.push(`/forum/post/${post.id}`)"
      >
        <!-- Pinned badge -->
        <div v-if="post.pinned" class="pin-badge">📌 置顶</div>

        <div class="post-header">
          <div class="post-author-avatar">
            {{ post.author?.display_name?.charAt(0) || '👤' }}
          </div>
          <div class="post-author-info">
            <span class="post-author-name">{{ post.author?.display_name || '匿名' }}</span>
            <span class="post-date">{{ formatDate(post.created_at) }}</span>
          </div>
          <span class="post-board-tag">{{ getBoardLabel(post.board) }}</span>
        </div>

        <h3 class="post-title">{{ post.title }}</h3>
        <p class="post-content" v-if="post.content">{{ truncate(post.content, 120) }}</p>

        <!-- Images preview -->
        <div v-if="post.images?.length" class="post-images">
          <div
            v-for="(img, i) in post.images.slice(0, 3)"
            :key="i"
            class="post-img"
          >
            <img :src="img" alt="" />
            <div v-if="i === 2 && post.images.length > 3" class="img-more">
              +{{ post.images.length - 3 }}
            </div>
          </div>
        </div>

        <div class="post-footer">
          <div class="post-stat">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
            </svg>
            <span>{{ post.likes }}</span>
          </div>
          <div class="post-stat">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
            </svg>
            <span>{{ post.comments }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Compose Modal -->
    <Teleport to="body">
      <Transition name="slide-up">
        <div v-if="showCompose" class="compose-overlay">
          <div class="compose-panel">
            <div class="compose-header">
              <button class="compose-cancel" @click="showCompose = false">取消</button>
              <span class="compose-title">发帖</span>
              <button class="compose-submit" :disabled="!canSubmit" @click="submitPost">发布</button>
            </div>

            <div class="compose-body">
              <select v-model="newPost.board" class="board-select">
                <option v-for="b in boards" :key="b.id" :value="b.id">{{ b.label }}</option>
              </select>
              <input
                v-model="newPost.title"
                type="text"
                class="compose-input title-input"
                placeholder="标题"
                maxlength="100"
              />
              <textarea
                v-model="newPost.content"
                class="compose-input content-input"
                placeholder="说点什么..."
                rows="6"
              ></textarea>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { api } from '@/api/client'

interface PostAuthor {
  display_name: string
  avatar_url: string
}

interface Post {
  id: number
  user_id: number
  board: string
  title: string
  content: string
  images: string[]
  likes: number
  comments: number
  pinned: boolean
  created_at: string
  author?: PostAuthor
}

const boards = [
  { id: 'all', label: '全部' },
  { id: 'general', label: '综合' },
  { id: 'chat', label: '闲聊' },
  { id: 'share', label: '分享' },
  { id: 'question', label: '提问' },
  { id: 'announce', label: '公告' },
]

const currentBoard = ref('all')
const posts = ref<Post[]>([])
const showCompose = ref(false)

const newPost = reactive({
  board: 'general',
  title: '',
  content: '',
})

const canSubmit = computed(() => newPost.title.trim().length > 0)

function getBoardLabel(boardId: string): string {
  return boards.find(b => b.id === boardId)?.label || boardId
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return '刚刚'
  if (diffMin < 60) return `${diffMin}分钟前`
  const diffHour = Math.floor(diffMin / 60)
  if (diffHour < 24) return `${diffHour}小时前`
  const diffDay = Math.floor(diffHour / 24)
  if (diffDay < 7) return `${diffDay}天前`
  return `${d.getMonth() + 1}/${d.getDate()}`
}

function truncate(text: string, max: number): string {
  return text.length > max ? text.slice(0, max) + '...' : text
}

async function fetchPosts() {
  try {
    const params = currentBoard.value !== 'all' ? `?board=${currentBoard.value}` : ''
    const res = await api.get<{ data: Post[] }>(`/api/posts${params}`)
    posts.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch posts:', e)
  }
}

async function submitPost() {
  if (!canSubmit.value) return
  try {
    await api.post('/api/posts', {
      board: newPost.board,
      title: newPost.title.trim(),
      content: newPost.content.trim(),
    })
    showCompose.value = false
    newPost.title = ''
    newPost.content = ''
    newPost.board = 'general'
    fetchPosts()
  } catch (e) {
    console.error('Failed to create post:', e)
  }
}

onMounted(() => {
  fetchPosts()
})
</script>

<style scoped>
.forum-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
}

/* Board Tabs */
.board-tabs {
  display: flex;
  gap: 0;
  padding: 0 12px;
  background: var(--bg-primary);
  border-bottom: 0.5px solid var(--separator);
  overflow-x: auto;
  flex-shrink: 0;
  -webkit-overflow-scrolling: touch;
}

.board-tabs::-webkit-scrollbar {
  display: none;
}

.tab-btn {
  padding: 10px 16px;
  border: none;
  background: none;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
  position: relative;
  transition: color 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.tab-btn.active {
  color: var(--brand-primary);
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background: var(--brand-primary);
  border-radius: 2px;
}

/* Posts List */
.posts-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  -webkit-overflow-scrolling: touch;
}

.post-card {
  background: var(--bg-primary);
  border-radius: 14px;
  padding: 14px 16px;
  cursor: pointer;
  transition: transform 0.15s var(--ease-ios);
}

.post-card:active {
  transform: scale(0.98);
}

.pin-badge {
  font-size: 12px;
  color: var(--ios-orange);
  font-weight: 600;
  margin-bottom: 8px;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.post-author-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: white;
  font-weight: 600;
  flex-shrink: 0;
}

.post-author-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.post-author-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.post-date {
  font-size: 12px;
  color: var(--text-tertiary);
}

.post-board-tag {
  font-size: 11px;
  color: var(--brand-primary);
  background: rgba(88, 86, 214, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
  flex-shrink: 0;
}

.post-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 6px;
  line-height: 1.4;
}

.post-content {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0 0 10px;
}

.post-images {
  display: flex;
  gap: 6px;
  margin-bottom: 10px;
}

.post-img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  background: var(--fill-tertiary);
}

.post-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.img-more {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.post-footer {
  display: flex;
  gap: 20px;
  padding-top: 8px;
  border-top: 0.5px solid var(--separator);
}

.post-stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--text-tertiary);
}

.post-stat svg {
  width: 16px;
  height: 16px;
}

/* Empty */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 20px;
}

.empty-emoji { font-size: 48px; margin-bottom: 12px; }
.empty-title { font-size: 17px; font-weight: 600; color: var(--text-secondary); margin-bottom: 4px; }
.empty-subtitle { font-size: 14px; color: var(--text-tertiary); }

/* Icon button */
.icon-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  color: var(--brand-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  -webkit-tap-highlight-color: transparent;
}

.icon-btn:active { opacity: 0.5; }
.icon-btn svg { width: 22px; height: 22px; }

/* Compose */
.compose-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 1000;
  display: flex;
  align-items: flex-end;
}

.compose-panel {
  width: 100%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.compose-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 0.5px solid var(--separator);
}

.compose-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
}

.compose-cancel,
.compose-submit {
  border: none;
  background: none;
  font-size: 16px;
  cursor: pointer;
  padding: 4px 8px;
  -webkit-tap-highlight-color: transparent;
}

.compose-cancel { color: var(--text-secondary); }
.compose-submit { color: var(--brand-primary); font-weight: 600; }
.compose-submit:disabled { color: var(--text-quaternary); }

.compose-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
}

.board-select {
  padding: 8px 12px;
  border: 1px solid var(--separator);
  border-radius: 10px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
}

.compose-input {
  border: none;
  background: none;
  color: var(--text-primary);
  font-size: 15px;
  outline: none;
  font-family: inherit;
}

.title-input {
  font-size: 18px;
  font-weight: 600;
  padding: 0;
}

.title-input::placeholder { color: var(--text-quaternary); font-weight: 400; }

.content-input {
  resize: none;
  line-height: 1.6;
  min-height: 120px;
}

.content-input::placeholder { color: var(--text-quaternary); }

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s var(--ease-ios);
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0.8;
}
</style>
