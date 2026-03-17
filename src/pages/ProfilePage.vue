<template>
  <div class="profile-page">
    <NavBar title="个人中心" :show-back="true" back-to="/">
      <template #right>
        <button class="nav-btn" @click="router.push('/customize')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <circle cx="12" cy="12" r="3" />
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
          </svg>
        </button>
      </template>
    </NavBar>

    <div class="profile-content">
<!-- 用户信息卡片 -->
  <div class="user-card">
    <div class="avatar-section" @click="triggerAvatarUpload">
      <div class="avatar">
        <img v-if="userAvatar" :src="userAvatar" alt="avatar" />
        <span v-else class="avatar-placeholder">{{ user?.username?.[0] || '?' }}</span>
      </div>
      <div class="avatar-edit-badge">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/>
          <circle cx="12" cy="13" r="4"/>
        </svg>
      </div>
      <div class="online-dot"></div>
    </div>
    <div class="user-info" @click="showEditProfile = true">
      <h2 class="display-name">{{ user?.displayName || user?.username || '未登录' }}</h2>
      <p v-if="showAccountName" class="account-name">账户名：{{ user?.username }}</p>
      <p class="user-id">ID: {{ user?.id || '—' }}</p>
    </div>
    <input type="file" ref="avatarInput" accept="image/*" style="display: none" @change="handleAvatarChange" />
  </div>

      <!-- 数据统计 -->
      <div class="stats-row">
        <div class="stat-item" v-for="stat in stats" :key="stat.label">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
      </div>

      <!-- 功能菜单 -->
      <div class="menu-section">
        <div class="menu-group">
          <div class="menu-item" v-for="item in menuItems" :key="item.label" @click="handleMenuItem(item)">
            <div class="menu-icon" :style="{ background: item.color }">
              <span>{{ item.icon }}</span>
            </div>
            <span class="menu-label">{{ item.label }}</span>
            <svg class="menu-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18l6-6-6-6" />
            </svg>
          </div>
        </div>
      </div>

      <!-- 管理员面板入口 -->
      <div class="menu-section" v-if="authStore.isAdmin">
        <div class="section-title">管理</div>
        <div class="menu-group">
          <div class="menu-item" v-for="item in adminItems" :key="item.label" @click="handleMenuItem(item)">
            <div class="menu-icon" :style="{ background: item.color }">
              <span>{{ item.icon }}</span>
            </div>
            <span class="menu-label">{{ item.label }}</span>
            <svg class="menu-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18l6-6-6-6" />
            </svg>
          </div>
        </div>
      </div>

      <!-- 快捷设置 -->
      <div class="menu-section">
        <div class="section-title">快捷设置</div>
        <div class="menu-group">
          <div class="setting-item">
            <div class="setting-left">
              <span class="setting-icon">☽</span>
              <span>深色模式</span>
            </div>
            <label class="toggle">
              <input type="checkbox" :checked="settingsStore.settings.darkMode" @change="toggleDarkMode" />
              <span class="toggle-slider"></span>
            </label>
          </div>
          <div class="setting-item">
            <div class="setting-left">
              <span class="setting-icon">▲</span>
              <span>通知提醒</span>
            </div>
            <label class="toggle">
              <input type="checkbox" :checked="settingsStore.settings.notifyEnabled" @change="settingsStore.settings.notifyEnabled = !settingsStore.settings.notifyEnabled" />
              <span class="toggle-slider"></span>
            </label>
          </div>
<div class="setting-item">
  <div class="setting-left">
    <span class="setting-icon">◉</span>
    <span>消息音效</span>
  </div>
  <label class="toggle">
    <input type="checkbox" :checked="settingsStore.settings.notificationSoundEnabled" @change="settingsStore.settings.notificationSoundEnabled = !settingsStore.settings.notificationSoundEnabled" />
    <span class="toggle-slider"></span>
  </label>
