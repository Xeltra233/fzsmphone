import api from './client'
import type {
  ListResponse,
  CreateResponse,
  MessageResponse,
  User,
  Character,
  CharacterInput,
  Conversation,
  ConversationInput,
  Message,
  MessageInput,
  Post,
  PostInput,
  Comment,
  CommentInput,
  WeiboPost,
  WeiboInput,
  Moment,
  MomentInput,
  Diary,
  DiaryInput,
  Persona,
  PersonaInput,
  WorldBookEntry,
  WorldBookInput,
  AppSettings,
  Restaurant,
  RestaurantInput,
  Order,
  OrderInput,
  Preset,
  PresetInput,
  CallRecord,
  CallRecordInput,
  SmsThread,
  SmsThreadInput,
  SmsMessage,
  SmsMessageInput,
  WalletInfo,
  WalletTransaction,
  WalletTransactionInput,
  GameRecordItem,
  GameRecordInput,
  Product,
  CartItem,
  ShoppingOrder,
  ShoppingOrderItem,
  CoupleSpace,
  Anniversary,
  CouplePhoto,
  CoupleTask,
} from './types'

// ========== 角色 API ==========
export const characterApi = {
  list(params?: { search?: string; tag?: string }) {
    const p: Record<string, string> = {}
    if (params?.search) p.search = params.search
    if (params?.tag) p.tag = params.tag
    return api.get<ListResponse<Character>>('/api/characters', p)
  },
  get(id: number) {
    return api.get<Character>(`/api/characters/${id}`)
  },
  create(data: CharacterInput) {
    return api.post<CreateResponse>('/api/characters', data)
  },
  update(id: number, data: CharacterInput) {
    return api.put<MessageResponse>(`/api/characters/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/characters/${id}`)
  },
}

// ========== 对话 API ==========
export const conversationApi = {
  list() {
    return api.get<ListResponse<Conversation>>('/api/conversations')
  },
  get(id: number) {
    return api.get<Conversation>(`/api/conversations/${id}`)
  },
  create(data: ConversationInput) {
    return api.post<CreateResponse>('/api/conversations', data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/conversations/${id}`)
  },
  listMessages(conversationId: number) {
    return api.get<ListResponse<Message>>(`/api/conversations/${conversationId}/messages`)
  },
  sendMessage(conversationId: number, data: MessageInput) {
    return api.post<CreateResponse>(`/api/conversations/${conversationId}/messages`, data)
  },
}

// ========== 帖子 API ==========
export const postApi = {
  list(params?: { board?: string; page?: number }) {
    const p: Record<string, string> = {}
    if (params?.board) p.board = params.board
    if (params?.page) p.page = String(params.page)
    return api.get<ListResponse<Post>>('/api/posts', p)
  },
  get(id: number) {
    return api.get<Post>(`/api/posts/${id}`)
  },
  create(data: PostInput) {
    return api.post<CreateResponse>('/api/posts', data)
  },
  update(id: number, data: Partial<PostInput>) {
    return api.put<MessageResponse>(`/api/posts/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/posts/${id}`)
  },
  listComments(postId: number) {
    return api.get<ListResponse<Comment>>(`/api/posts/${postId}/comments`)
  },
  createComment(postId: number, data: CommentInput) {
    return api.post<CreateResponse>(`/api/posts/${postId}/comments`, data)
  },
}

// ========== 微博 API ==========
export const weiboApi = {
  list(tab?: string) {
    const p: Record<string, string> = {}
    if (tab) p.tab = tab
    return api.get<ListResponse<WeiboPost>>('/api/weibo', p)
  },
  get(id: number) {
    return api.get<WeiboPost>(`/api/weibo/${id}`)
  },
  create(data: WeiboInput) {
    return api.post<CreateResponse>('/api/weibo', data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/weibo/${id}`)
  },
}

