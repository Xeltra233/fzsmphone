<template>
  <div class="characters-page">
    <NavBar title="角色管理" :show-back="true" back-to="/">
      <template #right>
        <button class="nav-btn" @click="showCreateModal = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8" />
        <line x1="21" y1="21" x2="16.65" y2="16.65" />
      </svg>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索角色..."
        class="search-input"
      />
      <button v-if="searchQuery" class="clear-btn" @click="searchQuery = ''">✕</button>
    </div>

    <!-- 标签筛选 -->
    <div class="tags-bar" v-if="allTags.length > 0">
      <div class="tags-scroll">
        <button
          class="tag-chip"
          :class="{ active: selectedTag === '' }"
          @click="selectedTag = ''"
        >全部</button>
        <button
          v-for="tag in allTags"
          :key="tag"
          class="tag-chip"
          :class="{ active: selectedTag === tag }"
          @click="selectedTag = tag"
        >{{ tag }}</button>
      </div>
    </div>

    <!-- 角色列表 -->
    <div class="characters-list">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>

      <div v-else-if="filteredCharacters.length === 0" class="empty-state">
        <span class="empty-icon">👥</span>
        <p>{{ searchQuery ? '没有找到匹配的角色' : '还没有角色，创建一个吧' }}</p>
      </div>

      <div v-else class="character-grid">
        <div
          v-for="char in filteredCharacters"
          :key="char.id"
          class="character-card"
          @click="viewCharacter(char)"
        >
          <div class="card-avatar">
            <img v-if="char.avatar_url" :src="char.avatar_url" :alt="char.name" />
            <span v-else class="avatar-letter">{{ char.name[0] }}</span>
          </div>
          <div class="card-info">
            <h3 class="card-name">{{ char.name }}</h3>
            <p class="card-desc">{{ char.description || '暂无描述' }}</p>
            <div class="card-tags" v-if="char.tags && char.tags.length > 0">
              <span v-for="tag in char.tags.slice(0, 3)" :key="tag" class="mini-tag">{{ tag }}</span>
            </div>
          </div>
          <div class="card-actions">
            <button class="action-btn chat-btn" @click.stop="startChat(char)" title="开始聊天">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
              </svg>
            </button>
            <button class="action-btn edit-btn" @click.stop="editCharacter(char)" title="编辑">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
              </svg>
            </button>
            <button class="action-btn delete-btn" @click.stop="confirmDelete(char)" title="删除">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建角色弹窗 -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
        <div class="create-modal">
          <div class="modal-header">
            <span>创建角色</span>
            <button class="close-btn" @click="showCreateModal = false">✕</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>角色名称 *</label>
              <input v-model="newChar.name" type="text" placeholder="输入角色名称" class="form-input" />
            </div>
            <div class="form-group">
              <label>头像链接</label>
              <input v-model="newChar.avatar_url" type="text" placeholder="https://..." class="form-input" />
            </div>
            <div class="form-group">
              <label>简介</label>
              <textarea v-model="newChar.description" placeholder="角色简介..." class="form-textarea" rows="2"></textarea>
            </div>
            <div class="form-group">
              <label>性格特征</label>
              <textarea v-model="newChar.personality" placeholder="描述角色的性格..." class="form-textarea" rows="2"></textarea>
            </div>
            <div class="form-group">
              <label>系统提示词</label>
              <textarea v-model="newChar.system_prompt" placeholder="角色的系统提示词..." class="form-textarea" rows="3"></textarea>
            </div>
            <div class="form-group">
              <label>开场白</label>
              <textarea v-model="newChar.greeting" placeholder="角色的第一句话..." class="form-textarea" rows="2"></textarea>
            </div>
            <div class="form-group">
              <label>标签（逗号分隔）</label>
              <input v-model="tagsInput" type="text" placeholder="可爱, 温柔, 猫娘" class="form-input" />
            </div>
            <div class="form-group toggle-group">
              <label>公开角色</label>
              <label class="toggle">
                <input type="checkbox" v-model="newChar.is_public" />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn-cancel" @click="showCreateModal = false">取消</button>
            <button class="btn-create" @click="handleCreate" :disabled="!newChar.name?.trim()">创建</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 删除确认 -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
        <div class="confirm-dialog">
          <p>确定要删除角色「{{ deleteTarget.name }}」吗？</p>
          <p class="confirm-sub">此操作不可撤销</p>
          <div class="confirm-actions">
            <button class="btn-cancel" @click="deleteTarget = null">取消</button>
            <button class="btn-delete" @click="handleDelete">删除</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore, type Character } from '@/stores/chat'
import NavBar from '@/components/common/NavBar.vue'

const router = useRouter()
const chatStore = useChatStore()

const searchQuery = ref('')
const selectedTag = ref('')
const showCreateModal = ref(false)
const deleteTarget = ref<Character | null>(null)
const tagsInput = ref('')
const loading = ref(false)

const newChar = ref<Partial<Character>>({
  name: '',
  avatar_url: '',
  description: '',
  personality: '',
  system_prompt: '',
  greeting: '',
  is_public: false,
  tags: [],
})

const allTags = computed(() => {
  const tagSet = new Set<string>()
  chatStore.characters.forEach(c => {
    c.tags?.forEach(t => tagSet.add(t))
  })
  return Array.from(tagSet).sort()
})

const filteredCharacters = computed(() => {
  let list = chatStore.characters
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(c =>
      c.name.toLowerCase().includes(q) ||
      c.description?.toLowerCase().includes(q) ||
      c.tags?.some(t => t.toLowerCase().includes(q))
    )
  }
  if (selectedTag.value) {
    list = list.filter(c => c.tags?.includes(selectedTag.value))
  }
  return list
})

