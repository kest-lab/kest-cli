# Workspace Module API

> 💡 This documentation is automatically synchronized with the source code.

## 🎯 Module Purpose

### Module Purpose

**Business Value:**
The `workspace` module provides a comprehensive set of endpoints to manage and interact with workspaces, which are collaborative environments for teams. This module enables users to create, read, update, and delete (CRUD) workspaces, as well as manage workspace members. By offering these functionalities, the module enhances team collaboration and organizational efficiency, ensuring that users can easily organize and control access to their workspaces.

**Scope:**
The `workspace` module manages the lifecycle of workspaces and their associated members. It includes endpoints for creating, listing, retrieving, updating, and deleting workspaces. Additionally, it provides endpoints for managing workspace members, such as adding, listing, updating roles, and removing members. This ensures that all aspects of workspace management are covered, from basic operations to more granular control over member permissions.

**Role:**
This module is a core module, as it provides essential functionality for managing collaborative environments. The ability to create and manage workspaces and their members is fundamental to any application that supports team-based collaboration. The `workspace` module serves as a foundational component, enabling other modules and features to build upon its capabilities.

## 📌 Overview

The `workspace` module provides the following API endpoints:

| Method | Path | Description |
| :--- | :--- | :--- |
| <kbd>POST</kbd> | `/workspaces` | Create a new workspace for the authenticated user |
| <kbd>GET</kbd> | `/workspaces` | List all workspaces that the authenticated user is a member of or has access to if they are a super admin |
| <kbd>GET</kbd> | `/workspaces/:id` | Retrieve a specific workspace by ID, accessible to users with at least viewer role or super admins |
| <kbd>PATCH</kbd> | `/workspaces/:id` | Update a workspace by ID. Only the owner, an admin, or a super admin can update a workspace. |
| <kbd>DELETE</kbd> | `/workspaces/:id` | Delete a workspace by its ID. Only the workspace owner or a super admin can delete a workspace. |
| <kbd>POST</kbd> | `/workspaces/:id/members` | Add a new member to a workspace |
| <kbd>GET</kbd> | `/workspaces/:id/members` | List members of a specific workspace |
| <kbd>PATCH</kbd> | `/workspaces/:id/members/:uid` | Update the role of a member in a workspace |
| <kbd>DELETE</kbd> | `/workspaces/:id/members/:uid` | Remove a member from a workspace |

---

## Create a new workspace for the authenticated user

**Endpoint:**
<kbd>POST</kbd> `/workspaces`

### 🛡️ Permissions

User must be authenticated with a valid JWT token

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
participant User
participant AuthMiddleware
participant Handler
participant Service
participant Repository

User->>AuthMiddleware: Request with JWT
AuthMiddleware->>Handler: Pass request to handler
Handler->>Handler: Bind JSON to CreateWorkspaceRequest
Handler->>AuthMiddleware: Get user from context
Handler->>Service: Call CreateWorkspace with req and user ID
Service->>Service: Generate or sanitize slug
Service->>Repository: FindBySlug
alt Slug exists
  Service->>Handler: Return error
else Slug does not exist
  Service->>Service: Set default visibility if not provided
  Service->>Repository: Create workspace
  Service->>Repository: Add owner as member
  alt Adding member fails
    Repository->>Repository: Delete workspace
    Service->>Handler: Return error
  else Adding member succeeds
    Service->>Handler: Return created workspace
  end
end
Handler->>User: Return response
```

### 📥 Request: `CreateWorkspaceRequest`

| JSON Field | Type | Required/Validation | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | `required,max=100` |  |
| `slug` | `string` | `required,max=50,alphanum` |  |
| `description` | `string` | `max=500` |  |
| `type` | `string` | `required,oneof=personal team public` |  |
| `visibility` | `string` | `oneof=private team public` |  |

**Request Example:**
```json
{
    "name": "My New Workspace",
    "slug": "my-new-workspace",
    "description": "This is my new workspace for personal projects.",
    "type": "personal",
    "visibility": "private"
  }
```

**Response Example:**
```json
{
    "id": 1,
    "name": "My New Workspace",
    "slug": "my-new-workspace",
    "description": "This is my new workspace for personal projects.",
    "type": "personal",
    "owner_id": 1,
    "visibility": "private"
  }
