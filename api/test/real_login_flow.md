# Kest Real Login & Auth Flow

This flow verifies the integration of production-grade Registration, Login, and JWT authentication.

## 1. Authentication Phase

### Register New User
```kest
POST /api/v1/register
{
  "username": "tester{{$randomInt}}",
  "password": "password123",
  "email": "tester{{$randomInt}}@kest.io",
  "nickname": "Test User",
  "phone": "13800000000"
}

[Captures]
capturedUsername: data.username

[Asserts]
status == 201
body.data.username != ""
```

### Login to Get Token
```kest
POST /api/v1/login
{
  "username": "{{capturedUsername}}",
  "password": "password123"
}

[Captures]
authToken: data.access_token

[Asserts]
status == 200
body.data.access_token != ""
```

### Verify Token (Get Profile)
```kest
GET /api/v1/users/profile
Authorization: Bearer {{authToken}}

[Asserts]
status == 200
```

## 2. Resource Access with JWT

### Create Project (Authenticated)
```kest
POST /api/v1/projects
Authorization: Bearer {{authToken}}
{
  "name": "JWT Protected Project",
  "slug": "jwt-{{$randomInt}}",
  "description": "Verified via real JWT"
}

[Captures]
projectID: data.id

[Asserts]
status == 201
```

### Verify Unauthorized Access (No Token)
```kest
GET /api/v1/projects/{{projectID}}

[Asserts]
status == 401
```

### Verify Unauthorized Access (Invalid Token)
```kest
GET /api/v1/projects/{{projectID}}
Authorization: Bearer invalid-token-here

[Asserts]
status == 401
```
