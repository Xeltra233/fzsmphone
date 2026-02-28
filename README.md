# fzsmphone

高仿 iOS 风格手机模拟器 Web 应用，专为角色扮演和沉浸式社交模拟打造。集成 AI 生成内容，覆盖聊天、社交、娱乐、生活、工具五大模块，45+ 功能页面。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Pinia |
| 后端 | Go 1.22 + Chi v5 + PostgreSQL 16 |
| AI | OpenAI 兼容 API（流式 SSE） |
| 认证 | Discord OAuth2 + JWT |
| 部署 | Docker / Zeabur 一键部署 |
| 图标 | CSS 渐变图标系统（GradientIcon） |

## 功能总览

### 聊天模块

| 功能 | 说明 |
|------|------|
| 好友私聊 | 1v1 AI 角色对话，流式打字机效果 |
| 群聊 | 多角色群组讨论，AI 按角色轮流发言 |
| AI 角色管理 | 创建/编辑角色（人设、系统提示、开场白、示例对话、标签） |
| 用户人设 | 多人设管理与切换 |
| 世界书 | 关键词触发背景知识库，支持角色绑定 |
| 预设系统 | 系统提示词模板，支持 `{{char}}`/`{{user}}` 变量、预填充 |

### 社交模块

| 功能 | 说明 |
|------|------|
| 微博 | 信息流 + 热搜 + 榜单 + 评论互动，AI 生成 |
| 朋友圈 (QZone) | 仿微信朋友圈动态，AI 生成点赞/评论 |
| 论坛 | AI 生成帖子与楼中楼回复，截断容错解析 |
| 小红书 | 笔记流 + 互动，AI 生成内容 |
| 知乎 | 问答 + 热榜，AI 生成回答 |
| 情侣空间 | 纪念日、照片墙、任务清单、心愿单 |

### 娱乐模块

| 功能 | 说明 |
|------|------|
| 赌场 | 老虎机、骰子、幸运转盘，金币系统 |
| 小游戏 | 猜数字等益智小游戏合集 |
| 直播大厅 | 模拟直播间，弹幕互动，送礼系统 |
| 短视频 (抖音) | 上下滑动短视频流 |
| 迷你剧场 | 短剧集播放，评论互动 |
| 一起听 | 音乐共享播放 |

### 生活模块

| 功能 | 说明 |
|------|------|
| 购物商城 | 商品浏览、购物车、收藏、下单 |
| 外卖点餐 | 餐厅列表、菜单、下单、订单追踪 |
| 餐厅 | 餐厅详情、分类筛选 |
| 钱包 | 余额管理、交易记录 |

### 工具模块

| 功能 | 说明 |
|------|------|
| 电话 | 通讯录、通话记录、拨号盘 |
| 短信 | 对话列表、收发消息 |
| 语音通话 | 仿真通话界面 |
| 视频通话 | 仿真视频通话界面 |
| 日记 | 心情日记，心情/天气标记 |
| 手机偷看 | AI 生成隐私发现，角色关联 |
| 反向偷看 | 对方偷看你的手机 |
| 线下约会 | AI 推荐约会计划 |
| 股票行情 | 模拟实时行情，自选股追踪 |
| 汇率换算 | 多币种实时换算 |

### 系统管理

| 功能 | 说明 |
|------|------|
| 自定义设置 | API 配置、主题、深色模式、壁纸、气泡样式、字体 |
| 功能管理 | 管理员功能开关面板，实时启用/禁用模块 |
| 用户管理 | 用户列表、角色分配 |
| 数据统计 | 系统数据概览 |
| 系统设置 | 全局配置管理 |

## 快速开始

```bash
git clone https://github.com/907606864-hub/fzsmphone.git
cd fzsmphone

# 启动数据库
docker-compose up -d postgres

# 后端
cp server/.env.example server/.env    # 编辑填入配置
cd server && go run ./cmd/api &       # 启动在 :8080

# 前端
cd .. && npm install && npm run dev   # 启动在 :3000
```

## 环境变量

| 变量 | 说明 |
|------|------|
| `DATABASE_URL` | PostgreSQL 连接串 |
| `JWT_SECRET` | JWT 签名密钥 |
| `DISCORD_CLIENT_ID` | Discord 应用 ID |
| `DISCORD_CLIENT_SECRET` | Discord 应用密钥 |
| `DISCORD_REDIRECT_URI` | OAuth 回调地址 |
| `DISCORD_GUILD_ID` | 限制登录的服务器 ID |
| `DISCORD_REQUIRED_ROLE_IDS` | 必须拥有的身份组 |
| `DISCORD_ADMIN_IDS` | 管理员 Discord ID |
| `FRONTEND_URL` | 前端地址 (CORS) |

完整列表见 `server/.env.example`

## 项目结构

```
fzsmphone/
  src/
    api/            # API 客户端和类型定义
    components/     # 可复用组件 (NavBar, GradientIcon 等)
    pages/          # 45+ 功能页面
      chat/         # 聊天模块 (6 页面)
      social/       # 社交模块 (7 页面)
      entertainment/# 娱乐模块 (6 页面)
      lifestyle/    # 生活模块 (5 页面)
      tools/        # 工具模块 (10 页面)
      admin/        # 管理模块 (4 页面)
    stores/         # Pinia 状态管理
    utils/          # 工具函数 (AI服务, 图标系统, 社交解析器)
    router/         # Vue Router 路由配置
  server/
    cmd/api/        # Go 主入口
    internal/
      database/     # 数据库连接和迁移
      handlers/     # API 路由处理器
      middleware/   # 认证中间件
      router/       # 路由注册
```

## 部署

项目已配置 `Dockerfile` + `zbpack.json`，支持 Zeabur 一键部署。也可 `docker-compose up -d` 本地启动全部服务。

## License

MIT
