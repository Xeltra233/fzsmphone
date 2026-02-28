<template>
  <span
    class="gradient-icon"
    :class="{ inline: inline }"
    :style="{
      width: size + 'px',
      height: size + 'px',
      fontSize: (size * 0.55) + 'px',
      borderRadius: shape === 'circle' ? '50%' : (size * 0.2) + 'px',
      background: iconDef.bg,
    }"
  >{{ iconDef.char }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getContentIcon } from '@/utils/contentIcons'

const props = withDefaults(defineProps<{
  name: string
  size?: number
  shape?: 'circle' | 'square'
  inline?: boolean
}>(), {
  size: 24,
  shape: 'circle',
  inline: false,
})

const iconDef = computed(() => getContentIcon(props.name))
</script>

<style scoped>
.gradient-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 700;
  flex-shrink: 0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  user-select: none;
  line-height: 1;
}

.gradient-icon.inline {
  display: inline-flex;
  vertical-align: middle;
}
</style>
