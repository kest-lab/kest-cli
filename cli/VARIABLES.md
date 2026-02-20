# Kest 变量系统完整指南

## 📋 目录

- [变量语法](#变量语法)
- [变量优先级](#变量优先级)
- [默认值语法](#默认值语法)
- [内置变量](#内置变量)
- [变量来源](#变量来源)
- [严格模式](#严格模式)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 变量语法

### 基本语法

在 Kest Flow 文件中使用双花括号引用变量：

```markdown
### Step: Login
POST /api/login
```json
{
  "username": "{{username}}",
  "password": "{{password}}"
}
```
```

### 默认值语法 (v1.1+)

使用管道符号 `|` 提供默认值：

```markdown
### Step: Login
POST /api/login
```json
{
  "username": "{{username | default: \"admin\"}}",
  "password": "{{password | default: \"Admin@123\"}}"
}
```
```

**优点**：
- 减少 `--var` 参数
- 测试更方便
- 文档即配置

---

## 变量优先级

变量解析按以下优先级（从高到低）：

### 1. CLI 参数 `--var` (最高优先级)

```bash
kest run flow.md --var username=test --var password=Test@123
```

**特点**：
- 覆盖所有其他来源
- 适合临时覆盖
- CI/CD 环境注入

### 2. Flow 内捕获 `[Captures]`

```markdown
### Step: Login
POST /api/login

[Captures]
- token = data.token
- user_id = data.user.id

### Step: Get Profile (使用捕获的变量)
GET /api/profile
Authorization: Bearer {{token}}
```

**特点**：
- 步骤执行时动态捕获
- 作用域：当前 flow 执行上下文
- 后续步骤可使用

### 3. 环境配置 `config.yaml` 中的 `environments.*.variables`

```yaml
# .kest/config.yaml
environments:
  dev:
    base_url: http://localhost:3000
    variables:
      api_key: dev_key_123
      db_name: test_db
  
  prod:
    base_url: https://api.example.com
    variables:
      api_key: prod_key_456
      db_name: production_db
```

```bash
# 切换环境
kest env set dev
kest run flow.md  # 使用 dev 环境的变量
```

**特点**：
- 环境切换时自动加载
- 作用域：当前环境
- 适合环境特定配置

### 4. 全局配置 `config.yaml` 中的 `variables`

```yaml
# .kest/config.yaml
variables:
  base_url: http://localhost:3000
  timeout: 5000
  retry_count: 3

environments:
  prod:
    variables:
      base_url: https://api.example.com  # 覆盖全局 base_url
```

**特点**：
- 所有环境共享
- 作用域：项目级别
- 适合通用配置

### 5. 默认值 `{{var | default: "value"}}` (最低优先级)

```markdown
{{username | default: "admin"}}
```

**特点**：
- 仅在变量未定义时使用
- 内嵌在 flow 文件中
- 提供后备值

---

## 优先级示例

```yaml
# .kest/config.yaml
variables:
  api_key: global_key  # 优先级 4

environments:
  dev:
    variables:
      api_key: dev_key  # 优先级 3
```

```markdown
### Step: API Call
GET /api/data
Authorization: Bearer {{api_key | default: "fallback_key"}}
```

**不同场景下的值**：

```bash
# 场景 1: 默认（dev 环境）
$ kest env set dev
$ kest run flow.md
# api_key = "dev_key" (来自环境配置)

# 场景 2: CLI 覆盖
$ kest run flow.md --var api_key=cli_key
# api_key = "cli_key" (CLI 参数最高优先级)

# 场景 3: 无环境配置
$ kest env set staging  # staging 环境没有 api_key
$ kest run flow.md
# api_key = "global_key" (来自全局配置)

# 场景 4: 完全未定义
# 假设删除了所有配置中的 api_key
$ kest run flow.md
# api_key = "fallback_key" (使用默认值)
```

---

## 默认值语法

### 支持的格式

```markdown
# 字符串默认值
{{password | default: "Admin@123"}}

# 数字默认值（作为字符串）
{{port | default: "8080"}}

# URL 默认值
{{base_url | default: "http://localhost:3000"}}

# 空字符串默认值
{{optional_field | default: ""}}
```

### 注意事项

1. **引号是必需的**：
   ```markdown
   ✅ {{var | default: "value"}}
   ❌ {{var | default: value}}
   ```

2. **空格可选**：
   ```markdown
   ✅ {{var | default: "value"}}
   ✅ {{var|default:"value"}}
   ✅ {{var  |  default:  "value"}}
   ```

3. **不支持嵌套**：
   ```markdown
   ❌ {{var | default: "{{other_var}}"}}
   ```

---

## 内置变量

Kest 提供以下内置变量：

### `$timestamp`

当前 Unix 时间戳（秒）

```markdown
### Step: Create Record
POST /api/records
```json
{
  "name": "Record_{{$timestamp}}",
  "created_at": {{$timestamp}}
}
```
```

**输出示例**：
```json
{
  "name": "Record_1708444800",
  "created_at": 1708444800
}
```

### `$randomInt`

随机整数（0-9999）

```markdown
### Step: Create User
POST /api/users
```json
{
  "username": "user_{{$randomInt}}",
  "email": "user{{$randomInt}}@example.com"
}
```
```

**输出示例**：
```json
{
  "username": "user_7234",
  "email": "user7234@example.com"
}
```

---

## 变量来源

### 1. 命令行注入

```bash
# 单个变量
kest run flow.md --var api_key=secret

# 多个变量
kest run flow.md \
  --var username=admin \
  --var password=Admin@123 \
  --var env=production
```

### 2. 配置文件

```yaml
# .kest/config.yaml
project_id: my-api-project
active_env: dev

variables:
  timeout: 5000
  retry_count: 3

environments:
  dev:
    base_url: http://localhost:3000
    variables:
      api_key: dev_key_123
      debug: true
  
  staging:
    base_url: https://staging.api.com
    variables:
      api_key: staging_key_456
      debug: false
  
  prod:
    base_url: https://api.com
    variables:
      api_key: prod_key_789
      debug: false
```

### 3. Flow 内捕获

```markdown
### Step 1: Login
POST /api/login
```json
{"username": "admin", "password": "Admin@123"}
```

[Captures]
- token = data.token
- user_id = data.user.id
- expires_at = data.expires_at

### Step 2: Get User Profile
GET /api/users/{{user_id}}
Authorization: Bearer {{token}}

[Captures]
- username = data.username
- email = data.email
```

### 4. 环境变量（未来支持）

```bash
export KEST_API_KEY=secret
kest run flow.md
# 在 flow 中使用 {{KEST_API_KEY}}
```

---

## 严格模式

### 启用严格验证

```bash
kest run flow.md --strict
```

**行为**：
- 在执行前验证所有变量
- 未定义的变量（无默认值）会导致错误
- 避免无意义的 API 请求

### 示例

```markdown
### Step: Login
POST /api/login
```json
{
  "username": "{{username}}",
  "password": "{{password}}"
}
```
```

**不使用 --strict**：
```bash
$ kest run flow.md
# 发送请求: {"username": "{{username}}", "password": "{{password}}"}
# 服务器返回: 401 Unauthorized
# 用户困惑：是变量问题还是密码错误？
```

**使用 --strict**：
```bash
$ kest run flow.md --strict
❌ Error: Required variables not provided: username, password

Hint: Use one of the following:
  1. --var username=<value> --var password=<value>
  2. Add to config.yaml:
     environments:
       dev:
         variables:
           username: "admin"
           password: "Admin@123"
  3. Use default values:
     {{username | default: "admin"}}
     {{password | default: "Admin@123"}}
```

---

## 最佳实践

### 1. 使用默认值简化测试

```markdown
### Step: Login
POST /api/login
```json
{
  "username": "{{username | default: \"admin\"}}",
  "password": "{{password | default: \"Admin@123\"}}"
}
```

[Captures]
- token = data.token
```

**好处**：
- 无需每次传 `--var`
- 测试更快速
- 文档即配置

### 2. 环境特定配置

```yaml
# .kest/config.yaml
environments:
  dev:
    base_url: http://localhost:3000
    variables:
      db_name: test_db
      log_level: debug
  
  prod:
    base_url: https://api.example.com
    variables:
      db_name: production_db
      log_level: error
```

### 3. 敏感信息使用 CLI 注入

```bash
# 不要在配置文件中硬编码敏感信息
kest run flow.md --var api_key=$PROD_API_KEY
```

### 4. 使用 --strict 捕获错误

```bash
# 开发时使用严格模式
kest run flow.md --strict

# CI/CD 中也使用严格模式
kest run tests/ --strict --fail-fast
```

### 5. 组合使用 --fail-fast

```bash
# 快速失败，避免浪费时间
kest run flow.md --strict --fail-fast
```

---

## 常见问题

### Q1: 变量未替换，请求失败怎么办？

**问题**：
```bash
$ kest run flow.md
❌ Step failed: 401 Unauthorized
```

**解决方案**：
```bash
# 1. 使用 --debug-vars 查看变量解析
$ kest run flow.md --debug-vars

# 2. 使用 --strict 提前发现问题
$ kest run flow.md --strict

# 3. 使用 -v 查看详细信息
$ kest run flow.md -v
```

### Q2: 如何查看当前可用的变量？

```bash
# 查看当前项目和环境的变量
$ kest vars

Variables for project my-api (env: dev):
  api_key = dev_key_123
  base_url = http://localhost:3000
  token = eyJhbGc...
  user_id = 123
```

### Q3: 变量优先级不清楚怎么办？

参考本文档的 [变量优先级](#变量优先级) 部分，记住：

**CLI > Flow 捕获 > 环境配置 > 全局配置 > 默认值**

### Q4: 如何在不同环境间切换？

```bash
# 查看当前环境
$ kest env

# 切换环境
$ kest env set staging

# 运行测试
$ kest run flow.md  # 自动使用 staging 环境的变量
```

### Q5: 默认值语法不工作？

检查以下几点：

1. **引号是否正确**：
   ```markdown
   ✅ {{var | default: "value"}}
   ❌ {{var | default: value}}
   ```

2. **Kest 版本是否支持**：
   ```bash
   $ kest --version
   # 需要 v1.1.0 或更高版本
   ```

3. **语法是否正确**：
   ```markdown
   ✅ {{var | default: "value"}}
   ❌ {{var || default: "value"}}
   ❌ {{var | default = "value"}}
   ```

---

## 调试技巧

### 1. 使用 --debug-vars

```bash
$ kest run flow.md --debug-vars

📝 Variable Resolution Debug:

Step 1: Login
  Request Body (before):
    {"username": "{{username}}", "password": "{{password}}"}
  
  Available variables:
    ✅ base_url = http://localhost:3000 (from config.yaml)
    ❌ username = <not defined>
    ❌ password = <not defined>
  
  Request Body (after):
    {"username": "{{username}}", "password": "{{password}}"}
    ⚠️  Unresolved variables: username, password
```

### 2. 使用 --strict 提前验证

```bash
$ kest run flow.md --strict
❌ Error: Required variables not provided: username, password
```

### 3. 使用 -v 查看详细输出

```bash
$ kest run flow.md -v
⚠️  Warning: Undefined variables in URL: user_id
```

---

## 更新日志

### v1.1.0 (当前版本)

- ✅ 添加默认值语法支持 `{{var | default: "value"}}`
- ✅ 添加严格验证模式 `--strict`
- ✅ 添加 `--fail-fast` 模式
- ✅ 改进变量未定义时的错误提示

### v1.0.0

- ✅ 基本变量替换 `{{var}}`
- ✅ 内置变量 `$timestamp`, `$randomInt`
- ✅ 变量捕获 `[Captures]`
- ✅ 环境配置支持

---

## 相关文档

- [Flow 指南](FLOW_GUIDE.md) - Flow 文件编写指南
- [配置指南](GUIDE.md) - 项目配置详解
- [FAQ](FAQ.md) - 常见问题解答

---

**最后更新**: 2026-02-20  
**版本**: v1.1.0
