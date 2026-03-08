/**
 * 统一应用注册表
 * 所有应用的定义、路由映射、feature gating 都从这里读取
 */

export interface AppDefinition {
  /** 显示名称 */
  name: string
  /** 图标 key，对应 appIcons */
  iconKey: string
  /** 背景渐变色 */
  color: string
  /** 路由路径 */
  route: string
  /** 路由名称（用于 feature gating） */
  routeName: string
  /** 功能开关 ID */
  featureId?: string
  /** 角标数 */
  badge?: number
  /** 仅管理员可见 */
  adminOnly?: boolean
  /** 是否在 Dock 栏显示 */
  dock?: boolean
  /** Dock 栏排序（数字越小越靠前） */
  dockOrder?: number
}

/**
 * 所有应用定义（主屏幕 + Dock）
 */
export const appRegistry: AppDefinition[] = [
  // === 主屏幕应用 ===
  { name: '聊天', iconKey: 'chat', color: 'linear-gradient(135deg, #5B6EF5, #8B5CF6)', route: '/friends', routeName: 'Friends', featureId: 'chat' },
  { name: '微博', iconKey: 'weibo', color: 'linear-gradient(135deg, #E6162D, #FF4757)', route: '/weibo', routeName: 'Weibo', featureId: 'weibo' },
  { name: '朋友圈', iconKey: 'qzone', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/qzone', routeName: 'QZone', featureId: 'qzone' },
  { name: '外卖', iconKey: 'takeaway', color: 'linear-gradient(135deg, #0AB4FF, #26D0CE)', route: '/takeaway', routeName: 'Takeaway', featureId: 'takeaway' },
  { name: '购物', iconKey: 'shopping', color: 'linear-gradient(135deg, #FF6348, #FF4757)', route: '/shopping', routeName: 'Shopping', featureId: 'shopping' },
  { name: '音乐', iconKey: 'music', color: 'linear-gradient(135deg, #FC5C7D, #6A82FB)', route: '/listen-together', routeName: 'ListenTogether', featureId: 'music' },
  { name: '直播', iconKey: 'live', color: 'linear-gradient(135deg, #A855F7, #EC4899)', route: '/live', routeName: 'Live', featureId: 'live' },
  { name: '日记', iconKey: 'diary', color: 'linear-gradient(135deg, #FECA57, #FF9FF3)', route: '/diary', routeName: 'Diary', featureId: 'diary' },
  { name: '游戏', iconKey: 'games', color: 'linear-gradient(135deg, #6C5CE7, #A29BFE)', route: '/games', routeName: 'Games', featureId: 'games' },
  { name: '钱包', iconKey: 'wallet', color: 'linear-gradient(135deg, #00B894, #55EFC4)', route: '/wallet', routeName: 'Wallet', featureId: 'wallet' },
  { name: '影院', iconKey: 'theater', color: 'linear-gradient(135deg, #2D3436, #636E72)', route: '/mini-theater', routeName: 'MiniTheater', featureId: 'mini_theater' },
  { name: '股票', iconKey: 'stock', color: 'linear-gradient(135deg, #00B16A, #1ABC9C)', route: '/stock', routeName: 'Stock', featureId: 'stock' },
  { name: '汇率', iconKey: 'currency', color: 'linear-gradient(135deg, #3498DB, #2980B9)', route: '/currency', routeName: 'CurrencyConverter', featureId: 'currency' },
  { name: '情侣', iconKey: 'couple', color: 'linear-gradient(135deg, #FD79A8, #E84393)', route: '/couple-space', routeName: 'CoupleSpace', featureId: 'couple_space' },
  { name: '知乎', iconKey: 'zhihu', color: 'linear-gradient(135deg, #0066FF, #3399FF)', route: '/zhihu', routeName: 'Zhihu', featureId: 'zhihu' },
  { name: '小红书', iconKey: 'xiaohongshu', color: 'linear-gradient(135deg, #FF2442, #FF6B81)', route: '/xiaohongshu', routeName: 'Xiaohongshu', featureId: 'xiaohongshu' },
  { name: '抖音', iconKey: 'douyin', color: 'linear-gradient(135deg, #25F4EE, #FE2C55)', route: '/douyin', routeName: 'Douyin', featureId: 'douyin' },
  { name: 'Discord', iconKey: 'discord', color: 'linear-gradient(135deg, #5865F2, #7289DA)', route: '/discord', routeName: 'Discord', featureId: 'discord' },
  { name: '邮箱', iconKey: 'email', color: 'linear-gradient(135deg, #4A90D9, #357ABD)', route: '/email', routeName: 'Email', featureId: 'email' },
  { name: '浏览器', iconKey: 'browser', color: 'linear-gradient(135deg, #4285F4, #34A853)', route: '/browser', routeName: 'Browser', featureId: 'browser' },
  { name: '地图', iconKey: 'map', color: 'linear-gradient(135deg, #34A853, #FBBC05)', route: '/map', routeName: 'Map', featureId: 'map' },
  { name: '日历', iconKey: 'calendar', color: 'linear-gradient(135deg, #FF6B6B, #EE5A24)', route: '/calendar', routeName: 'Calendar', featureId: 'calendar' },
  { name: '设置', iconKey: 'settings', color: 'linear-gradient(135deg, #636E72, #B2BEC3)', route: '/customize', routeName: 'Customize' },
  { name: '角色', iconKey: 'characters', color: 'linear-gradient(135deg, #00CEC9, #81ECEC)', route: '/characters', routeName: 'Characters', featureId: 'characters' },
  { name: '人设', iconKey: 'personas', color: 'linear-gradient(135deg, #E17055, #FAB1A0)', route: '/personas', routeName: 'UserPersonas', featureId: 'personas' },
  { name: '世界书', iconKey: 'worldbook', color: 'linear-gradient(135deg, #6C5CE7, #DDA0DD)', route: '/worldbook', routeName: 'WorldBook', featureId: 'worldbook' },
  { name: '预设', iconKey: 'preset', color: 'linear-gradient(135deg, #FDCB6E, #F39C12)', route: '/preset', routeName: 'Preset', featureId: 'preset' },
  { name: '管理', iconKey: 'admin', color: 'linear-gradient(135deg, #E74C3C, #C0392B)', route: '/admin/features', routeName: 'AdminFeatures', adminOnly: true },

  // === Dock 栏应用 ===
  { name: '电话', iconKey: 'phone', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/phone', routeName: 'Phone', featureId: 'phone', dock: true, dockOrder: 1 },
  { name: '短信', iconKey: 'sms', color: 'linear-gradient(135deg, #2ED573, #7BED9F)', route: '/sms', routeName: 'Sms', featureId: 'sms', dock: true, dockOrder: 2 },
  { name: '聊天', iconKey: 'chat', color: 'linear-gradient(135deg, #5B6EF5, #8B5CF6)', route: '/friends', routeName: 'Friends', featureId: 'chat', dock: true, dockOrder: 3 },
  { name: '个人', iconKey: 'profile', color: 'linear-gradient(135deg, #636E72, #B2BEC3)', route: '/profile', routeName: 'Profile', dock: true, dockOrder: 4 },
]

/**
 * 从注册表自动生成 routeName -> featureId 的映射
 * 供路由守卫进行 feature gating
 */
export function buildRouteFeatureMap(): Record<string, string> {
  const map: Record<string, string> = {}
  for (const app of appRegistry) {
    if (app.featureId && app.routeName) {
      map[app.routeName] = app.featureId
    }
  }
  // 额外补充注册表中没有直接出现的路由别名
  map['Chat'] = 'chat'
  map['GroupChat'] = 'chat'
  map['CharacterEdit'] = 'characters'
  map['Moments'] = 'qzone'
  map['Restaurant'] = 'takeaway'
  map['TakeawayOrders'] = 'takeaway'
  map['VoiceCall'] = 'voice_call'
  map['VideoCall'] = 'video_call'
  map['PhonePeek'] = 'phone_peek'
  map['ReversePhonePeek'] = 'reverse_phone_peek'
  map['OfflineDate'] = 'offline_date'
  map['Casino'] = 'casino'
  return map
}

/** 获取主屏幕应用（不含 Dock 专属） */
export function getHomeApps(): AppDefinition[] {
  return appRegistry.filter((app) => !app.dock)
}

/** 获取 Dock 栏应用 */
export function getDockApps(): AppDefinition[] {
  return appRegistry
    .filter((app) => app.dock)
    .sort((a, b) => (a.dockOrder ?? 99) - (b.dockOrder ?? 99))
}
