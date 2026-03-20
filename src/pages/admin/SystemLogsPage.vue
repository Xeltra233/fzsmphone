<template>
  <div class="logs-page">
    <NavBar title="系统日志" :show-back="true" back-to="/admin" />

    <div class="logs-content">
      <!-- 统计卡片 -->
      <div class="stats-row">
        <div class="stat-item">
          <span class="stat-value">{{ stats?.today_count || 0 }}</span>
          <span class="stat-label">今日</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ stats?.yesterday_count || 0 }}</span>
          <span class="stat-label">昨日</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ stats?.week_count || 0 }}</span>
          <span class="stat-label">本周</span>
        </div>
        <div class="stat-item error">
          <span class="stat-value">{{ stats?.error_count || 0 }}</span>
          <span class="stat-label">错误</span>
        </div>
      </div>

      <!-- 筛选器 -->
      <div class="filters">
        <select v-model="filterLevel" class="filter-select" @change="applyFilter">
          <option value="">全部级别</option>
          <option value="INFO">INFO</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
        </select>
        <select v-model="filterAction" class="filter-select" @change="applyFilter">
          <option value="">全部操作</option>
          <option v-for="ac in actionCounts" :key="ac.action" :value="ac.action">
            {{ ac.action }}
          </option>
        </select>
        <button class="refresh-btn" @click="refreshLogs">刷新</button>
      </div>

      <!-- 日志列表 -->
      <div class="logs-list" ref="logsListRef">
        <div
          v-for="log in logs"
          :key="log.id"
          class="log-item"
          :class="log.level.toLowerCase()"
        >
          <div class="log-header">
            <span class="log-level" :class="log.level.toLowerCase()">{{ log.level }}</span>
            <span class="log-action">{{ log.action }}</span>
            <span class="log-time">{{ formatTime(log.created_at) }}</span>
          </div>
          <div class="log-details" v-if="log.details">
            <pre>{{ formatDetails(log.details) }}</pre>
          </div>
          <div class="log-meta" v-if="log.username || log.ip_address">
            <span v-if="log.username">👤 {{ log.username }}</span>
            <span v-if="log.ip_address">🌐 {{ log.ip_address }}</span>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading">加载中...</div>
        
        <!-- 空状态 -->
        <div v-if="!loading && logs.length === 0" class="empty">
          暂无日志
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="logs.length < total" class="load-more">
        <button @click="loadMore" :disabled="loading">加载更多</button>
      </div>

      <!-- 清理日志 -->
      <div class="clear-section">
        <button class="clear-btn" @click="handleClearLogs">清理30天前日志</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useLogsStore, type LogEntry } from '@/stores/logs'
import NavBar from '@/components/common/NavBar.vue'
import { showToast } from '@/utils/toast'

const logsStore = useLogsStore()

const logs = computed(() => logsStore.logs)
const total = computed(() => logsStore.total)
const stats = computed(() => logsStore.stats)
const actionCounts = computed(() => logsStore.actionCounts)
const loading = computed(() => logsStore.loading)

const filterLevel = ref('')
const filterAction = ref('')
const offset = ref(0)
const limit = 20
const logsListRef = ref<HTMLElement>()

async function loadData(reset = false) {
  if (reset) {
    offset.value = 0
  }
  
  await logsStore.fetchLogs({
    level: filterLevel.value || undefined,
    action: filterAction.value || undefined,
    limit,
    offset: offset.value,
  })
}

async function refreshLogs() {
  await loadData(true)
}

async function loadMore() {
  offset.value += limit
  await logsStore.fetchLogs({
    level: filterLevel.value || undefined,
    action: filterAction.value || undefined,
    limit,
    offset: offset.value,
  })
}

function applyFilter() {
  loadData(true)
}

async function handleClearLogs() {
  if (!confirm('确定要清理30天前的日志吗？此操作不可恢复。')) {
    return
  }
  
  try {
    const result = await logsStore.clearLogs(30)
    showToast(`已清理 ${result.deleted_rows} 条日志`)
    await refreshLogs()
  } catch (e: any) {
    showToast(e.message || '清理失败')
  }
}

function formatTime(time: string): string {
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatDetails(details: Record<string, any>): string {
  try {
    return JSON.stringify(details, null, 2)
  } catch {
    return String(details)
  }
}

onMounted(async () => {
  await logsStore.fetchStats()
  await loadData()
})
</script>

<style scoped>
.logs-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #f7f9fc 0%, #eef3f8 100%);
  color: #0f172a;
}

.logs-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 16px;
}

.stat-item {
  background: #ffffff;
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  border: 1px solid rgba(148, 163, 184, 0.16);
}

.stat-item.error .stat-value {
  color: #ff6b6b;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
}

.stat-label {
  font-size: 11px;
  color: #64748b;
}

.filters {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.filter-select {
  flex: 1;
  padding: 8px 12px;
  background: #ffffff;
  border: 1px solid #d7dee8;
  border-radius: 6px;
  color: #0f172a;
  font-size: 13px;
}

.refresh-btn {
  padding: 8px 16px;
  background: #ffffff;
  border: 1px solid #d7dee8;
  border: none;
  border-radius: 6px;
  color: #0f172a;
  font-size: 13px;
  cursor: pointer;
}

.logs-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-item {
  background: #ffffff;
  border-radius: 8px;
  padding: 12px;
  border-left: 3px solid #666;
  border-top: 1px solid rgba(148, 163, 184, 0.12);
  border-right: 1px solid rgba(148, 163, 184, 0.12);
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
}

.log-item.info {
  border-left-color: #3498db;
}

.log-item.warn {
  border-left-color: #f39c12;
}

.log-item.error {
  border-left-color: #e74c3c;
}

.log-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.log-level {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
}

.log-level.info {
  background: rgba(52, 152, 219, 0.2);
  color: #3498db;
}

.log-level.warn {
  background: rgba(243, 156, 18, 0.2);
  color: #f39c12;
}

.log-level.error {
  background: rgba(231, 76, 60, 0.2);
  color: #e74c3c;
}

.log-action {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: #0f172a;
}

.log-time {
  font-size: 11px;
  color: #94a3b8;
}

.log-details pre {
  margin: 8px 0;
  padding: 8px;
  background: #f8fafc;
  border-radius: 4px;
  font-size: 11px;
  color: #475569;
  overflow-x: auto;
}

.log-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #64748b;
  margin-top: 6px;
}

.loading,
.empty {
  text-align: center;
  padding: 24px;
  color: #64748b;
}

.load-more {
  text-align: center;
  padding: 16px;
}

.load-more button {
  padding: 10px 24px;
  background: #ffffff;
  border: 1px solid #d7dee8;
  border: none;
  border-radius: 6px;
  color: #0f172a;
  cursor: pointer;
}

.clear-section {
  margin-top: 24px;
  text-align: center;
}

.clear-btn {
  padding: 12px 24px;
  background: rgba(231, 76, 60, 0.2);
  border: 1px solid rgba(231, 76, 60, 0.4);
  border-radius: 6px;
  color: #e74c3c;
  font-size: 13px;
  cursor: pointer;
}
</style>
