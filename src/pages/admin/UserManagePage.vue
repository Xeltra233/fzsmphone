<template>
  <div class="user-manage-page">
    <nav-bar title="用户管理" back-route="/profile" />

    <div class="page-content">
      <div v-if="!authStore.isAdmin" class="no-access">
        <p>⛔ 需要管理员权限</p>
      </div>

      <template v-else>
        <!-- 搜索栏 -->
        <div class="search-bar">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索用户名、ID或Discord ID..."
            class="search-input"
          />
          <button class="create-user-btn" @click="showCreateUserModal = true">添加用户</button>
        </div>

        <!-- 统计概览 -->
        <div class="stats-bar">
          <div class="stat-chip">
            <span class="stat-chip-label">总用户</span>
            <span class="stat-chip-value">{{ users.length }}</span>
          </div>
          <div class="stat-chip">
            <span class="stat-chip-label">管理员</span>
            <span class="stat-chip-value">{{ roleCounts.admin }}</span>
          </div>
          <div class="stat-chip">
            <span class="stat-chip-label">普通用户</span>
            <span class="stat-chip-value">{{ roleCounts.user }}</span>
          </div>
          <div class="stat-chip">
            <span class="stat-chip-label">已封禁</span>
            <span class="stat-chip-value stat-banned">{{ bannedCount }}</span>
          </div>
        </div>

        <div v-if="loading" class="loading">加载中...</div>

        <div v-else-if="filteredUsers.length === 0" class="empty-state">
          <p>{{ searchQuery ? '未找到匹配用户' : '暂无用户数据' }}</p>
        </div>

        <div v-else class="user-list">
          <div
            v-for="u in filteredUsers"
            :key="u.id"
            class="user-item"
            :class="{ 'user-banned': u.is_banned }"
          >
            <div class="user-avatar">
              <img v-if="u.avatar_url" :src="u.avatar_url" alt="" />
              <span v-else class="avatar-placeholder">{{ (u.display_name || u.username)?.[0] || '?' }}</span>
              <div v-if="u.is_banned" class="avatar-ban-badge">
                <svg viewBox="0 0 16 16" fill="none" width="12" height="12"><circle cx="8" cy="8" r="6.5" stroke="#ff6b6b" stroke-width="1.5"/><line x1="4" y1="4" x2="12" y2="12" stroke="#ff6b6b" stroke-width="1.5"/></svg>
              </div>
            </div>
<div class="user-info">
  <div class="user-name">
    {{ u.display_name || u.username }}
    <span v-if="u.is_super_admin" class="role-badge badge-super">超级管理员</span>
    <span v-if="u.is_banned" class="ban-tag">已封禁</span>
  </div>
  <div class="user-meta">
    <span class="user-username">@{{ u.username }}</span>
    <span class="user-credits" @click="openCreditsModal(u)" title="点击修改额度">
      额度: {{ u.credits || 0 }} | Token: {{ u.total_tokens || 0 }}
    </span>
    <span class="user-date">{{ formatDate(u.created_at) }}</span>
  </div>
  <div v-if="u.is_banned && u.ban_reason" class="user-ban-reason">
    原因：{{ u.ban_reason }}
    <span v-if="u.banned_until" class="banned-until"> 解除时间: {{ formatDate(u.banned_until) }}</span>
  </div>
