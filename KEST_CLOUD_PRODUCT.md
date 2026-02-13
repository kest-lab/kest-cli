# Kest Cloud 产品定义

> **kest.run** — 云端 API 开发与测试平台

---

## 一句话介绍

**Kest Cloud 是一个云端 API 全生命周期管理平台，集 API 文档、自动化测试、流程编排、错误监控于一体。**

---

## 🎯 Kest Cloud 里面有什么？

### 核心功能模块（6 大模块）

```
┌─────────────────────────────────────────────────────────┐
│                    Kest Cloud (kest.run)                 │
├─────────┬──────────┬──────────┬──────────┬──────────────┤
│  📄     │  🧪      │  🔄      │  🐛      │  👥          │
│  API    │  Test    │  Flow    │  Issue   │  Team        │
│  Docs   │  Runner  │  Editor  │  Tracker │  Collab      │
└─────────┴──────────┴──────────┴──────────┴──────────────┘
```

---

### 1. 📄 API 文档中心（API Docs）

**一句话**：在线管理和浏览 API 文档，支持从 CLI 自动同步。

**功能**：
- 手动创建 API 规范（Method + Path + 描述）
- 从 CLI 自动导入（`kest doc push`）
- 从 OpenAPI/Swagger 批量导入
- 导出为标准格式
- 按分类（Category）组织 API
- 请求/响应示例管理
- 在线 API 调试（Try it）

**后端支持**：
```
POST   /v1/projects/:id/api-specs           创建 API 规范
GET    /v1/projects/:id/api-specs           列表
GET    /v1/projects/:id/api-specs/:sid      详情
GET    /v1/projects/:id/api-specs/:sid/full 详情 + 示例
PATCH  /v1/projects/:id/api-specs/:sid      更新
DELETE /v1/projects/:id/api-specs/:sid      删除
POST   /v1/projects/:id/api-specs/import    批量导入
GET    /v1/projects/:id/api-specs/export    导出
POST   /v1/projects/:id/api-specs/:sid/examples  添加示例
```

**分类管理**：
```
GET    /v1/projects/:id/categories          分类列表
POST   /v1/projects/:id/categories          创建分类
PATCH  /v1/projects/:id/categories/:cid     更新
DELETE /v1/projects/:id/categories/:cid     删除
PUT    /v1/projects/:id/categories/sort     排序
```

---

### 2. 🧪 云端测试运行（Test Runner）

**一句话**：在云端运行 API 自动化测试，查看结果和历史。

**功能**：
- 在线创建测试用例
- 云端运行测试（不需要本地环境）
- 测试结果实时展示
- 测试历史记录
- 多环境切换（dev / staging / prod）
- 环境变量管理
- 从 API 规范一键生成测试用例

**后端支持**：
```
# 测试用例
GET    /v1/projects/:id/test-cases              列表
POST   /v1/projects/:id/test-cases              创建
GET    /v1/projects/:id/test-cases/:tcid        详情
PATCH  /v1/projects/:id/test-cases/:tcid        更新
DELETE /v1/projects/:id/test-cases/:tcid        删除
POST   /v1/projects/:id/test-cases/:tcid/run    ⚡ 云端运行
POST   /v1/projects/:id/test-cases/from-spec    从 API 规范生成

# 环境管理
GET    /v1/projects/:id/environments            列表
POST   /v1/projects/:id/environments            创建
PATCH  /v1/projects/:id/environments/:eid       更新
DELETE /v1/projects/:id/environments/:eid       删除
POST   /v1/projects/:id/environments/:eid/duplicate  复制
```

---

### 3. 🔄 测试流程编排（Flow Editor）

**一句话**：可视化编排多步骤 API 测试流程，支持条件分支和数据传递。

**功能**：
- 可视化流程编辑器（拖拽式）
- 多步骤 API 调用编排
- 步骤间数据传递
- 条件分支和循环
- 云端执行流程
- 执行历史和日志
- SSE 实时事件推送（执行进度）

