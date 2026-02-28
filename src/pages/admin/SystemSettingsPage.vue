<template>
  <div class="system-settings-page">
    <NavBar title="系统设置" :show-back="true" back-to="/profile" />

    <div class="page-content" v-if="authStore.isAdmin">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载设置中...</p>
      </div>

      <template v-else>
        <!-- 设置列表 -->
        <div class="settings-section" v-if="settingsList.length > 0">
          <div class="section-header">
            <h3>配置项 ({{ settingsList.length }})</h3>
          </div>
          <div class="settings-list">
            <div
              class="setting-card"
              v-for="(item, index) in settingsList"
              :key="index"
            >
              <div class="setting-header">
                <input
                  class="setting-key"
                  v-model="item.key"
                  placeholder="键名 (key)"
                  :disabled="item.isExisting"
                />
                <button class="delete-btn" @click="removeSetting(index)" title="删除">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M18 6L6 18M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="setting-value-wrap">
                <textarea
                  class="setting-value"
                  v-model="item.valueStr"
                  placeholder='值 (JSON格式，如 "hello" 或 {"a":1})'
                  rows="3"
                  @input="validateJson(item)"
                ></textarea>
                <span class="json-hint" :class="{ error: item.jsonError }">
                  {{ item.jsonError || 'JSON 格式正确' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <span class="empty-icon">⚙️</span>
          <p>暂无系统设置</p>
        </div>

        <!-- 操作按钮 -->
        <div class="actions">
          <button class="add-btn" @click="addSetting">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14" />
            </svg>
            添加配置项
          </button>
          <button class="save-btn" @click="saveSettings" :disabled="saving || hasJsonErrors">
            <template v-if="saving">
              <div class="btn-spinner"></div>
              保存中...
            </template>
            <template v-else>
              ▣ 保存所有设置
            </template>
          </button>
        </div>
      </template>

      <!-- 提示信息 -->
      <div class="toast" v-if="toast.show" :class="toast.type">
        {{ toast.message }}
      </div>
    </div>

    <div v-else class="no-access">
      <p>⛔ 无权限访问</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import apiClient from '@/api/client'
import NavBar from '@/components/common/NavBar.vue'

const authStore = useAuthStore()

interface SettingItem {
  key: string
  valueStr: string
  jsonError: string
  isExisting: boolean
}

const loading = ref(true)
const saving = ref(false)
const settingsList = ref<SettingItem[]>([])
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const hasJsonErrors = computed(() => {
  return settingsList.value.some(item => item.jsonError !== '' || item.key.trim() === '')
})

function validateJson(item: SettingItem) {
  if (!item.valueStr.trim()) {
    item.jsonError = '值不能为空'
    return
  }
  try {
    JSON.parse(item.valueStr)
    item.jsonError = ''
  } catch {
    item.jsonError = 'JSON 格式错误'
  }
}

function addSetting() {
  settingsList.value.push({
    key: '',
    valueStr: '""',
    jsonError: '',
    isExisting: false,
  })
}

function removeSetting(index: number) {
  settingsList.value.splice(index, 1)
}

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

async function fetchSettings() {
  loading.value = true
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    settingsList.value = []
    for (const [key, value] of Object.entries(data)) {
      settingsList.value.push({
        key,
        valueStr: JSON.stringify(value, null, 2),
        jsonError: '',
        isExisting: true,
      })
    }
  } catch (err: any) {
    showToast('加载设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  if (hasJsonErrors.value) {
    showToast('请先修复 JSON 格式错误', 'error')
    return
  }

  // Check for duplicate keys
  const keys = settingsList.value.map(i => i.key.trim())
  const uniqueKeys = new Set(keys)
  if (keys.some(k => k === '')) {
    showToast('键名不能为空', 'error')
    return
  }
  if (uniqueKeys.size !== keys.length) {
    showToast('存在重复的键名', 'error')
    return
  }

  saving.value = true
  try {
    const payload: Record<string, any> = {}
    for (const item of settingsList.value) {
      payload[item.key.trim()] = JSON.parse(item.valueStr)
    }
    await apiClient.put('/api/settings', payload)
    showToast('设置已保存')
    // Mark all as existing
    settingsList.value.forEach(item => {
      item.isExisting = true
    })
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (authStore.isAdmin) {
    fetchSettings()
  }
})
</script>

<style scoped>
.system-settings-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #0a0a0a;
  color: #fff;
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  position: relative;
}

.page-content::-webkit-scrollbar {
  display: none;
}

/* 加载 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  gap: 16px;
  color: rgba(255, 255, 255, 0.5);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #5B6EF5;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  gap: 12px;
}

.empty-icon {
  font-size: 48px;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.4);
  font-size: 15px;
}

/* Section */
.settings-section {
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.5);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
}

/* 设置卡片 */
.settings-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.setting-card {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 14px;
}

.setting-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.setting-key {
  flex: 1;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 14px;
  font-weight: 600;
  color: #5B6EF5;
  font-family: 'SF Mono', 'Fira Code', monospace;
  outline: none;
}

.setting-key:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.setting-key:focus {
  border-color: #5B6EF5;
}

.delete-btn {
  width: 32px;
  height: 32px;
  background: rgba(255, 59, 48, 0.15);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
}

.delete-btn svg {
  width: 16px;
  height: 16px;
  color: #FF3B30;
}

.delete-btn:hover {
  background: rgba(255, 59, 48, 0.3);
}

.setting-value-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-value {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 13px;
  color: #fff;
  font-family: 'SF Mono', 'Fira Code', monospace;
  resize: vertical;
  outline: none;
  min-height: 60px;
  box-sizing: border-box;
}

.setting-value:focus {
  border-color: #5B6EF5;
}

.json-hint {
  font-size: 11px;
  color: rgba(52, 199, 89, 0.8);
  padding-left: 4px;
}

.json-hint.error {
  color: #FF3B30;
}

/* 操作按钮 */
.actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 16px;
  padding-bottom: 32px;
}

.add-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px dashed rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn svg {
  width: 18px;
  height: 18px;
}

.add-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
  color: #fff;
}

.save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  border: none;
  border-radius: 12px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.save-btn:not(:disabled):hover {
  filter: brightness(1.1);
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Toast */
.toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  z-index: 1000;
  animation: fadeInUp 0.3s ease;
}

.toast.success {
  background: rgba(52, 199, 89, 0.9);
  color: #fff;
}

.toast.error {
  background: rgba(255, 59, 48, 0.9);
  color: #fff;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/* 无权限 */
.no-access {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.4);
  font-size: 16px;
}
</style>
