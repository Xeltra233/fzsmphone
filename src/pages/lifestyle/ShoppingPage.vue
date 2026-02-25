<template>
  <div class="shopping-page">
    <NavBar title="购物" back />

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input">
        <span class="search-icon">🔍</span>
        <input v-model="searchText" placeholder="搜索商品" @input="filterProducts" />
        <span v-if="searchText" class="clear-btn" @click="searchText = ''; filterProducts()">✕</span>
      </div>
      <div class="cart-icon" @click="showCart = true">
        🛒
        <span v-if="cartCount > 0" class="cart-badge">{{ cartCount }}</span>
      </div>
    </div>

    <!-- 分类标签 -->
    <div class="category-tabs">
      <div
        v-for="cat in categories"
        :key="cat.id"
        class="tab"
        :class="{ active: activeCategory === cat.id }"
        @click="activeCategory = cat.id; filterProducts()"
      >
        <span class="tab-icon">{{ cat.icon }}</span>
        <span class="tab-label">{{ cat.name }}</span>
      </div>
    </div>

    <!-- Banner -->
    <div v-if="activeCategory === 'all'" class="banner-area">
      <div class="banner" :style="{ background: banners[currentBanner].bg }">
        <div class="banner-text">
          <div class="banner-title">{{ banners[currentBanner].title }}</div>
          <div class="banner-sub">{{ banners[currentBanner].sub }}</div>
        </div>
        <div class="banner-emoji">{{ banners[currentBanner].emoji }}</div>
      </div>
      <div class="banner-dots">
        <span v-for="(_, i) in banners" :key="i" :class="{ active: i === currentBanner }" @click="currentBanner = i"></span>
      </div>
    </div>

    <!-- 商品列表 -->
    <div class="products-grid">
      <div
        v-for="product in filteredProducts"
        :key="product.id"
        class="product-card"
        @click="selectProduct(product)"
      >
        <div class="product-image">{{ product.emoji }}</div>
        <div class="product-info">
          <div class="product-name">{{ product.name }}</div>
          <div class="product-desc">{{ product.desc }}</div>
          <div class="product-bottom">
            <span class="product-price">¥{{ product.price.toFixed(2) }}</span>
            <span v-if="product.originalPrice" class="original-price">¥{{ product.originalPrice.toFixed(2) }}</span>
            <span class="product-sales">{{ product.sales }}人付款</span>
          </div>
        </div>
      </div>
      <div v-if="filteredProducts.length === 0" class="empty-state">
        <div class="empty-icon">🔍</div>
        <div class="empty-text">没有找到相关商品</div>
      </div>
    </div>

    <!-- 商品详情弹窗 -->
    <div v-if="selectedProduct" class="modal-overlay" @click.self="selectedProduct = null">
      <div class="product-detail">
        <div class="detail-header">
          <div class="detail-image">{{ selectedProduct.emoji }}</div>
          <span class="close-btn" @click="selectedProduct = null">✕</span>
        </div>
        <div class="detail-body">
          <h3>{{ selectedProduct.name }}</h3>
          <div class="detail-price">
            <span class="price-now">¥{{ selectedProduct.price.toFixed(2) }}</span>
            <span v-if="selectedProduct.originalPrice" class="price-original">¥{{ selectedProduct.originalPrice.toFixed(2) }}</span>
            <span class="discount-tag" v-if="selectedProduct.originalPrice">
              {{ Math.round((1 - selectedProduct.price / selectedProduct.originalPrice) * 100) }}% OFF
            </span>
          </div>
          <p class="detail-desc">{{ selectedProduct.fullDesc }}</p>
          <div class="detail-specs">
            <div class="spec-item" v-for="spec in selectedProduct.specs" :key="spec.label">
              <span class="spec-label">{{ spec.label }}</span>
              <span class="spec-value">{{ spec.value }}</span>
            </div>
          </div>
          <div class="quantity-row">
            <span>数量</span>
            <div class="quantity-control">
              <button @click="detailQty = Math.max(1, detailQty - 1)">−</button>
              <span>{{ detailQty }}</span>
              <button @click="detailQty++">+</button>
            </div>
          </div>
          <div class="detail-actions">
            <button class="btn-fav" @click="toggleFav(selectedProduct)">
              {{ selectedProduct.fav ? '❤️ 已收藏' : '🤍 收藏' }}
            </button>
            <button class="btn-cart" @click="addToCart(selectedProduct, detailQty); selectedProduct = null">
              加入购物车
            </button>
            <button class="btn-buy" @click="buyNow(selectedProduct, detailQty)">
              立即购买
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 购物车弹窗 -->
    <div v-if="showCart" class="modal-overlay" @click.self="showCart = false">
      <div class="cart-panel">
        <div class="cart-header">
          <h3>购物车</h3>
          <span class="close-btn" @click="showCart = false">✕</span>
        </div>
        <div v-if="cart.length === 0" class="cart-empty">
          <div class="empty-icon">🛒</div>
          <div>购物车是空的</div>
        </div>
        <div v-else class="cart-items">
          <div v-for="item in cart" :key="item.product.id" class="cart-item">
            <div class="cart-item-emoji">{{ item.product.emoji }}</div>
            <div class="cart-item-info">
              <div class="cart-item-name">{{ item.product.name }}</div>
              <div class="cart-item-price">¥{{ item.product.price.toFixed(2) }}</div>
            </div>
            <div class="cart-item-qty">
              <button @click="updateCartQty(item, -1)">−</button>
              <span>{{ item.qty }}</span>
              <button @click="updateCartQty(item, 1)">+</button>
            </div>
            <span class="cart-item-delete" @click="removeFromCart(item)">🗑️</span>
          </div>
        </div>
        <div v-if="cart.length > 0" class="cart-footer">
          <div class="cart-total">
            <span>合计：</span>
            <span class="total-price">¥{{ cartTotal.toFixed(2) }}</span>
          </div>
          <button class="btn-checkout" @click="checkout">结算 ({{ cartCount }})</button>
        </div>
      </div>
    </div>

    <!-- 订单成功弹窗 -->
    <div v-if="orderSuccess" class="modal-overlay" @click.self="orderSuccess = false">
      <div class="order-success">
        <div class="success-icon">✅</div>
        <h3>下单成功！</h3>
        <p>订单号：{{ orderNo }}</p>
        <p>支付金额：¥{{ orderAmount.toFixed(2) }}</p>
        <button @click="orderSuccess = false">好的</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import NavBar from '@/components/common/NavBar.vue'

