import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export interface Character {
  id: number
  name: string
  avatar_url: string
  description: string
  personality: string
  system_prompt: string
  greeting: string
  is_public: boolean
  tags: string[]
}

export interface Conversation {
  id: number
  character_id: number
  title: string
  is_group: boolean
  last_message: string
  last_at: string
  character?: Character
  // Local fields
  unread?: number
}

export interface ChatMessage {
  id: number
  conversation_id: number
  role: 'user' | 'assistant' | 'system'
  content: string
  msg_type: 'text' | 'image' | 'emoji' | 'redpacket' | 'transfer' | 'share'
  extra: Record<string, unknown>
  created_at: string
}

export const useChatStore = defineStore('chat', () => {
  // State
  const conversations = ref<Conversation[]>([])
  const characters = ref<Character[]>([])
  const currentMessages = ref<ChatMessage[]>([])
  const currentConversationId = ref<number | null>(null)
  const loading = ref(false)

  // Getters
  const sortedConversations = computed(() => {
    return [...conversations.value].sort((a, b) => {
      return new Date(b.last_at).getTime() - new Date(a.last_at).getTime()
    })
  })

  const currentConversation = computed(() => {
    return conversations.value.find(c => c.id === currentConversationId.value)
  })

  // Actions
  async function fetchConversations() {
    try {
      const res = await api.get<{ data: Conversation[] }>('/api/conversations')
      conversations.value = res.data || []
    } catch (e) {
      console.error('Failed to fetch conversations:', e)
    }
  }

  async function fetchCharacters() {
    try {
      const res = await api.get<{ data: Character[] }>('/api/characters')
      characters.value = res.data || []
    } catch (e) {
      console.error('Failed to fetch characters:', e)
    }
  }

  async function fetchMessages(conversationId: number) {
    currentConversationId.value = conversationId
    try {
      const res = await api.get<{ data: ChatMessage[] }>(`/api/conversations/${conversationId}/messages`)
      currentMessages.value = res.data || []
    } catch (e) {
      console.error('Failed to fetch messages:', e)
    }
  }

  async function sendMessage(conversationId: number, content: string, msgType = 'text') {
    const tempMsg: ChatMessage = {
      id: Date.now(),
      conversation_id: conversationId,
      role: 'user',
      content,
      msg_type: msgType as ChatMessage['msg_type'],
      extra: {},
      created_at: new Date().toISOString(),
    }
    currentMessages.value.push(tempMsg)

    // Update conversation preview
    const conv = conversations.value.find(c => c.id === conversationId)
    if (conv) {
      conv.last_message = content
      conv.last_at = tempMsg.created_at
    }

    try {
      await api.post(`/api/conversations/${conversationId}/messages`, {
        content,
        msg_type: msgType,
      })
    } catch (e) {
      console.error('Failed to send message:', e)
    }
  }

  async function createConversation(characterId: number): Promise<Conversation | null> {
    try {
      const res = await api.post<Conversation>('/api/conversations', {
        character_id: characterId,
      })
      if (res) {
        conversations.value.unshift(res)
        return res
      }
    } catch (e) {
      console.error('Failed to create conversation:', e)
    }
    return null
  }

  async function deleteConversation(id: number) {
    try {
      await api.delete(`/api/conversations/${id}`)
      conversations.value = conversations.value.filter(c => c.id !== id)
      if (currentConversationId.value === id) {
        currentConversationId.value = null
        currentMessages.value = []
      }
    } catch (e) {
      console.error('Failed to delete conversation:', e)
    }
  }

  async function createCharacter(data: Partial<Character>): Promise<Character | null> {
    try {
      const res = await api.post<Character>('/api/characters', data)
      if (res) {
        characters.value.push(res)
        return res
      }
    } catch (e) {
      console.error('Failed to create character:', e)
    }
    return null
  }

  async function updateCharacter(id: number, data: Partial<Character>) {
    try {
      await api.put(`/api/characters/${id}`, data)
      const idx = characters.value.findIndex(c => c.id === id)
      if (idx !== -1) {
        characters.value[idx] = { ...characters.value[idx], ...data } as Character
      }
    } catch (e) {
      console.error('Failed to update character:', e)
    }
  }

  async function deleteCharacter(id: number) {
    try {
      await api.delete(`/api/characters/${id}`)
      characters.value = characters.value.filter(c => c.id !== id)
    } catch (e) {
      console.error('Failed to delete character:', e)
    }
  }

  function addLocalMessage(msg: ChatMessage) {
    currentMessages.value.push(msg)
  }

  function clearCurrentChat() {
    currentConversationId.value = null
    currentMessages.value = []
  }

  return {
    // State
    conversations,
    characters,
    currentMessages,
    currentConversationId,
    loading,
    // Getters
    sortedConversations,
    currentConversation,
    // Actions
    fetchConversations,
    fetchCharacters,
    fetchMessages,
    sendMessage,
    createConversation,
    deleteConversation,
    createCharacter,
    updateCharacter,
    deleteCharacter,
    addLocalMessage,
    clearCurrentChat,
  }
})
