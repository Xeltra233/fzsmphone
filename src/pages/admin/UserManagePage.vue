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
            <span class="stat-chip-label">版主</span>
            <span class="stat-chip-value">{{ roleCounts.moderator }}</span>
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
                <span v-if="u.is_banned" class="ban-tag">已封禁</span>
              </div>
              <div class="user-meta">
                <span class="user-username">@{{ u.username }}</span>
                <span class="user-date">{{ formatDate(u.created_at) }}</span>
              </div>
              <div v-if="u.is_banned && u.ban_reason" class="user-ban-reason">
                原因：{{ u.ban_reason }}
              </div>
            </div>
            <div class="user-actions">
              <select
                :value="u.role"
                @change="changeRole(u, ($event.target as HTMLSelectElement).value)"
                class="role-select"
                :class="'role-' + u.role"
                :disabled="u.is_banned"
              >
                <option value="user">普通用户</option>
                <option value="moderator">版主</option>
                <option value="admin">管理员</option>
              </select>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api/client'

const authStore = useAuthStore()

interface UserRecord {
  id: number
  discord_id: string
  username: string
  display_name: string
  avatar_url: string
  role: string
  is_banned: boolean
  ban_reason: string
  banned_at: string | null
  created_at: string
  updated_at: string
}

const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const users = ref<UserRecord[]>([])

// Ban modal state
const showBanModal = ref(false)
const banTarget = ref<UserRecord | null>(null)
const banReason = ref('')

const roleCounts = computed(() => {
  const counts = { admin: 0, moderator: 0, user: 0 }
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
    await api.post(`/api/users/${banTarget.value.id}/ban`, { reason: banReason.value })
    banTarget.value.is_banned = true
    banTarget.value.ban_reason = banReason.value
    banTarget.value.banned_at = new Date().toISOString()
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
  } catch (e: any) {
    console.error('Failed to unban user:', e)
    alert('解除封禁失败: ' + (e.message || '未知错误'))
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
  background: var(--color-bg, #000);
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
  color: rgba(255, 255, 255, 0.5);
}

.loading {
  text-align: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.5);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.4);
  font-size: 14px;
}

/* Search */
.search-bar {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 0 12px;
  margin-bottom: 12px;
}

.search-icon {
  width: 18px;
  height: 18px;
  color: rgba(255, 255, 255, 0.35);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  color: #fff;
  font-size: 15px;
  padding: 10px 8px;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

/* Stats bar */
.stats-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.stat-chip {
  flex: 1;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 10px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.stat-chip-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.stat-chip-value {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.stat-chip-value.stat-banned {
  color: #ff6b6b;
}

/* User list */
.user-list {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  overflow: hidden;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
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
  color: #fff;
}

.avatar-ban-badge {
  position: absolute;
  bottom: -2px;
  right: -2px;
  font-size: 14px;
  background: #1a1a2e;
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
  color: #fff;
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

.user-meta {
  display: flex;
  gap: 8px;
  margin-top: 2px;
}

.user-username {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.user-date {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.25);
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
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  color: #fff;
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

.role-select.role-moderator {
  background: rgba(255, 149, 0, 0.2);
  border-color: rgba(255, 149, 0, 0.4);
  color: #ffb347;
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
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
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
  background: #1a1a2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
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
  color: #fff;
}

.ban-modal-username {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
}

.ban-modal-field {
  margin-bottom: 12px;
}

.ban-modal-field label {
  display: block;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 6px;
}

.ban-reason-input {
  width: 100%;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 10px;
  color: #fff;
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
  color: rgba(255, 255, 255, 0.25);
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
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 10px;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.12);
}

.btn-confirm-ban {
  flex: 1;
  padding: 10px;
  background: rgba(255, 59, 48, 0.8);
  border: none;
  border-radius: 10px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-confirm-ban:hover {
  background: rgba(255, 59, 48, 1);
}
</style>
