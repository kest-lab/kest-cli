# 🔒 Security & Authorization Flow

Test security features including unauthorized access, invalid tokens, and permission checks.

---

## Step 1: Attempt Unauthorized Access to Profile

Try to access protected endpoint without token.

```kest
GET /api/v1/users/profile

[Asserts]
status == 401
duration < 200ms
```

---

## Step 2: Attempt Unauthorized Access to Projects

Try to access projects without authentication.

```kest
GET /api/v1/projects

[Asserts]
status == 401
duration < 200ms
```

---

## Step 3: Attempt with Invalid Token

Try to access with a malformed token.

```kest
GET /api/v1/users/profile
Authorization: Bearer invalid-token-12345

[Asserts]
status == 401
duration < 300ms
```

---

## Step 4: Register Valid User

Create a valid user for testing.

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "security{{$timestamp}}",
  "email": "security{{$timestamp}}@test.io",
  "password": "SecurePass123!"
}

[Captures]
username: data.username

[Asserts]
status == 201
body.code == 0
duration < 1000ms
```

---

## Step 5: Login to Get Valid Token

Obtain a valid authentication token.

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{username}}",
  "password": "SecurePass123!"
}

[Captures]
valid_token: data.access_token
user_id: data.user.id

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 800ms
```

---

## Step 6: Create Project

Create a project for ownership testing.

```kest
POST /api/v1/projects
Authorization: Bearer {{valid_token}}
Content-Type: application/json

{
  "name": "Security Test Project",
  "slug": "security-{{$timestamp}}"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
duration < 1000ms
```

---

## Step 7: Access Own Project (Should Succeed)

Verify user can access their own project.

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{valid_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == {{project_id}}
duration < 500ms
```

---

## Step 8: Attempt Access to Non-Existent Project

Try to access a project that doesn't exist.

```kest
GET /api/v1/projects/999999
Authorization: Bearer {{valid_token}}

[Asserts]
status == 404
duration < 500ms
```

---

## Step 9: Test Password Reset Without Auth

Verify password reset endpoint is accessible without auth.

```kest
POST /api/v1/password/reset
Content-Type: application/json

{
  "email": "{{username}}@test.io"
}

[Asserts]
status == 200
duration < 800ms
```

---

## Step 10: Verify Token Expiration Handling

Test with an expired/invalid token format.

```kest
GET /api/v1/projects
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid

[Asserts]
status == 401
duration < 300ms
```

---

**✅ Security & Authorization Flow Complete**
