<template>
  <div class="home-page">
    <!-- 壁纸背景 -->
    <div class="wallpaper"></div>

    <!-- 内容区域 -->
    <div class="home-content">
      <!-- 顶部间距（状态栏+灵动岛） -->
      <div class="top-spacer"></div>

      <!-- 日期时间 -->
      <div class="datetime-widget">
        <div class="widget-date">{{ dateStr }}</div>
        <div class="widget-time">{{ timeStr }}</div>
      </div>

      <!-- App 图标网格 -->
      <div class="app-grid">
        <div
          v-for="app in apps"
          :key="app.name"
          class="app-item pressable"
          @click="openApp(app)"
        >
          <div class="app-icon" :style="{ background: app.color }">
            <span class="app-emoji">{{ app.icon }}</span>
            <div v-if="app.badge" class="app-badge">{{ app.badge }}</div>
          </div>
          <span class="app-name">{{ app.name }}</span>
        </div>
      </div>
    </div>

    <!-- Dock 栏 -->
    <div class="dock">
      <div class="dock-bg"></div>
      <div class="dock-apps">
        <div
          v-for="app in dockApps"
          :key="app.name"
          class="dock-item pressable"
          @click="openApp(app)"
        >
          <div class="app-icon" :style="{ background: app.color }">
            <span class="app-emoji">{{ app.icon }}</span>
            <div v-if="app.badge" class="app-badge">{{ app.badge }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface AppItem {
  name: string
  icon: string
  color: string
  route: string
  badge?: number
}

const dateStr = ref('')
const timeStr = ref('')
let timer: ReturnType<typeof setInterval> | null = null

function updateDateTime() {
  const now = new Date()
  const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  dateStr.value = `${weekDays[now.getDay()]}  ${now.getMonth() + 1}月${now.getDate()}日`
  timeStr.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

const apps: AppItem[] = [
  { name: '聊天', icon: '💬', color: 'linear-gradient(135deg, #5B6EF5, #8B5CF6)', route: '/friends' },
  { name: '论坛', icon: '📋', color: 'linear-gradient(135deg, #FF6B35, #F7931E)', route: '/forum' },
  { name: '微博', icon: '🔥', color: 'linear-gradient(135deg, #E6162D, #FF4757)', route: '/weibo' },
  { name: '朋友圈', icon: '🌈', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/qzone' },
  { name: '外卖', icon: '🛵', color: 'linear-gradient(135deg, #0AB4FF, #26D0CE)', route: '/takeaway' },
  { name: '购物', icon: '🛍️', color: 'linear-gradient(135deg, #FF6348, #FF4757)', route: '/shopping' },
  { name: '音乐', icon: '🎵', color: 'linear-gradient(135deg, #FC5C7D, #6A82FB)', route: '/listen-together' },
  { name: '直播', icon: '📺', color: 'linear-gradient(135deg, #A855F7, #EC4899)', route: '/live' },
  { name: '日记', icon: '📖', color: 'linear-gradient(135deg, #FECA57, #FF9FF3)', route: '/diary' },
  { name: '游戏', icon: '🎮', color: 'linear-gradient(135deg, #6C5CE7, #A29BFE)', route: '/games' },
  { name: '钱包', icon: '💰', color: 'linear-gradient(135deg, #00B894, #55EFC4)', route: '/wallet' },
  { name: '影院', icon: '🎬', color: 'linear-gradient(135deg, #2D3436, #636E72)', route: '/mini-theater' },
  { name: '股票', icon: '📈', color: 'linear-gradient(135deg, #00B16A, #1ABC9C)', route: '/stock' },
  { name: '汇率', icon: '💱', color: 'linear-gradient(135deg, #3498DB, #2980B9)', route: '/currency' },
  { name: '情侣', icon: '💕', color: 'linear-gradient(135deg, #FD79A8, #E84393)', route: '/couple-space' },
  { name: '设置', icon: '⚙️', color: 'linear-gradient(135deg, #636E72, #B2BEC3)', route: '/customize' },
  { name: '角色', icon: '👥', color: 'linear-gradient(135deg, #00CEC9, #81ECEC)', route: '/characters' },
  { name: '人设', icon: '🎭', color: 'linear-gradient(135deg, #E17055, #FAB1A0)', route: '/personas' },
  { name: '世界书', icon: '📚', color: 'linear-gradient(135deg, #6C5CE7, #DDA0DD)', route: '/worldbook' },
  { name: '预设', icon: '⚡', color: 'linear-gradient(135deg, #FDCB6E, #F39C12)', route: '/preset' },
]

const dockApps: AppItem[] = [
  { name: '电话', icon: '📞', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/phone' },
  { name: '短信', icon: '✉️', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/sms' },
  { name: '聊天', icon: '💬', color: 'linear-gradient(135deg, #5B6EF5, #8B5CF6)', route: '/friends' },
  { name: '个人', icon: '👤', color: 'linear-gradient(135deg, #636E72, #B2BEC3)', route: '/profile' },
]

function openApp(app: AppItem) {
  router.push(app.route)
}

onMounted(() => {
  updateDateTime()
  timer = setInterval(updateDateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.home-page {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 壁纸 */
.wallpaper {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    160deg,
    #0f0c29 0%,
    #1a1a3e 25%,
    #24243e 50%,
    #302b63 75%,
    #0f0c29 100%
  );
  z-index: 0;
}

.wallpaper::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 20% 30%, rgba(99, 102, 241, 0.15) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(168, 85, 247, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 50% 50%, rgba(236, 72, 153, 0.05) 0%, transparent 70%);
}

.home-content {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 20px;
  overflow-y: auto;
}

.home-content::-webkit-scrollbar {
  display: none;
}

/* 顶部间距 */
.top-spacer {
  height: 60px;
  flex-shrink: 0;
}

/* 日期时间小组件 */
.datetime-widget {
  text-align: center;
  padding: 8px 0 20px;
  flex-shrink: 0;
}

.widget-date {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 0.5px;
}

.widget-time {
  font-size: 64px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -2px;
  line-height: 1;
  margin-top: 2px;
  text-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
}

/* App 网格 */
.app-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px 12px;
  padding: 8px 0 20px;
}

.app-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.app-icon {
  width: 60px;
  height: 60px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
  transition: transform 0.15s ease;
}

.app-icon::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.25) 0%,
    transparent 50%
  );
  pointer-events: none;
}

.app-emoji {
  font-size: 28px;
  z-index: 1;
}

.app-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  min-width: 20px;
  height: 20px;
  background: var(--color-red);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 5px;
  z-index: 2;
}

.app-name {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.9);
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
  max-width: 68px;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Dock */
.dock {
  position: relative;
  z-index: 1;
  padding: 12px 24px 28px;
}

.dock-bg {
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  border-radius: 32px;
  margin: 0 12px;
}

.dock-apps {
  position: relative;
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 8px 16px;
}

.dock-item {
  cursor: pointer;
}

.dock-item .app-icon {
  width: 56px;
  height: 56px;
}
</style>
