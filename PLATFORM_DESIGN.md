# Kest API Platform - 产品设计文档

基于 YApi 的理念，结合 Kest CLI 的能力，打造内网 API 管理平台

---

## 📊 YApi 核心功能分析

### YApi 的主要特点

1. **接口管理**
   - 接口文档编写和展示
   - 支持多种数据格式（JSON、Form、XML）
   - Mock 数据服务
   - 接口自动化测试

2. **项目管理**
   - 多项目管理
   - 成员权限控制
   - 分组管理

3. **协作能力**
   - 在线编辑器
   - 评论和变更记录
   - 接口变更通知

4. **自动化**
   - 导入 Swagger/OpenAPI
   - 导出接口文档
   - 接口自动测试

### YApi 的不足

- ❌ 部署复杂（需要 MongoDB）
- ❌ 没有 CLI 工具
- ❌ Mock Server 功能简单
- ❌ 缺少 gRPC 支持
- ❌ 没有历史回放功能
- ❌ 性能测试功能弱

---

## 🚀 Kest API Platform 设计方案

### 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    Kest API Platform                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │   Kest CLI   │◄───┤  Kest Server │───►│   Web UI     │ │
│  │              │    │              │    │              │ │
│  │ - Test       │    │ - API        │    │ - Dashboard  │ │
│  │ - Record     │    │ - Storage    │    │ - Editor     │ │
│  │ - Generate   │    │ - Mock       │    │ - Viewer     │ │
│  └──────────────┘    └──────────────┘    └──────────────┘ │
│         │                    │                    │        │
│         └────────────────────┴────────────────────┘        │
│                              │                             │
│                    ┌─────────▼─────────┐                   │
│                    │   SQLite / PG     │                   │
│                    │                   │                   │
│                    │ - APIs            │                   │
│                    │ - Projects        │                   │
│                    │ - History         │                   │
│                    │ - Mock Data       │                   │
│                    └───────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 💡 核心功能设计

### 1. CLI 自动生成文档

```bash
# 从实际测试生成文档
kest doc generate --from-history --project my-api

# 从 scenario 生成文档
kest doc generate --from-scenario tests.kest

# 推送到平台
kest doc push --project my-api --version v1.0.0

# 导出为 Markdown/HTML
kest doc export --format markdown -o API.md
```

**工作流程**：
```
CLI 测试 → 自动记录 → 分析请求/响应 → 生成文档 → 推送到平台
```

---

### 2. 项目结构

```yaml
# kest-platform.yaml
platform:
  url: http://api-platform.internal:3000
  token: ${KEST_PLATFORM_TOKEN}

projects:
  - id: user-service
    name: 用户服务
    version: v1.0.0
    base_url: https://api.example.com
    
  - id: order-service
    name: 订单服务
    version: v1.0.0
    base_url: https://api.example.com
```

---

### 3. API 文档自动化

#### 从历史记录生成文档

```go
// internal/doc/generator.go
package doc

type APIDoc struct {
    Project     string
    Version     string
    Endpoints   []Endpoint
    GeneratedAt time.Time
}

type Endpoint struct {
    Method      string
    Path        string
    Summary     string
    Description string
    Request     RequestSpec
    Response    ResponseSpec
    Examples    []Example
}

// 从历史记录生成
func GenerateFromHistory(projectID string, limit int) (*APIDoc, error) {
    store, _ := storage.NewStore()
    records := store.GetHistory(limit, projectID)
    
    // 分析每个请求
    endpoints := make(map[string]*Endpoint)
    for _, r := range records {
        key := r.Method + " " + r.Path
        if ep, exists := endpoints[key]; exists {
            // 添加示例
            ep.Examples = append(ep.Examples, toExample(r))
        } else {
            // 创建新端点
            endpoints[key] = &Endpoint{
                Method: r.Method,
                Path: r.Path,
                Request: analyzeRequest(r),
                Response: analyzeResponse(r),
                Examples: []Example{toExample(r)},
            }
        }
    }
    
    return &APIDoc{
        Project: projectID,
        Endpoints: mapToSlice(endpoints),
        GeneratedAt: time.Now(),
    }, nil
}
```

---

### 4. Web 平台功能

#### 核心页面

1. **项目列表**
   ```
   ┌────────────────────────────────────┐
   │  🏠 Kest API Platform              │
   ├────────────────────────────────────┤
   │                                    │
   │  📦 用户服务 (user-service)        │
   │     v1.0.0 | 23 APIs | Updated 2h  │
   │                                    │
   │  🛒 订单服务 (order-service)       │
   │     v1.0.0 | 15 APIs | Updated 5h  │
   │                                    │
   │  💳 支付服务 (payment-service)     │
   │     v1.0.0 | 8 APIs | Updated 1d   │
   └────────────────────────────────────┘
   ```

