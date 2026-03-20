<template>
  <div class="data-stats-page">
    <nav-bar title="数据统计" back-route="/profile" />

    <div class="page-content">
      <div v-if="!authStore.isAdmin" class="no-access">
        <p>⛔ 需要管理员权限</p>
      </div>

      <template v-else>
        <div v-if="loading" class="loading">加载中...</div>

        <template v-else>
          <!-- 总览卡片 -->
          <div class="overview-cards">
            <div class="overview-card">
              <div class="card-icon" style="background: linear-gradient(135deg, #5B6EF5, #8B5CF6)">○</div>
              <div class="card-info">
                <span class="card-value">{{ totalUsers }}</span>
                <span class="card-label">总用户数</span>
              </div>
            </div>
            <div class="overview-card">
              <div class="card-icon" style="background: linear-gradient(135deg, #00B894, #55EFC4)">▦</div>
              <div class="card-info">
                <span class="card-value">{{ todayNewUsers }}</span>
                <span class="card-label">今日新增</span>
              </div>
            </div>
            <div class="overview-card">
              <div class="card-icon" style="background: linear-gradient(135deg, #FDCB6E, #F39C12)">↑</div>
              <div class="card-info">
                <span class="card-value">{{ weekNewUsers }}</span>
                <span class="card-label">本周新增</span>
              </div>
            </div>
            <div class="overview-card">
              <div class="card-icon" style="background: linear-gradient(135deg, #E17055, #FAB1A0)">▥</div>
              <div class="card-info">
                <span class="card-value">{{ monthNewUsers }}</span>
                <span class="card-label">本月新增</span>
              </div>
            </div>
          </div>

          <!-- 角色分布 -->
          <div class="section">
            <div class="section-header">
              <span class="section-title">角色分布</span>
            </div>
            <div class="section-card">
              <div v-for="role in roleDistribution" :key="role.name" class="role-item">
                <div class="role-info">
                  <span class="role-dot" :style="{ background: role.color }"></span>
                  <span class="role-name">{{ role.label }}</span>
                  <span class="role-count">{{ role.count }}</span>
                </div>
                <div class="role-bar-bg">
                  <div
                    class="role-bar-fill"
                    :style="{ width: role.percent + '%', background: role.color }"
                  ></div>
                </div>
                <span class="role-percent">{{ role.percent.toFixed(1) }}%</span>
              </div>
            </div>
          </div>

          <!-- 注册趋势 (最近30天) -->
          <div class="section">
            <div class="section-header">
              <span class="section-title">注册趋势 (最近30天)</span>
            </div>
            <div class="section-card chart-card">
              <div class="chart-container">
                <div class="chart-y-axis">
                  <span v-for="label in yAxisLabels" :key="label">{{ label }}</span>
                </div>
                <div class="chart-bars">
                  <div
                    v-for="day in registrationTrend"
                    :key="day.date"
                    class="chart-bar-wrapper"
                    :title="day.date + ': ' + day.count + '人'"
                  >
                    <div
                      class="chart-bar"
                      :style="{ height: day.barHeight + '%' }"
                    ></div>
                    <span class="chart-bar-label">{{ day.shortDate }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 最近注册用户 -->
          <div class="section">
            <div class="section-header">
              <span class="section-title">最近注册</span>
            </div>
            <div class="section-card">
              <div v-for="u in recentUsers" :key="u.id" class="recent-user-item">
                <div class="recent-user-avatar">
                  <img v-if="u.avatar_url" :src="u.avatar_url" alt="" />
                  <span v-else>{{ (u.display_name || u.username || '?')[0] }}</span>
                </div>
                <div class="recent-user-info">
                  <span class="recent-user-name">{{ u.display_name || u.username }}</span>
                  <span class="recent-user-date">{{ formatDate(u.created_at) }}</span>
                </div>
                <span class="recent-user-role" :class="'role-' + u.role">{{ roleLabel(u.role) }}</span>
              </div>
              <div v-if="recentUsers.length === 0" class="empty-hint">暂无数据</div>
            </div>
          </div>
        </template>
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

const loading = ref(false)

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

const users = ref<UserRecord[]>([])

const totalUsers = computed(() => users.value.length)

const todayNewUsers = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return users.value.filter(u => new Date(u.created_at) >= today).length
})

const weekNewUsers = computed(() => {
  const now = new Date()
  const weekAgo = new Date(now.getTime() - 7 * 86400000)
  weekAgo.setHours(0, 0, 0, 0)
  return users.value.filter(u => new Date(u.created_at) >= weekAgo).length
})

const monthNewUsers = computed(() => {
  const now = new Date()
  const monthAgo = new Date(now.getFullYear(), now.getMonth(), 1)
  return users.value.filter(u => new Date(u.created_at) >= monthAgo).length
})

const roleColors: Record<string, string> = {
  admin: '#E6162D',
  moderator: '#F39C12',
  user: '#3498DB',
}

