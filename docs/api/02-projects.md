# Projects API

## Overview

The Projects module manages API projects, including creation, configuration, project stats, and CLI sync credentials.

## Base Path

```
/v1
```

All project endpoints require authentication.

---

## 1. Create Project

### POST /projects

Create a new API project.

**Authentication**: Required

#### Request Headers

```
Content-Type: application/json
Authorization: Bearer <token>
```

#### Request Body

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | string | ✅ Yes | min: 1, max: 100 | Project name |
| `slug` | string | ❌ No | min: 1, max: 50 | URL-friendly slug (auto-generated if not provided) |
| `platform` | string | ❌ No | enum: go, javascript, python, java, ruby, php, csharp | Primary platform/language |

#### Example Request

```json
{
  "name": "My E-commerce API",
  "slug": "ecommerce-api",
  "platform": "javascript"
}
```

#### Response (201 Created)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "My E-commerce API",
    "slug": "ecommerce-api",
    "public_key": "pk_live_51H2K3j...kLmN",
    "dsn": "https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1",
    "platform": "javascript",
    "status": 1,
    "rate_limit_per_minute": 1000,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

#### Error Responses

- **400 Bad Request**: Validation failed
- **409 Conflict**: Project slug already exists

---

## 2. List Projects

### GET /projects

List all projects for the authenticated user.

**Authentication**: Required

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | ❌ No | 1 | Page number |
| `per_page` | integer | ❌ No | 20 | Items per page (max 100) |
| `search` | string | ❌ No | - | Search by name or slug |
| `platform` | string | ❌ No | - | Filter by platform |
| `status` | integer | ❌ No | - | Filter by status (0=inactive, 1=active) |

#### Example Request

```
GET /projects?page=1&per_page=10&status=1
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "My E-commerce API",
        "slug": "ecommerce-api",
        "platform": "javascript",
        "status": 1
      },
      {
        "id": 2,
        "name": "Mobile App Backend",
        "slug": "mobile-backend",
        "platform": "go",
        "status": 1
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 2,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

---

## 3. Get Project

### GET /projects/:id

Get detailed project information.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "My E-commerce API",
    "slug": "ecommerce-api",
    "public_key": "pk_live_51H2K3j...kLmN",
    "dsn": "https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1",
    "platform": "javascript",
    "status": 1,
    "rate_limit_per_minute": 1000,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

#### Error Responses

- **404 Not Found**: Project not found
- **403 Forbidden**: No access to project

---

## 4. Update Project

### PUT /projects/:id

Update project information.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Request Body

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | string | ❌ No | min: 1, max: 100 | Project name |
| `platform` | string | ❌ No | enum: go, javascript, python, java, ruby, php, csharp | Primary platform |
| `status` | integer | ❌ No | enum: 0, 1 | Project status (0=inactive, 1=active) |
| `rate_limit_per_minute` | integer | ❌ No | min: 0, max: 100000 | Rate limit per minute |

#### Example Request

```json
{
  "name": "Updated E-commerce API",
  "platform": "typescript",
  "rate_limit_per_minute": 2000
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Project updated successfully",
  "data": {
    "id": 1,
    "name": "Updated E-commerce API",
    "slug": "ecommerce-api",
    "public_key": "pk_live_51H2K3j...kLmN",
    "dsn": "https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1",
    "platform": "typescript",
    "status": 1,
    "rate_limit_per_minute": 2000,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

---

## 5. Patch Project

### PATCH /projects/:id

Partially update project information (same as PUT but only updates provided fields).

**Authentication**: Required

Same parameters and response as PUT /projects/:id.

---

## 6. Delete Project

### DELETE /projects/:id

Delete a project and all associated data.

**⚠️ Warning**: This action is irreversible and will delete all API specs, test cases, and test results.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

---

## 7. Generate CLI Token

### POST /projects/:id/cli-tokens

Generate a project-scoped CLI token for `kest sync` uploads.

**Authentication**: Required + project write access

#### Request Body

```json
{
  "name": "Catalog API CLI sync",
  "scopes": ["spec:write"]
}
```

#### Response (201 Created)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "kest_pat_3f3b7c...",
    "token_type": "bearer",
    "project_id": 12,
    "token_info": {
      "id": 5,
      "project_id": 12,
      "name": "Catalog API CLI sync",
      "token_prefix": "kest_pat_3f3b7c12",
      "scopes": ["spec:write"],
      "created_at": "2026-04-09T10:00:00Z"
    }
  }
}
```

#### Notes

- The full token value is returned once.
- Supported scopes: `spec:write`, `run:write`
- Use the returned token with `Authorization: Bearer <kest_pat_...>`

---

## 8. Upload Specs From CLI

### POST /projects/:id/cli/spec-sync

Upload API specs inferred from local CLI history.

**Authentication**: Project-scoped CLI token with `spec:write`

#### Request Body

```json
{
  "project_id": 12,
  "source": "cli",
  "specs": [
    {
      "method": "GET",
      "path": "/v1/users",
      "title": "List users",
      "summary": "List users",
      "version": "v1"
    }
  ]
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "created": 1,
    "updated": 0,
    "skipped": 0,
    "errors": []
  }
}
```

#### Notes

- The URL project ID and token scope must match.
- Common auth headers and secret-shaped JSON fields are redacted before examples are stored.

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Project deleted successfully",
  "data": null
}
```

