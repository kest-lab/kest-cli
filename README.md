# Kest - API Testing Platform for Vibe Coding

<div align="center">

**Test APIs at the Speed of Thought**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![Node Version](https://img.shields.io/badge/Node-20+-339933?logo=node.js)](https://nodejs.org/)

The only API testing platform designed from the ground up for **Vibe Coding**—where velocity meets quality, and documentation writes itself.

[🚀 Quick Start](#-quick-start) • [📖 Documentation](#-documentation) • [🤝 Contributing](#-contributing)

</div>

---

## 🎯 What is Kest?

Kest is a modern API testing platform built for AI-assisted development workflows. It combines:

- **🌐 Cloud Platform** (web/) - Centralized API documentation and team collaboration
- **⚡ CLI Tool** - Local-first test execution with blazing speed  
- **🔌 MCP Integration** - Native support for AI editors (Cursor, Windsurf, Cline)

### Built for Vibe Coding

Traditional tools weren't designed for AI-assisted development. Kest was.

- **Markdown-Native**: Write tests in Markdown—documentation and tests are one
- **Local-First**: Zero network latency, native-speed execution
- **Git-Native**: Tests version alongside code, auto-integrated into CI/CD
- **MCP-Ready**: Plug into AI editors via Model Context Protocol

---

## 📁 Repository Structure

```
kest/
├── web/              # Next.js web platform (TypeScript)
│   ├── src/
│   ├── public/
│   └── package.json
│
├── api/              # Go backend service
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── go.mod
│
├── cli/              # Kest CLI tool (separate repo: kest-lab/kest-cli)
│
├── docker/           # Docker configurations
│   ├── web.Dockerfile
│   ├── api.Dockerfile
│   └── docker-compose.yml
│
└── docs/             # Documentation
```

---

## 🚀 Quick Start

### Prerequisites

- **Node.js** 20+ (for web platform)
- **Go** 1.22+ (for API backend)
- **Docker** & Docker Compose (for deployment)
- **PostgreSQL** 15+ (for production)

### Local Development

#### 1. Clone the Repository

```bash
git clone git@github.com:kest-lab/kest.git
cd kest
```

#### 2. Setup Web Platform

```bash
cd web
npm install
cp .env.example .env.local
# Edit .env.local with your configuration
npm run dev
```

The web platform will be available at `http://localhost:3000`

#### 3. Setup API Backend

```bash
cd api
go mod download
cp .env.example .env
# Edit .env with your configuration
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### Docker Deployment

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

---

## 📖 Documentation

### Core Features

#### Cloud Platform (web/)
- **API Documentation Editor**: OpenAPI 3.0 compatible
- **Multi-Level Management**: Projects, Environments, Categories
- **Team Collaboration**: Real-time editing, permission control
- **Version History**: Track changes and rollback

#### CLI Tool
- **Local Execution**: Millisecond-level test completion
- **Session Logging**: Automatic request/response recording
- **JSON Formatting**: Pretty-print built-in
- **Git Integration**: Tests as code

#### MCP Integration
- **AI Editor Support**: Cursor, Windsurf, Cline
- **Skill System**: Shareable coding standards
- **Context Protocol**: Native integration

### Environment Variables

#### Web Platform (`web/.env.local`)

```env
# API Connection
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_SITE_URL=http://localhost:3000

# Database (if needed for serverless)
DATABASE_URL=postgresql://user:pass@localhost:5432/kest

# Auth (configure based on your setup)
NEXTAUTH_SECRET=your-secret-key
NEXTAUTH_URL=http://localhost:3000
```

#### API Backend (`api/.env`)

```env
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=kest
DB_USER=kest_user
DB_PASSWORD=your_password

# JWT
JWT_SECRET=your-jwt-secret
JWT_EXPIRATION=24h

# CORS
ALLOWED_ORIGINS=http://localhost:3000,https://kest.example.com
```

---

## 🛠️ Development

### Web Platform Stack

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS + shadcn/ui
- **State**: Zustand + React Query
- **i18n**: next-intl
- **Theme**: OKLCH color system

### API Backend Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15+
- **Architecture**: Domain-Driven Design (DDD)

### Code Standards

- **Frontend**: ESLint + Prettier
- **Backend**: golangci-lint
- **Commits**: Conventional Commits
- **Branching**: Git Flow

---

## 🐳 Deployment

### Production Deployment

```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose -f docker-compose.prod.yml up -d --scale api=3
```

### Cloud Platforms

#### Vercel (Web)
```bash
cd web
vercel --prod
```

#### Google Cloud Run (API)
```bash
cd api
gcloud run deploy kest-api \
  --source . \
  --platform managed \
  --region us-central1
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📜 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🔗 Links

- **Website**: https://kest.example.com
- **Documentation**: https://docs.kest.example.com
- **CLI Repository**: https://github.com/kest-lab/kest-cli
- **Issue Tracker**: https://github.com/kest-lab/kest/issues
- **Discussions**: https://github.com/kest-lab/kest/discussions

---

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=kest-lab/kest&type=Date)](https://star-history.com/#kest-lab/kest&Date)

---

<div align="center">

**Built with ❤️ for Vibe Coding**

[⬆ Back to Top](#kest---api-testing-platform-for-vibe-coding)

</div>