interface ProductSpec {
  label: string
  value: string
}

interface Product {
  id: number
  name: string
  emoji: string
  desc: string
  fullDesc: string
  price: number
  originalPrice?: number
  category: string
  sales: number
  fav: boolean
  specs: ProductSpec[]
}

interface CartItem {
  product: Product
  qty: number
}

const categories = [
  { id: 'all', name: '全部', icon: '🏠' },
  { id: 'digital', name: '数码', icon: '📱' },
  { id: 'fashion', name: '服饰', icon: '👗' },
  { id: 'beauty', name: '美妆', icon: '💄' },
  { id: 'food', name: '食品', icon: '🍪' },
  { id: 'home', name: '家居', icon: '🏡' },
  { id: 'gift', name: '礼物', icon: '🎁' },
]

const allProducts = ref<Product[]>([
  { id: 1, name: 'AirPods Pro 2', emoji: '🎧', desc: '主动降噪 / 自适应音频', fullDesc: '全新AirPods Pro 2代，搭载H2芯片，提供高达2倍的主动降噪性能，自适应通透模式，个性化空间音频。', price: 1799, originalPrice: 1999, category: 'digital', sales: 28340, fav: false, specs: [{ label: '芯片', value: 'Apple H2' }, { label: '续航', value: '6小时' }, { label: '充电', value: 'MagSafe / Lightning' }] },
  { id: 2, name: 'iPhone 15 Pro', emoji: '📱', desc: 'A17 Pro芯片 / 钛金属', fullDesc: 'iPhone 15 Pro，采用航空级钛金属设计，A17 Pro芯片带来变革性能。48MP主摄，USB-C接口。', price: 7999, originalPrice: 8999, category: 'digital', sales: 15670, fav: false, specs: [{ label: '芯片', value: 'A17 Pro' }, { label: '屏幕', value: '6.1" OLED' }, { label: '存储', value: '256GB' }] },
  { id: 3, name: 'MacBook Air M3', emoji: '💻', desc: 'M3芯片 / 18小时续航', fullDesc: 'MacBook Air搭载M3芯片，13.6英寸Liquid Retina显示屏，最长18小时电池续航。', price: 8999, originalPrice: 9999, category: 'digital', sales: 8930, fav: false, specs: [{ label: '芯片', value: 'Apple M3' }, { label: '内存', value: '8GB' }, { label: '续航', value: '18小时' }] },
  { id: 4, name: '法式连衣裙', emoji: '👗', desc: '春季新款 / 碎花设计', fullDesc: '法式优雅碎花连衣裙，甜美浪漫风格，收腰设计显瘦，适合约会/日常穿搭。', price: 299, originalPrice: 599, category: 'fashion', sales: 42100, fav: false, specs: [{ label: '面料', value: '雪纺' }, { label: '尺码', value: 'S/M/L/XL' }, { label: '季节', value: '春夏' }] },
  { id: 5, name: '情侣卫衣套装', emoji: '👫', desc: '加绒保暖 / 情侣款', fullDesc: '宽松休闲情侣卫衣套装，内里加绒保暖，多色可选，情侣一起穿超甜。', price: 199, originalPrice: 398, category: 'fashion', sales: 31500, fav: false, specs: [{ label: '面料', value: '棉+绒' }, { label: '尺码', value: 'M-3XL' }, { label: '颜色', value: '5色' }] },
  { id: 6, name: 'SK-II 神仙水', emoji: '✨', desc: '230ml / 护肤精华', fullDesc: 'SK-II护肤精华露，蕴含超过90%的天然活肤酵母精华PITERA™，改善肤质，晶莹剔透。', price: 1190, originalPrice: 1590, category: 'beauty', sales: 19800, fav: false, specs: [{ label: '容量', value: '230ml' }, { label: '肤质', value: '通用' }, { label: '功效', value: '焕肤提亮' }] },
  { id: 7, name: 'YSL小金条', emoji: '💋', desc: '#1966 / 正红色', fullDesc: 'YSL圣罗兰小金条口红，丝绒质地，显色持久，#1966经典正红色。', price: 320, originalPrice: 380, category: 'beauty', sales: 55600, fav: false, specs: [{ label: '色号', value: '#1966' }, { label: '质地', value: '丝绒' }, { label: '功效', value: '显色保湿' }] },
  { id: 8, name: '手工巧克力礼盒', emoji: '🍫', desc: '比利时进口 / 24颗装', fullDesc: '比利时进口手工巧克力礼盒，24颗精选混合口味，送礼自用皆宜。', price: 168, originalPrice: 258, category: 'food', sales: 12300, fav: false, specs: [{ label: '数量', value: '24颗' }, { label: '口味', value: '混合' }, { label: '保质期', value: '6个月' }] },
  { id: 9, name: '日式抹茶礼盒', emoji: '🍵', desc: '宇治抹茶 / 精装', fullDesc: '京都宇治抹茶系列礼盒，含抹茶粉、抹茶饼干、抹茶巧克力，抹茶控必入。', price: 128, category: 'food', sales: 8700, fav: false, specs: [{ label: '产地', value: '日本京都' }, { label: '包含', value: '3件套' }] },
  { id: 10, name: '香薰蜡烛', emoji: '🕯️', desc: '白茶 / 天然大豆蜡', fullDesc: '天然大豆蜡香薰蜡烛，白茶清香，营造温馨氛围，燃烧时间约50小时。', price: 89, originalPrice: 128, category: 'home', sales: 25400, fav: false, specs: [{ label: '香型', value: '白茶' }, { label: '燃烧', value: '50小时' }, { label: '重量', value: '200g' }] },
  { id: 11, name: '北欧ins抱枕', emoji: '🛋️', desc: '天鹅绒 / 45cm', fullDesc: '北欧简约风格抱枕，天鹅绒面料柔软亲肤，45x45cm，含芯。', price: 49, originalPrice: 79, category: 'home', sales: 67800, fav: false, specs: [{ label: '尺寸', value: '45x45cm' }, { label: '面料', value: '天鹅绒' }, { label: '填充', value: 'PP棉' }] },
  { id: 12, name: '定制情侣项链', emoji: '💎', desc: '纯银 / 刻字定制', fullDesc: '925纯银情侣项链，可定制刻字，锁骨链设计，送给最爱的人。', price: 299, originalPrice: 499, category: 'gift', sales: 14200, fav: false, specs: [{ label: '材质', value: '925银' }, { label: '链长', value: '40+5cm' }, { label: '服务', value: '免费刻字' }] },
  { id: 13, name: '永生花礼盒', emoji: '🌹', desc: '厄瓜多尔玫瑰 / 保存3年', fullDesc: '进口厄瓜多尔永生花礼盒，真花制作，保存3年不凋谢，精美礼盒装。', price: 259, originalPrice: 399, category: 'gift', sales: 21300, fav: false, specs: [{ label: '花材', value: '厄瓜多尔玫瑰' }, { label: '保存', value: '3年+' }, { label: '包装', value: '豪华礼盒' }] },
  { id: 14, name: 'Apple Watch Ultra 2', emoji: '⌚', desc: '钛金属 / 双频GPS', fullDesc: 'Apple Watch Ultra 2，49mm钛金属表壳，精确双频GPS，最长72小时续航。', price: 5999, originalPrice: 6499, category: 'digital', sales: 6100, fav: false, specs: [{ label: '尺寸', value: '49mm' }, { label: '芯片', value: 'S9 SiP' }, { label: '续航', value: '72小时' }] },
])

