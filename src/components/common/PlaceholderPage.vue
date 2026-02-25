<template>
  <div class="placeholder-page">
    <div class="placeholder-nav">
      <button class="back-btn pressable" @click="goBack">
        <svg width="10" height="18" viewBox="0 0 10 18" fill="none">
          <path d="M9 1L1 9L9 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>返回</span>
      </button>
      <h1 class="nav-title">{{ title }}</h1>
      <div class="nav-spacer"></div>
    </div>
    <div class="placeholder-body">
      <div class="placeholder-icon">{{ icon }}</div>
      <h2 class="placeholder-title">{{ title }}</h2>
      <p class="placeholder-desc">此功能正在开发中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

defineProps<{
  title: string
  icon: string
}>()

function goBack() {
  router.back()
}
</script>

<style scoped>
.placeholder-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
}

.placeholder-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 16px;
  margin-top: 54px;
  background: var(--bg-primary);
  border-bottom: 0.5px solid var(--separator);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--color-blue);
  font-size: 17px;
  padding: 8px 4px;
}

.nav-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
}

.nav-spacer {
  width: 60px;
}

.placeholder-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
}

.placeholder-icon {
  font-size: 64px;
  margin-bottom: 8px;
}

.placeholder-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
}

.placeholder-desc {
  font-size: 15px;
  color: var(--text-secondary);
}
</style>
