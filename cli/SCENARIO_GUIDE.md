# Kest Scenario Guide

## What is a Scenario?

A **Scenario** is Kest's test case file format, using the `.kest` extension. It's a plain-text file describing a sequence of API test steps.

**Equivalent in other tools**:
- Postman → Collection
- Hurl → Test File
- k6 → Script
- **Kest → Scenario**

---

## Scenario File Formats

Kest supports two formats:
1. **`.kest` (CLI style)**: Minimal one-line-per-request format inherited from shell commands.
2. **`.flow.md` (Markdown style)**: Declarative format combining documentation and testing. Supports multi-line JSON and structured assertions. See [FLOW_GUIDE.md](FLOW_GUIDE.md).

---

### 1. Markdown Style (.flow.md) — Documentation as Tests

The recommended approach. Write tests like API documentation. Use `.flow.md` extension.

#### Syntax

Two syntax variants are supported:
- **New syntax (recommended)**: ` ```flow / ```step / ```edge ` for complete flow graphs
- **Legacy syntax (compatible)**: ` ```kest ` for single request blocks

**New syntax example:**
```flow
@flow id=user-flow
@name User Flow
@version 1.0
```

```step
@id login
@name Login
POST /api/v1/login
[Asserts]
status == 200
```

```edge
@from login
@to profile
@on success
```

**Legacy syntax example:**
```kest
# 1. First line is always METHOD URL
POST /api/v1/projects
X-User-ID: 100
Content-Type: application/json

# 2. After a blank line: Request Body (supports multi-line JSON)
{
  "name": "My Project",
  "description": "Created from Markdown"
}

# 3. Variable captures
[Captures]
project_id: data.id

# 4. Assertions
[Asserts]
status == 201
body.name == "My Project"
duration < 500ms
```

#### Running

```bash
kest run my-api-doc.flow.md
```

---

### 2. CLI Style (.kest) — Fast One-Liners

Best for: small, quick, one-off API calls.

```kest
# 1. Register a new user
POST /api/register -d '{"email":"test@example.com","password":"123456"}' -a "status==201"

# 2. Login and capture token
POST /api/login -d '{"email":"test@example.com","password":"123456"}' -c "token=data.token" -a "status==200"

# 3. Use token to get user info
GET /api/profile -H "Authorization: Bearer {{token}}" -a "status==200" -a "body.email==test@example.com"

# 4. Performance test: search must respond < 500ms
GET /api/search?q=test --max-time 500 -a "status==200"

# 5. Retry flaky endpoint
POST /api/webhook -d '{"event":"test"}' --retry 3 --retry-delay 1000
```

### Supported Command Format

```kest
# HTTP methods
GET /path
POST /path -d '{"key":"value"}'
PUT /path -d '{"key":"value"}'
DELETE /path
PATCH /path -d '{"key":"value"}'

# Headers
GET /path -H "Authorization: Bearer token" -H "X-Custom: value"

# Query parameters
GET /path -q "page=1" -q "limit=10"

# Variable capture
POST /login -c "token=data.token" -c "userId=data.user.id"

# Assertions
GET /users -a "status==200" -a "body.length==10"

# Performance assertion
GET /api --max-time 1000

# Retry mechanism
POST /api --retry 3 --retry-delay 1000

# gRPC call
grpc localhost:50051 package.Service/Method -d '{"field":"value"}'

# Streaming response
POST /chat -d '{"stream":true}' --stream
```

---

## 4 Ways to Create Scenarios

### Option 1: Manual (Recommended)

Best for: small projects, quick prototyping, custom tests.

```bash
cat > user-flow.kest << 'EOF'
# Full user flow test
POST /register -d '{"email":"new@test.com"}' -a "status==201"
POST /login -d '{"email":"new@test.com"}' -c "token=data.token"
GET /profile -H "Authorization: Bearer {{token}}" -a "status==200"
EOF

kest run user-flow.kest
```

---

### Option 2: Generate from OpenAPI/Swagger

Best for: existing API docs, quick endpoint coverage.

```bash
# From local file
kest generate --from-openapi swagger.json -o api-tests.kest

# From remote URL
kest get https://petstore3.swagger.io/api/v3/openapi.json --no-record > openapi.json
kest generate --from-openapi openapi.json -o petstore.kest
```

