import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import api from './api/client'

// Styles
import './styles/variables.css'
import './styles/base.css'
import './styles/animations.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Fetch system settings and apply them
async function applySystemSettings() {
  try {
    const token = localStorage.getItem('token')
    if (!token) return
    
    const res = await api.get<any>('/api/settings')
    const data = res.data || res || {}
    
    // Apply custom title
    if (data.app_title) {
      document.title = data.app_title
    }
    
    // Apply custom favicon
    if (data.favicon) {
      const link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
      if (link) {
        link.href = data.favicon
      }
    }
  } catch (e) {
    // Silently fail - use defaults
  }
}

// Apply settings after router is ready
router.isReady().then(() => {
  applySystemSettings()
})

app.mount('#app')