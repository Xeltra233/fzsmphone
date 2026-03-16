const isDevMode = import.meta.env.DEV

const SECURITY_LOG_KEY = 'security_logs'
const MAX_LOGS = 100

interface SecurityLogEntry {
  timestamp: number
  event: string
  details?: Record<string, unknown>
  level: 'info' | 'warn' | 'error'
}

function logToStorage(entry: SecurityLogEntry) {
  try {
    const logs = JSON.parse(localStorage.getItem(SECURITY_LOG_KEY) || '[]')
    logs.push(entry)
    if (logs.length > MAX_LOGS) {
      logs.shift()
    }
    localStorage.setItem(SECURITY_LOG_KEY, JSON.stringify(logs))
  } catch {
    // Ignore storage errors
  }
}

export function logSecurityEvent(
  event: string,
  details?: Record<string, unknown>,
  level: 'info' | 'warn' | 'error' = 'info'
) {
  const entry: SecurityLogEntry = {
    timestamp: Date.now(),
    event,
    details,
    level,
  }

  if (isDevMode) {
    switch (level) {
      case 'error':
        console.error('[SECURITY]', event, details || '')
        break
      case 'warn':
        console.warn('[SECURITY]', event, details || '')
        break
      default:
        console.log('[SECURITY]', event, details || '')
    }
  }

  logToStorage(entry)
}

export function getSecurityLogs(): SecurityLogEntry[] {
  try {
    return JSON.parse(localStorage.getItem(SECURITY_LOG_KEY) || '[]')
  } catch {
    return []
  }
}

export function clearSecurityLogs() {
  localStorage.removeItem(SECURITY_LOG_KEY)
}

export function isSecureConnection(): boolean {
  return window.location.protocol === 'https:' || window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
}

export function getClientInfo(): Record<string, string> {
  return {
    userAgent: navigator.userAgent,
    language: navigator.language,
    platform: navigator.platform,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  }
}