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

  const legacy = localStorage.getItem(baseKey)
  if (legacy !== null) {
    localStorage.setItem(scopedKey, legacy)
  }
  return legacy
}

export function setScopedItem(baseKey: string, value: string) {
  localStorage.setItem(getScopedKey(baseKey), value)
}

export function removeScopedItem(baseKey: string) {
  localStorage.removeItem(getScopedKey(baseKey))
}
