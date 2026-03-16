interface RateLimitConfig {
  maxRequests: number
  windowMs: number
}

interface RateLimitEntry {
  count: number
  resetTime: number
}

const rateLimitStore = new Map<string, RateLimitEntry>()

function cleanupExpiredEntries() {
  const now = Date.now()
  for (const [key, entry] of rateLimitStore.entries()) {
    if (now > entry.resetTime) {
      rateLimitStore.delete(key)
    }
  }
}

setInterval(cleanupExpiredEntries, 60000)

export function checkRateLimit(key: string, config: RateLimitConfig): { allowed: boolean; remaining: number; resetIn: number } {
  const now = Date.now()
  let entry = rateLimitStore.get(key)

  if (!entry || now > entry.resetTime) {
    entry = {
      count: 0,
      resetTime: now + config.windowMs,
    }
    rateLimitStore.set(key, entry)
  }

  entry.count++
  const remaining = Math.max(0, config.maxRequests - entry.count)
  const resetIn = Math.max(0, entry.resetTime - now)

  return {
    allowed: entry.count <= config.maxRequests,
    remaining,
    resetIn,
  }
}

export function resetRateLimit(key: string) {
  rateLimitStore.delete(key)
}

export function createRateLimiter(config: RateLimitConfig) {
  return {
    check: (key: string) => checkRateLimit(key, config),
    reset: (key: string) => resetRateLimit(key),
  }
}

export const apiRateLimiter = createRateLimiter({
  maxRequests: 60,
  windowMs: 60000,
})

export const loginRateLimiter = createRateLimiter({
  maxRequests: 5,
  windowMs: 60000,
})