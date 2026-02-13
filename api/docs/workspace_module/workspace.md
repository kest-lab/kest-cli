# Workspace Module API

> 💡 This documentation describes the Workspace management API for team collaboration.

## 🎯 Module Purpose

The Workspace module provides multi-tenant workspace management with role-based access control. It enables teams to collaborate on projects within isolated workspaces, similar to Postman's workspace model.

**Key Features:**
- **Multi-tenant architecture**: Each workspace is isolated with its own projects and members
- **Role-based access**: Owner, Admin, Editor, and Viewer roles with different permissions
- **Super Admin support**: System-level administrators can access all workspaces
- **Flexible workspace types**: Personal, Team, and Public workspaces

## 📊 API Endpoints

### Workspace Management

#### GET /workspaces
List all workspaces accessible to the current user.

**Authentication**: Required  
**Permission**: Any authenticated user (sees only their workspaces)  
**Super Admin**: Can see all workspaces in the system

**Response Example:**
```json
[
  {
    "id": 1,
    "name": "Frontend Team",
    "slug": "frontend-team",
    "description": "Workspace for frontend development",
    "type": "team",
    "owner_id": 5,
    "visibility": "team",
    "created_at": "2026-02-07T22:00:00Z",
    "updated_at": "2026-02-07T22:00:00Z"
  }
]
```

---

#### POST /workspaces
Create a new workspace.

**Authentication**: Required  
**Request Body:**
```json
{
  "name": "Frontend Team",
  "slug": "frontend-team",
  "description": "Workspace for frontend development",
  "type": "team",
  "visibility": "team"
}
```

**Validation Rules:**
- `name`: Required, max 100 characters
- `slug`: Required, max 50 characters, alphanumeric only
- `type`: Required, one of: `personal`, `team`, `public`
- `visibility`: Optional, one of: `private`, `team`, `public`

**Response (201 Created)**:
```json
{
  "code": 0,
  "message": "created",
  "data": {
    "id": 2,
    "name": "Frontend Team",
    "slug": "frontend-team",
    "description": "Workspace for frontend development",
    "type": "team",
    "owner_id": 5,
    "visibility": "team",
    "created_at": "2026-02-07T22:15:00Z",
    "updated_at": "2026-02-07T22:15:00Z"
  }
}
```

---

#### GET /workspaces/:id
Get workspace details by ID.

**Authentication**: Required  
**Permission**: Must be a workspace member or super admin

**Path Parameters:**
- `id`: Workspace ID (integer)

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Frontend Team",
    "slug": "frontend-team",
    "description": "Workspace for frontend development",
    "type": "team",
    "owner_id": 5,
    "visibility": "team",
    "created_at": "2026-02-07T22:00:00Z",
    "updated_at": "2026-02-07T22:00:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found`: Workspace not found or access denied

---

#### PATCH /workspaces/:id
Update workspace details.

**Authentication**: Required  
**Permission**: Workspace Admin or Owner, or Super Admin

**Path Parameters:**
- `id`: Workspace ID (integer)

**Request Body:**
```json
{
  "name": "Updated Team Name",
  "description": "Updated description",
  "visibility": "public"
}
```

**All fields are optional**

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Updated Team Name",
    "slug": "frontend-team",
    "description": "Updated description",
    "type": "team",
    "owner_id": 5,
    "visibility": "public",
    "created_at": "2026-02-07T22:00:00Z",
    "updated_at": "2026-02-07T22:20:00Z"
  }
}
```

---

#### DELETE /workspaces/:id
Delete a workspace (only Owner or Super Admin).

**Authentication**: Required  
**Permission**: Workspace Owner or Super Admin only

**Path Parameters:**
- `id`: Workspace ID (integer)

**Response (204 No Content)**

**Error Responses:**
- `403 Forbidden`: Only workspace owner can delete workspace

---

### Member Management

#### GET /workspaces/:id/members
List all members of a workspace.

**Authentication**: Required  
**Permission**: Workspace member (any role) or Super Admin

**Path Parameters:**
- `id`: Workspace ID (integer)

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "workspace_id": 1,
      "user_id": 5,
      "role": "owner",
      "invited_by": 5,
      "joined_at": "2026-02-07T22:00:00Z"
    },
    {
      "id": 2,
      "workspace_id": 1,
      "user_id": 18,
      "role": "editor",
      "invited_by": 5,
      "joined_at": "2026-02-07T22:10:00Z"
    }
  ]
}
```

---

#### POST /workspaces/:id/members
Add a member to the workspace.

**Authentication**: Required  
**Permission**: Workspace Admin or Owner, or Super Admin

**Path Parameters:**
- `id`: Workspace ID (integer)

**Request Body:**
```json
{
  "user_id": 18,
  "role": "editor"
}
```

