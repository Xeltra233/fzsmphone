<template>
  <div class="takeaway-page">
    <NavBar title="外卖" back-to="/">
      <template #right>
        <button class="icon-btn" @click="$router.push('/takeaway/orders')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2" />
            <rect x="9" y="3" width="6" height="4" rx="2" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- Search -->
    <div class="search-wrap">
      <div class="search-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8" /><path d="M21 21l-4.35-4.35" />
        </svg>
        <input v-model="search" type="text" placeholder="搜索餐厅或菜品" />
      </div>
    </div>

    <!-- Categories -->
    <div class="categories">
      <button
        v-for="cat in categories"
        :key="cat.id"
        class="cat-btn"
        :class="{ active: currentCat === cat.id }"
        @click="currentCat = cat.id"
      >
        <span class="cat-emoji">{{ cat.emoji }}</span>
        <span class="cat-label">{{ cat.label }}</span>
      </button>
    </div>

    <!-- Restaurant List -->
    <div class="restaurant-list">
      <div v-if="filteredRestaurants.length === 0" class="empty-state">
        <div class="empty-emoji">◈</div>
        <div class="empty-title">暂无餐厅</div>
      </div>

      <div
        v-for="r in filteredRestaurants"
        :key="r.id"
        class="restaurant-card"
        @click="$router.push(`/takeaway/restaurant/${r.id}`)"
      >
        <div class="rest-img">
          <img v-if="r.image_url" :src="r.image_url" :alt="r.name" />
          <span v-else class="rest-img-placeholder">◈</span>
        </div>
        <div class="rest-info">
          <div class="rest-name">{{ r.name }}</div>
          <div class="rest-meta">
            <span class="rest-rating">★ {{ r.rating }}</span>
            <span class="rest-sep">·</span>
            <span class="rest-time">{{ r.delivery_time }}</span>
          </div>
          <div class="rest-tags">
            <span class="rest-cat">{{ r.category }}</span>
            <span v-if="r.min_order > 0" class="rest-min">¥{{ r.min_order }}起送</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { api } from '@/api/client'

interface Restaurant {
  id: number
  name: string
  category: string
  image_url: string
  rating: number
  delivery_time: string
  min_order: number
}

const categories = [
  { id: 'all', emoji: '▲', label: '全部' },
  { id: '快餐', emoji: '◆', label: '快餐' },
  { id: '中餐', emoji: '◈', label: '中餐' },
  { id: '火锅', emoji: '◎', label: '火锅' },
  { id: '烧烤', emoji: '◎', label: '烧烤' },
  { id: '甜品', emoji: '●', label: '甜品' },
  { id: '饮品', emoji: '◇', label: '饮品' },
]

const search = ref('')
const currentCat = ref('all')
const restaurants = ref<Restaurant[]>([])

const filteredRestaurants = computed(() => {
  let list = restaurants.value
  if (currentCat.value !== 'all') {
    list = list.filter(r => r.category === currentCat.value)
  }
  const q = search.value.trim().toLowerCase()
  if (q) {
    list = list.filter(r => r.name.toLowerCase().includes(q) || r.category.toLowerCase().includes(q))
  }
  return list
})

async function fetchRestaurants() {
  try {
    const res = await api.get<{ data: Restaurant[] }>('/api/restaurants')
    restaurants.value = res.data || []
  } catch { /* empty */ }
}

onMounted(fetchRestaurants)
</script>

<style scoped>
.takeaway-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
}

.search-wrap {
  padding: 8px 16px;
  background: var(--bg-primary);
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--fill-tertiary);
  border-radius: 10px;
  padding: 0 12px;
  height: 36px;
}

.search-box svg { width: 16px; height: 16px; color: var(--text-quaternary); flex-shrink: 0; }

.search-box input {
  flex: 1; border: none; background: none; font-size: 14px;
  color: var(--text-primary); outline: none;
}
.search-box input::placeholder { color: var(--text-quaternary); }

/* Categories */
.categories {
  display: flex;
  gap: 0;
  padding: 10px 12px;
  background: var(--bg-primary);
  border-bottom: 0.5px solid var(--separator);
  overflow-x: auto;
  flex-shrink: 0;
}
.categories::-webkit-scrollbar { display: none; }

.cat-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  border: none;
  background: none;
  cursor: pointer;
  min-width: 56px;
  border-radius: 12px;
  transition: all 0.15s;
}

.cat-btn.active { background: rgba(88, 86, 214, 0.1); }
.cat-emoji { font-size: 24px; }
.cat-label { font-size: 11px; color: var(--text-secondary); white-space: nowrap; }
.cat-btn.active .cat-label { color: var(--brand-primary); font-weight: 600; }

/* Restaurant List */
.restaurant-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.restaurant-card {
  display: flex;
  gap: 12px;
  background: var(--bg-primary);
  border-radius: 14px;
  padding: 12px;
  cursor: pointer;
  transition: transform 0.15s var(--ease-ios);
}

.restaurant-card:active { transform: scale(0.98); }

.rest-img {
  width: 80px;
  height: 80px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--fill-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.rest-img img { width: 100%; height: 100%; object-fit: cover; }
.rest-img-placeholder { font-size: 32px; }

.rest-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
}

.rest-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rest-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--text-secondary);
}

.rest-rating { color: var(--ios-orange); }
.rest-sep { color: var(--text-quaternary); }

.rest-tags {
  display: flex;
  gap: 8px;
  font-size: 12px;
}

.rest-cat {
  color: var(--text-tertiary);
  background: var(--fill-tertiary);
  padding: 1px 6px;
  border-radius: 4px;
}

.rest-min {
  color: var(--text-tertiary);
}

.empty-state { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; }
.empty-emoji { font-size: 48px; margin-bottom: 12px; }
.empty-title { font-size: 17px; font-weight: 600; color: var(--text-secondary); }

.icon-btn {
  width: 32px; height: 32px; border: none; background: none;
  color: var(--brand-primary); cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.icon-btn:active { opacity: 0.5; }
.icon-btn svg { width: 22px; height: 22px; }
</style>
