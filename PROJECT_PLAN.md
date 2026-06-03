# Caddy Server Manager - 项目总体规划

> 详细文档索引：
> - [需求分析文档](./docs/01-requirements.md) — 用户角色、功能需求、非功能需求、约束条件
> - [业务设计文档](./docs/02-business-design.md) — 业务流程、状态机、领域模型、业务规则
> - [功能设计文档](./docs/03-functional-design.md) — 页面规格、API 接口契约、Caddyfile 生成模板、数据流

---

## 一、项目现状分析

| 维度 | 状态 |
|------|------|
| **代码阶段** | GoLand 脚手架，仅含 hello-world 的 `main.go` |
| **go.mod** | 裸模块声明，无版本号、无任何依赖 |
| **源文件** | 1 个 (`main.go`)，无任何业务代码 |
| **架构** | 无内部包、无分层、无结构 |
| **模板** | 无任何 HTML 模板文件 |
| **配置** | 无配置文件、无 Caddyfile |
| **数据库** | 无 SQLite 初始化代码、无迁移逻辑 |
| **结论** | **空白项目，需从零搭建** |

---

## 二、设计方向确认

### 选定风格：Build（huashu-design #11 号风格）

**哲学内核**：精致的简单比复杂更难

**核心特征**：
- 奢侈品级留白（70%+ 页面区域）
- 微妙的字重对比（font-weight: 200 ~ 600）
- 单一强调色的战略使用（天青色 **#5DADE2**）
- 呼吸感的节奏、黄金比例

**与 Zolysoft UI 规范的融合**：
| 维度 | 取值 |
|------|------|
| 主色 | `#5DADE2`（天青色） |
| 辅色 | `#3498DB`（深天青，hover态） / `#D4E6F1`（浅天青，背景） |
| 功能色 | 成功 `#10B981` / 警告 `#F59E0B` / 危险 `#EF4444` / 信息 `#3B82F6` |
| 背景 | `#F8FAFB`（接近白，保持呼吸感） |
| 卡片 | `#FFFFFF` + 极柔阴影（Build 风格软阴影） |
| 字重 | 标题 500-600 / 正文 300-400（Build 风格微妙对比） |
| 间距 | 基础 4px，全局组件间统一 24px |
| 字体 | `-apple-system, BlinkMacSystemFont, "PingFang SC", "Microsoft YaHei"` |
| 图标 | 无 emoji，使用 SVG 简约线条图标（Lucide 风格） |

---

## 三、架构设计建议

```
caddy_server/
├── main.go                  # 入口：初始化 DB、配置、信号处理、启动服务器
├── go.mod
├── Dockerfile               # Alpine 镜像构建文件
├── docker-compose.yaml      # 本地开发/部署编排
├── config/
│   ├── config.go            # 配置结构体 + YAML/JSON 解析
│   └── config.yaml          # 默认配置文件（域名、端口、Caddyfile 路径等）
├── database/
│   ├── database.go          # SQLite 连接 + 迁移 + 初始化管理员检测
│   └── models.go            # 数据模型定义
├── middleware/
│   ├── auth.go              # Token 认证中间件（从请求头读取 Token，查询用户信息）
│   └── init.go              # 初始化拦截中间件（无管理员时拦截所有页面到 /init）
├── handlers/
│   ├── index.go             # 首页 / 登录页
│   ├── init.go              # 初始化管理员页面（GET + POST）
│   ├── dashboard.go         # 控制台页面（Caddy 状态、下载、统计）
│   ├── site.go              # 网站 CRUD（列表 / 新增 / 编辑 / 删除）
│   ├── caddy.go             # Caddy 服务控制（启停 / 刷新 / 日志）
│   ├── log.go               # 日志查看 API
│   └── auth.go              # 登录/登出处理
├── services/
│   ├── caddy.go             # Caddy 二进制管理、服务控制、Caddyfile 生成
│   ├── site.go              # 网站配置业务逻辑
│   └── log.go               # 日志收集服务
├── templates/
│   ├── layout.html          # 全局布局（导航、侧栏、主内容区）
│   ├── init.html            # 初始化管理员页面
│   ├── login.html           # 登录页
│   ├── dashboard.html       # 控制台首页
│   ├── site_list.html       # 网站列表页
│   ├── site_form.html       # 新增/编辑网站表单页
│   ├── caddy_manage.html    # Caddy 服务管理页
│   ├── system_log.html      # 日志查看页
│   └── error.html           # 错误页（404/500）
├── static/
│   ├── css/
│   │   └── style.css        # 全局样式（天青色主题 + Build 风格）
│   └── js/
│       └── app.js           # 前端交互逻辑（含 Token 管理）
├── caddy_config/
│   ├── Caddyfile            # 主 Caddyfile（引用子文件）
│   └── sites/               # 子配置文件目录（按域名拆分）
│       └── .gitkeep
└── data/
    └── .gitkeep              # SQLite 数据库存放目录（挂载 volume）
```

