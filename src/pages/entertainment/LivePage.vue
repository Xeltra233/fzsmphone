<template>
  <div class="live-page">
    <template v-if="!currentRoom">
      <NavBar title="直播" back>
        <template #right>
          <button class="gen-btn" @click="handleGenerate" :disabled="store.generating">
            {{ store.liveLoading ? 'AI生成中...' : 'AI生成' }}
          </button>
        </template>
      </NavBar>

      <!-- 分类标签 -->
      <div class="category-bar">
        <div
          v-for="cat in categories"
          :key="cat.id"
          class="cat-tab"
          :class="{ active: activeCategory === cat.id }"
          @click="activeCategory = cat.id"
        >
          {{ cat.icon }} {{ cat.name }}
        </div>
      </div>

      <!-- 推荐Banner -->
      <div class="featured-section">
        <div class="featured-card" @click="enterRoom(featuredStreamer)">
          <div class="featured-bg" :style="{ background: featuredStreamer.bgGradient }">
            <div class="featured-avatar">{{ featuredStreamer.avatar }}</div>
            <div class="live-badge">● LIVE</div>
          </div>
          <div class="featured-info">
            <div class="featured-name">{{ featuredStreamer.name }}</div>
            <div class="featured-title">{{ featuredStreamer.title }}</div>
            <div class="featured-viewers">◉ {{ formatNumber(featuredStreamer.viewers) }}</div>
          </div>
        </div>
      </div>

      <!-- 直播列表 -->
      <div class="streamers-grid">
        <div
          v-for="s in filteredStreamers"
          :key="s.id"
          class="streamer-card"
          @click="enterRoom(s)"
        >
          <div class="streamer-preview" :style="{ background: s.bgGradient }">
            <div class="streamer-avatar-big">{{ s.avatar }}</div>
            <div class="viewer-count">◉ {{ formatNumber(s.viewers) }}</div>
            <div v-if="s.isLive" class="live-dot"></div>
          </div>
          <div class="streamer-meta">
            <div class="streamer-name">{{ s.name }}</div>
            <div class="streamer-title">{{ s.title }}</div>
            <div class="streamer-tags">
              <span v-for="tag in s.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 直播间 -->
    <template v-else>
      <div class="live-room">
        <!-- 顶栏 -->
        <div class="room-top">
          <div class="room-streamer">
            <div class="room-avatar">{{ currentRoom.avatar }}</div>
            <div class="room-info">
              <div class="room-name">{{ currentRoom.name }}</div>
              <div class="room-viewers">◉ {{ formatNumber(roomViewers) }}</div>
            </div>
            <button
              class="follow-btn"
              :class="{ followed: currentRoom.followed }"
              @click="currentRoom.followed = !currentRoom.followed"
            >
              {{ currentRoom.followed ? '已关注' : '+ 关注' }}
            </button>
          </div>
          <div class="close-room" @click="leaveRoom">✕</div>
        </div>

        <!-- 直播画面 -->
        <div class="room-stage" :style="{ background: currentRoom.bgGradient }">
          <div class="stage-avatar">{{ currentRoom.avatar }}</div>
          <div class="stage-title">{{ currentRoom.title }}</div>
          <div v-if="showGiftAnimation" class="gift-animation">
            <div class="gift-fly">{{ lastGift }} ×{{ lastGiftCount }}</div>
          </div>
        </div>

        <!-- 弹幕/聊天区 -->
        <div class="room-chat" ref="chatRef">
          <div v-for="(msg, i) in chatMessages" :key="i" class="chat-msg" :class="msg.type">
            <template v-if="msg.type === 'system'">
              <span class="system-text">{{ msg.text }}</span>
            </template>
            <template v-else-if="msg.type === 'gift'">
              <span class="gift-user">{{ msg.user }}</span>
              <span class="gift-text">送出了 {{ msg.gift }} ×{{ msg.count }}</span>
            </template>
            <template v-else>
              <span class="chat-user" :style="{ color: msg.color }">{{ msg.user }}</span>
              <span class="chat-text">{{ msg.text }}</span>
            </template>
          </div>
        </div>

        <!-- 底部操作栏 -->
        <div class="room-bottom">
          <div class="chat-input-row">
            <input
              v-model="chatInput"
              placeholder="说点什么..."
              @keyup.enter="sendChat"
            />
            <button class="send-btn" @click="sendChat">发送</button>
          </div>
          <div class="action-bar">
            <div class="gift-list">
              <div
                v-for="gift in gifts"
                :key="gift.emoji"
                class="gift-item"
                :class="{ selected: selectedGift === gift.emoji }"
                @click="selectedGift = gift.emoji"
              >
                <div class="gift-emoji">{{ gift.emoji }}</div>
                <div class="gift-name">{{ gift.name }}</div>
                <div class="gift-price">{{ gift.price }}币</div>
              </div>
            </div>
            <button
              class="send-gift-btn"
              :disabled="!selectedGift"
              @click="sendGift"
            >
              ★ 送出
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'

