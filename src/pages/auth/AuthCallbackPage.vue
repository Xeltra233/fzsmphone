<template>
  <div class="callback-page">
    <div v-if="error" class="callback-error">
      <div class="error-icon">△</div>
      <p class="error-title">登录失败</p>
      <p class="error-message">{{ error }}</p>
      <button class="retry-btn pressable" @click="goBack">返回登录</button>
    </div>
    <div v-else class="callback-loading">
      <div class="loading-spinner animate-spin">⏳</div>
      <p>正在验证身份...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/stores/auth'
import api from '@/api/client'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const error = ref('')

function goBack() {
  router.replace('/login')
}

onMounted(async () => {
  const code = route.query.code as string

  if (!code) {
    error.value = '未收到授权码，请重新登录'
    return
  }

  try {
    // Send the OAuth code to the backend for verification
    const res = await api.post<{ token: string; user: User }>('/api/auth/discord', { code })

    if (res.token && res.user) {
      auth.setToken(res.token)
      auth.setUser(res.user)
      router.replace('/')
    } else {
      error.value = '登录响应异常，请重试'
    }
  } catch (err: any) {
    // Extract error message from backend response
    const msg = err?.data?.error || err?.message || '登录失败，请重试'
    if (msg.includes('不是指定') || msg.includes('服务器的成员')) {
      error.value = '你不是指定 Discord 服务器的成员，无法登录。请先加入服务器后重试。'
    } else if (msg.includes('无法验证')) {
      error.value = '无法验证服务器成员身份，请稍后重试。'
    } else {
      error.value = msg
    }
  }
})
</script>

<style scoped>
.callback-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: linear-gradient(160deg, #0f0c29, #302b63, #24243e);
}

.callback-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 15px;
}

.loading-spinner {
  font-size: 32px;
}

.callback-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 32px;
  text-align: center;
}

.error-icon {
  font-size: 48px;
}

.error-title {
  font-size: 20px;
  font-weight: 700;
  color: #ff6b6b;
}

.error-message {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  max-width: 280px;
  line-height: 1.5;
}

.retry-btn {
  margin-top: 16px;
  padding: 12px 24px;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  backdrop-filter: blur(10px);
  transition: background 0.2s ease;
}

.retry-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
