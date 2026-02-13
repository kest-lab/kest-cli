---
name: deployment
description: Deployment workflows, containerization, and cloud-native practices
version: 1.0.0
category: DevOps
tags: [deployment, docker, k8s, cloud-run, ci-cd]
author: ZGO Team
updated: 2026-01-24
---

# Deployment Skill

## 📋 Purpose

This skill defines the deployment workflows and infrastructure standards for the ZGO project. it ensures that applications are delivered consistently, securely, and with high availability across environments.

## 🎯 When to Use

- Building and containerizing an application
- Setting up CI/CD pipelines
- Deploying to Google Cloud Run or Kubernetes
- Managing environment-specific configurations
- performing blue-green or canary deployments

## ⚙️ Prerequisites

- [ ] Docker installed locally
- [ ] Access to cloud provider (Google Cloud recommended)
- [ ] Understands image registry (GCR/AR)
- [ ] Knowledge of multi-stage Docker builds

---

## 🏗️ Containerization

We use **Multi-Stage Docker Builds** to ensure small, secure production images.

### 1. Multi-Stage Dockerfile (Standard)

```dockerfile
# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Production Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./main"]
```

### 2. Best Practices
- ✅ Use `.dockerignore` to keep image small
- ✅ Use `distroless` or `alpine` for production to reduce attack surface
- ✅ Set `CGO_ENABLED=0` for static binary
- ✅ Include non-root user in the final image

---

## 🛰️ Cloud Run Deployment (Recommended)

Google Cloud Run is the preferred platform for ZGO's serverless microservices.

### 1. Manual Deployment
```bash
# Build and push to registry
docker build -t gcr.io/[PROJECT_ID]/zgo-app:v1 .
docker push gcr.io/[PROJECT_ID]/zgo-app:v1

# Deploy to Cloud Run
gcloud run deploy zgo-app \
  --image gcr.io/[PROJECT_ID]/zgo-app:v1 \
  --region us-central1 \
  --allow-unauthenticated
```

### 2. AI-Assisted Deployment (MCP)
If the `cloudrun` MCP server is enabled, use it to deploy directly from the context.

---

## 🔄 CI/CD Workflow

We follow a **GitOps** approach using GitHub Actions.

### Pipeline Stages
1. **Lint**: Run golangci-lint
2. **Test**: Execute unit and integration tests
3. **Build**: Create Docker image
4. **Scan**: Vulnerability scanning (Trivy/Snyk)
5. **Deploy**: Auto-deploy to Staging; Manual trigger to Production

---

## ⚙️ Environment Management

Configurations must be injected via **Environment Variables** for 12-factor compliance.

| Key | Example (Dev) | Example (Prod) | Description |
|-----|---------------|----------------|-------------|
| `APP_ENV` | `development` | `production` | Environment toggle |
| `DB_DSN` | `localhost:5432` | `db.production.local` | Database connection |
| `LOG_LEVEL` | `debug` | `info` | Log verbosity |
| `SECRET_KEY` | `dev-secret` | `KMS-stored-secret` | Security token |

---

## 🩺 Health Checks & Monitoring

### 1. Standard Endpoints
- `/healthz`: Liveness check (checks process status)
- `/readyz`: Readiness check (checks DB/Redis connectivity)

### 2. Observability
- **Metrics**: Exported via Prometheus (default port: 9090)
- **Traces**: Exported to Cloud Trace / Jaeger
- **Logs**: Structured JSON logs sent to stdout

---

## ✅ Verification Checklist

- [ ] Dockerfile uses multi-stage builds.
- [ ] No secrets are hardcoded in the Dockerfile.
- [ ] `.dockerignore` excludes large/sensitive files (`vendor`, `.env`, `tmp`).
- [ ] Liveness/Readiness probes are implemented.
- [ ] Image size is under 100MB (target: ~30-50MB).
- [ ] CI/CD pipeline passes all stages.

---

## 🔧 Automation Scripts

Available in `scripts/`:
- [`deploy-cloudrun.sh`](./scripts/deploy-cloudrun.sh) - Helper for GCR deployments.
- [`verify-image.sh`](./scripts/verify-image.sh) - Checks image compliance & size.

---

## 📚 Complete Examples

- [**GitHub Actions Workflow YAML**](./examples/ci-cd-gh-actions.yaml)
- [**Full Production Dockerfile**](./examples/Dockerfile.production)
- [**Health Check Implementation**](./examples/health-check.go)

---

## 🔗 Related Skills

- [`coding-standards`](../coding-standards/): For log and config standards.
- [`api-development`](../api-development/): For health check endpoint standards.

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team