</div>
<div class="user-actions">
  <select
    :value="u.role"
    @change="changeRole(u, ($event.target as HTMLSelectElement).value)"
    class="role-select"
    :class="'role-' + u.role"
    :disabled="u.is_banned || (u.is_super_admin && !authStore.isSuperAdmin) || (u.role === 'admin' && !authStore.isSuperAdmin) || (String(u.id) === String(authStore.user?.id))"
  >
    <option value="user" :disabled="u.role === 'super_admin'">普通用户</option>
    <option value="admin" :disabled="!authStore.isSuperAdmin || u.role === 'super_admin'">管理员</option>
  </select>
              <button
                v-if="authStore.isSuperAdmin && String(u.id) !== String(authStore.user?.id)"
                class="impersonate-btn"
                @click="impersonateUser(u)"
                title="代管登录"
              >
                代管
              </button>
              <button
                v-if="!u.is_banned && u.role !== 'admin'"
                class="ban-btn"
                @click="openBanModal(u)"
                title="封禁用户"
              >
                <svg viewBox="0 0 16 16" fill="none" width="14" height="14"><rect x="3" y="7" width="10" height="7" rx="1.5" stroke="currentColor" stroke-width="1.3"/><path d="M5 7V5a3 3 0 016 0v2" stroke="currentColor" stroke-width="1.3" fill="none"/></svg>
              </button>
              <button
                v-if="u.is_banned"
                class="unban-btn"
                @click="unbanUser(u)"
                title="解除封禁"
              >
                <svg viewBox="0 0 16 16" fill="none" width="14" height="14"><rect x="3" y="7" width="10" height="7" rx="1.5" stroke="currentColor" stroke-width="1.3"/><path d="M5 7V5a3 3 0 016 0" stroke="currentColor" stroke-width="1.3" fill="none"/></svg>
              </button>
            </div>
          </div>
        </div>

        <div v-if="saving" class="save-indicator">保存中...</div>
      </template>
    </div>

    <!-- 封禁弹窗 -->
    <teleport to="body">
      <div v-if="showBanModal" class="modal-overlay" @click.self="showBanModal = false">
        <div class="ban-modal">
          <div class="ban-modal-header">
            <span class="ban-modal-icon">
              <svg viewBox="0 0 24 24" fill="none" width="24" height="24"><circle cx="12" cy="12" r="9" stroke="#ff6b6b" stroke-width="2"/><line x1="7" y1="7" x2="17" y2="17" stroke="#ff6b6b" stroke-width="2"/></svg>
            </span>
            <h3>封禁用户</h3>
          </div>
          <div class="ban-modal-user">
            <div class="ban-modal-avatar">
              <img v-if="banTarget?.avatar_url" :src="banTarget?.avatar_url" alt="" />
              <span v-else class="avatar-placeholder">{{ (banTarget?.display_name || banTarget?.username)?.[0] || '?' }}</span>
            </div>
            <div>
              <div class="ban-modal-name">{{ banTarget?.display_name || banTarget?.username }}</div>
              <div class="ban-modal-username">@{{ banTarget?.username }}</div>
            </div>
          </div>
          <div class="ban-modal-field">
            <label>封禁原因（可选）</label>
            <textarea
              v-model="banReason"
              placeholder="请输入封禁原因..."
              rows="3"
              class="ban-reason-input"
            ></textarea>
          </div>
          <div class="ban-modal-field">
            <label>封禁时长</label>
            <select v-model="banDuration" class="ban-duration-select">
              <option value="0">永久封禁</option>
              <option value="1">1小时</option>
              <option value="3">3小时</option>
              <option value="24">1天</option>
              <option value="72">3天</option>
              <option value="168">7天</option>
              <option value="720">30天</option>
            </select>
          </div>
          <div class="ban-modal-warning">
            注意：封禁后该用户将无法登录和使用任何功能
          </div>
          <div class="ban-modal-actions">
            <button class="btn-cancel" @click="showBanModal = false">取消</button>
            <button class="btn-confirm-ban" @click="confirmBan">确认封禁</button>
          </div>
        </div>
      </div>
</teleport>

<!-- 额度编辑弹窗 -->
<teleport to="body">
  <div v-if="showCreditsModal" class="modal-overlay" @click.self="showCreditsModal = false">
    <div class="credits-modal">
      <div class="credits-modal-header">
        <h3>调整用户额度</h3>
        <button class="modal-close" @click="showCreditsModal = false">&times;</button>
      </div>
      <div class="credits-modal-user">
        <div class="credits-modal-avatar">
          <img v-if="creditsTarget?.avatar_url" :src="creditsTarget?.avatar_url" alt="" />
          <span v-else class="avatar-placeholder">{{ (creditsTarget?.display_name || creditsTarget?.username)?.[0] || '?' }}</span>
        </div>
        <div>
          <div class="credits-modal-name">{{ creditsTarget?.display_name || creditsTarget?.username }}</div>
          <div class="credits-modal-username">@{{ creditsTarget?.username }}</div>
        </div>
      </div>
      <div class="credits-modal-current">
        当前额度: <span class="credits-value">{{ creditsTarget?.credits || 0 }}</span>
      </div>
      <div class="credits-modal-field">
        <label>调整额度</label>
        <div class="credits-input-group">
          <button class="credits-btn minus" @click="creditsAdjust -= 100">-100</button>
          <button class="credits-btn minus" @click="creditsAdjust -= 10">-10</button>
          <input v-model.number="creditsAdjust" type="number" class="credits-input" placeholder="正数增加，负数减少" />
          <button class="credits-btn plus" @click="creditsAdjust += 10">+10</button>
          <button class="credits-btn plus" @click="creditsAdjust += 100">+100</button>
        </div>
      </div>
      <div class="credits-modal-field">
        <label>或者直接设置额度</label>
        <input v-model.number="creditsSet" type="number" class="credits-input" placeholder="直接设置额度" />
      </div>
      <div class="credits-modal-field">
        <label>角色卡总空间（MB）</label>
        <input v-model.number="characterQuotaMb" type="number" min="1" class="credits-input" placeholder="默认 10" />
      </div>
      <div class="credits-modal-actions">
        <button class="btn-cancel" @click="showCreditsModal = false">取消</button>
        <button class="btn-confirm-credits" @click="confirmCreditsAdjust">确认调整</button>
      </div>
    </div>
  </div>