**后端支持**：
```
# 流程管理
GET    /v1/projects/:id/flows                   列表
POST   /v1/projects/:id/flows                   创建
GET    /v1/projects/:id/flows/:fid              详情
PATCH  /v1/projects/:id/flows/:fid              更新
PUT    /v1/projects/:id/flows/:fid              保存完整流程
DELETE /v1/projects/:id/flows/:fid              删除

# 步骤管理
POST   /v1/projects/:id/flows/:fid/steps        添加步骤
PATCH  /v1/projects/:id/flows/:fid/steps/:sid    更新步骤
DELETE /v1/projects/:id/flows/:fid/steps/:sid    删除步骤

# 连线管理
POST   /v1/projects/:id/flows/:fid/edges        添加连线
DELETE /v1/projects/:id/flows/:fid/edges/:eid    删除连线

# 执行
POST   /v1/projects/:id/flows/:fid/run          ⚡ 云端运行
GET    /v1/projects/:id/flows/:fid/runs          执行历史
GET    /v1/projects/:id/flows/:fid/runs/:rid     执行详情
GET    /v1/projects/:id/flows/:fid/runs/:rid/events  📡 SSE 实时事件
```

---

### 4. 🐛 错误监控（Issue Tracker）

**一句话**：集成 Sentry SDK，自动捕获和追踪 API 错误。

**功能**：
- 自动错误捕获（兼容 Sentry SDK）
- 错误列表和详情
- 错误状态管理（解决 / 忽略 / 重新打开）
- 错误事件时间线
- 按项目分组

**后端支持**：
```
# 数据采集（Sentry SDK 兼容）
POST   /api/:project_id/envelope/     接收错误数据
POST   /api/:project_id/store/        接收错误事件（旧版）

# 问题管理
GET    /v1/projects/:id/issues                        问题列表
GET    /v1/projects/:id/issues/:fingerprint            问题详情
POST   /v1/projects/:id/issues/:fingerprint/resolve    标记解决
POST   /v1/projects/:id/issues/:fingerprint/ignore     忽略
POST   /v1/projects/:id/issues/:fingerprint/reopen     重新打开
GET    /v1/projects/:id/issues/:fingerprint/events     事件列表
```

---

### 5. 👥 团队协作（Team Collaboration）

**一句话**：多人协作管理 API 项目，基于角色的权限控制。

**功能**：
- 项目成员管理（邀请 / 移除）
- 角色权限控制（Admin / Write / Read）
- 项目级别的访问控制
- 用户管理

**后端支持**：
```
# 成员管理
GET    /v1/projects/:id/members          成员列表
POST   /v1/projects/:id/members          添加成员
PATCH  /v1/projects/:id/members/:uid     更新角色
DELETE /v1/projects/:id/members/:uid     移除成员

# 角色权限
GET    /v1/roles                         角色列表
POST   /v1/roles                         创建角色
POST   /v1/roles/assign                  分配角色
POST   /v1/roles/remove                  移除角色
GET    /v1/users/:id/roles               用户角色
GET    /v1/permissions                   权限列表
```

---

### 6. 📊 项目管理（Project Hub）

**一句话**：创建和管理 API 项目，统一入口。

**功能**：
- 创建 / 编辑 / 删除项目
- 项目 DSN（用于 SDK 集成）
- 项目概览仪表板
- 系统健康检查

**后端支持**：
```
# 项目
POST   /v1/projects                      创建
GET    /v1/projects                      列表
GET    /v1/projects/:id                  详情
PUT    /v1/projects/:id                  更新
DELETE /v1/projects/:id                  删除
GET    /v1/projects/:id/dsn              获取 DSN

# 系统
GET    /v1/health                        健康检查
GET    /v1/system-features               系统功能
GET    /v1/setup-status                  初始化状态
GET    /swagger/*                        API 文档
GET    /metrics                          监控指标
```

---

## 🏗️ 产品架构图

```
┌──────────────────────────────────────────────────────────────┐
│                     Kest Cloud (kest.run)                     │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ 📄 API   │  │ 🧪 Test  │  │ 🔄 Flow  │  │ 🐛 Issue │    │
│  │   Docs   │  │  Runner  │  │  Editor  │  │ Tracker  │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
│       │              │              │              │          │
│  ┌────┴──────────────┴──────────────┴──────────────┴─────┐   │
│  │              👥 Team & Project Management              │   │
│  └───────────────────────┬───────────────────────────────┘   │
│                          │                                    │
│  ┌───────────────────────┴───────────────────────────────┐   │
│  │                  🔐 Auth & Permissions                 │   │
│  └───────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
                           │
                    ┌──────┴──────┐
                    │  Kest CLI   │
                    │  (本地工具)  │
                    └─────────────┘
```

