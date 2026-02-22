# Kest Web 项目上手指南（面向 Next.js 新手）

你现在看到的是一个 Monorepo，包含 CLI、后端 API、前端 Web 三部分。  
先说结论：`web` 目前不是 Next.js，而是 **React + Vite + TypeScript**。

## 1. 仓库结构

```text
kest/
├── api/                  # Go 后端（Gin + PostgreSQL）
├── web/                  # 前端控制台（React + Vite）
├── cli/                  # 命令行工具
├── docs/                 # 文档
├── docker-compose.yml    # 一键起后端+数据库
└── WEB_ONBOARDING_GUIDE_ZH.md
```

## 2. web 目录怎么读

```text
web/src/
├── main.tsx              # 前端入口：挂载 Router + React Query
├── App.tsx               # 路由总表
├── pages/                # 页面级代码（路由页面）
├── components/           # 组件
│   ├── ui/               # 基础 UI 组件（Button/Input/Dialog...）
│   ├── features/         # 业务组件
│   ├── layout/           # 布局组件（如 AdminLayout）
│   └── common/           # 通用组件
├── hooks/                # 业务 hooks（含 React Query hooks）
├── services/             # API 服务层（按业务模块封装接口）
├── http/                 # axios 封装、拦截器、错误处理
├── store/                # Zustand 全局状态（登录态等）
├── config/               # 环境变量、QueryClient 配置
├── themes/               # 主题 token（亮色/暗色/原子变量）
└── types/                # TypeScript 类型定义
```

## 3. 你关心的三件事

### 3.1 什么管 CSS？

- 全局样式入口：`web/src/index.css`
- 主题体系：`web/src/themes/index.css` + `web/src/themes/light.css` + `web/src/themes/dark.css` + `web/src/themes/primitives.css`
- 组件样式：多数直接写在 TSX 里，使用 Tailwind utility class

### 3.2 什么管组件？

- 基础组件库：`web/src/components/ui/*`
- 页面拼装和业务组件：`web/src/components/features/*`、`web/src/pages/*`
- 布局框架：`web/src/components/layout/admin-layout.tsx`

### 3.3 什么管前后端连接？

调用链路是：

1. 页面或组件调用 hooks：`web/src/hooks/use-kest-api.ts`
2. hooks 调用服务层：`web/src/services/kest-api.service.ts`
3. 服务层调用统一请求实例：`web/src/http/request.ts`
4. 请求实例从 `web/src/config/env.ts` 读取 `VITE_API_URL`

其中 `request.ts` 里已经做了 token 注入和统一错误处理。

## 4. 路由和页面入口

- 应用入口：`web/src/main.tsx`
- 路由定义：`web/src/App.tsx`
- 常见页面：
  - 项目列表：`/projects` -> `web/src/pages/projects/index.tsx`
  - 项目详情：`/projects/:id` -> `web/src/pages/projects/detail.tsx`
  - 登录：`/login` -> `web/src/pages/auth/login.tsx`

## 5. 本地启动（建议顺序）

### 5.1 启后端 + 数据库

在仓库根目录执行：

```bash
docker-compose up -d
```

后端健康检查：

```bash
curl http://localhost:8025/v1/health
```

### 5.2 启前端 web

```bash
cd web
npm install
npm run dev
```

默认前端端口是 `3000`。

## 6. 前后端联调注意点

- 当前 `web/src/config/env.ts` 默认 `VITE_API_URL=/api`
- `web/vite.config.ts` 里把 `/api` 代理到 `https://api.kest.dev`

如果你要联调本地后端，建议在 `web/.env.local` 写：

```env
VITE_API_URL=http://localhost:8025
```

这样前端请求会直接打到本地 API。

## 7. 新手推荐学习路径

1. 先从路由看页面：`web/src/App.tsx`
2. 选一个页面纵向读到底：
   - 页面：`web/src/pages/projects/detail.tsx`
   - hooks：`web/src/hooks/use-kest-api.ts`
   - service：`web/src/services/kest-api.service.ts`
   - request：`web/src/http/request.ts`
3. 再看状态管理：`web/src/store/auth-store.ts`
4. 最后再改 UI：`web/src/components/ui/*` 和 `web/src/themes/*`

## 8. 如何逐步完成一个需求

1. 新增一个页面：在 `web/src/pages` 建页面，再在 `web/src/App.tsx` 挂路由。
2. 新增一个接口：先在 `web/src/services/kest-api.service.ts` 加方法，再在 `web/src/hooks/use-kest-api.ts` 包装 query/mutation。
3. 新增一个组件：优先放 `web/src/components/features`；可复用基础组件放 `web/src/components/ui`。
4. 调样式：优先改 token（`themes`），再改单页面 class，避免到处硬编码颜色。

## 9. 常用命令

```bash
# 前端
cd web
npm run dev
npm run build
npm run lint

# 后端（在 api/）
go run cmd/server/main.go migrate
go run cmd/server/main.go
```

## 10. 相关文档

- 后端说明：`api/README.md`
- 前端说明（历史版本）：`web/README.md`
- API 文档目录：`docs/api/`