</teleport>

<teleport to="body">
  <div v-if="showCreateUserModal" class="modal-overlay" @click.self="showCreateUserModal = false">
    <div class="credits-modal create-user-modal">
      <div class="credits-modal-header">
        <h3>添加用户</h3>
        <button class="modal-close" @click="showCreateUserModal = false">&times;</button>
      </div>
      <div class="credits-modal-field">
        <label>用户名</label>
        <input v-model.trim="createUserForm.username" type="text" class="credits-input" placeholder="输入用户名" />
      </div>
      <div class="credits-modal-field">
        <label>邮箱（可选）</label>
        <input v-model.trim="createUserForm.email" type="email" class="credits-input" placeholder="输入邮箱" />
      </div>
      <div class="credits-modal-field">
        <label>密码</label>
        <input v-model="createUserForm.password" type="text" class="credits-input" placeholder="至少 6 位" />
      </div>
      <div class="credits-modal-field">
        <label>角色</label>
        <select v-model="createUserForm.role" class="ban-duration-select">
          <option value="user">普通用户</option>
          <option v-if="authStore.isSuperAdmin" value="admin">管理员</option>
        </select>
      </div>
      <div class="credits-modal-field">
        <label>初始额度</label>
        <input v-model.number="createUserForm.credits" type="number" class="credits-input" placeholder="默认 0" />
      </div>
      <div class="credits-modal-actions">
        <button class="btn-cancel" @click="showCreateUserModal = false">取消</button>
        <button class="btn-confirm-credits" @click="createUser">确认创建</button>
      </div>
    </div>
  </div>
</teleport>
</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/client'
import router from '@/router'

const authStore = useAuthStore()

interface UserRecord {
  id: number
  discord_id: string
  username: string
  display_name: string
  avatar_url: string
  role: string
  is_super_admin: boolean
  is_banned: boolean
  ban_reason: string
  banned_at: string | null
  banned_until: string | null
  created_at: string
  updated_at: string
  credits: number
  total_tokens: number
  character_storage_quota: number
}

const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const users = ref<UserRecord[]>([])

// Ban modal state
const showBanModal = ref(false)
const banTarget = ref<UserRecord | null>(null)
const banReason = ref('')
const banDuration = ref('0')

// Credits modal state
const showCreditsModal = ref(false)
const creditsTarget = ref<UserRecord | null>(null)
const creditsAdjust = ref(0)
const creditsSet = ref(0)
const characterQuotaMb = ref(10)
const showCreateUserModal = ref(false)
const createUserForm = ref({
  username: '',
  email: '',
  password: '',
  role: 'user',
  credits: 0,
})

const roleCounts = computed(() => {
  const counts = { admin: 0, user: 0 }
  for (const u of users.value) {
    if (u.role in counts) {
      counts[u.role as keyof typeof counts]++
    }
  }
  return counts
})

const bannedCount = computed(() => users.value.filter(u => u.is_banned).length)

const filteredUsers = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return users.value
  return users.value.filter(
    (u) =>
      u.username.toLowerCase().includes(q) ||
      u.display_name?.toLowerCase().includes(q) ||
      String(u.id).includes(q) ||
      u.discord_id?.toLowerCase().includes(q)
  )
})

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

