# ✅ Working Smoke Test

基于实际 API 行为的可用测试流程。

---

## Step 1: User Registration

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "smoketest{{$timestamp}}",
  "email": "smoke{{$timestamp}}@kest.io",
  "password": "SmokeTest123!",
  "nickname": "Smoke Test User"
}

[Captures]
test_username: data.username
test_email: data.email

[Asserts]
status >= 200
status < 300
body.data.username exists
duration < 1000ms
```

---

## Step 2: User Login

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{test_username}}",
  "password": "SmokeTest123!"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 1000ms
```

---

## Step 3: Get User Profile

```kest
GET /api/v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.username == "{{test_username}}"
body.data.email == "{{test_email}}"
duration < 500ms
```

---

## Step 4: Create Project

```kest
POST /api/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Smoke Test Project {{$timestamp}}",
  "description": "Automated smoke test project"
}

[Captures]
project_id: data.id

[Asserts]
status >= 200
status < 300
body.data.id exists
duration < 1000ms
```

---

## Step 5: Get Project Details

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{project_id}}"
duration < 500ms
```

---

## Step 6: List All Projects

```kest
GET /api/v1/projects
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 7: Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 8: Verify Project Deletion

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
body.code == 404
duration < 500ms
```

---

**✅ Smoke Test Complete - 8 Core Steps**
**Tests: Registration → Login → Profile → Project CRUD**
