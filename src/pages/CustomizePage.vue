<template>
  <div class="customize-page">
    <NavBar title="自定义设置" back />

    <div class="settings-sections">
      <!-- 外观设置 -->
      <div class="section">
        <div class="section-header">🎨 外观设置</div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">深色模式</span>
            <span class="label-desc">开启后界面切换为深色主题</span>
          </div>
          <div class="toggle" :class="{ on: darkMode }" @click="darkMode = !darkMode">
            <div class="toggle-thumb"></div>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">主题色</span>
            <span class="label-desc">选择你喜欢的主题颜色</span>
          </div>
          <div class="color-options">
            <div
              v-for="color in themeColors"
              :key="color.value"
              class="color-dot"
              :style="{ background: color.value }"
              :class="{ active: selectedColor === color.value }"
              @click="selectedColor = color.value"
            >
              <span v-if="selectedColor === color.value">✓</span>
            </div>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">壁纸</span>
            <span class="label-desc">选择锁屏和主页壁纸</span>
          </div>
        </div>
        <div class="wallpaper-grid">
          <div
            v-for="wp in wallpapers"
            :key="wp.id"
            class="wallpaper-item"
            :class="{ active: selectedWallpaper === wp.id }"
            :style="{ background: wp.gradient }"
            @click="selectedWallpaper = wp.id"
          >
            <span class="wp-emoji">{{ wp.emoji }}</span>
            <span class="wp-name">{{ wp.name }}</span>
            <span v-if="selectedWallpaper === wp.id" class="wp-check">✓</span>
          </div>
        </div>
      </div>

      <!-- 字体设置 -->
      <div class="section">
        <div class="section-header">🔤 字体设置</div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">字体大小</span>
            <span class="label-desc">当前：{{ fontSizeLabel }}</span>
          </div>
        </div>
        <div class="slider-row">
          <span class="slider-label">A</span>
          <input
            type="range"
            min="0"
            max="4"
            :value="fontSizeIndex"
            @input="fontSizeIndex = Number(($event.target as HTMLInputElement).value)"
            class="slider"
          />
          <span class="slider-label large">A</span>
        </div>
        <div class="font-preview" :style="{ fontSize: fontSizes[fontSizeIndex].size + 'px' }">
          这是预览文本 Preview Text
        </div>
      </div>

      <!-- 通知设置 -->
      <div class="section">
        <div class="section-header">🔔 通知设置</div>

        <div class="setting-item" v-for="notif in notifications" :key="notif.key">
          <div class="setting-label">
            <span class="label-text">{{ notif.label }}</span>
            <span class="label-desc">{{ notif.desc }}</span>
          </div>
          <div class="toggle" :class="{ on: notif.enabled }" @click="notif.enabled = !notif.enabled">
            <div class="toggle-thumb"></div>
          </div>
        </div>
      </div>

      <!-- AI 模型设置 -->
      <div class="section">
        <div class="section-header">🤖 AI 模型设置</div>

        <div class="setting-item" @click="showModelPicker = true">
          <div class="setting-label">
            <span class="label-text">聊天模型</span>
            <span class="label-desc">{{ selectedModel }}</span>
          </div>
          <span class="arrow">›</span>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">回复创意度</span>
            <span class="label-desc">Temperature: {{ temperature.toFixed(1) }}</span>
          </div>
        </div>
        <div class="slider-row">
          <span class="slider-min">精准</span>
          <input
            type="range"
            min="0"
            max="20"
            :value="temperature * 10"
            @input="temperature = Number(($event.target as HTMLInputElement).value) / 10"
            class="slider"
          />
          <span class="slider-max">创意</span>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">最大回复长度</span>
            <span class="label-desc">{{ maxTokens }} tokens</span>
          </div>
        </div>
        <div class="slider-row">
          <span class="slider-min">短</span>
          <input
            type="range"
            min="1"
            max="40"
            :value="maxTokens / 100"
            @input="maxTokens = Number(($event.target as HTMLInputElement).value) * 100"
            class="slider"
          />
          <span class="slider-max">长</span>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">流式输出</span>
            <span class="label-desc">实时显示 AI 生成的文字</span>
          </div>
          <div class="toggle" :class="{ on: streamOutput }" @click="streamOutput = !streamOutput">
            <div class="toggle-thumb"></div>
          </div>
        </div>
      </div>

      <!-- 聊天设置 -->
      <div class="section">
        <div class="section-header">💬 聊天设置</div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">自动保存聊天</span>
            <span class="label-desc">每次对话自动保存到历史</span>
          </div>
          <div class="toggle" :class="{ on: autoSave }" @click="autoSave = !autoSave">
            <div class="toggle-thumb"></div>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">消息气泡样式</span>
            <span class="label-desc">{{ bubbleStyles[bubbleStyleIndex].name }}</span>
          </div>
          <span class="arrow" @click="bubbleStyleIndex = (bubbleStyleIndex + 1) % bubbleStyles.length">›</span>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">发送键行为</span>
            <span class="label-desc">{{ sendKeyBehavior === 'send' ? 'Enter 发送' : 'Enter 换行' }}</span>
          </div>
          <div class="toggle" :class="{ on: sendKeyBehavior === 'send' }" @click="sendKeyBehavior = sendKeyBehavior === 'send' ? 'newline' : 'send'">
            <div class="toggle-thumb"></div>
          </div>
        </div>
      </div>

      <!-- 存储与数据 -->
      <div class="section">
        <div class="section-header">💾 存储与数据</div>

        <div class="setting-item" @click="showClearConfirm = true">
          <div class="setting-label">
            <span class="label-text danger">清除所有数据</span>
            <span class="label-desc">清除本地缓存和聊天记录</span>
          </div>
          <span class="arrow">›</span>
        </div>

        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">导出数据</span>
            <span class="label-desc">导出聊天记录和设置</span>
          </div>
          <span class="arrow" @click="exportData">›</span>
        </div>
      </div>

      <!-- 保存按钮 -->
      <button class="save-btn" @click="saveSettings">
        💾 保存设置
      </button>
      <div v-if="saved" class="save-toast">✅ 设置已保存</div>
    </div>

    <!-- 模型选择器 -->
    <div v-if="showModelPicker" class="modal-overlay" @click.self="showModelPicker = false">
      <div class="picker-panel">
        <div class="picker-header">
          <h3>选择 AI 模型</h3>
          <span class="close-btn" @click="showModelPicker = false">✕</span>
        </div>
        <div class="model-list">
          <div
            v-for="model in models"
            :key="model.id"
            class="model-item"
            :class="{ active: selectedModel === model.id }"
            @click="selectedModel = model.id; showModelPicker = false"
          >
            <div class="model-info">
              <div class="model-name">{{ model.name }}</div>
              <div class="model-desc">{{ model.desc }}</div>
            </div>
            <span v-if="selectedModel === model.id" class="check-icon">✓</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 清除确认 -->
    <div v-if="showClearConfirm" class="modal-overlay" @click.self="showClearConfirm = false">
      <div class="confirm-dialog">
        <div class="confirm-icon">⚠️</div>
        <h3>确定要清除所有数据吗？</h3>
        <p>此操作不可恢复，将清除所有本地数据</p>
        <div class="confirm-actions">
          <button class="btn-cancel" @click="showClearConfirm = false">取消</button>
          <button class="btn-danger" @click="clearAllData">确认清除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

