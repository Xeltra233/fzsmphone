<template>
  <div class="listen-page">
    <NavBar title="一起听" :show-back="true" back-to="/">
      <template #right>
        <button class="nav-btn" @click="handleGenerate" :disabled="store.generating">
          {{ store.musicLoading ? 'AI生成中...' : 'AI生成' }}
        </button>
        <button class="nav-btn" @click="showSearch = !showSearch">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- 搜索栏 -->
    <div v-if="showSearch" class="search-bar">
      <input v-model="searchQuery" type="text" placeholder="搜索歌曲、歌手..." class="search-input" />
    </div>

    <!-- 当前播放 -->
    <div class="now-playing" v-if="currentTrack">
      <div class="album-art" :class="{ spinning: isPlaying && !getTrackImages(currentTrack).length }">
        <template v-if="getTrackImages(currentTrack).length">
          <div class="img-wrapper">
            <img :src="getTrackImages(currentTrack)[0]" class="album-gen-image" alt="" />
            <button v-if="getStoreTrack(currentTrack)?.imagePrompt" class="regen-btn" :disabled="store.regeneratingImages.has(`${getStoreTrack(currentTrack)?.id}-0`)" @click.stop="store.regenerateImage('music', getStoreTrack(currentTrack)!.id, 0)">
              <span v-if="store.regeneratingImages.has(`${getStoreTrack(currentTrack)?.id}-0`)" class="regen-spin"></span>
              <svg v-else viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10" /><polyline points="1 20 1 14 7 14" /><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" /></svg>
            </button>
          </div>
        </template>
        <template v-else>
          <div class="album-disc">
            <div class="disc-center"></div>
          </div>
        </template>
      </div>
      <div class="track-info">
        <h2 class="track-name">{{ currentTrack.name }}</h2>
        <p class="track-artist">{{ currentTrack.artist }}</p>
      </div>

      <!-- 进度条 -->
      <div class="progress-section">
        <span class="time-label">{{ formatDuration(currentTime) }}</span>
        <div class="progress-bar" @click="seekTo($event)">
          <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
          <div class="progress-thumb" :style="{ left: progressPercent + '%' }"></div>
        </div>
        <span class="time-label">{{ formatDuration(currentTrack.duration) }}</span>
      </div>

      <!-- 控制按钮 -->
      <div class="controls">
        <button class="ctrl-btn" @click="toggleShuffle" :class="{ active: shuffle }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="16 3 21 3 21 8" />
            <line x1="4" y1="20" x2="21" y2="3" />
            <polyline points="21 16 21 21 16 21" />
            <line x1="15" y1="15" x2="21" y2="21" />
            <line x1="4" y1="4" x2="9" y2="9" />
          </svg>
        </button>
        <button class="ctrl-btn" @click="prevTrack">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <polygon points="19 20 9 12 19 4 19 20" />
            <line x1="5" y1="19" x2="5" y2="5" stroke="currentColor" stroke-width="2" />
          </svg>
        </button>
        <button class="play-btn" @click="togglePlay">
          <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3" />
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="currentColor">
            <rect x="6" y="4" width="4" height="16" />
            <rect x="14" y="4" width="4" height="16" />
          </svg>
        </button>
        <button class="ctrl-btn" @click="nextTrack">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 4 15 12 5 20 5 4" />
            <line x1="19" y1="5" x2="19" y2="19" stroke="currentColor" stroke-width="2" />
          </svg>
        </button>
        <button class="ctrl-btn" @click="cycleRepeat" :class="{ active: repeat !== 'off' }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="17 1 21 5 17 9" />
            <path d="M3 11V9a4 4 0 0 1 4-4h14" />
            <polyline points="7 23 3 19 7 15" />
            <path d="M21 13v2a4 4 0 0 1-4 4H3" />
          </svg>
          <span v-if="repeat === 'one'" class="repeat-badge">1</span>
        </button>
      </div>

      <!-- 在线听众 -->
      <div class="listeners-section">
        <div class="listeners-label">♪ {{ listeners.length }} 人正在收听</div>
        <div class="listeners-avatars">
          <div v-for="l in listeners.slice(0, 5)" :key="l" class="listener-dot" :title="l">
            {{ l[0] }}
          </div>
          <div v-if="listeners.length > 5" class="listener-more">+{{ listeners.length - 5 }}</div>
        </div>
      </div>
    </div>

    <!-- 播放列表 -->
    <div class="playlist-section">
      <div class="section-header">
        <span>播放列表</span>
        <span class="track-count">{{ filteredTracks.length }} 首</span>
      </div>
      <div class="track-list">
        <div
          v-for="(track, idx) in filteredTracks"
          :key="track.id"
          class="track-item"
          :class="{ playing: currentTrack?.id === track.id }"
          @click="playTrack(track)"
        >
          <span class="track-index" :class="{ 'is-playing': currentTrack?.id === track.id }">
            <template v-if="currentTrack?.id === track.id && isPlaying">
              <span class="eq-bar"></span>
            </template>
            <template v-else>{{ idx + 1 }}</template>
          </span>
          <div class="track-detail">
            <div class="track-title">{{ track.name }}</div>
            <div class="track-subtitle">{{ track.artist }} · {{ track.album }}</div>
          </div>
          <span class="track-duration">{{ formatDuration(track.duration) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'

const store = useSocialAIStore()

const showSearch = ref(false)
const searchQuery = ref('')
const isPlaying = ref(false)
const currentTime = ref(0)
const shuffle = ref(false)
const repeat = ref<'off' | 'all' | 'one'>('off')
let playTimer: ReturnType<typeof setInterval> | null = null

interface Track {
  id: number
  name: string
  artist: string
  album: string
  duration: number
}

const tracks = ref<Track[]>([
  { id: 1, name: '晴天', artist: '周杰伦', album: '叶惠美', duration: 269 },
  { id: 2, name: '七里香', artist: '周杰伦', album: '七里香', duration: 296 },
  { id: 3, name: '稻香', artist: '周杰伦', album: '魔杰座', duration: 223 },
  { id: 4, name: '花海', artist: '周杰伦', album: '魔杰座', duration: 264 },
  { id: 5, name: '告白气球', artist: '周杰伦', album: '周杰伦的床边故事', duration: 215 },
  { id: 6, name: '起风了', artist: '买辣椒也用券', album: '起风了', duration: 325 },
  { id: 7, name: '海阔天空', artist: 'Beyond', album: '乐与怒', duration: 326 },
  { id: 8, name: '光辉岁月', artist: 'Beyond', album: '命运派对', duration: 282 },
  { id: 9, name: '夜曲', artist: '周杰伦', album: '11月的萧邦', duration: 226 },
  { id: 10, name: '简单爱', artist: '周杰伦', album: '范特西', duration: 270 },
])

const currentTrack = ref<Track | null>(tracks.value[0])

const listeners = ref(['小明', '小红', '小华', '小李', '小张', '小王', '小赵'])

const filteredTracks = computed(() => {
  if (!searchQuery.value) return tracks.value
  const q = searchQuery.value.toLowerCase()
  return tracks.value.filter(t =>
    t.name.toLowerCase().includes(q) ||
    t.artist.toLowerCase().includes(q) ||
    t.album.toLowerCase().includes(q)
  )
})

const progressPercent = computed(() => {
  if (!currentTrack.value) return 0
  return (currentTime.value / currentTrack.value.duration) * 100
})

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

function togglePlay() {
  isPlaying.value = !isPlaying.value
  if (isPlaying.value) startTimer()
  else stopTimer()
}

function startTimer() {
  stopTimer()
  playTimer = setInterval(() => {
    if (!currentTrack.value) return
    currentTime.value += 1
    if (currentTime.value >= currentTrack.value.duration) {
      if (repeat.value === 'one') {
        currentTime.value = 0
      } else {
        nextTrack()
      }
    }
  }, 1000)
}

function stopTimer() {
  if (playTimer) {
    clearInterval(playTimer)
    playTimer = null
  }
}

function playTrack(track: Track) {
  currentTrack.value = track
  currentTime.value = 0
  isPlaying.value = true
  startTimer()
}

function nextTrack() {
  if (!currentTrack.value) return
  const idx = tracks.value.findIndex(t => t.id === currentTrack.value!.id)
  const nextIdx = shuffle.value
    ? Math.floor(Math.random() * tracks.value.length)
    : (idx + 1) % tracks.value.length
  playTrack(tracks.value[nextIdx])
}

function prevTrack() {
  if (currentTime.value > 3) {
    currentTime.value = 0
    return
  }
  if (!currentTrack.value) return
  const idx = tracks.value.findIndex(t => t.id === currentTrack.value!.id)
  const prevIdx = idx <= 0 ? tracks.value.length - 1 : idx - 1
  playTrack(tracks.value[prevIdx])
}

function toggleShuffle() {
  shuffle.value = !shuffle.value
}

function cycleRepeat() {
  const order: ('off' | 'all' | 'one')[] = ['off', 'all', 'one']
  const idx = order.indexOf(repeat.value)
  repeat.value = order[(idx + 1) % 3]
}

function seekTo(e: MouseEvent) {
  if (!currentTrack.value) return
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = (e.clientX - rect.left) / rect.width
  currentTime.value = Math.floor(pct * currentTrack.value.duration)
}

async function handleGenerate() {
  await store.generateMusicContent()
  if (store.musicTracks.length > 0) {
    // Replace track list with AI-generated tracks
    tracks.value = store.musicTracks.map((t, i) => ({
      id: i + 100,
      name: t.name,
      artist: t.artist,
      album: t.album,
      duration: t.duration,
    }))
    if (tracks.value.length > 0) {
      currentTrack.value = tracks.value[0]
      currentTime.value = 0
    }
  }
}

function getTrackImages(t: Track | null): string[] {
  if (!t) return []
  const storeItem = store.musicTracks.find(mt => mt.name === t.name) as any
  return Array.isArray(storeItem?.images) ? storeItem.images.filter((x: string) => !!x) : []
}

function getStoreTrack(t: Track | null): any {
  if (!t) return null
  return store.musicTracks.find(mt => mt.name === t.name) || null
}

onMounted(() => {
  store.loadData('music')
  if (store.musicTracks.length > 0) {
    tracks.value = store.musicTracks.map((t, i) => ({
      id: i + 100,
      name: t.name,
      artist: t.artist,
      album: t.album,
      duration: t.duration,
    }))
    if (tracks.value.length > 0) {
      currentTrack.value = tracks.value[0]
    }
  }
})

onUnmounted(() => {
  stopTimer()
})
</script>

<style scoped>
.listen-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  overflow-y: auto;
}

.listen-page::-webkit-scrollbar {
  display: none;
}

/* 搜索 */
.search-bar {
  padding: 0 16px 8px;
}

.search-input {
  width: 100%;
  padding: 8px 14px;
  border: none;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  font-size: 15px;
  color: var(--text-primary);
  outline: none;
  box-sizing: border-box;
}

.search-input::placeholder {
  color: var(--text-tertiary);
}

/* 当前播放 */
.now-playing {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.album-art {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1a1a2e, #16213e);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  margin-bottom: 20px;
}

.album-art.spinning {
  animation: rotate 8s linear infinite;
}

@keyframes rotate {
  to { transform: rotate(360deg); }
}

.album-disc {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background: repeating-radial-gradient(
    circle,
    #2d2d44 0px,
    #2d2d44 2px,
    #1a1a2e 3px,
    #1a1a2e 8px
  );
  display: flex;
  align-items: center;
  justify-content: center;
}

.disc-center {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
}
.album-gen-image { width: 100%; height: 100%; object-fit: cover; border-radius: 50%; display: block; }
.album-art .img-wrapper { position: relative; width: 100%; height: 100%; border-radius: 50%; overflow: hidden; }
.regen-btn { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); width: 32px; height: 32px; border-radius: 50%; background: rgba(0,0,0,0.5); border: none; color: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center; opacity: 0; transition: opacity 0.2s; z-index: 2; }
.img-wrapper:hover .regen-btn { opacity: 1; }
.regen-btn:disabled { cursor: wait; opacity: 1 !important; }
.regen-spin { width: 12px; height: 12px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: rspin 0.8s linear infinite; }
@keyframes rspin { to { transform: rotate(360deg); } }

.track-info {
  text-align: center;
  margin-bottom: 16px;
}

.track-name {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.track-artist {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

/* 进度条 */
.progress-section {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 0 8px;
  margin-bottom: 20px;
}

.time-label {
  font-size: 12px;
  color: var(--text-tertiary);
  min-width: 36px;
  text-align: center;
}

.progress-bar {
  flex: 1;
  height: 4px;
  background: var(--bg-quaternary);
  border-radius: 2px;
  position: relative;
  cursor: pointer;
}

.progress-fill {
  height: 100%;
  background: var(--brand-primary);
  border-radius: 2px;
  transition: width 0.3s linear;
}

.progress-thumb {
  position: absolute;
  top: 50%;
  width: 12px;
  height: 12px;
  background: var(--brand-primary);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

/* 控制按钮 */
.controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.ctrl-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: none;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.ctrl-btn svg {
  width: 22px;
  height: 22px;
}

.ctrl-btn.active {
  color: var(--brand-primary);
}

.ctrl-btn:active {
  opacity: 0.5;
}

.repeat-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  font-size: 10px;
  font-weight: 700;
  color: var(--brand-primary);
}

.play-btn {
  width: 56px;
  height: 56px;
  border: none;
  background: var(--brand-primary);
  color: #fff;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(91, 110, 245, 0.4);
  transition: transform 0.15s;
}

.play-btn svg {
  width: 24px;
  height: 24px;
}

.play-btn:active {
  transform: scale(0.92);
}

/* 在线听众 */
.listeners-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 12px 16px;
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
}

.listeners-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.listeners-avatars {
  display: flex;
  gap: -4px;
}

.listener-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6C5CE7, #A29BFE);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  margin-left: -6px;
  border: 2px solid var(--bg-primary);
}

.listener-dot:first-child {
  margin-left: 0;
}

.listener-more {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: var(--text-tertiary);
  margin-left: -6px;
  border: 2px solid var(--bg-primary);
}

/* 播放列表 */
.playlist-section {
  padding: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 17px;
  font-weight: 600;
  color: var(--text-primary);
}

.track-count {
  font-size: 13px;
  font-weight: 400;
  color: var(--text-tertiary);
}

.track-list {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.track-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
  -webkit-tap-highlight-color: transparent;
}

.track-item:active {
  background: var(--bg-tertiary);
}

.track-item:not(:last-child) {
  border-bottom: 0.5px solid var(--separator);
}

.track-item.playing {
  background: var(--bg-tertiary);
}

.track-index {
  width: 24px;
  text-align: center;
  font-size: 14px;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.track-index.is-playing {
  color: var(--brand-primary);
}

.eq-bar {
  display: inline-block;
  width: 12px;
  height: 12px;
  background: var(--brand-primary);
  mask: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12'%3E%3Crect x='0' y='4' width='2' height='8' rx='1'/%3E%3Crect x='5' y='0' width='2' height='12' rx='1'/%3E%3Crect x='10' y='6' width='2' height='6' rx='1'/%3E%3C/svg%3E") no-repeat center;
}

.track-detail {
  flex: 1;
  min-width: 0;
}

.track-title {
  font-size: 15px;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.track-item.playing .track-title {
  color: var(--brand-primary);
}

.track-subtitle {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.track-duration {
  font-size: 12px;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

/* Nav btn */
.nav-btn {
  display: flex;
  align-items: center;
  border: none;
  background: none;
  color: var(--brand-primary);
  cursor: pointer;
  padding: 6px;
}

.nav-btn svg {
  width: 22px;
  height: 22px;
}
</style>
