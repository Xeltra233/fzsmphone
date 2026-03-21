import { useAuthStore } from '@/stores/auth'

function getScopedKey(baseKey: string) {
  try {
    const auth = useAuthStore()
    const userId = auth.user?.id || 'guest'
    return `${baseKey}:${userId}`
  } catch {
    return `${baseKey}:guest`
  }
}

export function getScopedItem(baseKey: string): string | null {
  const scopedKey = getScopedKey(baseKey)
  const scoped = localStorage.getItem(scopedKey)
  if (scoped !== null) return scoped

  return localStorage.getItem(baseKey)
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
