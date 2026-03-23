<template>
  <div class="worldbook-page">
    <NavBar title="世界书" back-to="/friends">
      <template #right>
        <button class="icon-btn" @click="toggleBatchMode" :title="batchMode ? '退出批量管理' : '批量管理'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <rect x="4" y="5" width="16" height="4" rx="1.5" />
            <rect x="4" y="10" width="16" height="4" rx="1.5" />
            <rect x="4" y="15" width="16" height="4" rx="1.5" />
          </svg>
        </button>
        <button class="icon-btn" @click="showBookMenu = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <path d="M12 5v14M5 12h14" />
          </svg>
        </button>
      </template>
    </NavBar>

    <!-- 世界书选择器 -->
    <div v-if="batchMode" class="batch-toolbar">
      <div class="batch-summary">已选 {{ selectedBookIds.length }} 本</div>
      <div class="batch-actions">
        <button class="batch-btn" @click="toggleSelectAllBooks">{{ allBooksSelected ? '取消全选' : '全选当前' }}</button>
        <button class="batch-btn" :disabled="selectedBookIds.length === 0" @click="batchExportBooks">批量导出</button>
        <button class="batch-btn danger" :disabled="selectedBookIds.length === 0" @click="batchDeleteBooks">批量删除</button>
      </div>
    </div>

    <div class="book-tabs">
      <div
        v-for="book in worldBooks"
        :key="book.id"
        class="book-tab"
        :class="{ active: currentBookId === book.id, selected: batchMode && selectedBookIds.includes(String(book.id)) }"
        @click="batchMode ? toggleBookSelection(book.id) : currentBookId = book.id"
      >
        <label v-if="batchMode" class="batch-check" @click.stop>
          <input type="checkbox" :checked="selectedBookIds.includes(String(book.id))" @change="toggleBookSelection(book.id)" />
        </label>
        <span class="book-tab-icon">▤</span>
        <span class="book-tab-name">{{ book.name }}</span>
        <span class="book-tab-count">{{ book.entries.length }}</span>
      </div>
      <div class="book-tab add-tab" @click="showBookMenu = true">
        <span>+</span>
      </div>
    </div>

    <div class="page-content" v-if="currentBook">
      <!-- 当前世界书信息 -->
      <div class="book-info-bar">
        <div class="book-info-left">
          <div class="book-title">{{ currentBook.name }}</div>
          <div class="book-meta">
            {{ currentBook.entries.length }} 个条目
            <span v-if="currentBook.bindChars.length > 0"> · 绑定 {{ currentBook.bindChars.length }} 个角色</span>
          </div>
        </div>
        <div class="book-info-actions">
          <button class="small-btn" @click="editBookName">✎</button>
          <button class="small-btn" @click="exportBook">↑</button>
          <button class="small-btn danger" @click="deleteBook" v-if="worldBooks.length > 1">✕</button>
        </div>
      </div>

      <!-- 搜索和添加 -->
      <div class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">◎</span>
          <input v-model="searchText" placeholder="搜索条目..." />
        </div>
        <button class="ai-batch-btn" @click="showAIBatchModal = true" title="AI批量创建条目">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="16" height="16">
            <path d="M15 4V2m0 2v2m0-2h2m-2 0h-2"/>
            <path d="M8 9l3-7 3 7"/>
            <path d="M6 15l3-3h6l3 3"/>
            <path d="M4 20h16"/>
          </svg>
          <span>AI创建</span>
        </button>
        <button class="add-entry-btn" @click="openEntryEditor(null)">+ 新条目</button>
      </div>

      <!-- 条目列表 -->
      <div v-if="filteredEntries.length === 0 && !showEntryForm" class="empty">
        <div class="empty-icon">▥</div>
        <p>还没有世界书条目</p>
        <p class="empty-hint">世界书条目会在聊天中根据关键词自动注入到AI的上下文中</p>
        <button class="empty-btn" @click="openEntryEditor(null)">创建第一个条目</button>
      </div>

      <div class="entry-list" v-if="filteredEntries.length > 0">
        <div
          v-for="entry in filteredEntries"
          :key="entry.id"
          class="entry-card"
          :class="{ disabled: !entry.enabled }"
        >
          <div class="entry-header" @click="toggleEntryExpand(entry.id)">
            <div class="entry-title-row">
              <span class="entry-status-light" :class="entryStatusClass(entry)" :title="entryStatusLabel(entry)"></span>
              <span class="entry-title">{{ entry.title }}</span>
              <span class="entry-status-tag">{{ entryStatusLabel(entry) }}</span>
            </div>
            <div class="entry-keywords-preview" v-if="entry.keywords.length > 0">
              <span class="kw-tag" v-for="kw in entry.keywords.slice(0, 5)" :key="kw">{{ kw }}</span>
              <span v-if="entry.keywords.length > 5" class="kw-more">+{{ entry.keywords.length - 5 }}</span>
            </div>
          </div>

          <div v-if="expandedEntryId === entry.id" class="entry-expand">
            <div class="entry-content-preview">{{ entry.content }}</div>
            <div class="entry-actions">
              <button class="act-btn" @click="toggleEntryEnabled(entry)">
                {{ entry.enabled ? '禁用' : '启用' }}
              </button>
              <button class="act-btn" @click="openEntryEditor(entry)">编辑</button>
              <button class="act-btn danger" @click="deleteEntry(entry.id)">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 无世界书时的空状态 -->
    <div v-if="worldBooks.length === 0" class="page-content">
      <div class="empty">
        <div class="empty-icon">▥</div>
        <p>还没有世界书</p>
        <p class="empty-hint">世界书用于在聊天中自动注入背景设定信息</p>
        <button class="empty-btn" @click="createNewBook">创建第一本世界书</button>
      </div>
    </div>

    <!-- 底部操作 -->
    <div class="bottom-bar">
      <button class="import-btn" @click="triggerImport">↓ 导入世界书</button>
      <input ref="importInput" type="file" accept=".json" style="display:none" @change="handleImport" />
    </div>

    <!-- 条目编辑弹窗 -->
    <div v-if="showEntryForm" class="modal-overlay" @click.self="showEntryForm = false">
      <div class="editor-panel">
        <div class="editor-header">
          <h3>{{ editingEntry ? '编辑条目' : '新建条目' }}</h3>
          <span class="close-btn" @click="showEntryForm = false">✕</span>
        </div>
        <div class="editor-body">
          <div class="form-group">
            <label>条目标题</label>
            <input v-model="entryForm.title" placeholder="条目名称" />
          </div>
          <div class="form-group">
            <label>触发关键词</label>
            <input v-model="keywordsInput" placeholder="多个关键词用逗号分隔，如：魔法,法术,施法" />
            <div class="form-hint">当聊天内容中出现这些关键词时，条目内容会自动注入到AI上下文中</div>
          </div>
          <div class="form-group">
            <label>条目内容</label>
            <textarea v-model="entryForm.content" rows="6" placeholder="当关键词被触发时，注入到AI上下文的背景信息..."></textarea>
          </div>
          <div class="form-row">
            <label class="form-check">
              <input type="checkbox" v-model="entryForm.enabled" />
              <span>启用此条目</span>
            </label>
          </div>
          <!-- 高级设置 -->
          <div class="advanced-toggle" @click="showAdvanced = !showAdvanced">
            <span>⚙️ 高级设置 (SillyTavern)</span>
            <span class="toggle-arrow" :class="{ expanded: showAdvanced }">›</span>
          </div>
          <template v-if="showAdvanced">
            <div class="form-group">
              <label>次要关键词 (Secondary Keys)</label>
              <input v-model="secondaryKeysInput" placeholder="次要关键词，逗号分隔（selective模式下需同时匹配）" />
              <div class="form-hint">当selective开启时，主关键词和次要关键词需同时匹配才触发</div>
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.constant" />
                <span>常驻注入 (Constant)</span>
              </label>
              <div class="form-hint inline-hint">无需关键词匹配，始终注入到上下文中</div>
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.selective" />
                <span>选择性匹配 (Selective)</span>
              </label>
              <div class="form-hint inline-hint">需同时匹配主关键词和次要关键词</div>
            </div>
            <div class="form-group compact">
              <label>选择性逻辑 (Selective Logic)</label>
              <select v-model.number="entryForm.selectiveLogic">
                <option :value="0">AND ANY — 任意一个次要关键词匹配即可</option>
                <option :value="1">NOT ALL — 非所有次要关键词匹配时触发</option>
                <option :value="2">NOT ANY — 无任何次要关键词匹配时触发</option>
                <option :value="3">AND ALL — 所有次要关键词都须匹配</option>
              </select>
              <div class="form-hint">控制次要关键词的匹配逻辑（仅Selective开启时有效）</div>
            </div>
            <div class="form-group compact">
              <label>角色 (Role)</label>
              <select v-model.number="entryForm.role">
                <option :value="0">System</option>
                <option :value="1">User</option>
                <option :value="2">Assistant</option>
              </select>
              <div class="form-hint">注入时以什么角色身份添加到消息列表</div>
            </div>
            <div class="form-group compact">
              <label>排序权重 (Order)</label>
              <input type="number" v-model.number="entryForm.order" min="0" max="999" />
            </div>
            <div class="form-group compact">
              <label>注入位置 (Position)</label>
              <select v-model.number="entryForm.position">
                <option :value="0">系统提示词之前</option>
                <option :value="1">系统提示词之后</option>
                <option :value="2">对话历史之前</option>
                <option :value="3">对话历史之后</option>
                <option :value="4">按深度插入 (D#)</option>
              </select>
            </div>
            <div class="form-group compact">
              <label>插入深度 (Depth)</label>
              <input type="number" v-model.number="entryForm.depth" min="0" max="999" />
              <div class="form-hint">从最新消息往前数第几条消息处插入（position=4时有效）</div>
            </div>
            <div class="form-group compact">
              <label>扫描深度 (Scan Depth)</label>
              <input type="number" v-model.number="entryForm.scanDepth" min="0" max="999" placeholder="留空=扫描全部" />
              <div class="form-hint">仅扫描最近N条消息进行关键词匹配（0或留空=全部扫描）</div>
            </div>
            <div class="form-group compact">
              <label>触发概率 (%) </label>
              <input type="number" v-model.number="entryForm.probability" min="0" max="100" />
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.useProbability" />
                <span>启用概率判定</span>
              </label>
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.caseSensitive" />
                <span>区分大小写 (Case Sensitive)</span>
              </label>
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.matchWholeWords" />
                <span>全词匹配 (Match Whole Words)</span>
              </label>
            </div>
            <div class="form-row">
              <label class="form-check">
                <input type="checkbox" v-model="entryForm.excludeRecursion" />
                <span>排除递归扫描</span>
              </label>
            </div>
          </template>
        </div>
        <div class="editor-footer">
          <button class="btn-cancel" @click="showEntryForm = false">取消</button>
          <button class="btn-save" @click="saveEntry">▣ 保存</button>
        </div>
      </div>
    </div>

    <!-- AI批量创建弹窗 -->
    <div v-if="showAIBatchModal" class="modal-overlay" @click.self="closeAIBatchModal">
      <div class="ai-batch-panel">
        <div class="editor-header">
          <h3>AI批量创建条目</h3>
          <span class="close-btn" @click="closeAIBatchModal">✕</span>
        </div>
        <div class="editor-body">
          <div class="form-group">
            <label>主题描述</label>
            <textarea
              v-model="aiBatchPrompt"
              rows="4"
              placeholder="描述你想创建的世界书条目主题，例如：&#10;一个奇幻魔法世界的设定，包含各种魔法体系、种族、地理区域等&#10;&#10;或者：&#10;一个赛博朋克城市的背景设定，包含各势力、科技、街区等"
            ></textarea>
          </div>
          <div class="form-group">
            <label>生成数量</label>
            <div class="count-selector">
              <button class="count-btn" @click="aiBatchCount = Math.max(1, aiBatchCount - 1)">-</button>
              <span class="count-value">{{ aiBatchCount }}</span>
              <button class="count-btn" @click="aiBatchCount = Math.min(15, aiBatchCount + 1)">+</button>
            </div>
          </div>
          <div v-if="aiBatchError" class="ai-batch-error">{{ aiBatchError }}</div>
        </div>
        <div class="editor-footer">
          <button class="btn-cancel" @click="closeAIBatchModal" :disabled="aiBatchLoading">取消</button>
          <button class="btn-save" @click="generateBatchEntries" :disabled="aiBatchLoading || !aiBatchPrompt.trim()">
            <span v-if="aiBatchLoading" class="spinner"></span>
            <span v-else>生成条目</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 世界书菜单弹窗 -->
    <div v-if="showBookMenu" class="modal-overlay" @click.self="showBookMenu = false">
      <div class="menu-panel">
        <div class="menu-header">
          <h3>世界书管理</h3>
          <span class="close-btn" @click="showBookMenu = false">✕</span>
        </div>
        <div class="menu-body">
          <div class="menu-item" @click="createNewBook">
            <span class="menu-icon">▤</span>
            <span>新建世界书</span>
          </div>
          <div class="menu-item" @click="triggerImport; showBookMenu = false">
            <span class="menu-icon">↓</span>
            <span>导入世界书</span>
          </div>
          <div class="menu-divider" v-if="worldBooks.length > 0"></div>
          <div class="menu-section-title" v-if="worldBooks.length > 0">已有世界书</div>
          <div
            v-for="book in worldBooks"
            :key="book.id"
            class="menu-item"
            :class="{ active: currentBookId === book.id }"
            @click="currentBookId = book.id; showBookMenu = false"
          >
            <span class="menu-icon">▤</span>
            <div class="menu-item-info">
              <span>{{ book.name }}</span>
              <span class="menu-item-meta">{{ book.entries.length }} 个条目</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { sendAIRequest } from '@/utils/aiService'
import { parseAIJsonArray } from '@/utils/aiJsonParser'
import { useSettingsStore } from '@/stores/settings'
import { worldBookApi } from '@/api/services'

interface WorldBookEntry {
  id: string
  title: string
  keywords: string[]
  content: string
  enabled: boolean
  // SillyTavern advanced fields
  keysecondary: string[]
  constant: boolean
  selective: boolean
  selectiveLogic: number    // 0=AND_ANY, 1=NOT_ALL, 2=NOT_ANY, 3=AND_ALL
  role: number              // 0=system, 1=user, 2=assistant
  scanDepth: number | null  // null = scan all messages
  caseSensitive: boolean
  matchWholeWords: boolean
  order: number
  position: number
  depth: number
  probability: number
  useProbability: boolean
  excludeRecursion: boolean
}

interface WorldBook {
  id: string
  name: string
  entries: WorldBookEntry[]
  bindChars: string[]  // 绑定的角色ID列表
  createdAt: string
}

const worldBooks = ref<WorldBook[]>([])
const currentBookId = ref<string | null>(null)
const searchText = ref('')
const expandedEntryId = ref<string | null>(null)
const showEntryForm = ref(false)
const showBookMenu = ref(false)
const batchMode = ref(false)
const selectedBookIds = ref<string[]>([])
const editingEntry = ref<WorldBookEntry | null>(null)
const keywordsInput = ref('')
const importInput = ref<HTMLInputElement | null>(null)
const showAdvanced = ref(false)
const secondaryKeysInput = ref('')

// AI批量创建状态
const showAIBatchModal = ref(false)
const aiBatchPrompt = ref('')
const aiBatchCount = ref(5)
const aiBatchLoading = ref(false)
const aiBatchError = ref('')
const settingsStore = useSettingsStore()

const entryForm = ref({
  title: '',
  keywords: [] as string[],
  content: '',
  enabled: true,
  // SillyTavern advanced fields
  constant: false,
  selective: false,
  selectiveLogic: 0,
  role: 0,
  scanDepth: null as number | null,
  caseSensitive: false,
  matchWholeWords: false,
  order: 0,
  position: 0,
  depth: 4,
  probability: 100,
  useProbability: true,
  excludeRecursion: false,
})

// 条目状态指示灯
function entryStatusClass(entry: WorldBookEntry): string {
  if (!entry.enabled) return 'status-disabled'      // 红灯：禁用
  if (entry.constant) return 'status-constant'       // 蓝灯：常驻
  if (entry.selective) return 'status-selective'      // 黄灯：选择性
  return 'status-normal'                              // 绿灯：普通
}

function entryStatusLabel(entry: WorldBookEntry): string {
  if (!entry.enabled) return '禁用'
  if (entry.constant) return '常驻'
  if (entry.selective) return '选择性'
  return '正常'
}

function parseStPosition(raw: unknown): number {
  if (typeof raw === 'number' && Number.isFinite(raw)) return raw
  if (typeof raw === 'string') {
    const n = Number(raw)
    if (Number.isFinite(n)) return n
    const key = raw.trim().toLowerCase()
    if (key === 'before_char' || key === 'before_prompt' || key === 'before_an' || key === 'before_system') return 0
    if (key === 'after_char' || key === 'after_prompt' || key === 'after_an' || key === 'after_system') return 1
    if (key === 'before_history' || key === 'before_chat' || key === 'before_messages') return 2
    if (key === 'after_history' || key === 'after_chat' || key === 'after_messages') return 3
    if (key === 'at_depth' || key === 'in_chat' || key === 'chat' || key === 'depth') return 4
  }
  return 0
}

const currentBook = computed(() => {
  if (!currentBookId.value) return null
  return worldBooks.value.find(b => b.id === currentBookId.value) || null
})

const allBooksSelected = computed(() => worldBooks.value.length > 0 && worldBooks.value.every(book => selectedBookIds.value.includes(String(book.id))))

const filteredEntries = computed(() => {
  if (!currentBook.value) return []
  let entries = currentBook.value.entries
  if (searchText.value.trim()) {
    const kw = searchText.value.trim().toLowerCase()
    entries = entries.filter(e =>
      e.title.toLowerCase().includes(kw) ||
      e.content.toLowerCase().includes(kw) ||
      e.keywords.some(k => k.toLowerCase().includes(kw))
    )
  }
  return entries
})

async function loadWorldBooks() {
  try {
    const res: any = await worldBookApi.list()
    const list = Array.isArray(res.data) ? res.data : []
    if (list.length === 0) {
      try {
        const saved = localStorage.getItem('worldBooks')
        const localBooks = saved ? JSON.parse(saved) : []
        if (Array.isArray(localBooks) && localBooks.length > 0) {
          for (const localBook of localBooks) {
            await worldBookApi.create({
              name: localBook.name || '未命名世界书',
              bind_chars: Array.isArray(localBook.bindChars) ? localBook.bindChars : [],
              entries: Array.isArray(localBook.entries) ? localBook.entries : [],
            } as any)
          }
          const migrated = await worldBookApi.list()
          const migratedList = Array.isArray(migrated.data) ? migrated.data : []
          worldBooks.value = migratedList.map((book: any) => ({
            id: String(book.id),
            name: book.name || '未命名世界书',
            entries: Array.isArray(book.entries) ? book.entries : [],
            bindChars: Array.isArray(book.bind_chars) ? book.bind_chars : [],
            createdAt: book.created_at || new Date().toISOString(),
          })) as any
          currentBookId.value = worldBooks.value[0]?.id || null
          return
        }
      } catch {
        // ignore local migration failure
      }
    }
    worldBooks.value = list.map((book: any) => ({
      id: String(book.id),
      name: book.name || '未命名世界书',
      entries: Array.isArray(book.entries) ? book.entries.map((entry: any, idx: number) => ({
        id: String(entry.id || `entry-${idx}`),
        title: entry.title || entry.key || `条目${idx + 1}`,
        keywords: Array.isArray(entry.keywords) ? entry.keywords : [],
        content: entry.content || '',
        enabled: entry.enabled !== false && entry.is_enabled !== false,
        keysecondary: Array.isArray(entry.keysecondary) ? entry.keysecondary : [],
        constant: !!entry.constant,
        selective: !!entry.selective,
        selectiveLogic: typeof entry.selectiveLogic === 'number' ? entry.selectiveLogic : 0,
        role: typeof entry.role === 'number' ? entry.role : 0,
        scanDepth: typeof entry.scanDepth === 'number' ? entry.scanDepth : null,
        caseSensitive: !!entry.caseSensitive,
        matchWholeWords: !!entry.matchWholeWords,
        order: typeof entry.order === 'number' ? entry.order : idx,
        position: typeof entry.position === 'number' ? entry.position : 0,
        depth: typeof entry.depth === 'number' ? entry.depth : 4,
        probability: typeof entry.probability === 'number' ? entry.probability : 100,
        useProbability: entry.useProbability !== false,
        excludeRecursion: !!entry.excludeRecursion,
      })) : [],
      bindChars: Array.isArray(book.bind_chars) ? book.bind_chars : [],
      createdAt: book.created_at || new Date().toISOString(),
    }))
    currentBookId.value = worldBooks.value[0]?.id || null
  } catch (err) {
    console.error('Failed to load worldbooks', err)
  }
}

async function persistBook(book: WorldBook) {
  const numericId = Number(book.id)
  if (!Number.isFinite(numericId)) return
  await worldBookApi.update(numericId, {
    name: book.name,
    bind_chars: book.bindChars,
    entries: book.entries,
  } as any)
}

// 世界书操作
async function createNewBook() {
  const name = prompt('输入世界书名称：', '新世界书')
  if (!name?.trim()) return

  const book: WorldBook = {
    id: `wb-${Date.now()}`,
    name: name.trim(),
    entries: [],
    bindChars: [],
    createdAt: new Date().toISOString(),
  }
  try {
    const res: any = await worldBookApi.create({ name: book.name, bind_chars: [], entries: [] } as any)
    if (res.id) book.id = String(res.id)
  } catch (err) {
    console.error('Failed to create worldbook', err)
  }
  worldBooks.value.unshift(book)
  currentBookId.value = book.id
  showBookMenu.value = false
}

function toggleBatchMode() {
  batchMode.value = !batchMode.value
  if (!batchMode.value) selectedBookIds.value = []
}

function toggleBookSelection(id: string | number) {
  const key = String(id)
  if (selectedBookIds.value.includes(key)) {
    selectedBookIds.value = selectedBookIds.value.filter(item => item !== key)
    return
  }
  selectedBookIds.value.push(key)
}

function toggleSelectAllBooks() {
  if (allBooksSelected.value) {
    selectedBookIds.value = []
    return
  }
  selectedBookIds.value = worldBooks.value.map(book => String(book.id))
}

async function batchDeleteBooks() {
  if (selectedBookIds.value.length === 0) return
  if (!confirm(`确定要删除选中的 ${selectedBookIds.value.length} 本世界书吗？`)) return
  for (const id of selectedBookIds.value) {
    const numericId = Number(id)
    if (Number.isFinite(numericId)) {
      await worldBookApi.delete(numericId)
    }
  }
  selectedBookIds.value = []
  batchMode.value = false
  await loadWorldBooks()
}

function batchExportBooks() {
  const selected = worldBooks.value.filter(book => selectedBookIds.value.includes(String(book.id)))
  if (selected.length === 0) return
  const blob = new Blob([JSON.stringify(selected, null, 2)], { type: 'application/json;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `worldbooks-batch-${Date.now()}.json`
  link.click()
  URL.revokeObjectURL(url)
}

async function editBookName() {
  if (!currentBook.value) return
  const name = prompt('修改世界书名称：', currentBook.value.name)
  if (!name?.trim()) return
  currentBook.value.name = name.trim()
  await persistBook(currentBook.value)
}

async function deleteBook() {
  if (!currentBook.value) return
  if (!confirm(`确定要删除世界书「${currentBook.value.name}」吗？其中的所有条目都会被删除。`)) return

  const deletingId = Number(currentBook.value.id)
  if (Number.isFinite(deletingId)) {
    try {
      await worldBookApi.delete(deletingId)
    } catch (err) {
      console.error('Failed to delete worldbook', err)
    }
  }

  const idx = worldBooks.value.findIndex(b => b.id === currentBookId.value)
  worldBooks.value.splice(idx, 1)

  if (worldBooks.value.length > 0) {
    currentBookId.value = worldBooks.value[0].id
  } else {
    currentBookId.value = null
  }
}

function exportBook() {
  if (!currentBook.value) return

  // 导出为 SillyTavern 兼容格式（使用实际存储的值）
  const exportData = {
    name: currentBook.value.name,
    entries: currentBook.value.entries.reduce((acc: any, entry: WorldBookEntry, idx: number) => {
      acc[idx] = {
        uid: idx,
        key: entry.keywords,
        keysecondary: entry.keysecondary || [],
        secondary_keys: entry.keysecondary || [],
        comment: entry.title,
        content: entry.content,
        constant: entry.constant ?? false,
        selective: entry.selective ?? false,
        selectiveLogic: entry.selectiveLogic ?? 0,
        insertion_order: entry.order ?? idx,
        order: entry.order ?? idx,
        position: entry.position ?? 0,
        disable: !entry.enabled,
        enabled: entry.enabled !== false,
        excludeRecursion: entry.excludeRecursion ?? false,
        probability: entry.probability ?? 100,
        useProbability: entry.useProbability !== false,
        depth: entry.depth ?? 4,
        role: entry.role ?? 0,
        scanDepth: entry.scanDepth ?? null,
        caseSensitive: entry.caseSensitive ?? false,
        matchWholeWords: entry.matchWholeWords ?? false,
        extensions: {
          position: entry.position ?? 0,
          depth: entry.depth ?? 4,
          probability: entry.probability ?? 100,
          useProbability: entry.useProbability !== false,
          exclude_recursion: entry.excludeRecursion ?? false,
          role: entry.role ?? 0,
          scan_depth: entry.scanDepth ?? null,
          case_sensitive: entry.caseSensitive ?? false,
          match_whole_words: entry.matchWholeWords ?? false,
          selectiveLogic: entry.selectiveLogic ?? 0,
        },
      }
      return acc
    }, {}),
  }

  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `worldbook-${currentBook.value.name}.json`
  a.click()
  URL.revokeObjectURL(url)
}

// 条目操作
function toggleEntryExpand(id: string) {
  expandedEntryId.value = expandedEntryId.value === id ? null : id
}

function openEntryEditor(entry: WorldBookEntry | null) {
  editingEntry.value = entry
  if (entry) {
    entryForm.value = {
      title: entry.title,
      keywords: [...entry.keywords],
      content: entry.content,
      enabled: entry.enabled,
      constant: entry.constant ?? false,
      selective: entry.selective ?? false,
      selectiveLogic: entry.selectiveLogic ?? 0,
      role: entry.role ?? 0,
      scanDepth: entry.scanDepth ?? null,
      caseSensitive: entry.caseSensitive ?? false,
      matchWholeWords: entry.matchWholeWords ?? false,
      order: entry.order ?? 0,
      position: entry.position ?? 0,
      depth: entry.depth ?? 4,
      probability: entry.probability ?? 100,
      useProbability: entry.useProbability !== false,
      excludeRecursion: entry.excludeRecursion ?? false,
    }
    keywordsInput.value = entry.keywords.join(', ')
    secondaryKeysInput.value = (entry.keysecondary || []).join(', ')
  } else {
    entryForm.value = {
      title: '',
      keywords: [],
      content: '',
      enabled: true,
      constant: false,
      selective: false,
      selectiveLogic: 0,
      role: 0,
      scanDepth: null,
      caseSensitive: false,
      matchWholeWords: false,
      order: 0,
      position: 0,
      depth: 4,
      probability: 100,
      useProbability: true,
      excludeRecursion: false,
    }
    keywordsInput.value = ''
    secondaryKeysInput.value = ''
  }
  showAdvanced.value = false
  showEntryForm.value = true
}

async function saveEntry() {
  if (!currentBook.value) return
  if (!entryForm.value.title.trim()) return

  const keywords = keywordsInput.value
    .split(/[,，]/)
    .map(k => k.trim())
    .filter(k => k.length > 0)

  const keysecondary = secondaryKeysInput.value
    .split(/[,，]/)
    .map(k => k.trim())
    .filter(k => k.length > 0)

  const entryData: Omit<WorldBookEntry, 'id'> = {
    title: entryForm.value.title.trim(),
    keywords,
    content: entryForm.value.content,
    enabled: entryForm.value.enabled,
    keysecondary,
    constant: entryForm.value.constant,
    selective: entryForm.value.selective,
    selectiveLogic: entryForm.value.selectiveLogic,
    role: entryForm.value.role,
    scanDepth: entryForm.value.scanDepth,
    caseSensitive: entryForm.value.caseSensitive,
    matchWholeWords: entryForm.value.matchWholeWords,
    order: entryForm.value.order,
    position: entryForm.value.position,
    depth: entryForm.value.depth,
    probability: entryForm.value.probability,
    useProbability: entryForm.value.useProbability,
    excludeRecursion: entryForm.value.excludeRecursion,
  }

  if (editingEntry.value) {
    // 更新
    const idx = currentBook.value.entries.findIndex(e => e.id === editingEntry.value!.id)
    if (idx >= 0) {
      currentBook.value.entries[idx] = {
        ...currentBook.value.entries[idx],
        ...entryData,
      }
    }
  } else {
    // 新建
    currentBook.value.entries.push({
      id: `entry-${Date.now()}`,
      ...entryData,
    })
  }

  await persistBook(currentBook.value)
  showEntryForm.value = false
}

async function deleteEntry(id: string) {
  if (!currentBook.value) return
  if (!confirm('确定要删除此条目？')) return
  currentBook.value.entries = currentBook.value.entries.filter(e => e.id !== id)
  if (expandedEntryId.value === id) expandedEntryId.value = null
  await persistBook(currentBook.value)
}

async function toggleEntryEnabled(entry: WorldBookEntry) {
  entry.enabled = !entry.enabled
  if (currentBook.value) {
    await persistBook(currentBook.value)
  }
}

// 导入
function triggerImport() {
  importInput.value?.click()
}

function handleImport(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async () => {
    try {
      const data = JSON.parse(reader.result as string)
      const bookName = data.name || file.name.replace(/\.json$/i, '')

      const entries: WorldBookEntry[] = []

      // 支持 SillyTavern 格式（保留所有高级字段）
      if (data.entries && typeof data.entries === 'object') {
        const entryObj = data.entries
        for (const key of Object.keys(entryObj)) {
          const e = entryObj[key]
          const ext = e.extensions || {}
          entries.push({
            id: `entry-${Date.now()}-${key}`,
            title: e.comment || e.title || `条目${key}`,
            keywords: Array.isArray(e.key)
              ? e.key
              : (Array.isArray(e.keys)
                ? e.keys
                : (Array.isArray(e.keywords)
                  ? e.keywords
                  : (typeof e.key === 'string'
                    ? e.key.split(',').map((k: string) => k.trim())
                    : []))),
            content: e.content || '',
            enabled: typeof e.enabled === 'boolean' ? e.enabled : !e.disable,
            keysecondary: Array.isArray(e.keysecondary)
              ? e.keysecondary
              : (Array.isArray(e.secondary_keys)
                ? e.secondary_keys
                : (typeof e.keysecondary === 'string'
                  ? e.keysecondary.split(',').map((k: string) => k.trim()).filter(Boolean)
                  : [])),
            constant: e.constant ?? false,
            selective: e.selective ?? true,
            selectiveLogic: typeof ext.selectiveLogic === 'number'
              ? ext.selectiveLogic
              : (typeof e.selectiveLogic === 'number' ? e.selectiveLogic : 0),
            role: typeof ext.role === 'number'
              ? ext.role
              : (typeof e.role === 'number' ? e.role : 0),
            scanDepth: typeof ext.scan_depth === 'number'
              ? ext.scan_depth
              : (typeof e.scanDepth === 'number' ? e.scanDepth : null),
            caseSensitive: typeof ext.case_sensitive === 'boolean'
              ? ext.case_sensitive
              : !!e.caseSensitive,
            matchWholeWords: typeof ext.match_whole_words === 'boolean'
              ? ext.match_whole_words
              : !!e.matchWholeWords,
            order: typeof e.insertion_order === 'number'
              ? e.insertion_order
              : (typeof e.order === 'number' ? e.order : parseInt(key) || 0),
            position: parseStPosition(ext.position ?? e.position),
            depth: typeof ext.depth === 'number'
              ? ext.depth
              : (typeof e.depth === 'number' ? e.depth : 4),
            probability: typeof ext.probability === 'number'
              ? ext.probability
              : (typeof e.probability === 'number' ? e.probability : 100),
            useProbability: typeof ext.useProbability === 'boolean'
              ? ext.useProbability
              : (e.useProbability !== false),
            excludeRecursion: typeof ext.exclude_recursion === 'boolean'
              ? ext.exclude_recursion
              : !!e.excludeRecursion,
          })
        }
      }
      // 支持我们自己的格式（直接数组）
      else if (Array.isArray(data)) {
        for (const e of data) {
          entries.push({
            id: `entry-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
            title: e.title || e.key || '未命名',
            keywords: Array.isArray(e.keywords) ? e.keywords : [],
            content: e.content || '',
            enabled: e.enabled !== false,
            keysecondary: Array.isArray(e.keysecondary) ? e.keysecondary : [],
            constant: e.constant ?? false,
            selective: e.selective ?? false,
            selectiveLogic: typeof e.selectiveLogic === 'number' ? e.selectiveLogic : 0,
            role: typeof e.role === 'number' ? e.role : 0,
            scanDepth: typeof e.scanDepth === 'number' ? e.scanDepth : null,
            caseSensitive: e.caseSensitive ?? false,
            matchWholeWords: e.matchWholeWords ?? false,
            order: e.order ?? 0,
            position: parseStPosition(e.position),
            depth: e.depth ?? 4,
            probability: e.probability ?? 100,
            useProbability: e.useProbability !== false,
            excludeRecursion: e.excludeRecursion ?? false,
          })
        }
      }

      const newBook: WorldBook = {
        id: `wb-${Date.now()}`,
        name: bookName,
        entries,
        bindChars: [],
        createdAt: new Date().toISOString(),
      }

      try {
        const res: any = await worldBookApi.create({ name: newBook.name, bind_chars: [], entries } as any)
        if (res.id) newBook.id = String(res.id)
      } catch (err) {
        console.error('Failed to import worldbook to backend', err)
      }

      worldBooks.value.unshift(newBook)
      currentBookId.value = newBook.id

      alert(`成功导入世界书「${bookName}」，共 ${entries.length} 个条目`)
    } catch {
      alert('导入失败：文件格式不正确')
    }
  }
  reader.readAsText(file)
  if (importInput.value) importInput.value.value = ''
}

// AI批量创建
function closeAIBatchModal() {
  if (aiBatchLoading.value) return
  showAIBatchModal.value = false
  aiBatchPrompt.value = ''
  aiBatchError.value = ''
  aiBatchCount.value = 5
}

async function generateBatchEntries() {
  if (!currentBook.value) {
    aiBatchError.value = '请先选择或创建一本世界书'
    return
  }
  if (!aiBatchPrompt.value.trim()) return

  aiBatchLoading.value = true
  aiBatchError.value = ''

  try {
    const systemPrompt = `你是一个世界观设定创作助手。用户会给你一个主题描述，你需要根据描述生成${aiBatchCount.value}个世界书条目。

每个条目包含：
- title: 条目标题（简洁明了）
- keywords: 触发关键词数组（3-6个相关关键词，当聊天中出现这些词时会触发该条目）
- content: 条目内容（详细的背景设定描述，150-400字）

请直接输出JSON数组格式，不要包含其他文字说明。
示例格式：
[
  {
    "title": "火焰魔法",
    "keywords": ["火焰", "火球", "燃烧", "火系魔法", "火焰术"],
    "content": "火焰魔法是元素魔法中最常见的一种..."
  }
]`

    const result = await sendAIRequest({
      apiKey: settingsStore.settings.apiSource === 'platform' ? '' : settingsStore.settings.apiKey,
      apiUrl: settingsStore.getApiUrl(),
      providerId: settingsStore.getPlatformProviderId(),
      model: settingsStore.settings.model,
      messages: [
        { role: 'system', content: systemPrompt },
        { role: 'user', content: `主题：${aiBatchPrompt.value.trim()}\n\n请生成${aiBatchCount.value}个世界书条目。` }
      ],
      temperature: 0.8,
      maxTokens: 4000,
      stream: false,
      timeout: settingsStore.settings.timeout || 120,
    })

    if (!result.content) {
      throw new Error('AI返回内容为空')
    }

    let content = result.content.trim()
    // 去除可能的markdown代码块包裹
    if (content.startsWith('```')) {
      content = content.replace(/^```(?:json)?\s*\n?/, '').replace(/\n?```\s*$/, '')
    }

    const parsed = parseAIJsonArray(content, { allowSingleObject: false })
    if (!Array.isArray(parsed) || parsed.length === 0) {
      throw new Error('AI返回的数据格式不正确')
    }

    const newEntries: WorldBookEntry[] = parsed.map((item: any, idx: number) => ({
      id: `entry-${Date.now()}-${idx}`,
      title: String(item.title || `条目${idx + 1}`),
      keywords: Array.isArray(item.keywords) ? item.keywords.map(String) : [],
      content: String(item.content || ''),
      enabled: true,
      keysecondary: [],
      constant: false,
      selective: false,
      selectiveLogic: 0,
      role: 0,
      scanDepth: null,
      caseSensitive: false,
      matchWholeWords: false,
      order: idx,
      position: 0,
      depth: 4,
      probability: 100,
      useProbability: true,
      excludeRecursion: false,
    }))

    currentBook.value.entries.push(...newEntries)
    await persistBook(currentBook.value)
    closeAIBatchModal()
  } catch (err: any) {
    aiBatchError.value = err.message || '生成失败，请重试'
  } finally {
    aiBatchLoading.value = false
  }
}

onMounted(() => {
  loadWorldBooks()
})
</script>

<style scoped>
.worldbook-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

/* 世界书标签 */
.book-tabs {
  display: flex;
  gap: 6px;
  padding: 8px 12px;
  overflow-x: auto;
  background: var(--bg-primary);
  border-bottom: 0.5px solid var(--separator);
}

.book-tabs::-webkit-scrollbar { display: none; }

.batch-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  gap: 12px;
  border-bottom: 0.5px solid var(--separator);
  background: var(--bg-secondary, #f7f7fa);
}

.batch-summary {
  font-size: 13px;
  color: var(--text-secondary);
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.batch-btn {
  border: none;
  border-radius: 8px;
  padding: 8px 12px;
  background: var(--fill-tertiary, rgba(118,118,128,0.12));
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
}

.batch-btn.danger {
  background: rgba(255, 59, 48, 0.12);
  color: #ff3b30;
}

.batch-check {
  display: flex;
  align-items: center;
}

.book-tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 16px;
  background: var(--bg-tertiary);
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
}

.book-tab.active {
  background: var(--color-primary, #007aff);
  color: #fff;
}

.book-tab.selected {
  box-shadow: inset 0 0 0 2px rgba(0, 122, 255, 0.35);
  background: rgba(0, 122, 255, 0.12);
}

.book-tab-icon { font-size: 14px; }
.book-tab-name { font-size: 13px; font-weight: 500; }

.book-tab-count {
  font-size: 11px;
  padding: 0 5px;
  border-radius: 8px;
  background: rgba(0,0,0,0.1);
  min-width: 16px;
  text-align: center;
}

.book-tab.active .book-tab-count {
  background: rgba(255,255,255,0.25);
}

.add-tab {
  width: 32px;
  justify-content: center;
  font-size: 18px;
  color: var(--text-tertiary);
}

/* 内容区 */
.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

/* 世界书信息栏 */
.book-info-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
  margin-bottom: 10px;
}

.book-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.book-meta {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.book-info-actions {
  display: flex;
  gap: 4px;
}

.small-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: var(--bg-tertiary);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.small-btn.danger { background: rgba(255, 59, 48, 0.1); }
.small-btn:active { opacity: 0.7; }

/* 工具栏 */
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
}

.search-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-tertiary);
  border-radius: 10px;
  padding: 8px 10px;
}

.search-wrap input {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  font-size: 14px;
  color: var(--text-primary);
}

.search-icon { font-size: 13px; opacity: 0.5; }

.add-entry-btn {
  padding: 8px 14px;
  border-radius: 10px;
  border: none;
  background: var(--color-primary, #007aff);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  flex-shrink: 0;
}

/* 条目卡片 */
.entry-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.entry-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  transition: opacity 0.2s;
}

.entry-card.disabled { opacity: 0.5; }

.entry-header {
  padding: 12px 14px;
  cursor: pointer;
}

.entry-title-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.entry-icon { font-size: 14px; }

.entry-title {
  flex: 1;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 蓝绿灯状态指示 */
.entry-status-light {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 4px currentColor;
}

.entry-status-light.status-normal {
  background: #34c759;
  color: #34c759;
}

.entry-status-light.status-constant {
  background: #007aff;
  color: #007aff;
}

.entry-status-light.status-selective {
  background: #ff9500;
  color: #ff9500;
}

.entry-status-light.status-disabled {
  background: #ff3b30;
  color: #ff3b30;
  opacity: 0.6;
}

.entry-status-tag {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.entry-keywords-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.kw-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 6px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.kw-more {
  font-size: 11px;
  color: var(--text-tertiary);
  padding: 2px 4px;
}

/* 展开区 */
.entry-expand {
  padding: 0 14px 12px;
  border-top: 1px solid var(--separator);
}

.entry-content-preview {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin-top: 8px;
  padding: 8px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  white-space: pre-wrap;
  max-height: 150px;
  overflow-y: auto;
}

.entry-actions {
  display: flex;
  gap: 6px;
  margin-top: 10px;
}

.act-btn {
  flex: 1;
  padding: 8px;
  border-radius: 8px;
  border: none;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
}

.act-btn.danger { background: rgba(255, 59, 48, 0.1); color: #ff3b30; }
.act-btn:active { opacity: 0.7; }

/* 空状态 */
.empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-tertiary);
}

.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-hint { font-size: 13px; margin-top: 4px; opacity: 0.7; }

.empty-btn {
  margin-top: 12px;
  padding: 8px 20px;
  border-radius: 18px;
  border: none;
  background: var(--color-primary, #007aff);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

/* 底部 */
.bottom-bar {
  padding: 8px 12px 16px;
}

.import-btn {
  width: 100%;
  padding: 10px;
  border-radius: 12px;
  border: 1px dashed var(--border-primary, var(--separator));
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
}

.import-btn:active { background: var(--bg-tertiary); }

/* 弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.editor-panel, .menu-panel {
  width: 100%;
  max-width: 393px;
  max-height: 85%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
}

.editor-header, .menu-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid var(--separator);
}

.editor-header h3, .menu-header h3 {
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
  border: none;
}

.editor-body, .menu-body {
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
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--separator);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  resize: vertical;
  box-sizing: border-box;
  font-family: inherit;
}

.form-hint {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 4px;
  line-height: 1.4;
}

.form-row {
  margin-bottom: 12px;
}

.form-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}

.editor-footer {
  display: flex;
  gap: 10px;
  padding: 12px 16px;
  border-top: 1px solid var(--separator);
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
.btn-save { background: var(--color-primary, #007aff); color: #fff; }

/* 菜单 */
.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 4px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.15s;
}

.menu-item:active { background: var(--bg-tertiary); }
.menu-item.active { background: rgba(0, 122, 255, 0.08); }

.menu-icon { font-size: 20px; flex-shrink: 0; }

.menu-item-info {
  display: flex;
  flex-direction: column;
  font-size: 15px;
  color: var(--text-primary);
}

.menu-item-meta {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 1px;
}

.menu-divider {
  height: 0.5px;
  background: var(--separator);
  margin: 8px 0;
}

.menu-section-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-tertiary);
  padding: 4px 4px 8px;
  text-transform: uppercase;
}

.icon-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  color: var(--brand-primary, var(--color-primary, #007aff));
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-btn:active { opacity: 0.5; }
.icon-btn svg { width: 22px; height: 22px; }

/* AI批量创建按钮 */
.ai-batch-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 10px;
  border: none;
  background: var(--bg-tertiary);
  color: var(--color-primary, #007aff);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  flex-shrink: 0;
}

.ai-batch-btn:active { opacity: 0.7; }

.ai-batch-btn svg {
  flex-shrink: 0;
}

/* AI批量创建弹窗 */
.ai-batch-panel {
  width: 100%;
  max-width: 393px;
  max-height: 85%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
}

.count-selector {
  display: flex;
  align-items: center;
  gap: 16px;
}

.count-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: 1px solid var(--separator);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.count-btn:active { opacity: 0.6; }

.count-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  min-width: 24px;
  text-align: center;
}

.ai-batch-error {
  padding: 10px 12px;
  background: rgba(255, 59, 48, 0.1);
  color: #ff3b30;
  border-radius: 10px;
  font-size: 13px;
  margin-top: 8px;
}

.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

/* 高级设置 */
.advanced-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  margin-bottom: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-primary, #007aff);
  border-top: 1px solid var(--separator);
  user-select: none;
}

.toggle-arrow {
  font-size: 16px;
  font-weight: 700;
  transition: transform 0.2s;
}

.toggle-arrow.expanded {
  transform: rotate(90deg);
}

.form-group.compact {
  margin-bottom: 8px;
}

.form-group.compact label {
  font-size: 12px;
  margin-bottom: 2px;
}

.form-group.compact input[type="number"] {
  width: 100px;
}

.form-hint.inline-hint {
  margin-top: 0;
  margin-left: 26px;
  margin-bottom: 4px;
}

.form-group select {
  width: 100%;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid var(--separator);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
  font-family: inherit;
  -webkit-appearance: none;
  appearance: none;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
