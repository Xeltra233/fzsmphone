<template>
  <div class="login-page">
    <div class="login-bg"></div>
    <div class="login-content">
      <!-- 封禁提示 -->
      <div v-if="isBanned" class="ban-alert">
        <div class="ban-icon">
          <svg viewBox="0 0 36 36" fill="none" width="36" height="36"><circle cx="18" cy="18" r="14" stroke="#ff6b6b" stroke-width="2.5"/><line x1="9" y1="9" x2="27" y2="27" stroke="#ff6b6b" stroke-width="2.5"/></svg>
        </div>
        <div class="ban-title">账号已被封禁</div>
        <div v-if="banReason" class="ban-reason">封禁原因：{{ banReason }}</div>
        <div class="ban-hint">如有疑问请联系管理员</div>
      </div>

      <div class="login-logo">
        <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="14" y="4" width="36" height="56" rx="8" stroke="#fff" stroke-width="2.5"/>
          <rect x="26" y="8" width="12" height="4" rx="2" fill="rgba(255,255,255,0.5)"/>
          <circle cx="32" cy="52" r="2" fill="rgba(255,255,255,0.4)"/>
          <rect x="18" y="16" width="28" height="30" rx="2" fill="rgba(255,255,255,0.08)"/>
        </svg>
      </div>
      <h1 class="login-title">贩子死妈小手机</h1>
      <p class="login-subtitle">{{ isBanned ? '您的账号已被限制访问' : '请登录以继续' }}</p>

      <!-- 登录方式切换 -->
      <div class="login-tabs">
        <button 
          class="tab-btn" 
          :class="{ active: loginMode === 'password' }" 
          @click="loginMode = 'password'"
        >
          账号登录
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: loginMode === 'discord' }" 
          @click="loginMode = 'discord'"
        >
          Discord
        </button>
      </div>

      <!-- 账号密码登录 -->
      <div v-if="loginMode === 'password'" class="login-form">
        <div class="input-group">
          <input 
            v-model="identifier" 
            type="text" 
            placeholder="邮箱或用户名" 
            class="login-input"
          />
        </div>
        <div class="input-group">
          <input 
            v-model="password" 
            type="password" 
            placeholder="密码" 
            class="login-input"
            @keyup.enter="handleLogin"
          />
        </div>
        
        <button class="login-btn pressable" @click="handleLogin" :disabled="loggingIn">
          <template v-if="loggingIn">
            <div class="btn-spinner"></div>
            登录中...
          </template>
          <template v-else>
            登录
          </template>
        </button>

        <div class="register-link">
          还没有账号？
          <span class="link" @click="showRegister = true">立即注册</span>
        </div>
      </div>

      <!-- Discord 登录 -->
      <div v-if="loginMode === 'discord'" class="login-form">
        <button class="discord-btn pressable" @click="loginWithDiscord">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 00-5.487 0 12.64 12.64 0 00-.617-1.25.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.057 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.041-.106 13.107 13.107 0 01-1.872-.892.077.077 0 01-.008-.128 10.2 10.2 0 00.372-.292.074.074 0 01.077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 01.078.01c.12.098.246.198.373.292a.077.077 0 01-.006.127 12.299 12.299 0 01-1.873.892.077.077 0 00-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 00.084.028 19.839 19.839 0 006.002-3.03.077.077 0 00.032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 00-.031-.03z"/>
          </svg>
          <span>使用 Discord 登录</span>
        </button>
        <p class="login-hint">需要加入指定 Discord 服务器才能使用</p>
      </div>

      <!-- 注册弹窗 -->
      <div v-if="showRegister" class="modal-overlay" @click.self="showRegister = false">
        <div class="modal-content">
          <div class="modal-header">
            <h2>用户注册</h2>
            <button class="close-btn" @click="showRegister = false">&times;</button>
          </div>
          <div class="modal-body">
            <div class="input-group">
              <input v-model="registerForm.username" type="text" placeholder="用户名" class="login-input" />
            </div>
            <div class="input-group">
              <input v-model="registerForm.email" type="email" placeholder="邮箱" class="login-input" />
            </div>
            <div class="input-group">
              <input v-model="registerForm.password" type="password" placeholder="密码（至少6位）" class="login-input" @keyup.enter="handleRegister" />
            </div>
            <button class="login-btn pressable" @click="handleRegister" :disabled="registering">
              <template v-if="registering">
                <div class="btn-spinner"></div>
                注册中...
              </template>
              <template v-else>
                注册
              </template>
            </button>
          </div>
        </div>
      </div>

      <!-- 提示信息 -->
      <div class="toast" v-if="toast.show" :class="toast.type">
        {{ toast.message }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client'

const router = useRouter()

const DISCORD_CLIENT_ID = import.meta.env.VITE_DISCORD_CLIENT_ID || ''
const REDIRECT_URI = import.meta.env.VITE_DISCORD_REDIRECT_URI || `${window.location.origin}/auth/callback`

interface AuthResponse {
  token: string
  user: {
    id: number
    discord_id: string
    username: string
    email: string
    display_name: string
    avatar_url: string
    role: string
    is_super_admin: boolean
  }
}

const DISCORD_CLIENT_ID = import.meta.env.VITE_DISCORD_CLIENT_ID || ''
const REDIRECT_URI = import.meta.env.VITE_DISCORD_REDIRECT_URI || `${window.location.origin}/auth/callback`

const isBanned = ref(false)
const banReason = ref('')
const loginMode = ref('password')
const showRegister = ref(false)
const loggingIn = ref(false)
const registering = ref(false)

const identifier = ref('')
const password = ref('')

const registerForm = ref({
  username: '',
  email: '',
  password: ''
})

const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  if (urlParams.get('banned') === '1') {
    isBanned.value = true
    banReason.value = urlParams.get('reason') || ''
  }
})

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

