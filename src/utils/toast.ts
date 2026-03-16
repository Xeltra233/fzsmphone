import { createApp, ref } from 'vue'

const toastApp = createApp({
  setup() {
    const show = ref(false)
    const message = ref('')
    const type = ref<'success' | 'error'>('success')

    function display(msg: string, t: 'success' | 'error' = 'success') {
      message.value = msg
      type.value = t
      show.value = true
      setTimeout(() => {
        show.value = false
      }, 2000)
    }

    return { show, message, type, display }
  },
  template: `
    <Teleport to="body">
      <div v-if="show" class="global-toast" :class="type">
        {{ message }}
      </div>
    </Teleport>
  `,
})

const toastContainer = document.createElement('div')
document.body.appendChild(toastContainer)
toastApp.mount(toastContainer)

let toastRef: any = null

toastApp.mount(toastContainer)

export function showToast(message: string, type: 'success' | 'error' = 'success') {
  const toastEl = document.querySelector('.global-toast') as HTMLElement
  if (toastEl) {
    toastEl.textContent = message
    toastEl.className = `global-toast ${type}`
    toastEl.style.display = 'block'
    setTimeout(() => {
      toastEl.style.display = 'none'
    }, 2000)
  }
}