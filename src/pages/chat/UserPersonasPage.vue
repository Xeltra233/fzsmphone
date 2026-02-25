<template>
  <div class="personas-page">
    <NavBar title="用户人设" back-to="/friends" />

    <div class="page-content">
      <!-- 添加按钮 -->
      <div class="add-bar">
        <button class="add-btn" @click="showForm = true" v-if="!showForm">
          <span class="add-icon">+</span> 新建人设
        </button>
      </div>

      <!-- 新建/编辑表单 -->
      <div class="persona-form" v-if="showForm">
        <div class="form-header">
          <span>{{ editingId ? '编辑人设' : '新建人设' }}</span>
          <button class="close-btn" @click="resetForm">✕</button>
        </div>
        <input v-model="form.name" placeholder="人设名称" class="form-input" />
        <textarea v-model="form.description" placeholder="人设描述..." class="form-textarea" rows="3" />
        <input v-model="form.avatar_url" placeholder="头像URL（可选）" class="form-input" />
        <label class="form-check">
          <input type="checkbox" v-model="form.is_default" />
          <span>设为默认人设</span>
        </label>
        <div class="form-actions">
          <button class="btn-cancel" @click="resetForm">取消</button>
          <button class="btn-save" @click="handleSave" :disabled="!form.name">
            {{ editingId ? '更新' : '创建' }}
          </button>
        </div>
      </div>

      <!-- 加载状态 -->
      <div class="loading" v-if="store.loading">加载中...</div>

      <!-- 空状态 -->
      <div class="empty" v-else-if="store.personas.length === 0 && !showForm">
        <div class="empty-icon">🎭</div>
        <p>还没有创建人设</p>
        <p class="empty-hint">人设是你在角色扮演中的身份设定</p>
      </div>

      <!-- 人设列表 -->
      <div class="persona-list" v-else>
        <div
          v-for="persona in store.personas"
          :key="persona.id"
          class="persona-card"
          :class="{ 'is-default': persona.is_default }"
        >
          <div class="persona-avatar">
            <img v-if="persona.avatar_url" :src="persona.avatar_url" alt="" />
            <span v-else class="avatar-placeholder">{{ persona.name[0] }}</span>
            <span v-if="persona.is_default" class="default-badge">默认</span>
          </div>
          <div class="persona-info">
            <div class="persona-name">{{ persona.name }}</div>
            <div class="persona-desc">{{ persona.description || '暂无描述' }}</div>
          </div>
          <div class="persona-actions">
            <button class="action-btn edit" @click="startEdit(persona)">✏️</button>
            <button class="action-btn delete" @click="handleDelete(persona.id)">🗑️</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { usePersonasStore } from '@/stores/personas'
import type { Persona, PersonaInput } from '@/api/types'

const store = usePersonasStore()

const showForm = ref(false)
const editingId = ref<number | null>(null)
const form = ref<PersonaInput>({
  name: '',
  description: '',
  avatar_url: '',
  is_default: false,
})

onMounted(() => {
  store.fetchPersonas()
})

function startEdit(persona: Persona) {
  editingId.value = persona.id
  form.value = {
    name: persona.name,
    description: persona.description,
    avatar_url: persona.avatar_url,
    is_default: persona.is_default,
  }
  showForm.value = true
}

function resetForm() {
  showForm.value = false
  editingId.value = null
  form.value = { name: '', description: '', avatar_url: '', is_default: false }
}

async function handleSave() {
  if (!form.value.name) return
  try {
    if (editingId.value) {
      await store.updatePersona(editingId.value, form.value)
    } else {
      await store.createPersona(form.value)
    }
    resetForm()
  } catch (e) {
    console.error('保存人设失败', e)
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定删除此人设？')) return
  try {
    await store.deletePersona(id)
  } catch (e) {
    console.error('删除人设失败', e)
  }
}
</script>

<style scoped>
.personas-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.add-bar {
  margin-bottom: 16px;
}

.add-btn {
  width: 100%;
  padding: 12px;
  border: 2px dashed var(--border-color);
  border-radius: 12px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn:active {
  background: var(--bg-secondary);
}

.add-icon {
  font-size: 18px;
  margin-right: 4px;
}

/* 表单 */
.persona-form {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: var(--text-secondary);
  cursor: pointer;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  margin-bottom: 8px;
  box-sizing: border-box;
}

.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  margin-bottom: 8px;
  resize: vertical;
  box-sizing: border-box;
}

.form-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 12px;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.btn-cancel {
  padding: 8px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
}

.btn-save {
  padding: 8px 20px;
  border: none;
  border-radius: 8px;
  background: var(--accent-blue);
  color: white;
  font-size: 14px;
  cursor: pointer;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 状态 */
.loading, .empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.empty-hint {
  font-size: 13px;
  margin-top: 4px;
  opacity: 0.7;
}

/* 人设卡片 */
.persona-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.persona-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
  transition: transform 0.15s;
}

.persona-card.is-default {
  border: 1px solid var(--accent-blue);
}

.persona-card:active {
  transform: scale(0.98);
}

.persona-avatar {
  position: relative;
  width: 48px;
  height: 48px;
  flex-shrink: 0;
}

.persona-avatar img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--accent-blue), var(--accent-purple, #a855f7));
  color: white;
  font-size: 20px;
  font-weight: 600;
}

.default-badge {
  position: absolute;
  bottom: -2px;
  right: -4px;
  background: var(--accent-blue);
  color: white;
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 8px;
}

.persona-info {
  flex: 1;
  min-width: 0;
}

.persona-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-primary);
}

.persona-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.persona-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.action-btn {
  background: none;
  border: none;
  font-size: 16px;
  padding: 6px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.15s;
}

.action-btn:active {
  background: var(--bg-primary);
}
</style>
