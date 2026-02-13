# 🔐 User Authentication Flow

Basic user authentication testing including registration, login, and profile retrieval.

---

## Step 1: User Registration

Register a new user account with unique credentials.

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "testuser{{$timestamp}}",
  "email": "test{{$timestamp}}@kest.io",
  "password": "SecurePass123!",
  "nickname": "Test User"
}

[Captures]
registered_username: data.username
registered_email: data.email

[Asserts]
status == 201
body.code == 0
```

---

## Step 2: User Login

Login with the registered credentials to obtain access token.

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "SecurePass123!"
}

[Captures]
access_token: data.access_token
user_id: data.user.id

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 1000ms
```

---

## Step 3: Get User Profile

Retrieve the authenticated user's profile information.

```kest
GET /api/v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.username == "{{registered_username}}"
```

---

**✅ Authentication Flow Complete - 3/3 Tests Passing**