### 技术选型理由

| 选择 | 理由 |
|------|------|
| **Gin** | 高性能、中间件生态成熟、模板渲染内建 |
| **SQLite + mattn/go-sqlite3** | 零配置、嵌入式、适合中小规模管理工具 |
| **html/template** | Gin 原生支持，模板继承简单 |
| **ini/config 库** | 配置文件解析，支持 YAML |
| **os/exec** | 管理 Caddy 进程（启停/刷新/hot-reload） |
| **golang.org/x/crypto/bcrypt** | 管理员密码哈希 |

---

## 四、数据库设计（遵循 Zolysoft 数据库规范）

> 由于是 SQLite，部分 PostgreSQL 特性（JSONB、触发器）需用替代方案或代码层实现。

### 表设计

#### t_user（管理员账号表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 雪花ID（用自增替代，单机 SQLite） |
| username | VARCHAR(64) NOT NULL | 管理员用户名 |
| password_hash | VARCHAR(255) NOT NULL | bcrypt 密码哈希 |
| created_by | BIGINT DEFAULT 0 | 创建人ID |
| created_time | BIGINT NOT NULL | 创建时间（毫秒时间戳） |
| update_by | BIGINT DEFAULT 0 | 最后修改人ID |
| update_time | BIGINT DEFAULT 0 | 最后修改时间 |
| delete_time | BIGINT DEFAULT 0 | 逻辑删除标记 |

#### t_site（代理网站表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| name | VARCHAR(128) NOT NULL | 网站名称（管理标识） |
| domain | VARCHAR(255) NOT NULL | 域名 |
| scheme | VARCHAR(8) DEFAULT 'https' | 请求类型：http / https |
| proxy_to | VARCHAR(512) NOT NULL | 代理目标地址 |
| compression | TINYINT DEFAULT 0 | 是否启用压缩：0=否, 1=是 |
| request_headers | TEXT | 请求头配置（JSON 格式，如 `[{"name":"X-Real-IP","value":"{remote_host}"}]`） |
| path_routes | TEXT | 分地址代理（JSON 格式，如 `[{"path":"/api","to":"http://backend:8080"}]`） |
| extra_config | TEXT | 额外 Caddyfile 片段 |
| status | VARCHAR(16) DEFAULT 'enabled' | enabled=启用, disabled=停用 |
| created_by | BIGINT DEFAULT 0 | 创建人ID |
| created_time | BIGINT NOT NULL | 创建时间 |
| update_by | BIGINT DEFAULT 0 | 最后修改人ID |
| update_time | BIGINT DEFAULT 0 | 最后修改时间 |
| delete_time | BIGINT DEFAULT 0 | 逻辑删除标记 |

#### t_config（系统配置表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| config_key | VARCHAR(128) NOT NULL UNIQUE | 配置键 |
| config_value | TEXT NOT NULL | 配置值 |
| description | VARCHAR(256) | 配置说明 |
| created_by | BIGINT DEFAULT 0 | 创建人ID |
| created_time | BIGINT NOT NULL | 创建时间 |
| update_by | BIGINT DEFAULT 0 | 最后修改人ID |
| update_time | BIGINT DEFAULT 0 | 最后修改时间 |
| delete_time | BIGINT DEFAULT 0 | 逻辑删除标记 |

预设配置项：
- `caddyfile_path` — 主 Caddyfile 路径
- `caddyfile_sites_dir` — 子配置文件目录
- `caddy_binary_path` — Caddy 二进制路径
- `server_domain` — 管理工具自身域名

#### t_caddy_log（Caddy 运行日志表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| log_level | VARCHAR(16) | 日志级别：info / warn / error |
| log_message | TEXT | 日志内容 |
| log_source | VARCHAR(64) | 来源（caddy stdout / stderr） |
| created_time | BIGINT NOT NULL | 记录时间（毫秒时间戳） |

#### t_token（登录令牌表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| user_id | BIGINT NOT NULL | 关联 t_user.id |
| token | VARCHAR(64) NOT NULL UNIQUE | 登录令牌（UUID 生成） |
| expires_at | BIGINT NOT NULL | 过期时间（毫秒时间戳，默认 24 小时有效） |
| created_by | BIGINT DEFAULT 0 | 创建人ID |
| created_time | BIGINT NOT NULL | 创建时间（毫秒时间戳） |
| update_by | BIGINT DEFAULT 0 | 最后修改人ID |
| update_time | BIGINT DEFAULT 0 | 最后修改时间 |
| delete_time | BIGINT DEFAULT 0 | 逻辑删除标记 |

