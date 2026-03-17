<template>
  <div class="leaderboard-page">
    <NavBar title="综合排行榜" :show-back="true" />
    
    <div class="leaderboard-content">
      <div class="leaderboard-list">
        <div 
          class="leaderboard-item" 
          v-for="(user, index) in leaderboard" 
          :key="user.id"
          @click="viewUserProfile(user.id)"
        >
          <div class="rank-container">
            <div class="rank-badge" :class="getRankClass(index)">
              {{ index + 1 }}
            </div>
          </div>
          
          <div class="user-avatar">
            <img v-if="user.avatar_url" :src="user.avatar_url" :alt="user.username" />
            <span v-else class="avatar-placeholder">{{ user.username.charAt(0) }}</span>
          </div>
          
          <div class="user-info">
            <div class="user-name">{{ user.display_name || user.username }}</div>
            <div class="user-stats">
              <span class="token-count">{{ user.total_tokens }} Tokens</span>
            </div>
          </div>
          
          <div class="online-status">
            <div class="online-dot" :class="{ online: isUserOnline(user.last_active_at) }"></div>
          </div>
        </div>
        
        <div v-if="leaderboard.length === 0" class="empty-state">
          <div class="empty-icon">🏆</div>
          <p>暂无排行榜数据</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/common/NavBar.vue'

const router = useRouter()
const authStore = useAuthStore()

interface LeaderboardUser {
  id: number
  username: string
  display_name: string
  avatar_url: string
  total_tokens: number
  last_active_at: string
  created_at: string
}

const leaderboard = ref<LeaderboardUser[]>([])

onMounted(async () => {
  await loadLeaderboard()
})

async function loadLeaderboard() {
  try {
    const response = await fetch('/api/leaderboard', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (response.ok) {
      const data = await response.json()
      leaderboard.value = data.data || []
    }
  } catch (error) {
    console.error('Failed to load leaderboard:', error)
  }
}

function getRankClass(index: number) {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return 'normal'
}

function viewUserProfile(userId: number) {
  // TODO: 实现查看用户详情功能
  console.log('View user profile:', userId)
}

function isUserOnline(lastActive: string): boolean {
  const lastActiveTime = new Date(lastActive).getTime()
  const fiveMinutesAgo = Date.now() - 5 * 60 * 1000
  return lastActiveTime > fiveMinutesAgo
}
</script>

<style scoped>
.leaderboard-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

.leaderboard-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.leaderboard-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.leaderboard-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: background 0.15s;
}

.leaderboard-item:active {
  background: var(--bg-tertiary);
}

.rank-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  flex-shrink: 0;
}

.rank-badge {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
}

.rank-badge.gold {
  background: linear-gradient(135deg, #FFD700, #FFA500);
}

.rank-badge.silver {
  background: linear-gradient(135deg, #C0C0C0, #A9A9A9);
}

.rank-badge.bronze {
  background: linear-gradient(135deg, #CD7F32, #A0522D);
}

.rank-badge.normal {
  background: var(--bg-quaternary);
  color: var(--text-primary);
}

.user-avatar {
  width: 48px;
  height: 48px;
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
  font-size: 20px;
  font-weight: 700;
  color: #fff;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-stats {
  display: flex;
  align-items: center;
  gap: 8px;
}

.token-count {
  font-size: 14px;
  color: var(--text-secondary);
}

.online-status {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.online-dot {
  width: 10px;
  height: 10px;
  background: var(--bg-quaternary);
  border-radius: 50%;
}

.online-dot.online {
  background: var(--color-green);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-tertiary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}
</style>