// ========== 朋友圈 API ==========
export const momentApi = {
  list() {
    return api.get<ListResponse<Moment>>('/api/moments')
  },
  create(data: MomentInput) {
    return api.post<CreateResponse>('/api/moments', data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/moments/${id}`)
  },
}

// ========== 日记 API ==========
export const diaryApi = {
  list() {
    return api.get<ListResponse<Diary>>('/api/diaries')
  },
  get(id: number) {
    return api.get<Diary>(`/api/diaries/${id}`)
  },
  create(data: DiaryInput) {
    return api.post<CreateResponse>('/api/diaries', data)
  },
  update(id: number, data: DiaryInput) {
    return api.put<MessageResponse>(`/api/diaries/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/diaries/${id}`)
  },
}

// ========== 用户 API ==========
export const userApi = {
  list() {
    return api.get<ListResponse<User>>('/api/users')
  },
  get(id: number) {
    return api.get<User>(`/api/users/${id}`)
  },
  updateRole(id: number, role: string) {
    return api.patch<MessageResponse>(`/api/users/${id}/role`, { role })
  },
  ban(id: number, reason?: string) {
    return api.post<MessageResponse>(`/api/users/${id}/ban`, { reason: reason || '' })
  },
  unban(id: number) {
    return api.post<MessageResponse>(`/api/users/${id}/unban`)
  },
}

// ========== 用户人设 API ==========
export const personaApi = {
  list() {
    return api.get<ListResponse<Persona>>('/api/personas')
  },
  create(data: PersonaInput) {
    return api.post<CreateResponse>('/api/personas', data)
  },
  update(id: number, data: PersonaInput) {
    return api.put<MessageResponse>(`/api/personas/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/personas/${id}`)
  },
}

// ========== 世界书 API ==========
export const worldBookApi = {
  list() {
    return api.get<ListResponse<WorldBookEntry>>('/api/worldbook')
  },
  get(id: number) {
    return api.get<WorldBookEntry>(`/api/worldbook/${id}`)
  },
  create(data: WorldBookInput) {
    return api.post<CreateResponse>('/api/worldbook', data)
  },
  update(id: number, data: WorldBookInput) {
    return api.put<MessageResponse>(`/api/worldbook/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/worldbook/${id}`)
  },
}

// ========== 应用设置 API ==========
export const settingsApi = {
  get() {
    return api.get<{ data: AppSettings }>('/api/settings')
  },
  update(data: AppSettings) {
    return api.put<MessageResponse>('/api/settings', data)
  },
}

// ========== 餐厅 API ==========
export const restaurantApi = {
  list(category?: string) {
    const p: Record<string, string> = {}
    if (category) p.category = category
    return api.get<ListResponse<Restaurant>>('/api/restaurants', p)
  },
  get(id: number) {
    return api.get<Restaurant>(`/api/restaurants/${id}`)
  },
  create(data: RestaurantInput) {
    return api.post<CreateResponse>('/api/restaurants', data)
  },
}

// ========== 订单 API ==========
export const orderApi = {
  list() {
    return api.get<ListResponse<Order>>('/api/orders')
  },
  get(id: number) {
    return api.get<Order>(`/api/orders/${id}`)
  },
  create(data: OrderInput) {
    return api.post<CreateResponse>('/api/orders', data)
  },
  updateStatus(id: number, status: string) {
    return api.patch<MessageResponse>(`/api/orders/${id}/status`, { status })
  },
}

// ========== 预设 API ==========
export const presetApi = {
  list() {
    return api.get<ListResponse<Preset>>('/api/presets')
  },
  get(id: number) {
    return api.get<Preset>(`/api/presets/${id}`)
  },
  create(data: PresetInput) {
    return api.post<CreateResponse>('/api/presets', data)
  },
  update(id: number, data: PresetInput) {
    return api.put<MessageResponse>(`/api/presets/${id}`, data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/presets/${id}`)
  },
}

// ========== 通话记录 API ==========
export const callApi = {
  list() {
    return api.get<ListResponse<CallRecord>>('/api/calls')
  },
  create(data: CallRecordInput) {
    return api.post<CreateResponse>('/api/calls', data)
  },
  delete(id: number) {
    return api.delete<MessageResponse>(`/api/calls/${id}`)
  },
  clear() {
    return api.delete<MessageResponse>('/api/calls')
  },
}

// ========== 短信 API ==========
export const smsApi = {
  listThreads() {
    return api.get<ListResponse<SmsThread>>('/api/sms/threads')
  },
  createThread(data: SmsThreadInput) {
    return api.post<CreateResponse>('/api/sms/threads', data)
  },
  deleteThread(id: number) {
    return api.delete<MessageResponse>(`/api/sms/threads/${id}`)
  },
  listMessages(threadId: number) {
    return api.get<ListResponse<SmsMessage>>(`/api/sms/threads/${threadId}/messages`)
  },
  createMessage(threadId: number, data: SmsMessageInput) {
    return api.post<CreateResponse>(`/api/sms/threads/${threadId}/messages`, data)
  },
}

