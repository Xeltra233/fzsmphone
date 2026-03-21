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

### 当前重点能力

- **公开系统设置**: 登录页可通过公开接口读取应用名、公告和 favicon，未登录状态也能正确显示自定义站点信息
- **综合排行榜**: 提供 `/leaderboard` 页面与 `/api/leaderboard` 接口，按用户 Token 消耗展示前 50 名
- **签到与额度系统**: 支持每日签到、连续签到加成、邀请码奖励、兑换码兑换、模型显示名称映射
- **OAuth 配置**: 支持在后台配置 Discord OAuth（Client ID、Client Secret、Redirect URI）
- **多供应商 API**: 支持 Gemini、OpenAI、OpenRouter、DeepSeek 及自定义兼容接口
- **后台管理**: 管理员可查看用户、封禁/解封；超级管理员可调整额度、修改系统设置、配置全局 API

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
| admin (管理员) | 可查看管理区、管理用户、封禁/解封普通用户 |
| user (普通用户) | 只能使用应用，设置个人 API 配置 |

### API 配置说明

- **全局 API 设置** (仅超级管理员可配置)
  - 设置后影响整个站点的所有用户
  - 用户未配置个人 API 时使用全局配置
  - 可维护默认模型、可用模型列表，并为每个模型配置显示名称
  - 支持通过当前 API Key + API 地址一键拉取上游模型列表

- **个人 API 设置** (所有用户可配置)
  - 仅影响本人
  - 优先级高于全局配置
  - 用户设置页会优先显示后台维护的模型显示名称

## 环境变量

> 提示：生产环境下后端会自动执行数据库迁移。若历史数据存在重复用户名、重复邮箱或重复邀请码，新增唯一索引迁移会失败，需要先清洗数据。

### 必需变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `JWT_SECRET` | JWT 签名密钥，必须设置一个随机字符串 | `your-random-secret-key-min-32-chars` |
| `DATABASE_URL` | PostgreSQL 数据库连接字符串 | `postgres://user:pass@host:5432/db?sslmode=disable` |

> ⚠️ **重要**: 如果没有设置 `JWT_SECRET`，应用将无法启动！

### 后端 (Go) 完整配置

```env
# 基础配置 (必需)
PORT=8080                  # 服务器端口
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable  # 数据库连接
JWT_SECRET=your-secret-key # JWT 签名密钥
FRONTEND_URL=http://localhost:3000  # 前端地址
ENVIRONMENT=development    # 环境：development 或 production

# Discord OAuth2 登录 (可选)
DISCORD_CLIENT_ID=你的Discord应用Client ID
DISCORD_CLIENT_SECRET=你的Discord应用Client Secret
DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
DISCORD_GUILD_ID=你的Discord服务器ID
DISCORD_BOT_TOKEN=你的Discord Bot Token
DISCORD_REQUIRED_ROLE_IDS=必需角色ID1,必需角色ID2
DISCORD_ROLE_MATCH=all    # 角色匹配模式：all=需全部拥有，any=拥有任一即可
DISCORD_AUTH_FAIL_MESSAGE=请先加入指定服务器并完成身份组验证
DISCORD_ADMIN_IDS=管理员Discord ID1,管理员Discord ID2
```

### 前端 (Vite)

```env
VITE_API_URL=http://localhost:8080
VITE_DISCORD_CLIENT_ID=your-client-id
VITE_DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
```

### Docker 部署 (单端口)

本项目支持**单端口部署**，后端服务同时提供 API 和前端静态页面。

**必需的环境变量**：

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `ENVIRONMENT` | 设为 `production` 开启单端口模式 | `production` |
| `JWT_SECRET` | JWT 签名密钥 | `your-random-secret` |
| `DATABASE_URL` | PostgreSQL 连接字符串 | `postgres://...` |

设置 `ENVIRONMENT=production` 后，后端会自动从 `./dist` 文件夹提供前端页面，只需 **8080 端口**即可访问完整应用。

### 如何获取必需的环境变量

#### JWT_SECRET
生成一个随机的密钥，可以使用以下命令生成：
```bash
openssl rand -base64 32
```
或使用在线随机字符串生成器。

#### DATABASE_URL
如果使用 Docker Compose 部署，数据库连接字符串为：
```
postgres://fzsmphone:fzsmphone_dev_123@postgres:5432/fzsmphone?sslmode=disable
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
- 本地开发：打开 http://localhost:3000
- 生产部署：打开 http://your-domain:8080
- 使用邮箱注册或登录
- 或使用 Discord 账号登录

### 数据库迁移注意事项

项目当前会为以下字段建立唯一约束：

- `LOWER(username)`
- `LOWER(email)`，但排除空字符串
- `invite_code`，但排除空字符串

如果线上已有脏数据，建议先执行以下 SQL 查重：

```sql
SELECT LOWER(username) AS normalized_username, COUNT(*) AS cnt, ARRAY_AGG(id ORDER BY id) AS user_ids
FROM users
GROUP BY LOWER(username)
HAVING COUNT(*) > 1;

