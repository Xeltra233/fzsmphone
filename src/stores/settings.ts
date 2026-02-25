import { defineStore } from 'pinia'
import { ref } from 'vue'
import { settingsApi } from '@/api/services'
import type { AppSettings } from '@/api/types'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<AppSettings>({})
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSettings() {
    loading.value = true
    error.value = null
    try {
      const res = await settingsApi.get()
      settings.value = res.data || {}
    } catch (e: any) {
      error.value = e.message || '加载设置失败'
    } finally {
      loading.value = false
    }
  }

  async function updateSettings(data: AppSettings) {
    await settingsApi.update(data)
    // 合并更新到本地
    Object.assign(settings.value, data)
  }

  function getSetting<T = unknown>(key: string, defaultValue?: T): T {
    return (settings.value[key] as T) ?? (defaultValue as T)
  }

  return {
    settings,
    loading,
    error,
    fetchSettings,
    updateSettings,
    getSetting,
  }
})
