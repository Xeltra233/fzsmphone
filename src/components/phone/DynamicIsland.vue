<template>
  <div class="dynamic-island" :class="{ expanded: isExpanded }" @click="toggle">
    <div class="island-camera"></div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isExpanded = ref(false)

function toggle() {
  isExpanded.value = !isExpanded.value
  if (isExpanded.value) {
    setTimeout(() => {
      isExpanded.value = false
    }, 3000)
  }
}
</script>

<style scoped>
.dynamic-island {
  position: absolute;
  top: 12px;
  left: 50%;
  transform: translateX(-50%);
  width: 126px;
  height: 37px;
  background: #000;
  border-radius: 22px;
  z-index: 200;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.dynamic-island::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(
    circle at 30% 30%,
    rgba(255, 255, 255, 0.06) 0%,
    transparent 70%
  );
  border-radius: inherit;
}

.dynamic-island.expanded {
  width: 340px;
  height: 80px;
  border-radius: 36px;
}

.island-camera {
  width: 12px;
  height: 12px;
  background: radial-gradient(circle, #1a1a3e 30%, #0a0a1e 70%);
  border-radius: 50%;
  position: absolute;
  right: 28px;
  box-shadow: inset 0 0 2px rgba(255, 255, 255, 0.1);
}
</style>
