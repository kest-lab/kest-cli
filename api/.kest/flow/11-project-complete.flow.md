# 🏗️ Project Module Complete Flow

Complete CRUD testing for project management including DSN retrieval.

---

## Step 1: User Login (Setup)

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

[Captures]
access_token: data.access_token
user_id: data.user.id

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
```

---

## Step 2: Create Project

```kest
POST /api/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Kest Test Project {{$timestamp}}",
  "description": "Automated test project for Kest API testing",
  "platform": "web"
}

[Captures]
project_id: data.id
project_name: data.name

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name exists
duration < 1000ms
```

---

## Step 3: Get Project Details

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{project_id}}"
body.data.name == "{{project_name}}"
body.data.description exists
duration < 500ms
```

---

## Step 4: List All Projects

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

## Step 5: Update Project

```kest
PUT /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Updated Kest Project",
  "description": "Updated description for testing",
  "platform": "mobile"
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 6: Verify Project Update

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "Updated Kest Project"
body.data.description == "Updated description for testing"
body.data.platform == "mobile"
```

---

## Step 7: Get Project DSN

```kest
GET /api/v1/projects/{{project_id}}/dsn
Authorization: Bearer {{access_token}}

[Captures]
project_dsn: data.dsn

[Asserts]
status == 200
body.code == 0
body.data.dsn exists
duration < 500ms
```

---

## Step 8: Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 9: Verify Project Deletion (Negative Test)

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
duration < 500ms
```

---

**✅ Project Module Complete - 9 Steps**
