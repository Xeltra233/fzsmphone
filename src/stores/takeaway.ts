import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { restaurantApi, orderApi } from '@/api/services'
import type { Restaurant, Order, OrderInput } from '@/api/types'

export const useTakeawayStore = defineStore('takeaway', () => {
  // 餐厅
  const restaurants = ref<Restaurant[]>([])
  const currentRestaurant = ref<Restaurant | null>(null)
  const restaurantsLoading = ref(false)

  // 订单
  const orders = ref<Order[]>([])
  const ordersLoading = ref(false)

  // 购物车
  const cart = ref<Map<string, { name: string; price: number; quantity: number }>>(new Map())

  const cartItems = computed(() => Array.from(cart.value.values()))
  const cartTotal = computed(() =>
    cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
  )
  const cartCount = computed(() =>
    cartItems.value.reduce((sum, item) => sum + item.quantity, 0)
  )

  // ===== 餐厅 =====
  async function fetchRestaurants(category?: string) {
    restaurantsLoading.value = true
    try {
      const res = await restaurantApi.list(category)
      restaurants.value = res.data || []
    } finally {
      restaurantsLoading.value = false
    }
  }

  async function fetchRestaurant(id: number) {
    currentRestaurant.value = await restaurantApi.get(id)
  }

  // ===== 订单 =====
  async function fetchOrders() {
    ordersLoading.value = true
    try {
      const res = await orderApi.list()
      orders.value = res.data || []
    } finally {
      ordersLoading.value = false
    }
  }

  async function createOrder(data: OrderInput) {
    const res = await orderApi.create(data)
    clearCart()
    await fetchOrders()
    return res
  }

  async function cancelOrder(id: number) {
    await orderApi.updateStatus(id, 'cancelled')
    const order = orders.value.find(o => o.id === id)
    if (order) order.status = 'cancelled'
  }

  // ===== 购物车 =====
  function addToCart(name: string, price: number) {
    const existing = cart.value.get(name)
    if (existing) {
      existing.quantity++
    } else {
      cart.value.set(name, { name, price, quantity: 1 })
    }
    // 触发响应式更新
    cart.value = new Map(cart.value)
  }

  function removeFromCart(name: string) {
    const existing = cart.value.get(name)
    if (!existing) return
    if (existing.quantity > 1) {
      existing.quantity--
    } else {
      cart.value.delete(name)
    }
    cart.value = new Map(cart.value)
  }

  function clearCart() {
    cart.value = new Map()
  }

  // 活跃订单（非已完成/已取消）
  const activeOrders = computed(() =>
    orders.value.filter(o => !['delivered', 'cancelled'].includes(o.status))
  )

  return {
    // 餐厅
    restaurants,
    currentRestaurant,
    restaurantsLoading,
    fetchRestaurants,
    fetchRestaurant,
    // 订单
    orders,
    ordersLoading,
    activeOrders,
    fetchOrders,
    createOrder,
    cancelOrder,
    // 购物车
    cart,
    cartItems,
    cartTotal,
    cartCount,
    addToCart,
    removeFromCart,
    clearCart,
  }
})