2. **API 详情页**
   ```
   ┌─────────────────────────────────────────────┐
   │  POST /api/users                            │
   ├─────────────────────────────────────────────┤
   │                                             │
   │  📝 Description: 创建新用户                 │
   │                                             │
   │  📤 Request:                                │
   │    Content-Type: application/json          │
   │    {                                        │
   │      "email": "string",                     │
   │      "name": "string"                       │
   │    }                                        │
   │                                             │
   │  📥 Response: 201 Created                   │
   │    {                                        │
   │      "id": 123,                             │
   │      "email": "test@example.com",           │
   │      "created_at": "2026-01-30T..."         │
   │    }                                        │
   │                                             │
   │  🧪 Test Examples (3):                      │
   │    - 成功创建用户 (200ms)                   │
   │    - 邮箱已存在 (45ms)                      │
   │    - 参数验证失败 (12ms)                    │
   │                                             │
   │  [Try it] [Copy cURL] [Generate Test]      │
   └─────────────────────────────────────────────┘
   ```

3. **历史记录和统计**
   ```
   ┌──────────────────────────────────────┐
   │  📊 API Statistics                   │
   ├──────────────────────────────────────┤
   │                                      │
   │  Total Requests: 1,234               │
   │  Success Rate: 98.5%                 │
   │  Avg Response Time: 234ms            │
   │                                      │
   │  Most Used APIs:                     │
   │  1. GET /users (456 calls)           │
   │  2. POST /login (234 calls)          │
   │  3. GET /orders (189 calls)          │
   │                                      │
   │  Performance Trends: [Chart]         │
   └──────────────────────────────────────┘
   ```

---

### 5. Mock Server 功能

```bash
# 启动 Mock Server（基于历史数据）
kest mock start --project user-service --port 8080

# 配置 Mock 规则
kest mock add --path "/users/:id" --response user.json --delay 100ms

# 智能 Mock（基于真实响应）
kest mock smart --from-history
```

**Mock Server 特性**：
- ✅ 基于真实历史响应
- ✅ 支持动态数据
- ✅ 可配置延迟和错误
- ✅ 支持 gRPC Mock

---

### 6. CLI 命令设计

```bash
# === 文档管理 ===
kest doc generate           # 从历史生成文档
kest doc push              # 推送到平台
kest doc pull              # 从平台拉取
kest doc export            # 导出文档

# === 平台集成 ===
kest platform login        # 登录平台
kest platform status       # 查看同步状态
kest platform sync         # 同步数据

# === Mock Server ===
kest mock start            # 启动 Mock Server
kest mock add              # 添加 Mock 规则
kest mock list             # 列出 Mock 规则

# === 项目管理 ===
kest project list          # 列出项目
kest project create        # 创建项目
kest project switch        # 切换项目
```

---

## 🏗️ 技术栈建议

### 后端（Kest Server）

```go
// 技术选型
- 框架: Gin (轻量、快速)
- 数据库: SQLite (单机) / PostgreSQL (生产)
- 认证: JWT
- API: RESTful + gRPC
```

**核心模块**：
```
kest-server/
├── internal/
│   ├── api/          # REST API handlers
│   ├── grpc/         # gRPC services
│   ├── storage/      # 数据库层
│   ├── doc/          # 文档生成
│   ├── mock/         # Mock Server
│   └── auth/         # 认证授权
├── web/              # 前端构建输出
└── cmd/
    └── server/       # 启动入口
```

---

### 前端（Web UI）

```typescript
// 技术选型
- 框架: Next.js 14 (React)
- UI: shadcn/ui + Tailwind CSS
- 状态: Zustand
- API: TanStack Query
- 编辑器: Monaco Editor
```

**页面结构**：
```
web/
├── app/
│   ├── projects/           # 项目列表
│   ├── projects/[id]/      # 项目详情
│   │   ├── apis/           # API 列表
│   │   ├── history/        # 历史记录
│   │   ├── mock/           # Mock 管理
│   │   └── settings/       # 设置
│   └── docs/               # 文档中心
└── components/
    ├── api-editor/         # API 编辑器
    ├── request-viewer/     # 请求查看器
    ├── mock-config/        # Mock 配置
    └── charts/             # 统计图表
```

---

## 📋 实施路线图

### Phase 1: CLI 增强（2周）

- [ ] 实现 `kest doc generate` 命令
- [ ] 从历史记录分析生成 API 文档
- [ ] 支持导出 Markdown/OpenAPI 格式
- [ ] 添加 `kest platform` 命令框架

