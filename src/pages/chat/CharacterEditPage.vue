<template>
  <div class="character-edit-page">
    <NavBar :title="isNew ? '创建角色' : '编辑角色'" :show-back="true" back-to="/characters">
      <template #right>
        <button class="nav-btn save-btn" @click="handleSave" :disabled="!canSave">
          <span>保存</span>
        </button>
      </template>
    </NavBar>

    <div class="edit-content" v-if="!loading">
      <!-- 头像预览 -->
      <div class="avatar-preview">
        <div class="avatar-circle">
          <img v-if="form.avatar_url" :src="form.avatar_url" :alt="form.name" />
          <span v-else class="avatar-letter">{{ form.name?.[0] || '?' }}</span>
        </div>
        <button class="change-avatar-btn" @click="showAvatarInput = !showAvatarInput">更换头像</button>
      </div>

      <div v-if="showAvatarInput" class="avatar-url-input">
        <input
          v-model="form.avatar_url"
          type="text"
          placeholder="输入头像图片链接..."
          class="form-input"
        />
      </div>

      <!-- 基本信息 -->
      <div class="section">
        <div class="section-title">基本信息</div>
        <div class="form-card">
          <div class="form-row">
            <label>名称</label>
            <input v-model="form.name" type="text" placeholder="角色名称" class="row-input" />
          </div>
          <div class="form-row">
            <label>简介</label>
            <input v-model="form.description" type="text" placeholder="一句话介绍" class="row-input" />
          </div>
          <div class="form-row">
            <label>标签</label>
            <input v-model="tagsInput" type="text" placeholder="逗号分隔" class="row-input" />
          </div>
          <div class="form-row toggle-row">
            <label>公开</label>
            <label class="toggle">
              <input type="checkbox" v-model="form.is_public" />
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- 性格设定 -->
      <div class="section">
        <div class="section-title">性格设定</div>
        <div class="form-card">
          <div class="form-block">
            <label>性格特征</label>
            <textarea
              v-model="form.personality"
              placeholder="描述角色的性格特点、说话方式等..."
              class="block-textarea"
              rows="4"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- 系统提示词 -->
      <div class="section">
        <div class="section-title">系统提示词</div>
        <div class="form-card">
          <div class="form-block">
            <label>System Prompt</label>
            <textarea
              v-model="form.system_prompt"
              placeholder="定义角色行为的系统级指令..."
              class="block-textarea code-font"
              rows="6"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- 开场白 -->
      <div class="section">
        <div class="section-title">开场白</div>
        <div class="form-card">
          <div class="form-block">
            <label>Greeting</label>
            <textarea
              v-model="form.greeting"
              placeholder="角色发送的第一条消息..."
              class="block-textarea"
              rows="3"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- 预览 -->
      <div class="section" v-if="form.greeting">
        <div class="section-title">消息预览</div>
        <div class="preview-card">
          <div class="preview-avatar">
            <img v-if="form.avatar_url" :src="form.avatar_url" />
            <span v-else>{{ form.name?.[0] || '?' }}</span>
          </div>
          <div class="preview-bubble">
            <div class="preview-name">{{ form.name || '角色' }}</div>
            <div class="preview-text">{{ form.greeting }}</div>
          </div>
        </div>
      </div>

      <!-- 危险区域 -->
      <div class="section" v-if="!isNew">
        <div class="section-title danger-title">危险操作</div>
        <div class="form-card">
          <button class="danger-btn" @click="confirmDelete">删除此角色</button>
        </div>
      </div>

      <div class="bottom-spacer"></div>
    </div>

    <div v-else class="loading-state">
      <div class="spinner"></div>
    </div>

    <!-- 删除确认 -->
    <Teleport to="#phone-overlay">
      <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="showDeleteConfirm = false">
        <div class="confirm-dialog">
          <p>确定要删除角色「{{ form.name }}」吗？</p>
          <p class="confirm-sub">相关的对话记录也将被清除</p>
          <div class="confirm-actions">
            <button class="btn-cancel" @click="showDeleteConfirm = false">取消</button>
            <button class="btn-delete" @click="handleDelete">删除</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore, type Character } from '@/stores/chat'
import NavBar from '@/components/common/NavBar.vue'

const route = useRoute()
const router = useRouter()
const chatStore = useChatStore()

const loading = ref(false)
const showAvatarInput = ref(false)
const showDeleteConfirm = ref(false)
const tagsInput = ref('')

const form = ref({
  name: '',
  avatar_url: '',
  description: '',
  personality: '',
  system_prompt: '',
  greeting: '',
  is_public: false,
})

