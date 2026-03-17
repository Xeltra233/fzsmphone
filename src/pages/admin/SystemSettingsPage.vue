<template>
  <div class="system-settings-page">
    <NavBar title="系统设置" :show-back="true" back-to="/profile" />

    <div class="page-content" v-if="authStore.isSuperAdmin || authStore.isAdmin">
      <!-- 标签切换 -->
      <div class="settings-tabs">
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'credits' }"
          @click="activeTab = 'credits'"
        >
          额度设置
        </button>
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'system' }"
          @click="activeTab = 'system'"
        >
          系统配置
        </button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'api' }"
  @click="activeTab = 'api'"
>
  API设置
</button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'oauth' }"
  @click="activeTab = 'oauth'"
>
  OAuth配置
</button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'contact' }"
  @click="activeTab = 'contact'"
>
  联系方式
</button>
</div>

      <!-- 额度设置 -->
      <template v-if="activeTab === 'credits'">
        <div v-if="!authStore.isSuperAdmin" class="permission-notice">
          <span class="notice-icon">ℹ️</span>
          您的账户为管理员，只能查看额度设置。如需修改，请联系超级管理员。
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>

        <template v-else>
          <div class="settings-section">
            <div class="section-header">
              <h3>签到设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">每日签到奖励</span>
                  <span class="label-desc">用户每日签到可获得的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.signin_daily_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">连续签到加成</span>
                  <span class="label-desc">连续签到时每日额外增加的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.signin_streak_bonus"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度/天</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">签到功能</span>
                  <span class="label-desc">是否启用每日签到功能</span>
                </div>
                <label class="toggle">
                  <input
                    type="checkbox"
                    v-model="creditSettings.signin_enabled"
                    :disabled="!authStore.isSuperAdmin"
                  />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <div class="section-header">
              <h3>邀请设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">邀请奖励额度</span>
                  <span class="label-desc">成功邀请好友后奖励的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.invite_reward_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">邀请功能</span>
                  <span class="label-desc">是否启用邀请好友功能</span>
                </div>
                <label class="toggle">
                  <input
                    type="checkbox"
                    v-model="creditSettings.invite_enabled"
                    :disabled="!authStore.isSuperAdmin"
                  />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <div class="section-header">
              <h3>新用户设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">新用户初始额度</span>
                  <span class="label-desc">新注册用户自动获得的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.default_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>
            </div>
          </div>

          <button
            class="save-btn"
            @click="saveCreditSettings"
            :disabled="saving || !authStore.isSuperAdmin"
          >
            <template v-if="saving">
              <div class="btn-spinner"></div>
              保存中...
            </template>
            <template v-else>
              ▣ 保存设置
            </template>
          </button>
        </template>
      </template>

<!-- 系统配置 -->
  <template v-if="activeTab === 'system'">
    <div v-if="!authStore.isSuperAdmin" class="permission-notice">
      <span class="notice-icon">ℹ️</span>
      您的账户为管理员，只能查看系统配置。如需修改，请联系超级管理员。
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>基础设置</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">应用名称</span>
            <span class="label-desc">显示在登录页和状态栏</span>
          </div>
          <input
            v-model="systemSettings.app_name"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈小手机"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">网页标题</span>
            <span class="label-desc">浏览器标签页显示的标题</span>
          </div>
          <input
            v-model="systemSettings.app_title"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈小手机"
          />
        </div>
      </div>
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>公告提示</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">公告内容</span>
            <span class="label-desc">显示在登录页顶部，支持换行</span>
          </div>
          <textarea
            v-model="systemSettings.announcement"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="欢迎使用..."
            rows="3"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">每日提示</span>
            <span class="label-desc">显示在额度页面的小技巧</span>
          </div>
          <textarea
            v-model="systemSettings.tips"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="每日签到可获得额外额度..."
            rows="2"
          />
        </div>
      </div>
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>版权声明</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">免责声明</span>
            <span class="label-desc">显示在页面底部，禁止商业声明</span>
          </div>
          <textarea
            v-model="systemSettings.disclaimer"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈保留所有权益..."
            rows="3"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">群二维码路径</span>
            <span class="label-desc">放在public目录下，如 qun_qrcode.jpg</span>
          </div>
          <input
            v-model="systemSettings.qun_qrcode"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
