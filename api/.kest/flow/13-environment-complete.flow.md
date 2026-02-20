# 🌍 Environment Module Complete Flow

Complete CRUD testing for environment management including duplication.

---

## Step 1: User Login & Create Project (Setup)

```kest
POST /v1/login
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
POST /v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Environment Test Project {{$timestamp}}",
  "description": "Project for testing environments"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: Create Development Environment

```kest
POST /v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Development",
  "base_url": "https://dev.example.com",
  "variables": {
    "API_KEY": "dev-key-123",
    "DEBUG": "true"
  }
}

[Captures]
env_id: data.id
env_name: data.name

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "Development"
duration < 1000ms
```

---

## Step 4: Get Environment Details

```kest
GET /v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{env_id}}"
body.data.name == "Development"
body.data.base_url == "https://dev.example.com"
duration < 500ms
```

---

## Step 5: List All Environments

```kest
GET /v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 6: Update Environment

```kest
PATCH /v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Development Updated",
  "base_url": "https://dev-v2.example.com",
  "variables": {
    "API_KEY": "dev-key-456",
    "DEBUG": "false"
  }
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 7: Verify Environment Update

```kest
GET /v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "Development Updated"
body.data.base_url == "https://dev-v2.example.com"
```

---

## Step 8: Duplicate Environment

```kest
POST /v1/projects/{{project_id}}/environments/{{env_id}}/duplicate
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Staging"
}

[Captures]
duplicated_env_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "Staging"
duration < 1000ms
```

---

## Step 9: Verify Duplicated Environment

```kest
GET /v1/projects/{{project_id}}/environments/{{duplicated_env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "Staging"
body.data.base_url == "https://dev-v2.example.com"
```

---

## Step 10: Delete Original Environment

```kest
DELETE /v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 11: Delete Duplicated Environment

```kest
DELETE /v1/projects/{{project_id}}/environments/{{duplicated_env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 12: Cleanup - Delete Project

```kest
DELETE /v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Environment Module Complete - 12 Steps**