const charId = computed(() => {
  const id = route.params.id
  return id ? Number(id) : null
})

const isNew = computed(() => !charId.value)

const canSave = computed(() => form.value.name.trim().length > 0)

async function loadCharacter() {
  if (!charId.value) return
  loading.value = true

  // Ensure characters are loaded
  if (chatStore.characters.length === 0) {
    await chatStore.fetchCharacters()
  }

  const char = chatStore.characters.find(c => c.id === charId.value)
  if (char) {
    form.value = {
      name: char.name,
      avatar_url: char.avatar_url || '',
      description: char.description || '',
      personality: char.personality || '',
      system_prompt: char.system_prompt || '',
      greeting: char.greeting || '',
      is_public: char.is_public || false,
    }
    tagsInput.value = char.tags?.join(', ') || ''
  }

  loading.value = false
}

async function handleSave() {
  if (!canSave.value) return

  const tags = tagsInput.value
    .split(/[,，]/)
    .map(t => t.trim())
    .filter(Boolean)

  const data: Partial<Character> = {
    ...form.value,
    tags,
  }

  if (isNew.value) {
    const result = await chatStore.createCharacter(data)
    if (result) {
      router.replace(`/character/edit/${result.id}`)
    }
  } else {
    await chatStore.updateCharacter(charId.value!, data)
    router.back()
  }
}

function confirmDelete() {
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!charId.value) return
  await chatStore.deleteCharacter(charId.value)
  showDeleteConfirm.value = false
  router.replace('/characters')
}

onMounted(() => {
  loadCharacter()
})
</script>

<style scoped>
.character-edit-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

.edit-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.edit-content::-webkit-scrollbar {
  display: none;
}

/* 头像预览 */
.avatar-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.avatar-circle {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #6C5CE7, #A29BFE);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(108, 92, 231, 0.3);
}

.avatar-circle img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-letter {
  font-size: 36px;
  font-weight: 700;
  color: #fff;
}

.change-avatar-btn {
  border: none;
  background: none;
  color: var(--brand-primary);
  font-size: 14px;
  cursor: pointer;
  padding: 4px 8px;
}

.avatar-url-input {
  margin: 0 0 16px;
}

/* Section */
.section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 4px 8px;
}

.danger-title {
  color: var(--color-red);
}

.form-card {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

/* 行内表单 */
.form-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
}

.form-row:not(:last-child) {
  border-bottom: 0.5px solid var(--separator);
}

.form-row label {
  width: 60px;
  font-size: 15px;
  color: var(--text-primary);
  flex-shrink: 0;
}

.row-input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  font-size: 15px;
  color: var(--text-primary);
  text-align: right;
  padding: 0;
}

.row-input::placeholder {
  color: var(--text-quaternary);
}

.toggle-row {
  justify-content: space-between;
}

/* 块级表单 */
.form-block {
  padding: 12px 16px;
}

.form-block label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.block-textarea {
  width: 100%;
  border: 1px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.5;
  padding: 10px 12px;
  outline: none;
  resize: vertical;
  font-family: inherit;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.block-textarea:focus {
  border-color: var(--brand-primary);
}

.code-font {
  font-family: 'SF Mono', 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

/* 预览 */
.preview-card {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  padding: 14px;
  display: flex;
  gap: 10px;
}

.preview-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #6C5CE7, #A29BFE);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.preview-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-bubble {
  flex: 1;
  background: var(--bg-tertiary);
  border-radius: 16px 16px 16px 4px;
  padding: 10px 14px;
}

.preview-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--brand-primary);
  margin-bottom: 4px;
}

.preview-text {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.4;
  white-space: pre-wrap;
}

/* 危险按钮 */
.danger-btn {
  width: 100%;
  padding: 14px;
  border: none;
  background: none;
  color: var(--color-red);
  font-size: 16px;
  cursor: pointer;
  text-align: center;
}

.danger-btn:active {
  background: var(--bg-tertiary);
}

/* Form input */
.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 15px;
  outline: none;
  box-sizing: border-box;
}

.form-input:focus {
  border-color: var(--brand-primary);
}

/* Loading */
.loading-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
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

.bottom-spacer {
  height: 40px;
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

/* 确认弹窗 */
.modal-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.confirm-dialog {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: 24px;
  margin: 0 32px;
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

/* Save btn */
.save-btn {
  display: flex;
  align-items: center;
  border: none;
  background: none;
  color: var(--brand-primary);
  font-size: 17px;
  font-weight: 600;
  cursor: pointer;
  padding: 6px;
}

.save-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