const searchText = ref('')
const activeCategory = ref('all')
const filteredProducts = ref<Product[]>([])
const selectedProduct = ref<Product | null>(null)
const detailQty = ref(1)
const cart = ref<CartItem[]>([])
const showCart = ref(false)
const orderSuccess = ref(false)
const orderNo = ref('')
const orderAmount = ref(0)
const currentBanner = ref(0)
let bannerTimer: ReturnType<typeof setInterval> | null = null

const banners = [
  { title: '春季大促', sub: '全场低至5折', emoji: '🌸', bg: 'linear-gradient(135deg, #ff9a9e, #fecfef)' },
  { title: '情侣专区', sub: '甜蜜好物推荐', emoji: '💕', bg: 'linear-gradient(135deg, #a18cd1, #fbc2eb)' },
  { title: '数码狂欢', sub: '新品首发优惠', emoji: '🔥', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
]

const cartCount = computed(() => cart.value.reduce((s, i) => s + i.qty, 0))
const cartTotal = computed(() => cart.value.reduce((s, i) => s + i.product.price * i.qty, 0))

function filterProducts() {
  let list = allProducts.value
  if (activeCategory.value !== 'all') {
    list = list.filter(p => p.category === activeCategory.value)
  }
  if (searchText.value.trim()) {
    const kw = searchText.value.trim().toLowerCase()
    list = list.filter(p => p.name.toLowerCase().includes(kw) || p.desc.toLowerCase().includes(kw))
  }
  filteredProducts.value = list
}

function selectProduct(p: Product) {
  selectedProduct.value = p
  detailQty.value = 1
}

function toggleFav(p: Product) {
  p.fav = !p.fav
}

function addToCart(p: Product, qty: number) {
  const existing = cart.value.find(i => i.product.id === p.id)
  if (existing) {
    existing.qty += qty
  } else {
    cart.value.push({ product: p, qty })
  }
}

function updateCartQty(item: CartItem, delta: number) {
  item.qty += delta
  if (item.qty <= 0) {
    cart.value = cart.value.filter(i => i.product.id !== item.product.id)
  }
}

function removeFromCart(item: CartItem) {
  cart.value = cart.value.filter(i => i.product.id !== item.product.id)
}

function buyNow(p: Product, qty: number) {
  orderAmount.value = p.price * qty
  orderNo.value = 'SP' + Date.now().toString(36).toUpperCase()
  selectedProduct.value = null
  orderSuccess.value = true
}

function checkout() {
  orderAmount.value = cartTotal.value
  orderNo.value = 'SP' + Date.now().toString(36).toUpperCase()
  cart.value = []
  showCart.value = false
  orderSuccess.value = true
}

onMounted(() => {
  filterProducts()
  bannerTimer = setInterval(() => {
    currentBanner.value = (currentBanner.value + 1) % banners.length
  }, 4000)
})

onUnmounted(() => {
  if (bannerTimer) clearInterval(bannerTimer)
})
</script>

<style scoped>
.shopping-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow-y: auto;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  background: var(--bg-primary);
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-tertiary);
  border-radius: 20px;
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
.clear-btn { cursor: pointer; opacity: 0.5; font-size: 12px; }

