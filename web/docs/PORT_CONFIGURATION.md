# 端口配置指南

## 🚀 快速开始

### 1. 复制环境变量文件

```bash
cp .env.example .env.local
```

### 2. 编辑端口配置

编辑 `.env.local` 文件：

```bash
# Server Configuration
PORT="3000"              # 开发服务器端口
HOSTNAME="localhost"     # 服务器主机名

# Frontend Configuration
NEXT_PUBLIC_APP_URL="http://localhost:3000"
```

## 📋 可用脚本

### 开发模式

```bash
# 默认开发模式（端口 3000）
pnpm dev

# 自定义端口开发模式
pnpm dev:custom
```

### 生产模式

```bash
# 默认生产模式（使用环境变量端口）
pnpm start

# 自定义端口生产模式
pnpm start:custom
```

## 🔧 常见配置示例

### 本地开发

```bash
# .env.local
PORT="3000"
HOSTNAME="localhost"
NEXT_PUBLIC_APP_URL="http://localhost:3000"
```

### 团队开发（避免端口冲突）

```bash
# .env.local
PORT="3001"
HOSTNAME="localhost"
NEXT_PUBLIC_APP_URL="http://localhost:3001"
```

### 内网部署

```bash
# .env.local
PORT="8080"
HOSTNAME="0.0.0.0"
NEXT_PUBLIC_APP_URL="http://192.168.1.100:8080"
```

## 🌐 环境变量说明

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `PORT` | 服务器端口 | `3000` | `3001` |
| `HOSTNAME` | 服务器主机名 | `localhost` | `0.0.0.0` |
| `NEXT_PUBLIC_APP_URL` | 前端应用URL | `http://localhost:3000` | `http://localhost:3001` |

## 🚨 注意事项

1. **端口冲突** - 确保选择的端口没有被其他应用占用
2. **防火墙** - 如果使用 `0.0.0.0`，确保防火墙允许相应端口
3. **环境变量** - 生产环境请使用 `.env.local` 而不是 `.env.example`
4. **URL匹配** - 确保 `NEXT_PUBLIC_APP_URL` 与实际端口和主机名匹配

## 🔍 故障排除

### 端口被占用

```bash
# 查看端口占用
lsof -i :3000

# 杀死占用进程
kill -9 <PID>
```

### 无法访问

- 检查 `HOSTNAME` 设置：
  - `localhost` - 仅本机访问
  - `0.0.0.0` - 允许外部访问
- 检查防火墙设置
- 确认 `NEXT_PUBLIC_APP_URL` 正确
