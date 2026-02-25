<template>
  <div class="qzone-page">
    <NavBar title="朋友圈" :show-back="true" back-to="/">
      <template #right>
        <button class="nav-btn" @click="showCompose = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
            <line x1="12" y1="8" x2="12" y2="16" />
            <line x1="8" y1="12" x2="16" y2="12" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- 封面区域 -->
    <div class="cover-section">
      <div class="cover-image"></div>
      <div class="cover-user">
        <span class="cover-name">{{ username }}</span>
        <div class="cover-avatar">
          <img v-if="userAvatar" :src="userAvatar" alt="" />
          <span v-else>{{ username[0] || '?' }}</span>
        </div>
      </div>
    </div>

    <!-- 动态列表 -->
    <div class="moments-list">
      <div v-if="moments.length === 0" class="empty-state">
        <p>还没有动态，发一条吧~</p>
      </div>

      <div v-for="moment in moments" :key="moment.id" class="moment-item">
        <div class="moment-avatar">
          <img v-if="moment.avatar" :src="moment.avatar" :alt="moment.author" />
          <span v-else class="avatar-letter">{{ moment.author[0] }}</span>
        </div>
        <div class="moment-body">
          <div class="moment-author">{{ moment.author }}</div>
          <div class="moment-text">{{ moment.content }}</div>

          <!-- 图片网格 -->
          <div v-if="moment.images && moment.images.length > 0" class="moment-images" :class="imageGridClass(moment.images.length)">
            <div
              v-for="(img, idx) in moment.images"
              :key="idx"
              class="moment-img"
              @click="previewImage(img)"
            >
              <img :src="img" alt="" />
            </div>
          </div>

          <!-- 底部信息 -->
          <div class="moment-meta">
            <span class="moment-time">{{ formatTime(moment.created_at) }}</span>
            <div class="moment-actions">
              <button class="meta-btn" @click="toggleLike(moment)">
                <svg viewBox="0 0 24 24" :fill="moment.liked ? 'var(--color-red)' : 'none'" :stroke="moment.liked ? 'var(--color-red)' : 'currentColor'" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
                </svg>
                <span v-if="moment.likes > 0">{{ moment.likes }}</span>
              </button>
              <button class="meta-btn" @click="openComments(moment)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                </svg>
                <span v-if="moment.comments > 0">{{ moment.comments }}</span>
              </button>
            </div>
          </div>

          <!-- 互动区 -->
          <div v-if="moment.showInteraction" class="interaction-box">
            <!-- 点赞列表 -->
            <div v-if="moment.likedBy && moment.likedBy.length > 0" class="likes-row">
              <svg viewBox="0 0 24 24" fill="var(--color-red)" stroke="none" width="14" height="14">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" />
              </svg>
              <span>{{ moment.likedBy.join('，') }}</span>
            </div>
            <!-- 评论列表 -->
            <div v-if="moment.commentList && moment.commentList.length > 0" class="comments-list">
              <div v-for="(c, ci) in moment.commentList" :key="ci" class="comment-row">
                <span class="comment-author">{{ c.author }}：</span>
                <span class="comment-text">{{ c.text }}</span>
              </div>
            </div>
            <!-- 评论输入 -->
            <div class="comment-input-row">
              <input
                v-model="moment.commentDraft"
                type="text"
                placeholder="评论..."
                class="comment-input"
                @keydown.enter="submitComment(moment)"
              />
              <button class="send-comment" @click="submitComment(moment)" :disabled="!moment.commentDraft?.trim()">发送</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 发布弹窗 -->
    <Teleport to="body">
      <div v-if="showCompose" class="modal-overlay" @click.self="showCompose = false">
        <div class="compose-modal">
          <div class="modal-header">
            <button class="cancel-text" @click="showCompose = false">取消</button>
            <span>发朋友圈</span>
            <button class="publish-btn" @click="publishMoment" :disabled="!composeText.trim()">发布</button>
          </div>
          <div class="compose-body">
            <textarea
              v-model="composeText"
              placeholder="这一刻的想法..."
              class="compose-textarea"
              rows="5"
              autofocus
            ></textarea>
            <div class="compose-images" v-if="composeImages.length > 0">
              <div v-for="(img, idx) in composeImages" :key="idx" class="compose-img-item">
                <img :src="img" alt="" />
                <button class="remove-img" @click="composeImages.splice(idx, 1)">✕</button>
              </div>
            </div>
            <div class="compose-tools">
              <button class="tool-btn" @click="addImage">🖼️ 图片</button>
              <button class="tool-btn" @click="addLocation">📍 位置</button>
              <button class="tool-btn" @click="toggleVisibility">{{ visibility === 'public' ? '🌐 公开' : visibility === 'friends' ? '👥 好友可见' : '🔒 仅自己' }}</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 图片预览 -->
    <Teleport to="body">
      <div v-if="previewImg" class="preview-overlay" @click="previewImg = ''">
        <img :src="previewImg" alt="" class="preview-full" />
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/common/NavBar.vue'