</div>
<div class="setting-item">
  <div class="setting-left">
    <span class="setting-icon">🕐</span>
    <span>时间格式</span>
  </div>
  <select class="time-format-select" v-model="settingsStore.settings.timeFormat" @change="settingsStore.settings.timeFormat = settingsStore.settings.timeFormat">
    <option value="24h">24小时制</option>
    <option value="12h">12小时制</option>
  </select>
</div>
<div class="setting-item" @click="router.push('/customize')">
            <div class="setting-left">
              <span class="setting-icon">⚙️</span>
              <span>全部设置</span>
            </div>
            <svg class="menu-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 18l6-6-6-6" />
            </svg>
          </div>
        </div>
      </div>

<!-- 退出登录 -->
  <button class="logout-btn" @click="handleLogout">退出登录</button>

  <div class="version-info">
    <p class="copyright">贩子死妈小手机版权所有</p>
    <p class="version">v1.0.0</p>
    <p class="custom-disclaimer" v-if="customDisclaimer">{{ customDisclaimer }}</p>
  </div>
</div>

  <!-- 编辑资料弹窗 -->
  <div v-if="showEditProfile" class="modal-overlay" @click="showEditProfile = false">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>编辑资料</h3>
        <button class="modal-close" @click="showEditProfile = false">&times;</button>
      </div>
<div class="modal-body">
  <div class="form-group">
    <label>显示名称</label>
    <input type="text" v-model="editDisplayName" placeholder="输入显示名称" maxlength="50" />
  </div>
  <div class="form-group">
    <label>头像链接</label>
    <input type="text" v-model="editAvatarUrl" placeholder="输入图片URL，如 https://xxx.com/avatar.jpg" maxlength="500" />
  </div>
</div>
      <div class="modal-footer">
        <button class="btn-cancel" @click="showEditProfile = false">取消</button>
        <button class="btn-confirm" @click="saveProfile" :disabled="saving">保存</button>
      </div>
    </div>
  </div>

  <!-- 联系方式弹窗 -->
  <div v-if="showContactModal" class="modal-overlay" @click="showContactModal = false">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>官方群聊</h3>
        <button class="modal-close" @click="showContactModal = false">&times;</button>
      </div>
      <div class="modal-body contact-body">
        <img :src="qrcodeUrl" alt="官方群二维码" class="contact-qrcode-img" @error="qrcodeLoadError = true" />
        <div v-if="qrcodeLoadError" class="contact-qrcode-placeholder">群二维码加载失败</div>
        <div class="contact-info" v-if="!qrcodeLoadError">
          <span class="contact-label">进群密码：</span>
          <span class="contact-value">贩子死妈</span>
        </div>
      </div>
    </div>
  </div>

  <!-- 我的邀请弹窗 -->
  <div v-if="showInviteModal" class="modal-overlay" @click="showInviteModal = false">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>我的邀请</h3>
        <button class="modal-close" @click="showInviteModal = false">&times;</button>
      </div>
      <div class="modal-body invite-body" v-if="!loadingInvite">
        <div class="invite-code-section">
          <span class="invite-code-label">我的邀请码</span>
          <div class="invite-code-value">{{ inviteData.code || '加载中...' }}</div>
        </div>
        <div class="invite-stats">
          <div class="invite-stat">
            <span class="stat-num">{{ inviteData.invitees?.length || 0 }}</span>
            <span class="stat-text">已邀请</span>
          </div>
          <div class="invite-stat">
            <span class="stat-num">{{ inviteData.totalRewards || 0 }}</span>
            <span class="stat-text">累计获得额度</span>
          </div>
        </div>
        <div class="invite-list" v-if="inviteData.invitees?.length > 0">
          <div class="invite-list-title">邀请记录</div>
          <div class="invite-item" v-for="invitee in inviteData.invitees" :key="invitee.id">
            <span class="invitee-name">{{ invitee.username }}</span>
            <span class="invitee-reward">+{{ invitee.reward }}</span>
          </div>
        </div>
        <div class="invite-tip">
          好友每次签到，您将获得其获得额度的15%作为返利
        </div>
      </div>
      <div class="modal-body invite-body loading" v-else>
        <div class="loading-spinner"></div>
      </div>
    </div>
  </div>
