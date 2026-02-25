import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const usePhoneStore = defineStore('phone', () => {
  // === 电池 ===
  const batteryLevel = ref(85)
  const isCharging = ref(false)
  let batteryTimer: ReturnType<typeof setInterval> | null = null

  function startBatteryDrain() {
    stopBatteryTimer()
    batteryTimer = setInterval(() => {
      if (!isCharging.value && batteryLevel.value > 0) {
        batteryLevel.value = Math.max(0, batteryLevel.value - 1)
      } else if (isCharging.value && batteryLevel.value < 100) {
        batteryLevel.value = Math.min(100, batteryLevel.value + 1)
      }
    }, isCharging.value ? 2000 : 60000)
  }

  function toggleCharging() {
    isCharging.value = !isCharging.value
    startBatteryDrain()
  }

  function stopBatteryTimer() {
    if (batteryTimer) {
      clearInterval(batteryTimer)
      batteryTimer = null
    }
  }

  // === 控制中心 ===
  const brightness = ref(100)
  const volume = ref(50)
  const wifiEnabled = ref(true)
  const airplaneMode = ref(false)
  const bluetoothEnabled = ref(false)

  const networkEnabled = computed(() => !airplaneMode.value && wifiEnabled.value)

  // === 主题 ===
  const theme = ref<'light' | 'dark'>('dark')

  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  function initTheme() {
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  return {
    // 电池
    batteryLevel,
    isCharging,
    startBatteryDrain,
    toggleCharging,
    stopBatteryTimer,
    // 控制中心
    brightness,
    volume,
    wifiEnabled,
    airplaneMode,
    bluetoothEnabled,
    networkEnabled,
    // 主题
    theme,
    toggleTheme,
    initTheme,
  }
})
