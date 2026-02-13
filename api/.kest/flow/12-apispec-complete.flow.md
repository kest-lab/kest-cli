# 📋 API Spec Module Complete Flow

Complete CRUD testing for API specifications and examples.

---

## Step 1: User Login & Create Project (Setup)

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
```

---

## Step 2: Create Test Project

```kest
POST /api/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "API Spec Test Project {{$timestamp}}",
  "description": "Project for testing API specifications"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: Create API Specification

```kest
POST /api/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "User API",
  "method": "POST",
  "path": "/api/users",
  "description": "Create a new user",
  "request_body": {
    "username": "string",
    "email": "string"
  },
  "response_body": {
    "id": "string",
    "username": "string",
    "email": "string"
  }
}

[Captures]
spec_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "User API"
duration < 1000ms
```

---

## Step 4: Get API Specification

```kest
GET /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{spec_id}}"
body.data.name == "User API"
body.data.method == "POST"
duration < 500ms
```

---

## Step 5: List API Specifications

```kest
GET /api/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 6: Update API Specification

```kest
PATCH /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Updated User API",
  "description": "Updated description for user creation API"
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 7: Get Spec with Examples

```kest
GET /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}/full
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{spec_id}}"
body.data.name == "Updated User API"
duration < 500ms
```

---

## Step 8: Create API Example

```kest
POST /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}/examples
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Success Example",
  "request": {
    "username": "john_doe",
    "email": "john@example.com"
  },
  "response": {
    "id": "123",
    "username": "john_doe",
    "email": "john@example.com"
  },
  "status_code": 201
}

[Captures]
example_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
duration < 1000ms
```

---

## Step 9: Export API Specs

```kest
GET /api/v1/projects/{{project_id}}/api-specs/export
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 2000ms
```

---

## Step 10: Delete API Specification

```kest
DELETE /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 11: Verify Deletion (Negative Test)

```kest
GET /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
duration < 500ms
```

---

## Step 12: Cleanup - Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ API Spec Module Complete - 12 Steps**
