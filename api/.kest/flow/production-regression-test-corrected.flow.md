# 🔄 Production API Regression Test (Corrected)

Test the production API at https://api.kest.dev using the correct v1 path.

---

## Step 1: Health Check

```kest
GET https://api.kest.dev/v1/health

[Asserts]
status == 200
body.status == "ok"
body.version == "v1"
duration < 1000ms
```

---

## Step 2: User Registration

```kest
POST https://api.kest.dev/v1/register
Content-Type: application/json

{
  "username": "test_user_{{$randomInt}}",
  "password": "TestPass123",
  "email": "test_{{$timestamp}}@example.com",
  "nickname": "Test User"
}

[Captures]
registered_user_id: data.id
registered_username: data.username

[Asserts]
status in [200, 201]
body.code == 0
body.data.id exists
body.data.username exists
```

---

## Step 3: User Login

```kest
POST https://api.kest.dev/v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "TestPass123"
}

[Captures]
access_token: data.access_token
user_id: data.user.id

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
body.data.user.id exists
```

---

## Step 4: Create Project

```kest
POST https://api.kest.dev/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Regression Test Project {{$timestamp}}",
  "platform": "javascript"
}

[Captures]
project_id: data.id
project_name: data.name
project_dsn: data.dsn

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name exists
body.data.dsn exists
```

---

## Step 5: List Projects

```kest
GET https://api.kest.dev/v1/projects
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.items is array
body.data.items[0].id exists
```

---

## Step 6: Get Project Details

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == {{project_id}}
body.data.name == {{project_name}}
```

---

## Step 7: Get Project DSN

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/dsn
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.dsn exists
body.data.public_key exists
```

---

## Step 8: Update Project

```kest
PUT https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Updated Project Name {{$timestamp}}",
  "rate_limit_per_minute": 2000
}

[Asserts]
status == 200
body.code == 0
body.data.name contains "Updated"
```

---

## Step 9: Get User Profile

```kest
GET https://api.kest.dev/v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == {{user_id}}
body.data.username == {{registered_username}}
```

---

## Step 10: Update User Profile

```kest
PUT https://api.kest.dev/v1/users/profile
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "nickname": "Updated Nickname {{$timestamp}}",
  "bio": "Updated bio for regression test"
}

[Asserts]
status == 200
body.code == 0
body.data.nickname contains "Updated"
```

---

## Step 11: Test API Specifications

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 12: Test Environments

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 13: Test Test Cases

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/test-cases
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 14: Test Categories

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 15: Test Members

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/members
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 16: Test Permissions

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/permissions
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 17: Test Issues

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}/issues
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 404, 501]  // May not be implemented
```

---

## Step 18: Test System Endpoints

```kest
GET https://api.kest.dev/v1/version

[Asserts]
status == 200
body.version exists
```

```kest
GET https://api.kest.dev/v1/system/info
Authorization: Bearer {{access_token}}

[Asserts]
status in [200, 401, 404, 501]  // May require admin or not implemented
```

---

## Step 19: Delete Project (Cleanup)

```kest
DELETE https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 20: Delete User Account (Cleanup)

```kest
DELETE https://api.kest.dev/v1/users/account
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Test Summary

This regression test covers:
- ✅ Authentication (register, login, profile)
- ✅ Project CRUD operations
- ✅ DSN generation
- ❓ API Specifications (may not be implemented)
- ❓ Environments (may not be implemented)
- ❓ Test Cases (may not be implemented)
- ❓ Categories (may not be implemented)
- ❓ Members (may not be implemented)
- ❓ Permissions (may not be implemented)
- ❓ Issues (may not be implemented)
- ✅ System endpoints (partial)
