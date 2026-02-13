# Kest 平台 CI/CD 实施工作文档

## 一、部署模式

**单体部署（Monolithic）**：Go 后端通过 `go:embed` 将 Vite 前端构建产物嵌入二进制，一个容器同时提供 API + 前端。

## 二、项目结构

| 组件 | 技术栈 | 说明 |
|------|--------|------|
| **API + 前端** | Go 1.23 + Gin + Vite/React（嵌入） | `cmd/main.go` 使用 `//go:embed web/dist` |
| **数据库** | PostgreSQL 16 | docker-compose 管理 |
| **Dockerfile** | `api/Dockerfile` | 4 阶段多级构建 |
| **CI/CD** | `.github/workflows/deploy.yml` | GitHub Actions → ghcr.io → SSH 部署 |

## 三、整体架构

```
push to main
     │
     ▼
GitHub Actions CI
     │
     └── Build Kest Image (前端+后端一体)
              │
              ▼
         Push to ghcr.io/<repo>:latest + :sha
              │
              ▼
         Deploy Job (SSH 到服务器)
              │
              ├── docker compose pull
              └── docker compose up -d
```

## 四、已完成的工作

### 4.1 清理冗余文件

- [x] 删除 `web/package-lock.json`，统一使用 pnpm
- [x] 删除 `docker/api.Dockerfile`（与 `api/Dockerfile` 重复）
- [x] 删除 `docker/api.dockerignore`
- [x] 删除 `docker/web.Dockerfile`（单体模式不需要独立前端容器）
- [x] 删除 `docker/web.dockerignore`

### 4.2 完善 `api/Dockerfile`

4 阶段多级构建，最大化缓存复用：

```
Stage 1: web-deps     → COPY package.json + pnpm-lock.yaml → pnpm install（npm 依赖缓存层）
Stage 2: web-builder  → COPY 前端源码 → pnpm build → 产出 web/dist
Stage 3: go-deps      → COPY go.mod + go.sum → go mod download（Go 依赖缓存层）
Stage 4: builder      → COPY Go 源码 + web/dist → go build（go:embed 嵌入前端）
Stage 5: runtime      → alpine + 二进制文件
```

**缓存策略**：依赖文件（`package.json`/`pnpm-lock.yaml`/`go.mod`/`go.sum`）单独 COPY 并安装，源码变更时依赖层直接命中缓存。

### 4.3 更新 `docker-compose.yml`

- 移除 `web` 服务（单体模式不需要）
- API build context 改为项目根目录 `.`，dockerfile 指向 `api/Dockerfile`
- 健康检查改为 `curl -f http://localhost:8080/health`
- `ALLOWED_ORIGINS` 默认值改为 `http://localhost:8080`

### 4.4 创建 `.github/workflows/deploy.yml`

- **触发**：push 到 `main` 分支
- **Job 1 `build-and-push`**：构建镜像 → 推送到 ghcr.io（`latest` + commit SHA 双标签）
- **Job 2 `deploy`**：SSH 到服务器 → `docker compose pull` → `docker compose up -d`
- **缓存加速**：Docker Buildx + `cache-from/cache-to: type=gha`

---

## 五、待完成的工作（需手动操作）

### 5.1 GitHub 仓库配置

在 **Settings → Secrets and variables → Actions → Secrets** 中添加：

| Secret | 说明 |
|--------|------|
| `SERVER_HOST` | 服务器 IP 或域名 |
| `SERVER_USER` | SSH 登录用户名 |
| `SERVER_SSH_KEY` | SSH 私钥内容（用于免密登录） |
| `DEPLOY_PATH`（可选） | 服务器上项目目录，默认 `/opt/kest` |

> `GITHUB_TOKEN` 由 GitHub 自动提供，无需手动配置。

### 5.2 服务器初始化

#### 环境要求

- Docker Engine 24+
- Docker Compose v2
- 能访问 ghcr.io

#### 初始化步骤

```bash
# 1. 登录 GitHub Container Registry
echo $GITHUB_PAT | docker login ghcr.io -u <github-username> --password-stdin

# 2. 创建项目目录
sudo mkdir -p /opt/kest
cd /opt/kest

# 3. 放置 docker-compose.yml 和 .env 文件

# 4. 首次启动
docker compose up -d
```

#### 服务器 `.env` 文件

```env
DB_PASSWORD=<生产数据库密码>
JWT_SECRET=<生产JWT密钥>
GIN_MODE=release
ALLOWED_ORIGINS=https://your-domain.com
```

### 5.3 验证

推送代码到 `main` 分支，检查：
1. GitHub Actions 是否成功构建并推送镜像
2. 服务器是否自动拉取并更新容器
3. 访问 `http://<server>:8080` 确认前后端正常工作
4. 访问 `http://<server>:8080/health` 确认健康检查通过

---

## 六、可选增强

- **Watchtower 自动更新**：在服务器上运行 [Watchtower](https://github.com/containrrr/watchtower) 容器，自动监听 Registry 中的镜像更新并重启容器，可替代 SSH deploy 步骤
- **多环境支持**：增加 `staging` 分支 → 部署到测试环境
- **镜像版本回滚**：每次构建同时打 `latest` 和 `sha` 标签，出问题时可快速回滚到指定 commit 的镜像
- **健康检查 + 通知**：deploy 后自动 curl 健康检查接口，失败时发送 Slack/飞书通知
- **Dependabot**：自动更新 Go / npm 依赖，保持安全