```

**Handler Implementation:**
`workspace.CreateWorkspace`

---

## List all workspaces that the authenticated user is a member of or has access to if they are a super admin

**Endpoint:**
<kbd>GET</kbd> `/workspaces`

### 🛡️ Permissions

JWT token required. Super admins can see all workspaces, regular users can only see workspaces they are members of.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler
    participant Service
    participant Repository
    participant Database

    User->>Handler: GET /workspaces (with JWT)
    Handler->>Service: ListWorkspaces(currentUser.ID, currentUser.IsSuperAdmin)
    Service->>Repository: ListByUserID(currentUser.ID, currentUser.IsSuperAdmin)
    alt isSuperAdmin == true
        Repository->>Database: Find all workspaces (order by created_at DESC)
    else isSuperAdmin == false
        Repository->>Database: Find workspaces where user is a member (order by created_at DESC)
    end
    Database-->>Repository: return workspaces
    Repository-->>Service: return workspaces
    Service-->>Handler: return workspaces
    Handler-->>User: return workspaces

```

**Request Example:**
```json
{}
```

**Response Example:**
```json
[
    {
      "id": 1,
      "name": "Workspace 1",
      "slug": "workspace-1",
      "description": "This is the first workspace.",
      "type": "team",
      "visibility": "private",
      "created_at": "2023-10-01T12:00:00Z"
    },
    {
      "id": 2,
      "name": "Workspace 2",
      "slug": "workspace-2",
      "description": "This is the second workspace.",
      "type": "personal",
      "visibility": "public",
      "created_at": "2023-10-02T12:00:00Z"
    }
  ]
```

**Handler Implementation:**
`workspace.ListWorkspaces`

---

## Retrieve a specific workspace by ID, accessible to users with at least viewer role or super admins

**Endpoint:**
<kbd>GET</kbd> `/workspaces/:id`

### 🛡️ Permissions

⚠️ SECURITY RISK: The endpoint checks if the user has at least a viewer role or is a super admin. However, it does not explicitly check for a valid JWT token in the middleware.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
participant User
participant Middleware
participant Handler
participant Service
participant Repository
User->>Middleware: Request with workspace ID
Middleware->>Handler: Pass context with user and workspace ID
Handler->>Service: Call GetWorkspace with id, userID, isSuperAdmin
Service->>Repository: HasPermission with id, userID, RoleViewer, isSuperAdmin
alt Permission granted
Repository-->>Service: true
Service->>Repository: FindByID with id
Repository-->>Service: WorkspacePO
Service-->>Handler: FromWorkspace(WorkspacePO)
Handler-->>User: Response with workspace details
else Permission denied
Repository-->>Service: false
Service-->>Handler: Error - workspace not found or access denied
Handler-->>User: 404 Not Found
end

```

**Request Example:**
```json
{}
```

**Response Example:**
```json
{
    "id": 1,
    "name": "Example Workspace",
    "slug": "example-workspace",
    "description": "This is an example workspace.",
    "type": "team",
    "visibility": "private"
  }
```

**Handler Implementation:**
`workspace.GetWorkspace`

---

## Update a workspace by ID. Only the owner, an admin, or a super admin can update a workspace.

**Endpoint:**
<kbd>PATCH</kbd> `/workspaces/:id`

### 🛡️ Permissions

⚠️ SECURITY RISK: Only the owner, an admin, or a super admin can update a workspace. The current implementation checks for these roles but does not explicitly validate the JWT token or user session.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
participant User
participant Handler
participant Service
participant Repository
participant DB
User->>Handler: PATCH /workspaces/:id with UpdateWorkspaceRequest
Handler->>Service: UpdateWorkspace(uint(id), &req, currentUser.ID, currentUser.IsSuperAdmin)
Service->>Repository: HasPermission(id, userID, RoleAdmin, isSuperAdmin)
Repository->>DB: Check permission in database
DB-->>Repository: Return permission result
Repository-->>Service: Return permission result
alt Permission granted
Service->>Repository: FindByID(id)
Repository->>DB: Fetch workspace by ID
DB-->>Repository: Return workspace
Repository-->>Service: Return workspace
Service->>Repository: Update(workspace)
Repository->>DB: Save updated workspace
DB-->>Repository: Return save result
Repository-->>Service: Return save result
Service-->>Handler: Return updated workspace
Handler-->>User: Return updated workspace
else Permission denied
Service-->>Handler: Return error (insufficient permissions)
Handler-->>User: Return error (insufficient permissions)
end
```

