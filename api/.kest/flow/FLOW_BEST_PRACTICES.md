# 📚 Kest Flow 最佳实践指南

基于实际测试经验总结的 Flow 文件编写最佳实践。

---

## 🎯 核心原则

1. **每个步骤应该独立可验证**
2. **使用变量传递数据，避免硬编码**
3. **断言应该灵活，不要过于严格**
4. **包含清理步骤，避免测试数据污染**

---

## 📝 Flow 文件结构

### 标准模板

```markdown
# 模块名称 Flow

模块功能描述

---

## Step 1: 前置条件（如登录）

```kest
POST /api/login
Content-Type: application/json

{
  "username": "{{test_user}}",
  "password": "{{test_password}}"
}

[Captures]
access_token: data.token

[Asserts]
status >= 200
status < 300
body.data.token exists
duration < 1000ms
```

---

## Step 2: 主要测试逻辑

```kest
POST /api/resource
Authorization: Bearer {{access_token}}

{
  "name": "Test {{$timestamp}}"
}

[Captures]
resource_id: data.id

[Asserts]
status >= 200
status < 300
body.data.id exists
```

---

## Step N: 清理步骤

```kest
DELETE /api/resource/{{resource_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
status < 300
```
```

---

## ✅ 变量使用最佳实践

### 1. 使用 Config 变量存储常量

**好的做法** ✅:
```yaml
# .kest/config.yaml
environments:
  local:
    variables:
      admin_email: admin@example.com
      admin_password: admin123
```

```kest
POST /api/login
{
  "email": "{{admin_email}}",
  "password": "{{admin_password}}"
}
```

**不好的做法** ❌:
```kest
POST /api/login
{
  "email": "admin@example.com",  # 硬编码
  "password": "admin123"
}
```

### 2. 使用动态变量避免冲突

**好的做法** ✅:
```kest
POST /api/users
{
  "username": "testuser{{$timestamp}}",
  "email": "test{{$randomInt}}@example.com"
}
```

**不好的做法** ❌:
```kest
POST /api/users
{
  "username": "testuser",  # 重复运行会冲突
  "email": "test@example.com"
}
```

### 3. 捕获所有需要的 ID

**好的做法** ✅:
```kest
POST /api/projects
{
  "name": "Test Project"
}

[Captures]
project_id: data.id
project_name: data.name
created_at: data.created_at

# 后续步骤可以使用
GET /api/projects/{{project_id}}
```

**不好的做法** ❌:
```kest
POST /api/projects
{
  "name": "Test Project"
}

# 没有捕获 ID，后续步骤无法使用
GET /api/projects/123  # 硬编码 ID
```

---

## 🎯 断言最佳实践

### 1. 使用范围断言而非精确值

**好的做法** ✅:
```kest
[Asserts]
status >= 200
status < 300
body.data exists
duration < 1000ms
```

**不好的做法** ❌:
```kest
[Asserts]
status == 200  # 太严格，201 也是成功
body.code == 0  # API 可能返回不同的成功码
```

### 2. 检查关键字段存在性

**好的做法** ✅:
```kest
[Asserts]
body.data.id exists
body.data.name exists
body.data.email exists
```

**不好的做法** ❌:
```kest
[Asserts]
body.data.id == "123"  # 精确值会变化
body.data.name == "Test"  # 可能有时间戳后缀
```

### 3. 使用变量进行比较

**好的做法** ✅:
```kest
POST /api/users
{
  "username": "user{{$timestamp}}"
}

[Captures]
created_username: data.username

GET /api/users/profile
[Asserts]
body.data.username == "{{created_username}}"
```

---

## 🔄 数据流转最佳实践

### 1. 链式测试

```kest
## Step 1: 创建资源
POST /api/projects
[Captures]
project_id: data.id

## Step 2: 使用资源
GET /api/projects/{{project_id}}

## Step 3: 更新资源
PUT /api/projects/{{project_id}}

## Step 4: 删除资源
DELETE /api/projects/{{project_id}}

## Step 5: 验证删除
GET /api/projects/{{project_id}}
[Asserts]
status == 404
```

### 2. 多资源关联

```kest
## 创建项目
POST /api/projects
[Captures]
project_id: data.id

## 创建环境（关联项目）
POST /api/projects/{{project_id}}/environments
[Captures]
env_id: data.id

## 创建测试用例（关联项目和环境）
POST /api/projects/{{project_id}}/test-cases
{
  "environment_id": "{{env_id}}"
}
[Captures]
testcase_id: data.id
```

---

## 🧹 清理步骤最佳实践

### 1. 总是包含清理步骤

**好的做法** ✅:
```kest
## 创建测试数据
POST /api/items
[Captures]
item_id: data.id

## 执行测试
GET /api/items/{{item_id}}

## 清理 - 删除测试数据
DELETE /api/items/{{item_id}}
```

### 2. 逆序清理

```kest
## 创建顺序
POST /api/projects
[Captures]
project_id: data.id

POST /api/projects/{{project_id}}/categories
[Captures]
category_id: data.id

## 清理顺序（逆序）
DELETE /api/projects/{{project_id}}/categories/{{category_id}}
DELETE /api/projects/{{project_id}}
```

### 3. 验证清理成功

```kest
## 删除资源
DELETE /api/projects/{{project_id}}
[Asserts]
status >= 200
status < 300

## 验证删除
GET /api/projects/{{project_id}}
[Asserts]
status == 404
```

---

## ⚡ 性能测试最佳实践

### 1. 添加性能断言

```kest
GET /api/users
[Asserts]
status == 200
duration < 500ms  # 列表查询应该快速

GET /api/users/{{user_id}}/details
[Asserts]
status == 200
duration < 1000ms  # 详情查询可以稍慢
```

### 2. 使用并行模式

```bash
# 串行执行（默认）
kest run test.flow.md

# 并行执行（更快）
kest run test.flow.md --parallel --jobs 8
```

---

## 🔒 安全测试最佳实践

### 1. 测试未授权访问

```kest
## 正常访问
GET /api/admin/users
Authorization: Bearer {{admin_token}}
[Asserts]
status == 200

## 未授权访问
GET /api/admin/users
# 不带 Authorization header
[Asserts]
status == 401
```

### 2. 测试权限边界

```kest
## 用户 A 创建资源
POST /api/projects
Authorization: Bearer {{user_a_token}}
[Captures]
project_id: data.id

## 用户 B 尝试访问（应该失败）
GET /api/projects/{{project_id}}
Authorization: Bearer {{user_b_token}}
[Asserts]
status == 403
```

---

## 📦 模块化最佳实践

### 1. 拆分大型 Flow

**不好的做法** ❌:
```
01-complete-test.flow.md  (100+ 步骤)
```

**好的做法** ✅:
```
01-auth-flow.flow.md          (登录、注册)
02-project-crud.flow.md       (项目 CRUD)
03-environment-crud.flow.md   (环境 CRUD)
99-smoke-test.flow.md         (核心功能快速验证)
```

### 2. 创建可重用的前置条件

```kest
# 在多个 flow 中重用
## 前置条件：登录
POST /api/login
{
  "username": "{{test_user}}",
  "password": "{{test_password}}"
}

[Captures]
access_token: data.token

[Asserts]
status == 200
body.data.token exists
```

---

## 🐛 调试最佳实践

### 1. 使用 kest show last

```bash
# 运行测试
kest run test.flow.md

# 查看最后一次请求详情
kest show last
```

### 2. 使用 kest history

```bash
# 查看所有历史请求
kest history

# 查看特定记录
kest show 123
```

### 3. 启用日志

```yaml
# .kest/config.yaml
log_enabled: true
```

日志保存在 `.kest/logs/` 目录。

---

## 📊 测试报告最佳实践

### 1. 使用描述性的步骤名称

**好的做法** ✅:
```markdown
## Step 1: 用户注册 - 创建新账号
## Step 2: 用户登录 - 获取访问令牌
## Step 3: 获取用户资料 - 验证登录状态
```

**不好的做法** ❌:
```markdown
## Step 1
## Step 2
## Step 3
```

### 2. 添加注释说明

```kest
## Step 5: 创建项目
# 注意：项目名称使用时间戳避免冲突
# 预期：返回 201 Created 和项目 ID

POST /api/projects
{
  "name": "Project {{$timestamp}}"
}
```

---

## 🎓 常见错误和解决方案

### 错误 1: 变量未定义

**问题**:
```
Request: GET /api/projects/{{project_id}}
Response: 400 Bad Request
```

**解决**:
```kest
# 确保在使用前捕获变量
POST /api/projects
[Captures]
project_id: data.id  # ← 必须先捕获

GET /api/projects/{{project_id}}
```

### 错误 2: 断言过于严格

**问题**:
```kest
[Asserts]
status == 200  # 失败：实际返回 201
```

**解决**:
```kest
[Asserts]
status >= 200
status < 300  # 接受所有 2xx 状态码
```

### 错误 3: 忘记清理测试数据

**问题**: 多次运行测试导致数据冲突

**解决**:
```kest
## 创建
POST /api/users
{
  "username": "user{{$timestamp}}"  # 使用动态值
}
[Captures]
user_id: data.id

## 清理
DELETE /api/users/{{user_id}}
```

---

## 🚀 高级技巧

### 1. 条件断言（未来功能）

```kest
[Asserts]
if body.data.type == "premium":
  body.data.features.length > 10
else:
  body.data.features.length > 3
```

### 2. 循环测试（未来功能）

```kest
[Loop]
count: 10

POST /api/items
{
  "name": "Item {{$loop_index}}"
}
```

### 3. 数据驱动测试（未来功能）

```kest
[DataSource]
file: test-data.csv

POST /api/users
{
  "username": "{{data.username}}",
  "email": "{{data.email}}"
}
```

---

## 📚 参考资源

- Kest CLI 文档: `.kest/flow/README.md`
- 改进记录: `/Users/stark/item/kest/kest-cli/IMPROVEMENTS.md`
- 测试报告: `.kest/flow/FINAL_TEST_REPORT.md`
- 示例 Flow: `.kest/flow/99-working-smoke-test.flow.md`

---

**最后更新**: 2026-02-02  
**版本**: v1.0  
**状态**: ✅ 已验证
