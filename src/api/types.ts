// ========== 通用 ==========
export interface ListResponse<T> {
  data: T[]
  message?: string
}

export interface CreateResponse {
  id: number
  message: string
}

export interface MessageResponse {
  message: string
}

// ========== 用户 ==========
export interface User {
  id: number
  discord_id: string
  username: string
  display_name: string
  avatar_url: string
  role: string
  created_at: string
  updated_at: string
}

// ========== 角色 ==========
export interface Character {
  id: number
  user_id: number
  name: string
  avatar_url: string
  description: string
  personality: string
  system_prompt: string
  greeting: string
  is_public: boolean
  tags: string[]
  created_at: string
  updated_at: string
}

export interface CharacterInput {
  name: string
  avatar_url: string
  description: string
  personality: string
  system_prompt: string
  greeting: string
  is_public: boolean
  tags: string[]
}

// ========== 对话 ==========
export interface Conversation {
  id: number
  user_id: number
  character_id: number
  title: string
  is_group: boolean
  last_message: string
  last_at: string
  created_at: string
  updated_at: string
  character?: {
    name: string
    avatar_url: string
  }
}

export interface ConversationInput {
  character_id: number
  title: string
  is_group?: boolean
}

// ========== 消息 ==========
export interface Message {
  id: number
  conversation_id: number
  role: string
  content: string
  msg_type: string
  created_at: string
}

export interface MessageInput {
  role?: string
  content: string
  msg_type?: string
}

// ========== 帖子 ==========
export interface Post {
  id: number
  user_id: number
  board: string
  title: string
  content: string
  images: string[]
  likes: number
  comments: number
  pinned: boolean
  created_at: string
  updated_at: string
  author?: {
    username: string
    avatar_url: string
  }
}

export interface PostInput {
  board: string
  title: string
  content: string
  images?: string[]
}

// ========== 评论 ==========
export interface Comment {
  id: number
  post_id: number
  user_id: number
  parent_id?: number
  content: string
  likes: number
  created_at: string
  author_name: string
  author_avatar: string
}

export interface CommentInput {
  content: string
  parent_id?: number
}

// ========== 微博 ==========
export interface WeiboPost {
  id: number
  user_id: number
  content: string
  images: string[]
  likes: number
  reposts: number
  comments: number
  is_hot: boolean
  created_at: string
  author_name: string
  author_avatar: string
}

export interface WeiboInput {
  content: string
  images?: string[]
}

// ========== 朋友圈 ==========
export interface Moment {
  id: number
  user_id: number
  content: string
  images: string[]
  location: string
  visibility: string
  likes: number
  created_at: string
  author_name: string
  author_avatar: string
}

export interface MomentInput {
  content: string
  images?: string[]
  location?: string
  visibility?: string
}

// ========== 日记 ==========
export interface Diary {
  id: number
  user_id: number
  title: string
  content: string
  mood: string
  weather: string
  tags: string[]
  is_private: boolean
  created_at: string
  updated_at: string
}

export interface DiaryInput {
  title: string
  content: string
  mood: string
  weather: string
  tags?: string[]
  is_private?: boolean
}

// ========== 用户人设 ==========
export interface Persona {
  id: number
  user_id: number
  name: string
  description: string
  avatar_url: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface PersonaInput {
  name: string
  description: string
  avatar_url?: string
  is_default?: boolean
}

// ========== 世界书 ==========
export interface WorldBookEntry {
  id: number
  user_id: number
  key: string
  content: string
  keywords: string[]
  is_enabled: boolean
  priority: number
  created_at: string
  updated_at: string
}

export interface WorldBookInput {
  key: string
  content: string
  keywords?: string[]
  is_enabled?: boolean
  priority?: number
}

// ========== 应用设置 ==========
export interface AppSettings {
  [key: string]: unknown
}

// ========== 餐厅 ==========
export interface Restaurant {
  id: number
  name: string
  category: string
  image_url: string
  rating: number
  delivery_time: string
  min_order: number
  items: MenuItem[]
  created_at: string
}

export interface MenuItem {
  name: string
  price: number
  image_url?: string
  description?: string
}

export interface RestaurantInput {
  name: string
  category: string
  image_url?: string
  rating?: number
  delivery_time?: string
  min_order?: number
  items?: MenuItem[]
}

// ========== 订单 ==========
export interface Order {
  id: number
  user_id: number
  restaurant_id: number
  items: OrderItem[]
  total: number
  status: string
  address: string
  note: string
  created_at: string
  updated_at: string
  restaurant_name: string
  restaurant_image: string
}

export interface OrderItem {
  name: string
  price: number
  quantity: number
}

export interface OrderInput {
  restaurant_id: number
  items: OrderItem[]
  total: number
  address: string
  note?: string
}
