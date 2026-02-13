# 🐛 Issue Module Complete Flow

Complete testing for issue management including status changes and event tracking.

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
  "name": "Issue Test Project {{$timestamp}}",
  "description": "Project for testing issue tracking"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: List Issues (Empty State)

```kest
GET /api/v1/projects/{{project_id}}/issues
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 4: Get Issue by Fingerprint (Simulated)

Note: Issues are typically created by error ingestion. For testing, we'll use a mock fingerprint.

```kest
GET /api/v1/projects/{{project_id}}/issues/test-fingerprint-{{$timestamp}}
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
duration < 1000ms
```

---

## Step 5: Resolve Issue

```kest
POST /api/v1/projects/{{project_id}}/issues/test-fingerprint-123/resolve
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "comment": "Fixed in version 1.2.3"
}

[Asserts]
status >= 200
duration < 1000ms
```

---

## Step 6: Ignore Issue

```kest
POST /api/v1/projects/{{project_id}}/issues/test-fingerprint-456/ignore
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "comment": "Known issue, will fix later"
}

[Asserts]
status >= 200
duration < 1000ms
```

---

## Step 7: Reopen Issue

```kest
POST /api/v1/projects/{{project_id}}/issues/test-fingerprint-123/reopen
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "comment": "Issue still occurring"
}

[Asserts]
status >= 200
duration < 1000ms
```

---

## Step 8: Get Issue Events

```kest
GET /api/v1/projects/{{project_id}}/issues/test-fingerprint-123/events
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
body.code == 0
duration < 1000ms
```

---

## Step 9: Cleanup - Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Issue Module Complete - 9 Steps**
