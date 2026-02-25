import { defineStore } from 'pinia'
import { ref } from 'vue'
import { momentApi } from '@/api/services'
import type { Moment } from '@/api/types'

export const useMomentsStore = defineStore('moments', () => {
  const moments = ref<Moment[]>([])
  const loading = ref(false)

  async function fetchMoments() {
    loading.value = true
    try {
      const res = await momentApi.list()
      moments.value = res.data || []
    } catch (e) {
      console.error('获取朋友圈失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function createMoment(data: { content: string; images?: string[]; location?: string; visibility?: string }) {
    try {
      await momentApi.create(data)
      await fetchMoments()
    } catch (e) {
      console.error('发布朋友圈失败:', e)
    }
  }

  async function deleteMoment(id: number) {
    try {
      await momentApi.delete(id)
      moments.value = moments.value.filter(m => m.id !== id)
    } catch (e) {
      console.error('删除朋友圈失败:', e)
    }
  }

  return {
    moments,
    loading,
    fetchMoments,
    createMoment,
    deleteMoment,
  }
})