const roleLabels: Record<string, string> = {
  admin: '管理员',
  moderator: '版主',
  user: '普通用户',
}

function roleLabel(role: string) {
  return roleLabels[role] || role
}

const roleDistribution = computed(() => {
  const counts: Record<string, number> = {}
  for (const u of users.value) {
    counts[u.role] = (counts[u.role] || 0) + 1
  }
  const total = users.value.length || 1
  return Object.entries(counts).map(([name, count]) => ({
    name,
    label: roleLabels[name] || name,
    count,
    percent: (count / total) * 100,
    color: roleColors[name] || '#888',
  })).sort((a, b) => b.count - a.count)
})

const registrationTrend = computed(() => {
  const days: { date: string; shortDate: string; count: number; barHeight: number }[] = []
  const now = new Date()
  const countMap: Record<string, number> = {}

  for (const u of users.value) {
    const d = new Date(u.created_at)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    countMap[key] = (countMap[key] || 0) + 1
  }

  for (let i = 29; i >= 0; i--) {
    const d = new Date(now.getTime() - i * 86400000)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    days.push({
      date: key,
      shortDate: `${d.getMonth() + 1}/${d.getDate()}`,
      count: countMap[key] || 0,
      barHeight: 0,
    })
  }

  const maxCount = Math.max(...days.map(d => d.count), 1)
  for (const day of days) {
    day.barHeight = (day.count / maxCount) * 100
  }

  return days
})

const yAxisLabels = computed(() => {
  const maxCount = Math.max(...registrationTrend.value.map(d => d.count), 1)
  return [maxCount, Math.round(maxCount * 0.5), 0].map(String)
})

const recentUsers = computed(() => {
  return [...users.value]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 10)
})

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

onMounted(async () => {
  if (!authStore.isAdmin) return
  loading.value = true
  try {
    const res: any = await api.get('/api/users')
    users.value = Array.isArray(res) ? res : (res.data || [])
  } catch (e: any) {
    console.error('Failed to load users:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.data-stats-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(180deg, #f7f9fc 0%, #eef3f8 100%);
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
  color: #64748b;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #64748b;
}

/* Overview Cards */
.overview-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 24px;
}

.overview-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

.card-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.card-info {
  display: flex;
  flex-direction: column;
}

.card-value {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
}

.card-label {
  font-size: 12px;
  color: #64748b;
}

/* Section */
.section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  padding: 0 4px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-card {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

/* Role Distribution */
.role-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid #e5e7eb;
}

.role-item:last-child {
  border-bottom: none;
}

.role-info {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 110px;
}

.role-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.role-name {
  font-size: 14px;
  color: #0f172a;
}

.role-count {
  font-size: 12px;
  color: #64748b;
  margin-left: auto;
}

.role-bar-bg {
  flex: 1;
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  overflow: hidden;
}

.role-bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.5s ease;
}

.role-percent {
  font-size: 12px;
  color: #64748b;
  min-width: 45px;
  text-align: right;
}

/* Chart */
.chart-card {
  padding: 16px;
}

.chart-container {
  display: flex;
  gap: 8px;
  height: 160px;
}

.chart-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  font-size: 10px;
  color: #94a3b8;
  padding-bottom: 18px;
  min-width: 20px;
  text-align: right;
}

.chart-bars {
  flex: 1;
  display: flex;
  align-items: flex-end;
  gap: 2px;
  border-bottom: 1px solid #dbe3ee;
  padding-bottom: 18px;
  position: relative;
}

.chart-bar-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  height: 100%;
  position: relative;
}

.chart-bar {
  width: 100%;
  max-width: 12px;
  background: linear-gradient(180deg, #5B6EF5, #8B5CF6);
  border-radius: 2px 2px 0 0;
  min-height: 1px;
  transition: height 0.5s ease;
}

.chart-bar-label {
  position: absolute;
  bottom: -18px;
  font-size: 8px;
  color: rgba(255, 255, 255, 0.3);
  white-space: nowrap;
  transform: rotate(-45deg);
  transform-origin: top center;
}

.chart-bar-wrapper:nth-child(odd) .chart-bar-label {
  display: none;
}

/* Recent Users */
.recent-user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.recent-user-item:last-child {
  border-bottom: none;
}

.recent-user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.recent-user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recent-user-avatar span {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.recent-user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.recent-user-name {
  font-size: 14px;
  color: #fff;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recent-user-date {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
}

.recent-user-role {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
  flex-shrink: 0;
}

.role-admin {
  background: rgba(230, 22, 45, 0.15);
  color: #ff6b7a;
}

.role-moderator {
  background: rgba(243, 156, 18, 0.15);
  color: #f5c842;
}

.role-user {
  background: rgba(52, 152, 219, 0.15);
  color: #5dade2;
}

.empty-hint {
  text-align: center;
  padding: 24px;
  color: rgba(255, 255, 255, 0.3);
  font-size: 14px;
}
</style>
