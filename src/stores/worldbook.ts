import { defineStore } from 'pinia'
import { ref } from 'vue'
import { worldBookApi } from '@/api/services'
import type { WorldBookEntry, WorldBookInput } from '@/api/types'

export const useWorldBookStore = defineStore('worldbook', () => {
  const entries = ref<WorldBookEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchEntries() {
    loading.value = true
    error.value = null
    try {
      const res = await worldBookApi.list()
      entries.value = res.data || []
    } catch (e: any) {
      error.value = e.message || '加载世界书失败'
    } finally {
      loading.value = false
    }
  }

  async function createEntry(data: WorldBookInput) {
    const res = await worldBookApi.create(data)
    await fetchEntries()
    return res
  }

  async function updateEntry(id: number, data: WorldBookInput) {
    await worldBookApi.update(id, data)
    await fetchEntries()
  }

  async function deleteEntry(id: number) {
    await worldBookApi.delete(id)
    entries.value = entries.value.filter(e => e.id !== id)
  }

  async function toggleEnabled(id: number) {
    const entry = entries.value.find(e => e.id === id)
    if (!entry) return
    await worldBookApi.update(id, {
      key: entry.key,
      content: entry.content,
      keywords: entry.keywords,
      is_enabled: !entry.is_enabled,
      priority: entry.priority,
    })
    entry.is_enabled = !entry.is_enabled
  }

  return {
    entries,
    loading,
    error,
    fetchEntries,
    createEntry,
    updateEntry,
    deleteEntry,
    toggleEnabled,
  }
})
