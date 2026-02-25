import { defineStore } from 'pinia'
import { ref } from 'vue'
import { diaryApi } from '@/api/services'
import type { Diary, DiaryInput } from '@/api/types'

export const useDiaryStore = defineStore('diary', () => {
  const diaries = ref<Diary[]>([])
  const currentDiary = ref<Diary | null>(null)
  const loading = ref(false)

  async function fetchDiaries() {
    loading.value = true
    try {
      const res = await diaryApi.list()
      diaries.value = res.data || []
    } catch (e) {
      console.error('获取日记失败:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchDiary(id: number) {
    try {
      currentDiary.value = await diaryApi.get(id)
    } catch (e) {
      console.error('获取日记详情失败:', e)
    }
  }

  async function createDiary(data: DiaryInput) {
    try {
      const res = await diaryApi.create(data)
      await fetchDiaries()
      return res.id
    } catch (e) {
      console.error('创建日记失败:', e)
      return null
    }
  }

  async function updateDiary(id: number, data: DiaryInput) {
    try {
      await diaryApi.update(id, data)
      const idx = diaries.value.findIndex(d => d.id === id)
      if (idx !== -1) {
        diaries.value[idx] = { ...diaries.value[idx], ...data } as Diary
      }
      if (currentDiary.value?.id === id) {
        currentDiary.value = { ...currentDiary.value, ...data } as Diary
      }
    } catch (e) {
      console.error('更新日记失败:', e)
    }
  }

  async function deleteDiary(id: number) {
    try {
      await diaryApi.delete(id)
      diaries.value = diaries.value.filter(d => d.id !== id)
      if (currentDiary.value?.id === id) {
        currentDiary.value = null
      }
    } catch (e) {
      console.error('删除日记失败:', e)
    }
  }

  return {
    diaries,
    currentDiary,
    loading,
    fetchDiaries,
    fetchDiary,
    createDiary,
    updateDiary,
    deleteDiary,
  }
})