---

## 🖥️ 页面结构设计

### 导航菜单

```
Kest Cloud
├── 🏠 Dashboard          项目概览、最近活动
├── 📄 API Docs           API 文档中心
│   ├── 分类浏览
│   ├── 搜索
│   └── Try it（在线调试）
├── 🧪 Test Runner        自动化测试
│   ├── 测试用例列表
│   ├── 运行测试
│   └── 测试报告
├── 🔄 Flows              测试流程编排
│   ├── 流程列表
│   ├── 可视化编辑器
│   └── 执行历史
├── 🐛 Issues             错误监控
│   ├── 问题列表
│   ├── 问题详情
│   └── 事件时间线
├── ⚙️ Settings           项目设置
│   ├── 基本信息
│   ├── 环境管理
│   ├── 成员管理
│   └── DSN / 集成
└── 👤 Account            账户
    ├── 个人资料
    └── 登出
```

---

## 📝 对外介绍文案

### 简短版（一段话）

> **Kest Cloud** 是一个云端 API 全生命周期管理平台。在这里，你可以管理 API 文档、运行自动化测试、编排测试流程、监控 API 错误，并与团队实时协作。配合 Kest CLI 使用，实现本地开发到云端管理的无缝衔接。

### 功能介绍版

> **Kest Cloud** (kest.run) 为你的 API 提供一站式云端服务：
>
> - **📄 API 文档** — 在线管理 API 文档，支持从 CLI 自动同步，告别过时文档
> - **🧪 云端测试** — 无需本地环境，一键在云端运行 API 自动化测试
> - **🔄 流程编排** — 可视化编排多步骤测试流程，支持条件分支和数据传递
> - **🐛 错误监控** — 兼容 Sentry SDK，自动捕获和追踪 API 异常
> - **👥 团队协作** — 基于角色的权限控制，多人协作管理 API 项目
>
> 配合 **Kest CLI** 使用，从本地开发到云端管理，一气呵成。

### Landing Page 版

```
Kest Cloud — 你的 API，尽在掌控

📄 文档即代码
   API 文档自动同步，永远保持最新。
   支持 OpenAPI 导入，在线浏览和调试。

🧪 云端一键测试
   不需要本地环境，在云端直接运行 API 测试。
   测试结果实时展示，历史记录随时查看。

🔄 可视化流程
   拖拽式编排多步骤测试流程。
   条件分支、数据传递、循环执行，全部可视化。

🐛 错误不遗漏
   兼容 Sentry SDK，自动捕获 API 异常。
   问题追踪、状态管理、事件时间线，一目了然。

👥 团队无缝协作
   邀请团队成员，按角色分配权限。
   项目级别管理，安全可控。

⚡ CLI + Cloud，双剑合璧
   本地用 CLI 开发测试，云端用 Cloud 管理协作。
   数据自动同步，工作流无缝衔接。
```

---

## 🔗 CLI 与 Cloud 的关系

```
┌─────────────────────────────────────────────────────┐
│                   开发者工作流                         │
└─────────────────────────────────────────────────────┘

  本地开发                          云端管理
  ────────                          ────────
  
  $ kest get /api/users    ──────>  📄 自动记录为文档
  $ kest post /api/users   ──────>  📄 自动生成示例
  $ kest test run          ──────>  🧪 测试结果同步到云端
  $ kest doc push          ──────>  📄 文档推送到 Cloud
  $ kest flow run          ──────>  🔄 流程执行记录同步
  
                                    🐛 SDK 自动上报错误
                                    👥 团队成员在线查看
                                    📊 Dashboard 汇总展示
```

### 具体场景

**场景 1：API 文档管理**
```bash
# 本地：CLI 测试 API
$ kest post /api/users -d '{"email":"test@test.com"}'
# → 自动记录请求和响应

# 本地：推送文档到云端
$ kest doc push --project my-api

# 云端：团队成员在 kest.run 浏览文档
# 云端：在线 Try it 调试 API
```

