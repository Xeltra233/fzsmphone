<template>
  <div class="browser-page">
    <NavBar title="浏览器" back />

    <!-- Search Bar -->
    <div class="search-wrap">
      <div class="search-bar">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="searchQuery" class="search-input" placeholder="搜索..." @keydown.enter="doSearch" />
      </div>
    </div>

    <!-- Empty Hint -->
    <div v-if="store.browserResults.length === 0 && !store.browserLoading" class="empty-state">
      <div class="empty-icon"><svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="2"><circle cx="24" cy="24" r="20"/><line x1="4" y1="24" x2="44" y2="24"/><path d="M24 4a30.6 30.6 0 0 1 8 20 30.6 30.6 0 0 1-8 20 30.6 30.6 0 0 1-8-20 30.6 30.6 0 0 1 8-20z"/></svg></div>
      <div class="empty-title">搜索一下</div>
      <div class="empty-sub">输入关键词，AI将生成符合世界观的搜索结果</div>
      <button class="gen-btn" :disabled="store.generating" @click="store.generateBrowserContent()">
        ◆ 生成推荐内容
      </button>
    </div>

    <!-- Loading -->
    <div v-else-if="store.browserLoading" class="loading-wrap">
      <div class="loading-spinner"></div>
      <span>搜索中...</span>
    </div>

    <!-- Results List -->
    <div v-else-if="!selectedResult" class="results-list">
      <div v-for="r in store.browserResults" :key="r.id" class="result-item" @click="selectedResult = r">
        <div class="result-title">{{ r.title }}</div>
        <div class="result-url">{{ r.url }}</div>
        <div class="result-summary">{{ r.summary }}</div>
      </div>
    </div>

    <!-- Article Detail -->
    <div v-else class="article-view">
      <button class="back-link" @click="selectedResult = null">← 返回结果</button>
      <h1 class="article-title">{{ selectedResult.title }}</h1>
      <div class="article-url">{{ selectedResult.url }}</div>
      <div class="article-body">{{ selectedResult.content }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useSocialAIStore } from '@/stores/socialAI'
import type { BrowserResult } from '@/utils/socialParsers'

const store = useSocialAIStore()
const searchQuery = ref('')
const selectedResult = ref<BrowserResult | null>(null)

onMounted(() => { store.loadData('browser') })

function doSearch() {
  if (!searchQuery.value.trim()) return
  store.generateBrowserContent(`用户搜索了：${searchQuery.value}`)
}
</script>

<style scoped>
.browser-page { height: 100%; display: flex; flex-direction: column; background: var(--bg-primary); }
.search-wrap { padding: 8px 16px; }
.search-bar { display: flex; align-items: center; gap: 8px; background: var(--bg-tertiary); border-radius: 12px; padding: 0 14px; }
.search-icon { width: 18px; height: 18px; color: var(--text-tertiary); flex-shrink: 0; }
.search-input { flex: 1; border: none; background: none; padding: 12px 0; font-size: 15px; color: var(--text-primary); outline: none; }

.empty-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 8px; padding: 40px 20px; }
.empty-icon { width: 56px; height: 56px; color: var(--text-tertiary); }
.empty-icon svg { width: 100%; height: 100%; }
.empty-title { font-size: 18px; font-weight: 700; color: var(--text-primary); }
.empty-sub { font-size: 13px; color: var(--text-tertiary); text-align: center; }
.gen-btn { padding: 10px 24px; border: none; border-radius: 20px; background: var(--brand-primary); color: #fff; font-size: 14px; font-weight: 600; cursor: pointer; margin-top: 8px; }
.gen-btn:disabled { opacity: 0.5; }
.loading-wrap { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; color: var(--text-secondary); }
.loading-spinner { width: 32px; height: 32px; border: 3px solid var(--separator); border-top-color: var(--brand-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.results-list { flex: 1; overflow-y: auto; padding: 0 16px 80px; }
.result-item { padding: 14px 0; border-bottom: 0.5px solid var(--separator); cursor: pointer; }
.result-title { font-size: 16px; font-weight: 600; color: var(--brand-primary); }
.result-url { font-size: 12px; color: var(--text-tertiary); margin-top: 2px; }
.result-summary { font-size: 14px; color: var(--text-secondary); margin-top: 6px; line-height: 1.5; }

.article-view { flex: 1; overflow-y: auto; padding: 0 16px 80px; }
.back-link { background: none; border: none; color: var(--brand-primary); font-size: 14px; padding: 12px 0; cursor: pointer; }
.article-title { font-size: 22px; font-weight: 700; color: var(--text-primary); margin: 0; }
.article-url { font-size: 12px; color: var(--text-tertiary); margin-top: 4px; }
.article-body { padding: 16px 0; font-size: 15px; color: var(--text-primary); line-height: 1.7; white-space: pre-wrap; }
</style>
