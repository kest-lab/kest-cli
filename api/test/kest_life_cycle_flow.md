# Kest Full Life Cycle Test Flow

This flow verifies the complete journey from project creation to test execution, including RBAC enforcement.

## 1. Setup Phase

### Create Project
```kest
POST /api/v1/projects
X-User-ID: 100
{
  "name": "E-Commerce Suite",
  "slug": "eco-{{$randomInt}}",
  "description": "Core commerce engine"
}

[Captures]
projectID: data.id

[Asserts]
status == 201
```

### Invite Admin
```kest
POST /api/v1/projects/{{projectID}}/members
X-User-ID: 100
{
  "user_id": 101,
  "role": "admin"
}

[Asserts]
status == 201
```

### Invite Reader
```kest
POST /api/v1/projects/{{projectID}}/members
X-User-ID: 100
{
  "user_id": 102,
  "role": "read"
}

[Asserts]
status == 201
```

## 2. Resource Definition

### Create Environment
```kest
POST /api/v1/projects/{{projectID}}/environments
X-User-ID: 101
{
  "name": "Staging",
  "base_url": "http://localhost:8080",
  "variables": {
    "apiKey": "stg-secret-123"
  },
  "headers": {
    "X-User-ID": "101"
  }
}

[Asserts]
status == 201
```

### Create Category
```kest
POST /api/v1/projects/{{projectID}}/categories
X-User-ID: 101
{
  "name": "Payment Service",
  "description": "API relating to checkout"
}

[Captures]
catID: data.id

[Asserts]
status == 201
```

## 3. API Specification & Testing

### Import API Spec (Single Endpoint)
```kest
POST /api/v1/projects/{{projectID}}/api-specs
X-User-ID: 101
{
  "method": "GET",
  "path": "/api/v1/projects/{{projectID}}/categories",
  "summary": "List categories",
  "version": "v1"
}

[Captures]
specID: data.id

[Asserts]
status == 201
```

### Create Test Case from Spec
```kest
POST /api/v1/projects/{{projectID}}/test-cases/from-spec
X-User-ID: 101
{
  "api_spec_id": {{specID}},
  "name": "Smoke Test for Categories",
  "env": "Staging"
}

[Captures]
tcid: data.id

[Asserts]
status == 201
```

### Execute Test Case (Self-Running)
```kest
POST /api/v1/projects/{{projectID}}/test-cases/{{tcid}}/run
X-User-ID: 101

[Asserts]
status == 200
body.data.status == "pass"
```

## 4. RBAC Verification

### Reader cannot delete category (Expected: 403)
```kest
DELETE /api/v1/projects/{{projectID}}/categories/{{catID}}
X-User-ID: 102

[Asserts]
status == 403
```

### Admin can delete category (Expected: 204)
```kest
DELETE /api/v1/projects/{{projectID}}/categories/{{catID}}
X-User-ID: 101

[Asserts]
status == 204
```

> [!NOTE]
> All endpoints verified with project-nested structure.
