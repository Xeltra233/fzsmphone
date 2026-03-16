import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/client'

interface UserSettings {
  credits: number
  total_tokens: number
  streak: number
  can_signin: boolean
  invite_code: string
  invite_count: number
  rewards_claimed: number
  settings: {
    default_credits: number
    signin_daily_credits: number
    signin_streak_bonus: number
    invite_reward_credits: number
    signin_enabled: boolean
    invite_enabled: boolean
  }
  model_display_names: Record<string, string>
}

interface SigninResult {
  credits_earned: number
  streak_bonus: number
  total_credits: number
  streak: number
}

interface InviteCodeResult {
  invite_code: string
  invite_count: number
  rewards_claimed: number
}

interface RedeemResult {
  message: string
  credits_earned: number
}

interface Coupon {
  code: string
  credits: number
  max_uses: number
  current_uses: number
  expires_at: string | null
  created_at: string
}

export const useCreditsStore = defineStore('credits', () => {
  const settings = ref<UserSettings | null>(null)
  const loading = ref(false)

  async function fetchSettings() {
    loading.value = true
    try {
      const data = await api.get<UserSettings>('/api/credits/settings')
      settings.value = data
    } catch (e) {
      console.error('Failed to fetch credits settings:', e)
    } finally {
      loading.value = false
    }
  }

  async function signIn(): Promise<SigninResult | null> {
    try {
      const data = await api.post<SigninResult>('/api/credits/signin', {})
      await fetchSettings()
      return data
    } catch (e: any) {
      throw new Error(e.message || '签到失败')
    }
  }

  async function getInviteCode(): Promise<InviteCodeResult> {
    return await api.get<InviteCodeResult>('/api/credits/invite-code')
  }

  async function useInviteCode(code: string): Promise<{ message: string; credits_earned: number }> {
    return await api.post<{ message: string; credits_earned: number }>('/api/credits/invite-code', { invite_code: code })
  }

  async function redeemCoupon(code: string): Promise<RedeemResult> {
    return await api.post<RedeemResult>('/api/credits/redeem', { code })
  }

  async function consumeCredits(tokens: number, model: string) {
    return await api.post<{ cost: number; credits: number; tokens: number }>('/api/credits/consume', { tokens, model })
  }

  async function updateModelDisplayName(model: string, displayName: string) {
    return await api.put<{ message: string }>('/api/credits/model-display-name', { model, display_name: displayName })
  }

  return {
    settings,
    loading,
    fetchSettings,
    signIn,
    getInviteCode,
    useInviteCode,
    redeemCoupon,
    consumeCredits,
    updateModelDisplayName,
  }
})