import { createApp, ref, computed, onMounted, onUnmounted } from 'vue'

const toastState = ref({
  show: false,
  message: '',
  type: 'info' as 'success' | 'error' | 'warning' | 'info',
  duration: 2000,
})

const toastApp = createApp({
  setup() {
    const state = toastState
    const visible = computed(() => state.value.show)
    const message = computed(() => state.value.message)
    const toastType = computed(() => state.value.type)
    const duration = computed(() => state.value.duration)

    let timer: ReturnType<typeof setTimeout> | null = null

    function hide() {
      state.value.show = false
    }

    function display(msg: string, t: typeof state.value.type = 'info', dur = 2000) {
      if (timer) clearTimeout(timer)
      state.value.message = msg
      state.value.type = t
      state.value.duration = dur
      state.value.show = true
      timer = setTimeout(() => {
        state.value.show = false
      }, dur)
    }

    onMounted(() => {
      window.showToast = display
      window.hideToast = hide
    })

    onUnmounted(() => {
      delete window.showToast
      delete window.hideToast
    })

    return { visible, message, toastType, duration }
  },
  template: `
    <Teleport to="body">
      <Transition name="toast">
        <div v-if="visible" class="global-toast" :class="toastType">
          <span class="toast-icon">
            <template v-if="toastType === 'success'">&#10003;</template>
            <template v-else-if="toastType === 'error'">&#10007;</template>
            <template v-else-if="toastType === 'warning'">&#9888;</template>
            <template v-else>&#8505;</template>
          </span>
          <span class="toast-message">{{ message }}</span>
        </div>
      </Transition>
    </Teleport>
  `,
})

const toastContainer = document.createElement('div')
toastContainer.id = 'toast-container'
document.body.appendChild(toastContainer)
toastApp.mount(toastContainer)

interface ToastFunction {
  (message: string, type?: 'success' | 'error' | 'warning' | 'info', duration?: number): void
}

declare global {
  interface Window {
    showToast?: ToastFunction
    hideToast?: () => void
  }
}

export function showToast(message: string, type: 'success' | 'error' | 'warning' | 'info' = 'info', duration = 2000) {
  if (window.showToast) {
    window.showToast(message, type, duration)
  }
}

export function hideToast() {
  if (window.hideToast) {
    window.hideToast()
  }
}

export function showSuccess(message: string, duration?: number) {
  showToast(message, 'success', duration)
}

export function showError(message: string, duration?: number) {
  showToast(message, 'error', duration)
}

export function showWarning(message: string, duration?: number) {
  showToast(message, 'warning', duration)
}

export function showInfo(message: string, duration?: number) {
  showToast(message, 'info', duration)
}