const store = useSocialAIStore()

interface Streamer {
  id: number
  name: string
  avatar: string
  title: string
  category: string
  viewers: number
  isLive: boolean
  bgGradient: string
  tags: string[]
  followed: boolean
}

interface ChatMessage {
  type: 'chat' | 'system' | 'gift'
  user?: string
  text?: string
  color?: string
  gift?: string
  count?: number
}

const categories = [
  { id: 'all', name: '全部', icon: '▲' },
  { id: 'gaming', name: '游戏', icon: '▶' },
  { id: 'music', name: '音乐', icon: '♪' },
  { id: 'talk', name: '聊天', icon: '◌' },
  { id: 'food', name: '美食', icon: '◈' },
  { id: 'outdoor', name: '户外', icon: '△' },
]

const activeCategory = ref('all')

const streamers = ref<Streamer[]>([
  { id: 1, name: '小鱼儿', avatar: '◎', title: '今天一起打排位！', category: 'gaming', viewers: 52340, isLive: true, bgGradient: 'linear-gradient(135deg, #667eea, #764ba2)', tags: ['王者荣耀', '排位'], followed: false },
  { id: 2, name: '音乐少女', avatar: '♪', title: '深夜弹唱 治愈你的心', category: 'music', viewers: 18900, isLive: true, bgGradient: 'linear-gradient(135deg, #f093fb, #f5576c)', tags: ['弹唱', '吉他'], followed: false },
  { id: 3, name: '美食大作战', avatar: '◈', title: '挑战10人份炸鸡！', category: 'food', viewers: 31200, isLive: true, bgGradient: 'linear-gradient(135deg, #4facfe, #00f2fe)', tags: ['美食', '挑战'], followed: false },
  { id: 4, name: '户外探险家', avatar: '△', title: '徒步穿越无人区 Day3', category: 'outdoor', viewers: 8760, isLive: true, bgGradient: 'linear-gradient(135deg, #43e97b, #38f9d7)', tags: ['户外', '探险'], followed: false },
  { id: 5, name: '甜甜酱', avatar: '◇', title: '和大家聊聊天~', category: 'talk', viewers: 42100, isLive: true, bgGradient: 'linear-gradient(135deg, #fa709a, #fee140)', tags: ['日常', '互动'], followed: true },
  { id: 6, name: '电竞小王子', avatar: '◎', title: 'LOL大师局 冲冲冲', category: 'gaming', viewers: 27800, isLive: true, bgGradient: 'linear-gradient(135deg, #a18cd1, #fbc2eb)', tags: ['英雄联盟', '大师'], followed: false },
  { id: 7, name: '钢琴诗人', avatar: '♪', title: '古典音乐之夜', category: 'music', viewers: 9500, isLive: true, bgGradient: 'linear-gradient(135deg, #ffecd2, #fcb69f)', tags: ['钢琴', '古典'], followed: false },
  { id: 8, name: '旅行日记', avatar: '➤', title: '东京街头漫步', category: 'outdoor', viewers: 15300, isLive: true, bgGradient: 'linear-gradient(135deg, #89f7fe, #66a6ff)', tags: ['旅行', '日本'], followed: false },
])

const featuredStreamer = computed(() => {
  const live = streamers.value.filter(s => s.isLive)
  return live.reduce((a, b) => a.viewers > b.viewers ? a : b, live[0])
})

const filteredStreamers = computed(() => {
  let list = streamers.value
  if (activeCategory.value !== 'all') {
    list = list.filter(s => s.category === activeCategory.value)
  }
  return list.filter(s => s.id !== featuredStreamer.value?.id)
})

