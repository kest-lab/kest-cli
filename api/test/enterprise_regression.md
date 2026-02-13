# Enterprise Regression Test Suite

This suite verifies high-level enterprise requirements: Multi-tenancy isolation, Role hierarchy, and complex resource management.

## 1. Multi-Project Setup
We create two distinct projects to test isolation.

### Project Alpha (Team A)
```kest
POST /api/v1/projects
X-User-ID: 100
{
  "name": "Project Alpha",
  "slug": "alpha-{{$randomInt}}",
  "description": "Internal Core API"
}

[Captures]
alphaID: data.id

[Asserts]
status == 201
```

### Project Beta (Team B)
```kest
POST /api/v1/projects
X-User-ID: 200
{
  "name": "Project Beta",
  "slug": "beta-{{$randomInt}}",
  "description": "Experimental Sandbox"
}

[Captures]
betaID: data.id

[Asserts]
status == 201
```

## 2. Team Collaboration (Project Alpha)
Owner (100) invites Admin (101) and Reader (102).

### Add Admin
```kest
POST /api/v1/projects/{{alphaID}}/members
X-User-ID: 100
{
  "user_id": 101,
  "role": "admin"
}

[Asserts]
status == 201
```

### Add Reader
```kest
POST /api/v1/projects/{{alphaID}}/members
X-User-ID: 100
{
  "user_id": 102,
  "role": "read"
}

[Asserts]
status == 201
```

## 3. Strict Isolation Verification (Tenant Check)
User 200 (Owner of Beta) should NOT be able to see Project Alpha's members.

```kest
GET /api/v1/projects/{{alphaID}}/members
X-User-ID: 200

[Asserts]
status == 403
```

## 4. Role Hierarchy Verification (Permission Check)

### Admin can create resources (Expected: Pass)
```kest
POST /api/v1/projects/{{alphaID}}/categories
X-User-ID: 101
{
  "name": "Admin Created Category"
}

[Captures]
catID: data.id

[Asserts]
status == 201
```

### Reader CANNOT create resources (Expected: Forbidden)
```kest
POST /api/v1/projects/{{alphaID}}/categories
X-User-ID: 102
{
  "name": "Illegal Category"
}

[Asserts]
status == 403
```

### Reader CANNOT add members (Expected: Forbidden)
```kest
POST /api/v1/projects/{{alphaID}}/members
X-User-ID: 102
{
  "user_id": 103,
  "role": "read"
}

[Asserts]
status == 403
```

## 5. Cross-Project Resource Access (Negative)
User 101 (Admin of Alpha) tries to create a category in Project Beta.

```kest
POST /api/v1/projects/{{betaID}}/categories
X-User-ID: 101
{
  "name": "Stealth Category"
}

[Asserts]
status == 403
```

## 6. Cleanup & Final State
Owner (100) can delete the category created by Admin.

```kest
DELETE /api/v1/projects/{{alphaID}}/categories/{{catID}}
X-User-ID: 100

[Asserts]
status == 204
```
