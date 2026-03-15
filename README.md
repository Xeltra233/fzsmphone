# fzsmphone

高仿 iOS 风格手机模拟器 Web 应用，专为角色扮演和沉浸式社交模拟打造。集成 AI 生成内容，覆盖聊天、社交、娱乐、生活、工具五大模块，45+ 功能页面。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Pinia |
| 后端 | Go 1.22 + Chi v5 + PostgreSQL 16 |
| AI | OpenAI 兼容 API（流式 SSE） |
| 认证 | Discord OAuth2 / 邮箱密码 + JWT |
| 部署 | Docker / Zeabur 一键部署 |
| 图标 | CSS 渐变图标系统（GradientIcon） |

## 功能特性

- **聊天模块**: AI 对话、角色创建、人格设定、世界书
- **社交模块**: 微博、小红书、知乎、QQ 空间、Discord
- **娱乐模块**: 抖音、直播、在线游戏、迷你剧场
- **生活模块**: 外卖、购物、钱包、日记
- **工具模块**: 电话、短信、邮箱、浏览器、地图

## 部署

### Zeabur 一键部署

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/7PNREM)

## 认证与用户系统

### 登录方式

系统支持两种登录方式：

1. **邮箱/用户名 + 密码登录**
   - 用户可自行注册账号
   - 支持邮箱或用户名登录
   - 密码使用 bcrypt 加密存储

2. **Discord OAuth2 登录**
   - 需要配置 Discord 应用
   - 可设置服务器成员验证
   - 可设置必需角色

### 用户分级

| 角色 | 权限 |
|------|------|
| super_admin (超级管理员) | 最高权限，可设置全局 API 配置影响整个站点，可管理所有用户 |
| admin (管理员) | 可管理用户、查看系统设置 |
| user (普通用户) | 只能使用应用，设置个人 API 配置 |

### API 配置说明

- **全局 API 设置** (仅超级管理员可配置)
  - 设置后影响整个站点的所有用户
  - 用户未配置个人 API 时使用全局配置

- **个人 API 设置** (所有用户可配置)
  - 仅影响本人
  - 优先级高于全局配置

## 环境变量

### Zeabur 部署必需变量

在 Zeabur 控制台的 **Environment Variables** 中必须设置以下变量：

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `JWT_SECRET` | JWT 签名密钥，必须设置一个随机字符串 | `your-random-secret-key-min-32-chars` |
| `DATABASE_URL` | PostgreSQL 数据库连接字符串 | `postgres://user:pass@host:5432/db?sslmode=disable` |

> ⚠️ **重要**: 如果没有设置 `JWT_SECRET`，应用将无法启动！

### 后端 (Go) 完整配置

```env
# 基础配置 (必需)
PORT=8080
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
JWT_SECRET=your-secret-key
FRONTEND_URL=http://localhost:3000
ENVIRONMENT=development

# Discord OAuth2 (可选)
DISCORD_CLIENT_ID=your-client-id
DISCORD_CLIENT_SECRET=your-client-secret
DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
DISCORD_GUILD_ID=your-guild-id
DISCORD_BOT_TOKEN=your-bot-token
DISCORD_REQUIRED_ROLE_IDS=role1,role2
DISCORD_ROLE_MATCH=all
DISCORD_AUTH_FAIL_MESSAGE=请先加入指定服务器并完成身份组验证
DISCORD_ADMIN_IDS=admin-discord-id1,admin-discord-id2
```

### 前端 (Vite)

```env
VITE_API_URL=http://localhost:8080
VITE_DISCORD_CLIENT_ID=your-client-id
VITE_DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
```

### 如何获取必需的环境变量

#### JWT_SECRET
生成一个随机的密钥，可以使用以下命令生成：
```bash
openssl rand -base64 32
```
或使用在线随机字符串生成器。

#### DATABASE_URL
如果使用 Zeabur 内置的 PostgreSQL 数据库，在 Zeabur 控制台的 **Variables** 中可以找到 `POSTGRES_URL`，格式为：
```
postgres://username:password@host:5432/databasename?sslmode=disable
```

## 快速开始

### 本地开发

1. **克隆项目**
   ```bash
   git clone https://github.com/Xeltra233/fzsmphone.git
   cd fzsmphone
   ```

2. **启动后端**
   ```bash
   cd server
   cp .env.example .env
   # 编辑 .env 配置数据库和 JWT 密钥
   go run cmd/api/main.go
   ```

3. **启动前端**
   ```bash
   npm install
   npm run dev
   ```

4. **访问应用**
   - 打开 http://localhost:3000
   - 使用邮箱注册或登录
   - 或使用 Discord 账号登录

### Docker 部署

```bash
docker-compose up -d
```

## 使用教程

### 1. 注册与登录

访问应用后，你可以选择：
- **账号登录**: 输入邮箱/用户名和密码登录
- **立即注册**: 点击注册链接，填写用户名、邮箱、密码
- **Discord 登录**: 如已配置 Discord OAuth2

### 2. 超级管理员设置

如果你需要将用户提升为超级管理员：

1. 使用管理员账号登录
2. 进入 `/admin/users` 用户管理页面
3. 找到目标用户，点击设置超级管理员

或者直接在数据库中执行：
```sql
UPDATE users SET is_super_admin = true WHERE username = 'yourUsername';
```

### 3. 配置 API

#### 超级管理员配置全局 API

1. 登录超级管理员账号
2. 进入 `/admin/settings` 系统设置
3. 切换到 "API设置" 标签
4. 在"全局API设置"中填入：
   - **API Key**: 你的 API Key
   - **API URL**: API 地址 (如 `https://api.openai.com/v1/chat/completions`)
   - **默认模型**: 模型名称 (如 `gpt-4o-mini`)
5. 点击保存

#### 普通用户配置个人 API

1. 登录账号
2. 进入 `/admin/settings` 系统设置
3. 切换到 "API设置" 标签
4. 在"个人API设置"中填入自己的 API 配置
5. 点击保存

### 4. 系统配置 (仅超级管理员)

超级管理员可以修改系统配置：

1. 进入 `/admin/settings`
2. 切换到"系统配置"标签
3. 添加或修改配置项
4. 点击保存

### 5. 用户管理 (管理员)

管理员可以：
- 查看所有用户
- 修改用户角色
- 封禁/解封用户

## 目录结构

```
fzsmphone/
├── src/                    # Vue 前端
│   ├── api/               # API 客户端
│   ├── components/        # 公共组件
│   ├── layouts/           # 布局
│   ├── pages/             # 页面组件
│   ├── router/            # 路由
│   ├── stores/            # Pinia 状态管理
│   ├── styles/            # 样式
│   └── utils/             # 工具函数
├── server/                # Go 后端
│   ├── cmd/              # 入口文件
│   └── internal/         # 内部包
│       ├── config/       # 配置
│       ├── database/     # 数据库
│       ├── handlers/     # 处理器
│       ├── middleware/   # 中间件
│       ├── models/       # 数据模型
│       ├── router/       # 路由
│       └── ws/           # WebSocket
├── public/               # 静态资源
└── package.json          # 前端依赖
```

## 许可证

MIT License