### 📥 Request: `UpdateWorkspaceRequest`

| JSON Field | Type | Required/Validation | Description |
| :--- | :--- | :--- | :--- |
| `name` | `*string` | `omitempty,max=100` |  |
| `description` | `*string` | `omitempty,max=500` |  |
| `visibility` | `*string` | `omitempty,oneof=private team public` |  |

**Request Example:**
```json
{
    "name": "Updated Workspace Name",
    "description": "This is an updated description for the workspace.",
    "visibility": "public"
  }
```

**Response Example:**
```json
{
    "id": 1,
    "name": "Updated Workspace Name",
    "description": "This is an updated description for the workspace.",
    "visibility": "public"
  }
```

**Handler Implementation:**
`workspace.UpdateWorkspace`

---

## Delete a workspace by its ID. Only the workspace owner or a super admin can delete a workspace.

**Endpoint:**
<kbd>DELETE</kbd> `/workspaces/:id`

### 🛡️ Permissions

⚠️ SECURITY RISK: This endpoint requires the user to be authenticated and either the workspace owner or a super admin.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler
    participant Service
    participant Repository
    User->>Handler: DELETE /workspaces/:id
    Handler->>Service: DeleteWorkspace(id, userID, isSuperAdmin)
    Service->>Repository: FindByID(id)
    Repository-->>Service: WorkspacePO (if found) or Error
    alt Workspace found and user is authorized
        Service->>Repository: Delete(id)
        Repository-->>Service: Success or Error
        Service-->>Handler: Success or Error
        Handler-->>User: 204 No Content or 403 Forbidden
    else Workspace not found or user is not authorized
        Service-->>Handler: Error
        Handler-->>User: 403 Forbidden
    end
```

**Request Example:**
```json
{}
```

**Response Example:**
```json
{}
```

**Handler Implementation:**
`workspace.DeleteWorkspace`

---

## Add a new member to a workspace

**Endpoint:**
<kbd>POST</kbd> `/workspaces/:id/members`

### 🛡️ Permissions

⚠️ SECURITY RISK: Only users with admin, owner, or super admin roles can add members to a workspace. The user must be authenticated and have the necessary permissions.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler
    participant Service
    participant Repository
    participant DB

    User->>Handler: POST /workspaces/:id/members (AddMemberRequest)
    Handler->>Handler: Parse workspace ID
    Handler->>Handler: Bind JSON request
    Handler->>Service: AddMember(workspaceID, req, currentUser.ID, currentUser.IsSuperAdmin)
    Service->>Repository: HasPermission(workspaceID, currentUser.ID, RoleAdmin, currentUser.IsSuperAdmin)
    Repository->>DB: Check user role in workspace
    DB-->>Repository: Return user role
    Repository-->>Service: Return permission check result
    Service->>Repository: FindMember(workspaceID, req.UserID)
    Repository->>DB: Check if user is already a member
    DB-->>Repository: Return existing member
    Repository-->>Service: Return existing member check result
    Service->>Repository: AddMember(WorkspaceMemberPO)
    Repository->>DB: Create new member
    DB-->>Repository: Return creation result
    Repository-->>Service: Return creation result
    Service-->>Handler: Return success or error
    Handler-->>User: Response (Created or BadRequest)
```

### 📥 Request: `AddMemberRequest`

| JSON Field | Type | Required/Validation | Description |
| :--- | :--- | :--- | :--- |
| `user_id` | `uint` | `required` |  |
| `role` | `string` | `required,oneof=owner admin editor viewer` |  |

**Request Example:**
```json
{
    "user_id": 123,
    "role": "editor"
  }
```

**Response Example:**
```json
{
    "message": "member added successfully"
  }
```

**Handler Implementation:**
`workspace.AddMember`

---

## List members of a specific workspace

**Endpoint:**
<kbd>GET</kbd> `/workspaces/:id/members`

### 🛡️ Permissions

