<template>
  <div class="post-detail-page">
    <NavBar :title="post?.title || '帖子详情'" />

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <span>加载中...</span>
    </div>

    <template v-else-if="post">
      <!-- 帖子内容 -->
      <div class="post-card">
        <div class="post-header">
          <div class="author-avatar">
            <img v-if="post.author?.avatar_url" :src="post.author.avatar_url" alt="" />
            <span v-else>{{ post.author?.username?.charAt(0) || '?' }}</span>
          </div>
          <div class="author-info">
            <span class="author-name">{{ post.author?.username || '匿名' }}</span>
            <span class="post-time">{{ formatTime(post.created_at) }}</span>
          </div>
          <span v-if="post.pinned" class="pin-badge">📌 置顶</span>
        </div>

        <div class="post-board">
          <span class="board-tag">{{ post.board || '综合' }}</span>
        </div>

        <h2 class="post-title">{{ post.title }}</h2>
        <div class="post-content">{{ post.content }}</div>

        <div v-if="post.images?.length" class="post-images">
          <img
            v-for="(img, idx) in post.images"
            :key="idx"
            :src="img"
            class="post-image"
            @click="previewImage(img)"
          />
        </div>

        <div class="post-stats">
          <button class="stat-btn" @click="handleLike">
            <span>👍</span>
            <span>{{ post.likes || 0 }}</span>
          </button>
          <div class="stat-item">
            <span>💬</span>
            <span>{{ comments.length }}</span>
          </div>
        </div>
      </div>

      <!-- 评论区 -->
      <div class="comments-section">
        <div class="comments-header">
          <h3>评论 ({{ comments.length }})</h3>
        </div>

        <!-- 评论输入 -->
        <div class="comment-input-area">
          <div v-if="replyTo" class="reply-hint">
            <span>回复 {{ replyTo.author_name }}：</span>
            <button class="cancel-reply" @click="replyTo = null">✕</button>
          </div>
          <div class="input-row">
            <input
              v-model="commentText"
              type="text"
              :placeholder="replyTo ? `回复 ${replyTo.author_name}...` : '写评论...'"
              class="comment-input"
              @keyup.enter="submitComment"
            />
            <button
              class="send-btn"
              :disabled="!commentText.trim() || submitting"
              @click="submitComment"
            >
              发送
            </button>
          </div>
        </div>

        <!-- 评论列表 -->
        <div v-if="comments.length === 0" class="empty-comments">
          <span>💬</span>
          <p>暂无评论，来抢沙发吧~</p>
        </div>

        <div v-else class="comments-list">
          <div
            v-for="comment in comments"
            :key="comment.id"
            class="comment-item"
          >
            <div class="comment-avatar">
              <img v-if="comment.author_avatar" :src="comment.author_avatar" alt="" />
              <span v-else>{{ comment.author_name?.charAt(0) || '?' }}</span>
            </div>
            <div class="comment-body">
              <div class="comment-meta">
                <span class="comment-author">{{ comment.author_name }}</span>
                <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
              </div>
              <div v-if="comment.parent_id" class="reply-tag">
                回复
              </div>
              <p class="comment-text">{{ comment.content }}</p>
              <div class="comment-actions">
                <button class="action-btn" @click="handleCommentLike(comment)">
                  👍 {{ comment.likes || 0 }}
                </button>
                <button class="action-btn" @click="replyTo = comment">
                  💬 回复
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="error-state">
      <span>😕</span>
      <p>帖子不存在或已被删除</p>
      <button class="back-btn" @click="$router.back()">返回</button>
    </div>

    <!-- 图片预览 -->
    <div v-if="previewImg" class="image-preview-overlay" @click="previewImg = null">
      <img :src="previewImg" class="preview-image" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { postApi } from '@/api/services'
import type { Post, Comment } from '@/api/types'

const route = useRoute()
const postId = Number(route.params.postId)

const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const loading = ref(true)
const commentText = ref('')
const replyTo = ref<Comment | null>(null)
const submitting = ref(false)
const previewImg = ref<string | null>(null)

