<template>
  <div class="couple-space-page">
    <NavBar title="情侣空间" />

    <!-- 头部横幅 -->
    <div class="couple-banner">
      <div class="banner-bg"></div>
      <div class="couple-avatars">
        <div class="avatar-wrapper">
          <div class="avatar">💑</div>
          <span class="avatar-name">{{ partnerA }}</span>
        </div>
        <div class="heart-icon">
          <span class="heart-beat">❤️</span>
          <span class="days-count">{{ daysTogether }} 天</span>
        </div>
        <div class="avatar-wrapper">
          <div class="avatar">💕</div>
          <span class="avatar-name">{{ partnerB }}</span>
        </div>
      </div>
      <div class="anniversary-text">在一起的第 {{ daysTogether }} 天</div>
    </div>

    <!-- 功能网格 -->
    <div class="features-grid">
      <div
        v-for="feat in features"
        :key="feat.label"
        class="feature-card"
        @click="handleFeature(feat.action)"
      >
        <span class="feature-icon">{{ feat.icon }}</span>
        <span class="feature-label">{{ feat.label }}</span>
      </div>
    </div>

    <!-- 纪念日列表 -->
    <div class="section">
      <div class="section-header">
        <h3>纪念日</h3>
        <button class="add-btn" @click="showAddAnniversary = true">+ 添加</button>
      </div>
      <div v-if="anniversaries.length === 0" class="empty-hint">
        还没有纪念日，快来添加吧~
      </div>
      <div v-else class="anniversary-list">
        <div v-for="ann in anniversaries" :key="ann.id" class="anniversary-item">
          <div class="ann-icon">{{ ann.icon }}</div>
          <div class="ann-info">
            <span class="ann-title">{{ ann.title }}</span>
            <span class="ann-date">{{ ann.date }}</span>
          </div>
          <div class="ann-countdown">
            <span class="countdown-num">{{ ann.daysLeft }}</span>
            <span class="countdown-label">{{ ann.daysLeft >= 0 ? '天后' : '天前' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 情侣相册 -->
    <div class="section">
      <div class="section-header">
        <h3>我们的相册</h3>
        <button class="add-btn" @click="addPhoto">+ 上传</button>
      </div>
      <div v-if="photos.length === 0" class="empty-hint">
        还没有照片，快来记录甜蜜瞬间~
      </div>
      <div v-else class="photo-grid">
        <div v-for="(photo, idx) in photos" :key="idx" class="photo-item">
          <div class="photo-placeholder">📷</div>
          <span class="photo-date">{{ photo.date }}</span>
        </div>
      </div>
    </div>

    <!-- 情侣任务 -->
    <div class="section">
      <div class="section-header">
        <h3>情侣任务</h3>
      </div>
      <div class="task-list">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-item"
          :class="{ completed: task.done }"
          @click="toggleTask(task)"
        >
          <div class="task-check">{{ task.done ? '✅' : '⬜' }}</div>
          <span class="task-text">{{ task.text }}</span>
        </div>
      </div>
    </div>

    <!-- 添加纪念日弹窗 -->
    <div v-if="showAddAnniversary" class="modal-overlay" @click.self="showAddAnniversary = false">
      <div class="modal-panel">
        <h3>添加纪念日</h3>
        <div class="form-group">
          <label>名称</label>
          <input v-model="newAnn.title" placeholder="例如：在一起纪念日" />
        </div>
        <div class="form-group">
          <label>日期</label>
          <input v-model="newAnn.date" type="date" />
        </div>
        <div class="form-group">
          <label>图标</label>
          <div class="icon-picker">
            <span
              v-for="icon in annIcons"
              :key="icon"
              class="icon-option"
              :class="{ selected: newAnn.icon === icon }"
              @click="newAnn.icon = icon"
            >{{ icon }}</span>
          </div>
        </div>
        <div class="modal-actions">
          <button class="cancel-btn" @click="showAddAnniversary = false">取消</button>
          <button class="confirm-btn" @click="addAnniversary">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

const partnerA = ref('我')
const partnerB = ref('TA')
const daysTogether = ref(365)

const features = [
  { icon: '📅', label: '纪念日', action: 'anniversary' },
  { icon: '📸', label: '相册', action: 'album' },
  { icon: '💌', label: '情书', action: 'letter' },
  { icon: '🎯', label: '心愿单', action: 'wish' },
  { icon: '🎮', label: '小游戏', action: 'game' },
  { icon: '📍', label: '足迹', action: 'footprint' },
]

interface Anniversary {
  id: number
  title: string
  date: string
  icon: string
  daysLeft: number
}

const anniversaries = ref<Anniversary[]>([
  { id: 1, title: '在一起纪念日', date: '2024-02-14', icon: '❤️', daysLeft: 354 },
  { id: 2, title: '第一次旅行', date: '2024-05-01', icon: '✈️', daysLeft: 65 },
])

const annIcons = ['❤️', '🎂', '✈️', '💍', '🌹', '🎄', '🎆', '🏠']

const showAddAnniversary = ref(false)
const newAnn = reactive({
  title: '',
  date: '',
  icon: '❤️',
})

function addAnniversary() {
  if (!newAnn.title || !newAnn.date) return
  const target = new Date(newAnn.date)
  const now = new Date()
  const daysLeft = Math.ceil((target.getTime() - now.getTime()) / 86400000)
  anniversaries.value.push({
    id: Date.now(),
    title: newAnn.title,
    date: newAnn.date,
    icon: newAnn.icon,
    daysLeft,
  })
  newAnn.title = ''
  newAnn.date = ''
  newAnn.icon = '❤️'
  showAddAnniversary.value = false
}

interface Photo {
  date: string
  url?: string
}

const photos = ref<Photo[]>([
  { date: '2024-12-25' },
  { date: '2024-10-01' },
  { date: '2024-08-15' },
])

function addPhoto() {
  photos.value.unshift({ date: new Date().toISOString().slice(0, 10) })
}

interface CoupleTask {
  id: number
  text: string
  done: boolean
}

const tasks = ref<CoupleTask[]>([
  { id: 1, text: '一起看日出', done: false },
  { id: 2, text: '一起做一顿饭', done: true },
  { id: 3, text: '一起去游乐园', done: false },
  { id: 4, text: '一起拍一组情侣照', done: false },
  { id: 5, text: '一起写一封信给未来的自己', done: true },
  { id: 6, text: '一起看一场星空', done: false },
])

function toggleTask(task: CoupleTask) {
  task.done = !task.done
}

function handleFeature(action: string) {
  // 各功能入口，目前仅滚动到对应区域
  console.log('Feature:', action)
}
</script>

<style scoped>
.couple-space-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

/* 横幅 */
.couple-banner {
  position: relative;
  padding: 30px 16px 20px;
  text-align: center;
  overflow: hidden;
}

.banner-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #ff6b9d 0%, #c44dff 50%, #6e8efb 100%);
  opacity: 0.15;
}

.couple-avatars {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 12px;
}

.avatar-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  border: 3px solid rgba(255, 107, 157, 0.3);
}

