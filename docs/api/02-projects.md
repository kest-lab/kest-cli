# Projects API

## Overview

The Projects module manages project records for authenticated users.

- Creating a project auto-generates a slug when `slug` is omitted.
- The creator is automatically added to the project as `owner`.
- Project-specific routes are protected by project role checks.

## Base Path

```text
/v1
```

All endpoints below require authentication.

## Endpoint Summary

| Method | Path | Access |
|--------|------|--------|
| `POST` | `/projects` | Authenticated user |
| `GET` | `/projects` | Authenticated user |
| `GET` | `/projects/:id` | Project role `read` or higher |
| `PUT` | `/projects/:id` | Project role `write` or higher |
| `PATCH` | `/projects/:id` | Project role `write` or higher |
| `DELETE` | `/projects/:id` | Project role `admin` or higher |
| `GET` | `/projects/:id/stats` | Project role `read` or higher |

## 1. Create Project

### POST /projects

Create a new project.

#### Request Headers

```text
Content-Type: application/json
Authorization: Bearer <token>
```

#### Request Body

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | string | Yes | `min=1`, `max=100` | Project name |
| `slug` | string | No | `min=1`, `max=50` | Optional unique slug |
| `platform` | string | No | `go`, `javascript`, `python`, `java`, `ruby`, `php`, `csharp` | Primary platform |

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
  "message": "created",
  "data": {
    "id": 1,
    "name": "My E-commerce API",
    "slug": "ecommerce-api",
    "platform": "javascript",
    "status": 1,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

#### Error Responses

- `400 Bad Request`: Invalid request parameters
- `409 Conflict`: Project slug already exists

## 2. List Projects

### GET /projects

List projects visible to the authenticated user.

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | No | `1` | Page number |
| `per_page` | integer | No | `20` | Items per page, capped at `100` |

#### Example Request

```text
GET /projects?page=1&per_page=10
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 2,
        "name": "Mobile App Backend",
        "slug": "mobile-backend",
        "platform": "go",
        "status": 1
      },
      {
        "id": 1,
        "name": "My E-commerce API",
        "slug": "ecommerce-api",
        "platform": "javascript",
        "status": 1
      }
    ],
    "meta": {
      "total": 2,
      "page": 1,
      "per_page": 10,
      "pages": 1
    }
  }
}
```

#### Notes

- Results are restricted to projects where the authenticated user is a member.
- The current implementation does not support `search`, `platform`, or `status` filters on this endpoint.

## 3. Get Project

### GET /projects/:id

Get project details.

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Project ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "My E-commerce API",
    "slug": "ecommerce-api",
    "platform": "javascript",
    "status": 1,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

#### Error Responses

- `400 Bad Request`: Invalid ID format
- `403 Forbidden`: Permission denied
- `404 Not Found`: Project not found

## 4. Update Project

### PUT /projects/:id

Replace the editable fields of a project.

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Project ID |

#### Request Body

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | string | No | `min=1`, `max=100` | Project name |
| `platform` | string | No | `go`, `javascript`, `python`, `java`, `ruby`, `php`, `csharp` | Primary platform |
| `status` | integer | No | `0`, `1` | Project status |

#### Example Request

```json
{
  "name": "Updated E-commerce API",
  "platform": "python",
  "status": 1
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Updated E-commerce API",
    "slug": "ecommerce-api",
    "platform": "python",
    "status": 1,
    "created_at": "2024-02-05T01:00:00Z"
  }
}
```

#### Error Responses

- `400 Bad Request`: Invalid ID format or request body
- `403 Forbidden`: Permission denied
- `404 Not Found`: Project not found

## 5. Patch Project

### PATCH /projects/:id

Partially update a project.

The accepted request body and response shape are the same as `PUT /projects/:id`.

## 6. Delete Project

### DELETE /projects/:id

Delete a project.

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Project ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "project deleted"
  }
}
```

#### Error Responses

- `400 Bad Request`: Invalid ID format
- `403 Forbidden`: Permission denied
- `404 Not Found`: Project not found

## 7. Get Project Stats

### GET /projects/:id/stats

Return aggregate counters for a project.

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Project ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "api_spec_count": 12,
    "flow_count": 3,
    "environment_count": 2,
    "member_count": 5,
    "category_count": 8
  }
}
```

#### Error Responses

- `400 Bad Request`: Invalid ID format
- `403 Forbidden`: Permission denied
- `500 Internal Server Error`: Failed to load project statistics

## Usage Examples

### JavaScript (Fetch API)

```javascript
const token = 'your-jwt-token';

async function createProject() {
  const response = await fetch('http://localhost:8025/v1/projects', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      name: 'My API Project',
      platform: 'javascript',
    }),
  });

  return response.json();
}

async function listProjects() {
  const response = await fetch('http://localhost:8025/v1/projects?page=1&per_page=10', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.json();
}

async function getProjectStats(projectId) {
  const response = await fetch(`http://localhost:8025/v1/projects/${projectId}/stats`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return response.json();
}
```

### cURL

```bash
# Create project
curl -X POST http://localhost:8025/v1/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "name": "My API Project",
    "platform": "go"
  }'

# List projects
curl -X GET "http://localhost:8025/v1/projects?page=1&per_page=10" \
  -H "Authorization: Bearer TOKEN"

# Get project details
curl -X GET http://localhost:8025/v1/projects/1 \
  -H "Authorization: Bearer TOKEN"

# Patch project
curl -X PATCH http://localhost:8025/v1/projects/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "name": "Updated Project Name",
    "status": 1
  }'

# Get project stats
curl -X GET http://localhost:8025/v1/projects/1/stats \
  -H "Authorization: Bearer TOKEN"
```

## Current Implementation Notes

- Project responses currently expose `id`, `name`, `slug`, `platform`, `status`, and `created_at`.
- This module does not currently expose `public_key`, `dsn`, or `rate_limit_per_minute`.
- The documented contract here is based on the current handlers, DTOs, and routes in `internal/modules/project`.
