<template>
  <div class="callback-page">
    <div class="loading-spinner animate-spin">⏳</div>
    <p>正在登录...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

onMounted(async () => {
  const token = route.query.token as string
  if (token) {
    auth.setToken(token)
    await auth.fetchUser()
    router.replace('/')
  } else {
    router.replace('/login')
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
  gap: 16px;
  background: linear-gradient(160deg, #0f0c29, #302b63, #24243e);
  color: rgba(255, 255, 255, 0.7);
  font-size: 15px;
}
.loading-spinner {
  font-size: 32px;
}
</style>
