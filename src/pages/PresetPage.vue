<template>
  <div class="preset-page">
    <NavBar title="预设管理" back />

    <!-- 搜索和筛选 -->
    <div class="search-bar">
      <div class="search-input">
        <span class="search-icon">🔍</span>
        <input v-model="searchText" placeholder="搜索预设..." />
      </div>
      <button class="add-btn" @click="openEditor(null)">+ 新建</button>
    </div>

    <!-- 分类标签 -->
    <div class="category-tabs">
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

    <!-- 内置推荐 -->
    <div v-if="activeCategory === 'all' && !searchText" class="recommend-section">
      <div class="section-title">⭐ 推荐预设</div>
      <div class="recommend-scroll">
        <div
          v-for="p in recommendedPresets"
          :key="p.id"
          class="recommend-card"
          :style="{ background: p.gradient }"
          @click="applyPreset(p)"
        >
          <div class="rec-emoji">{{ p.emoji }}</div>
          <div class="rec-name">{{ p.name }}</div>
          <div class="rec-desc">{{ p.shortDesc }}</div>
          <div class="rec-badge">{{ p.category }}</div>
        </div>
      </div>
    </div>

    <!-- 预设列表 -->
    <div class="preset-list">
      <div class="section-title">
        {{ activeCategory === 'all' ? '全部预设' : categories.find(c => c.id === activeCategory)?.name }}
        <span class="count">({{ filteredPresets.length }})</span>
      </div>

      <div v-if="filteredPresets.length === 0" class="empty-state">
        <div class="empty-icon">📦</div>
        <div class="empty-text">暂无预设</div>
        <button class="empty-btn" @click="openEditor(null)">创建第一个预设</button>
      </div>

      <div
        v-for="preset in filteredPresets"
        :key="preset.id"
        class="preset-card"
      >
        <div class="preset-header" @click="toggleExpand(preset.id)">
          <div class="preset-icon">{{ preset.emoji }}</div>
          <div class="preset-info">
            <div class="preset-name">{{ preset.name }}</div>
            <div class="preset-meta">
              <span class="preset-category">{{ preset.category }}</span>
              <span class="preset-date">{{ preset.updatedAt }}</span>
            </div>
          </div>
          <div class="preset-actions-quick">
            <button class="action-apply" @click.stop="applyPreset(preset)" title="应用">▶</button>
            <span class="expand-arrow" :class="{ expanded: expandedId === preset.id }">›</span>
          </div>
        </div>

        <div v-if="expandedId === preset.id" class="preset-expand">
          <div class="expand-section">
            <div class="expand-label">描述</div>
            <div class="expand-value">{{ preset.description || '无描述' }}</div>
          </div>
          <div class="expand-section">
            <div class="expand-label">系统提示词</div>
            <div class="expand-value code">{{ truncate(preset.systemPrompt, 200) }}</div>
          </div>
          <div v-if="preset.greeting" class="expand-section">
            <div class="expand-label">开场白</div>
            <div class="expand-value">{{ preset.greeting }}</div>
          </div>
          <div class="expand-section">
            <div class="expand-label">参数</div>
            <div class="expand-params">
              <span class="param-tag">Temperature: {{ preset.temperature }}</span>
              <span class="param-tag">MaxTokens: {{ preset.maxTokens }}</span>
              <span v-if="preset.model" class="param-tag">Model: {{ preset.model }}</span>
            </div>
          </div>
          <div class="expand-actions">
            <button class="btn-edit" @click="openEditor(preset)">✏️ 编辑</button>
            <button class="btn-duplicate" @click="duplicatePreset(preset)">📋 复制</button>
            <button class="btn-export" @click="exportPreset(preset)">📤 导出</button>
            <button class="btn-delete" @click="deletePreset(preset.id)">🗑️ 删除</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑器弹窗 -->
    <div v-if="showEditor" class="modal-overlay" @click.self="showEditor = false">
      <div class="editor-panel">
        <div class="editor-header">
          <h3>{{ editingPreset ? '编辑预设' : '新建预设' }}</h3>
          <span class="close-btn" @click="showEditor = false">✕</span>
        </div>
        <div class="editor-body">
          <div class="form-group">
            <label>名称</label>
            <input v-model="form.name" placeholder="预设名称" />
          </div>
          <div class="form-group">
            <label>Emoji 图标</label>
            <div class="emoji-picker">
              <span
                v-for="e in emojiOptions"
                :key="e"
                class="emoji-option"
                :class="{ active: form.emoji === e }"
                @click="form.emoji = e"
              >
                {{ e }}
              </span>
            </div>
          </div>
          <div class="form-group">
            <label>分类</label>
            <select v-model="form.category">
              <option v-for="cat in categories.filter(c => c.id !== 'all')" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>描述</label>
            <input v-model="form.description" placeholder="简短描述" />
          </div>
          <div class="form-group">
            <label>系统提示词</label>
            <textarea v-model="form.systemPrompt" rows="5" placeholder="你是一个..."></textarea>
          </div>
          <div class="form-group">
            <label>开场白</label>
            <textarea v-model="form.greeting" rows="2" placeholder="（可选）角色的第一句话"></textarea>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Temperature</label>
              <input type="number" v-model.number="form.temperature" min="0" max="2" step="0.1" />
            </div>
            <div class="form-group half">
              <label>Max Tokens</label>
              <input type="number" v-model.number="form.maxTokens" min="100" max="4000" step="100" />
            </div>
          </div>
          <div class="form-group">
            <label>指定模型（可选）</label>
            <input v-model="form.model" placeholder="如 gpt-4o, claude-3.5-sonnet" />
          </div>
        </div>
        <div class="editor-footer">
          <button class="btn-cancel" @click="showEditor = false">取消</button>
          <button class="btn-save" @click="savePreset">💾 保存</button>
        </div>
      </div>
    </div>

    <!-- 应用成功提示 -->
    <div v-if="applyToast" class="apply-toast">
      ✅ 已应用预设「{{ applyToast }}」
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

