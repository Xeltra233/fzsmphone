<template>
  <div class="status-bar">
    <div class="status-left">
      <span class="time">{{ phone.currentTime }}</span>
    </div>
    <div class="status-right">
      <svg v-if="phone.wifiEnabled" class="status-icon" viewBox="0 0 24 24" fill="currentColor">
        <path d="M1.29 7.05c5.95-5.43 14.47-5.43 20.42 0l-1.41 1.41c-4.95-4.52-12.65-4.52-17.6 0L1.29 7.05zM5.12 10.88c3.75-3.43 9.01-3.43 12.76 0l-1.41 1.41c-2.76-2.52-6.18-2.52-8.94 0l-1.41-1.41zM8.95 14.71c1.56-1.43 3.54-1.43 5.1 0l-1.41 1.41c-.56-.51-1.72-.51-2.28 0l-1.41-1.41zM12 18l-1.41-1.41c.78-.72 2.04-.72 2.82 0L12 18z"/>
      </svg>
      <svg v-if="phone.airplaneMode" class="status-icon" viewBox="0 0 24 24" fill="currentColor">
        <path d="M21 16v-2l-8-5V3.5a1.5 1.5 0 00-3 0V9l-8 5v2l8-2.5V19l-2 1.5V22l3.5-1 3.5 1v-1.5L13 19v-5.5l8 2.5z"/>
      </svg>
      <div class="battery-wrap">
        <div class="battery-body">
          <div
            class="battery-fill"
            :style="{
              width: `${phone.batteryLevel}%`,
              background: phone.batteryLevel <= 20 ? 'var(--color-red)' : 'var(--color-green)',
            }"
          ></div>
        </div>
        <div class="battery-cap"></div>
        <svg v-if="phone.isCharging" class="charging-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M11 21h-1l1-7H7.5c-.88 0-.33-.75-.31-.78C8.48 10.94 10.42 7.54 13.01 3h1l-1 7h3.51c.4 0 .62.19.4.66C12.97 17.55 11 21 11 21z"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePhoneStore } from '@/stores/phone'

const phone = usePhoneStore()
</script>

<style scoped>
.status-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 54px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 0 28px 8px;
  z-index: 100;
  pointer-events: none;
  color: var(--status-bar-foreground, var(--text-primary));
}

.status-left .time {
  font-size: 16px;
  font-weight: 600;
  color: currentColor;
  letter-spacing: -0.3px;
}

.status-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-icon {
  width: 16px;
  height: 16px;
  color: currentColor;
}

/* 电池 */
.battery-wrap {
  display: flex;
  align-items: center;
  gap: 1px;
  position: relative;
}

.battery-body {
  width: 25px;
  height: 12px;
  border: 1.5px solid currentColor;
  border-radius: 3px;
  padding: 1px;
  display: flex;
  align-items: stretch;
}

.battery-fill {
  border-radius: 1px;
  transition: width 0.5s ease, background 0.3s ease;
}

.battery-cap {
  width: 2px;
  height: 5px;
  background: currentColor;
  border-radius: 0 1px 1px 0;
  opacity: 0.4;
}

.charging-icon {
  position: absolute;
  width: 10px;
  height: 10px;
  left: 50%;
  top: 50%;
  transform: translate(-60%, -50%);
  color: currentColor;
}
</style>