### Phase 2: 基础平台（4周）

- [ ] 搭建 Kest Server 基础框架
- [ ] 实现项目和 API 管理 API
- [ ] 实现用户认证和权限
- [ ] 数据库 Schema 设计
- [ ] CLI 与 Server 集成

### Phase 3: Web UI（4周）

- [ ] 项目列表和详情页面
- [ ] API 文档展示页面
- [ ] 历史记录查看
- [ ] API 测试界面（类似 Postman）
- [ ] 统计和图表

### Phase 4: Mock Server（2周）

- [ ] Mock Server 核心引擎
- [ ] 智能 Mock（基于历史）
- [ ] Mock 管理界面
- [ ] gRPC Mock 支持

### Phase 5: 高级功能（4周）

- [ ] 团队协作功能
- [ ] Webhook 通知
- [ ] API 变更检测
- [ ] 性能趋势分析
- [ ] CI/CD 集成

---

## 🎯 与 YApi 的对比

| 功能 | YApi | Kest Platform |
|------|------|---------------|
| 部署复杂度 | 高（MongoDB） | 低（SQLite/单二进制） |
| CLI 工具 | ❌ | ✅ 强大的 CLI |
| gRPC 支持 | ❌ | ✅ 完整支持 |
| 自动文档生成 | 部分 | ✅ 从测试自动生成 |
| Mock Server | 基础 | ✅ 智能 Mock |
| 历史回放 | ❌ | ✅ 完整历史 |
| 性能测试 | ❌ | ✅ 内置支持 |
| 并行测试 | ❌ | ✅ 支持 |
| Streaming | ❌ | ✅ 支持 |
| 私有化部署 | ✅ | ✅ 更简单 |

---

## 💻 关键代码示例

### 1. 文档生成 API

```go
// internal/api/doc.go
func (h *Handler) GenerateDoc(c *gin.Context) {
    req := struct {
        ProjectID string `json:"project_id"`
        Limit     int    `json:"limit"`
    }{}
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 从历史生成文档
    doc, err := h.docService.GenerateFromHistory(req.ProjectID, req.Limit)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, doc)
}
```

### 2. CLI 推送命令

```go
// internal/cli/doc.go
var docPushCmd = &cobra.Command{
    Use:   "push",
    Short: "Push API documentation to platform",
    RunE: func(cmd *cobra.Command, args []string) error {
        // 生成文档
        doc, err := doc.GenerateFromHistory(projectID, 100)
        if err != nil {
            return err
        }
        
        // 推送到平台
        client := platform.NewClient(platformURL, token)
        return client.Push(doc)
    },
}
```

### 3. Mock Server

```go
// internal/mock/server.go
type MockServer struct {
    store   *storage.Store
    rules   map[string]*Rule
}

func (m *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    key := r.Method + " " + r.URL.Path
    
    // 检查自定义规则
    if rule, ok := m.rules[key]; ok {
        time.Sleep(rule.Delay)
        w.WriteHeader(rule.Status)
        w.Write([]byte(rule.Response))
        return
    }
    
    // 从历史中查找
    record := m.store.FindLatestMatch(r.Method, r.URL.Path)
    if record != nil {
        w.WriteHeader(record.ResponseStatus)
        w.Write([]byte(record.ResponseBody))
        return
    }
    
    w.WriteHeader(404)
}
```

---

## 🚀 快速开始（MVP）

### 最小可行产品包含：

1. **CLI 工具**
   - ✅ 已有测试功能
   - ✅ 添加文档生成
   - ✅ 添加平台推送

2. **Server**
   - ✅ 项目管理 API
   - ✅ 文档存储和展示
   - ✅ 基础认证

3. **Web UI**
   - ✅ 项目列表
   - ✅ API 文档展示
   - ✅ 简单的在线测试

### 预期时间：**6-8周**

---

## 📚 相关文档

创建以下文档：
1. `PLATFORM_DESIGN.md` - 详细设计文档
2. `API_SPEC.md` - Platform API 规范
3. `DEPLOYMENT.md` - 部署指南

---

**优势总结**：

✅ **基于 Kest CLI**：已有强大的测试功能  
✅ **自动化文档**：从真实测试自动生成，始终准确  
✅ **轻量部署**：单二进制 + SQLite，5分钟部署  
✅ **完整功能**：REST + gRPC + Streaming 全支持  
✅ **开发者友好**：CLI-first，AI 友好  
✅ **企业级**：权限控制、团队协作、私有化部署  

这将是一个比 YApi 更现代、更强大的 API 管理平台！🚀
