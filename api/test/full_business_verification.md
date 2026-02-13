# Kest Platform Business Logic Verification

This document provides a comprehensive verification of the Kest API platform, covering all core modules and RBAC enforcement.

## 1. Project Management (Owner Flow)
The owner (User 100) creates a project and manages its core metadata.

```kest
POST /api/v1/projects
X-User-ID: 100
{
  "name": "Enterprise Project",
  "slug": "ent-proj-{{$randomInt}}",
  "platform": "mobile"
}

[Captures]
projectID: data.id

[Asserts]
status == 201
body.name == "Enterprise Project"
```

## 2. RBAC & Membership
The owner (100) adds an Admin (101) and a Read-Only (102) user.

### Add Admin User
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

### Add Read-Only User
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

## 3. Category Management (Admin Flow)
User 101 (Admin) creates categories for API organization.

```kest
POST /api/v1/projects/{{projectID}}/categories
X-User-ID: 101
{
  "name": "Authentication"
}

[Captures]
categoryID: data.id

[Asserts]
status == 201
```

## 4. Environment Management (Admin Flow)
User 101 (Admin) configures a staging environment.

```kest
POST /api/v1/projects/{{projectID}}/environments
X-User-ID: 101
{
  "name": "Staging",
  "base_url": "https://staging.api.example.com",
  "project_id": {{projectID}}
}

[Captures]
envID: data.id

[Asserts]
status == 201
```

## 5. API Specification Management
Admin (101) creates an API specification.

```kest
POST /api/v1/api-specs
X-User-ID: 101
{
  "project_id": {{projectID}},
  "category_id": {{categoryID}},
  "name": "Login API",
  "method": "POST",
  "path": "/api/v1/login",
  "title": "User Login"
}

[Captures]
specID: data.id

[Asserts]
status == 201
```

## 6. Test Case Management
Owner (100) creates a test case for the specification.

```kest
POST /api/v1/test-cases
X-User-ID: 100
{
  "project_id": {{projectID}},
  "api_spec_id": {{specID}},
  "name": "Success Login Test",
  "method": "POST",
  "url": "http://localhost:8080/health",
  "asserts": [
    {
      "type": "status",
      "operator": "==",
      "expected": "200"
    }
  ]
}

[Captures]
testCaseID: data.id

[Asserts]
status == 201
```

## 7. Execution (Running Test Case)
Running the test case should return a success result.

```kest
POST /api/v1/test-cases/{{testCaseID}}/run
X-User-ID: 101

[Asserts]
status == 200
body.success == true
```

## 8. Role Enforcement (Forbidden Actions)
Read-Only user (102) should not be able to delete anything.

```kest
DELETE /api/v1/projects/{{projectID}}
X-User-ID: 102

[Asserts]
status == 403
```

## 9. Cleanup (Optional/Verification)
Ensure data integrity by listing members.

```kest
GET /api/v1/projects/{{projectID}}/members
X-User-ID: 102

[Asserts]
status == 200
body.length >= 3
```
