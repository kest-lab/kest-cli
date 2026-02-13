# Apispec Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `apispec` module provides the following API endpoints:

### {
  "summary": "Retrieves a paginated list of API specifications for a given project, optionally filtered by version.",
  "request_example": "",
  "response_example": {
    "items": [
      {
        "project_id": 1,
        "category_id": 2,
        "method": "GET",
        "path": "/example/path",
        "summary": "Example summary",
        "description": "This is an example description",
        "tags": ["tag1", "tag2"],
        "request_body": {},
        "parameters": [],
        "responses": {
          "200": {
            "description": "OK"
          }
        },
        "version": "1.0.0",
        "is_public": true
      }
    ],
    "meta": {
      "total": 100,
      "page": 1,
      "page_size": 20,
      "total_pages": 5
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/api-specs`

**Handler Implementation:**
`apispec.ListSpecs`

---

### {
  "summary": "Creates a new API specification for a given project. The endpoint allows specifying the HTTP method, path, summary, description, tags, request body, parameters, and responses for the API.",
  "request_example": {
    "method": "GET",
    "path": "/users",
    "summary": "Retrieve a list of users",
    "description": "This endpoint returns a list of all users in the system.",
    "tags": ["user", "list"],
    "request_body": {
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "limit": {
                "type": "integer"
              },
              "offset": {
                "type": "integer"
              }
            }
          }
        }
      }
    },
    "parameters": [
      {
        "name": "limit",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer"
        }
      },
      {
        "name": "offset",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer"
        }
      }
    ],
    "responses": {
      "200": {
        "description": "A list of users",
        "content": {
          "application/json": {
            "schema": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer"
                  },
                  "name": {
                    "type": "string"
                  }
                }
              }
            }
          }
        }
      }
    },
    "version": "1.0.0",
    "is_public": true
  },
  "response_example": {
    "id": 1,
    "project_id": 1,
    "category_id": null,
    "method": "GET",
    "path": "/users",
    "summary": "Retrieve a list of users",
    "description": "This endpoint returns a list of all users in the system.",
    "tags": ["user", "list"],
    "request_body": {
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "limit": {
                "type": "integer"
              },
              "offset": {
                "type": "integer"
              }
            }
          }
        }
      }
    },
    "parameters": [
      {
        "name": "limit",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer"
        }
      },
      {
        "name": "offset",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer"
        }
      }
    ],
    "responses": {
      "200": {
        "description": "A list of users",
        "content": {
          "application/json": {
            "schema": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "id": {
                    "type": "integer"
                  },
                  "name": {
                    "type": "string"
                  }
                }
              }
            }
          }
        }
      }
    },
    "version": "1.0.0",
    "is_public": true
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/api-specs`

**Handler Implementation:**
`apispec.CreateSpec`

---

### {
  "summary": "This endpoint allows the import of multiple API specifications into a specified project. Each specification includes details such as HTTP method, path, summary, and more.",
  "request_example": {
    "specs": [
      {
        "project_id": 1,
        "category_id": 2,
        "method": "GET",
        "path": "/users",
        "summary": "Retrieve a list of users",
        "description": "This API returns a list of all registered users in the system.",
        "tags": ["user", "list"],
        "request_body": null,
        "parameters": [],
        "responses": {
          "200": {
            "description": "A successful response with a list of users"
          }
        },
        "version": "1.0.0",
        "is_public": true
      },
      {
        "project_id": 1,
        "method": "POST",
        "path": "/users",
        "summary": "Create a new user",
        "description": "This API creates a new user in the system.",
        "tags": ["user", "create"],
        "request_body": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "name": {
                    "type": "string"
                  },
                  "email": {
                    "type": "string"
                  }
                }
              }
            }
          }
        },
        "parameters": [],
        "responses": {
          "201": {
            "description": "User created successfully"
          }
        },
        "version": "1.0.0",
        "is_public": false
      }
    ]
  },
  "response_example": {
    "message": "Specs imported successfully"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/api-specs/import`

**Handler Implementation:**
`apispec.ImportSpecs`

---

### {
  "summary": "Exports the API specifications of a project in a specified format, such as JSON. The endpoint allows users to download or view the API specs for further use or documentation.",
  "request_example": "",
  "response_example": {
    "api_spec_id": 1,
    "name": "Example API",
    "method": "GET",
    "path": "/example",
    "summary": "A brief summary of the example API.",
    "description": "A detailed description of the example API.",
    "tags": ["example", "test"],
    "request_body": {},
    "parameters": [],
    "responses": {
      "200": {
        "description": "Success"
      }
    },
    "version": "1.0.0",
    "is_public": false
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/api-specs/export`

**Handler Implementation:**
`apispec.ExportSpecs`

---

### {
  "summary": "Retrieves the details of a specific API specification by its ID for a given project.",
  "request_example": "",
  "response_example": {
    "project_id": 1,
    "category_id": 2,
    "method": "GET",
    "path": "/users",
    "summary": "Get a list of users",
    "description": "This endpoint returns a list of all users in the system.",
    "tags": ["user", "list"],
    "request_body": {
      "content_type": "application/json",
      "example": "{\"name\": \"John\"}"
    },
    "parameters": [
      {
        "name": "page",
        "in": "query",
        "description": "Page number",
        "required": true,
        "schema": {
          "type": "integer"
        }
      }
    ],
    "responses": {
      "200": {
        "description": "Successful response",
        "content_type": "application/json",
        "example": "[{\"id\": 1, \"name\": \"John\"}, {\"id\": 2, \"name\": \"Jane\"}]"
      },
      "400": {
        "description": "Bad request",
        "content_type": "application/json",
        "example": "{\"error\": \"Invalid page number\"}"
      }
    },
    "version": "1.0.0",
    "is_public": true
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/api-specs/:sid`

**Handler Implementation:**
`apispec.GetSpec`

---

### {
  "summary": "Retrieves a full API specification including examples for a given project and spec ID.",
  "request_example": "",
  "response_example": {
    "category_id": 1,
    "summary": "A sample API endpoint",
    "description": "This is a detailed description of the API endpoint, explaining its purpose and usage.",
    "tags": ["tag1", "tag2"],
    "request_body": {
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "id": {
                "type": "integer"
              },
              "name": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "parameters": [
      {
        "name": "param1",
        "in": "query",
        "required": true,
        "schema": {
          "type": "string"
        }
      }
    ],
    "responses": {
      "200": {
        "description": "Success",
        "content": {
          "application/json": {
            "schema": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "integer"
                },
                "name": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "is_public": true,
    "examples": [
      {
        "name": "Example 1",
        "request_headers": {
          "Content-Type": "application/json"
        },
        "request_body": "{\"id\": 1, \"name\": \"John Doe\"}",
        "response_status": 200,
        "response_body": "{\"id\": 1, \"name\": \"John Doe\", \"message\": \"Welcome!\"}",
        "duration_ms": 150
      }
    ]
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/api-specs/:sid/full`

**Handler Implementation:**
`apispec.GetSpecWithExamples`

---

### {
  "summary": "Update an existing API specification for a project. This endpoint allows modification of various fields such as summary, description, tags, request body, parameters, responses, and visibility.",
  "request_example": {
    "category_id": 2,
    "summary": "Updated summary of the API",
    "description": "This is an updated detailed description of what the API does.",
    "tags": ["tag1", "tag2"],
    "request_body": {
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "age": {
                "type": "integer"
              }
            }
          }
        }
      }
    },
    "parameters": [
      {
        "name": "Authorization",
        "in": "header",
        "required": true
      }
    ],
    "responses": {
      "200": {
        "description": "Success response",
        "content": {
          "application/json": {
            "example": {
              "message": "Success"
            }
          }
        }
      }
    },
    "is_public": true
  },
  "response_example": {
    "id": 1,
    "project_id": 1,
    "category_id": 2,
    "method": "GET",
    "path": "/users",
    "summary": "Updated summary of the API",
    "description": "This is an updated detailed description of what the API does.",
    "tags": ["tag1", "tag2"],
    "request_body": {
      "content": {
        "application/json": {
          "schema": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "age": {
                "type": "integer"
              }
            }
          }
        }
      }
    },
    "parameters": [
      {
        "name": "Authorization",
        "in": "header",
        "required": true
      }
    ],
    "responses": {
      "200": {
        "description": "Success response",
        "content": {
          "application/json": {
            "example": {
              "message": "Success"
            }
          }
        }
      }
    },
    "version": "1.0.0",
    "is_public": true
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/projects/:id/api-specs/:sid`

**Handler Implementation:**
`apispec.UpdateSpec`

---

### Deletes a specific API specification from a project.

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id/api-specs/:sid`

**Handler Implementation:**
`apispec.DeleteSpec`

---

### {
  "summary": "Creates a new example for a specific API specification within a project. This endpoint allows users to add a realistic example of how the API can be used, including request and response details.",
  "request_example": {
    "name": "Example for GET /users",
    "request_headers": {
      "Authorization": "Bearer token123"
    },
    "request_body": "{\"id\": 1}",
    "response_status": 200,
    "response_body": "{\"id\": 1, \"name\": \"John Doe\"}",
    "duration_ms": 50
  },
  "response_example": {
    "id": 1,
    "api_spec_id": 123,
    "name": "Example for GET /users",
    "request_headers": {
      "Authorization": "Bearer token123"
    },
    "request_body": "{\"id\": 1}",
    "response_status": 200,
    "response_body": "{\"id\": 1, \"name\": \"John Doe\"}",
    "duration_ms": 50,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/api-specs/:sid/examples`

**Handler Implementation:**
`apispec.CreateExample`

---