.cart-icon {
  position: relative;
  font-size: 22px;
  cursor: pointer;
}

.cart-badge {
  position: absolute;
  top: -6px;
  right: -8px;
  background: #ff3b30;
  color: #fff;
  font-size: 10px;
  min-width: 16px;
  height: 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.category-tabs {
  display: flex;
  gap: 4px;
  padding: 4px 12px 8px;
  overflow-x: auto;
}

.category-tabs::-webkit-scrollbar { display: none; }

.tab {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 12px;
  border-radius: 12px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
}

.tab.active {
  background: var(--color-primary);
}

.tab-icon { font-size: 18px; }
.tab-label { font-size: 11px; color: var(--text-secondary); white-space: nowrap; }
.tab.active .tab-label { color: #fff; font-weight: 600; }

.banner-area {
  padding: 0 16px 8px;
}

.banner {
  border-radius: 16px;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #fff;
}

.banner-title { font-size: 20px; font-weight: 700; }
.banner-sub { font-size: 13px; margin-top: 4px; opacity: 0.9; }
.banner-emoji { font-size: 40px; }

.banner-dots {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-top: 8px;
}

.banner-dots span {
  width: 6px;
  height: 6px;
  border-radius: 3px;
  background: var(--text-quaternary);
  cursor: pointer;
  transition: all 0.3s;
}

.banner-dots span.active {
  width: 18px;
  background: var(--color-primary);
}

.products-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  padding: 8px 12px 20px;
}

.product-card {
  background: var(--bg-secondary);
  border-radius: 14px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s;
}

.product-card:active { transform: scale(0.97); }

.product-image {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 50px;
  background: var(--bg-tertiary);
}

.product-info {
  padding: 10px;
}

.product-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-bottom {
  display: flex;
  align-items: baseline;
  gap: 4px;
  flex-wrap: wrap;
}

.product-price {
  font-size: 16px;
  font-weight: 700;
  color: #ff3b30;
}

.original-price {
  font-size: 11px;
  color: var(--text-quaternary);
  text-decoration: line-through;
}

.product-sales {
  font-size: 10px;
  color: var(--text-quaternary);
  margin-left: auto;
}

/* 详情弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100;
}

.product-detail {
  width: 100%;
  max-width: 393px;
  max-height: 85%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  overflow-y: auto;
}

.detail-header {
  position: relative;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 80px;
  background: var(--bg-tertiary);
}

.close-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 28px;
  height: 28px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
}

.detail-body {
  padding: 16px;
}

.detail-body h3 {
  font-size: 18px;
  margin: 0 0 8px;
  color: var(--text-primary);
}

.detail-price {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 12px;
}

.price-now { font-size: 22px; font-weight: 700; color: #ff3b30; }
.price-original { font-size: 13px; color: var(--text-quaternary); text-decoration: line-through; }

.discount-tag {
  font-size: 11px;
  background: #ff3b30;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
}

.detail-desc {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 12px;
}

.detail-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.spec-item {
  background: var(--bg-tertiary);
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 12px;
}

.spec-label { color: var(--text-tertiary); margin-right: 4px; }
.spec-value { color: var(--text-primary); font-weight: 500; }

.quantity-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-size: 14px;
  color: var(--text-primary);
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.quantity-control button {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  border: 1px solid var(--border-primary);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quantity-control span { font-weight: 600; min-width: 20px; text-align: center; }

.detail-actions {
  display: flex;
  gap: 8px;
}

.detail-actions button {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.btn-fav { background: var(--bg-tertiary); color: var(--text-primary); }
.btn-cart { background: #ff9500; color: #fff; }
.btn-buy { background: #ff3b30; color: #fff; }

/* 购物车 */
.cart-panel {
  width: 100%;
  max-width: 393px;
  max-height: 70%;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  display: flex;
  flex-direction: column;
}

.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-primary);
}

