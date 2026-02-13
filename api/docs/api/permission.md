# Permission Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `permission` module provides the following API endpoints:

### {
  "description": "This endpoint allows the creation of a new role. It requires a JSON payload with the role's name, which is mandatory, and optionally, the display name, description, and whether it is a default role. Upon successful creation, the newly created role details are returned in the response.",
  "request_example": {
    "name": "admin",
    "display_name": "Administrator",
    "description": "Role for system administrators with full access to all features.",
    "is_default": false
  },
  "response_example": {
    "id": 1,
    "name": "admin",
    "display_name": "Administrator",
    "description": "Role for system administrators with full access to all features.",
    "is_default": false,
    "created_at": "2023-10-05T14:25:30Z",
    "updated_at": "2023-10-05T14:25:30Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/roles`

**Handler Implementation:**
`permission.CreateRole`

---

### {
  "description": "The GET /roles endpoint is used to retrieve a list of all roles available in the system. It returns an array of role objects, each containing details such as the name, display name, description, and whether it is set as a default role. No request body is required for this operation.",
  "request_example": {},
  "response_example": {
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "display_name": "Administrator",
        "description": "Has full access to all features and settings.",
        "is_default": false
      },
      {
        "id": 2,
        "name": "user",
        "display_name": "Standard User",
        "description": "Can view and edit their own profile information.",
        "is_default": true
      },
      {
        "id": 3,
        "name": "moderator",
        "display_name": "Moderator",
        "description": "Can manage user-generated content and moderate discussions.",
        "is_default": false
      }
    ]
  }
}

**Endpoint:**
<kbd>GET</kbd> `/roles`

**Handler Implementation:**
`permission.ListRoles`

---

### {
  "description": "This endpoint retrieves the details of a specific role by its ID. It first checks if the provided role ID is valid, then fetches the role from the service. If the role is found, it returns the role details in the response. If the role ID is invalid or the role is not found, it returns an appropriate error message.",
  "request": {
    "method": "GET",
    "url": "/roles/123",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response": {
    "status_code": 200,
    "body": {
      "id": 123,
      "name": "admin",
      "display_name": "Administrator",
      "description": "Full access to all resources and operations.",
      "is_default": false
    }
  },
  "error_responses": [
    {
      "status_code": 400,
      "body": {
        "message": "Invalid role ID"
      }
    },
    {
      "status_code": 404,
      "body": {
        "message": "Role not found"
      }
    }
  ]
}

**Endpoint:**
<kbd>GET</kbd> `/roles/:id`

**Handler Implementation:**
`permission.GetRole`

---

### {
  "description": "This endpoint updates an existing role by its ID. It allows modification of the role's name, display name, and description. The role ID must be a valid positive integer. The request body should contain at least one of the fields to update. If the role is successfully updated, the full updated role details are returned in the response.",
  "request_example": {
    "method": "PUT",
    "url": "/roles/123",
    "body": {
      "name": "editor",
      "display_name": "Content Editor",
      "description": "A role for users who can edit content."
    }
  },
  "response_example": {
    "id": 123,
    "name": "editor",
    "display_name": "Content Editor",
    "description": "A role for users who can edit content.",
    "is_default": false,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-05T14:30:00Z"
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/roles/:id`

**Handler Implementation:**
`permission.UpdateRole`

---

### {
  "description": "This endpoint is used to delete a role with a specified ID. It first attempts to parse the role ID from the URL parameter. If the ID is invalid, it returns a 400 Bad Request response. If the role deletion is successful, it returns a 204 No Content status. In case of an error during the deletion process, it returns a 500 Internal Server Error.",
  "request_example": {
    "method": "DELETE",
    "url": "/roles/123"
  },
  "response_example": {
    "status_code": 204,
    "body": ""
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/roles/:id`

**Handler Implementation:**
`permission.DeleteRole`

---

### {
  "description": "This endpoint assigns a specific role to a user. It requires the user ID and the role ID to be provided in the request body. Upon successful assignment, it returns a 204 No Content status, indicating that the operation was successful without any content to return.",
  "request_example": {
    "user_id": 1,
    "role_id": 2
  },
  "response_example": "No response body, as the server responds with a 204 No Content status."
}

**Endpoint:**
<kbd>POST</kbd> `/roles/assign`

**Handler Implementation:**
`permission.AssignRole`

---

### {
  "description": "This endpoint is used to remove a specific role from a user. It requires the user ID and the role ID to be provided in the request body. Upon successful removal of the role, the server responds with a 204 No Content status, indicating that the operation was successful without returning any content.",
  "request_example": {
    "user_id": 1,
    "role_id": 2
  },
  "response_example": "No response body, only a 204 No Content status code is returned upon success."
}

**Endpoint:**
<kbd>POST</kbd> `/roles/remove`

**Handler Implementation:**
`permission.RemoveRole`

---

### {
  "description": "This endpoint retrieves the roles assigned to a specific user. It requires a valid user ID as a path parameter. The response includes the user ID and a list of roles associated with that user.",
  "request_example": {
    "method": "GET",
    "url": "/users/123/roles",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response_example": {
    "status_code": 200,
    "body": {
      "UserID": 123,
      "Roles": [
        {
          "Name": "admin",
          "DisplayName": "Administrator",
          "Description": "Full access to all resources and features."
        },
        {
          "Name": "editor",
          "DisplayName": "Editor",
          "Description": "Access to edit and manage content."
        }
      ]
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/users/:id/roles`

**Handler Implementation:**
`permission.GetUserRoles`

---

### {
  "description": "The GET /permissions endpoint retrieves a list of all available permissions from the system. It returns a JSON array containing the details of each permission. No request body is required for this operation, as it simply fetches and lists the permissions.",
  "request_example": {},
  "response_example": [
    {
      "id": 1,
      "name": "view_user",
      "display_name": "View User",
      "description": "Allows the user to view other users' profiles."
    },
    {
      "id": 2,
      "name": "edit_user",
      "display_name": "Edit User",
      "description": "Allows the user to edit other users' profiles."
    },
    {
      "id": 3,
      "name": "delete_user",
      "display_name": "Delete User",
      "description": "Allows the user to delete other users' accounts."
    }
  ]
}

**Endpoint:**
<kbd>GET</kbd> `/permissions`

**Handler Implementation:**
`permission.ListPermissions`

---