</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePhoneStore } from '@/stores/phone'
import { useSettingsStore } from '@/stores/settings'
import NavBar from '@/components/common/NavBar.vue'

const router = useRouter()
const authStore = useAuthStore()
const phoneStore = usePhoneStore()
const settingsStore = useSettingsStore()

const API_BASE = import.meta.env.VITE_API_URL || ''

const bio = ref('')
const showEditProfile = ref(false)
const showContactModal = ref(false)
const showInviteModal = ref(false)
const qrcodeLoadError = ref(false)
const editDisplayName = ref('')
const editAvatarUrl = ref('')
const saving = ref(false)
const avatarInput = ref<HTMLInputElement | null>(null)
const inviteData = ref<{ code: string; invitees: any[]; totalRewards: number }>({
  code: '',
  invitees: [],
  totalRewards: 0
})
const loadingInvite = ref(false)
const customDisclaimer = ref('')

const user = computed(() => authStore.user)
const userAvatar = computed(() => {
  const avatar = user.value?.avatar
  if (!avatar) return ''
  if (avatar.startsWith('http')) return avatar
  if (avatar.startsWith('/avatar_')) return avatar
  return ''
})
const qrcodeUrl = computed(() => {
  // 使用GitHub raw URL作为默认
  return 'https://raw.githubusercontent.com/Xeltra233/fzsmphone/master/public/qun_qrcode.jpg'
})
const showAccountName = computed(() => {
  return user.value?.displayName && user.value?.displayName !== user.value?.username
})

interface StatItem {
  label: string
  value: string | number
}

const stats = ref<StatItem[]>([
  { label: '对话', value: 0 },
  { label: '帖子', value: 0 },
  { label: '日记', value: 0 },
  { label: '天数', value: 0 },
])

interface MenuItem {
  icon: string
  label: string
  color: string
  route?: string
  action?: string
}

const menuItems: MenuItem[] = [
  { icon: '○', label: '编辑资料', color: 'linear-gradient(135deg, #5B6EF5, #8B5CF6)', action: 'editProfile' },
  { icon: '◈', label: '我的人设', color: 'linear-gradient(135deg, #E17055, #FAB1A0)', route: '/personas' },
  { icon: '○', label: '角色管理', color: 'linear-gradient(135deg, #00CEC9, #81ECEC)', route: '/characters' },
  { icon: '▥', label: '世界书', color: 'linear-gradient(135deg, #6C5CE7, #DDA0DD)', route: '/worldbook' },
  { icon: '⚡', label: '预设管理', color: 'linear-gradient(135deg, #FDCB6E, #F39C12)', route: '/preset' },
  { icon: '✧', label: '主题定制', color: 'linear-gradient(135deg, #636E72, #B2BEC3)', route: '/customize' },
  { icon: '¤', label: '我的钱包', color: 'linear-gradient(135deg, #00B894, #55EFC4)', route: '/wallet' },
  { icon: '◈', label: '我的额度', color: 'linear-gradient(135deg, #667eea, #764ba2)', route: '/credits' },
  { icon: '💬', label: '官方群聊', color: 'linear-gradient(135deg, #74B9FF, #0984E3)', action: 'showContact' },
  { icon: '∞', label: '我的邀请', color: 'linear-gradient(135deg, #FD79A8, #E84393)', action: 'showInvite' },
]

const adminItems: MenuItem[] = [
  { icon: '■', label: '功能管理', color: 'linear-gradient(135deg, #FF6B6B, #EE5A24)', route: '/admin/features' },
  { icon: '○', label: '用户管理', color: 'linear-gradient(135deg, #E6162D, #FF4757)', route: '/admin/users' },
  { icon: '▥', label: '数据统计', color: 'linear-gradient(135deg, #3498DB, #2980B9)', route: '/admin/stats' },
  { icon: '⚙', label: '系统设置', color: 'linear-gradient(135deg, #636E72, #95A5A6)', route: '/admin/settings' },
]