// 直播间
const currentRoom = ref<Streamer | null>(null)
const chatMessages = ref<ChatMessage[]>([])
const chatInput = ref('')
const chatRef = ref<HTMLElement | null>(null)
const roomViewers = ref(0)
const selectedGift = ref('')
const showGiftAnimation = ref(false)
const lastGift = ref('')
const lastGiftCount = ref(1)

let autoMsgTimer: ReturnType<typeof setInterval> | null = null
let viewerTimer: ReturnType<typeof setInterval> | null = null

const gifts = [
  { emoji: '✿', name: '玫瑰', price: 1 },
  { emoji: '●', name: '蛋糕', price: 5 },
  { emoji: '◇', name: '钻石', price: 10 },
  { emoji: '➤', name: '火箭', price: 50 },
  { emoji: '♛', name: '皇冠', price: 100 },
  { emoji: '★', name: '城堡', price: 500 },
]

const botNames = ['路人甲', '小明同学', '开心果', '追梦人', '夜猫子', '可爱多', '小太阳', '棉花糖', '大橙子', '樱花落']
const botColors = ['#ff6b6b', '#4ecdc4', '#45b7d1', '#96ceb4', '#ffeaa7', '#dfe6e9', '#fd79a8', '#6c5ce7', '#00b894', '#e17055']
const botChats = [
  '666666', '好厉害！', '哈哈哈', '冲冲冲！', '太棒了',
  '来了来了', '签到打卡', '主播好可爱', '加油！', '♥♥♥',
  '第一次来', '老粉报到', '好好听', '牛啊', '笑死我了',
  '主播晚上好', '什么时候下播', '求翻牌', '太强了吧', '啊啊啊啊',
]

