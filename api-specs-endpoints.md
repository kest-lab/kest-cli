# API Specifications 端点总结

共 **9 个** API 端点，所有端点均需要认证（JWT Token），且作用于特定项目下。

基础路径：`/v1/projects/:id/api-specs`

---

## 端点列表

| # | 方法 | 路径 | 描述 |
|---|------|------|------|
| 1 | GET | `/projects/:id/api-specs` | 获取 API 规范列表 |
| 2 | POST | `/projects/:id/api-specs` | 创建 API 规范 |
| 3 | GET | `/projects/:id/api-specs/:sid` | 获取单个 API 规范详情 |
| 4 | GET | `/projects/:id/api-specs/:sid/full` | 获取 API 规范详情（含示例） |
| 5 | PATCH | `/projects/:id/api-specs/:sid` | 更新 API 规范 |
| 6 | DELETE | `/projects/:id/api-specs/:sid` | 删除 API 规范 |
| 7 | POST | `/projects/:id/api-specs/import` | 导入 API 规范 |
| 8 | GET | `/projects/:id/api-specs/export` | 导出 API 规范 |
| 9 | POST | `/projects/:id/api-specs/:sid/examples` | 创建请求/响应示例 |

---

## 端点详细说明

### 1. 获取 API 规范列表

- **方法**：`GET /projects/:id/api-specs`
- **描述**：获取指定项目下的所有 API 规范，支持分页、搜索和筛选
- **路径参数**：`id` — 项目 ID
- **查询参数**：
  - `page` — 页码（默认 1）
  - `per_page` — 每页数量（默认 20，最大 100）
  - `search` — 按名称或版本搜索
  - `version` — 按版本筛选
  - `status` — 按状态筛选（draft / published / deprecated）
- **响应**：返回 `items` 数组和 `pagination` 分页信息

### 2. 创建 API 规范

- **方法**：`POST /projects/:id/api-specs`
- **描述**：在指定项目下创建一个新的 API 规范
- **路径参数**：`id` — 项目 ID
- **请求体**：
  - `name`（必填）— 规范名称，1-100 字符
  - `version`（必填）— 版本号，语义化版本格式（如 1.0.0）
  - `description` — 描述，最大 500 字符
  - `spec_type` — 规范类型（openapi / swagger / postman）
  - `content`（必填）— OpenAPI/Swagger 规范 JSON 内容
  - `status` — 状态（默认 draft）
- **响应**：返回创建的 API 规范对象（201 Created）

### 3. 获取单个 API 规范详情

- **方法**：`GET /projects/:id/api-specs/:sid`
- **描述**：获取指定 API 规范的详细信息
- **路径参数**：
  - `id` — 项目 ID
  - `sid` — 规范 ID
- **响应**：返回完整的 API 规范对象，包含 `examples_count` 字段

### 4. 获取 API 规范详情（含示例）

- **方法**：`GET /projects/:id/api-specs/:sid/full`
- **描述**：获取 API 规范的完整信息，包括所有请求/响应示例
- **路径参数**：
  - `id` — 项目 ID
  - `sid` — 规范 ID
- **响应**：返回 API 规范对象及其关联的 `examples` 数组

### 5. 更新 API 规范

- **方法**：`PATCH /projects/:id/api-specs/:sid`
- **描述**：更新已有的 API 规范，支持部分更新
- **路径参数**：
  - `id` — 项目 ID
  - `sid` — 规范 ID
- **请求体**（均为可选）：
  - `name` — 规范名称
  - `version` — 版本号
  - `description` — 描述
  - `content` — 规范内容
  - `status` — 状态（draft / published / deprecated）
- **响应**：返回更新后的 API 规范对象

### 6. 删除 API 规范

- **方法**：`DELETE /projects/:id/api-specs/:sid`
- **描述**：删除指定的 API 规范及其所有关联示例。**此操作不可逆**
- **路径参数**：
  - `id` — 项目 ID
  - `sid` — 规范 ID
- **响应**：返回成功消息

### 7. 导入 API 规范

- **方法**：`POST /projects/:id/api-specs/import`
- **描述**：从多种来源批量导入 API 规范
- **路径参数**：`id` — 项目 ID
- **请求体**：
  - `source`（必填）— 导入来源（url / file / postman / swagger）
  - `data`（必填）— 导入数据（URL 字符串或 JSON 对象）
  - `options` — 导入选项（如 `auto_publish`、`prefix`）
- **响应**：返回导入结果，包含 `imported`（导入数）、`updated`（更新数）、`errors`（错误列表）

### 8. 导出 API 规范

- **方法**：`GET /projects/:id/api-specs/export`
- **描述**：将 API 规范导出为多种格式
- **路径参数**：`id` — 项目 ID
- **查询参数**：
  - `format` — 导出格式（json / yaml / postman，默认 json）
  - `ids` — 指定导出的规范 ID（逗号分隔，不填则导出全部）
  - `include_examples` — 是否包含示例（默认 false）
- **响应**：返回导出的规范数据

### 9. 创建请求/响应示例

- **方法**：`POST /projects/:id/api-specs/:sid/examples`
- **描述**：为指定 API 规范创建一个请求/响应示例
- **路径参数**：
  - `id` — 项目 ID
  - `sid` — 规范 ID
- **请求体**：
  - `path`（必填）— API 路径（如 /users）
  - `method`（必填）— HTTP 方法（GET / POST 等）
  - `status_code`（必填）— 响应状态码
  - `request_headers` — 请求头
  - `request_body` — 请求体
  - `response_headers` — 响应头
  - `response_body` — 响应体
  - `description` — 示例描述
- **响应**：返回创建的示例对象（201 Created）