**Generated file example**:
```kest
# Generated from swagger.json
# Project: My API

# Get user by ID
GET /users/{id} -a "status==200"

# Create new user
POST /users -d '{}' -a "status==200"

# Update user
PUT /users/{id} -d '{}' -a "status==200"
```

**After generating, optimize the file**:
1. Replace placeholder `{}` with real data
2. Add variable captures (`-c`)
3. Add performance assertions (`--max-time`)
4. Add retry for flaky endpoints (`--retry`)

---

### Option 3: Convert from History (Recommended!)

Best for: solidifying manual test sessions into repeatable scenarios.

```bash
# 1. Test manually
kest post /login -d '{"user":"admin"}' -c "token=data.token"
kest get /profile -H "Authorization: Bearer {{token}}"
kest get /orders

# 2. Review history
kest history
# ID    TIME                 METHOD URL                    STATUS DURATION
# -------------------------------------------------------------------------
# #12   10:23:45 today       GET    /orders                200    123ms
# #11   10:23:40 today       GET    /profile               200    45ms
# #10   10:23:30 today       POST   /login                 200    234ms

# 3. Organize into a scenario
cat > my-workflow.kest << 'EOF'
# Workflow from manual testing
POST /login -d '{"user":"admin"}' -c "token=data.token"
GET /profile -H "Authorization: Bearer {{token}}"
GET /orders
EOF
```

---

### Option 4: AI-Assisted Generation

Best for: complex scenarios, rapid prototyping.

**Method A: Ask AI directly**
```
You: Generate a Kest scenario file to test an e-commerce checkout flow:
1. User login
2. Browse products
3. Add to cart
4. Place order
5. Check order status

AI: (generates .kest file)
```

**Method B: Use `kest gen`**
```bash
kest gen "test e-commerce checkout: login, browse, add to cart, order, check status"
```

---

## Scenario Templates

### Template 1: Basic CRUD

```kest
# Full CRUD test
# Create
POST /api/items -d '{"name":"test","price":100}' -c "itemId=data.id" -a "status==201"

# Read (list)
GET /api/items -a "status==200" --max-time 500

# Read (single)
GET /api/items/{{itemId}} -a "status==200" -a "body.name==test"

# Update
PUT /api/items/{{itemId}} -d '{"name":"updated","price":200}' -a "status==200"

# Delete
DELETE /api/items/{{itemId}} -a "status==204"

# Verify deletion
GET /api/items/{{itemId}} -a "status==404"
```

---

### Template 2: Authentication Flow

```kest
# Full authentication test
# 1. Register
POST /api/auth/register -d '{"email":"test@example.com","password":"pass123"}' -a "status==201"

# 2. Login
POST /api/auth/login -d '{"email":"test@example.com","password":"pass123"}' -c "accessToken=tokens.access" -c "refreshToken=tokens.refresh" -a "status==200"

# 3. Access protected resource
GET /api/protected -H "Authorization: Bearer {{accessToken}}" -a "status==200"

# 4. Refresh token
POST /api/auth/refresh -d '{"refresh_token":"{{refreshToken}}"}' -c "newAccessToken=tokens.access" -a "status==200"

# 5. Use new token
GET /api/protected -H "Authorization: Bearer {{newAccessToken}}" -a "status==200"

# 6. Logout
POST /api/auth/logout -H "Authorization: Bearer {{newAccessToken}}" -a "status==200"
```

---

### Template 3: Performance Benchmark

```kest
# Performance benchmark
# All endpoints must respond within specified time

# Health check < 100ms
GET /api/health --max-time 100 -a "status==200"

# Home page < 500ms
GET /api/home --max-time 500 -a "status==200"

# Search < 1000ms
GET /api/search?q=test --max-time 1000 -a "status==200"

# List query < 800ms
GET /api/products?page=1&limit=20 --max-time 800 -a "status==200"

# Detail page < 300ms
GET /api/products/123 --max-time 300 -a "status==200"
```

---

### Template 4: Stability Testing (Retry)

```kest
# Flaky API tests
# Webhook notification (may timeout)
POST /api/webhooks/notify -d '{"event":"order.created"}' --retry 5 --retry-delay 2000 -a "status==200"

# Third-party API (may fail)
GET /api/external/data --retry 3 --retry-delay 1000 -a "status==200"

# Eventual consistency check (needs multiple attempts)
GET /api/async/status --retry 10 --retry-delay 500 -a "body.status==completed"
```

