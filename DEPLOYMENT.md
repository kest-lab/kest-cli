# Kest API 部署指南

## 📋 部署前检查清单

通过当前的 Dockerfile 部署后，API **可以访问**，但需要正确配置环境变量。

### ✅ Dockerfile 已就绪
- ✅ 纯 API 模式（不包含前端）
- ✅ 多阶段构建优化
- ✅ 支持云平台环境变量
- ✅ 健康检查配置

### ⚠️ 必须配置的环境变量

#### 1. 数据库配置（必需）

**PostgreSQL（推荐）**:
```bash
DB_DRIVER=postgres
DB_HOST=your-postgres-host.com
DB_PORT=5432
DB_NAME=kest
DB_USERNAME=your-username
DB_PASSWORD=your-password
DB_SSL_MODE=require
```

**MySQL**:
```bash
DB_DRIVER=mysql
DB_HOST=your-mysql-host.com
DB_PORT=3306
DB_NAME=kest
DB_USERNAME=your-username
DB_PASSWORD=your-password
```

#### 2. 安全密钥（必需）

生成强密钥：
```bash
openssl rand -base64 32
```

配置：
```bash
APP_KEY=your-generated-key-here
JWT_SECRET=your-generated-jwt-secret-here
JWT_EXPIRE_DAYS=7
```

#### 3. 服务器配置

```bash
APP_ENV=production
GIN_MODE=release
APP_DEBUG=false
```

---

## 🚀 部署到云平台

### Render 部署

1. **创建 Web Service**
   - Repository: `kest-labs/kest`
   - Branch: `main`
   - Build Command: 自动检测 Dockerfile
   - Start Command: 自动（使用 Dockerfile CMD）

2. **配置环境变量**
   
   在 Render Dashboard 中添加：
   ```
   DB_DRIVER=postgres
   DB_HOST=<从 Render PostgreSQL 获取>
   DB_PORT=5432
   DB_NAME=kest
   DB_USERNAME=<从 Render PostgreSQL 获取>
   DB_PASSWORD=<从 Render PostgreSQL 获取>
   
   APP_KEY=<生成的密钥>
   JWT_SECRET=<生成的密钥>
   APP_ENV=production
   GIN_MODE=release
   ```

3. **添加 PostgreSQL 数据库**
   - 在 Render 创建 PostgreSQL 实例
   - 自动获取连接信息
   - 配置到环境变量

### Zeabur 部署

1. **导入项目**
   - 连接 GitHub 仓库
   - 选择 `kest` 项目
   - Zeabur 自动检测 Dockerfile

2. **添加 PostgreSQL 服务**
   - 在同一项目中添加 PostgreSQL
   - Zeabur 自动注入数据库环境变量

3. **配置环境变量**
   ```
   APP_KEY=<生成的密钥>
   JWT_SECRET=<生成的密钥>
   APP_ENV=production
   GIN_MODE=release
   ```

### Railway 部署

1. **New Project from GitHub**
   - 选择 `kest` 仓库
   - Railway 自动检测 Dockerfile

2. **添加 PostgreSQL**
   - Add Plugin → PostgreSQL
   - 自动配置数据库连接

3. **配置环境变量**
   - 在 Variables 标签页添加必需的环境变量

---

## 🔍 部署后验证

### 1. 健康检查

```bash
curl https://your-app.com/v1/health
```

**期望响应**:
```json
{
  "status": "ok",
  "version": "v1"
}
```

### 2. 数据库连接检查

```bash
curl https://your-app.com/health
```

应该显示数据库连接正常。

### 3. API 端点测试

**注册用户**:
```bash
curl -X POST https://your-app.com/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test@123"
  }'
```

**登录**:
```bash
curl -X POST https://your-app.com/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test@123"
  }'
```

---

## 📊 可访问的 API 端点

部署后，以下端点可以正常访问：

### 公开端点
- `GET /v1/health` - 健康检查
- `POST /v1/register` - 用户注册
- `POST /v1/login` - 用户登录
- `POST /v1/password/reset` - 密码重置

### 认证端点（需要 JWT Token）
- `GET /v1/users/profile` - 获取用户信息
- `GET /v1/projects` - 项目列表
- `POST /v1/projects` - 创建项目
- `GET /v1/projects/:id/environments` - 环境列表
- `POST /v1/projects/:id/environments` - 创建环境
- `GET /v1/projects/:id/categories` - 分类列表
- `POST /v1/projects/:id/categories` - 创建分类
- 等等...

---

## ⚠️ 常见问题

### Q1: 部署后 API 返回 500 错误

**原因**: 数据库连接失败

**解决**:
1. 检查数据库环境变量是否正确
2. 确认数据库服务正在运行
3. 检查网络连接和防火墙规则

### Q2: JWT Token 无效

**原因**: JWT_SECRET 未配置或不一致

**解决**:
1. 确保 `JWT_SECRET` 环境变量已设置
2. 重启服务使新配置生效

### Q3: CORS 错误

**原因**: 前端域名未在 CORS 白名单中

**解决**:
```bash
CORS_ALLOW_ORIGINS=https://your-frontend.com,https://app.kest.dev
```

### Q4: 数据库迁移

**首次部署需要运行迁移**:

部署后，数据库表会自动创建（通过 GORM AutoMigrate）。

如果需要手动迁移，可以：
1. 连接到数据库
2. 运行迁移工具（如果有）

---

## 🎯 性能优化建议

### 1. 数据库连接池
```bash
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300
```

### 2. 启用 Redis 缓存（可选）
```bash
REDIS_ENABLED=true
REDIS_HOST=your-redis-host
REDIS_PORT=6379
```

### 3. 启用 CDN（如果有静态资源）

### 4. 配置日志级别
```bash
LOG_LEVEL=info  # 生产环境使用 info 或 warn
LOG_FORMAT=json # JSON 格式便于日志分析
```

---

## 📝 总结

### ✅ 可以部署并访问 API

**前提条件**:
1. ✅ 配置数据库连接
2. ✅ 设置 JWT 密钥
3. ✅ 配置基本环境变量

**部署后**:
- ✅ 所有 API 端点可以正常访问
- ✅ 健康检查正常
- ✅ 用户认证功能正常
- ✅ CRUD 操作正常

**不包含**:
- ❌ 前端界面（需要单独部署）
- ❌ WebSocket 支持（如果需要）

---

## 🔗 相关文档

- 环境变量模板: `.env.production.example`
- API 文档: `/swagger` (部署后访问)
- Flow 测试: `api/.kest/flow/`

---

**部署成功后，你的 API 将完全可用！** 🎉