onMounted(async () => {
  await loadPost()
})

async function loadPost() {
  loading.value = true
  try {
    const [postRes, commentsRes] = await Promise.all([
      postApi.get(postId),
      postApi.listComments(postId),
    ])
    post.value = postRes
    comments.value = commentsRes.data || []
  } catch (e) {
    console.error('加载帖子失败', e)
    post.value = null
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!commentText.value.trim() || submitting.value) return
  submitting.value = true
  try {
    await postApi.createComment(postId, {
      content: commentText.value.trim(),
      parent_id: replyTo.value?.id,
    })
    commentText.value = ''
    replyTo.value = null
    // 重新加载评论
    const res = await postApi.listComments(postId)
    comments.value = res.data || []
  } catch (e) {
    console.error('发送评论失败', e)
  } finally {
    submitting.value = false
  }
}

function handleLike() {
  if (post.value) {
    post.value.likes = (post.value.likes || 0) + 1
  }
}

function handleCommentLike(comment: Comment) {
  comment.likes = (comment.likes || 0) + 1
}

function previewImage(url: string) {
  previewImg.value = url
}

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.post-detail-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.loading-state,
.error-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-secondary);
}

.error-state span {
  font-size: 48px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-blue);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.back-btn {
  padding: 8px 24px;
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.post-card {
  background: var(--bg-secondary);
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-blue);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.author-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.post-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.pin-badge {
  font-size: 12px;
  color: var(--accent-orange);
  background: rgba(255, 149, 0, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
}

.post-board {
  margin-bottom: 8px;
}

.board-tag {
  font-size: 12px;
  color: var(--accent-blue);
  background: rgba(0, 122, 255, 0.1);
  padding: 2px 10px;
  border-radius: 10px;
}

.post-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 10px;
  line-height: 1.4;
}

.post-content {
  font-size: 15px;
  color: var(--text-primary);
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.post-images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 8px;
  margin-top: 12px;
}

.post-image {
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
  border-radius: 8px;
  cursor: pointer;
}

.post-stats {
  display: flex;
  gap: 20px;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.stat-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 16px;
  transition: background 0.2s;
}

.stat-btn:active {
  background: var(--bg-tertiary);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: var(--text-secondary);
}

/* 评论区 */
.comments-section {
  flex: 1;
}

.comments-header {
  padding: 14px 16px 8px;
}

.comments-header h3 {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.comment-input-area {
  padding: 8px 16px 12px;
  border-bottom: 1px solid var(--border-color);
}

.reply-hint {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--accent-blue);
  margin-bottom: 6px;
}

.cancel-reply {
  background: none;
  border: none;
  color: var(--text-tertiary);
  font-size: 14px;
  cursor: pointer;
}

.input-row {
  display: flex;
  gap: 8px;
}

.comment-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--border-color);
  border-radius: 20px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.comment-input:focus {
  border-color: var(--accent-blue);
}

.send-btn {
  padding: 8px 16px;
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
}

.send-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.empty-comments {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 16px;
  color: var(--text-tertiary);
}

.empty-comments span {
  font-size: 40px;
  margin-bottom: 8px;
}

.empty-comments p {
  font-size: 14px;
  margin: 0;
}

.comments-list {
  padding: 0 16px;
}

.comment-item {
  display: flex;
  gap: 10px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent-purple);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-author {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.comment-time {
  font-size: 11px;
  color: var(--text-tertiary);
}

.reply-tag {
  font-size: 11px;
  color: var(--accent-blue);
  margin-bottom: 4px;
}

.comment-text {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.5;
  margin: 0 0 6px;
  word-break: break-word;
}

.comment-actions {
  display: flex;
  gap: 16px;
}

.action-btn {
  background: none;
  border: none;
  font-size: 12px;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: 2px 0;
}

.action-btn:active {
  color: var(--accent-blue);
}

/* 图片预览 */
.image-preview-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}

.preview-image {
  max-width: 90%;
  max-height: 90%;
  object-fit: contain;
  border-radius: 8px;
}
</style>