---

### Template 5: gRPC + REST Mixed

```kest
# Mixed protocol test
# REST login
POST /api/login -d '{"email":"test@example.com"}' -c "token=data.token"

# gRPC call
grpc localhost:50051 user.UserService/GetProfile -d '{"token":"{{token}}"}' -p user.proto

# REST query
GET /api/orders -H "Authorization: Bearer {{token}}"

# gRPC create order
grpc localhost:50051 order.OrderService/Create -d '{"items":[{"id":1}]}' -p order.proto
```

---

## Running Scenarios

### Basic Execution

```bash
# Sequential (default)
kest run my-scenario.kest

# Parallel (fast)
kest run my-scenario.kest --parallel --jobs 8

# With specific environment
kest env set staging
kest run my-scenario.kest
```

### Advanced Options

```bash
# Verbose output
kest run tests.kest -v

# Pipe results to file
kest run tests.kest --parallel > test-results.log

# CI/CD mode
kest run tests.kest --quiet --output json
```

---

## Best Practices

### 1. File Organization

```
project/
├── .kest/
│   ├── config.yaml
│   └── logs/
├── scenarios/
│   ├── smoke-tests.kest      # Smoke tests
│   ├── auth-flow.kest         # Authentication flow
│   ├── user-crud.kest         # User CRUD
│   ├── order-flow.kest        # Order flow
│   └── performance.kest       # Performance benchmarks
└── README.md
```

### 2. Naming Conventions

```kest
# Good naming
# Scenario: User registration and first login
# Test: POST /register should return 201

# Avoid
# test1
# stuff
```

### 3. Comment Style

```kest
# ===================================
# Scenario: E-commerce checkout flow
# Author: stark
# Created: 2026-01-30
# Dependencies: Requires staging env
# ===================================

# Step 1: User login
# Expected: Returns access_token
POST /login -d '{"email":"test@example.com"}' -c "token=data.token"

# Step 2: Browse products (must respond < 500ms)
GET /products --max-time 500 -a "status==200"
```

### 4. Variable Management

```kest
# Use meaningful variable names
POST /login -c "accessToken=data.access" -c "userId=data.user.id"

# Avoid cryptic names
POST /login -c "t=data.access" -c "id=data.user.id"
```

### 5. Layered Assertions

```kest
# Basic assertion
GET /users -a "status==200"

# Business logic assertion
GET /users -a "status==200" -a "body.length==10"

# Performance assertion
GET /users -a "status==200" --max-time 500

# Combined assertion
GET /users -a "status==200" -a "body.length==10" --max-time 500
```

---

## Comparison with Other Formats

| Feature | Kest Scenario | Postman Collection | Hurl | k6 Script |
|---|---|---|---|---|
| Format | Plain text | JSON | Plain text | JavaScript |
| Variables | ✅ | ✅ | ✅ | ✅ |
| Assertions | ✅ | ✅ | ✅ | ✅ |
| Git-friendly | ✅ | ❌ | ✅ | ✅ |
| AI generation | ✅ | ❌ | ⚠️ | ⚠️ |
| Performance | ✅ | ❌ | ✅ | ✅ |
| gRPC | ✅ | ✅ | ❌ | ❌ |
| Parallel | ✅ | ❌ | ✅ | ✅ |

---

## Roadmap

1. **Auto-generate from history**
   ```bash
   kest history export --from 10 --to 15 -o workflow.kest
   ```

2. **Conditional execution**
   ```kest
   # if status == 200
   POST /next-step
   ```

3. **Loops**
   ```kest
   # for i in 1..10
   GET /items/{{i}}
   ```

4. **Sub-scenario imports**
   ```kest
   # import auth-flow.kest
   POST /protected-action
   ```

---

## Recommended Workflow

### Development Phase

```bash
# 1. Explore APIs manually
kest post /login -d '{}' -c "token=..."
kest get /profile -H "Authorization: ..."

# 2. Save to scenario
vim dev-tests.kest

# 3. Run and verify
kest run dev-tests.kest
```

### CI/CD Phase

```bash
# Smoke tests
kest run smoke-tests.kest --parallel --jobs 8

# Full test suite
kest run tests/ --parallel --quiet --output json
```

---

*Keep Every Step Tested.* 🦅
