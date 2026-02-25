<template>
  <div class="orders-page">
    <NavBar title="我的订单" back-to="/takeaway" />

    <div class="page-content">
      <!-- 加载 -->
      <div class="loading" v-if="store.ordersLoading">加载中...</div>

      <!-- 空状态 -->
      <div class="empty" v-else-if="store.orders.length === 0">
        <div class="empty-icon">📦</div>
        <p>暂无订单</p>
        <p class="empty-hint">去点些好吃的吧~</p>
        <button class="go-order-btn" @click="$router.push('/restaurant')">去点餐</button>
      </div>

      <!-- 订单列表 -->
      <div class="order-list" v-else>
        <div
          v-for="order in store.orders"
          :key="order.id"
          class="order-card"
        >
          <div class="order-header">
            <div class="order-restaurant">
              <img v-if="order.restaurant_image" :src="order.restaurant_image" class="rest-thumb" />
              <span class="rest-thumb-placeholder" v-else>🏪</span>
              <span class="rest-name">{{ order.restaurant_name || '未知餐厅' }}</span>
            </div>
            <span class="order-status" :class="order.status">{{ statusText(order.status) }}</span>
          </div>

          <div class="order-items">
            <div v-for="(item, idx) in order.items" :key="idx" class="order-item">
              <span class="item-name">{{ item.name }} × {{ item.quantity }}</span>
              <span class="item-price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
            </div>
          </div>

          <div class="order-footer">
            <div class="order-total">
              合计 <strong>¥{{ order.total.toFixed(2) }}</strong>
            </div>
            <div class="order-time">{{ formatTime(order.created_at) }}</div>
          </div>

          <div class="order-actions" v-if="order.status === 'pending'">
            <button class="cancel-btn" @click="handleCancel(order.id)">取消订单</button>
          </div>

          <!-- 配送进度 -->
          <div class="delivery-progress" v-if="isActive(order.status)">
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: progressWidth(order.status) }"></div>
            </div>
            <div class="progress-steps">
              <span :class="{ done: stepDone(order.status, 0) }">已下单</span>
              <span :class="{ done: stepDone(order.status, 1) }">商家确认</span>
              <span :class="{ done: stepDone(order.status, 2) }">制作中</span>
              <span :class="{ done: stepDone(order.status, 3) }">配送中</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'
import { useTakeawayStore } from '@/stores/takeaway'

const store = useTakeawayStore()

onMounted(() => {
  store.fetchOrders()
})

function statusText(status: string): string {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    preparing: '制作中',
    delivering: '配送中',
    delivered: '已送达',
    cancelled: '已取消',
  }
  return map[status] || status
}

function isActive(status: string): boolean {
  return ['pending', 'confirmed', 'preparing', 'delivering'].includes(status)
}

function progressWidth(status: string): string {
  const map: Record<string, string> = {
    pending: '10%',
    confirmed: '35%',
    preparing: '60%',
    delivering: '85%',
  }
  return map[status] || '0%'
}

function stepDone(status: string, step: number): boolean {
  const order = ['pending', 'confirmed', 'preparing', 'delivering']
  const idx = order.indexOf(status)
  return idx >= step
}

function formatTime(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    const month = d.getMonth() + 1
    const day = d.getDate()
    const h = d.getHours().toString().padStart(2, '0')
    const m = d.getMinutes().toString().padStart(2, '0')
    return `${month}月${day}日 ${h}:${m}`
  } catch {
    return dateStr
  }
}

async function handleCancel(id: number) {
  if (!confirm('确定取消此订单？')) return
  try {
    await store.cancelOrder(id)
  } catch (e) {
    console.error('取消订单失败', e)
  }
}
</script>

<style scoped>
.orders-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.loading, .empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary);
}

.empty-icon { font-size: 48px; margin-bottom: 12px; }
.empty-hint { font-size: 13px; margin-top: 4px; opacity: 0.7; }

.go-order-btn {
  margin-top: 16px;
  padding: 10px 28px;
  border: none;
  border-radius: 20px;
  background: #ff6b35;
  color: white;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
}

/* 订单列表 */
.order-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.order-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 14px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.order-restaurant {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rest-thumb {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  object-fit: cover;
}

.rest-thumb-placeholder {
  font-size: 20px;
}

.rest-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-primary);
}

.order-status {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 10px;
  font-weight: 600;
}

.order-status.pending { background: #fef3c7; color: #92400e; }
.order-status.confirmed { background: #dbeafe; color: #1e40af; }
.order-status.preparing { background: #fce7f3; color: #9d174d; }
.order-status.delivering { background: #d1fae5; color: #065f46; }
.order-status.delivered { background: #e5e7eb; color: #374151; }
.order-status.cancelled { background: #fee2e2; color: #991b1b; }

/* 订单商品 */
.order-items {
  border-top: 1px solid var(--border-color);
  border-bottom: 1px solid var(--border-color);
  padding: 8px 0;
  margin-bottom: 8px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 0;
  font-size: 13px;
}

.item-name { color: var(--text-primary); }
.item-price { color: var(--text-secondary); }

/* 底部 */
.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.order-total {
  color: var(--text-primary);
}

.order-total strong {
  color: #ff6b35;
  font-size: 16px;
}

.order-time {
  color: var(--text-tertiary, #999);
  font-size: 12px;
}

/* 操作 */
.order-actions {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}

.cancel-btn {
  padding: 6px 16px;
  border: 1px solid var(--border-color);
  border-radius: 16px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
}

.cancel-btn:active { background: var(--bg-primary); }

/* 配送进度 */
.delivery-progress {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.progress-bar {
  height: 4px;
  background: var(--bg-primary);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #ff6b35, #ff9a5c);
  border-radius: 2px;
  transition: width 0.5s ease;
}

.progress-steps {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text-tertiary, #999);
}

.progress-steps .done {
  color: #ff6b35;
  font-weight: 600;
}
</style>
