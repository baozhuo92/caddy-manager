# Caddy Manager

Caddy Web 服务器管理面板 — 基于官方 REST API 开发的现代化管理界面。

接口文档：https://caddyserver.com/docs/api

## 功能

| 模块 | 说明 |
|------|------|
| **仪表盘** | 服务器状态概览、快捷操作入口 |
| **配置管理** | GET/POST/PUT/PATCH/DELETE `/config/[path]`，路径导航、JSON 编辑器 |
| **服务器控制** | 停止服务 (`POST /stop`)、加载配置 (`POST /load`)、配置适配 (`POST /adapt`) |
| **证书管理** | PKI CA 信息查看、证书链 (`GET /pki/ca/<id>`) |
| **反向代理** | 上游服务器 CRUD、实时状态监控 (`GET /reverse_proxy/upstreams`) |

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建**: Vite 8
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP**: Axios
- **设计系统**: UI/UX PRO MAX — 工具实用风格，科技蓝调色板
  - 主色: `#1E40AF` | 辅色: `#3B82F6` | 强调色: `#0EA5E9`
  - 字体: Space Grotesk (标题) / DM Sans (正文) / Fira Code (代码)

## 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器 (默认 http://localhost:5173)
npm run dev
```

开发服务器已配置代理，将 `/api/*` 请求转发到 Caddy API 地址 (`http://localhost:2019`)。

```bash
# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

## 配置 Caddy API 地址

启动后在页面顶部的输入框中修改 Caddy API 地址（默认 `http://localhost:2019`），或修改 `vite.config.js` 中的 `target` 字段。

## 项目结构

```
src/
├── assets/          # 全局样式
│   ├── base.css     # CSS 变量与基础重置
│   └── main.css     # 样式入口
├── components/
│   ├── common/      # 通用组件 (Btn, Card, Input, Table, Badge, Message, JsonEditor)
│   └── proxy/       # 反向代理相关组件
├── router/
│   └── index.js     # 路由配置
├── services/
│   └── caddyApi.js  # Caddy API 封装
├── stores/
│   ├── config.js    # 配置管理状态
│   ├── server.js    # 服务器控制状态
│   ├── certificates.js  # 证书管理状态
│   └── proxy.js     # 反向代理状态
└── views/
    ├── Dashboard.vue
    ├── Config.vue
    ├── Server.vue
    ├── Certificates.vue
    └── Proxy.vue
```

## 待开发内容

- **反向代理站点管理** — 完整的站点级 CRUD：新增、修改、删除反向代理站点配置（域名、上游服务器、TLS、健康检查等）
- **页面功能 BUG 修复** — 各模块交互细节优化与异常处理
- **配置导入/导出** — 可视化对比、历史版本管理
- **日志实时查看** — 对接 Caddy 日志端点
- **暗色模式** — 跟随系统 / 手动切换

## API 端点覆盖

- [x] `GET /config/[path]` — 导出配置
- [x] `POST /config/[path]` — 设置/追加配置
- [x] `PUT /config/[path]` — 创建/插入配置
- [x] `PATCH /config/[path]` — 替换配置
- [x] `DELETE /config/[path]` — 删除配置
- [x] `POST /load` — 加载配置
- [x] `POST /stop` — 停止服务
- [x] `POST /adapt` — 配置适配
- [x] `GET /pki/ca/<id>` — PKI CA 信息
- [x] `GET /pki/ca/<id>/certificates` — 证书链
- [x] `GET /reverse_proxy/upstreams` — 上游状态

## UI 设计规范

- **风格**: 工具实用风格 (Data-Dense Dashboard)
- **配色**: 科技蓝调色板，避免紫色
- **组件**: SVG 图标 (Lucide 风格)，统一 150-300ms 过渡动画
- **可访问性**: WCAG AA 标准，键盘导航支持


# 功能调整
- 根据配置管理获取到的json配置分析出json结构
- 使用分析出来的配置结构用于新增、修改、删除代理
- 当前页面上有很多功能是不必要的，比如说：配置管理的编辑和删除按钮以及最下面的快捷路径
- 这个管理工具最基础的用法就是通过它去配置反向代理（需要增加反向代理的优先级最高），其他的基本上只需要看看就行
- 反向代理需要可以配置跨域、basicauth 、websocket、自定义头、域名、选择类型http/https/all等