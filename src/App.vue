<template>
  <PhoneLayout>
    <router-view v-slot="{ Component }">
      <transition :name="transitionName">
        <component :is="Component" />
      </transition>
    </router-view>
  </PhoneLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import PhoneLayout from '@/layouts/PhoneLayout.vue'
import { usePhoneStore } from '@/stores/phone'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'

const phone = usePhoneStore()
const auth = useAuthStore()
const settingsStore = useSettingsStore()
const router = useRouter()

// === iOS 风格转场动画 ===
const transitionName = ref('fade')

// 为路由分配深度级别，用于判断前进/后退
const routeDepthMap: Record<string, number> = {
  Home: 0,
  Login: 0,
  // 一级页面
  Friends: 1,
  Characters: 1,
  UserPersonas: 1,
  WorldBook: 1,
  Weibo: 1,
  QZone: 1,
  Moments: 1,
  CoupleSpace: 1,
  Takeaway: 1,
  Shopping: 1,
  Wallet: 1,
  ListenTogether: 1,
  Live: 1,
  Games: 1,
  Casino: 1,
  MiniTheater: 1,
  Diary: 1,
  Phone: 1,
  Sms: 1,
  Stock: 1,
  CurrencyConverter: 1,
  Profile: 1,
  Customize: 1,
  Preset: 1,
  Zhihu: 1,
  Xiaohongshu: 1,
  Douyin: 1,
  Discord: 1,
  Email: 1,
  Browser: 1,
  Map: 1,
  Calendar: 1,
  AdminFeatures: 1,
  AdminUsers: 1,
  AdminStats: 1,
  AdminSettings: 1,
  // 二级页面
  Chat: 2,
  GroupChat: 2,
  CharacterEdit: 2,
  Restaurant: 2,
  TakeawayOrders: 2,
  VoiceCall: 2,
  VideoCall: 2,
  PhonePeek: 2,
  ReversePhonePeek: 2,
  OfflineDate: 2,
}

const removeRouteTransitionGuard = router.beforeEach((to, from) => {
  const toDepth = routeDepthMap[to.name as string] ?? 1
  const fromDepth = routeDepthMap[from.name as string] ?? 1

  if (from.name === undefined || from.name === null) {
    // 初次加载，使用淡入
    transitionName.value = 'fade'
  } else if (toDepth > fromDepth) {
    // 前进：从右侧推入
    transitionName.value = 'slide-right'
  } else if (toDepth < fromDepth) {
    // 后退：从左侧推入
    transitionName.value = 'slide-left'
  } else {
    // 同级切换：淡入淡出
    transitionName.value = 'fade'
  }
})

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
  removeRouteTransitionGuard()
  phone.stopBatteryTimer()
  phone.stopClock()
})
</script>