function formatNumber(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

function enterRoom(s: Streamer) {
  currentRoom.value = s
  roomViewers.value = s.viewers
  chatMessages.value = [
    { type: 'system', text: `欢迎来到 ${s.name} 的直播间` },
    { type: 'system', text: '请文明观看，遵守直播间规则' },
  ]

  // 自动弹幕
  autoMsgTimer = setInterval(() => {
    const rand = Math.random()
    if (rand < 0.15) {
      // 系统消息
      const name = botNames[Math.floor(Math.random() * botNames.length)]
      chatMessages.value.push({ type: 'system', text: `${name} 进入了直播间` })
    } else if (rand < 0.25) {
      // 礼物
      const name = botNames[Math.floor(Math.random() * botNames.length)]
      const gift = gifts[Math.floor(Math.random() * gifts.length)]
      const count = [1, 1, 1, 5, 10, 66][Math.floor(Math.random() * 6)]
      chatMessages.value.push({ type: 'gift', user: name, gift: gift.emoji, count })
    } else {
      // 普通弹幕
      const name = botNames[Math.floor(Math.random() * botNames.length)]
      const color = botColors[Math.floor(Math.random() * botColors.length)]
      const text = botChats[Math.floor(Math.random() * botChats.length)]
      chatMessages.value.push({ type: 'chat', user: name, text, color })
    }
    // 限制消息数量
    if (chatMessages.value.length > 100) {
      chatMessages.value = chatMessages.value.slice(-60)
    }
    scrollChat()
  }, 1500 + Math.random() * 2000)

  // 观众数波动
  viewerTimer = setInterval(() => {
    roomViewers.value += Math.floor(Math.random() * 20) - 8
    if (roomViewers.value < 100) roomViewers.value = 100
  }, 3000)
}

function leaveRoom() {
  currentRoom.value = null
  if (autoMsgTimer) { clearInterval(autoMsgTimer); autoMsgTimer = null }
  if (viewerTimer) { clearInterval(viewerTimer); viewerTimer = null }
}

const chatReactions = [
  '说得好！', '哈哈哈', '同意', '+1', '真的吗', '666', '确实', '笑死',
  '有道理', '我也觉得', '不错不错', '厉害了', '赞同', '可以可以', '懂了懂了',
]
const giftThanks = [
  '谢谢大佬！', '感谢老板！', '大佬大气！', '爱你哟~', '太感动了！',
  '谢谢支持！', '大佬牛！', '老板大气！', '么么哒~', '感谢打赏！',
]
const streamerReactions = [
  '谢谢 {user} 的 {gift}！', '哇 {user} 太豪了！', '{user} 是真爱粉！',
  '感谢 {user}！爱你~', '{user} 太给力了！',
]

function triggerBotReactions(context: 'chat' | 'gift', extra?: { user?: string; gift?: string }) {
  const count = 1 + Math.floor(Math.random() * 3)
  for (let i = 0; i < count; i++) {
    const delay = 500 + Math.random() * 2000
    setTimeout(() => {
      const name = botNames[Math.floor(Math.random() * botNames.length)]
      const color = botColors[Math.floor(Math.random() * botColors.length)]
      let text: string
      if (context === 'gift') {
        text = giftThanks[Math.floor(Math.random() * giftThanks.length)]
      } else {
        text = chatReactions[Math.floor(Math.random() * chatReactions.length)]
      }
      chatMessages.value.push({ type: 'chat', user: name, text, color })
      scrollChat()
    }, delay)
  }
  // Streamer reacts to gifts
  if (context === 'gift' && currentRoom.value && extra) {
    setTimeout(() => {
      const tpl = streamerReactions[Math.floor(Math.random() * streamerReactions.length)]
      const text = tpl.replace('{user}', extra.user || '').replace('{gift}', extra.gift || '')
      chatMessages.value.push({ type: 'chat', user: currentRoom.value!.name, text, color: '#ffd700' })
      scrollChat()
    }, 800 + Math.random() * 1500)
  }
}

function sendChat() {
  const text = chatInput.value.trim()
  if (!text) return
  chatMessages.value.push({ type: 'chat', user: '我', text, color: '#007aff' })
  chatInput.value = ''
  scrollChat()
  triggerBotReactions('chat')
}

function sendGift() {
  if (!selectedGift.value) return
  const gift = gifts.find(g => g.emoji === selectedGift.value)
  if (!gift) return

  chatMessages.value.push({ type: 'gift', user: '我', gift: gift.emoji, count: 1 })
  lastGift.value = gift.emoji
  lastGiftCount.value = 1
  showGiftAnimation.value = true
  setTimeout(() => { showGiftAnimation.value = false }, 2000)
  scrollChat()
  triggerBotReactions('gift', { user: '我', gift: gift.name })
}

function scrollChat() {
  nextTick(() => {
    if (chatRef.value) {
      chatRef.value.scrollTop = chatRef.value.scrollHeight
    }
  })
}

async function handleGenerate() {
  await store.generateLiveContent()
  if (store.liveStreamers.length > 0) {
    streamers.value = store.liveStreamers.map((s, i) => ({
      id: i + 100,
      name: s.name,
      avatar: s.avatar,
      title: s.title,
      category: s.category,
      viewers: s.viewers,
      isLive: s.isLive,
      bgGradient: s.bgGradient,
      tags: s.tags,
      followed: false,
    }))
  }
}

onMounted(() => {
  store.loadData('live')
  if (store.liveStreamers.length > 0) {
    streamers.value = store.liveStreamers.map((s, i) => ({
      id: i + 100,
      name: s.name,
      avatar: s.avatar,
      title: s.title,
      category: s.category,
      viewers: s.viewers,
      isLive: s.isLive,
      bgGradient: s.bgGradient,
      tags: s.tags,
      followed: false,
    }))
  }
})

onUnmounted(() => {
  if (autoMsgTimer) clearInterval(autoMsgTimer)
  if (viewerTimer) clearInterval(viewerTimer)
})
</script>

<style scoped>
.live-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.category-bar {
  display: flex;
  gap: 4px;
  padding: 8px 12px;
  overflow-x: auto;
}

.category-bar::-webkit-scrollbar { display: none; }

.cat-tab {
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 13px;
  white-space: nowrap;
  cursor: pointer;
  color: var(--text-secondary);
  background: var(--bg-tertiary);
  transition: all 0.2s;
  flex-shrink: 0;
}

.cat-tab.active {
  background: var(--color-primary);
  color: #fff;
  font-weight: 600;
}

/* Featured */
.featured-section {
  padding: 8px 16px;
}

.featured-card {
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  background: var(--bg-secondary);
}

.featured-bg {
  height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.featured-avatar { font-size: 60px; }

.live-badge {
  position: absolute;
  top: 10px;
  left: 10px;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 10px;
  font-weight: 600;
}

.featured-info {
  padding: 12px;
}

.featured-name { font-size: 16px; font-weight: 700; color: var(--text-primary); }
.featured-title { font-size: 13px; color: var(--text-secondary); margin-top: 2px; }
.featured-viewers { font-size: 12px; color: var(--text-tertiary); margin-top: 4px; }

/* Grid */
.streamers-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  padding: 8px 12px 20px;
}

.streamer-card {
  background: var(--bg-secondary);
  border-radius: 14px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s;
}

.streamer-card:active { transform: scale(0.97); }

.streamer-preview {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.streamer-avatar-big { font-size: 40px; }

.viewer-count {
  position: absolute;
  bottom: 6px;
  right: 6px;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 8px;
}

.live-dot {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 8px;
  height: 8px;
  border-radius: 4px;
  background: #ff3b30;
  animation: pulse-dot 1.5s infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.3); }
}

