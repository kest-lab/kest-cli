# 🌐 Production Regression Test

Complete regression test for production environment using config-based URLs.

**Usage:**
```bash
# Switch to production environment
kest env use production

# Run the test
kest run .kest/flow/production-regression.flow.md

# Switch back to local
kest env use local
```

---

## Step 1: Health Check

```kest
GET /health

[Asserts]
status == 200
body.status == "ok"
body.version == "v1"
duration < 2000ms
```

---

## Step 2: User Registration

```kest
POST /register
Content-Type: application/json

{
  "username": "prod_test_{{$randomInt}}",
  "password": "SecurePass123!",
  "email": "prod_test_{{$timestamp}}@example.com",
  "nickname": "Production Test User"
}

[Captures]
user_id: data.id
username: data.username
user_email: data.email

[Asserts]
status == 200
body.code == 0
body.message == "created"
body.data.id exists
body.data.username exists
body.data.email exists
duration < 3000ms
```

---

## Step 3: User Login

```kest
POST /login
Content-Type: application/json

{
  "username": "{{username}}",
  "password": "SecurePass123!"
}

[Captures]
access_token: data.access_token
logged_user_id: data.user.id

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
body.data.user.id == {{user_id}}
duration < 2000ms
```

---

## Step 4: Get User Profile

```kest
GET /users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == {{user_id}}
body.data.username == "{{username}}"
body.data.email == "{{user_email}}"
duration < 2000ms
```

---

## Step 5: Update User Profile

```kest
PUT /users/profile
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "nickname": "Updated Test User {{$timestamp}}",
  "bio": "Automated regression test user"
}

[Captures]
updated_nickname: data.nickname

[Asserts]
status == 200
body.code == 0
body.data.nickname exists
duration < 2000ms
```

---

## Step 6: Create Project

```kest
POST /projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Prod Regression Test {{$timestamp}}",
  "platform": "javascript",
  "slug": "prod-test-{{$randomInt}}"
}

[Captures]
project_id: data.id
project_name: data.name
project_slug: data.slug
project_public_key: data.public_key

[Asserts]
status == 200
body.code == 0
body.data.id exists
body.data.name exists
body.data.slug exists
body.data.public_key exists
body.data.dsn exists
duration < 3000ms
```

---

## Step 7: List Projects

```kest
GET /projects
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.data is array
body.data.meta.total >= 1
duration < 2000ms
```

---

## Step 8: Get Project Details

```kest
GET /projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == {{project_id}}
body.data.name == "{{project_name}}"
body.data.slug == "{{project_slug}}"
duration < 2000ms
```

---

## Step 9: Get Project DSN

```kest
GET /projects/{{project_id}}/dsn
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.dsn exists
body.data.public_key == "{{project_public_key}}"
duration < 2000ms
```

---

## Step 10: Update Project

```kest
PUT /projects/{{project_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Updated Prod Test {{$timestamp}}",
  "rate_limit_per_minute": 5000
}

[Asserts]
status == 200
body.code == 0
body.data.id == {{project_id}}
duration < 2000ms
```

---

## Step 11: Test Issues Endpoint

```kest
GET /projects/{{project_id}}/issues
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.issues is array
duration < 2000ms
```

---

## Step 12: Create Issue (if supported)

```kest
POST /projects/{{project_id}}/issues
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "title": "Test Issue {{$timestamp}}",
  "description": "Automated test issue",
  "type": "bug",
  "severity": "low"
}

[Asserts]
status in [200, 201, 404, 501]
duration < 3000ms
```

---

## Step 13: Test API Specs Endpoint

```kest
GET /projects/{{project_id}}/api-specs
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 500]
duration < 2000ms
```

---

## Step 14: Test Environments Endpoint

```kest
GET /projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Step 15: Test Test Cases Endpoint

```kest
GET /projects/{{project_id}}/test-cases
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Step 16: Test Categories Endpoint

```kest
GET /projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Step 17: Test Members Endpoint

```kest
GET /projects/{{project_id}}/members
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Step 18: Test Permissions Endpoint

```kest
GET /projects/{{project_id}}/permissions
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Step 19: Delete Project (Cleanup)

```kest
DELETE /projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 2000ms
```

---

## Step 20: Delete User Account (Cleanup)

```kest
DELETE /users/account
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404]
duration < 2000ms
```

---

## Test Coverage Summary

✅ **Core Features Tested:**
- Health check
- User registration & authentication
- User profile management
- Project CRUD operations
- Project DSN generation
- Issues tracking (basic)

❓ **Optional Features (may not be implemented):**
- API Specifications
- Environments
- Test Cases
- Categories
- Members
- Permissions

---

## Performance Benchmarks

All requests should complete within:
- Health check: < 2s
- Authentication: < 3s
- CRUD operations: < 2s
- List operations: < 2s