const authStore = useAuthStore()

const username = computed(() => authStore.user?.username || '我')
const userAvatar = computed(() => authStore.user?.avatar || '')

const showCompose = ref(false)
const composeText = ref('')
const composeImages = ref<string[]>([])
const visibility = ref<'public' | 'friends' | 'private'>('public')
const previewImg = ref('')

interface Comment {
  author: string
  text: string
}

interface Moment {
  id: number
  author: string
  avatar: string
  content: string
  images: string[]
  likes: number
  liked: boolean
  likedBy: string[]
  comments: number
  commentList: Comment[]
  commentDraft: string
  showInteraction: boolean
  created_at: string
}

const moments = ref<Moment[]>([
  {
    id: 1,
    author: '小明',
    avatar: '',
    content: '今天天气真好啊，出去散步了一圈 🌞',
    images: [],
    likes: 3,
    liked: false,
    likedBy: ['小红', '小华', '小李'],
    comments: 2,
    commentList: [
      { author: '小红', text: '天气好心情就好~' },
      { author: '小华', text: '在哪里散步呀' },
    ],
    commentDraft: '',
    showInteraction: false,
    created_at: new Date(Date.now() - 3600000).toISOString(),
  },
  {
    id: 2,
    author: '小红',
    avatar: '',
    content: '分享一首好听的歌 🎵 最近单曲循环停不下来',
    images: [],
    likes: 5,
    liked: true,
    likedBy: ['小明', '小华', '小李', '我', '小张'],
    comments: 1,
    commentList: [
      { author: '小明', text: '什么歌？分享一下' },
    ],
    commentDraft: '',
    showInteraction: false,
    created_at: new Date(Date.now() - 7200000).toISOString(),
  },
  {
    id: 3,
    author: '小华',
    avatar: '',
    content: '新学的拿铁拉花 ☕ 虽然有点歪但是第一次还算成功吧哈哈',
    images: [],
    likes: 8,
    liked: false,
    likedBy: ['小红', '小明'],
    comments: 0,
    commentList: [],
    commentDraft: '',
    showInteraction: false,
    created_at: new Date(Date.now() - 86400000).toISOString(),
  },
])

function imageGridClass(count: number): string {
  if (count === 1) return 'grid-1'
  if (count === 2) return 'grid-2'
  if (count === 4) return 'grid-4'
  return 'grid-3'
}

function formatTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}天前`
  return new Date(iso).toLocaleDateString('zh-CN')
}

function toggleLike(moment: Moment) {
  if (moment.liked) {
    moment.likes--
    moment.liked = false
    moment.likedBy = moment.likedBy.filter(n => n !== username.value)
  } else {
    moment.likes++
    moment.liked = true
    moment.likedBy.push(username.value)
  }
}

function openComments(moment: Moment) {
  moment.showInteraction = !moment.showInteraction
}

function submitComment(moment: Moment) {
  const text = moment.commentDraft?.trim()
  if (!text) return
  moment.commentList.push({ author: username.value, text })
  moment.comments++
  moment.commentDraft = ''
  moment.showInteraction = true
}

function publishMoment() {
  if (!composeText.value.trim()) return
  const newMoment: Moment = {
    id: Date.now(),
    author: username.value,
    avatar: userAvatar.value,
    content: composeText.value,
    images: [...composeImages.value],
    likes: 0,
    liked: false,
    likedBy: [],
    comments: 0,
    commentList: [],
    commentDraft: '',
    showInteraction: false,
    created_at: new Date().toISOString(),
  }
  moments.value.unshift(newMoment)
  composeText.value = ''
  composeImages.value = []
  showCompose.value = false
}

function addImage() {
  const url = prompt('输入图片链接:')
  if (url) composeImages.value.push(url)
}

function addLocation() {
  // Placeholder
  composeText.value += ' 📍 某个有趣的地方'
}

function toggleVisibility() {
  const order: ('public' | 'friends' | 'private')[] = ['public', 'friends', 'private']
  const idx = order.indexOf(visibility.value)
  visibility.value = order[(idx + 1) % 3]
}

function previewImage(url: string) {
  previewImg.value = url
}

onMounted(() => {
  // Could fetch from API
})
</script>

<style scoped>
.qzone-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

/* 封面 */
.cover-section {
  position: relative;
  height: 200px;
  flex-shrink: 0;
}

.cover-image {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.cover-user {
  position: absolute;
  bottom: -24px;
  right: 16px;
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.cover-name {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
  margin-bottom: 8px;
}

.cover-avatar {
  width: 64px;
  height: 64px;
  border-radius: 10px;
  overflow: hidden;
  border: 3px solid var(--bg-primary);
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
}

.cover-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 动态列表 */
.moments-list {
  flex: 1;
  overflow-y: auto;
  padding-top: 36px;
}

.moments-list::-webkit-scrollbar {
  display: none;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: var(--text-tertiary);
  font-size: 14px;
}

.moment-item {
  display: flex;
  gap: 10px;
  padding: 16px;
  border-bottom: 0.5px solid var(--separator);
}

.moment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  overflow: hidden;
  background: linear-gradient(135deg, #00CEC9, #81ECEC);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.moment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-letter {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.moment-body {
  flex: 1;
  min-width: 0;
}

.moment-author {
  font-size: 15px;
  font-weight: 600;
  color: var(--brand-primary);
  margin-bottom: 4px;
}

.moment-text {
  font-size: 15px;
  color: var(--text-primary);
  line-height: 1.5;
  margin-bottom: 8px;
  white-space: pre-wrap;
  word-break: break-word;
}

/* 图片网格 */
.moment-images {
  display: grid;
  gap: 4px;
  margin-bottom: 8px;
  border-radius: 8px;
  overflow: hidden;
}

.grid-1 { grid-template-columns: 1fr; max-width: 200px; }
.grid-2 { grid-template-columns: 1fr 1fr; max-width: 240px; }
.grid-3 { grid-template-columns: 1fr 1fr 1fr; max-width: 280px; }
.grid-4 { grid-template-columns: 1fr 1fr; max-width: 240px; }

.moment-img {
  aspect-ratio: 1;
  overflow: hidden;
  cursor: pointer;
}

.moment-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 底部 */
.moment-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.moment-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.moment-actions {
  display: flex;
  gap: 12px;
}

.meta-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: none;
  color: var(--text-tertiary);
  font-size: 12px;
  cursor: pointer;
  padding: 4px;
}

.meta-btn svg {
  width: 16px;
  height: 16px;
}

.meta-btn:active {
  opacity: 0.5;
}

/* 互动区 */
.interaction-box {
  margin-top: 8px;
  background: var(--bg-tertiary);
  border-radius: 6px;
  padding: 8px 10px;
}

.likes-row {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  font-size: 13px;
  color: var(--brand-primary);
  padding-bottom: 6px;
  border-bottom: 0.5px solid var(--separator);
  margin-bottom: 6px;
}

.likes-row svg {
  flex-shrink: 0;
  margin-top: 2px;
}

.comments-list {
  padding-bottom: 6px;
}

.comment-row {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary);
}

.comment-author {
  color: var(--brand-primary);
  font-weight: 500;
}

.comment-input-row {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.comment-input {
  flex: 1;
  border: none;
  background: var(--bg-primary);
  border-radius: 14px;
  padding: 6px 12px;
  font-size: 13px;
  color: var(--text-primary);
  outline: none;
}

.comment-input::placeholder {
  color: var(--text-quaternary);
}

.send-comment {
  border: none;
  background: none;
  color: var(--brand-primary);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0 4px;
}

.send-comment:disabled {
  opacity: 0.4;
}

/* 发布弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 1000;
}

.compose-modal {
  background: var(--bg-primary);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  width: 100%;
  max-height: 80%;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 0.5px solid var(--separator);
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
}

.cancel-text {
  border: none;
  background: none;
  color: var(--text-secondary);
  font-size: 16px;
  cursor: pointer;
}

.publish-btn {
  border: none;
  background: var(--brand-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  padding: 6px 16px;
  border-radius: 16px;
  cursor: pointer;
}

.publish-btn:disabled {
  opacity: 0.4;
}

.compose-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.compose-textarea {
  width: 100%;
  border: none;
  background: none;
  outline: none;
  font-size: 16px;
  color: var(--text-primary);
  line-height: 1.6;
  resize: none;
  font-family: inherit;
  box-sizing: border-box;
}

.compose-textarea::placeholder {
  color: var(--text-quaternary);
}

.compose-images {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.compose-img-item {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
}

.compose-img-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-img {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.6);
  border: none;
  color: #fff;
  font-size: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.compose-tools {
  display: flex;
  gap: 12px;
  margin-top: 16px;
  padding-top: 12px;
  border-top: 0.5px solid var(--separator);
}

.tool-btn {
  border: none;
  background: var(--bg-tertiary);
  padding: 6px 12px;
  border-radius: 14px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
}

/* 图片预览 */
.preview-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  cursor: pointer;
}

.preview-full {
  max-width: 90%;
  max-height: 90%;
  object-fit: contain;
  border-radius: 4px;
}

/* Nav btn */
.nav-btn {
  display: flex;
  align-items: center;
  border: none;
  background: none;
  color: var(--brand-primary);
  cursor: pointer;
  padding: 6px;
}

.nav-btn svg {
  width: 22px;
  height: 22px;
}
</style>