// 外观
const darkMode = ref(false)
const selectedColor = ref('#007aff')
const selectedWallpaper = ref('default')

const themeColors = [
  { value: '#007aff', name: '蓝色' },
  { value: '#ff3b30', name: '红色' },
  { value: '#34c759', name: '绿色' },
  { value: '#ff9500', name: '橙色' },
  { value: '#af52de', name: '紫色' },
  { value: '#ff2d55', name: '粉色' },
  { value: '#5856d6', name: '靛蓝' },
  { value: '#00c7be', name: '青色' },
]

const wallpapers = [
  { id: 'default', name: '默认', emoji: '🌊', gradient: 'linear-gradient(135deg, #667eea, #764ba2)' },
  { id: 'sunset', name: '日落', emoji: '🌅', gradient: 'linear-gradient(135deg, #fa709a, #fee140)' },
  { id: 'forest', name: '森林', emoji: '🌲', gradient: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
  { id: 'ocean', name: '海洋', emoji: '🐳', gradient: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
  { id: 'night', name: '夜空', emoji: '🌌', gradient: 'linear-gradient(135deg, #0c0c1d, #1a1a3e)' },
  { id: 'cherry', name: '樱花', emoji: '🌸', gradient: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
]

// 字体
const fontSizeIndex = ref(2)
const fontSizes = [
  { label: '极小', size: 12 },
  { label: '小', size: 13 },
  { label: '标准', size: 14 },
  { label: '大', size: 16 },
  { label: '极大', size: 18 },
]
const fontSizeLabel = ref('标准')

// 通知
const notifications = ref([
  { key: 'msg', label: '消息通知', desc: '新消息时显示推送通知', enabled: true },
  { key: 'sound', label: '通知声音', desc: '收到通知时播放提示音', enabled: true },
  { key: 'vibrate', label: '震动反馈', desc: '触摸操作时震动反馈', enabled: false },
  { key: 'badge', label: '角标显示', desc: '未读消息角标提醒', enabled: true },
])

// AI 模型
const selectedModel = ref('gpt-4o')
const temperature = ref(0.7)
const maxTokens = ref(2000)
const streamOutput = ref(true)
const showModelPicker = ref(false)

const models = [
  { id: 'gpt-4o', name: 'GPT-4o', desc: '最强大的模型，适合复杂对话' },
  { id: 'gpt-4o-mini', name: 'GPT-4o Mini', desc: '轻量快速，日常对话首选' },
  { id: 'claude-3.5-sonnet', name: 'Claude 3.5 Sonnet', desc: '擅长创意写作和角色扮演' },
  { id: 'claude-3-haiku', name: 'Claude 3 Haiku', desc: '快速响应，性价比高' },
  { id: 'deepseek-v3', name: 'DeepSeek V3', desc: '国产大模型，中文表现优秀' },
  { id: 'gemini-pro', name: 'Gemini Pro', desc: 'Google 出品，多模态能力' },
]

// 聊天设置
const autoSave = ref(true)
const bubbleStyleIndex = ref(0)
const bubbleStyles = [
  { name: '圆角气泡' },
  { name: '方角气泡' },
  { name: '极简风格' },
]
const sendKeyBehavior = ref<'send' | 'newline'>('send')

// 操作
const showClearConfirm = ref(false)
const saved = ref(false)

function saveSettings() {
  // 模拟保存
  saved.value = true
  setTimeout(() => { saved.value = false }, 2000)
}

function clearAllData() {
  showClearConfirm.value = false
  localStorage.clear()
  saved.value = true
  setTimeout(() => { saved.value = false }, 2000)
}

function exportData() {
  const data = {
    darkMode: darkMode.value,
    themeColor: selectedColor.value,
    wallpaper: selectedWallpaper.value,
    fontSize: fontSizeIndex.value,
    model: selectedModel.value,
    temperature: temperature.value,
    maxTokens: maxTokens.value,
    exportTime: new Date().toISOString(),
  }
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'fzsmphone-settings.json'
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.customize-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.settings-sections {
  padding: 8px 12px 24px;
}

.section {
  background: var(--bg-secondary);
  border-radius: 14px;
  margin-bottom: 12px;
  overflow: hidden;
}

.section-header {
  padding: 12px 16px;
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-secondary);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-secondary);
  cursor: pointer;
}

.setting-item:last-child { border-bottom: none; }

.setting-label {
  flex: 1;
  min-width: 0;
}

.label-text {
  font-size: 14px;
  color: var(--text-primary);
  display: block;
}

.label-text.danger { color: #ff3b30; }

.label-desc {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 1px;
  display: block;
}

.arrow {
  font-size: 18px;
  color: var(--text-quaternary);
  cursor: pointer;
}

/* Toggle */
.toggle {
  width: 48px;
  height: 28px;
  border-radius: 14px;
  background: var(--bg-tertiary);
  position: relative;
  cursor: pointer;
  transition: background 0.3s;
  flex-shrink: 0;
}

.toggle.on {
  background: #34c759;
}

.toggle-thumb {
  width: 24px;
  height: 24px;
  border-radius: 12px;
  background: #fff;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.3s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.toggle.on .toggle-thumb {
  transform: translateX(20px);
}

/* Color options */
.color-options {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.color-dot {
  width: 28px;
  height: 28px;
  border-radius: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.color-dot.active {
  border-color: #fff;
  box-shadow: 0 0 0 2px currentColor;
  transform: scale(1.1);
}

/* Wallpaper */
.wallpaper-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  padding: 8px 16px 12px;
}

.wallpaper-item {
  height: 80px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  cursor: pointer;
  position: relative;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.wallpaper-item.active {
  border-color: var(--color-primary);
}

.wp-emoji { font-size: 24px; }
.wp-name { font-size: 10px; color: #fff; font-weight: 500; }

.wp-check {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 18px;
  height: 18px;
  border-radius: 9px;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Slider */
.slider-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 16px 12px;
}

.slider-label { font-size: 14px; color: var(--text-tertiary); }
.slider-label.large { font-size: 20px; }
.slider-min, .slider-max { font-size: 11px; color: var(--text-tertiary); white-space: nowrap; }

.slider {
  flex: 1;
  -webkit-appearance: none;
  height: 4px;
  border-radius: 2px;
  background: var(--bg-tertiary);
  outline: none;
}

.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 22px;
  height: 22px;
  border-radius: 11px;
  background: var(--color-primary);
  cursor: pointer;
  border: 2px solid #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.font-preview {
  padding: 8px 16px 12px;
  color: var(--text-secondary);
  text-align: center;
}

/* Model picker */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.picker-panel {
  width: 100%;
  max-width: 393px;
  max-height: 60%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
}

.picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-primary);
}

.picker-header h3 {
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

.model-list {
  overflow-y: auto;
  padding: 8px;
}

.model-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.model-item:active, .model-item.active {
  background: var(--bg-tertiary);
}

.model-name { font-size: 14px; font-weight: 600; color: var(--text-primary); }
.model-desc { font-size: 12px; color: var(--text-tertiary); margin-top: 2px; }

.check-icon {
  color: var(--color-primary);
  font-size: 18px;
  font-weight: 700;
}

/* Confirm dialog */
.confirm-dialog {
  background: var(--bg-primary);
  border-radius: 20px;
  padding: 24px;
  text-align: center;
  margin: auto;
  width: 280px;
}

.confirm-icon { font-size: 40px; margin-bottom: 8px; }

.confirm-dialog h3 {
  margin: 0 0 6px;
  font-size: 17px;
  color: var(--text-primary);
}

.confirm-dialog p {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0 0 16px;
}

.confirm-actions {
  display: flex;
  gap: 10px;
}

.confirm-actions button {
  flex: 1;
  padding: 10px;
  border-radius: 12px;
  border: none;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel { background: var(--bg-tertiary); color: var(--text-primary); }
.btn-danger { background: #ff3b30; color: #fff; }

/* Save */
.save-btn {
  width: 100%;
  padding: 14px;
  border-radius: 14px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  margin-top: 4px;
}

.save-toast {
  text-align: center;
  padding: 10px;
  color: #34c759;
  font-size: 14px;
  font-weight: 600;
}
</style>
