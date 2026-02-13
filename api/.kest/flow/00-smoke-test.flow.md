# 🚀 Quick Smoke Test Flow

Fast smoke test to verify all critical endpoints are working.

---

## Step 1: Health Check

Verify the API server is running.

```kest
GET /health

[Asserts]
status == 200
body.status == "ok"
```

---

## Step 2: Register User

Quick user registration.

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "smoke{{$timestamp}}",
  "email": "smoke{{$timestamp}}@test.io",
  "password": "test123"
}

[Captures]
username: data.username

[Asserts]
status == 201
body.code == 0
```

---

## Step 3: Login

Quick login to get token.

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{username}}",
  "password": "test123"
}

[Captures]
token: data.access_token

[Asserts]
status == 200
body.code == 0
```

---

## Step 4: Get Profile

Verify authenticated endpoint works.

```kest
GET /api/v1/users/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 5: Create Project

Verify project creation works.

```kest
POST /api/v1/projects
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Smoke Test",
  "slug": "smoke-{{$timestamp}}"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
```

---

## Step 6: List Projects

Verify project listing works.

```kest
GET /api/v1/projects
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Smoke Test Complete - All Critical Endpoints Working**