function handleMenuItem(item: MenuItem) {
  if (item.route) {
    router.push(item.route)
  } else if (item.action === 'editProfile') {
    editDisplayName.value = user.value?.displayName || user.value?.username || ''
    showEditProfile.value = true
  } else if (item.action === 'showContact') {
    qrcodeLoadError.value = false
    showContactModal.value = true
  } else if (item.action === 'showInvite') {
    loadInviteData()
    showInviteModal.value = true
  }
}

async function loadInviteData() {
  if (!user.value) return
  loadingInvite.value = true
  try {
    const response = await fetch(`${API_BASE}/api/credits/invite-info`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (response.ok) {
      inviteData.value = await response.json()
    }
  } catch (err) {
    console.error('Failed to load invite data:', err)
  } finally {
    loadingInvite.value = false
  }
}

function toggleDarkMode() {
  settingsStore.settings.darkMode = !settingsStore.settings.darkMode
  phoneStore.toggleTheme()
}

async function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function triggerAvatarUpload() {
  avatarInput.value?.click()
}

async function handleAvatarChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file || !user.value) return

  const formData = new FormData()
  formData.append('avatar', file)

  try {
    const response = await fetch(`/api/users/${user.value.id}/avatar`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: formData
    })
    if (response.ok) {
      const data = await response.json()
      authStore.user!.avatar = data.avatar_url
      localStorage.setItem('user', JSON.stringify(authStore.user))
    }
  } catch (err) {
    console.error('Failed to upload avatar:', err)
  }
  target.value = ''
}

function openEditProfile() {
  editDisplayName.value = user.value?.displayName || user.value?.username || ''
  editAvatarUrl.value = user.value?.avatar || ''
  showEditProfile.value = true
}

async function saveProfile() {
  if (!user.value) return
  saving.value = true
  try {
    const response = await fetch(`/api/users/${user.value.id}/profile`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        display_name: editDisplayName.value,
        avatar_url: editAvatarUrl.value
      })
    })
    if (response.ok) {
      authStore.user!.displayName = editDisplayName.value
      if (editAvatarUrl.value) {
        authStore.user!.avatar = editAvatarUrl.value
      }
      localStorage.setItem('user', JSON.stringify(authStore.user))
      showEditProfile.value = false
      editAvatarUrl.value = ''
    }
  } catch (err) {
    console.error('Failed to save profile:', err)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  // Sync dark mode state
  if (settingsStore.settings.darkMode !== (phoneStore.theme === 'dark')) {
    document.documentElement.setAttribute('data-theme', settingsStore.settings.darkMode ? 'dark' : 'light')
  }

  // Calculate days since first use
  const firstUse = localStorage.getItem('first_use')
  if (!firstUse) {
    localStorage.setItem('first_use', new Date().toISOString())
    stats.value[3].value = 1
  } else {
    const days = Math.ceil((Date.now() - new Date(firstUse).getTime()) / 86400000)
    stats.value[3].value = days
  }

  // Load custom disclaimer from settings
  try {
    const response = await fetch(`${API_BASE}/api/settings/key?key=disclaimer`, {
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    })
    if (response.ok) {
      const data = await response.json()
      if (data.value) {
        customDisclaimer.value = data.value
      }
    }
  } catch (err) {
    console.error('Failed to load disclaimer:', err)
  }
})
</script>

<style scoped>
.profile-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

.profile-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.profile-content::-webkit-scrollbar {
  display: none;
}

