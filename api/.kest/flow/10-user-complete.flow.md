# 👤 User Module Complete Flow

Complete CRUD testing for user authentication and profile management.

---

## Step 1: User Registration

```kest
POST /v1/register
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
user_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.username exists
body.data.email exists
duration < 1000ms
```

---

## Step 2: User Login

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "SecurePass123!"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
body.data.user.username == "{{registered_username}}"
duration < 1000ms
```

---

## Step 3: Get User Profile

```kest
GET /v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.username == "{{registered_username}}"
body.data.email == "{{registered_email}}"
duration < 500ms
```

---

## Step 4: Update User Profile

```kest
PUT /v1/users/profile
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "nickname": "Updated Test User",
  "bio": "This is my updated bio"
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 5: Verify Profile Update

```kest
GET /v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.nickname == "Updated Test User"
body.data.bio == "This is my updated bio"
```

---

## Step 6: Change Password

```kest
PUT /v1/users/password
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "old_password": "SecurePass123!",
  "new_password": "NewSecurePass456!"
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 7: Login with New Password

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "NewSecurePass456!"
}

[Captures]
new_access_token: data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 1000ms
```

---

## Step 8: List Users (Admin)

```kest
GET /v1/users
Authorization: Bearer {{new_access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 9: Get User Info by ID

```kest
GET /v1/users/{{user_id}}/info
Authorization: Bearer {{new_access_token}}

[Asserts]
status == 200
body.code == 0
body.data.username == "{{registered_username}}"
duration < 500ms
```

---

## Step 10: Delete Account

```kest
DELETE /v1/users/account
Authorization: Bearer {{new_access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 11: Verify Account Deletion (Negative Test)

```kest
GET /v1/users/profile
Authorization: Bearer {{new_access_token}}

[Asserts]
status == 401
duration < 500ms
```

---

**✅ User Module Complete - 11 Steps**
