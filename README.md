# fzsmphone

> iOS 风格手机模拟器 Web 应用 — 集聊天、社交、工具、生活、娱乐于一体的全功能平台

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite |
| 状态管理 | Pinia (15 stores) |
| 后端 | Go 1.22 + Chi v5 |
| 数据库 | PostgreSQL 16 + pgx/v5 |
| 实时通信 | WebSocket (nhooyr.io/websocket) |
| 认证 | Discord OAuth + JWT |
| 容器化 | Docker Compose |

## 项目结构

```
fzsmphone/
├── src/                          # Vue 3 前端
│   ├── api/                      # API 客户端 & 服务层
│   │   ├── client.ts             # 统一 HTTP 客户端
│   │   ├── services.ts           # 12 个领域 API 模块
│   │   └── types.ts              # TypeScript 类型定义
│   ├── components/               # 可复用组件
│   │   ├── common/               # NavBar, PlaceholderPage
│   │   └── phone/                # StatusBar, DynamicIsland
│   ├── layouts/                  # PhoneLayout (iPhone 外壳)
│   ├── pages/                    # 35+ 页面
│   │   ├── auth/                 # 登录 & OAuth 回调
│   │   ├── chat/                 # 聊天、角色、群聊、人设、世界书
│   │   ├── entertainment/        # 一起听歌、游戏、直播、赌场、迷你剧场
│   │   ├── home/                 # 主屏幕
│   │   ├── lifestyle/            # 外卖、餐厅、订单、钱包、购物
│   │   ├── social/               # 论坛、微博、QQ空间、帖子详情、情侣空间
│   │   └── tools/                # 日记、电话、短信、股票、汇率、通话、偷看手机、线下约会
│   ├── router/                   # Vue Router (35+ 路由 + 认证守卫)
│   ├── stores/                   # Pinia 状态管理 (15 stores)
│   └── styles/                   # CSS 变量主题 + 动画
├── server/                       # Go 后端
│   ├── cmd/api/                  # 入口 main.go
│   └── internal/
│       ├── config/               # 环境配置
│       ├── database/             # PostgreSQL 连接 + 14 SQL 迁移
│       ├── handlers/             # 11 个 HTTP 处理器
│       ├── middleware/           # JWT 认证 + 响应助手
│       ├── models/               # 数据模型
│       ├── router/               # Chi 路由注册
│       └── ws/                   # WebSocket Hub
├── docker-compose.yml            # PostgreSQL + API 容器
├── .env                          # 环境变量 (不提交)
└── package.json                  # 前端依赖
```

## 功能模块 (35+)

### 💬 聊天
- 好友列表 / 私聊 / 群聊
- AI 角色管理（创建/编辑/删除）
- 用户人设系统
- 世界书（Lore Book）

### 🌐 社交
- 论坛（发帖/评论/搜索）
- 微博（发布/时间线）
- QQ 空间（动态流）
- 帖子详情 & 评论
- 情侣空间

### 🛠 工具
- 日记本
- 电话 / 短信
- 股票行情
- 汇率换算
- 语音通话 / 视频通话
- 偷看手机 / 反向偷看
- 线下约会规划

### 🍜 生活
- 外卖点餐 / 餐厅浏览
- 订单管理
- 钱包
- 购物商城

### 🎮 娱乐
- 一起听歌
- 小游戏合集
- 直播间
- 赌场
- 迷你剧场

### ⚙️ 系统
- 个人主页 & 资料编辑
- 自定义主题
- 预设管理
- 管理后台（内嵌手机 UI）

## 快速开始

### 环境要求

- Node.js 18+
- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose（推荐）

### 1. 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone <repo-url> fzsmphone
cd fzsmphone

# 创建环境变量
cp server/.env.example .env
# 编辑 .env 填入 Discord OAuth 等配置

# 启动 PostgreSQL + 后端
docker-compose up -d

# 安装前端依赖 & 启动开发服务器
npm install
npm run dev
```

前端: http://localhost:3000  
后端 API: http://localhost:8080

### 2. 本地开发（不用 Docker）

```bash
# 启动 PostgreSQL（确保已安装并运行）
# 创建数据库: fzsmphone

# 后端
cd server
cp .env.example .env
# 编辑 .env
go run ./cmd/api

# 前端（新终端）
cd ..
npm install
npm run dev
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `PORT` | 后端端口 | `8080` |
| `DATABASE_URL` | PostgreSQL 连接串 | - |
| `JWT_SECRET` | JWT 签名密钥 | - |
| `DISCORD_CLIENT_ID` | Discord OAuth 客户端 ID | - |
| `DISCORD_CLIENT_SECRET` | Discord OAuth 客户端密钥 | - |
| `DISCORD_REDIRECT_URI` | Discord OAuth 回调地址 | - |
| `FRONTEND_URL` | 前端 URL（CORS） | `http://localhost:3000` |
| `ENVIRONMENT` | 运行环境 | `development` |

## API 端点

### 认证
- `POST /api/auth/discord` — Discord OAuth 回调
- `GET /api/auth/me` — 获取当前用户

### 聊天
- `GET/POST /api/characters` — 角色列表/创建
- `GET/PUT/DELETE /api/characters/:id` — 角色详情/更新/删除
- `GET/POST /api/conversations` — 对话列表/创建
- `GET/DELETE /api/conversations/:id` — 对话详情/删除
- `GET/POST /api/conversations/:id/messages` — 消息列表/发送

### 社交
- `GET/POST /api/posts` — 帖子列表/创建
- `GET/PUT/DELETE /api/posts/:id` — 帖子详情/更新/删除
- `GET/POST /api/posts/:id/comments` — 评论列表/创建
- `GET/POST /api/weibo` — 微博列表/创建
- `GET/DELETE /api/weibo/:id` — 微博详情/删除
- `GET/POST /api/moments` — 动态列表/创建
- `DELETE /api/moments/:id` — 删除动态

### 工具
- `GET/POST /api/diaries` — 日记列表/创建
- `GET/PUT/DELETE /api/diaries/:id` — 日记详情/更新/删除
- `GET/POST /api/personas` — 人设列表/创建
- `PUT/DELETE /api/personas/:id` — 人设更新/删除
- `GET/POST /api/worldbook` — 世界书列表/创建
- `GET/PUT/DELETE /api/worldbook/:id` — 世界书详情/更新/删除

### 生活
- `GET/POST /api/restaurants` — 餐厅列表/创建
- `GET /api/restaurants/:id` — 餐厅详情
- `GET/POST /api/orders` — 订单列表/创建
- `GET /api/orders/:id` — 订单详情
- `PATCH /api/orders/:id/status` — 更新订单状态

### 设置 & 用户
- `GET/PUT /api/settings` — 获取/更新设置
- `GET /api/users` — 用户列表
- `GET /api/users/:id` — 用户详情
- `PATCH /api/users/:id/role` — 更新用户角色

### WebSocket
- `GET /api/ws` — WebSocket 连接

## 构建 & 部署

```bash
# 前端构建
npm run build
# 产出: dist/

# 后端构建
cd server
go build -o fzsmphone ./cmd/api

# Docker 部署
docker-compose -f docker-compose.yml up -d
```

## 主题系统

支持亮色/暗色主题，通过 CSS 变量实现：

```css
:root {
  --bg-primary: #ffffff;
  --text-primary: #1d1d1f;
  --accent: #007aff;
  /* ... 更多变量 */
}

[data-theme="dark"] {
  --bg-primary: #000000;
  --text-primary: #f5f5f7;
  --accent: #0a84ff;
}
```

## License

MIT