/* 用户卡片 */
.user-card {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.avatar-section {
  position: relative;
  flex-shrink: 0;
  cursor: pointer;
}

.avatar-edit-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 24px;
  height: 24px;
  background: var(--color-primary, #5B6EF5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--bg-primary);
}

.avatar-edit-badge svg {
  width: 12px;
  height: 12px;
  color: #fff;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
}

.online-dot {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 14px;
  height: 14px;
  background: var(--color-green);
  border-radius: 50%;
  border: 3px solid var(--bg-primary);
}

.user-info {
  flex: 1;
  min-width: 0;
}

.display-name {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.user-id {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 2px 0 6px;
}

.account-name {
  font-size: 12px;
  color: var(--text-quaternary);
  margin: 0;
}

.user-bio {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 统计行 */
.stats-row {
  display: flex;
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  padding: 16px;
  margin-bottom: 16px;
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 菜单 */
.menu-section {
  margin-bottom: 16px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 4px 8px;
}

.menu-group {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  cursor: pointer;
  transition: background 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.menu-item:active {
  background: var(--bg-tertiary);
}

.menu-item:not(:last-child) {
  border-bottom: 0.5px solid var(--separator);
}

.menu-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.menu-label {
  flex: 1;
  font-size: 16px;
  color: var(--text-primary);
}

.menu-arrow {
  width: 16px;
  height: 16px;
  color: var(--text-quaternary);
  flex-shrink: 0;
}

/* 退出按钮 */
.logout-btn {
  width: 100%;
  padding: 14px;
  background: var(--bg-primary);
  border: none;
  border-radius: var(--radius-lg);
  color: var(--color-red);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 16px;
  transition: background 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.logout-btn:active {
  background: var(--bg-tertiary);
}

.version-info {
  text-align: center;
  font-size: 12px;
  color: var(--text-quaternary);
  padding-bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.version-info .copyright {
  font-weight: 600;
  color: var(--text-tertiary);
  margin: 0;
}

.version-info .version {
  margin: 0;
}

.version-info .custom-disclaimer {
  margin: 4px 0 0;
  font-size: 11px;
  color: var(--text-quaternary);
}

/* 快捷设置项 */
.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
}

.setting-item:not(:last-child) {
  border-bottom: 0.5px solid var(--separator);
}

.setting-item:active {
  background: var(--bg-tertiary);
}

.setting-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.setting-icon {
  font-size: 18px;
}

.setting-item span {
  font-size: 16px;
  color: var(--text-primary);
}

/* Toggle switch */
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

.time-format-select {
  padding: 6px 12px;
  border: 0.5px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
}

.time-format-select:focus {
  outline: none;
  border-color: var(--color-primary, #5B6EF5);
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  width: 90%;
  max-width: 400px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 0.5px solid var(--separator);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 0.5px solid var(--separator);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 16px;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary, #5B6EF5);
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 0.5px solid var(--separator);
}

.btn-cancel, .btn-confirm {
  flex: 1;
  padding: 12px;
  border-radius: var(--radius-md);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  border: none;
}

.btn-cancel {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.btn-confirm {
  background: var(--color-primary, #5B6EF5);
  color: #fff;
}

.btn-confirm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.user-info {
  cursor: pointer;
}

.contact-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 24px;
}

.contact-qrcode-img {
  width: 200px;
  height: 200px;
  object-fit: contain;
  border-radius: 12px;
  border: 1px solid var(--separator);
}

.contact-qrcode-placeholder {
  color: var(--text-tertiary);
  font-size: 14px;
  padding: 40px;
}

.contact-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 16px;
}

.contact-label {
  color: var(--text-secondary);
}

.contact-value {
  color: var(--color-primary, #5B6EF5);
  font-weight: 600;
}

.invite-body {
  padding: 20px;
}

.invite-body.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.invite-code-section {
  text-align: center;
  margin-bottom: 20px;
}

.invite-code-label {
  font-size: 14px;
  color: var(--text-tertiary);
  display: block;
  margin-bottom: 8px;
}

.invite-code-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-primary, #5B6EF5);
  letter-spacing: 2px;
}

.invite-stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
}

.invite-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-num {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-text {
  font-size: 12px;
  color: var(--text-tertiary);
}

.invite-list {
  margin-bottom: 16px;
}

.invite-list-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.invite-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  margin-bottom: 8px;
}

.invitee-name {
  color: var(--text-primary);
}

.invitee-reward {
  color: var(--color-green);
  font-weight: 600;
}

.invite-tip {
  text-align: center;
  font-size: 12px;
  color: var(--text-quaternary);
  padding: 12px;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid var(--separator);
  border-top-color: var(--color-primary, #5B6EF5);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