// ========== 钱包 API ==========
export const walletApi = {
  get() {
    return api.get<WalletInfo>('/api/wallet')
  },
  setBalance(balance: number) {
    return api.put<WalletInfo>('/api/wallet/balance', { balance })
  },
  listTransactions() {
    return api.get<ListResponse<WalletTransaction>>('/api/wallet/transactions')
  },
  createTransaction(data: WalletTransactionInput) {
    return api.post<{ id: number; balance: number; message: string }>('/api/wallet/transaction', data)
  },
}

// ========== 游戏记录 API ==========
export const gameApi = {
  listRecords() {
    return api.get<ListResponse<GameRecordItem>>('/api/games/records')
  },
  createRecord(data: GameRecordInput) {
    return api.post<{ id: number; message: string }>('/api/games/records', data)
  },
  clearRecords() {
    return api.delete<{ message: string }>('/api/games/records')
  },
}

// ========== 购物 API ==========
export const shoppingApi = {
  listProducts(category?: string) {
    const p: Record<string, string> = {}
    if (category) p.category = category
    return api.get<ListResponse<Product>>('/api/shopping/products', p)
  },
  getProduct(id: number) {
    return api.get<Product>(`/api/shopping/products/${id}`)
  },
  listCart() {
    return api.get<ListResponse<CartItem>>('/api/shopping/cart')
  },
  addToCart(data: { product_id: number; quantity: number }) {
    return api.post<CreateResponse>('/api/shopping/cart', data)
  },
  updateCartItem(id: number, quantity: number) {
    return api.put<MessageResponse>(`/api/shopping/cart/${id}`, { quantity })
  },
  removeFromCart(id: number) {
    return api.delete<MessageResponse>(`/api/shopping/cart/${id}`)
  },
  clearCart() {
    return api.delete<MessageResponse>('/api/shopping/cart')
  },
  listOrders() {
    return api.get<ListResponse<ShoppingOrder>>('/api/shopping/orders')
  },
  createOrder(data: { items: ShoppingOrderItem[]; total: number; order_no: string }) {
    return api.post<CreateResponse>('/api/shopping/orders', data)
  },
  listFavorites() {
    return api.get<ListResponse<number>>('/api/shopping/favorites')
  },
  addFavorite(productId: number) {
    return api.post<MessageResponse>('/api/shopping/favorites', { product_id: productId })
  },
  removeFavorite(productId: number) {
    return api.delete<MessageResponse>(`/api/shopping/favorites/${productId}`)
  },
}

// ========== 情侣空间 API ==========
export const coupleApi = {
  getSpace() {
    return api.get<CoupleSpace>('/api/couple/space')
  },
  updateSpace(data: { partner_a?: string; partner_b?: string; start_date?: string }) {
    return api.put<MessageResponse>('/api/couple/space', data)
  },
  listAnniversaries() {
    return api.get<ListResponse<Anniversary>>('/api/couple/anniversaries')
  },
  createAnniversary(data: { title: string; date: string; icon: string }) {
    return api.post<CreateResponse>('/api/couple/anniversaries', data)
  },
  deleteAnniversary(id: number) {
    return api.delete<MessageResponse>(`/api/couple/anniversaries/${id}`)
  },
  listPhotos() {
    return api.get<ListResponse<CouplePhoto>>('/api/couple/photos')
  },
  createPhoto(data: { date: string; url?: string }) {
    return api.post<CreateResponse>('/api/couple/photos', data)
  },
  deletePhoto(id: number) {
    return api.delete<MessageResponse>(`/api/couple/photos/${id}`)
  },
  listTasks() {
    return api.get<ListResponse<CoupleTask>>('/api/couple/tasks')
  },
  createTask(data: { text: string }) {
    return api.post<CreateResponse>('/api/couple/tasks', data)
  },
  toggleTask(id: number) {
    return api.put<MessageResponse>(`/api/couple/tasks/${id}/toggle`, {})
  },
  deleteTask(id: number) {
    return api.delete<MessageResponse>(`/api/couple/tasks/${id}`)
  },
}
