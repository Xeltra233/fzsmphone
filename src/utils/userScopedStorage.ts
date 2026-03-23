import { useAuthStore } from '@/stores/auth'

function resolveStoredUserId(): string | null {
  try {
    const raw = localStorage.getItem('user')
    if (!raw) return null
    const parsed = JSON.parse(raw)
    const id = parsed?.id
    return id ? String(id) : null
  } catch {
    return null
  }
}

function getScopedKey(baseKey: string) {
  try {
    const auth = useAuthStore()
    const userId = auth.user?.id || resolveStoredUserId() || 'guest'
    return `${baseKey}:${userId}`
  } catch {
    return `${baseKey}:${resolveStoredUserId() || 'guest'}`
  }
}

export function getScopedItem(baseKey: string): string | null {
  const scopedKey = getScopedKey(baseKey)
  const scoped = localStorage.getItem(scopedKey)
  if (scoped !== null) return scoped

  const baseValue = localStorage.getItem(baseKey)
  if (baseValue !== null) {
    try {
      localStorage.setItem(scopedKey, baseValue)
      localStorage.removeItem(baseKey)
    } catch {
      // ignore migration failure
    }
    return baseValue
  }

  const guestKey = `${baseKey}:guest`
  const guestValue = localStorage.getItem(guestKey)
  if (guestValue !== null) {
    try {
      localStorage.setItem(scopedKey, guestValue)
      if (scopedKey !== guestKey) localStorage.removeItem(guestKey)
    } catch {
      // ignore migration failure
    }
    return guestValue
  }

  return null
}

export function setScopedItem(baseKey: string, value: string) {
  const scopedKey = getScopedKey(baseKey)
  try {
    localStorage.setItem(scopedKey, value)
  } catch {
    localStorage.removeItem(baseKey)
    localStorage.setItem(scopedKey, value)
  }
}

export function removeScopedItem(baseKey: string) {
  localStorage.removeItem(getScopedKey(baseKey))
}
