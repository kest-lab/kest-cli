# User Auth & Profile Flow Test 🧪

This flow tests the complete lifecycle of a user: registration, login, and profile retrieval.

## Step 1: User Registration
Register a new test user.

```kest
POST /v1/register
Content-Type: application/json

{
  "username": "tester{{$timestamp}}",
  "email": "test-flow-{{$timestamp}}@example.com",
  "password": "password123",
  "nickname": "FlowTester"
}

[Asserts]
status == 201
```

## Step 2: User Login
Login with the registered credentials to obtain a token.

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "tester{{$timestamp}}",
  "password": "password123"
}

[Captures]
token: data.access_token

[Asserts]
status == 200
body.data.access_token != ""
```

## Step 3: Get Profile
Use the captured token to fetch the user profile.

```kest
GET /v1/users/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.data.email != ""
```
