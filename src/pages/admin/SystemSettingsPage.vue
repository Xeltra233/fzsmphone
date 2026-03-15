<template>
  <div class="system-settings-page">
    <NavBar title="系统设置" :show-back="true" back-to="/profile" />

    <div class="page-content" v-if="authStore.isSuperAdmin || authStore.isAdmin">
      <!-- 标签切换 -->
      <div class="settings-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: isSystemTab }" 
          @click="activeTab = 'system'"
        >
          系统配置
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: isApiTab }" 
          @click="activeTab = 'api'"
        >
          API设置
        </button>
      </div>

<!-- 系统配置 -->
<template v-if="isSystemTab">
        <div v-if="!authStore.isSuperAdmin" class="permission-notice">
          <span class="notice-icon">ℹ️</span>
          您的账户为管理员，只能查看系统配置。如需修改，请联系超级管理员。
        </div>

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
                    :disabled="item.isExisting || !authStore.isSuperAdmin"
                  />
                  <button v-if="authStore.isSuperAdmin" class="delete-btn" @click="removeSetting(index)" title="删除">
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
                    :disabled="!authStore.isSuperAdmin"
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

<!-- API设置 -->
<template v-if="isApiTab">
          <div class="api-settings-section">
            <div v-if="apiSettingsLoading" class="loading-state">
              <div class="spinner"></div>
              <p>加载API设置...</p>
            </div>
            <template v-else>
              <!-- 超级管理员可以看到全局设置 -->
              <div v-if="authStore.isSuperAdmin && apiSettings?.global" class="settings-card">
                <div class="card-header">
                  <h3>全局API设置</h3>
                  <span class="badge">全局</span>
                </div>
                <div class="card-body">
                  <div class="input-group">
                    <label>API Key</label>
                    <input v-model="apiForm.globalApiKey" type="password" placeholder="留空则使用用户个人设置" class="setting-input" />
                  </div>
                  <div class="input-group">
                    <label>API URL</label>
                    <input v-model="apiForm.globalApiUrl" placeholder="https://api.openai.com/v1/chat/completions" class="setting-input" />
                  </div>
                  <div class="input-group">
                    <label>默认模型</label>
                    <input v-model="apiForm.globalModel" placeholder="gpt-4o-mini" class="setting-input" />
                  </div>
                </div>
              </div>

              <!-- 个人API设置 (所有用户) -->
              <div class="settings-card">
                <div class="card-header">
                  <h3>个人API设置</h3>
                  <span class="badge">个人</span>
                </div>
                <div class="card-body">
                  <div class="input-group">
                    <label>API Key</label>
                    <input v-model="apiForm.personalApiKey" type="password" placeholder="请输入您的API Key" class="setting-input" />
                  </div>
                  <div class="input-group">
                    <label>API URL</label>
                    <input v-model="apiForm.personalApiUrl" placeholder="https://api.openai.com/v1/chat/completions" class="setting-input" />
                  </div>
                  <div class="input-group">
                    <label>默认模型</label>
                    <input v-model="apiForm.personalModel" placeholder="gpt-4o-mini" class="setting-input" />
                  </div>
                </div>
              </div>

              <button class="save-btn" @click="saveApiSettings" :disabled="apiSaving">
                <template v-if="apiSaving">
                  <div class="btn-spinner"></div>
                  保存中...
                </template>
                <template v-else>
                  ▣ 保存API设置
                </template>
              </button>
            </template>
          </div>
        </template>
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

const activeTab = ref<'system' | 'api'>('system')

const isApiTab = computed(() => activeTab.value === 'api')
const isSystemTab = computed(() => activeTab.value === 'system')

const apiSettingsLoading = ref(false)
const apiSaving = ref(false)
const apiSettings = ref<Record<string, any> | null>(null)
const apiForm = ref({
  globalApiKey: '',
  globalApiUrl: '',
  globalModel: '',
  personalApiKey: '',
  personalApiUrl: '',
  personalModel: '',
})

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

async function fetchApiSettings() {
  apiSettingsLoading.value = true
  try {
    const res = await apiClient.get<Record<string, any>>('/api/settings/api') as Record<string, any>
    apiSettings.value = res
    if (authStore.isSuperAdmin && res.global) {
      apiForm.value.globalApiKey = res.global.api_key || ''
      apiForm.value.globalApiUrl = res.global.api_url || ''
      apiForm.value.globalModel = res.global.model || ''
    }
    if (res.personal) {
      apiForm.value.personalApiKey = res.personal.api_key || ''
      apiForm.value.personalApiUrl = res.personal.api_url || ''
      apiForm.value.personalModel = res.personal.model || ''
    }
  } catch (err: any) {
    showToast('加载API设置失败', 'error')
  } finally {
    apiSettingsLoading.value = false
  }
}

async function saveApiSettings() {
  apiSaving.value = true
  try {
    if (authStore.isSuperAdmin) {
      await apiClient.put('/api/settings/api', {
        api_key: apiForm.value.globalApiKey,
        api_url: apiForm.value.globalApiUrl,
        model: apiForm.value.globalModel,
        is_global: true,
      })
    }
    await apiClient.put('/api/settings/api', {
      api_key: apiForm.value.personalApiKey,
      api_url: apiForm.value.personalApiUrl,
      model: apiForm.value.personalModel,
      is_global: false,
    })
    showToast('API设置已保存')
  } catch (err: any) {
    showToast('保存API设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    apiSaving.value = false
  }
}

onMounted(() => {
  if (authStore.isAdmin) {
    fetchSettings()
    fetchApiSettings()
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

/* 标签切换 */
.settings-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.tab-btn {
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: rgba(91, 110, 245, 0.3);
  border-color: #5B6EF5;
  color: #fff;
}

/* 权限提示 */
.permission-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(91, 110, 245, 0.15);
  border: 1px solid rgba(91, 110, 245, 0.3);
  border-radius: 12px;
  margin-bottom: 16px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.notice-icon {
  font-size: 16px;
}

/* API设置区域 */
.api-settings-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.settings-card {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.card-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.badge {
  padding: 4px 10px;
  background: rgba(91, 110, 245, 0.2);
  border-radius: 12px;
  font-size: 11px;
  color: #5B6EF5;
}

.card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-group label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

.setting-input {
  padding: 10px 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  outline: none;
}

.setting-input:focus {
  border-color: #5B6EF5;
}

.setting-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
