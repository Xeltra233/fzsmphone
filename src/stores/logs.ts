import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/client'

export interface LogEntry {
  id: number
  level: string
  action: string
  user_id?: number
  username?: string
  details?: Record<string, any>
  ip_address?: string
  created_at: string
}

export interface LogStats {
  today_count: number
  yesterday_count: number
  week_count: number
  error_count: number
}

export interface LogActionCount {
  action: string
  count: number
}

export interface LogsResponse {
  logs: LogEntry[]
  total: number
  limit: number
  offset: number
}

interface StatsResponse {
  stats: LogStats
  action_counts: LogActionCount[]
}

export const useLogsStore = defineStore('logs', () => {
  const logs = ref<LogEntry[]>([])
  const total = ref(0)
  const stats = ref<LogStats | null>(null)
  const actionCounts = ref<LogActionCount[]>([])
  const loading = ref(false)

  async function fetchLogs(options: {
    action?: string
    level?: string
    user_id?: number
    limit?: number
    offset?: number
  } = {}) {
    loading.value = true
    try {
      const params = new URLSearchParams()
      if (options.action) params.set('action', options.action)
      if (options.level) params.set('level', options.level)
      if (options.user_id) params.set('user_id', String(options.user_id))
      if (options.limit) params.set('limit', String(options.limit))
      if (options.offset) params.set('offset', String(options.offset))

      const query = params.toString() ? `?${params.toString()}` : ''
      const data = await api.get<LogsResponse>(`/api/logs${query}`)
      logs.value = data.logs
      total.value = data.total
    } catch (e) {
      console.error('Failed to fetch logs:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchStats() {
    try {
      const data = await api.get<StatsResponse>('/api/logs/stats')
      stats.value = data.stats
      actionCounts.value = data.action_counts
    } catch (e) {
      console.error('Failed to fetch log stats:', e)
    }
  }

async function clearLogs(days: number = 30) {
  return await api.delete<{ message: string; deleted_rows: number }>(`/api/logs?days=${days}`)
}

  return {
    logs,
    total,
    stats,
    actionCounts,
    loading,
    fetchLogs,
    fetchStats,
    clearLogs,
  }
})