#### t_login_log（登录记录表，Zolysoft 规范必备）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| user_id | BIGINT NOT NULL | 管理员ID |
| login_ip | VARCHAR(64) | 登录IP |
| login_time | BIGINT NOT NULL | 登录时间 |
| logout_time | BIGINT DEFAULT 0 | 登出时间 |
| user_agent | VARCHAR(512) | 浏览器UA |
| created_time | BIGINT NOT NULL | 记录创建时间 |

---

## 五、功能模块任务拆分

### 第一阶段：项目基础搭建（预计 2-3 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 1.1 | 初始化 Go 模块，添加 Gin、SQLite、bcrypt 等依赖 | `go.mod` + `go.sum` | P0 |
| 1.2 | 创建项目目录结构 | 完整目录树 | P0 |
| 1.3 | 实现数据库连接 + 自动迁移逻辑 | `database/database.go` | P0 |
| 1.4 | 实现数据模型定义（所有表结构） | `database/models.go` | P0 |
| 1.5 | 实现配置文件读取（YAML 格式） | `config/config.go` | P0 |
| 1.6 | 创建全局 CSS 样式（天青色 Build 风格主题） | `static/css/style.css` | P0 |

### 第二阶段：管理员初始化与认证（预计 1-2 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 2.1 | 实现初始化拦截中间件（无管理员时所有请求重定向到 /init，/init 和 /login 和静态资源除外） | `middleware/init.go` | P0 |
| 2.2 | 实现 Token 认证中间件（读 Header Token → 查 t_token + 联查 t_user → 注入 gin.Context；失败返回 401 JSON 或重定向到 /login?redirect=原路径） | `middleware/auth.go` | P0 |
| 2.3 | 实现初始化页面（首次设置管理员账号密码） | `handlers/init.go` + `templates/init.html` | P0 |
| 2.4 | 实现登录/登出功能（登录生成 UUID Token 存 t_token，登出删除 Token；前端管理 localStorage） | `handlers/auth.go` + `templates/login.html` | P0 |
| 2.5 | 实现全局布局模板（侧栏导航 + 主内容区） | `templates/layout.html` | P0 |

### 第三阶段：Caddy 服务管理（预计 2-3 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 3.1 | Caddy 二进制检测与一键下载 | `services/caddy.go` | P0 |
| 3.2 | Caddy 服务启/停控制（os/exec 进程管理） | `services/caddy.go` | P0 |
| 3.3 | Caddy 服务 reload（更新 Caddyfile 后热刷新） | `services/caddy.go` | P0 |
| 3.4 | Caddy 运行日志收集入库 | `services/log.go` | P0 |
| 3.5 | 控制台页面（Caddy 状态、统计、操作入口） | `handlers/dashboard.go` + `templates/dashboard.html` | P0 |
| 3.6 | 日志查看页面 | `handlers/log.go` + `templates/system_log.html` | P1 |

### 第四阶段：网站代理管理（预计 3-4 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 4.1 | 网站列表页面（分页、搜索、数量统计） | `handlers/site.go` + `templates/site_list.html` | P0 |
| 4.2 | 新增代理网站表单页面 | `handlers/site.go` + `templates/site_form.html` | P0 |
| 4.3 | 编辑代理网站表单页面 | （复用 site_form.html） | P0 |
| 4.4 | 删除代理网站（软删除 + 清理子配置文件） | `services/site.go` | P0 |
| 4.5 | Caddyfile 子文件生成引擎 | `services/caddy.go` | P0 |
| 4.6 | Caddyfile 配置管理页（主文件位置、子文件目录） | `handlers/dashboard.go` + `templates/caddy_manage.html` | P1 |

### 第五阶段：配置管理页面（预计 1-2 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 5.1 | 系统配置管理页面（域名、路径等系统级配置） | `handlers/dashboard.go` + `templates/caddy_manage.html` | P1 |
| 5.2 | 配置保存与 Caddyfile 重新生成联动 | `services/caddy.go` | P1 |

### 第六阶段：收尾与打磨（预计 1 小时）