interface Preset {
  id: number
  name: string
  emoji: string
  category: string
  description: string
  shortDesc: string
  systemPrompt: string
  greeting: string
  temperature: number
  maxTokens: number
  model: string
  gradient: string
  updatedAt: string
  isBuiltin: boolean
}

const categories = [
  { id: 'all', name: '全部', icon: '📋' },
  { id: 'roleplay', name: '角色扮演', icon: '🎭' },
  { id: 'assistant', name: '助手', icon: '🤖' },
  { id: 'creative', name: '创意写作', icon: '✍️' },
  { id: 'learning', name: '学习教育', icon: '📚' },
  { id: 'fun', name: '趣味', icon: '🎮' },
]

const activeCategory = ref('all')
const searchText = ref('')
const expandedId = ref<number | null>(null)
const showEditor = ref(false)
const editingPreset = ref<Preset | null>(null)
const applyToast = ref('')

const emojiOptions = ['🎭', '🤖', '✍️', '📚', '🎮', '💕', '🌸', '⚔️', '🔮', '🎵', '🌍', '💼', '🧠', '🎨', '🐱', '👻']

const form = ref({
  name: '',
  emoji: '🎭',
  category: '角色扮演',
  description: '',
  systemPrompt: '',
  greeting: '',
  temperature: 0.7,
  maxTokens: 2000,
  model: '',
})