async function fetchUsers() {
  loading.value = true
  try {
    const res: any = await api.get('/api/users')
    users.value = res.data || res || []
  } catch (e: any) {
    console.error('Failed to fetch users:', e)
    alert('加载用户列表失败: ' + (e.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

async function changeRole(user: UserRecord, newRole: string) {
  if (newRole === user.role) return
  const oldRole = user.role

  const currentUser = authStore.user
  const currentIsSuperAdmin = currentUser?.isSuperAdmin || authStore.isSuperAdmin

  if (user.is_super_admin && !currentIsSuperAdmin) {
    alert('无法修改超级管理员的权限')
    return
  }

  if (user.role === 'super_admin' && !currentIsSuperAdmin) {
    alert('无法修改超级管理员的权限')
    return
  }

  if (user.role === 'admin' && !currentIsSuperAdmin) {
    alert('只有超级管理员可以修改管理员')
    return
  }

  if (String(user.id) === String(currentUser?.id)) {
    alert('无法修改自己的角色')
    return
  }

  saving.value = true
  try {
    await api.patch(`/api/users/${user.id}/role`, { role: newRole })
    user.role = newRole
  } catch (e: any) {
    console.error('Failed to update role:', e)
    alert('修改角色失败: ' + (e.message || '未知错误'))
    user.role = oldRole
  } finally {
    saving.value = false
  }
}

function openBanModal(user: UserRecord) {
  banTarget.value = user
  banReason.value = ''
  showBanModal.value = true
}

async function confirmBan() {
  if (!banTarget.value) return
  saving.value = true
  try {
    const durationHours = parseInt(banDuration.value)
    await api.post(`/api/users/${banTarget.value.id}/ban`, { 
      reason: banReason.value,
      duration_hours: durationHours
    })
    banTarget.value.is_banned = true
    banTarget.value.ban_reason = banReason.value
    banTarget.value.banned_at = new Date().toISOString()
    if (durationHours > 0) {
      const until = new Date()
      until.setHours(until.getHours() + durationHours)
      banTarget.value.banned_until = until.toISOString()
    }
    showBanModal.value = false
  } catch (e: any) {
    console.error('Failed to ban user:', e)
    alert('封禁失败: ' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

async function unbanUser(user: UserRecord) {
  if (!confirm(`确定要解除 ${user.display_name || user.username} 的封禁吗？`)) return
  saving.value = true
  try {
    await api.post(`/api/users/${user.id}/unban`)
    user.is_banned = false
    user.ban_reason = ''
    user.banned_at = null
    user.banned_until = null
  } catch (e: any) {
    console.error('Failed to unban user:', e)
    alert('解除封禁失败: ' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

function openCreditsModal(user: UserRecord) {
  creditsTarget.value = user
  creditsAdjust.value = 0
  creditsSet.value = 0
  characterQuotaMb.value = Math.max(1, Math.round((user.character_storage_quota || 10485760) / 1024 / 1024))
  showCreditsModal.value = true
}

async function confirmCreditsAdjust() {
  if (!creditsTarget.value) return
  saving.value = true
  try {
    let newCredits: number
    if (creditsSet.value !== 0) {
      newCredits = creditsSet.value
    } else {
      newCredits = (creditsTarget.value.credits || 0) + creditsAdjust.value
    }
    const nextQuota = Math.max(1, Math.round(characterQuotaMb.value)) * 1024 * 1024
    await api.patch(`/api/users/${creditsTarget.value.id}/credits`, {
      credits: newCredits,
      character_storage_quota: nextQuota,
    })
    creditsTarget.value.credits = newCredits
    creditsTarget.value.character_storage_quota = nextQuota
    showCreditsModal.value = false
  } catch (e: any) {
    console.error('Failed to adjust credits:', e)
    alert('调整额度失败: ' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

async function createUser() {
  if (!createUserForm.value.username || !createUserForm.value.password) {
    alert('请填写用户名和密码')
    return
  }

  saving.value = true
  try {
    await api.post('/api/users', createUserForm.value)
    showCreateUserModal.value = false
    createUserForm.value = { username: '', email: '', password: '', role: 'user', credits: 0 }
    await fetchUsers()
  } catch (e: any) {
    console.error('Failed to create user:', e)
    alert('创建用户失败: ' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

async function impersonateUser(user: UserRecord) {
  if (!confirm(`将以 ${user.display_name || user.username} 身份进入账户，是否继续？`)) return
  saving.value = true
  try {
    const originalSession = {
      id: authStore.user?.id,
      username: authStore.user?.username,
      displayName: authStore.user?.displayName,
      token: localStorage.getItem('token'),
      user: localStorage.getItem('user'),
    }
    const res: any = await api.post(`/api/users/${user.id}/impersonate`)
    sessionStorage.setItem('impersonator_user', JSON.stringify(originalSession))
    authStore.applyAuthPayload(res)
    await router.push('/profile')
  } catch (e: any) {
    console.error('Failed to impersonate user:', e)
    alert('代管登录失败: ' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (authStore.isAdmin) {
    fetchUsers()
  }
})
</script>

<style scoped>
.user-manage-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.page-content::-webkit-scrollbar {
  display: none;
}

.no-access {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  font-size: 16px;
  color: var(--text-secondary);
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
  font-size: 14px;
}

/* Search */
.search-bar {
  display: flex;
  align-items: center;
  background: var(--bg-primary);
  border-radius: 10px;
  padding: 0 12px;
  margin-bottom: 12px;
  border: 1px solid var(--separator);
  box-shadow: var(--shadow-sm);
}

.create-user-btn {
  flex: 0 0 auto;
  border: none;
  border-radius: 10px;
  padding: 8px 12px;
  background: rgba(91, 110, 245, 0.14);
  color: #5B6EF5;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.search-icon {
  width: 18px;
  height: 18px;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 15px;
  padding: 10px 8px;
}

.search-input::placeholder {
  color: var(--text-tertiary);
}

/* Stats bar */
.stats-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.stat-chip {
  flex: 1;
  background: var(--bg-primary);
  border-radius: 10px;
  padding: 10px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  border: 1px solid var(--separator);
  box-shadow: var(--shadow-sm);
}

.stat-chip-label {
  font-size: 11px;
  color: var(--text-secondary);
}

.stat-chip-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-chip-value.stat-banned {
  color: #ff6b6b;
}

/* User list */
.user-list {
  background: var(--bg-primary);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--separator);
  box-shadow: var(--shadow-sm);
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-bottom: 1px solid var(--separator);
  transition: background 0.2s;
}

.user-item:last-child {
  border-bottom: none;
}

.user-item.user-banned {
  background: rgba(255, 59, 48, 0.06);
  border-left: 3px solid rgba(255, 59, 48, 0.5);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 16px;
  font-weight: 700;
  color: var(--bg-primary);
}

.avatar-ban-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  font-size: 14px;
  background: var(--bg-primary);
  border-radius: 50%;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  align-items: center;
  gap: 6px;
}

.ban-tag {
  font-size: 10px;
  padding: 1px 6px;
  background: rgba(255, 59, 48, 0.25);
  color: #ff6b6b;
  border-radius: 4px;
  font-weight: 600;
  flex-shrink: 0;
}

.role-badge {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 4px;
  font-weight: 600;
  flex-shrink: 0;
}

.badge-super {
  background: linear-gradient(135deg, #FFD700, #FFA500);
  color: #1a1a1a;
}

.user-meta {
  display: flex;
  gap: 8px;
  margin-top: 2px;
}

.user-username {
  font-size: 12px;
  color: var(--text-secondary);
}

.user-date {
  font-size: 12px;
  color: var(--text-tertiary);
}

.user-ban-reason {
  font-size: 11px;
  color: rgba(255, 107, 107, 0.7);
  margin-top: 3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Actions */
.user-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.role-select {
  appearance: none;
  -webkit-appearance: none;
  background: var(--bg-primary);
  border: 1px solid var(--separator);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 12px;
  padding: 6px 10px;
  cursor: pointer;
  outline: none;
}

.role-select:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.role-select.role-admin {
  background: rgba(255, 59, 48, 0.2);
  border-color: rgba(255, 59, 48, 0.4);
  color: #ff6b6b;
}

.role-select.role-user {
  background: rgba(52, 199, 89, 0.15);
  border-color: rgba(52, 199, 89, 0.3);
  color: #34c759;
}

.ban-btn, .unban-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.impersonate-btn {
  height: 32px;
  padding: 0 10px;
  border-radius: 8px;
  border: 1px solid rgba(255, 149, 0, 0.25);
  background: rgba(255, 149, 0, 0.12);
  color: #ff9500;
  font-size: 12px;
  cursor: pointer;
}

.ban-btn {
  background: rgba(255, 59, 48, 0.15);
  border: 1px solid rgba(255, 59, 48, 0.3);
}

.ban-btn:hover {
  background: rgba(255, 59, 48, 0.3);
  transform: scale(1.05);
}

.unban-btn {
  background: rgba(52, 199, 89, 0.15);
  border: 1px solid rgba(52, 199, 89, 0.3);
}

.unban-btn:hover {
  background: rgba(52, 199, 89, 0.3);
  transform: scale(1.05);
}

/* Save indicator */
.save-indicator {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--text-primary);
  color: var(--bg-primary);
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 13px;
  z-index: 100;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.ban-modal {
  background: var(--bg-primary);
  border: 1px solid var(--separator);
  border-radius: 20px;
  padding: 24px;
  width: 90%;
  max-width: 380px;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.ban-modal-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.ban-modal-icon {
  font-size: 24px;
}

.ban-modal-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: #ff6b6b;
}

.ban-modal-user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  margin-bottom: 16px;
}

.ban-modal-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ban-modal-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.ban-modal-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.ban-modal-username {
  font-size: 13px;
  color: var(--text-secondary);
}

.ban-modal-field {
  margin-bottom: 12px;
}

.ban-modal-field label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.ban-reason-input {
  width: 100%;
  background: var(--bg-secondary);
  border: 1px solid var(--separator);
  border-radius: 10px;
  color: var(--text-primary);
  font-size: 14px;
  padding: 10px 12px;
  outline: none;
  resize: none;
  font-family: inherit;
  box-sizing: border-box;
}

.ban-reason-input:focus {
  border-color: rgba(255, 59, 48, 0.5);
}

.ban-reason-input::placeholder {
  color: var(--text-tertiary);
}

.ban-modal-warning {
  font-size: 12px;
  color: rgba(255, 149, 0, 0.8);
  background: rgba(255, 149, 0, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  margin-bottom: 16px;
}

.ban-modal-actions {
  display: flex;
  gap: 10px;
}

.btn-cancel {
  flex: 1;
  padding: 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--separator);
  border-radius: 10px;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-cancel:hover {
  background: var(--fill-secondary);
}

.btn-confirm-ban {
  flex: 1;
  padding: 10px;
  background: rgba(255, 59, 48, 0.8);
  border-radius: 10px;
  color: var(--bg-primary);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-confirm-ban:hover {
  background: rgba(255, 59, 48, 1);
}

.ban-duration-select {
  width: 100%;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--separator);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
}

.credits-modal {
  background: var(--bg-primary);
  border-radius: 16px;
  padding: 20px;
  width: 90%;
  max-width: 400px;
  color: var(--text-primary);
}

.create-user-modal {
  max-width: 440px;
}

.credits-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.credits-modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.credits-modal-user {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.credits-modal-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--fill-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.credits-modal-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.credits-modal-name {
  font-weight: 600;
}

.credits-modal-username {
  font-size: 13px;
  color: var(--text-secondary);
}

.credits-modal-current {
  text-align: center;
  padding: 12px;
  background: rgba(91, 110, 245, 0.1);
  border-radius: 8px;
  margin-bottom: 16px;
}

.credits-value {
  color: #5B6EF5;
  font-weight: 700;
  font-size: 18px;
}

.credits-modal-field {
  margin-bottom: 16px;
}

.credits-modal-field label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.credits-input-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: stretch;
}

.credits-btn {
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
  transition: opacity 0.2s;
  flex: 0 0 auto;
  min-width: 56px;
}

.credits-btn.minus {
  background: rgba(255, 59, 48, 0.2);
  color: #ff3b30;
}

.credits-btn.plus {
  background: rgba(52, 199, 89, 0.2);
  color: #34c759;
}

.credits-btn:hover {
  opacity: 0.8;
}

.credits-input {
  flex: 1 1 100%;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--separator);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  text-align: center;
  min-width: 0;
}

.credits-modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.btn-confirm-credits {
  flex: 1;
  padding: 10px;
  background: rgba(91, 110, 245, 0.8);
  border-radius: 10px;
  color: var(--bg-primary);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.btn-confirm-credits:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.user-credits {
  cursor: pointer;
  color: #5B6EF5;
  font-size: 12px;
}

.user-credits:hover {
  text-decoration: underline;
}

.banned-until {
  color: rgba(255, 149, 0, 0.8);
  font-size: 11px;
}
</style>
