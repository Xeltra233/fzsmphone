<template>
  <div class="map-page">
    <NavBar title="地图" back />

    <!-- Empty -->
    <div v-if="store.mapLocations.length === 0 && !store.mapLoading" class="empty-state">
      <div class="empty-icon"><svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 12v32l14-8 16 8 14-8V4l-14 8-16-8-14 8z"/><line x1="16" y1="4" x2="16" y2="36"/><line x1="32" y1="12" x2="32" y2="44"/></svg></div>
      <div class="empty-title">探索地图</div>
      <div class="empty-sub">让AI生成当前世界的地点信息</div>
      <button class="gen-btn" :disabled="store.generating" @click="store.generateMapContent()">
        ◆ AI 生成地图
      </button>
    </div>

    <!-- Loading -->
    <div v-else-if="store.mapLoading" class="loading-wrap">
      <div class="loading-spinner"></div>
      <span>正在探索...</span>
    </div>

    <!-- Location List -->
    <div v-else-if="!selectedLocation" class="loc-list">
      <div v-for="loc in store.mapLocations" :key="loc.id" class="loc-card" @click="selectedLocation = loc">
        <div class="loc-icon" :style="{ background: typeColor(loc.type) }">
          <span>{{ typeEmoji(loc.type) }}</span>
        </div>
        <div class="loc-info">
          <div class="loc-name">{{ loc.name }}</div>
          <div class="loc-type">{{ loc.type }}</div>
          <div class="loc-chars" v-if="loc.characters.length && loc.characters[0] !== '无'">
            <span v-for="c in loc.characters.slice(0,3)" :key="c" class="char-dot">{{ c }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Location Detail -->
    <div v-else class="loc-detail">
      <button class="back-link" @click="selectedLocation = null">← 返回地图</button>
      <div class="hero" :style="{ background: typeColor(selectedLocation.type) }">
        <span class="hero-emoji">{{ typeEmoji(selectedLocation.type) }}</span>
        <h2 class="hero-name">{{ selectedLocation.name }}</h2>
        <span class="hero-type">{{ selectedLocation.type }}</span>
      </div>
      <div class="detail-section">
        <div class="section-title">◈ 描述</div>
        <div class="section-content">{{ selectedLocation.description }}</div>
      </div>
      <div class="detail-section" v-if="selectedLocation.characters.length && selectedLocation.characters[0] !== '无'">
        <div class="section-title">◎ 在场角色</div>
        <div class="char-tags">
          <span v-for="c in selectedLocation.characters" :key="c" class="char-tag">{{ c }}</span>
        </div>
      </div>
      <div class="detail-section" v-if="selectedLocation.events.length && selectedLocation.events[0] !== '平静'">
        <div class="section-title">※ 正在发生</div>
        <div class="event-list">
          <div v-for="e in selectedLocation.events" :key="e" class="event-item">{{ e }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'
import type { MapLocation } from '@/utils/socialParsers'

const store = useSocialAIStore()
const selectedLocation = ref<MapLocation | null>(null)

onMounted(() => { store.loadData('map') })

const TYPE_MAP: Record<string, { symbol: string; color: string }> = {
  '住宅': { symbol: '⌂', color: 'linear-gradient(135deg,#a8e6cf,#55efc4)' },
  '商业': { symbol: '◈', color: 'linear-gradient(135deg,#ffeaa7,#fdcb6e)' },
  '学校': { symbol: '▣', color: 'linear-gradient(135deg,#74b9ff,#0984e3)' },
  '公园': { symbol: '♧', color: 'linear-gradient(135deg,#55efc4,#00b894)' },
  '餐厅': { symbol: '◎', color: 'linear-gradient(135deg,#fab1a0,#e17055)' },
  '景点': { symbol: '△', color: 'linear-gradient(135deg,#dfe6e9,#b2bec3)' },
  '办公': { symbol: '■', color: 'linear-gradient(135deg,#a29bfe,#6c5ce7)' },
  '娱乐': { symbol: '♦', color: 'linear-gradient(135deg,#fd79a8,#e84393)' },
}

function typeEmoji(type: string) { return TYPE_MAP[type]?.symbol || '◇' }
function typeColor(type: string) { return TYPE_MAP[type]?.color || 'linear-gradient(135deg,#636e72,#b2bec3)' }
</script>

<style scoped>
.map-page { height: 100%; display: flex; flex-direction: column; background: var(--bg-primary); }
.empty-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; padding: 40px 20px; }
.empty-icon { width: 56px; height: 56px; color: var(--text-tertiary); }
.empty-icon svg { width: 100%; height: 100%; }
.empty-title { font-size: 18px; font-weight: 700; color: var(--text-primary); }
.empty-sub { font-size: 13px; color: var(--text-tertiary); }
.gen-btn { padding: 10px 24px; border: none; border-radius: 20px; background: var(--brand-primary); color: #fff; font-size: 14px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.gen-btn:disabled { opacity: 0.5; }
.loading-wrap { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; color: var(--text-secondary); }
.loading-spinner { width: 32px; height: 32px; border: 3px solid var(--separator); border-top-color: var(--brand-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.loc-list { flex: 1; overflow-y: auto; padding: 8px 16px 80px; display: flex; flex-direction: column; gap: 8px; }
.loc-card { display: flex; align-items: center; gap: 12px; padding: 14px; background: var(--bg-secondary); border-radius: 14px; cursor: pointer; }
.loc-icon { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 20px; flex-shrink: 0; }
.loc-info { flex: 1; min-width: 0; }
.loc-name { font-size: 15px; font-weight: 600; color: var(--text-primary); }
.loc-type { font-size: 12px; color: var(--text-tertiary); margin-top: 2px; }
.loc-chars { display: flex; gap: 6px; margin-top: 4px; }
.char-dot { font-size: 11px; padding: 2px 8px; background: var(--fill-tertiary); border-radius: 10px; color: var(--text-secondary); }

.loc-detail { flex: 1; overflow-y: auto; padding: 0 0 80px; }
.back-link { background: none; border: none; color: var(--brand-primary); font-size: 14px; padding: 12px 16px; cursor: pointer; }
.hero { padding: 24px; text-align: center; }
.hero-emoji { font-size: 36px; }
.hero-name { font-size: 22px; font-weight: 700; color: #fff; margin: 8px 0 4px; }
.hero-type { font-size: 13px; color: rgba(255,255,255,0.7); }
.detail-section { padding: 12px 16px; }
.section-title { font-size: 14px; font-weight: 700; color: var(--text-primary); margin-bottom: 6px; }
.section-content { font-size: 15px; color: var(--text-secondary); line-height: 1.6; }
.char-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.char-tag { padding: 4px 12px; background: var(--brand-primary); color: #fff; border-radius: 12px; font-size: 13px; }
.event-list { display: flex; flex-direction: column; gap: 6px; }
.event-item { padding: 8px 12px; background: var(--fill-tertiary); border-radius: 10px; font-size: 14px; color: var(--text-secondary); }
</style>
