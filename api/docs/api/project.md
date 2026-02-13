# Project Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `project` module provides the following API endpoints:

### {
  "description": "This endpoint allows a user to create a new project. The request must include the project name, and optionally, a slug and platform. The userID is obtained from the context, which is set by the authentication middleware. If the project creation is successful, a response with the created project details is returned. If a project with the same slug already exists, a conflict error is returned.",
  "request_example": {
    "name": "My New Project",
    "slug": "my-new-project",
    "platform": "go"
  },
  "response_example": {
    "id": 1,
    "name": "My New Project",
    "slug": "my-new-project",
    "platform": "go",
    "created_at": "2023-10-04T14:25:00Z",
    "updated_at": "2023-10-04T14:25:00Z",
    "owner_id": 1,
    "url": "http://example.com/projects/1"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects`

**Handler Implementation:**
`project.Create`

---

### {
  "description": "This endpoint retrieves a paginated list of projects. It allows clients to specify the page number and the number of items per page using query parameters. The response includes a list of project items and pagination metadata such as total number of items, current page, items per page, and total number of pages.",
  "request_example": {
    "method": "GET",
    "url": "/projects?page=1&per_page=20"
  },
  "response_example": {
    "items": [
      {
        "id": 1,
        "name": "Project A",
        "slug": "project-a",
        "platform": "go",
        "status": 1,
        "rate_limit_per_minute": 5000
      },
      {
        "id": 2,
        "name": "Project B",
        "slug": "project-b",
        "platform": "javascript",
        "status": 0,
        "rate_limit_per_minute": 10000
      }
    ],
    "meta": {
      "total": 100,
      "page": 1,
      "per_page": 20,
      "pages": 5
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects`

**Handler Implementation:**
`project.List`

---

### {
  "description": "This endpoint retrieves a specific project by its unique identifier. If the project is found, it returns the project details in a success response. If the project is not found, it returns a 404 Not Found error. In case of any other errors, it returns a 500 Internal Server Error.",
  "request": {
    "method": "GET",
    "url": "/projects/{id}",
    "parameters": {
      "path": {
        "id": "The unique identifier of the project to retrieve."
      }
    }
  },
  "response": {
    "success": {
      "status": 200,
      "body": {
        "id": "12345",
        "name": "My Project",
        "slug": "my-project",
        "platform": "go",
        "status": 1,
        "rate_limit_per_minute": 1000
      }
    },
    "not_found": {
      "status": 404,
      "body": {
        "error": "Project not found"
      }
    },
    "internal_server_error": {
      "status": 500,
      "body": {
        "error": "Internal server error"
      }
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id`

**Handler Implementation:**
`project.Get`

---

### {
  "description": "This endpoint allows you to update an existing project by providing its ID. You can modify the project's name, platform, status, and rate limit per minute. The `name` field is a required string with a minimum length of 1 and a maximum length of 100 characters. The `platform` field is optional and must be one of the following: go, javascript, python, java, ruby, php, or csharp. The `status` and `rate_limit_per_minute` fields are also optional; `status` must be either 0 or 1, and `rate_limit_per_minute` must be between 0 and 100000.",
  "request_example": {
    "name": "Updated Project Name",
    "platform": "python",
    "status": 1,
    "rate_limit_per_minute": 5000
  },
  "response_example": {
    "id": "123456789",
    "name": "Updated Project Name",
    "slug": "updated-project-name",
    "platform": "python",
    "status": 1,
    "rate_limit_per_minute": 5000,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-02T12:00:00Z"
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/projects/:id`

**Handler Implementation:**
`project.Update`

---

### {
  "description": "This endpoint allows updating an existing project by its ID. It accepts a JSON payload containing the fields to be updated, such as the project name, platform, status, and rate limit per minute. The platform must be one of the specified options (go, javascript, python, java, ruby, php, csharp). The status can be either 0 or 1, and the rate limit per minute, if provided, must be within the range of 0 to 100,000.",
  "request_example": {
    "name": "Updated Project Name",
    "platform": "python",
    "status": 1,
    "rate_limit_per_minute": 5000
  },
  "response_example": {
    "id": "12345",
    "name": "Updated Project Name",
    "slug": "updated-project-name",
    "platform": "python",
    "status": 1,
    "rate_limit_per_minute": 5000,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-02T12:00:00Z"
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/projects/:id`

**Handler Implementation:**
`project.Update`

---

### {
  "description": "This endpoint deletes a project with the specified ID. If the project does not exist, it returns a 404 Not Found error. If the deletion is successful, it returns a success message.",
  "request": {
    "method": "DELETE",
    "url": "/projects/123",
    "body": "N/A (No request body required for DELETE method)"
  },
  "response": {
    "success": {
      "status_code": 200,
      "body": {
        "message": "project deleted"
      }
    },
    "not_found": {
      "status_code": 404,
      "body": {
        "error": "Project with ID 123 not found"
      }
    },
    "internal_server_error": {
      "status_code": 500,
      "body": {
        "error": "An unexpected error occurred while deleting the project"
      }
    }
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id`

**Handler Implementation:**
`project.Delete`

---

### {
  "description": "This endpoint retrieves the DSN (Data Source Name) for a specific project identified by its ID. The DSN is used to configure error reporting or other data collection mechanisms. The response includes the DSN, public key, and project ID.",
  "request_example": null,
  "response_example": {
    "dsn": "http://public_key@host/project_id",
    "public_key": "your_public_key_here",
    "project_id": "12345"
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/dsn`

**Handler Implementation:**
`project.GetDSN`

---