.avatar-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.heart-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.heart-beat {
  font-size: 28px;
  animation: heartBeat 1.5s ease-in-out infinite;
}

@keyframes heartBeat {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

.days-count {
  font-size: 12px;
  color: var(--accent-pink);
  font-weight: 700;
}

.anniversary-text {
  position: relative;
  font-size: 14px;
  color: var(--text-secondary);
}

/* 功能网格 */
.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  padding: 16px;
}

.feature-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 8px;
  background: var(--bg-secondary);
  border-radius: 14px;
  cursor: pointer;
  transition: transform 0.2s;
}

.feature-card:active {
  transform: scale(0.95);
}

.feature-icon {
  font-size: 28px;
}

.feature-label {
  font-size: 12px;
  color: var(--text-primary);
  font-weight: 500;
}

/* 区块 */
.section {
  padding: 0 16px 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.add-btn {
  background: none;
  border: none;
  font-size: 13px;
  color: var(--accent-pink);
  cursor: pointer;
  font-weight: 600;
}

.empty-hint {
  text-align: center;
  padding: 24px;
  color: var(--text-tertiary);
  font-size: 13px;
}

/* 纪念日 */
.anniversary-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.anniversary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  background: var(--bg-secondary);
  border-radius: 14px;
}

.ann-icon {
  font-size: 28px;
}

.ann-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.ann-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.ann-date {
  font-size: 12px;
  color: var(--text-tertiary);
}

.ann-countdown {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.countdown-num {
  font-size: 20px;
  font-weight: 800;
  color: var(--accent-pink);
}

.countdown-label {
  font-size: 11px;
  color: var(--text-tertiary);
}

/* 相册 */
.photo-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.photo-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.photo-placeholder {
  width: 100%;
  aspect-ratio: 1;
  background: var(--bg-secondary);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
}

.photo-date {
  font-size: 11px;
  color: var(--text-tertiary);
}

/* 情侣任务 */
.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: var(--bg-secondary);
  border-radius: 12px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.task-item.completed {
  opacity: 0.5;
}

.task-check {
  font-size: 18px;
}

.task-text {
  font-size: 14px;
  color: var(--text-primary);
}

.task-item.completed .task-text {
  text-decoration: line-through;
  color: var(--text-tertiary);
}

/* 弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 24px;
}

.modal-panel {
  width: 100%;
  max-width: 340px;
  background: var(--bg-secondary);
  border-radius: 20px;
  padding: 24px;
}

.modal-panel h3 {
  margin: 0 0 16px;
  font-size: 17px;
  color: var(--text-primary);
}

.form-group {
  margin-bottom: 14px;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
}

.form-group input:focus {
  border-color: var(--accent-pink);
}

.icon-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.icon-option {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  border-radius: 10px;
  background: var(--bg-primary);
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}

.icon-option.selected {
  border-color: var(--accent-pink);
  background: rgba(255, 107, 157, 0.1);
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 18px;
}

.cancel-btn,
.confirm-btn {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  cursor: pointer;
}

.cancel-btn {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.confirm-btn {
  background: var(--accent-pink);
  color: #fff;
}
</style>