User must be authenticated and have at least 'viewer' role or be a super admin to access the workspace members.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler
    participant Service
    participant Repository
    participant Database

    User->>Handler: GET /workspaces/:id/members
    Handler->>Service: ListMembers(workspaceID, currentUser.ID, currentUser.IsSuperAdmin)
    Service->>Repository: HasPermission(workspaceID, userID, RoleViewer, isSuperAdmin)
    alt isSuperAdmin
        Repository-->>Service: true
    else
        Repository->>Database: Check user role in workspace
        Database-->>Repository: userRole
        Repository-->>Service: RoleLevel[userRole] >= RoleLevel[RoleViewer]
    end
    alt hasPermission
        Service->>Repository: ListMembers(workspaceID)
        Repository->>Database: Fetch members for workspace
        Database-->>Repository: members
        Repository-->>Service: members
        Service-->>Handler: members
        Handler-->>User: 200 OK with members
    else
        Service-->>Handler: 404 Not Found or Access Denied
        Handler-->>User: 404 Not Found or Access Denied
    end
```

**Request Example:**
```json
{}
```

**Response Example:**
```json
[
    {
      "user_id": 1,
      "role": "owner",
      "joined_at": "2023-10-01T12:00:00Z"
    },
    {
      "user_id": 2,
      "role": "admin",
      "joined_at": "2023-10-02T12:00:00Z"
    }
  ]
```

**Handler Implementation:**
`workspace.ListMembers`

---

## Update the role of a member in a workspace

**Endpoint:**
<kbd>PATCH</kbd> `/workspaces/:id/members/:uid`

### 🛡️ Permissions

⚠️ SECURITY RISK: Only super admins, workspace owners, and admins can update member roles. Regular users cannot change the role of the workspace owner.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant S as Service
    participant R as Repository
    C->>H: PATCH /workspaces/:id/members/:uid
    H->>H: Parse workspaceID and userID from path
    H->>H: Bind JSON to UpdateMemberRoleRequest
    H->>S: UpdateMemberRole(workspaceID, userID, req.Role, currentUser.ID, currentUser.IsSuperAdmin)
    S->>R: FindByID(workspaceID)
    R-->>S: Return workspace
    S->>S: Check if target user is the owner
    S->>S: If isSuperAdmin, update role directly
    S->>R: HasPermission(workspaceID, requesterID, RoleAdmin, false)
    R-->>S: Return permission check result
    S->>R: UpdateMemberRole(workspaceID, userID, role)
    R-->>S: Return update result
    S-->>H: Return success or error
    H-->>C: Response with message or error
```

### 📥 Request: `UpdateMemberRoleRequest`

| JSON Field | Type | Required/Validation | Description |
| :--- | :--- | :--- | :--- |
| `role` | `string` | `required,oneof=owner admin editor viewer` |  |

**Request Example:**
```json
{
    "role": "editor"
  }
```

**Response Example:**
```json
{
    "message": "role updated successfully"
  }
```

**Handler Implementation:**
`workspace.UpdateMemberRole`

---

## Remove a member from a workspace

**Endpoint:**
<kbd>DELETE</kbd> `/workspaces/:id/members/:uid`

### 🛡️ Permissions

User must be authenticated and have the 'admin' role or be a super admin to remove members. The owner cannot be removed.

### 🗺️ Logic Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler
    participant Service
    participant Repository
    participant DB

    User->>Handler: DELETE /workspaces/:id/members/:uid
    Handler->>Service: RemoveMember(workspaceID, targetUserID, currentUser.ID, currentUser.IsSuperAdmin)
    Service->>Repository: FindByID(workspaceID)
    Repository->>DB: Query workspace by ID
    DB-->>Repository: Return workspace
    Repository-->>Service: Return workspace

    alt isSuperAdmin
        Service->>Repository: RemoveMember(workspaceID, targetUserID)
        Repository->>DB: Delete member from workspace_members table
        DB-->>Repository: Return success/error
        Repository-->>Service: Return success/error
    else
        Service->>Repository: HasPermission(workspaceID, requesterID, RoleAdmin, false)
        Repository->>DB: Check permission
        DB-->>Repository: Return permission result
        Repository-->>Service: Return permission result
        alt hasPermission
            Service->>Repository: RemoveMember(workspaceID, targetUserID)
            Repository->>DB: Delete member from workspace_members table
            DB-->>Repository: Return success/error
            Repository-->>Service: Return success/error
        else
            Service-->>Handler: Forbidden (insufficient permissions)
        end
    end

    Service-->>Handler: NoContent/Forbidden
    Handler-->>User: 204 No Content/403 Forbidden

```

**Request Example:**
```json
{}
```

**Response Example:**
```json
{}
```

**Handler Implementation:**
`workspace.RemoveMember`

---