function loginWithDiscord() {
  const params = new URLSearchParams({
    client_id: DISCORD_CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    response_type: 'code',
    scope: 'identify guilds guilds.members.read',
  })
  window.location.href = `https://discord.com/api/oauth2/authorize?${params.toString()}`
}

async function handleLogin() {
  if (!identifier.value || !password.value) {
    showToast('请输入用户名/邮箱和密码', 'error')
    return
  }
  
  loggingIn.value = true
  try {
    const res = await apiClient.post<AuthResponse>('/api/auth/login', {
      identifier: identifier.value,
      password: password.value
    }) as AuthResponse

    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    window.location.href = '/'
  } catch (err: any) {
    showToast(err.response?.data?.error || '登录失败', 'error')
  } finally {
    loggingIn.value = false
  }
}

async function handleRegister() {
  if (!registerForm.value.username || !registerForm.value.email || !registerForm.value.password) {
    showToast('请填写完整信息', 'error')
    return
  }
  
  if (registerForm.value.password.length < 6) {
    showToast('密码长度至少6位', 'error')
    return
  }
  
  registering.value = true
  try {
    const res = await apiClient.post<AuthResponse>('/api/auth/register', registerForm.value) as AuthResponse

    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    showRegister.value = false
    window.location.href = '/'
  } catch (err: any) {
    showToast(err.response?.data?.error || '注册失败', 'error')
  } finally {
    registering.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(160deg, #0f0c29, #302b63, #24243e);
}

.login-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
}

.login-logo {
  width: 80px;
  height: 80px;
  margin-bottom: 8px;
}

.login-logo svg {
  width: 100%;
  height: 100%;
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
}

.login-subtitle {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 16px;
}

/* Tabs */
.login-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.tab-btn {
  padding: 8px 20px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: rgba(91, 110, 245, 0.3);
  border-color: #5B6EF5;
  color: #fff;
}

/* Form */
.login-form {
  width: 100%;
  max-width: 280px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-group {
  width: 100%;
}

.login-input {
  width: 100%;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  color: #fff;
  font-size: 15px;
  outline: none;
  box-sizing: border-box;
}

.login-input:focus {
  border-color: #5B6EF5;
}

.login-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  color: #fff;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.login-btn:hover:not(:disabled) {
  filter: brightness(1.1);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.discord-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px 28px;
  background: #5865F2;
  color: #fff;
  border-radius: 12px;
  font-size: 17px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(88, 101, 242, 0.4);
  border: none;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.discord-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(88, 101, 242, 0.5);
}

.login-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin-top: 16px;
  text-align: center;
}

.register-link {
  text-align: center;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 8px;
}

.link {
  color: #5B6EF5;
  cursor: pointer;
}

.link:hover {
  text-decoration: underline;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal-content {
  background: #1a1a2e;
  border-radius: 16px;
  width: 90%;
  max-width: 320px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h2 {
  margin: 0;
  font-size: 18px;
  color: #fff;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  font-size: 24px;
  cursor: pointer;
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Toast */
.toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  z-index: 1000;
  animation: fadeInUp 0.3s ease;
}

.toast.success {
  background: rgba(52, 199, 89, 0.9);
  color: #fff;
}

.toast.error {
  background: rgba(255, 59, 48, 0.9);
  color: #fff;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Ban alert */
.ban-alert {
  background: rgba(255, 59, 48, 0.12);
  border: 1px solid rgba(255, 59, 48, 0.35);
  border-radius: 16px;
  padding: 20px 28px;
  text-align: center;
  margin-bottom: 16px;
  animation: banShake 0.5s ease;
}

@keyframes banShake {
  0%, 100% { transform: translateX(0); }
  20% { transform: translateX(-6px); }
  40% { transform: translateX(6px); }
  60% { transform: translateX(-4px); }
  80% { transform: translateX(4px); }
}

.ban-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}

.ban-title {
  font-size: 18px;
  font-weight: 700;
  color: #ff6b6b;
  margin-bottom: 6px;
}

.ban-reason {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 4px;
}

.ban-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}
</style>