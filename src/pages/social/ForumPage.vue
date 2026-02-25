<template>
  <div class="forum-page">
    <NavBar title="论坛" back-to="/">
      <template #right>
        <div class="header-actions">
          <button class="hdr-btn" @click="showPromptEditor = true" title="提示词设置">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18"><path d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6z"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
          </button>
          <button class="hdr-btn" @click="handleRefresh" :disabled="store.generating" title="刷新">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18" :class="{ spinning: store.generating }"><path d="M23 4v6h-6"/><path d="M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
          </button>
          <button class="hdr-btn compose-btn" @click="showCompose = true" title="发帖">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18"><path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
          </button>
        </div>
      </template>
    </NavBar>

    <!-- 帖子列表 -->
    <div class="forum-content" ref="contentRef">
      <!-- 全局错误提示（所有状态下都显示） -->
      <div v-if="store.lastError" class="error-toast" @click="store.lastError = ''">
        ⚠️ {{ store.lastError }}
      </div>

      <!-- 加载中 -->
      <div v-if="store.forumLoading && store.forumThreads.length === 0" class="loading-state">
        <div class="loading-spinner">💬</div>
        <div class="loading-text">正在生成论坛内容...</div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="store.forumThreads.length === 0" class="empty-state">
        <div class="empty-icon">💬</div>
        <div class="empty-text">暂无帖子</div>
        <div class="empty-hint">点击右上角刷新按钮让AI生成帖子～</div>
        <button class="generate-btn" @click="handleRefresh" :disabled="store.generating">
          {{ store.generating ? '生成中...' : '✨ 生成论坛内容' }}
        </button>
      </div>

      <!-- 帖子卡片列表 -->
      <div v-else class="thread-list">

        <!-- 顶部加载指示 -->
        <div v-if="store.generating" class="top-loading">
          <span class="dot-loading">生成中</span>
        </div>

        <div
          v-for="thread in store.forumThreads"
          :key="thread.id"
          class="thread-item"
          @click="openThread(thread.id)"
        >
          <div class="thread-header">
            <div class="author-avatar" :style="{ background: getColor(thread.author) }">
              {{ thread.author[0] || '?' }}
            </div>
            <div class="thread-author">
              <div class="author-name">{{ thread.author }}</div>
            </div>
            <div class="thread-id-tag">t{{ thread.id }}</div>
            <button class="delete-btn" @click.stop="handleDelete(thread.id)" title="删除">×</button>
          </div>

          <div class="post-body">
            <h3 class="thread-title">{{ thread.title }}</h3>
            <div class="thread-content">{{ truncate(thread.content, 100) }}</div>
          </div>

          <div class="thread-stats">
            <button class="action-btn" :class="{ liked: thread.isLiked }" @click.stop="store.toggleForumLike(thread.id)">
              <span>{{ thread.isLiked ? '❤️' : '🤍' }}</span>
              <span>{{ thread.likes }}</span>
            </button>
            <button class="action-btn">
              <span>💬</span>
              <span>{{ thread.replies.length }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 发帖弹窗 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="showCompose" class="compose-overlay" @click.self="showCompose = false">
          <div class="compose-panel">
            <div class="dialog-header">
              <button class="dialog-close" @click="showCompose = false">取消</button>
              <h3>发新帖</h3>
              <button class="dialog-submit" :disabled="!newTitle.trim() || store.generating" @click="submitPost">
                {{ store.generating ? '发送中...' : '✈' }}
              </button>
            </div>
            <div class="dialog-body">
              <input v-model="newTitle" type="text" class="post-title-input" placeholder="请输入帖子标题..." maxlength="100" />
              <textarea v-model="newContent" class="post-content-input" placeholder="分享你的想法..." rows="5"></textarea>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 提示词编辑弹窗 -->
    <Teleport to="#phone-overlay">
      <Transition name="slide-up">
        <div v-if="showPromptEditor" class="compose-overlay" @click.self="showPromptEditor = false">
          <div class="compose-panel prompt-panel">
            <div class="dialog-header">
              <button class="dialog-close" @click="showPromptEditor = false">关闭</button>
              <h3>论坛提示词</h3>
              <button class="dialog-submit" @click="savePrompt">保存</button>
            </div>
            <div class="dialog-body">
              <p class="prompt-hint">自定义AI生成论坛内容的提示词。支持变量：<code v-pre>{{characters}}</code> <code v-pre>{{user}}</code> <code v-pre>{{action}}</code></p>
              <textarea v-model="promptText" class="prompt-textarea" rows="12"></textarea>
              <div class="prompt-actions">
                <button class="reset-btn" @click="resetPrompt">恢复默认</button>
                <button class="clear-data-btn" @click="clearForumData">清除所有帖子</button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'
import { getAvatarColor, truncateText } from '@/utils/socialParsers'
import { getPromptTemplate, setPromptTemplate, resetPromptTemplate } from '@/utils/socialPrompts'

const router = useRouter()
const store = useSocialAIStore()

const showCompose = ref(false)
const showPromptEditor = ref(false)
const newTitle = ref('')
const newContent = ref('')
const promptText = ref('')
const contentRef = ref<HTMLElement>()

function getColor(name: string) {
  return getAvatarColor(name)
}

function truncate(text: string, max: number) {
  return truncateText(text, max)
}

function openThread(threadId: string) {
  router.push(`/forum/post/${threadId}`)
}

async function handleRefresh() {
  await store.generateForumContent()
}

async function submitPost() {
  if (!newTitle.value.trim()) return
  showCompose.value = false
  await store.forumPost(newTitle.value.trim(), newContent.value.trim())
  newTitle.value = ''
  newContent.value = ''
}

function handleDelete(threadId: string) {
  if (confirm('确定删除这个帖子吗？')) {
    store.deleteForumThread(threadId)
  }
}

function savePrompt() {
  setPromptTemplate('forum', promptText.value)
  showPromptEditor.value = false
}

function resetPrompt() {
  resetPromptTemplate('forum')
  promptText.value = getPromptTemplate('forum')
}

function clearForumData() {
  if (confirm('确定清除所有帖子数据吗？')) {
    store.clearData('forum')
  }
}

onMounted(() => {
  store.loadData('forum')
  promptText.value = getPromptTemplate('forum')
})
</script>

<style scoped>
.forum-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fcfaf7;
  font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', sans-serif;
}