placeholder="qun_qrcode.jpg"
      />
      <div class="file-upload-area" v-if="authStore.isSuperAdmin">
        <input type="file" id="qrcode-upload" accept="image/*" @change="handleQrcodeUpload" hidden />
        <label for="qrcode-upload" class="upload-btn">上传二维码图片</label>
        <span v-if="uploadingQrcode" class="upload-status">上传中...</span>
      </div>
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">网站图标</span>
        <span class="label-desc">浏览器标签页图标，显示在浏览器标签上</span>
      </div>
      <div class="favicon-preview">
        <img v-if="systemSettings.favicon" :src="systemSettings.favicon" alt="favicon" class="favicon-img" @error="faviconError = true" />
        <span v-else class="favicon-placeholder">未设置</span>
      </div>
      <div class="file-upload-area" v-if="authStore.isSuperAdmin">
        <input type="file" id="favicon-upload" accept="image/*,.ico" @change="handleFaviconUpload" hidden />
        <label for="favicon-upload" class="upload-btn">上传图标</label>
        <span v-if="uploadingFavicon" class="upload-status">上传中...</span>
      </div>
      <input
        v-model="systemSettings.favicon"
        class="setting-input full-width"
        :disabled="!authStore.isSuperAdmin"
        placeholder="/icon.png"
      />
    </div>
  </div>
</div>

    <button
      class="save-btn"
      @click="saveSystemSettings"
      :disabled="saving || !authStore.isSuperAdmin"
    >
      <template v-if="saving">
        <div class="btn-spinner"></div>
        保存中...
      </template>
      <template v-else>
        ▣ 保存设置
      </template>
    </button>
  </template>

<!-- API设置 -->
<template v-if="activeTab === 'api'">
<div class="api-settings-section">
  <div v-if="apiSettingsLoading" class="loading-state">
    <div class="spinner"></div>
    <p>加载API设置...</p>
  </div>
  <template v-else>
    <!-- 全局API设置 -->
    <div class="settings-card">
      <div class="card-header">
        <h3>全局API设置</h3>
        <span class="badge">全局</span>
      </div>
      <div class="card-body">
        <div class="input-group">
          <label>API Key <span class="input-tip">全局用户共用</span></label>
          <input v-model="apiForm.globalApiKey" type="password" placeholder="留空则强制用户填写个人API" class="setting-input" />
        </div>
        <div class="input-group">
          <label>API 地址</label>
          <select v-model="apiForm.globalApiUrl" class="setting-input">
            <option value="https://generativelanguage.googleapis.com/v1beta/openai/chat/completions">Gemini AI Studio</option>
            <option value="https://api.openai.com/v1/chat/completions">OpenAI</option>
            <option value="https://openrouter.ai/api/v1/chat/completions">OpenRouter (Claude)</option>
            <option value="https://api.deepseek.com/chat/completions">DeepSeek</option>
            <option value="custom">自定义地址</option>
          </select>
        </div>
        <div v-if="apiForm.globalApiUrl === 'custom'" class="input-group">
          <label>自定义API地址</label>
          <input v-model="apiForm.globalCustomUrl" placeholder="https://api.example.com/v1/chat/completions" class="setting-input" />
        </div>
        <div class="input-group">
          <label>默认模型 <span class="input-tip">用户未选择时的默认</span></label>
          <input v-model="apiForm.globalModel" placeholder="gpt-4o-mini" class="setting-input" />
        </div>