const presets = ref<Preset[]>([
  { id: 1, name: '温柔女友', emoji: '💕', category: '角色扮演', description: '温柔体贴的女朋友角色', shortDesc: '甜蜜温柔的陪伴', systemPrompt: '你是用户的女朋友，性格温柔体贴，说话方式甜美可爱。你会关心对方的日常生活，分享有趣的事情，偶尔撒娇。回复时使用亲密的语气，适当加入emoji表情。', greeting: '亲爱的～你今天过得怎么样呀？我好想你呢 💕', temperature: 0.8, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #ff9a9e, #fecfef)', updatedAt: '2024-12-20', isBuiltin: true },
  { id: 2, name: '霸道总裁', emoji: '👔', category: '角色扮演', description: '高冷霸道的总裁角色', shortDesc: '傲娇又宠溺', systemPrompt: '你是一位年轻的集团总裁，外表高冷内心温柔。你对用户有特别的感情但不善表达，经常用傲娇的方式关心对方。说话简洁有力，偶尔会展现柔软的一面。', greeting: '又来找我了？...坐吧，我刚好有空。', temperature: 0.7, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #667eea, #764ba2)', updatedAt: '2024-12-18', isBuiltin: true },
  { id: 3, name: '知心好友', emoji: '🤝', category: '角色扮演', description: '能倾听和理解你的好朋友', shortDesc: '最懂你的朋友', systemPrompt: '你是用户最好的朋友，善于倾听和理解。你会认真回应对方的情绪，提供真诚的建议，有时候也会开玩笑活跃气氛。你们之间没有距离感。', greeting: '嘿！最近怎么样？有什么想聊的吗～', temperature: 0.8, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #43e97b, #38f9d7)', updatedAt: '2024-12-15', isBuiltin: true },
  { id: 4, name: '创意写手', emoji: '✍️', category: '创意写作', description: '帮你写小说、诗歌、文案', shortDesc: '灵感与文字的魔法师', systemPrompt: '你是一位才华横溢的创意写作助手。你擅长各种文体：小说、散文、诗歌、广告文案、剧本等。你会根据用户的需求提供创意灵感，帮助构思情节、打磨文字。你的语言优美而富有想象力。', greeting: '你好！今天想创作什么呢？一个动人的故事，一首诗，还是...？', temperature: 0.9, maxTokens: 3000, model: '', gradient: 'linear-gradient(135deg, #a18cd1, #fbc2eb)', updatedAt: '2024-12-12', isBuiltin: true },
  { id: 5, name: '英语老师', emoji: '📚', category: '学习教育', description: '耐心的英语教学', shortDesc: '你的专属英语教师', systemPrompt: '你是一位经验丰富且耐心的英语老师。你会用中英双语教学，解释语法、词汇和表达方式。你善于举例说明，会纠正用户的错误并给出改进建议。教学风格轻松有趣。', greeting: 'Hello! 今天我们来学什么呢？Feel free to ask me anything! 😊', temperature: 0.5, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #4facfe, #00f2fe)', updatedAt: '2024-12-10', isBuiltin: true },
  { id: 6, name: '心理咨询师', emoji: '🧠', category: '助手', description: '专业的心理疏导', shortDesc: '倾听你的内心', systemPrompt: '你是一位专业的心理咨询师，擅长认知行为疗法和正念技术。你善于倾听，用温和的方式引导用户探索内心感受。你不会做诊断，但会提供有效的情绪管理建议和心理支持。', greeting: '你好，很高兴你来找我聊天。现在感觉怎么样？我在这里认真听你说。', temperature: 0.6, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #ffecd2, #fcb69f)', updatedAt: '2024-12-08', isBuiltin: true },
  { id: 7, name: '古风侠客', emoji: '⚔️', category: '角色扮演', description: '武侠世界中的侠客', shortDesc: '仗剑走天涯', systemPrompt: '你是一位武侠世界中的侠客，性格豪爽仗义，武艺高强。你行走江湖，路见不平拔刀相助。说话风格为古风，用词典雅。你对友人真诚，对敌人果决。', greeting: '在下刚从塞北归来，正想找一处酒肆歇脚。阁下不如一同饮上几杯？', temperature: 0.8, maxTokens: 2000, model: '', gradient: 'linear-gradient(135deg, #89f7fe, #66a6ff)', updatedAt: '2024-12-05', isBuiltin: true },
  { id: 8, name: '猫娘', emoji: '🐱', category: '趣味', description: '可爱的猫娘角色', shortDesc: '喵～超可爱的', systemPrompt: '你是一个可爱的猫娘，会在说话时加入"喵"和各种可爱的语气词。你很黏人，喜欢被摸头和夸奖。性格活泼调皮，有时会撒娇耍赖。你对主人非常忠诚和依赖。', greeting: '主人回来了喵～今天好想你哦！能摸摸我的头吗？喵～ 🐱', temperature: 0.9, maxTokens: 1500, model: '', gradient: 'linear-gradient(135deg, #fa709a, #fee140)', updatedAt: '2024-12-01', isBuiltin: true },
])

const recommendedPresets = computed(() => presets.value.filter(p => p.isBuiltin).slice(0, 4))

const filteredPresets = computed(() => {
  let list = presets.value
  if (activeCategory.value !== 'all') {
    const catName = categories.find(c => c.id === activeCategory.value)?.name || ''
    list = list.filter(p => p.category === catName)
  }
  if (searchText.value.trim()) {
    const kw = searchText.value.trim().toLowerCase()
    list = list.filter(p =>
      p.name.toLowerCase().includes(kw) ||
      p.description.toLowerCase().includes(kw) ||
      p.systemPrompt.toLowerCase().includes(kw)
    )
  }
  return list
})

function toggleExpand(id: number) {
  expandedId.value = expandedId.value === id ? null : id
}

function truncate(s: string, len: number): string {
  return s.length > len ? s.slice(0, len) + '...' : s
}

function openEditor(preset: Preset | null) {
  editingPreset.value = preset
  if (preset) {
    form.value = {
      name: preset.name,
      emoji: preset.emoji,
      category: preset.category,
      description: preset.description,
      systemPrompt: preset.systemPrompt,
      greeting: preset.greeting,
      temperature: preset.temperature,
      maxTokens: preset.maxTokens,
      model: preset.model,
    }
  } else {
    form.value = {
      name: '',
      emoji: '🎭',
      category: '角色扮演',
      description: '',
      systemPrompt: '',
      greeting: '',
      temperature: 0.7,
      maxTokens: 2000,
      model: '',
    }
  }
  showEditor.value = true
}

function savePreset() {
  if (!form.value.name.trim() || !form.value.systemPrompt.trim()) return

  const now = new Date()
  const dateStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

  if (editingPreset.value) {
    const idx = presets.value.findIndex(p => p.id === editingPreset.value!.id)
    if (idx >= 0) {
      presets.value[idx] = {
        ...presets.value[idx],
        ...form.value,
        shortDesc: form.value.description.slice(0, 20),
        updatedAt: dateStr,
      }
    }
  } else {
    presets.value.unshift({
      id: Date.now(),
      ...form.value,
      shortDesc: form.value.description.slice(0, 20),
      gradient: 'linear-gradient(135deg, #667eea, #764ba2)',
      updatedAt: dateStr,
      isBuiltin: false,
    })
  }

  showEditor.value = false
}

function duplicatePreset(p: Preset) {
  const now = new Date()
  const dateStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
  presets.value.unshift({
    ...p,
    id: Date.now(),
    name: p.name + ' (副本)',
    isBuiltin: false,
    updatedAt: dateStr,
  })
}

function deletePreset(id: number) {
  presets.value = presets.value.filter(p => p.id !== id)
  if (expandedId.value === id) expandedId.value = null
}

function applyPreset(p: Preset) {
  applyToast.value = p.name
  setTimeout(() => { applyToast.value = '' }, 2000)
}

function exportPreset(p: Preset) {
  const data = {
    name: p.name,
    emoji: p.emoji,
    category: p.category,
    description: p.description,
    systemPrompt: p.systemPrompt,
    greeting: p.greeting,
    temperature: p.temperature,
    maxTokens: p.maxTokens,
    model: p.model,
  }
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `preset-${p.name}.json`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.preset-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.search-bar {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-tertiary);
  border-radius: 18px;
  padding: 8px 12px;
}

.search-input input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  font-size: 14px;
  color: var(--text-primary);
}

.search-icon { font-size: 14px; opacity: 0.5; }

.add-btn {
  padding: 8px 16px;
  border-radius: 18px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  flex-shrink: 0;
}

.category-tabs {
  display: flex;
  gap: 4px;
  padding: 4px 12px 8px;
  overflow-x: auto;
}

.category-tabs::-webkit-scrollbar { display: none; }

.cat-tab {
  padding: 5px 12px;
  border-radius: 14px;
  font-size: 12px;
  white-space: nowrap;
  cursor: pointer;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  transition: all 0.2s;
  flex-shrink: 0;
}

.cat-tab.active {
  background: var(--color-primary);
  color: #fff;
  font-weight: 600;
}

/* 推荐 */
.section-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  padding: 8px 16px 6px;
}

.section-title .count {
  font-size: 13px;
  font-weight: 400;
  color: var(--text-tertiary);
}

.recommend-scroll {
  display: flex;
  gap: 10px;
  padding: 0 12px 8px;
  overflow-x: auto;
}

.recommend-scroll::-webkit-scrollbar { display: none; }

.recommend-card {
  flex-shrink: 0;
  width: 140px;
  padding: 14px;
  border-radius: 16px;
  color: #fff;
  cursor: pointer;
  position: relative;
  transition: transform 0.2s;
}

.recommend-card:active { transform: scale(0.96); }

.rec-emoji { font-size: 28px; margin-bottom: 6px; }
.rec-name { font-size: 14px; font-weight: 700; margin-bottom: 2px; }
.rec-desc { font-size: 11px; opacity: 0.85; }

.rec-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 9px;
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.25);
}

/* 预设列表 */
.preset-list {
  padding-bottom: 20px;
}

.preset-card {
  margin: 4px 12px;
  background: var(--bg-secondary);
  border-radius: 14px;
  overflow: hidden;
}

.preset-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  cursor: pointer;
}