.cart-header h3 {
  margin: 0;
  font-size: 17px;
  color: var(--text-primary);
}

.cart-header .close-btn {
  position: static;
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.cart-empty {
  padding: 40px;
  text-align: center;
  color: var(--text-tertiary);
  font-size: 14px;
}

.cart-empty .empty-icon { font-size: 40px; margin-bottom: 8px; }

.cart-items {
  flex: 1;
  overflow-y: auto;
  padding: 8px 16px;
}

.cart-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-secondary);
}

.cart-item-emoji { font-size: 32px; }

.cart-item-info { flex: 1; }
.cart-item-name { font-size: 14px; font-weight: 500; color: var(--text-primary); }
.cart-item-price { font-size: 13px; color: #ff3b30; font-weight: 600; }

.cart-item-qty {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cart-item-qty button {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  border: 1px solid var(--border-primary);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cart-item-qty span { font-size: 13px; font-weight: 600; min-width: 16px; text-align: center; }
.cart-item-delete { cursor: pointer; font-size: 16px; }

.cart-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-primary);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cart-total { font-size: 14px; color: var(--text-primary); }
.total-price { font-size: 18px; font-weight: 700; color: #ff3b30; }

.btn-checkout {
  padding: 10px 24px;
  background: #ff3b30;
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

/* 订单成功 */
.order-success {
  background: var(--bg-primary);
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  margin: auto;
  width: 280px;
}

.success-icon { font-size: 48px; margin-bottom: 12px; }

.order-success h3 {
  margin: 0 0 8px;
  font-size: 18px;
  color: var(--text-primary);
}

.order-success p {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 4px 0;
}

.order-success button {
  margin-top: 16px;
  padding: 10px 32px;
  background: var(--color-primary);
  color: #fff;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: var(--text-tertiary);
}

.empty-state .empty-icon { font-size: 40px; margin-bottom: 8px; }
.empty-state .empty-text { font-size: 14px; }
</style>
