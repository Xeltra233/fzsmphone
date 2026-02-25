<template>
  <div class="worldbook-page">
    <NavBar title="世界书" back-to="/friends" />

    <div class="page-content">
      <!-- 添加按钮 -->
      <div class="add-bar">
        <button class="add-btn" @click="showForm = true" v-if="!showForm">
          <span class="add-icon">+</span> 新建条目
        </button>
      </div>

      <!-- 新建/编辑表单 -->
      <div class="entry-form" v-if="showForm">
        <div class="form-header">
          <span>{{ editingId ? '编辑条目' : '新建条目' }}</span>
          <button class="close-btn" @click="resetForm">✕</button>
        </div>
        <input v-model="form.key" placeholder="条目名称 / 关键字" class="form-input" />
        <textarea v-model="form.content" placeholder="条目内容（当关键字被触发时插入的文本）..." class="form-textarea" rows="4" />
        <input v-model="keywordsInput" placeholder="触发关键词（逗号分隔）" class="form-input" />
        <div class="form-row">
          <label class="form-check">
            <input type="checkbox" v-model="form.is_enabled" />
            <span>启用</span>
          </label>
          <div class="priority-input">
            <span class="priority-label">优先级</span>
            <input type="number" v-model.number="form.priority" min="0" max="100" class="form-input-small" />
          </div>
        </div>
        <div class="form-actions">
          <button class="btn-cancel" @click="resetForm">取消</button>
          <button class="btn-save" @click="handleSave" :disabled="!form.key">
            {{ editingId ? '更新' : '创建' }}
          </button>
        </div>
      </div>

      <!-- 加载状态 -->
      <div class="loading" v-if="store.loading">加载中...</div>

      <!-- 空状态 -->
      <div class="empty" v-else-if="store.entries.length === 0 && !showForm">
        <div class="empty-icon">📚</div>
        <p>还没有世界书条目</p>
        <p class="empty-hint">世界书用于在聊天中自动注入背景设定信息</p>
      </div>

      <!-- 条目列表 -->
      <div class="entry-list" v-else>
        <div
          v-for="entry in store.entries"
          :key="entry.id"
          class="entry-card"
          :class="{ disabled: !entry.is_enabled }"
        >
          <div class="entry-header">
            <div class="entry-key">
              <span class="key-icon">📖</span>
              {{ entry.key }}
            </div>
            <div class="entry-priority">P{{ entry.priority }}</div>
          </div>
          <div class="entry-content">{{ entry.content }}</div>
          <div class="entry-keywords" v-if="entry.keywords.length > 0">
            <span class="keyword-tag" v-for="kw in entry.keywords" :key="kw">{{ kw }}</span>
          </div>
          <div class="entry-footer">
            <button class="toggle-btn" :class="{ active: entry.is_enabled }" @click="handleToggle(entry.id)">
              {{ entry.is_enabled ? '✅ 已启用' : '❌ 已禁用' }}
            </button>
            <div class="entry-actions">
              <button class="action-btn" @click="startEdit(entry)">✏️</button>
              <button class="action-btn" @click="handleDelete(entry.id)">🗑️</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useWorldBookStore } from '@/stores/worldbook'
import type { WorldBookEntry, WorldBookInput } from '@/api/types'

const store = useWorldBookStore()

const showForm = ref(false)
const editingId = ref<number | null>(null)
const keywordsInput = ref('')
const form = ref<WorldBookInput>({
  key: '',
  content: '',
  keywords: [],
  is_enabled: true,
  priority: 0,
})

onMounted(() => {
  store.fetchEntries()
})

function startEdit(entry: WorldBookEntry) {
  editingId.value = entry.id
  form.value = {
    key: entry.key,
    content: entry.content,
    keywords: [...entry.keywords],
    is_enabled: entry.is_enabled,
    priority: entry.priority,
  }
  keywordsInput.value = entry.keywords.join(', ')
  showForm.value = true
}

function resetForm() {
  showForm.value = false
  editingId.value = null
  keywordsInput.value = ''
  form.value = { key: '', content: '', keywords: [], is_enabled: true, priority: 0 }
}

async function handleSave() {
  if (!form.value.key) return
  // 解析关键词
  form.value.keywords = keywordsInput.value
    .split(',')
    .map(k => k.trim())
    .filter(k => k.length > 0)

  try {
    if (editingId.value) {
      await store.updateEntry(editingId.value, form.value)
    } else {
      await store.createEntry(form.value)
    }
    resetForm()
  } catch (e) {
    console.error('保存条目失败', e)
  }
}

async function handleDelete(id: number) {
  if (!confirm('确定删除此条目？')) return
  try {
    await store.deleteEntry(id)
  } catch (e) {
    console.error('删除条目失败', e)
  }
}

async function handleToggle(id: number) {
  try {
    await store.toggleEnabled(id)
  } catch (e) {
    console.error('切换状态失败', e)
  }
}
</script>

<style scoped>
.worldbook-page {
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

.add-bar { margin-bottom: 16px; }

.add-btn {
  width: 100%;
  padding: 12px;
  border: 2px dashed var(--border-color);
  border-radius: 12px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 15px;
  cursor: pointer;
}

.add-btn:active { background: var(--bg-secondary); }
.add-icon { font-size: 18px; margin-right: 4px; }

/* 表单 */
.entry-form {
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

.form-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.form-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}

.priority-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.priority-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.form-input-small {
  width: 60px;
  padding: 6px 8px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  text-align: center;
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

.btn-save:disabled { opacity: 0.5; cursor: not-allowed; }

/* 状态 */
.loading, .empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-hint { font-size: 13px; margin-top: 4px; opacity: 0.7; }

/* 条目卡片 */
.entry-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.entry-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 14px;
  transition: opacity 0.2s;
}

.entry-card.disabled {
  opacity: 0.5;
}

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.entry-key {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.key-icon { font-size: 14px; }

.entry-priority {
  background: var(--accent-blue);
  color: white;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
}

.entry-content {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-bottom: 8px;
  max-height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.entry-keywords {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 8px;
}

.keyword-tag {
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 6px;
}

.entry-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toggle-btn {
  background: none;
  border: none;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.15s;
}

.toggle-btn:active { background: var(--bg-primary); }
.toggle-btn.active { color: var(--accent-green, #22c55e); }

.entry-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  background: none;
  border: none;
  font-size: 16px;
  padding: 4px 6px;
  cursor: pointer;
  border-radius: 6px;
}

.action-btn:active { background: var(--bg-primary); }
</style>
