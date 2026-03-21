import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import api from '@/api/client'

export interface ManagedModelOption {
  id: string
  enabled: boolean
  displayName: string
}

export interface ChatSettings {
  // API 配置
  apiKey: string
  apiUrl: string
  customApiUrl: string
  model: string
  // 社交内容 API（可选独立配置）
  socialApiKey: string
  socialApiUrl: string
  socialModel: string
  enableYamlParsing: boolean
  // 对话参数
  temperature: number
  maxLength: number
  contextSize: number       // 发送给 AI 的最大历史消息条数
  streamEnabled: boolean
  enableSplit: boolean
  timeout: number
  // 心跳（主动消息）
  heartbeatEnabled: boolean
  heartbeatInterval: number
  // 通知
  notifyEnabled: boolean
  notificationSoundEnabled: boolean
  notificationSound: string
  // 外观
  darkMode: boolean
  themeColor: string
  wallpaper: string
  wallpaperImage: string
  fontSize: number
  customFontFamily: string
  // 聊天
  autoSave: boolean
  bubbleStyle: number
  sendKeyBehavior: 'send' | 'newline'
  // 用户资料
  userName: string
  userAvatar: string
  // 时间格式
  timeFormat: '12h' | '24h'
}

const DEFAULT_SETTINGS: ChatSettings = {
  apiKey: '',
  apiUrl: 'https://api.openai.com/v1/chat/completions',
  customApiUrl: '',
  model: 'gpt-4o-mini',
  socialApiKey: '',
  socialApiUrl: '',
  socialModel: '',
  enableYamlParsing: false,
  temperature: 0.9,
  maxLength: 4000,
  contextSize: 20,
  streamEnabled: true,
  enableSplit: true,
  timeout: 60,
  heartbeatEnabled: false,
  heartbeatInterval: 5,
  notifyEnabled: true,
  notificationSoundEnabled: true,
  notificationSound: 'default',
  darkMode: false,
  themeColor: '#007aff',
  wallpaper: 'default',
  wallpaperImage: '',
  fontSize: 2,
  customFontFamily: '',
  autoSave: true,
  bubbleStyle: 0,
  sendKeyBehavior: 'send',
  userName: '',
  userAvatar: '',
  timeFormat: '24h',
}

const STORAGE_KEY = 'fzsm-settings'

export const useSettingsStore = defineStore('settings', () => {
  // 从 localStorage 加载
  function loadFromStorage(): ChatSettings {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        return { ...DEFAULT_SETTINGS, ...parsed }
      }
    } catch (e) {
      console.error('加载设置失败:', e)
    }
    return { ...DEFAULT_SETTINGS }
  }

  const settings = ref<ChatSettings>(loadFromStorage())
  const loading = ref(false)
  const availableModels = ref<ManagedModelOption[]>([])
  const defaultModel = ref('')

  // 自动保存到 localStorage
  watch(settings, (val) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(val))
    } catch (e) {
      console.error('保存设置失败:', e)
    }
  }, { deep: true })

  function updateSettings(partial: Partial<ChatSettings>) {
    Object.assign(settings.value, partial)
  }

  function resetSettings() {
    settings.value = { ...DEFAULT_SETTINGS }
  }

  function getSetting<K extends keyof ChatSettings>(key: K): ChatSettings[K] {
    return settings.value[key]
  }

  // 获取实际的 API URL
  function getApiUrl(): string {
    if (settings.value.apiUrl === 'custom') {
      return settings.value.customApiUrl
    }
    return settings.value.apiUrl
  }

  function normalizeModelId(value: unknown): string {
    return typeof value === 'string' ? value.trim() : ''
  }

  function normalizeDisplayName(value: unknown): string {
    return typeof value === 'string' ? value.trim() : ''
  }

  function parseManagedModels(value: unknown): ManagedModelOption[] {
    if (Array.isArray(value)) {
      return value
        .map((item) => {
          if (typeof item === 'string') {
            const id = normalizeModelId(item)
            if (!id) return null
            return { id, enabled: true, displayName: '' }
          }
          if (!item || typeof item !== 'object') return null
          const model = item as Record<string, unknown>
          const id = normalizeModelId(model.id)
          if (!id) return null
          return {
            id,
            enabled: model.enabled !== false,
            displayName: normalizeDisplayName(model.display_name ?? model.displayName),
          }
        })
        .filter((item): item is ManagedModelOption => Boolean(item))
    }

    if (typeof value !== 'string') return []
    return value
      .split(',')
      .map((item) => normalizeModelId(item))
      .filter(Boolean)
      .map((id) => ({ id, enabled: true, displayName: '' }))
  }

  function getModelDisplayName(modelId: string): string {
    const normalized = normalizeModelId(modelId)
    if (!normalized) return ''
    const model = availableModels.value.find((item) => item.id.toLowerCase() === normalized.toLowerCase())
    return model?.displayName || normalized
  }

  async function fetchAvailableModels() {
    try {
      const res = await api.get<Record<string, any>>('/api/settings/api')
      availableModels.value = parseManagedModels(res.available_models || [])
      defaultModel.value = normalizeModelId(res.default_model)
      return availableModels.value
    } catch (e) {
      console.error('加载模型配置失败:', e)
      return availableModels.value
    }
  }

  // 拉取模型列表（通过后端代理避免 CORS）
  async function fetchModels(): Promise<string[]> {
    const apiKey = settings.value.apiKey
    if (!apiKey) throw new Error('请先输入 API Key')

    const apiUrl = getApiUrl()
    if (!apiUrl) throw new Error('请先配置 API 地址')

    const API_BASE = import.meta.env.VITE_API_URL || ''
    const token = localStorage.getItem('token') || ''

    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 15000)

    try {
      const response = await fetch(`${API_BASE}/api/ai/models`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          apiUrl: apiUrl,
          apiKey: apiKey,
        }),
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('鉴权失败，请手动输入模型名')
        }
        throw new Error(`HTTP ${response.status}`)
      }

      const data = await response.json()
      let modelList: string[] = []

      if (Array.isArray(data?.data)) {
        modelList = data.data.map((item: any) => item.id || item.name || item)
      } else if (Array.isArray(data?.models)) {
        modelList = data.models.map((item: any) => item.id || item.name || item)
      } else if (Array.isArray(data)) {
        modelList = data
      }

      return modelList
        .filter((v: any) => typeof v === 'string')
        .sort()
    } catch (err: any) {
      clearTimeout(timeoutId)
      if (err.name === 'AbortError') {
        throw new Error('请求超时，请检查网络')
      }
      throw err
    }
  }

  return {
    settings,
    loading,
    availableModels,
    defaultModel,
    updateSettings,
    resetSettings,
    getSetting,
    getApiUrl,
    getModelDisplayName,
    fetchAvailableModels,
    fetchModels,
  }
})
