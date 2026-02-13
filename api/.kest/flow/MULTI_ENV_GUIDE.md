# 🌍 多环境测试指南

## 环境配置

Kest 支持多环境配置，所有环境定义在 `.kest/config.yaml` 中：

```yaml
environments:
  # 本地开发环境
  local:
    base_url: http://127.0.0.1:8080
  
  # 开发环境
  dev:
    base_url: https://dev-api.kest.io
  
  # 预发布环境
  staging:
    base_url: https://staging-api.kest.io
  
  # 生产环境
  production:
    base_url: https://api.kest.io
```

---

## 环境切换

### 查看所有环境

```bash
kest env list
```

### 切换环境

```bash
# 切换到开发环境
kest env use dev

# 切换到预发布环境
kest env use staging

# 切换到生产环境
kest env use production

# 切换回本地环境
kest env use local
```

---

## 运行测试

### 在当前环境运行

```bash
# 使用当前激活的环境（默认是 local）
kest run .kest/flow/00-smoke-test.flow.md
```

### 在特定环境运行

```bash
# 方法1：先切换环境，再运行
kest env use staging
kest run .kest/flow/00-smoke-test.flow.md

# 方法2：运行后切回原环境
kest env use staging
kest run .kest/flow/00-smoke-test.flow.md
kest env use local
```

---

## Flow 文件编写规范

### ✅ 正确：使用相对路径

```kest
# Health check
GET /health

# API endpoints
POST /api/v1/register
POST /api/v1/login
GET /api/v1/users/profile
```

**优点**：
- 自动适配不同环境的 base_url
- Flow 文件可以在任何环境运行
- 无需修改代码

### ❌ 错误：硬编码 URL

```kest
# ❌ 不要这样做
GET http://127.0.0.1:8080/health
POST http://127.0.0.1:8080/api/v1/login
```

**缺点**：
- 只能在本地运行
- 无法在其他环境测试
- 需要手动修改 URL

---

## 完整示例

### 1. 本地开发测试

```bash
# 确保在 local 环境
kest env use local

# 运行所有测试
./.kest/flow/run-all-flows.sh
```

### 2. 预发布环境验证

```bash
# 切换到 staging
kest env use staging

# 运行冒烟测试
kest run .kest/flow/00-smoke-test.flow.md

# 运行完整测试
kest run .kest/flow/01-auth-flow.flow.md
kest run .kest/flow/02-project-flow.flow.md
```

### 3. 生产环境健康检查

```bash
# 切换到 production
kest env use production

# 只运行健康检查（不要运行会创建数据的测试！）
kest get /health -a "status=200"
```

---

## 环境变量隔离

每个环境的变量（如 token、ID）是**独立存储**的：

```bash
# 在 local 环境登录
kest env use local
kest post /api/v1/login -d '{"username":"test"}' -c "token=data.access_token"

# 在 staging 环境登录（不会影响 local 的 token）
kest env use staging
kest post /api/v1/login -d '{"username":"test"}' -c "token=data.access_token"

# 查看当前环境的变量
kest vars
```

---

## 最佳实践

### 1. 开发阶段
- 使用 `local` 环境
- 频繁运行测试
- 快速迭代

### 2. 提交前验证
- 切换到 `dev` 环境
- 运行完整测试套件
- 确保所有测试通过

### 3. 发布前检查
- 切换到 `staging` 环境
- 运行完整测试
- 验证新功能

### 4. 生产监控
- 使用 `production` 环境
- **只运行只读测试**（GET 请求）
- 定期健康检查

---

## 注意事项

⚠️ **生产环境警告**：
- 不要在生产环境运行会创建/修改/删除数据的测试
- 只运行健康检查和只读查询
- 使用专门的测试账号

⚠️ **环境隔离**：
- 每个环境的数据是独立的
- 变量不会跨环境共享
- 确保在正确的环境运行测试

---

**Keep Every Step Tested! 🦅**
