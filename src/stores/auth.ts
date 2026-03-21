import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/client'

export interface User {
  id: string
  discordId: string
  username: string
  displayName: string
  email: string
  avatar: string
  role: string
  isSuperAdmin: boolean
  approved: boolean
  banned: boolean
}

interface UserResponse {
  id: number
  discord_id: string
  username: string
  email: string
  display_name: string
  avatar_url: string
  role: string
  is_super_admin: boolean
  is_banned: boolean
  ban_reason: string
}

interface AuthPayload {
  token: string
  user: UserResponse
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin' || user.value?.role === 'super_admin' || user.value?.isSuperAdmin === true)
  const isSuperAdmin = computed(() => user.value?.role === 'super_admin' || user.value?.isSuperAdmin === true)

  async function fetchUser() {
    if (!token.value) return
    loading.value = true
    try {
      const data = await api.get<UserResponse>('/api/auth/me')
      user.value = {
        id: String(data.id),
        discordId: data.discord_id || '',
        username: data.username,
        displayName: data.display_name || '',
        email: data.email || '',
        avatar: data.avatar_url || '',
        role: data.role || 'user',
        isSuperAdmin: data.is_super_admin || false,
        approved: true,
        banned: data.is_banned || false
      }
    } catch {
      logout()
    } finally {
      loading.value = false
    }
  }

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(userData: User) {
    user.value = userData
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function applyAuthPayload(payload: AuthPayload) {
    token.value = payload.token
    localStorage.setItem('token', payload.token)
    user.value = {
      id: String(payload.user.id),
      discordId: payload.user.discord_id || '',
      username: payload.user.username,
      displayName: payload.user.display_name || '',
      email: payload.user.email || '',
      avatar: payload.user.avatar_url || '',
      role: payload.user.role || 'user',
      isSuperAdmin: payload.user.is_super_admin || false,
      approved: true,
      banned: payload.user.is_banned || false,
    }
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    user,
    token,
    loading,
    isLoggedIn,
    isAdmin,
    isSuperAdmin,
    fetchUser,
    setToken,
    setUser,
    applyAuthPayload,
    logout,
  }
})
