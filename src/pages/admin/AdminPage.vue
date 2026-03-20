<template>
  <div class="admin-page">
    <NavBar title="管理面板" :show-back="true" back-to="/" />

    <div class="admin-content">
      <!-- 管理员信息 -->
      <div class="admin-header">
        <div class="admin-avatar">
          {{ user?.username?.charAt(0)?.toUpperCase() || 'A' }}
        </div>
        <div class="admin-info">
          <div class="admin-name">{{ user?.username }}</div>
          <div class="admin-role">{{ isSuperAdmin ? '超级管理员' : '管理员' }}</div>
        </div>
      </div>

      <!-- 功能卡片 -->
      <div class="admin-grid">
        <div class="admin-card" @click="router.push('/admin/features')">
          <div class="card-icon" style="background: linear-gradient(135deg, #FF6B6B, #EE5A24)">■</div>
          <div class="card-title">功能管理</div>
          <div class="card-desc">开关功能模块</div>
        </div>

        <div class="admin-card" @click="router.push('/admin/users')">
          <div class="card-icon" style="background: linear-gradient(135deg, #E6162D, #FF4757)">○</div>
          <div class="card-title">用户管理</div>
          <div class="card-desc">管理注册用户</div>
        </div>

        <div class="admin-card" @click="router.push('/admin/stats')">
          <div class="card-icon" style="background: linear-gradient(135deg, #3498DB, #2980B9)">▥</div>
          <div class="card-title">数据统计</div>
          <div class="card-desc">查看运营数据</div>
        </div>

<div class="admin-card" @click="router.push('/admin/settings')">
  <div class="card-icon" style="background: linear-gradient(135deg, #636E72, #95A5A6)">⚙</div>
  <div class="card-title">系统设置</div>
  <div class="card-desc">全局配置和API</div>
</div>

<div class="admin-card" @click="router.push('/admin/logs')">
  <div class="card-icon" style="background: linear-gradient(135deg, #9B59B6, #8E44AD)">📋</div>
  <div class="card-title">系统日志</div>
  <div class="card-desc">查看操作记录</div>
</div>
      </div>

      <!-- 快速操作 -->
      <div class="quick-actions">
        <div class="action-title">快速操作</div>
        <div class="action-list">
          <div class="action-item" @click="router.push('/characters')">
            <span class="action-icon">👤</span>
            <span>角色管理</span>
            <span class="action-arrow">›</span>
          </div>
          <div class="action-item" @click="router.push('/worldbook')">
            <span class="action-icon">📚</span>
            <span>世界书</span>
            <span class="action-arrow">›</span>
          </div>
          <div class="action-item" @click="router.push('/preset')">
            <span class="action-icon">⚡</span>
            <span>预设管理</span>
            <span class="action-arrow">›</span>
          </div>
          <div class="action-item" @click="router.push('/')">
            <span class="action-icon">🏠</span>
            <span>返回首页</span>
            <span class="action-arrow">›</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/common/NavBar.vue'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const isSuperAdmin = computed(() => authStore.isSuperAdmin)
</script>

<style scoped>
.admin-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #f7f9fc 0%, #eef3f8 100%);
  color: #1f2937;
}

.admin-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.admin-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, #ffffff, #f4f8fc);
  border-radius: 16px;
  margin-bottom: 20px;
  border: 1px solid rgba(148, 163, 184, 0.18);
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}

.admin-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.admin-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.admin-role {
  font-size: 13px;
  color: #64748b;
}

.admin-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.admin-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

.admin-card:active {
  transform: scale(0.98);
}

.card-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-bottom: 10px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 4px;
}

.card-desc {
  font-size: 12px;
  color: #64748b;
}

.quick-actions {
  background: #ffffff;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

.action-title {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 12px;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 10px;
  cursor: pointer;
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.action-item:active {
  background: #eef4ff;
}

.action-icon {
  font-size: 18px;
}

.action-item span:nth-child(2) {
  flex: 1;
  font-size: 14px;
}

.action-arrow {
  color: #94a3b8;
  font-size: 20px;
}
</style>
