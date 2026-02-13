# 🌐 线上API回归测试指南

## 📖 概述

本指南说明如何使用Kest Flow通过配置化方式对线上API进行回归测试。

---

## 🔧 配置说明

### 环境配置文件

配置文件位置：`.kest/config.yaml`

```yaml
environments:
  local:
    base_url: http://127.0.0.1:8080
  production:
    base_url: https://api.kest.dev/v1  # 线上API地址
```

### 关键配置项

- `active_env`: 当前激活的环境
- `base_url`: API基础URL，所有相对路径都会基于此URL
- `variables`: 环境特定的变量

---

## 🚀 快速开始

### 1. 查看可用环境

```bash
kest env list
```

输出示例：
```
Available environments:
  * local       (active)
    dev
    staging
    production
```

### 2. 切换到线上环境

```bash
kest env use production
```

### 3. 运行回归测试

```bash
kest run .kest/flow/production-regression.flow.md
```

### 4. 切换回本地环境

```bash
kest env use local
```

---

## 📝 编写可配置的Flow文件

### ✅ 正确做法：使用相对路径

```kest
# 健康检查
GET /health

# 用户注册
POST /register
Content-Type: application/json

{
  "username": "test_user",
  "password": "password123"
}

# 用户登录
POST /login
Content-Type: application/json

{
  "username": "test_user",
  "password": "password123"
}

[Captures]
access_token: data.access_token

# 获取项目列表
GET /projects
Authorization: Bearer {{access_token}}
```

**优点**：
- 自动适配不同环境的base_url
- 同一个Flow文件可以在任何环境运行
- 无需修改代码

### ❌ 错误做法：硬编码完整URL

```kest
# ❌ 不要这样做
GET https://api.kest.dev/v1/health
POST https://api.kest.dev/v1/register
```

**缺点**：
- 无法切换环境
- 每个环境需要单独的Flow文件
- 维护困难

---

## 🎯 完整测试流程

### 步骤1：准备测试环境

```bash
# 1. 确认当前环境
kest env list

# 2. 切换到目标环境
kest env use production
```

### 步骤2：运行测试

```bash
# 运行单个测试
kest run .kest/flow/production-regression.flow.md

# 运行所有测试
kest run .kest/flow/

# 并行运行（提高速度）
kest run .kest/flow/ --parallel --jobs 4
```

### 步骤3：查看测试结果

```bash
# 查看最新测试日志
kest logs

# 查看详细日志
cat .kest/logs/kest.log

# 查看测试报告
cat .kest/flow/PRODUCTION_REGRESSION_REPORT.md
```

### 步骤4：恢复环境

```bash
# 切换回本地环境
kest env use local
```

---

## 📊 测试结果解读

### 成功的测试

```
✓ 01:42:08 [GET] /health                           805ms
✓ 01:42:09 [POST] /login                            361ms
✓ 01:42:09 [GET] /users/profile                    301ms
```

- ✓ 表示测试通过
- 时间戳显示执行时间
- 响应时间在右侧显示

### 失败的测试

```
✗ 01:42:08 [POST] /register                         358ms
    Error: assertion failed: status == 200
    Response Body Sample:
      {
        "code": 0,
        "message": "created",
        ...
```

- ✗ 表示测试失败
- 显示失败原因
- 提供响应体样本

---

## 🔍 常见问题排查

### 问题1：所有测试都返回404

**原因**：base_url配置错误

**解决**：
```bash
# 检查配置
cat .kest/config.yaml

# 确认production环境的base_url是否正确
# 应该是：https://api.kest.dev/v1
# 而不是：https://api.kest.dev 或 https://api.kest.dev/api/v1
```

### 问题2：断言失败 - status in [200, 404]

**原因**：Kest不支持 `in` 语法

**解决**：
```kest
# ❌ 错误
[Asserts]
status in [200, 404]

# ✅ 正确
[Asserts]
status == 200
```

### 问题3：响应格式与预期不符

**原因**：API实际响应与文档不一致