<div class="input-group">
            <label>可用模型列表 <span class="input-tip">供用户选择，用英文逗号分隔</span></label>
            <input v-model="apiForm.globalModelList" placeholder="gpt-4o,gpt-4o-mini,claude-3-5-sonnet" class="setting-input" />
            <span class="input-desc">设置后全局用户可在这些模型中选择使用</span>
          </div>
          <div class="input-row">
            <div class="input-group">
              <label>Temperature <span class="input-tip">创造性</span></label>
              <input v-model.number="apiForm.globalTemperature" type="number" step="0.1" min="0" max="2" class="setting-input" />
            </div>
            <div class="input-group">
              <label>Max Length <span class="input-tip">最大回复长度</span></label>
              <input v-model.number="apiForm.globalMaxLength" type="number" min="100" max="32000" class="setting-input" />
            </div>
          </div>
          <div class="input-row">
            <div class="input-group">
              <label>Context Size <span class="input-tip">历史消息数</span></label>
              <input v-model.number="apiForm.globalContextSize" type="number" min="1" max="100" class="setting-input" />
            </div>
            <div class="input-group">
              <label>Timeout <span class="input-tip">超时秒数</span></label>
              <input v-model.number="apiForm.globalTimeout" type="number" min="10" max="300" class="setting-input" />
            </div>
          </div>
          <div class="input-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="apiForm.globalStreamEnabled" />
              <span>启用流式输出</span>
            </label>
          </div>
          <div class="input-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="apiForm.globalEnableSplit" />
              <span>启用长消息分段</span>
            </label>
          </div>
          <div class="input-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="apiForm.globalEnableYamlParsing" />
              <span>启用YAML解析</span>
            </label>
          </div>
        </div>
      </div>

    <!-- 社交内容全局API -->
    <div class="settings-card">
      <div class="card-header">
        <h3>社交内容全局API</h3>
        <span class="badge">全局</span>
      </div>
      <div class="card-body">
        <div class="input-group">
          <label>社交内容 API Key <span class="input-tip">用于微博、邮箱等社交内容生成</span></label>
          <input v-model="apiForm.globalSocialApiKey" type="password" placeholder="留空则使用主API" class="setting-input" />
        </div>
        <div class="input-group">
          <label>社交内容 API 地址</label>
          <input v-model="apiForm.globalSocialApiUrl" placeholder="留空则使用主API地址" class="setting-input" />
        </div>
        <div class="input-group">
          <label>社交内容默认模型</label>
          <input v-model="apiForm.globalSocialModel" placeholder="gpt-4o-mini" class="setting-input" />
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

<!-- OAuth配置 -->
<div v-show="activeTab === 'oauth'" class="settings-section">
  <div class="section-header">
    <h3>Discord OAuth设置</h3>
  </div>
  <div class="settings-list">
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Client ID</span>
        <span class="label-desc">Discord应用程序的Client ID</span>
      </div>
      <input v-model="oauthSettings.discord_client_id" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="your-discord-client-id" />
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Client Secret</span>
        <span class="label-desc">Discord应用程序的Client Secret</span>
      </div>
      <input v-model="oauthSettings.discord_client_secret" type="password" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="your-discord-client-secret" />
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Redirect URI</span>
        <span class="label-desc">OAuth回调地址，需在Discord开发者门户配置</span>
      </div>
      <input v-model="oauthSettings.discord_redirect_uri" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="https://yourdomain.com/api/auth/discord/callback" />
    </div>
  </div>
</div>
<button v-show="activeTab === 'oauth'" class="save-btn" @click="saveOauthSettings" :disabled="saving || !authStore.isSuperAdmin">
  <template v-if="saving">
    <div class="btn-spinner"></div>
    保存中...
  </template>
  <template v-else>
    ▣ 保存OAuth设置
  </template>
</button>
</div>

<!-- 联系方式 -->
  <div v-show="activeTab === 'contact'" class="settings-section">
    <div class="section-header">
      <h3>官方群聊</h3>
    </div>
    <div class="settings-list">
      <div class="contact-qrcode">
        <img :src="'/qun_qrcode.jpg'" alt="官方群二维码" class="qrcode-image" @error="qrcodeError = true" />
        <div v-if="qrcodeError" class="qrcode-placeholder">群二维码加载失败，请联系管理员</div>
        <div class="contact-password">
          <span class="password-label">进群密码：</span>
          <span class="password-value">贩子死妈</span>
        </div>
      </div>
    </div>
  </div>

<!-- 提示信息 -->
<div class="toast" v-if="toast.show" :class="toast.type">
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import apiClient from '@/api/client'
import NavBar from '@/components/common/NavBar.vue'

const authStore = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const qrcodeError = ref(false)
const faviconError = ref(false)
const uploadingQrcode = ref(false)
const uploadingFavicon = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const activeTab = ref<'credits' | 'system' | 'api' | 'oauth' | 'contact'>('credits')

const creditSettings = reactive({
  default_credits: 1000,
  signin_daily_credits: 10,
  signin_streak_bonus: 5,
  invite_reward_credits: 100,
  signin_enabled: true,
  invite_enabled: true,
})

const systemSettings = reactive({
  app_name: '',
  app_title: '',
  announcement: '',
  tips: '',
  disclaimer: '',
  qun_qrcode: '',
  favicon: '/icon.png',
})

const oauthSettings = reactive({
  discord_client_id: '',
  discord_client_secret: '',
  discord_redirect_uri: '',
})

