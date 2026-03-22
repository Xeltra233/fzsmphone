<template>
  <div class="phone-peek-page">
    <NavBar title="偷看TA手机" />

    <div class="peek-content">
      <!-- 状态显示 -->
      <div class="peek-status-card">
        <div class="peek-icon" :class="{ peeking: isPeeking }">
          {{ isPeeking ? '◉' : '◎' }}
        </div>
        <h2>{{ isPeeking ? '正在偷看...' : '偷看TA的手机' }}</h2>
        <p>模拟查看对方手机的各项信息</p>
      </div>

      <!-- 偷看选项 -->
      <div class="peek-options">
        <div
          v-for="option in peekOptions"
          :key="option.label"
          class="peek-option"
          @click="peekAt(option)"
        >
          <span class="option-icon">{{ option.icon }}</span>
          <div class="option-info">
            <span class="option-label">{{ option.label }}</span>
            <span class="option-desc">{{ option.desc }}</span>
          </div>
          <span class="option-arrow">›</span>
        </div>
      </div>

      <!-- 偷看结果 -->
      <div v-if="peekResult" class="peek-result">
        <div class="result-header">
          <span>{{ peekResult.icon }} {{ peekResult.title }}</span>
          <button class="close-result" @click="peekResult = null">✕</button>
        </div>
        <div class="result-body">
          <div v-for="(item, idx) in peekResult.items" :key="idx" class="result-item">
            <span class="result-label">{{ item.label }}</span>
            <span class="result-value">{{ item.value }}</span>
          </div>
        </div>
        <div class="result-hint">
          <span>△</span>
          <span>以上内容纯属角色扮演虚构</span>
        </div>
      </div>

      <!-- 发现日志 -->
      <div v-if="discoveries.length > 0" class="discoveries-section">
        <h3>发现记录</h3>
        <div class="discovery-list">
          <div v-for="d in discoveries" :key="d.id" class="discovery-item">
            <span class="disc-icon">{{ d.icon }}</span>
            <div class="disc-info">
              <span class="disc-title">{{ d.title }}</span>
              <span class="disc-time">{{ d.time }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSettingsStore } from '@/stores/settings'
import { useCharactersStore } from '@/stores/characters'
import { sendAIRequest, getCharacterById } from '@/utils/aiService'

const settingsStore = useSettingsStore()
const charactersStore = useCharactersStore()
const isPeeking = ref(false)

interface PeekOption {
  icon: string
  label: string
  desc: string
  category: string
}

interface PeekResultItem {
  label: string
  value: string
}

interface PeekResultData {
  icon: string
  title: string
  items: PeekResultItem[]
}

interface Discovery {
  id: number
  icon: string
  title: string
  time: string
}

const peekOptions: PeekOption[] = [
  { icon: '◌', label: '聊天记录', desc: '查看最近聊天', category: 'chat' },
  { icon: '▣', label: '相册', desc: '查看最近照片', category: 'photos' },
  { icon: '▢', label: '通话记录', desc: '查看最近通话', category: 'calls' },
  { icon: '◎', label: '浏览记录', desc: '查看浏览历史', category: 'browser' },
  { icon: '◆', label: '位置信息', desc: '查看去过的地方', category: 'location' },
  { icon: '☐', label: '购物记录', desc: '查看购买历史', category: 'shopping' },
]

const peekResult = ref<PeekResultData | null>(null)
const discoveries = ref<Discovery[]>([])

const mockResults: Record<string, PeekResultData> = {
  chat: {
    icon: '◌', title: '聊天记录',
    items: [
      { label: '最近联系人', value: '闺蜜、同事小王、妈妈' },
      { label: '最后聊天', value: '5分钟前跟闺蜜聊天' },
      { label: '聊天内容', value: '"今晚一起吃火锅吗？"' },
      { label: '未读消息', value: '3条未读' },
    ],
  },
  photos: {
    icon: '▣', title: '相册',
    items: [
      { label: '最新照片', value: '今天拍的美食照片' },
      { label: '相册数量', value: '2,847张照片' },
      { label: '最近删除', value: '有3张最近删除的照片' },
      { label: '自拍数量', value: '568张自拍' },
    ],
  },
  calls: {
    icon: '▢', title: '通话记录',
    items: [
      { label: '最近通话', value: '妈妈 - 12分钟前' },
      { label: '通话时长', value: '8分32秒' },
      { label: '未接来电', value: '外卖小哥 x2' },
      { label: '今日通话', value: '5次' },
    ],
  },
  browser: {
    icon: '◎', title: '浏览记录',
    items: [
      { label: '最近搜索', value: '"周末去哪玩"' },
      { label: '常用网站', value: '淘宝、微博、知乎' },
      { label: '最近浏览', value: '旅游攻略 - 3小时前' },
      { label: '书签数量', value: '42个收藏' },
    ],
  },
  location: {
    icon: '◆', title: '位置信息',
    items: [
      { label: '当前位置', value: '家（Wi-Fi连接中）' },
      { label: '今日足迹', value: '公司→超市→家' },
      { label: '常去地点', value: '公司、商场、健身房' },
      { label: '上次外出', value: '2小时前到家' },
    ],
  },
  shopping: {
    icon: '☐', title: '购物记录',
    items: [
      { label: '最近购买', value: '一件连衣裙 ¥299' },
      { label: '待收货', value: '2个包裹在路上' },
      { label: '收藏夹', value: '15件商品' },
      { label: '本月消费', value: '¥2,456' },
    ],
  },
}

