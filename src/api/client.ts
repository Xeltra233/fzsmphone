import router from '@/router'

const BASE_URL = import.meta.env.VITE_API_URL || ''

// ========== 统一错误类型 ==========
export class ApiError extends Error {
  status: number
  data: Record<string, unknown>

  constructor(message: string, status: number, data: Record<string, unknown> = {}) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }

  /** 是否为认证失败 */
  get isUnauthorized(): boolean {
    return this.status === 401
  }

  /** 是否为封禁 */
  get isBanned(): boolean {
    return this.status === 403 && !!this.data.banned
  }
}

interface RequestOptions extends RequestInit {
  params?: Record<string, string>
  /** 请求超时（秒），默认 30 */
  timeout?: number
}

// 401 跳转防抖锁：避免并发请求重复跳转
let redirectingTo401 = false

class ApiClient {
  private baseUrl: string
  private defaultTimeout: number

  constructor(baseUrl: string, defaultTimeout = 30) {
    this.baseUrl = baseUrl
    this.defaultTimeout = defaultTimeout
  }

  private getToken(): string | null {
    return localStorage.getItem('token')
  }

  private buildUrl(path: string, params?: Record<string, string>): string {
    const url = new URL(`${this.baseUrl}${path}`, window.location.origin)
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        url.searchParams.set(key, value)
      })
    }
    return url.toString()
  }

  private async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const { params, timeout, ...fetchOptions } = options
    const url = this.buildUrl(path, params)
    const timeoutMs = (timeout ?? this.defaultTimeout) * 1000

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(fetchOptions.headers as Record<string, string> || {}),
    }

    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    // AbortController 超时机制
    const controller = new AbortController()
    let didTimeout = false
    const existingSignal = fetchOptions.signal
    const onAbort = () => controller.abort()
    if (existingSignal) {
      if (existingSignal.aborted) {
        controller.abort()
      } else {
        existingSignal.addEventListener('abort', onAbort, { once: true })
      }
    }
    const timeoutId = setTimeout(() => {
      didTimeout = true
      controller.abort()
    }, timeoutMs)

    try {
      const response = await fetch(url, {
        ...fetchOptions,
        headers,
        signal: controller.signal,
      })

      clearTimeout(timeoutId)
      if (existingSignal) {
        existingSignal.removeEventListener('abort', onAbort)
      }

      // 错误 body 只读取一次，避免 403 分支与通用错误分支重复消费 body
      let parsedErrorData: Record<string, unknown> | null = null
      const getErrorData = async (): Promise<Record<string, unknown>> => {
        if (parsedErrorData) return parsedErrorData
        const raw = await response.text().catch(() => '')
        if (!raw) {
          parsedErrorData = { error: 'Request failed' }
          return parsedErrorData
        }
        try {
          const parsed = JSON.parse(raw)
          parsedErrorData = (parsed && typeof parsed === 'object')
            ? parsed as Record<string, unknown>
            : { error: String(parsed) }
        } catch {
          parsedErrorData = { error: raw }
        }
        return parsedErrorData
      }

      // 401 未认证 → 通过 router 跳转，保持 SPA 状态
      if (response.status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        if (!redirectingTo401) {
          redirectingTo401 = true
          try {
            await router.replace('/login')
          } finally {
            // 重置锁，延迟一小段避免竞态
            setTimeout(() => { redirectingTo401 = false }, 500)
          }
        }
        throw new ApiError('Unauthorized', 401)
      }

      // 403 封禁
      if (response.status === 403) {
        const data = await getErrorData()
        if (data.banned) {
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          const reason = encodeURIComponent((data.reason as string) || '')
          if (!redirectingTo401) {
            redirectingTo401 = true
            try {
              await router.replace(`/login?banned=1&reason=${reason}`)
            } finally {
              setTimeout(() => { redirectingTo401 = false }, 500)
            }
          }
          throw new ApiError((data.error as string) || 'Account banned', 403, data)
        }
      }

      if (!response.ok) {
        const errorData = await getErrorData()
        throw new ApiError(
          (errorData.error as string) || `HTTP ${response.status}`,
          response.status,
          errorData,
        )
      }

      return response.json()
  } catch (err: any) {
      clearTimeout(timeoutId)
      if (existingSignal) {
        existingSignal.removeEventListener('abort', onAbort)
      }
      if (err instanceof ApiError) throw err
      if (err.name === 'AbortError') {
        if (didTimeout) {
          throw new ApiError(`请求超时（${timeout ?? this.defaultTimeout}秒）`, 0, { timeout: true })
        }
        throw new ApiError('请求已取消', 0, { aborted: true })
      }
      throw err
    }
  }

  async get<T>(path: string, params?: Record<string, string>): Promise<T> {
    return this.request<T>(path, { method: 'GET', params })
  }

  async post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async put<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async patch<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async delete<T>(path: string): Promise<T> {
    return this.request<T>(path, { method: 'DELETE' })
  }
}

export const api = new ApiClient(BASE_URL)
export default api