| # | 任务 | 产出 | 优先级 |
|---|------|------|--------|
| 6.1 | 路由注册与 `main.go` 组装 | `main.go` | P0 |
| 6.2 | 404/500 错误页面 | `templates/error.html` | P1 |
| 6.3 | 前端 JS 交互（Token 管理 localStorage 存取、请求拦截 401 跳转登录+记录来源页、表单验证、Toast、加载遮罩、API 请求统一带 Token Header） | `static/js/app.js` | P1 |
| 6.4 | 全局状态覆盖（加载中/空数据/错误/正常） | 前端模板 | P1 |

---

## 六、关键设计决策

### 1. 使用 html/template 而非前端 SPA
理由：管理工具型项目，无需 SEO、无需复杂前端状态，服务端渲染 + 少量 JS 交互即可满足需求。保持技术栈一致性。

### 2. Token 认证机制

```
前端                               后端
 │                                  │
 │  POST /login (账号+密码)          │
 │ ─────────────────────────────────> │  验证密码 → 生成 UUID Token → 存入 t_token（24h 有效期）
 │  ←─ { token: "xxx", user: {...}}  │
 │                                  │
 │  存 token 到 localStorage          │
 │                                  │
 │  GET /sites (Header: Token: xxx) │
 │ ─────────────────────────────────> │  中间件查 t_token → 未过期 → 查 t_user → set gin.Context
 │  ←─ 200 正常返回                   │
 │                                  │
 │  或: Token 过期/不存在             │
 │  ←─ 401 { error: "TOKEN_EXPIRED" }│
 │                                  │
 │  localStorage.removeItem("token") │
 │  记录当前路径 → 跳转 /login?redirect=/sites
 │                                  │
 │  登录成功后跳回 redirect 所指页面    │
```

**登出**：POST /logout → 删除 t_token 记录 → 前端清除 localStorage。

**过期清理**：启动时/定时删除 `expires_at < now` 的过期 Token，防止表膨胀。

**Token 存放位置**：
- 前端：`localStorage`（键名 `auth_token`），请求时通过 JS 设置 Header `Token: xxx`
- 后端验证：从请求头 `Token` 字段读取，查 `t_token WHERE token = ? AND expires_at > ? AND delete_time = 0`，联表查用户信息注入 `gin.Context`。

### 3. Caddyfile 子文件生成策略
每个代理网站对应一个 `.conf` 子文件，放在 `caddy_config/sites/` 目录下。修改网站配置时重写对应子文件，然后执行 `caddy reload`（或 `kill -SIGUSR1`）刷新配置。

示例生成格式：
```caddyfile
# Auto-generated by Caddy Server Manager
# Site: example.com

example.com {
    reverse_proxy http://localhost:8080 {
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
    }
    encode gzip
}
```

### 4. 日志收集
Caddy 以子进程方式启动，通过捕获 `stdout`/`stderr` 管道实时写入 `t_caddy_log` 表。保留最近 N 条日志，超过阈值自动清理。

### 5. Caddy 下载策略
控制台检测 `server/caddy` 二进制是否存在，若无则提供一键下载按钮，从官方 GitHub Release 下载对应平台的二进制文件到 `server/` 目录。

### 6. 前端 Token 管理流程（app.js 核心逻辑）

```
// 登录成功：存 token + 跳转
localStorage.setItem('auth_token', data.token);
window.location.href = redirectUrl || '/';

// 所有 API 请求：统一带 Token Header
fetch(url, {
  headers: { 'Token': localStorage.getItem('auth_token') }
});

// 401 拦截：清除 token + 记录来源 + 跳转登录
if (response.status === 401) {
  localStorage.removeItem('auth_token');
  window.location.href = '/login?redirect=' + encodeURIComponent(window.location.pathname);
}

// 登出：清除 token + 跳转登录
localStorage.removeItem('auth_token');
window.location.href = '/login';
```

**后端中间件对 API 请求和页面请求的处理差异**：
- API 请求（`/api/*` 或 `Content-Type: application/json`）：Token 无效返回 `401 JSON { error: "TOKEN_EXPIRED" }`
- 页面请求（text/html）：Token 无效返回 `302 Redirect /login?redirect=原路径`

### 7. Docker 部署方案

#### 镜像基座
`alpine:3.21`，Caddy 通过 `apk add caddy` 预装，无需在管理页下载。

#### Caddy 运行方式：os/exec 子进程模式
Go 程序以 `exec.Command` 启动 Caddy 子进程，负责完整的生命周期管理。

| 操作 | 命令/方式 | 说明 |
|------|----------|------|
| 启动 | `caddy run --config <path>` | 前台运行，捕获 stdout/stderr |
| 停止 | `cmd.Process.Signal(syscall.SIGTERM)` | 优雅停止 |
| 重载配置 | `caddy reload --config <path>` | Caddy 原生热重载，不中断流量 |
| 日志收集 | `cmd.StdoutPipe()` + `cmd.StderrPipe()` | 实时扫描 + 批量写入 t_caddy_log |
| 进程退出 | Go 作为 PID 1 监听 SIGTERM -> 先停 Caddy -> 再退出 | 保证优雅关闭 |