---

## Usage Examples

### JavaScript (Fetch API)

```javascript
const token = 'your-jwt-token';

// Create a new project
const createProject = async () => {
  const response = await fetch('http://localhost:8025/projects', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      name: 'My API Project',
      platform: 'javascript'
    })
  });
  
  const data = await response.json();
  console.log('Project created:', data.data);
  return data.data;
};

// List projects
const listProjects = async () => {
  const response = await fetch('http://localhost:8025/projects?page=1&per_page=10', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  const data = await response.json();
  console.log('Projects:', data.data.items);
  return data.data;
};

// Generate a CLI token
const createCliToken = async (projectId) => {
  const response = await fetch(`http://localhost:8025/v1/projects/${projectId}/cli-tokens`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      name: 'Payments API CLI sync',
      scopes: ['spec:write']
    })
  });
  
  const data = await response.json();
  console.log('CLI token:', data.data.token);
  return data.data;
};
```

### cURL

```bash
# Create project
curl -X POST http://localhost:8025/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "name": "My API Project",
    "platform": "go"
  }'

# List projects
curl -X GET "http://localhost:8025/projects?page=1&per_page=10" \
  -H "Authorization: Bearer TOKEN"

# Get project details
curl -X GET http://localhost:8025/projects/1 \
  -H "Authorization: Bearer TOKEN"

# Update project
curl -X PUT http://localhost:8025/projects/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "name": "Updated Project Name",
    "rate_limit_per_minute": 5000
  }'

# Generate CLI token
curl -X POST http://localhost:8025/v1/projects/1/cli-tokens \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "name": "Payments API CLI sync",
    "scopes": ["spec:write"]
  }'
```

---

## CLI Configuration

After generating a project-scoped CLI token, run this inside your Kest project:

```bash
kest sync config \
  --platform-url "https://api.kest.dev/v1" \
  --platform-token "kest_pat_..." \
  --project-id "1"
```

Then push local request history:

```bash
kest sync push
```

---

## Rate Limits

Each project has its own rate limit:

- **Default**: 1000 requests per minute
- **Maximum**: 100,000 requests per minute
- **Burst**: Up to 2x the rate limit for short bursts

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1641234567
```

---

## Security Considerations

1. **CLI Token Scope**: Generate CLI tokens per project and keep them scoped as tightly as possible.
2. **One-Time Copy**: The full CLI token is only returned once. Store it securely in `.kest/config.yaml`.
3. **HTTPS**: Always use HTTPS for production CLI uploads.
4. **Environment Separation**: Use different projects for dev, staging, and production.
5. **Access Control**: Only members with project write access should generate upload tokens.

---

## Testing

Run the project tests:

```bash
# Unit tests
go test ./internal/modules/project/...

# Integration tests
go test ./tests/feature/project_test.go
```