.streamer-meta {
  padding: 8px 10px;
}

.streamer-name { font-size: 14px; font-weight: 600; color: var(--text-primary); }
.streamer-title {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.streamer-tags {
  display: flex;
  gap: 4px;
  margin-top: 4px;
}

.tag {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 6px;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
}

/* ===== 直播间 ===== */
.live-room {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #000;
}

.room-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 54px 12px 8px;
  background: rgba(0, 0, 0, 0.6);
}

.room-streamer {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.room-avatar { font-size: 28px; }

.room-info {
  flex: 1;
}

.room-name { font-size: 14px; font-weight: 600; color: #fff; }
.room-viewers { font-size: 11px; color: rgba(255, 255, 255, 0.6); }

.follow-btn {
  padding: 4px 12px;
  border-radius: 14px;
  border: none;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  background: #ff3b30;
  color: #fff;
}

.follow-btn.followed {
  background: rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.7);
}

.close-room {
  width: 28px;
  height: 28px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
  margin-left: 8px;
}

.room-stage {
  height: 220px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.stage-avatar { font-size: 70px; }
.stage-title {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  margin-top: 8px;
}

.gift-animation {
  position: absolute;
  left: 20px;
  bottom: 30px;
  animation: gift-fly 2s ease-out forwards;
}

.gift-fly {
  font-size: 32px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}

@keyframes gift-fly {
  0% { opacity: 1; transform: translateY(0) scale(1); }
  50% { opacity: 1; transform: translateY(-30px) scale(1.3); }
  100% { opacity: 0; transform: translateY(-60px) scale(0.8); }
}

.room-chat {
  flex: 1;
  overflow-y: auto;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.4);
}

.chat-msg {
  padding: 3px 0;
  font-size: 13px;
  line-height: 1.4;
}

.system-text {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
}

.chat-user {
  font-weight: 600;
  margin-right: 6px;
}

.chat-text {
  color: rgba(255, 255, 255, 0.9);
}

.gift-user {
  color: #ffd700;
  font-weight: 600;
  margin-right: 6px;
}

.gift-text {
  color: #ffd700;
}

.room-bottom {
  background: rgba(0, 0, 0, 0.7);
  padding: 8px 12px 12px;
}

.chat-input-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.chat-input-row input {
  flex: 1;
  padding: 8px 12px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 13px;
  outline: none;
}

.chat-input-row input::placeholder { color: rgba(255, 255, 255, 0.4); }

.send-btn {
  padding: 8px 16px;
  border-radius: 18px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.action-bar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.gift-list {
  flex: 1;
  display: flex;
  gap: 6px;
  overflow-x: auto;
}

.gift-list::-webkit-scrollbar { display: none; }

.gift-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  padding: 4px 8px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.1);
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
  border: 2px solid transparent;
}

.gift-item.selected {
  border-color: #ffd700;
  background: rgba(255, 215, 0, 0.15);
}

.gift-emoji { font-size: 20px; }
.gift-name { font-size: 9px; color: rgba(255, 255, 255, 0.7); }
.gift-price { font-size: 9px; color: rgba(255, 215, 0, 0.8); }

.send-gift-btn {
  padding: 8px 14px;
  border-radius: 16px;
  border: none;
  background: linear-gradient(135deg, #ffd700, #ff8c00);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  flex-shrink: 0;
}

.send-gift-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.gen-btn {
  padding: 4px 12px;
  border-radius: 14px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.gen-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
