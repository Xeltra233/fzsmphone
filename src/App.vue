<template>
  <PhoneLayout>
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </PhoneLayout>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import PhoneLayout from '@/layouts/PhoneLayout.vue'
import { usePhoneStore } from '@/stores/phone'
import { useAuthStore } from '@/stores/auth'

const phone = usePhoneStore()
const auth = useAuthStore()

onMounted(() => {
  phone.initTheme()
  phone.startBatteryDrain()

  if (auth.token) {
    auth.fetchUser()
  }
})

onUnmounted(() => {
  phone.stopBatteryTimer()
})
</script>
