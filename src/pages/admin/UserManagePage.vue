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
            placeholder="搜索用户名或ID..."
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
            <span class="stat-chip-label">普通</span>
            <span class="stat-chip-value">{{ roleCounts.user }}</span>
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
          >
            <div class="user-avatar">
              <img v-if="u.avatar_url" :src="u.avatar_url" alt="" />
              <span v-else class="avatar-placeholder">{{ (u.display_name || u.username)?.[0] || '?' }}</span>
            </div>
            <div class="user-info">
              <div class="user-name">{{ u.display_name || u.username }}</div>
              <div class="user-meta">
                <span class="user-username">@{{ u.username }}</span>
                <span class="user-date">{{ formatDate(u.created_at) }}</span>
              </div>
            </div>
            <div class="role-selector">
              <select
                :value="u.role"
                @change="changeRole(u, ($event.target as HTMLSelectElement).value)"
                class="role-select"
                :class="'role-' + u.role"
              >
                <option value="user">普通用户</option>
                <option value="moderator">版主</option>
                <option value="admin">管理员</option>
              </select>
            </div>
          </div>
        </div>

        <div v-if="saving" class="save-indicator">保存中...</div>
      </template>
    </div>
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
  created_at: string
  updated_at: string
}

const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const users = ref<UserRecord[]>([])

const roleCounts = computed(() => {
  const counts = { admin: 0, moderator: 0, user: 0 }
  for (const u of users.value) {
    if (u.role in counts) {
      counts[u.role as keyof typeof counts]++
    }
  }
  return counts
})

const filteredUsers = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return users.value
  return users.value.filter(
    (u) =>
      u.username.toLowerCase().includes(q) ||
      u.display_name?.toLowerCase().includes(q) ||
      String(u.id).includes(q)
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
}

.user-item:last-child {
  border-bottom: none;
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

/* Role selector */
.role-selector {
  flex-shrink: 0;
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
</style>