const apiSettingsLoading = ref(false)
const apiSaving = ref(false)
const apiSettings = ref<Record<string, any> | null>(null)
const apiForm = ref({
  globalApiKey: '',
  globalApiUrl: '',
  globalCustomUrl: '',
  globalModel: '',
  globalModelList: '',
  globalTemperature: 0.9,
  globalMaxLength: 4000,
  globalContextSize: 20,
  globalTimeout: 60,
  globalStreamEnabled: true,
  globalEnableSplit: true,
  globalEnableYamlParsing: false,
  globalSocialApiKey: '',
  globalSocialApiUrl: '',
  globalSocialModel: '',
})

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

async function fetchCreditSettings() {
  loading.value = true
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    
    if (data.default_credits !== undefined) creditSettings.default_credits = Number(data.default_credits)
    if (data.signin_daily_credits !== undefined) creditSettings.signin_daily_credits = Number(data.signin_daily_credits)
    if (data.signin_streak_bonus !== undefined) creditSettings.signin_streak_bonus = Number(data.signin_streak_bonus)
    if (data.invite_reward_credits !== undefined) creditSettings.invite_reward_credits = Number(data.invite_reward_credits)
    if (data.signin_enabled !== undefined) creditSettings.signin_enabled = data.signin_enabled === true || data.signin_enabled === 'true'
    if (data.invite_enabled !== undefined) creditSettings.invite_enabled = data.invite_enabled === true || data.invite_enabled === 'true'
  } catch (err: any) {
    showToast('加载设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    loading.value = false
  }
}

async function fetchSystemSettings() {
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    systemSettings.app_name = data.app_name || ''
    systemSettings.app_title = data.app_title || ''
    systemSettings.announcement = data.announcement || ''
    systemSettings.tips = data.tips || ''
    systemSettings.disclaimer = data.disclaimer || ''
    systemSettings.qun_qrcode = data.qun_qrcode || ''
    systemSettings.favicon = data.favicon || '/icon.png'
  } catch (err: any) {
    console.error('Failed to load system settings:', err)
  }
}

async function saveCreditSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      default_credits: creditSettings.default_credits,
      signin_daily_credits: creditSettings.signin_daily_credits,
      signin_streak_bonus: creditSettings.signin_streak_bonus,
      invite_reward_credits: creditSettings.invite_reward_credits,
      signin_enabled: creditSettings.signin_enabled,
      invite_enabled: creditSettings.invite_enabled,
    })
    showToast('额度设置已保存')
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

async function handleQrcodeUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  uploadingQrcode.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res: any = await apiClient.upload('/api/settings/upload', formData)
    systemSettings.qun_qrcode = res.path
    showToast('二维码上传成功')
  } catch (err: any) {
    showToast('上传失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    uploadingQrcode.value = false
  }
}

async function handleFaviconUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  uploadingFavicon.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res: any = await apiClient.upload('/api/settings/upload', formData)
    systemSettings.favicon = res.path
    faviconError.value = false
    showToast('图标上传成功')
  } catch (err: any) {
    showToast('上传失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    uploadingFavicon.value = false
  }
}

async function saveSystemSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      app_name: systemSettings.app_name,
      app_title: systemSettings.app_title,
      announcement: systemSettings.announcement,
      tips: systemSettings.tips,
      disclaimer: systemSettings.disclaimer,
      qun_qrcode: systemSettings.qun_qrcode,
      favicon: systemSettings.favicon,
    })
    showToast('系统设置已保存')
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

async function fetchOauthSettings() {
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    oauthSettings.discord_client_id = data.discord_client_id || ''
    oauthSettings.discord_client_secret = data.discord_client_secret || ''
    oauthSettings.discord_redirect_uri = data.discord_redirect_uri || ''
  } catch (err: any) {
    console.error('Failed to load OAuth settings:', err)
  }
}

