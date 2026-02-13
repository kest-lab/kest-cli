# Projects API

## Overview

The Projects module manages API projects, including creation, configuration, and DSN (Data Source Name) generation.

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

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Project deleted successfully",
  "data": null
}
```

---

## 7. Get Project DSN

### GET /projects/:id/dsn

Get the Data Source Name (DSN) for a project. This is used to configure SDKs and send data to Kest.

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
    "dsn": "https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1",
    "public_key": "pk_live_51H2K3j...kLmN",
    "project_id": 1,
    "environment": "production"
  }
}
```

#### DSN Format

The DSN contains all necessary information to connect your application to Kest:

```
https://api.kest.com/v1/ingest?public_key=<public_key>&project_id=<project_id>
```

For development environment:

```
https://api-dev.kest.com/v1/ingest?public_key=<public_key>&project_id=<project_id>
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

// Get project DSN
const getProjectDSN = async (projectId) => {
  const response = await fetch(`http://localhost:8025/projects/${projectId}/dsn`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  const data = await response.json();
  console.log('DSN:', data.data.dsn);
  return data.data.dsn;
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

# Get DSN
curl -X GET http://localhost:8025/projects/1/dsn \
  -H "Authorization: Bearer TOKEN"
```

---

## SDK Configuration

Once you have your project's DSN, you can configure the Kest SDK in your application:

### JavaScript/TypeScript

```javascript
import Kest from '@kest-lab/kest-js';

const kest = new Kest({
  dsn: 'https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1'
});
```

### Go

```go
import "github.com/kest-lab/kest-go"

kest.Init(kest.Config{
    DSN: "https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1",
})
```

### Python

```python
from kest import Kest

kest = Kest(
    dsn="https://api.kest.com/v1/ingest?public_key=pk_live_51H2K3j...kLmN&project_id=1"
)
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

1. **Public Key**: Treat your project's public key like a password
2. **HTTPS**: Always use HTTPS for production DSNs
3. **Environment Separation**: Use different projects for different environments
4. **Access Control**: Only authorized users should be able to view/update projects

---

## Testing

Run the project tests:

```bash
# Unit tests
go test ./internal/modules/project/...

# Integration tests
go test ./tests/feature/project_test.go
```