/* 头部按钮 */
.header-actions {
  display: flex;
  gap: 4px;
  align-items: center;
}

.hdr-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  color: var(--brand-primary, #5856d6);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  -webkit-tap-highlight-color: transparent;
}

.hdr-btn:active { opacity: 0.5; }
.hdr-btn:disabled { opacity: 0.3; }

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 内容区域 */
.forum-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  -webkit-overflow-scrolling: touch;
}

.forum-content::-webkit-scrollbar { display: none; }

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #8f8a82;
}

.loading-spinner {
  font-size: 48px;
  margin-bottom: 16px;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.1); opacity: 1; }
}

.loading-text { font-size: 14px; }

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  text-align: center;
  color: #b8b2a9;
  padding: 20px;
}

.empty-icon { font-size: 48px; margin-bottom: 16px; opacity: 0.6; }
.empty-text { font-size: 18px; font-weight: 600; color: #8f8a82; margin-bottom: 8px; }
.empty-hint { font-size: 14px; margin-bottom: 24px; line-height: 1.5; }

.generate-btn {
  background: linear-gradient(135deg, #e8a0bf 0%, #e1c195 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 25px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(184, 169, 154, 0.3);
  transition: all 0.3s ease;
}

.generate-btn:active { transform: scale(0.95); }
.generate-btn:disabled { opacity: 0.5; }

/* 错误提示 */
.error-toast {
  background: #fff3cd;
  color: #856404;
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 13px;
  cursor: pointer;
}

/* 顶部加载 */
.top-loading {
  text-align: center;
  padding: 8px;
  color: #e8a0bf;
  font-size: 13px;
}

.dot-loading::after {
  content: '';
  animation: dots 1.5s steps(4) infinite;
}

@keyframes dots {
  0% { content: ''; }
  25% { content: '.'; }
  50% { content: '..'; }
  75% { content: '...'; }
}

/* 帖子卡片 */
.thread-item {
  background: #ffffff;
  border-radius: 16px;
  padding: 12px;
  cursor: pointer;
  border: 1px solid #f1ede8;
  box-shadow: 0 4px 15px rgba(184, 169, 154, 0.15);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.thread-item:active {
  transform: scale(0.98);
}

.thread-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  position: relative;
}

.author-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.thread-author { flex: 1; }

.author-name {
  font-weight: 600;
  font-size: 13px;
  color: #5c554e;
}

.thread-id-tag {
  background: #f1ede8;
  color: #b8b2a9;
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 99px;
  font-weight: 500;
}

.delete-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: #9ca3af;
  color: white;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.delete-btn:active { opacity: 1; }

.post-body { margin-bottom: 8px; }

.thread-title {
  font-size: 15px;
  font-weight: 600;
  color: #5c554e;
  line-height: 1.3;
  margin: 0 0 4px;
}

.thread-content {
  font-size: 13px;
  line-height: 1.4;
  color: #8f8a82;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.thread-stats {
  display: flex;
  gap: 20px;
}

.action-btn {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  color: #8f8a82;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.2s ease;
  -webkit-tap-highlight-color: transparent;
}

.action-btn:active { background: #e9ecef; }
.action-btn.liked { background: rgba(231, 76, 60, 0.1); border-color: #e74c3c; color: #e74c3c; }

/* 发帖弹窗 */
.compose-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.compose-panel {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-height: 80vh;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.prompt-panel {
  max-height: 90vh;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.dialog-header h3 {
  margin: 0;
  font-size: 17px;
  color: #333;
  font-weight: 600;
}

.dialog-close, .dialog-submit {
  border: none;
  background: none;
  font-size: 15px;
  cursor: pointer;
  padding: 4px 8px;
  -webkit-tap-highlight-color: transparent;
}

.dialog-close { color: #8f8a82; }
.dialog-submit { color: #e8a0bf; font-weight: 600; }
.dialog-submit:disabled { opacity: 0.4; }

.dialog-body { padding: 20px; }

.post-title-input {
  width: 100%;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 12px 16px;
  font-size: 16px;
  margin-bottom: 12px;
  outline: none;
  box-sizing: border-box;
  color: #333;
  background: #fafbfc;
}

.post-title-input:focus { border-color: #e8a0bf; }

.post-content-input {
  width: 100%;
  min-height: 100px;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 12px 16px;
  font-size: 14px;
  resize: vertical;
  outline: none;
  box-sizing: border-box;
  font-family: inherit;
  color: #333;
  background: #fafbfc;
}

.post-content-input:focus { border-color: #e8a0bf; }

/* 提示词编辑 */
.prompt-hint {
  font-size: 12px;
  color: #8f8a82;
  margin: 0 0 12px;
  line-height: 1.5;
}

.prompt-hint code {
  background: #f1ede8;
  padding: 1px 4px;
  border-radius: 4px;
  font-size: 11px;
}

.prompt-textarea {
  width: 100%;
  min-height: 200px;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 12px;
  font-size: 13px;
  resize: vertical;
  outline: none;
  box-sizing: border-box;
  font-family: inherit;
  color: #333;
  line-height: 1.5;
  background: #fafbfc;
}

.prompt-textarea:focus { border-color: #e8a0bf; }

.prompt-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.reset-btn, .clear-data-btn {
  flex: 1;
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 13px;
  cursor: pointer;
  border: 1px solid #e9ecef;
  background: #f8f9fa;
  color: #666;
  transition: all 0.2s;
}

.reset-btn:active, .clear-data-btn:active { background: #e9ecef; }
.clear-data-btn { color: #e74c3c; border-color: #fed7d7; }

/* 过渡 */
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(20px) scale(0.95); opacity: 0; }
</style>
