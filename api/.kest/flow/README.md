# 🦅 Kest API Flow Tests

This directory contains comprehensive flow test files for the Kest API project.

## 📁 Flow Files

| File | Description | Complexity | Duration |
|------|-------------|------------|----------|
| `00-smoke-test.flow.md` | Quick smoke test for critical endpoints | Low | ~5s |
| `01-auth-flow.flow.md` | Complete authentication lifecycle | Medium | ~8s |
| `02-project-flow.flow.md` | Project CRUD operations | Medium | ~12s |
| `03-apispec-flow.flow.md` | API specification management | High | ~15s |
| `04-performance-flow.flow.md` | Performance benchmarks with SLA | Medium | ~10s |
| `05-security-flow.flow.md` | Security and authorization tests | Medium | ~8s |

---

## 🚀 Quick Start

### Run Individual Flow

```bash
# Quick smoke test (recommended first)
kest run .kest/flow/00-smoke-test.flow.md

# Authentication flow
kest run .kest/flow/01-auth-flow.flow.md

# Project management
kest run .kest/flow/02-project-flow.flow.md

# API specifications
kest run .kest/flow/03-apispec-flow.flow.md

# Performance benchmarks
kest run .kest/flow/04-performance-flow.flow.md

# Security tests
kest run .kest/flow/05-security-flow.flow.md
```

### Run All Flows Sequentially

```bash
# Run all tests in order
for flow in .kest/flow/*.flow.md; do
  echo "Running $flow..."
  kest run "$flow"
done
```

### Run All Flows in Parallel

```bash
# Fast parallel execution (use with caution - may cause conflicts)
kest run .kest/flow/ --parallel --jobs 3
```

---

## 📋 Prerequisites

1. **Server Running**: Ensure the API server is running on `http://127.0.0.1:8080`
   ```bash
   # In kest-api directory
   go run cmd/server/main.go
   # OR
   ./test-server
   ```

2. **Database Ready**: Ensure PostgreSQL is running and migrations are applied
   ```bash
   make migrate-up
   ```

3. **Kest CLI Installed**: Ensure you have the latest Kest CLI
   ```bash
   cd kest-cli
   go build -o ~/go/bin/kest cmd/kest/main.go
   ```

---

## 🎯 Test Coverage

### Authentication (01-auth-flow.flow.md)
- ✅ User registration
- ✅ User login
- ✅ Get profile
- ✅ Update profile
- ✅ Change password
- ✅ Login with new password
- ✅ Get user info by ID

### Project Management (02-project-flow.flow.md)
- ✅ Create project
- ✅ List projects
- ✅ Get project details
- ✅ Update project
- ✅ Get project DSN
- ✅ Delete project
- ✅ Verify deletion

### API Specifications (03-apispec-flow.flow.md)
- ✅ Create API spec
- ✅ List API specs
- ✅ Get single spec
- ✅ Update spec
- ✅ Create example
- ✅ Get spec with examples
- ✅ Export specs
- ✅ Delete spec

### Performance (04-performance-flow.flow.md)
- ✅ Health check < 100ms
- ✅ Profile retrieval < 300ms
- ✅ Project list < 500ms
- ✅ Project creation < 1000ms
- ✅ All CRUD operations within SLA

### Security (05-security-flow.flow.md)
- ✅ Unauthorized access blocked
- ✅ Invalid token rejected
- ✅ Token validation
- ✅ Permission checks
- ✅ Non-existent resource handling

---

## 🛠 Troubleshooting

### Server Not Running
```bash
Error: connection refused

Solution:
cd /Users/stark/item/kest/kest-api
./test-server
```

### Database Connection Failed
```bash
Error: database connection failed

Solution:
# Check PostgreSQL is running
docker ps | grep postgres

# Run migrations
make migrate-up
```

### Token Expiration
```bash
Error: 401 Unauthorized

Solution:
# Tokens may expire - rerun the flow from the beginning
kest run .kest/flow/01-auth-flow.flow.md
```

### Port Already in Use
```bash
Error: bind: address already in use

Solution:
# Kill existing process
lsof -i :8080 -t | xargs kill -9

# Restart server
./test-server
```

---

## 📊 Expected Results

### Successful Run
```
🚀 Running 7 test(s) from 01-auth-flow.flow.md

╭─────────────────────────────────────────────────────────────────────╮
│                        TEST SUMMARY                                 │
├─────────────────────────────────────────────────────────────────────┤
│ ✓ POST     /v1/register                       234ms │
│ ✓ POST     /v1/login                          178ms │
│ ✓ GET      /v1/users/profile                   92ms │
│ ✓ PUT      /v1/users/profile                  156ms │
│ ✓ PUT      /v1/users/password                 189ms │
│ ✓ POST     /v1/login                          165ms │
│ ✓ GET      /v1/users/:id/info                  87ms │
├─────────────────────────────────────────────────────────────────────┤
│ Total: 7  │  Passed: 7  │  Failed: 0  │  Time: 1101ms │
╰─────────────────────────────────────────────────────────────────────╯

✓ All tests passed!
```

---

## 🔄 CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run Kest Flow Tests
  run: |
    cd kest-api
    kest run .kest/flow/00-smoke-test.flow.md
    kest run .kest/flow/01-auth-flow.flow.md
    kest run .kest/flow/02-project-flow.flow.md
```

---

## 📝 Writing New Flows

1. Create a new `.flow.md` file in this directory
2. Follow the naming convention: `##-name-flow.flow.md`
3. Use the existing flows as templates
4. Include proper documentation and assertions
5. Test locally before committing

Example structure:
```markdown
# 🎯 Your Flow Title

Description of what this flow tests.

---

## Step 1: Step Name

Description of this step.

\`\`\`kest
METHOD /path
Header: value

{
  "body": "data"
}

[Captures]
variable: path.to.value

[Asserts]
status == 200
duration < 500ms
\`\`\`
```

---

**Keep Every Step Tested! 🦅**
