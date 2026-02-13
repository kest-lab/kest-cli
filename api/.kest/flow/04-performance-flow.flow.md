# ⚡ Performance Benchmark Flow

Performance testing for critical endpoints with strict duration requirements.

---

## Step 1: Health Check Performance

Health endpoint must respond in under 100ms.

```kest
GET /api/v1/health

[Asserts]
status == 200
duration < 100ms
```

---

## Step 2: Setup - Login

Get authentication token for protected endpoints.

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "perf{{$timestamp}}",
  "password": "test123"
}

[Captures]
token: data.access_token

[Asserts]
status == 200
body.code == 0
duration < 800ms
```

---

## Step 3: Profile Retrieval Performance

Profile endpoint must respond in under 300ms.

```kest
GET /api/v1/users/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
duration < 300ms
```

---

## Step 4: Project List Performance

Project listing must respond in under 500ms.

```kest
GET /api/v1/projects
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
duration < 500ms
```

---

## Step 5: Create Project Performance

Project creation must complete in under 1 second.

```kest
POST /api/v1/projects
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Perf Test Project",
  "slug": "perf-{{$timestamp}}"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
duration < 1000ms
```

---

## Step 6: Get Project Performance

Single project retrieval must respond in under 300ms.

```kest
GET /api/v1/projects/{{project_id}}
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
duration < 300ms
```

---

## Step 7: Update Project Performance

Project update must complete in under 800ms.

```kest
PUT /api/v1/projects/{{project_id}}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Updated Perf Project"
}

[Asserts]
status == 200
body.code == 0
duration < 800ms
```

---

## Step 8: API Spec Creation Performance

API spec creation must complete in under 1 second.

```kest
POST /api/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "Perf Test API",
  "method": "GET",
  "path": "/test"
}

[Captures]
spec_id: data.id

[Asserts]
status == 201
body.code == 0
duration < 1000ms
```

---

## Step 9: API Spec List Performance

API spec listing must respond in under 500ms.

```kest
GET /api/v1/projects/{{project_id}}/api-specs
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
duration < 500ms
```

---

## Step 10: Delete Performance

Deletion operations must complete in under 500ms.

```kest
DELETE /api/v1/projects/{{project_id}}/api-specs/{{spec_id}}
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
duration < 500ms
```

---

**✅ Performance Benchmark Complete - All Endpoints Within SLA**