**Validation Rules:**
- `user_id`: Required (integer)
- `role`: Required, one of: `owner`, `admin`, `editor`, `viewer`

**Response (201 Created):**
```json
{
  "code": 0,
  "message": "created",
  "data": {
    "message": "member added successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: User is already a member
- `403 Forbidden`: Insufficient permissions

---

#### PATCH /workspaces/:id/members/:uid
Update a member's role.

**Authentication**: Required  
**Permission**: Workspace Admin or Owner, or Super Admin

**Path Parameters:**
- `id`: Workspace ID (integer)
- `uid`: User ID (integer)

**Request Body:**
```json
{
  "role": "admin"
}
```

**Validation Rules:**
- `role`: Required, one of: `owner`, `admin`, `editor`, `viewer`

**Restrictions:**
- Cannot change workspace owner's role
- Cannot promote to owner (workspace can only have one owner)

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "role updated successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Cannot change workspace owner's role
- `403 Forbidden`: Insufficient permissions

---

#### DELETE /workspaces/:id/members/:uid
Remove a member from the workspace.

**Authentication**: Required  
**Permission**: Workspace Admin or Owner, or Super Admin

**Path Parameters:**
- `id`: Workspace ID (integer)
- `uid`: User ID (integer)

**Restrictions:**
- Cannot remove workspace owner
- Super Admin can remove any member except owner

**Response (204 No Content)**

**Error Responses:**
- `400 Bad Request`: Cannot remove workspace owner
- `403 Forbidden`: Insufficient permissions

---

## 🎭 Permissions

### Role Hierarchy

| Role | Level | Permissions |
|------|-------|-------------|
| **Owner** | 40 | Full control, can delete workspace |
| **Admin** | 30 | Manage members, manage projects, update settings |
| **Editor** | 20 | Create/edit projects and resources |
| **Viewer** | 10 | Read-only access |

### Super Admin

Super Admins have system-level access and can:
- View all workspaces regardless of membership
- Manage any workspace's members
- Delete any workspace
- Override all permission checks

### Permission Matrix

| Action | Owner | Admin | Editor | Viewer | Super Admin |
|--------|-------|-------|--------|--------|-------------|
| Delete workspace | ✅ | ❌ | ❌ | ❌ | ✅ |
| Update workspace | ✅ | ✅ | ❌ | ❌ | ✅ |
| Invite members | ✅ | ✅ | ❌ | ❌ | ✅ |
| Remove members | ✅ | ✅ | ❌ | ❌ | ✅ |
| Change member roles | ✅ | ✅ | ❌ | ❌ | ✅ |
| Create projects | ✅ | ✅ | ✅ | ❌ | ✅ |
| Edit projects | ✅ | ✅ | ✅ | ❌ | ✅ |
| View all resources | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 🔐 Authentication

All workspace endpoints require authentication via Bearer token:

```bash
Authorization: Bearer <access_token>
```

The user information is extracted from the JWT token and used for permission checks.

---

## 📝 Workflow Examples

### Creating a Team Workspace

```bash
# 1. Create workspace
curl -X POST https://api.kest.dev/v1/workspaces \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frontend Team",
    "slug": "frontend-team",
    "type": "team",
    "visibility": "team"
  }'

# 2. Invite team members
curl -X POST https://api.kest.dev/v1/workspaces/1/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 18,
    "role": "editor"
  }'

# 3. List members
curl https://api.kest.dev/v1/workspaces/1/members \
  -H "Authorization: Bearer $TOKEN"
```

### Managing Member Roles

```bash
# Promote to admin
curl -X PATCH https://api.kest.dev/v1/workspaces/1/members/18 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# Remove member
curl -X DELETE https://api.kest.dev/v1/workspaces/1/members/18 \
  -H "Authorization: Bearer $TOKEN"
```

---

## ⚠️ Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Invalid request data | Validation failed |
| 401 | user not authenticated | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Workspace not found | Workspace doesn't exist or no access |
| 409 | Conflict | Duplicate slug or member already exists |

---

## 🗺️ Business Context

**When to use:**
- Setting up team collaboration spaces
- Organizing projects by department or team
- Managing access to multiple projects

**Related modules:**
- `project`: Projects are created within workspaces
- `member`: Workspace members can access workspace projects
- `permission`: Role-based permissions control access

---

## 📊 Data Model

```
Workspace
├── id (PK)
├── name
├── slug (unique)
├── description
├── type (personal|team|public)
├── owner_id (FK → users)
├── visibility (private|team|public)
└── settings (JSON)

WorkspaceMember
├── id (PK)
├── workspace_id (FK → workspaces)
├── user_id (FK → users)
├── role (owner|admin|editor|viewer)
├── invited_by (FK → users)
└── joined_at
```

---

**Last Updated**: 2026-02-07  
**Module Version**: 1.0.0  
**Status**: ✅ Implemented, pending integration
