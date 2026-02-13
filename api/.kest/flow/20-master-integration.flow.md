# 🚀 Master Integration Flow

Complete end-to-end integration test covering all major modules.

---

## Step 1: User Registration & Login

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "master{{$timestamp}}",
  "email": "master{{$timestamp}}@kest.io",
  "password": "MasterPass123!",
  "nickname": "Master Test User"
}

[Captures]
user_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 2: Login

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{user_id}}",
  "password": "MasterPass123!"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
```

---

## Step 3: Create Project

```kest
POST /api/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Master Integration Project {{$timestamp}}",
  "description": "Complete integration test project",
  "platform": "web"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 4: Create Category

```kest
POST /api/v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Smoke Tests",
  "description": "Critical path tests",
  "sort_order": 1
}

[Captures]
category_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 5: Create Environment

```kest
POST /api/v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Production",
  "base_url": "https://api.production.com",
  "variables": {
    "API_KEY": "prod-key",
    "TIMEOUT": "5000"
  }
}

[Captures]
env_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 6: Create API Specification

```kest
POST /api/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Login API",
  "method": "POST",
  "path": "/v1/login",
  "description": "User authentication endpoint",
  "request_body": {
    "username": "string",
    "password": "string"
  },
  "response_body": {
    "access_token": "string",
    "user": "object"
  }
}

[Captures]
spec_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 7: Create Test Case from Spec

```kest
POST /api/v1/projects/{{project_id}}/test-cases/from-spec
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "spec_id": "{{spec_id}}",
  "name": "Login Test Case",
  "category_id": "{{category_id}}"
}

[Captures]
testcase_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 8: Run Test Case

```kest
POST /api/v1/projects/{{project_id}}/test-cases/{{testcase_id}}/run
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "environment_id": "{{env_id}}"
}

[Asserts]
status == 200
body.code == 0
duration < 3000ms
```

---

## Step 9: List Project Issues

```kest
GET /api/v1/projects/{{project_id}}/issues
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 10: Get Project DSN

```kest
GET /api/v1/projects/{{project_id}}/dsn
Authorization: Bearer {{access_token}}

[Captures]
project_dsn: data.dsn

[Asserts]
status == 200
body.code == 0
body.data.dsn exists
```

---

## Step 11: Update User Profile

```kest
PUT /api/v1/users/profile
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "nickname": "Master Integration User",
  "bio": "Running complete integration tests"
}

[Asserts]
status == 200
body.code == 0
```

---

## Step 12: List All Projects

```kest
GET /api/v1/projects
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
```

---

## Step 13: Export API Specs

```kest
GET /api/v1/projects/{{project_id}}/api-specs/export
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 2000ms
```

---

## Step 14: Cleanup - Delete Test Case

```kest
DELETE /api/v1/projects/{{project_id}}/test-cases/{{testcase_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 15: Cleanup - Delete API Spec

```kest
DELETE /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 16: Cleanup - Delete Environment

```kest
DELETE /api/v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 17: Cleanup - Delete Category

```kest
DELETE /api/v1/projects/{{project_id}}/categories/{{category_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 18: Cleanup - Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 19: Cleanup - Delete User Account

```kest
DELETE /api/v1/users/account
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Master Integration Flow Complete - 19 Steps**
**Covers: User, Project, Category, Environment, API Spec, Test Case, Issue modules**
