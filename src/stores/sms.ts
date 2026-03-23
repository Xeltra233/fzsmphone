import { defineStore } from 'pinia'
import { getScopedItem, setScopedItem } from '@/utils/userScopedStorage'
import { ref, computed } from 'vue'
import { getCharacterById } from '@/utils/aiService'
import { smsApi } from '@/api/services'
import { useCharactersStore } from '@/stores/characters'

export interface SmsMessage {
  id: string
  from: 'me' | 'other'
  text: string
  time: string
  timestamp: number
}

export interface SmsConversation {
  id: string
  characterId: string
  name: string
  number: string
  avatar?: string
  color: string
  lastMsg: string
  time: string
  unread: number
  messages: SmsMessage[]
  timestamp: number
}

const SMS_STORAGE_KEY = 'fzsm-sms-conversations'

export const useSmsStore = defineStore('sms', () => {
  const conversations = ref<SmsConversation[]>([])
  const currentConversationId = ref<string | null>(null)

  const sortedConversations = computed(() => {
    return [...conversations.value].sort((a, b) => b.timestamp - a.timestamp)
  })

  const currentConversation = computed(() => {
    return conversations.value.find(c => c.id === currentConversationId.value) || null
  })

  // 从 API 或 localStorage 加载
  async function loadConversations() {
    try {
      const res = await smsApi.listThreads()
      if (res.data && res.data.length > 0) {
        // Map API threads to local format
        const threads = res.data
        const convs: SmsConversation[] = []
        for (const t of threads) {
          const msgs = await loadThreadMessages((t as any).id)
          convs.push({
            id: String((t as any).id),
            characterId: (t as any).character_id || '',
            name: (t as any).recipient || '未命名',
            number: (t as any).number || '',
            avatar: '',
            color: '#007aff',
            lastMsg: msgs.length > 0 ? msgs[msgs.length - 1].text.slice(0, 50) : '',
            time: msgs.length > 0 ? msgs[msgs.length - 1].time : '',
            unread: 0,
            messages: msgs,
            timestamp: new Date((t as any).created_at || Date.now()).getTime(),
          })
        }
        conversations.value = convs
        return
      }
    } catch { /* API failed, fallback */ }
    // Fallback: localStorage
    try {
      const saved = getScopedItem(SMS_STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        if (Array.isArray(parsed)) conversations.value = parsed
      }
    } catch { /* ignore */ }

    if (conversations.value.length === 0) {
      initFromCharacters()
    }
  }

  async function loadThreadMessages(threadId: number): Promise<SmsMessage[]> {
    try {
      const res = await smsApi.listMessages(threadId)
      if (res.data) {
        return res.data.map((m: any) => ({
          id: String(m.id),
          from: m.sender === 'user' ? 'me' as const : 'other' as const,
          text: m.content || '',
          time: new Date(m.created_at || Date.now()).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', hour12: false }),
          timestamp: new Date(m.created_at || Date.now()).getTime(),
        }))
      }
    } catch { /* ignore */ }
    return []
  }

  function saveConversations() {
    try {
      const toSave = conversations.value.map(c => ({
        ...c,
        messages: c.messages.slice(-100),
      }))
      setScopedItem(SMS_STORAGE_KEY, JSON.stringify(toSave))
    } catch { /* ignore */ }
  }

  // 从角色列表初始化短信对话
  function initFromCharacters() {
    try {
      const roleChars = useCharactersStore().charItems
      const colors = ['#ff2d55', '#007aff', '#ff9500', '#5856d6', '#34c759', '#af52de', '#ff6348']

      roleChars.forEach((char: any, index: number) => {
        const charId = String(char.id)
        const existing = conversations.value.find(c => c.characterId === charId)
        if (!existing) {
          conversations.value.push({
            id: `sms-${charId}`,
            characterId: charId,
            name: char.name || '未命名',
            number: generatePhoneNumber(charId),
            avatar: char.avatar_url || '',
            color: colors[index % colors.length],
            lastMsg: '',
            time: '',
            unread: 0,
            messages: [],
            timestamp: Date.now() - index * 1000,
          })
        }
      })
      saveConversations()
    } catch {
      // ignore
    }
  }

  // 根据角色ID生成伪手机号
  function generatePhoneNumber(charId: string): string {
    let hash = 0
    for (let i = 0; i < charId.length; i++) {
      hash = ((hash << 5) - hash) + charId.charCodeAt(i)
      hash |= 0
    }
    const num = Math.abs(hash) % 90000000 + 10000000
    return `138${num}`
  }

  // 打开对话
  function openConversation(id: string) {
    currentConversationId.value = id
    const conv = conversations.value.find(c => c.id === id)
    if (conv) {
      conv.unread = 0
      saveConversations()
    }
  }

  function closeConversation() {
    currentConversationId.value = null
  }

  // 发送用户消息
  async function sendMessage(conversationId: string, text: string): Promise<SmsMessage> {
    const conv = conversations.value.find(c => c.id === conversationId)
    if (!conv) throw new Error('对话不存在')

    const now = new Date()
    const msg: SmsMessage = {
      id: `msg-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
      from: 'me',
      text,
      time: formatTime(now),
      timestamp: now.getTime(),
    }
    conv.messages.push(msg)
    conv.lastMsg = text.slice(0, 50)
    conv.time = formatTime(now)
    conv.timestamp = now.getTime()
    saveConversations()

    // Sync to API
    const threadId = parseInt(conversationId)
    if (!isNaN(threadId)) {
      try { await smsApi.createMessage(threadId, { role: 'user', content: text }) } catch { /* ignore */ }
    }
    return msg
  }

  // 添加AI回复消息
  async function addAIReply(conversationId: string, text: string): Promise<SmsMessage> {
    const conv = conversations.value.find(c => c.id === conversationId)
    if (!conv) throw new Error('对话不存在')

    const now = new Date()
    const msg: SmsMessage = {
      id: `msg-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
      from: 'other',
      text,
      time: formatTime(now),
      timestamp: now.getTime(),
    }
    conv.messages.push(msg)
    conv.lastMsg = text.slice(0, 50)
    conv.time = formatTime(now)
    conv.timestamp = now.getTime()

    if (currentConversationId.value !== conversationId) {
      conv.unread = (conv.unread || 0) + 1
    }

    saveConversations()

    // Sync to API
    const threadId = parseInt(conversationId)
    if (!isNaN(threadId)) {
      try { await smsApi.createMessage(threadId, { role: 'character', content: text }) } catch { /* ignore */ }
    }
    return msg
  }

  // 创建新对话（通过角色）
  function createConversation(characterId: string): SmsConversation {
    const existing = conversations.value.find(c => c.characterId === characterId)
    if (existing) return existing

    const char = getCharacterById(characterId)
    const conv: SmsConversation = {
      id: `sms-${characterId}`,
      characterId,
      name: char?.name || '未命名',
      number: generatePhoneNumber(characterId),
      avatar: char?.avatar || '',
      color: '#007aff',
      lastMsg: '',
      time: '',
      unread: 0,
      messages: [],
      timestamp: Date.now(),
    }
    conversations.value.unshift(conv)
    saveConversations()
    return conv
  }

  // 创建自定义对话（非角色关联）
  function createCustomConversation(name: string, number: string): SmsConversation {
    const conv: SmsConversation = {
      id: `sms-custom-${Date.now()}`,
      characterId: '',
      name,
      number,
      color: '#34c759',
      lastMsg: '',
      time: '',
      unread: 0,
      messages: [],
      timestamp: Date.now(),
    }
    conversations.value.unshift(conv)
    saveConversations()
    return conv
  }

  // 删除对话
  async function deleteConversation(id: string) {
    conversations.value = conversations.value.filter(c => c.id !== id)
    if (currentConversationId.value === id) {
      currentConversationId.value = null
    }
    saveConversations()
    const threadId = parseInt(id)
    if (!isNaN(threadId)) {
      try { await smsApi.deleteThread(threadId) } catch { /* ignore */ }
    }
  }

  // 刷新角色信息
  function refreshCharacters() {
    for (const conv of conversations.value) {
      if (conv.characterId) {
        const char = getCharacterById(conv.characterId)
        if (char) {
          conv.name = char.name || conv.name
          conv.avatar = char.avatar || conv.avatar
        }
      }
    }
    saveConversations()
  }

  // 获取对话的AI消息历史（用于构建上下文）
  function getMessageHistory(conversationId: string, limit = 20) {
    const conv = conversations.value.find(c => c.id === conversationId)
    if (!conv) return []
    return conv.messages.slice(-limit)
  }

  // 总未读数
  const totalUnread = computed(() => {
    return conversations.value.reduce((sum, c) => sum + (c.unread || 0), 0)
  })

  function formatTime(date: Date): string {
    const h = date.getHours().toString().padStart(2, '0')
    const m = date.getMinutes().toString().padStart(2, '0')
    return `${h}:${m}`
  }

  return {
    conversations,
    currentConversationId,
    sortedConversations,
    currentConversation,
    totalUnread,
    loadConversations,
    saveConversations,
    openConversation,
    closeConversation,
    sendMessage,
    addAIReply,
    createConversation,
    createCustomConversation,
    deleteConversation,
    refreshCharacters,
    getMessageHistory,
  }
})
