# Members API

## Overview

The Members module manages project members, roles, and permissions for collaborative project management.

## Base Path

```
/v1/projects/:id/members
```

All member endpoints require authentication and are scoped to a specific project.

---

## 1. List Project Members

### GET /projects/:id/members

List all members of a project.

**Authentication**: Required (Project Read access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | ❌ No | 1 | Page number |
| `per_page` | integer | ❌ No | 20 | Items per page (max 100) |
| `role` | string | ❌ No | - | Filter by role |
| `search` | string | ❌ No | - | Search by name or email |

#### Example Request

```
GET /projects/1/members?page=1&per_page=10
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
        "user_id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "nickname": "John Doe",
        "avatar": "https://example.com/avatar.jpg",
        "role": "owner",
        "status": "active",
        "invited_at": "2024-02-05T01:00:00Z",
        "joined_at": "2024-02-05T01:00:00Z",
        "last_active_at": "2024-02-05T02:00:00Z"
      },
      {
        "id": 2,
        "user_id": 2,
        "username": "jane_smith",
        "email": "jane@example.com",
        "nickname": "Jane Smith",
        "avatar": "https://example.com/avatar2.jpg",
        "role": "admin",
        "status": "active",
        "invited_at": "2024-02-05T01:30:00Z",
        "joined_at": "2024-02-05T01:35:00Z",
        "last_active_at": "2024-02-05T01:45:00Z"
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

## 2. Invite Member

### POST /projects/:id/members

Invite a new member to the project.

**Authentication**: Required (Project Admin access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Request Body

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `email` | string | ✅ Yes | valid email | Email address of user to invite |
| `role` | string | ✅ Yes | enum: admin, write, read | Member role |
| `message` | string | ❌ No | max: 500 | Custom invitation message |

#### Example Request

```json
{
  "email": "newmember@example.com",
  "role": "write",
  "message": "Please join our project for API testing"
}
```

#### Response (201 Created)

```json
{
  "code": 0,
  "message": "Invitation sent successfully",
  "data": {
    "id": 3,
    "email": "newmember@example.com",
    "role": "write",
    "status": "pending",
    "invited_at": "2024-02-05T02:00:00Z",
    "expires_at": "2024-02-12T02:00:00Z"
  }
}
```

---

## 3. Get Member

### GET /projects/:id/members/:mid

Get details of a specific project member.

**Authentication**: Required (Project Read access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |
| `mid` | integer | ✅ Yes | Member ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "nickname": "John Doe",
    "avatar": "https://example.com/avatar.jpg",
    "role": "owner",
    "status": "active",
    "permissions": [
      "project.read",
      "project.write",
      "project.admin",
      "member.manage",
      "testcase.run"
    ],
    "invited_at": "2024-02-05T01:00:00Z",
    "joined_at": "2024-02-05T01:00:00Z",
    "last_active_at": "2024-02-05T02:00:00Z",
    "activity_summary": {
      "tests_run": 45,
      "test_cases_created": 12,
      "last_login": "2024-02-05T01:55:00Z"
    }
  }
}
```

---

## 4. Update Member Role

### PATCH /projects/:id/members/:mid

Update a member's role.

**Authentication**: Required (Project Admin access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |
| `mid` | integer | ✅ Yes | Member ID |

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `role` | string | ✅ Yes | New role (admin, write, read) |

#### Example Request

```json
{
  "role": "admin"
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Member role updated successfully",
  "data": {
    "id": 2,
    "role": "admin",
    "updated_at": "2024-02-05T02:30:00Z"
  }
}
```

---

## 5. Remove Member

### DELETE /projects/:id/members/:mid

Remove a member from the project.

**Authentication**: Required (Project Admin access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |
| `mid` | integer | ✅ Yes | Member ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Member removed successfully",
  "data": null
}
```

---

## 6. Accept Invitation

### POST /projects/:id/members/accept

Accept a project invitation.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token` | string | ✅ Yes | Invitation token |

#### Example Request

```json
{
  "token": "inv_1234567890abcdef"
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Invitation accepted successfully",
  "data": {
    "member_id": 3,
    "project_id": 1,
    "role": "write",
    "joined_at": "2024-02-05T02:30:00Z"
  }
}
```

---

## 7. Decline Invitation

### POST /projects/:id/members/decline

Decline a project invitation.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `token` | string | ✅ Yes | Invitation token |
| `reason` | string | ❌ No | Decline reason |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Invitation declined",
  "data": null
}
```

---

## 8. Leave Project

### DELETE /projects/:id/members/me

Leave the current project.

**Authentication**: Required

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | ✅ Yes | Project ID |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Left project successfully",
  "data": null
}
```

---

## Roles and Permissions

### Owner

- Full control over the project
- Can delete the project
- Can manage all members
- All permissions

### Admin

- Manage project settings
- Invite/remove members (except owner)
- Manage all resources
- Permissions: project.*, member.*, testcase.*, apispec.*, environment.*

### Write

- Create and edit test cases
- View and edit API specifications
- Run tests
- Permissions: project.read, project.write, testcase.*, apispec.read, environment.read

### Read

- View-only access
- Can run existing tests
- Permissions: project.read, testcase.read, apispec.read, environment.read

---

## Usage Examples

### JavaScript (Fetch API)

```javascript
const token = 'your-jwt-token';
const projectId = 1;

// Invite member
const inviteMember = async () => {
  const response = await fetch(`http://localhost:8025/projects/${projectId}/members`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      email: 'newmember@example.com',
      role: 'write',
      message: 'Join our API testing team!'
    })
  });
  
  return await response.json();
};

// List members
const listMembers = async () => {
  const response = await fetch(`http://localhost:8025/projects/${projectId}/members`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
};

// Update member role
const updateMemberRole = async (memberId) => {
  const response = await fetch(`http://localhost:8025/projects/${projectId}/members/${memberId}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      role: 'admin'
    })
  });
  
  return await response.json();
};
```

### cURL

```bash
# Invite member
curl -X POST http://localhost:8025/projects/1/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "email": "newmember@example.com",
    "role": "write"
  }'

# List members
curl -X GET "http://localhost:8025/projects/1/members?page=1&per_page=10" \
  -H "Authorization: Bearer TOKEN"

# Update member role
curl -X PATCH http://localhost:8025/projects/1/members/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "role": "admin"
  }'

# Remove member
curl -X DELETE http://localhost:8025/projects/1/members/2 \
  -H "Authorization: Bearer TOKEN"
```

---

## Security Considerations

1. **Invitation Tokens**: Tokens expire after 7 days
2. **Role Hierarchy**: Users cannot promote themselves to higher roles
3. **Owner Protection**: Owner cannot be removed by other admins
4. **Audit Trail**: All member actions are logged
5. **Email Verification**: Invitations are sent to verified emails only