SELECT LOWER(email) AS normalized_email, COUNT(*) AS cnt, ARRAY_AGG(id ORDER BY id) AS user_ids
FROM users
WHERE email <> ''
GROUP BY LOWER(email)
HAVING COUNT(*) > 1;

SELECT invite_code, COUNT(*) AS cnt, ARRAY_AGG(id ORDER BY id) AS user_ids
FROM users
WHERE invite_code <> ''
GROUP BY invite_code
HAVING COUNT(*) > 1;
```

### Docker 部署

```bash
docker-compose up -d
```

## 使用教程

### 1. 注册与登录

访问应用后，你可以选择：
- **账号登录**: 输入邮箱/用户名和密码登录
- **立即注册**: 点击注册链接，填写用户名、邮箱、密码，可选填写邀请码
- **Discord 登录**: 如已配置 Discord OAuth2

登录页会读取公开设置接口，因此未登录时也能显示自定义应用名、公告和站点图标。

### 2. 超级管理员设置

当前版本没有前台“设置超级管理员”入口。

如需将用户提升为超级管理员，请直接在数据库中执行：

```sql
UPDATE users SET is_super_admin = true WHERE username = 'yourUsername';
```

如果希望该用户在前端权限判断里也完全一致，建议同时同步角色字段：

```sql
UPDATE users
SET role = 'super_admin', is_super_admin = true
WHERE username = 'yourUsername';
```

> ⚠️ 不建议直接修改现有管理员或超级管理员账号，操作前请先备份数据库。

### 3. 配置 API

#### 超级管理员配置全局 API

1. 登录超级管理员账号
2. 进入 `/admin/settings` 系统设置
3. 切换到 "API设置" 标签
4. 在"全局API设置"中填入：
   - **API Key**: 你的 API Key
   - **API URL**: API 地址 (如 `https://api.openai.com/v1/chat/completions`)
   - **默认模型**: 模型名称 (如 `gpt-4o-mini`)
   - **可用模型列表**: 可手动添加，也可点击“拉取模型”从上游接口导入
   - **显示名称**: 可为每个模型单独设置前台展示名，例如把 `gpt-4o-mini` 显示为 `GPT 4o Mini`
5. 点击保存

保存后：

- 用户设置页的模型下拉框会优先显示你配置的模型显示名称
- 聊天页的信息面板会显示当前模型的展示名，并保留原始模型 ID 作为补充
- 若用户未配置主模型，前台会优先使用后台设置的默认模型

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
- 封禁/解封用户

超级管理员额外可以：
- 修改用户角色
- 调整用户额度
- 修改全局系统设置
- 配置全局 API 与 OAuth 参数

### 6. 排行榜与额度

- `/leaderboard`：展示基于 `total_tokens` 的综合排行榜
- `/credits`：展示剩余额度、累计 Token、连续签到、邀请码与兑换码入口
- 签到按钮会根据系统开关与当天是否已签到显示 `立即签到`、`已签到` 或 `未开启`

### 7. 数据备份与导入

用户可在设置页的“数据”标签中使用一键导出/导入。

导出内容包含：

- 角色卡
- 预设
- 用户人设
- 日记
- 聊天会话与聊天记录本地缓存
- 短信、通话、钱包、游戏记录
- 世界书、本地社交数据、模板变量、图片生成配置与普通设置

导出内容不包含：

- `apiKey`
- `socialApiKey`
- 生图服务密钥（如 `novelai.key`、`openai.key`、`gemini.key`）
- 登录 token、用户认证态等敏感信息

导入行为说明：

- 导入前会先显示预览，估算可新增数量与重复数量
- 导入会优先合并本地数据，而不是直接覆盖整个浏览器存储
- 对角色卡、预设、人设、日记、短信、通话、游戏记录会尽量按内容签名去重
- 同一份备份重复导入时，系统会尽量跳过已存在内容，但不能保证 100% 无重复

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

本项目仅供个人学习交流使用，**禁止任何形式的商业用途**。
你可以自由使用、修改、分发本项目，但不得用于商业目的。
如需商用请联系作者获取授权。

## 贡献者

感谢以下 contributors 的贡献：

<a href="https://github.com/Xeltra233/fzsmphone/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Xeltra233/fzsmphone&max=12" alt="contributors" />
</a>

## 合作者

| 头像 | 名称 | 角色 |
|------|------|------|
| <img src="https://avatars.githubusercontent.com/u/243946305?v=4" width="40" height="40" style="border-radius: 50%;" /> | Xeltra233 | 项目创始人 |
| <img src="https://avatars.githubusercontent.com/u/234566466?v=4" width="40" height="40" style="border-radius: 50%;" /> | 907606864-hub | 合作开发者 |
| <img src="https://avatars.githubusercontent.com/u/268621292?v=4" width="40" height="40" style="border-radius: 50%;" /> | dog-the-king | 合作开发者 |

---

<p align="center">
  贩子死妈小手机版权所有 © 2026
</p>