**场景 2：自动化测试**
```bash
# 本地：CLI 运行测试
$ kest test run --suite smoke-test
# → 结果自动同步到云端

# 云端：在 kest.run 查看测试报告
# 云端：直接在云端运行测试（CI/CD 集成）
# 云端：查看测试历史和趋势
```

**场景 3：流程编排**
```bash
# 云端：在 kest.run 可视化编排测试流程
#   Step 1: POST /api/users → 创建用户
#   Step 2: POST /api/auth/login → 登录获取 token
#   Step 3: GET /api/users/me → 验证用户信息

# 云端：一键运行流程
# 云端：查看每步执行结果和耗时

# 本地：也可以用 CLI 运行
$ kest flow run --id 123
```

**场景 4：错误监控**
```javascript
// 在你的应用中集成 Sentry SDK
import * as Sentry from '@sentry/node';
Sentry.init({ dsn: 'https://key@kest.run/api/1/envelope/' });

// API 错误自动上报到 Kest Cloud
// 在 kest.run 查看错误详情、管理状态
```

---

## 📊 与竞品对比

| 功能 | Kest Cloud | Postman | Swagger UI | Sentry |
|------|-----------|---------|------------|--------|
| API 文档管理 | ✅ | ✅ | ✅ | ❌ |
| 在线 API 调试 | ✅ | ✅ | ✅ | ❌ |
| 自动化测试 | ✅ | ✅ | ❌ | ❌ |
| 测试流程编排 | ✅ | ⚠️ 有限 | ❌ | ❌ |
| 错误监控 | ✅ | ❌ | ❌ | ✅ |
| CLI 工具 | ✅ | ⚠️ Newman | ❌ | ✅ |
| 团队协作 | ✅ | ✅ | ❌ | ✅ |
| 自托管 | ✅ | ❌ | ✅ | ✅ |
| 开源 | ✅ | ❌ | ✅ | ✅ |
| **一体化** | **✅** | ❌ | ❌ | ❌ |

**Kest Cloud 的独特优势**：把 Postman + Swagger + Sentry 的核心功能整合到一个平台。

---

## 🎯 目标用户

### 主要用户
- **后端开发者** — 管理 API 文档、运行测试
- **前端开发者** — 查看 API 文档、使用 Mock
- **QA 工程师** — 编排测试流程、运行自动化测试
- **技术负责人** — 监控 API 质量、管理团队

### 使用场景
- **个人开发者** — 免费使用，管理个人项目
- **小团队** — 协作管理 API，共享文档和测试
- **企业** — 自托管部署，完整的权限控制

---

## 💰 定价建议（未来）

| 计划 | 价格 | 项目数 | 成员数 | 功能 |
|------|------|--------|--------|------|
| **Free** | $0 | 3 | 1 | 基础功能 |
| **Pro** | $9/月 | 无限 | 5 | 全部功能 |
| **Team** | $29/月 | 无限 | 20 | 全部功能 + 优先支持 |
| **Enterprise** | 联系我们 | 无限 | 无限 | 自托管 + SLA |

---

## 🚀 上线计划

### Phase 1：MVP（当前）
- ✅ 用户认证（登录/注册）
- ✅ 项目管理
- ✅ API 文档管理（手动创建）
- ⏳ API 文档浏览和搜索

### Phase 2：核心功能（2-4 周）
- ⏳ 云端测试运行
- ⏳ 环境管理 UI
- ⏳ 测试结果展示
- ⏳ API 在线调试（Try it）

### Phase 3：高级功能（4-8 周）
- ⏳ Flow 可视化编辑器
- ⏳ Issue 错误监控面板
- ⏳ 团队成员管理 UI
- ⏳ CLI 集成

### Phase 4：完善（8-12 周）
- ⏳ Dashboard 仪表板
- ⏳ 通知系统
- ⏳ 导入/导出
- ⏳ 自托管部署文档

---

**产品名称**: Kest Cloud  
**域名**: kest.run  
**定位**: 云端 API 全生命周期管理平台  
**口号**: Your API, Under Control.
