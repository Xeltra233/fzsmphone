<template>
  <PhoneLayout>
    <router-view v-slot="{ Component }">
      <transition name="fade">
        <component :is="Component" />
      </transition>
    </router-view>
  </PhoneLayout>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watchEffect } from 'vue'
import PhoneLayout from '@/layouts/PhoneLayout.vue'
import { usePhoneStore } from '@/stores/phone'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'

const phone = usePhoneStore()
const auth = useAuthStore()
const settingsStore = useSettingsStore()

// 壁纸渐变映射
const wallpaperGradients: Record<string, string> = {
  default: 'linear-gradient(160deg, #0f0c29 0%, #1a1a3e 25%, #24243e 50%, #302b63 75%, #0f0c29 100%)',
  sunset: 'linear-gradient(135deg, #fa709a, #fee140)',
  forest: 'linear-gradient(135deg, #43e97b, #38f9d7)',
  ocean: 'linear-gradient(135deg, #4facfe, #00f2fe)',
  night: 'linear-gradient(135deg, #0c0c1d, #1a1a3e)',
  cherry: 'linear-gradient(135deg, #ffecd2, #fcb69f)',
}

// 字号级别映射 (0-4)
const fontSizeMap: Record<number, number> = {
  0: 12,
  1: 13,
  2: 14,
  3: 16,
  4: 18,
}

// 实时应用外观设置到 CSS 变量
watchEffect(() => {
  const s = settingsStore.settings
  const root = document.documentElement

  // 主题色
  if (s.themeColor) {
    root.style.setProperty('--brand-primary', s.themeColor)
    root.style.setProperty('--color-blue', s.themeColor)
    root.style.setProperty('--text-link', s.themeColor)
  }

  // 壁纸
  const gradient = wallpaperGradients[s.wallpaper] || wallpaperGradients.default
  root.style.setProperty('--wallpaper-gradient', gradient)

  // 自定义壁纸图片
  if (s.wallpaperImage) {
    root.style.setProperty('--wallpaper-image', `url(${s.wallpaperImage})`)
  } else {
    root.style.setProperty('--wallpaper-image', 'none')
  }

  // 字号
  const baseFontSize = fontSizeMap[s.fontSize] ?? 14
  root.style.setProperty('--font-size-base', `${baseFontSize}px`)
})

onMounted(() => {
  // 以 settingsStore 的 darkMode 为准，同步 phoneStore 的主题
  const savedDark = settingsStore.settings.darkMode
  phone.theme = savedDark ? 'dark' : 'light'
  phone.initTheme()
  phone.startBatteryDrain()
  phone.startClock()

  if (auth.token) {
    auth.fetchUser()
  }
})

onUnmounted(() => {
  phone.stopBatteryTimer()
  phone.stopClock()
})
</script>
