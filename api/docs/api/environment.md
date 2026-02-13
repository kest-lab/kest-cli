# Environment Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `environment` module provides the following API endpoints:

### {
  "description": "The GET / endpoint in the environment module is used to retrieve a list of all environments associated with a project. It does not require any request body, and the response will include details such as the name, display name, base URL, variables, and headers for each environment.",
  "request_example": {},
  "response_example": {
    "environments": [
      {
        "id": 1,
        "project_id": 101,
        "name": "production",
        "display_name": "Production Environment",
        "base_url": "https://api.production.com",
        "variables": {
          "API_KEY": "prod-12345",
          "TIMEOUT": 30
        },
        "headers": {
          "Authorization": "Bearer prod-token-123"
        }
      },
      {
        "id": 2,
        "project_id": 101,
        "name": "staging",
        "display_name": "Staging Environment",
        "base_url": "https://api.staging.com",
        "variables": {
          "API_KEY": "stage-67890",
          "TIMEOUT": 20
        },
        "headers": {
          "Authorization": "Bearer stage-token-456"
        }
      }
    ]
  }
}

**Endpoint:**
<kbd>GET</kbd> `/`

**Handler Implementation:**
`environment.unknown`

---

### {
  "description": "This endpoint is used to create a new environment within a specified project. The request should include the project ID, a unique name for the environment, and optionally, a display name, base URL, variables, and headers. This allows for setting up different configurations (like development, staging, or production) with specific settings for each.",
  "request_example": {
    "project_id": 1,
    "name": "production",
    "display_name": "Production Environment",
    "base_url": "https://api.production.example.com",
    "variables": {
      "API_KEY": "prod-key-12345",
      "TIMEOUT": 30
    },
    "headers": {
      "Authorization": "Bearer prod-token-12345"
    }
  },
  "response_example": {
    "id": 1,
    "project_id": 1,
    "name": "production",
    "display_name": "Production Environment",
    "base_url": "https://api.production.example.com",
    "variables": {
      "API_KEY": "prod-key-12345",
      "TIMEOUT": 30
    },
    "headers": {
      "Authorization": "Bearer prod-token-12345"
    },
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/`

**Handler Implementation:**
`environment.unknown`

---

### ```json
{
  "description": "This endpoint retrieves the details of a specific environment by its ID. The response includes all the information associated with the environment, such as its name, display name, base URL, variables, and headers.",
  "request": {
    "method": "GET",
    "url": "/environment/123",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response": {
    "status": 200,
    "body": {
      "id": 123,
      "project_id": 456,
      "name": "production",
      "display_name": "Production Environment",
      "base_url": "https://api.example.com",
      "variables": {
        "API_KEY": "abc123",
        "TIMEOUT": 30
      },
      "headers": {
        "Authorization": "Bearer xyz789",
        "Content-Type": "application/json"
      }
    }
  }
}
```

**Endpoint:**
<kbd>GET</kbd> `/:eid`

**Handler Implementation:**
`environment.unknown`

---

### {
  "description": "This endpoint allows for updating an existing environment by its ID. You can modify the environment's name, display name, base URL, variables, and headers. Only the fields that need to be updated should be included in the request. If a field is not provided, it will retain its current value.",
  "request_example": {
    "name": "UpdatedEnvironmentName",
    "display_name": "Updated Display Name",
    "base_url": "https://updated-base-url.com",
    "variables": {
      "new_variable": "new_value"
    },
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response_example": {
    "status": "success",
    "message": "Environment updated successfully",
    "data": {
      "id": 1,
      "project_id": 123,
      "name": "UpdatedEnvironmentName",
      "display_name": "Updated Display Name",
      "base_url": "https://updated-base-url.com",
      "variables": {
        "existing_variable": "existing_value",
        "new_variable": "new_value"
      },
      "headers": {
        "Authorization": "Bearer token",
        "Content-Type": "application/json"
      }
    }
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/:eid`

**Handler Implementation:**
`environment.unknown`

---

### ```json
{
  "description": "This endpoint allows the deletion of a specific environment identified by its unique ID (eid). Upon successful execution, the specified environment will be removed from the system. No request body is needed for this operation as it is a DELETE request, and the response will typically be a status code indicating the success or failure of the operation.",
  "request_example": {
    "method": "DELETE",
    "url": "/environment/12345",
    "headers": {
      "Authorization": "Bearer your_access_token"
    }
  },
  "response_example": {
    "status_code": 204,
    "body": null
  }
}
```

**Endpoint:**
<kbd>DELETE</kbd> `/:eid`

**Handler Implementation:**
`environment.unknown`

---

### {
  "description": "This endpoint allows the duplication of an existing environment within a project. The new environment will be created with the same settings as the original, except for the name and any variables that are overridden. The `:eid` in the URL represents the ID of the environment to be duplicated.",
  "request_example": {
    "method": "POST",
    "url": "/123/duplicate",
    "body": {
      "name": "NewEnvironmentName",
      "override_vars": {
        "variable1": "new_value1",
        "variable2": "new_value2"
      }
    }
  },
  "response_example": {
    "status": 201,
    "body": {
      "id": 456,
      "project_id": 789,
      "name": "NewEnvironmentName",
      "display_name": "Original Environment Display Name",
      "base_url": "https://original-base-url.com",
      "variables": {
        "variable1": "new_value1",
        "variable2": "new_value2",
        "variable3": "original_value3"
      },
      "headers": {
        "header1": "value1",
        "header2": "value2"
      }
    }
  }
}

**Endpoint:**
<kbd>POST</kbd> `/:eid/duplicate`

**Handler Implementation:**
`environment.unknown`

---