.preset-icon { font-size: 28px; flex-shrink: 0; }

.preset-info { flex: 1; min-width: 0; }

.preset-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.preset-meta {
  display: flex;
  gap: 8px;
  margin-top: 2px;
}

.preset-category {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
}

.preset-date {
  font-size: 11px;
  color: var(--text-quaternary);
}

.preset-actions-quick {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.action-apply {
  width: 30px;
  height: 30px;
  border-radius: 15px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.expand-arrow {
  font-size: 18px;
  color: var(--text-quaternary);
  transition: transform 0.2s;
}

.expand-arrow.expanded { transform: rotate(90deg); }

/* 展开区 */
.preset-expand {
  padding: 0 14px 12px;
  border-top: 1px solid var(--border-secondary);
}

.expand-section {
  margin-top: 10px;
}

.expand-label {
  font-size: 11px;
  color: var(--text-tertiary);
  font-weight: 600;
  margin-bottom: 3px;
  text-transform: uppercase;
}

.expand-value {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.expand-value.code {
  font-family: monospace;
  font-size: 12px;
  background: var(--bg-tertiary);
  padding: 8px;
  border-radius: 8px;
  white-space: pre-wrap;
  word-break: break-all;
}

.expand-params {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.param-tag {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 6px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.expand-actions {
  display: flex;
  gap: 6px;
  margin-top: 12px;
}

.expand-actions button {
  flex: 1;
  padding: 8px;
  border-radius: 10px;
  border: none;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
}

.btn-edit { background: var(--color-primary); color: #fff; }
.btn-duplicate { background: var(--bg-tertiary); color: var(--text-primary); }
.btn-export { background: var(--bg-tertiary); color: var(--text-primary); }
.btn-delete { background: rgba(255, 59, 48, 0.1); color: #ff3b30; }

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-tertiary);
}

.empty-icon { font-size: 40px; margin-bottom: 8px; }
.empty-text { font-size: 14px; margin-bottom: 12px; }

.empty-btn {
  padding: 8px 20px;
  border-radius: 18px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

/* 编辑器 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.editor-panel {
  width: 100%;
  max-width: 393px;
  max-height: 90%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-primary);
}

.editor-header h3 {
  margin: 0;
  font-size: 17px;
  color: var(--text-primary);
}

.close-btn {
  width: 28px;
  height: 28px;
  border-radius: 14px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
}

.editor-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
}

.form-group {
  margin-bottom: 12px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--border-primary);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  resize: vertical;
  box-sizing: border-box;
}

.form-group textarea {
  font-family: inherit;
}

.form-row {
  display: flex;
  gap: 10px;
}

.form-group.half { flex: 1; }

.emoji-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.emoji-option {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  cursor: pointer;
  background: var(--bg-tertiary);
  border: 2px solid transparent;
  transition: all 0.2s;
}

.emoji-option.active {
  border-color: var(--color-primary);
  background: var(--bg-secondary);
}

.editor-footer {
  display: flex;
  gap: 10px;
  padding: 12px 16px;
  border-top: 1px solid var(--border-primary);
}

.editor-footer button {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  border: none;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel { background: var(--bg-tertiary); color: var(--text-primary); }
.btn-save { background: var(--color-primary); color: #fff; }

/* 应用提示 */
.apply-toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  border-radius: 20px;
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  z-index: 200;
  white-space: nowrap;
}
</style>
