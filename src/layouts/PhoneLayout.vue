<template>
  <div class="phone-container">
    <div class="iphone-frame">
      <!-- 侧按钮装饰 -->
      <div class="side-button side-button-silent"></div>
      <div class="side-button side-button-vol-up"></div>
      <div class="side-button side-button-vol-down"></div>
      <div class="side-button side-button-power"></div>

      <!-- 屏幕区域 -->
      <div
        class="screen"
        :style="{
          filter: phone.brightness < 100 ? `brightness(${phone.brightness / 100})` : undefined,
          '--status-bar-foreground': statusBarForeground,
        }"
      >
        <!-- 状态栏 -->
        <StatusBar />

        <!-- Dynamic Island -->
        <DynamicIsland />

        <!-- 页面内容 -->
        <div class="screen-content">
          <div v-if="renderError" class="error-boundary">
            <div class="error-icon">⚠️</div>
            <div class="error-title">页面加载异常</div>
            <div class="error-message">{{ renderError.message }}</div>
            <button class="error-retry" @click="handleRetry">返回首页</button>
          </div>
          <slot v-else />
        </div>

        <!-- 底部指示器 -->
        <div class="home-indicator-wrap">
          <div class="home-indicator"></div>
        </div>

        <!-- Teleport 目标：所有弹窗/菜单渲染到此处，限制在手机屏幕内 -->
        <div id="phone-overlay"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onErrorCaptured, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StatusBar from '@/components/phone/StatusBar.vue'
import DynamicIsland from '@/components/phone/DynamicIsland.vue'
import { usePhoneStore } from '@/stores/phone'

const phone = usePhoneStore()
const route = useRoute()
const router = useRouter()

// 从路由 meta 读取状态栏颜色，不再硬编码路由名
const statusBarForeground = computed(() => {
  return (route.meta.statusBarColor as string) || 'var(--text-primary)'
})

// 错误边界：捕获子组件渲染错误，防止白屏
const renderError = ref<Error | null>(null)
onErrorCaptured((err: Error) => {
  renderError.value = err
  console.error('[PhoneLayout] 子组件渲染错误:', err)
  return false // 阻止错误继续向上传播
})

function handleRetry() {
  renderError.value = null
  router.push('/')
}
</script>

<style scoped>
.phone-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100vh;
}

.iphone-frame {
  position: relative;
  width: 393px;
  height: 852px;
  border-radius: 55px;
  background: linear-gradient(145deg, #2a2a2a 0%, #1a1a1a 50%, #0d0d0d 100%);
  padding: 3px;
  box-shadow: var(--shadow-phone);
  overflow: visible;
}

/* 侧按钮 */
.side-button {
  position: absolute;
  background: #2a2a2a;
  border-radius: 2px;
}

.side-button-silent {
  left: -3px;
  top: 160px;
  width: 3px;
  height: 28px;
}

.side-button-vol-up {
  left: -3px;
  top: 210px;
  width: 3px;
  height: 50px;
}

.side-button-vol-down {
  left: -3px;
  top: 270px;
  width: 3px;
  height: 50px;
}

.side-button-power {
  right: -3px;
  top: 230px;
  width: 3px;
  height: 70px;
}

/* 屏幕 */
.screen {
  width: 100%;
  height: 100%;
  border-radius: 52px;
  background: var(--bg-secondary);
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: background-color 0.3s ease;
}

.screen-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 底部 Home Indicator */
.home-indicator-wrap {
  position: absolute;
  bottom: 8px;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  pointer-events: none;
  z-index: 50;
}

.home-indicator {
  width: 134px;
  height: 5px;
  background: var(--text-primary);
  border-radius: 100px;
  opacity: 0.3;
}

/* 错误边界 */
.error-boundary {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  text-align: center;
  gap: 12px;
}

.error-icon {
  font-size: 48px;
  line-height: 1;
}

.error-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.error-message {
  font-size: 13px;
  color: var(--text-secondary);
  max-width: 280px;
  word-break: break-word;
}

.error-retry {
  margin-top: 8px;
  padding: 10px 28px;
  background: var(--brand-primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-full);
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.error-retry:active {
  opacity: 0.8;
}

/* 响应式 - 小屏直接全屏 */
@media (max-width: 430px) {
  .phone-container {
    align-items: stretch;
  }

  .iphone-frame {
    width: 100%;
    height: 100vh;
    border-radius: 0;
    padding: 0;
    box-shadow: none;
  }

  .screen {
    border-radius: 0;
  }

  .side-button {
    display: none;
  }
}
</style>