const DISCOVERY_KEY = 'phone-peek-discoveries'

async function peekAt(option: PeekOption) {
  isPeeking.value = true

  // Try AI generation
  const s = settingsStore.settings
  let aiResult: PeekResultData | null = null

  if (settingsStore.hasChatProviderAccess()) {
    try {
      const currentId = localStorage.getItem('currentPersonaId') || ''
      const character = currentId ? getCharacterById(currentId) : null
      const charName = character?.name || 'TA'

      const prompt = `You are simulating a phone peek feature. The user is peeking at ${charName}'s phone, specifically their ${option.label}.
Generate 4 realistic, fun items for this category. Return ONLY a JSON array like:
[{"label":"...","value":"..."},...]
Keep values short (under 20 chars). Use Chinese. Category: ${option.category}`

      const response = await sendAIRequest({
        apiKey: s.apiSource === 'platform' ? '' : s.apiKey,
        apiUrl: settingsStore.getApiUrl(),
        providerId: settingsStore.getPlatformProviderId(),
        model: s.model,
        messages: [{ role: 'user', content: prompt }],
        temperature: 0.9,
        maxTokens: 300,
        stream: false,
        timeout: s.timeout,
      })

      if (response.content) {
        const jsonMatch = response.content.match(/\[.*\]/s)
        if (jsonMatch) {
          const items = JSON.parse(jsonMatch[0]) as PeekResultItem[]
          if (items.length > 0) {
            aiResult = { icon: option.icon, title: option.label, items }
          }
        }
      }
    } catch { /* fallback to mock */ }
  }

  setTimeout(() => {
    isPeeking.value = false
    peekResult.value = aiResult || mockResults[option.category] || null
    const now = new Date()
    const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
    discoveries.value.unshift({
      id: Date.now(),
      icon: option.icon,
      title: `查看了${option.label}`,
      time: timeStr,
    })
    if (discoveries.value.length > 20) discoveries.value = discoveries.value.slice(0, 20)
    saveDiscoveries()
  }, aiResult ? 300 : 1500)
}

function saveDiscoveries() {
  try {
    localStorage.setItem(DISCOVERY_KEY, JSON.stringify(discoveries.value))
  } catch { /* ignore */ }
}

onMounted(() => {
  charactersStore.fetchCharacters()
  try {
    const saved = localStorage.getItem(DISCOVERY_KEY)
    if (saved) discoveries.value = JSON.parse(saved)
  } catch { /* ignore */ }
})
</script>

<style scoped>
.phone-peek-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.peek-content {
  padding: 16px;
}

.peek-status-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  background: var(--bg-secondary);
  border-radius: 20px;
  margin-bottom: 16px;
}

.peek-icon {
  font-size: 48px;
  transition: transform 0.3s;
}

.peek-icon.peeking {
  animation: peekAnim 0.5s ease-in-out infinite;
}

@keyframes peekAnim {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

.peek-status-card h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.peek-status-card p {
  margin: 0;
  font-size: 13px;
  color: var(--text-tertiary);
}

.peek-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.peek-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--bg-secondary);
  border-radius: 14px;
  cursor: pointer;
  transition: transform 0.15s;
}

.peek-option:active {
  transform: scale(0.98);
}

.option-icon {
  font-size: 28px;
}

.option-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.option-label {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.option-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

.option-arrow {
  font-size: 20px;
  color: var(--text-tertiary);
}

/* 偷看结果 */
.peek-result {
  background: var(--bg-secondary);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 16px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-color);
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-result {
  background: none;
  border: none;
  font-size: 16px;
  color: var(--text-tertiary);
  cursor: pointer;
}

.result-body {
  padding: 12px 16px;
}

.result-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.result-item:last-child {
  border-bottom: none;
}

.result-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.result-value {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: 500;
  text-align: right;
  max-width: 60%;
}

.result-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: rgba(255, 149, 0, 0.08);
  font-size: 12px;
  color: var(--accent-orange);
}

/* 发现记录 */
.discoveries-section {
  margin-top: 8px;
}

.discoveries-section h3 {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 10px;
}

.discovery-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.discovery-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.disc-icon {
  font-size: 20px;
}

.disc-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.disc-title {
  font-size: 13px;
  color: var(--text-primary);
}

.disc-time {
  font-size: 11px;
  color: var(--text-tertiary);
}
</style>