**示例**：
```json
// 文档中的格式
{
  "data": {
    "items": [...],
    "pagination": {...}
  }
}

// 实际返回的格式
{
  "data": {
    "data": [...],
    "meta": {...}
  }
}
```

**解决**：根据实际响应调整断言

### 问题4：认证失败

**原因**：Token未正确捕获或传递

**解决**：
```kest
# 确保登录后捕获token
POST /login
...

[Captures]
access_token: data.access_token

# 在后续请求中使用
GET /projects
Authorization: Bearer {{access_token}}
```

---

## 📈 性能基准

### 预期响应时间

| 操作类型 | 预期时间 | 可接受时间 |
|---------|---------|-----------|
| 健康检查 | < 500ms | < 2000ms |
| 认证操作 | < 1000ms | < 3000ms |
| 查询操作 | < 500ms | < 2000ms |
| 创建操作 | < 1000ms | < 3000ms |
| 更新操作 | < 500ms | < 2000ms |
| 删除操作 | < 500ms | < 2000ms |

### 性能问题排查

如果响应时间超过预期：

1. **检查网络延迟**
   ```bash
   ping api.kest.dev
   ```

2. **检查API服务器负载**
   - 查看服务器监控
   - 检查数据库性能

3. **优化测试**
   - 使用并行执行
   - 减少不必要的断言

---

## 🛠️ 高级用法

### 使用环境变量

```yaml
# .kest/config.yaml
environments:
  production:
    base_url: https://api.kest.dev/v1
    variables:
      admin_username: admin
      admin_password: ${PROD_ADMIN_PASSWORD}  # 从环境变量读取
```

### 条件测试

```kest
# 只在特定环境运行的测试
GET /admin/debug
Authorization: Bearer {{access_token}}

[Asserts]
# 生产环境应该返回403
status == 403
```

### 数据驱动测试

```kest
# 测试多个用户
POST /register
Content-Type: application/json

{
  "username": "user_{{$randomInt}}",
  "email": "user_{{$timestamp}}@example.com",
  "password": "Test123!"
}
```

---

## 📋 测试清单

### 上线前检查

- [ ] 所有核心功能测试通过
- [ ] 性能测试达标
- [ ] 安全测试通过
- [ ] 错误处理正确
- [ ] 文档与实现一致

### 定期回归测试

- [ ] 每日：烟雾测试（核心功能）
- [ ] 每周：完整回归测试
- [ ] 发布前：全面测试
- [ ] 发布后：验证测试

---

## 🔐 安全注意事项

1. **不要在Flow文件中硬编码敏感信息**
   ```kest
   # ❌ 错误
   POST /login
   {
     "username": "admin",
     "password": "real_password_123"
   }
   
   # ✅ 正确
   POST /login
   {
     "username": "{{admin_username}}",
     "password": "{{admin_password}}"
   }
   ```

2. **使用测试账户**
   - 不要使用真实用户数据
   - 测试后清理测试数据

3. **限制生产环境测试**
   - 避免大量并发测试
   - 避免破坏性操作
   - 使用只读测试账户

---

## 📚 相关文档

- [Flow文件编写指南](./FLOW_BEST_PRACTICES.md)
- [多环境配置指南](./MULTI_ENV_GUIDE.md)
- [日志查看指南](./LOGS_GUIDE.md)
- [快速参考](./QUICK_REFERENCE.md)
- [生产环境测试报告](./PRODUCTION_REGRESSION_REPORT.md)

---

## 🎓 总结

通过配置化的方式进行线上回归测试的关键点：

1. ✅ 使用相对路径编写Flow文件
2. ✅ 在`.kest/config.yaml`中配置环境
3. ✅ 使用`kest env use`切换环境
4. ✅ 测试后切换回本地环境
5. ✅ 定期运行回归测试
6. ✅ 及时查看和分析测试结果

**一键回归测试命令**：
```bash
kest env use production && \
kest run .kest/flow/production-regression.flow.md && \
kest env use local
```