async function saveOauthSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      discord_client_id: oauthSettings.discord_client_id,
      discord_client_secret: oauthSettings.discord_client_secret,
      discord_redirect_uri: oauthSettings.discord_redirect_uri,
    })
    showToast('OAuth设置已保存')
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
          apiForm.value.globalCustomUrl = res.global.custom_url || ''
          apiForm.value.globalModel = res.global.model || ''
          apiForm.value.globalModelList = res.global.model_list || ''
          apiForm.value.globalTemperature = res.global.temperature ?? 0.9
          apiForm.value.globalMaxLength = res.global.max_length ?? 4000
          apiForm.value.globalContextSize = res.global.context_size ?? 20
          apiForm.value.globalTimeout = res.global.timeout ?? 60
          apiForm.value.globalStreamEnabled = res.global.stream_enabled ?? true
          apiForm.value.globalEnableSplit = res.global.enable_split ?? true
          apiForm.value.globalEnableYamlParsing = res.global.enable_yaml_parsing ?? false
          apiForm.value.globalSocialApiKey = res.global.social_api_key || ''
          apiForm.value.globalSocialApiUrl = res.global.social_api_url || ''
          apiForm.value.globalSocialModel = res.global.social_model || ''
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
          custom_url: apiForm.value.globalCustomUrl,
          model: apiForm.value.globalModel,
          model_list: apiForm.value.globalModelList,
          temperature: apiForm.value.globalTemperature,
          max_length: apiForm.value.globalMaxLength,
          context_size: apiForm.value.globalContextSize,
          timeout: apiForm.value.globalTimeout,
          stream_enabled: apiForm.value.globalStreamEnabled,
          enable_split: apiForm.value.globalEnableSplit,
          enable_yaml_parsing: apiForm.value.globalEnableYamlParsing,
          social_api_key: apiForm.value.globalSocialApiKey,
          social_api_url: apiForm.value.globalSocialApiUrl,
          social_model: apiForm.value.globalSocialModel,
          is_global: true,
        })
    }
    showToast('API设置已保存')
  } catch (err: any) {
    showToast('保存API设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    apiSaving.value = false
  }
}

onMounted(() => {
  if (authStore.isAdmin) {
    fetchCreditSettings()
    fetchSystemSettings()
    fetchApiSettings()
    fetchOauthSettings()
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

/* Section */
.settings-section {
  margin-bottom: 20px;
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

/* 设置项 */
.settings-list {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label-text {
  font-size: 15px;
  color: #fff;
}

.label-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.setting-input-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.setting-input {
  width: 80px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  text-align: right;
  outline: none;
}

.setting-input:focus {
  border-color: #5B6EF5;
}

.setting-input.full-width {
  width: 100%;
  text-align: left;
}

.setting-item.vertical {
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
}

.setting-item.vertical .setting-label {
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.setting-item.vertical .setting-input {
  width: 100%;
  box-sizing: border-box;
}

.setting-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-suffix {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
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
  background: var(--bg-quaternary, rgba(255, 255, 255, 0.1));
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
  background: #34C759;
}

.toggle input:checked + .toggle-slider::after {
  transform: translateX(20px);
}

/* 操作按钮 */
.save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  width: 100%;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  border: none;
  border-radius: 12px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 20px;
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
  from { opacity: 0; transform: translateX(-50%) translateY(10px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
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
  flex-wrap: wrap;
}

.tab-btn {
  padding: 8px 16px;
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
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.input-tip {
  font-size: 11px;
  color: rgba(91, 110, 245, 0.8);
  margin-left: 6px;
}

.input-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 4px;
  display: block;
}

.input-row {
  display: flex;
  gap: 12px;
}

.input-row .input-group {
  flex: 1;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #fff;
  font-size: 14px;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #5B6EF5;
}

.qrcode-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  width: 100%;
}

.contact-qrcode {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  gap: 12px;
}

.contact-password {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
}

.password-label {
  color: rgba(255, 255, 255, 0.6);
}

.password-value {
  color: #5B6EF5;
  font-weight: 600;
}

.qrcode-image {
  max-width: 200px;
  max-height: 200px;
  border-radius: 12px;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.qrcode-placeholder {
  color: rgba(255, 255, 255, 0.4);
  font-size: 14px;
  padding: 20px;
}

.password-display {
  font-size: 18px;
  font-weight: 600;
  color: #5B6EF5;
  padding: 12px 20px;
  background: rgba(91, 110, 245, 0.1);
  border-radius: 8px;
  letter-spacing: 2px;
}

.file-upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 8px 0;
}

.upload-btn {
  display: inline-block;
  padding: 8px 16px;
  background: rgba(91, 110, 245, 0.2);
  border: 1px solid rgba(91, 110, 245, 0.4);
  border-radius: 6px;
  color: #5B6EF5;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-btn:hover {
  background: rgba(91, 110, 245, 0.3);
  border-color: #5B6EF5;
}

.upload-status {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.favicon-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.favicon-img {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  object-fit: contain;
  background: rgba(255, 255, 255, 0.1);
}

.favicon-placeholder {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
}
</style>