function viewCharacter(char: Character) {
  router.push(`/character/edit/${char.id}`)
}

function editCharacter(char: Character) {
  router.push(`/character/edit/${char.id}`)
}

async function startChat(char: Character) {
  // Check if conversation exists for this character
  const existing = chatStore.conversations.find(c => c.character_id === char.id)
  if (existing) {
    router.push(`/chat/${existing.id}`)
  } else {
    const conv = await chatStore.createConversation(char.id)
    if (conv) {
      router.push(`/chat/${conv.id}`)
    }
  }
}

function confirmDelete(char: Character) {
  deleteTarget.value = char
}

async function handleDelete() {
  if (!deleteTarget.value) return
  await chatStore.deleteCharacter(deleteTarget.value.id)
  deleteTarget.value = null
}

async function handleCreate() {
  if (!newChar.value.name?.trim()) return

  const tags = tagsInput.value
    .split(/[,，]/)
    .map(t => t.trim())
    .filter(Boolean)

  const result = await chatStore.createCharacter({
    ...newChar.value,
    tags,
  })

  if (result) {
    showCreateModal.value = false
    newChar.value = {
      name: '',
      avatar_url: '',
      description: '',
      personality: '',
      system_prompt: '',
      greeting: '',
      is_public: false,
      tags: [],
    }
    tagsInput.value = ''
  }
}

onMounted(async () => {
  loading.value = true
  await chatStore.fetchCharacters()
  loading.value = false
})
</script>

<style scoped>
.characters-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

/* 搜索栏 */
.search-bar {
  display: flex;
  align-items: center;
  margin: 8px 16px;
  padding: 0 12px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  height: 36px;
}

.search-icon {
  width: 16px;
  height: 16px;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  padding: 0 8px;
  font-size: 15px;
  color: var(--text-primary);
}

.search-input::placeholder {
  color: var(--text-tertiary);
}

.clear-btn {
  background: var(--bg-quaternary);
  border: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  font-size: 10px;
  color: var(--text-tertiary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 标签筛选 */
.tags-bar {
  padding: 4px 16px 8px;
}

.tags-scroll {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tags-scroll::-webkit-scrollbar {
  display: none;
}

.tag-chip {
  padding: 4px 12px;
  border-radius: 14px;
  border: 1px solid var(--separator);
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-size: 13px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.tag-chip.active {
  background: var(--brand-primary);
  color: #fff;
  border-color: var(--brand-primary);
}

/* 角色列表 */
.characters-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 16px 16px;
}

.characters-list::-webkit-scrollbar {
  display: none;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 60px 0;
  color: var(--text-tertiary);
  font-size: 14px;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--separator);
  border-top-color: var(--brand-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 60px 0;
}

.empty-icon {
  font-size: 48px;
}

.empty-state p {
  font-size: 14px;
  color: var(--text-tertiary);
}

/* 角色卡片网格 */
.character-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.character-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  padding: 14px;
  cursor: pointer;
  transition: background 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.character-card:active {
  background: var(--bg-tertiary);
}

.card-avatar {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  overflow: hidden;
  background: linear-gradient(135deg, #6C5CE7, #A29BFE);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-letter {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
}

.card-info {
  flex: 1;
  min-width: 0;
}

.card-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 2px;
}

.card-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0 0 4px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-tags {
  display: flex;
  gap: 4px;
}

.mini-tag {
  padding: 1px 6px;
  background: var(--bg-tertiary);
  border-radius: 4px;
  font-size: 11px;
  color: var(--text-tertiary);
}

.card-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn:active {
  transform: scale(0.9);
}

.chat-btn {
  color: var(--brand-primary);
}

.edit-btn {
  color: var(--color-orange);
}

.delete-btn {
  color: var(--color-red);
}

/* 创建弹窗 */
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

.create-modal {
  background: var(--bg-primary);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  width: 100%;
  max-height: 85%;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 0.5px solid var(--separator);
  flex-shrink: 0;
}

.close-btn {
  background: var(--bg-tertiary);
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.modal-body::-webkit-scrollbar {
  display: none;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 15px;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-input:focus {
  border-color: var(--brand-primary);
}

.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 15px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-textarea:focus {
  border-color: var(--brand-primary);
}

.toggle-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toggle-group label:first-child {
  margin-bottom: 0;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 0.5px solid var(--separator);
  flex-shrink: 0;
}

.btn-cancel {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: var(--radius-md);
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
}

.btn-create {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: var(--radius-md);
  background: var(--brand-primary);
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-create:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 删除确认 */
.confirm-dialog {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: 24px;
  margin: auto 32px;
  text-align: center;
}

.confirm-dialog p {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.confirm-sub {
  font-size: 13px !important;
  color: var(--text-tertiary) !important;
  font-weight: 400 !important;
  margin-bottom: 20px !important;
}

.confirm-actions {
  display: flex;
  gap: 12px;
}

.btn-delete {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: var(--radius-md);
  background: var(--color-red);
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
}

/* Toggle */
.toggle {
  position: relative;
  width: 51px;
  height: 31px;
  cursor: pointer;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
  position: absolute;
}

.toggle-slider {
  position: absolute;
  inset: 0;
  background: var(--bg-quaternary);
  border-radius: 31px;
  transition: background 0.3s;
}

.toggle-slider::after {
  content: '';
  position: absolute;
  left: 2px;
  top: 2px;
  width: 27px;
  height: 27px;
  background: #fff;
  border-radius: 50%;
  transition: transform 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle input:checked + .toggle-slider {
  background: var(--color-green);
}

.toggle input:checked + .toggle-slider::after {
  transform: translateX(20px);
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