#### Dockerfile

```dockerfile
FROM alpine:3.21

RUN apk add --no-cache caddy curl ca-certificates

RUN mkdir -p /app/server /app/caddy_config/sites /app/data

COPY caddy_server /app/caddy_server
COPY config.yaml /app/config/config.yaml
COPY templates/ /app/templates/
COPY static/ /app/static/

WORKDIR /app

ENV CONFIG_PATH=/app/config/config.yaml

EXPOSE 80 443

ENTRYPOINT ["./caddy_server"]
```

#### 配置文件 `config.yaml`

```yaml
server:
  port: 8080

database:
  path: /app/data/caddy.db

caddy:
  binary_path: /usr/sbin/caddy            # apk 安装路径
  caddyfile_path: /app/caddy_config/Caddyfile
  caddyfile_sites_dir: /app/caddy_config/sites
  server_domain: ""                        # 管理工具自身域名，空则不在 Caddyfile 中配置
```

#### Docker Compose 示例

```yaml
services:
  caddy-manager:
    build: .
    container_name: caddy-manager
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - ./config.yaml:/app/config/config.yaml
      - ./caddy_config:/app/caddy_config
      - ./data:/app/data
      - caddy_data:/root/.local/share/caddy   # Caddy 证书等数据持久化
    restart: unless-stopped

volumes:
  caddy_data:
```

#### 关键调整

| 调整点 | 说明 |
|--------|------|
| PID 1 信号处理 | Go 程序作为容器 PID 1 需显式监听 `os.Signal`，收到 SIGTERM 时先 `cmd.Process.Signal(syscall.SIGTERM)` 停 Caddy，再 `os.Exit(0)` |
| 下载功能 | Docker 中 Caddy 已预装，管理页下载按钮替换为提示文字 "Docker 环境，Caddy 已预装"，通过环境变量 `RUN_ENV=docker` 控制显示 |
| 配置文件路径 | 通过环境变量 `CONFIG_PATH` 注入，默认 `/app/config/config.yaml` |
| Caddy 证书目录 | 挂载 volume 到 `/root/.local/share/caddy` 持久化 Let's Encrypt 证书，避免容器重启导致证书丢失 |
| 端口分离 | 管理工具自身监听 `8080`，Caddy 监听 `80/443`，互不冲突 |

---

## 七、路由设计

```
GET  /init              → 初始化管理员页面（无管理员时强制重定向至此）
POST /init              → 提交管理员账号密码
GET  /login             → 登录页（接收 ?redirect= 参数，登录成功后跳转回原页面）
POST /login             → 登录提交 → 返回 Token JSON
GET  /logout            → 登出 → 删除 Token → 跳转 /login
GET  /                  → 控制台首页（需 Token 认证）
GET  /sites             → 网站列表
GET  /sites/new         → 新增网站
POST /sites/new         → 保存新增网站
GET  /sites/:id/edit    → 编辑网站
POST /sites/:id/edit    → 保存编辑
POST /sites/:id/delete  → 删除网站
GET  /caddy/manage      → Caddy 管理页（启停/刷新/配置）
POST /caddy/start       → 启动 Caddy
POST /caddy/stop        → 停止 Caddy
POST /caddy/reload      → 刷新 Caddy 配置
POST /caddy/download    → 下载 Caddy 二进制
GET  /logs              → 日志查看页
GET  /api/logs          → 日志数据 API（分页，前端 fetch 带 Token Header）
GET  /api/stats         → 统计数据 API
GET  /api/caddy/status  → Caddy 运行状态 API
GET  /static/*          → 静态文件服务
```

---

## 八、建议的执行顺序

按依赖关系，推荐按以下顺序编码：

1. **项目骨架**：`go.mod` 依赖 → 目录结构 → `config` → `database`（迁移 + 模型）
2. **主题样式**：`static/css/style.css`（先确定视觉基调，后续模板直接引用）
3. **初始化流程**：中间件 → init 页面 → login 页面 → layout
4. **Caddy 核心**：services 层 → dashboard → Caddy 管理页
5. **站点管理**：site CRUD + Caddyfile 生成
6. **日志系统**：日志收集 + 查看页
7. **收尾整合**：`main.go` 路由 → 错误页 → 前端 JS

---

> **下一步**：请确认以上分析和计划。确认后我将开始第一阶段编码。
