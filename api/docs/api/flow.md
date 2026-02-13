# Flow Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `flow` module provides the following API endpoints:

### {
  "description": "This endpoint retrieves a list of all flows associated with a specific project. It requires the project ID as a path parameter. The response includes an array of flow objects and the total number of flows.",
  "request_example": {
    "method": "GET",
    "url": "/projects/123/flows"
  },
  "response_example": {
    "items": [
      {
        "name": "User Registration Flow",
        "description": "Flow for user registration process",
        "steps": [
          {
            "name": "Step 1: Enter User Details",
            "sort_order": 1,
            "method": "POST",
            "url": "https://api.example.com/register",
            "headers": "Content-Type: application/json",
            "body": "{\"username\": \"testuser\", \"password\": \"testpass\"}",
            "captures": "",
            "asserts": "status_code == 200",
            "position_x": 100,
            "position_y": 100
          },
          {
            "name": "Step 2: Verify Email",
            "sort_order": 2,
            "method": "GET",
            "url": "https://api.example.com/verify-email",
            "headers": "",
            "body": "",
            "captures": "",
            "asserts": "status_code == 200",
            "position_x": 300,
            "position_y": 100
          }
        ],
        "edges": [
          {
            "source_step_id": 1,
            "target_step_id": 2,
            "variable_mapping": ""
          }
        ]
      }
    ],
    "total": 1
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/flows`

**Handler Implementation:**
`flow.ListFlows`

---

### {
  "description": "This endpoint allows the creation of a new flow within a specified project. A flow is a sequence of steps and edges (connections between steps) that define a particular process or workflow. The request must include the name of the flow, and optionally, a description. Upon successful creation, the response will contain details of the newly created flow.",
  "request_example": {
    "name": "New API Testing Flow",
    "description": "A flow for testing our public APIs"
  },
  "response_example": {
    "id": 1,
    "project_id": 123,
    "name": "New API Testing Flow",
    "description": "A flow for testing our public APIs",
    "created_at": "2023-10-05T14:25:36Z",
    "updated_at": "2023-10-05T14:25:36Z",
    "steps": [],
    "edges": []
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/flows`

**Handler Implementation:**
`flow.CreateFlow`

---

### {
  "description": "This endpoint retrieves a specific flow within a project by its unique identifier. It requires the project ID and the flow ID as part of the URL. If the flow is found, it returns the flow details; otherwise, it returns an error indicating that the flow was not found or another internal server error.",
  "request": {
    "method": "GET",
    "url": "/projects/123/flows/456",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response": {
    "status": 200,
    "body": {
      "id": 456,
      "name": "Example Flow",
      "description": "This is an example flow for demonstration purposes.",
      "steps": [
        {
          "id": 1,
          "name": "Step 1",
          "sort_order": 1,
          "method": "GET",
          "url": "https://api.example.com/data",
          "headers": "{\"Authorization\": \"Bearer token\"}",
          "body": "",
          "captures": "{}",
          "asserts": "{}",
          "position_x": 100,
          "position_y": 100
        },
        {
          "id": 2,
          "name": "Step 2",
          "sort_order": 2,
          "method": "POST",
          "url": "https://api.example.com/submit",
          "headers": "{\"Content-Type\": \"application/json\"}",
          "body": "{\"key\": \"value\"}",
          "captures": "{}",
          "asserts": "{}",
          "position_x": 300,
          "position_y": 100
        }
      ],
      "edges": [
        {
          "id": 1,
          "source_step_id": 1,
          "target_step_id": 2,
          "variable_mapping": ""
        }
      ]
    }
  },
  "error_responses": [
    {
      "status": 404,
      "body": {
        "error": "flow not found"
      }
    },
    {
      "status": 500,
      "body": {
        "error": "internal server error"
      }
    }
  ]
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/flows/:fid`

**Handler Implementation:**
`flow.GetFlow`

---

### {
  "description": "This endpoint allows updating the name and description of a specific flow within a project. The flow is identified by its ID (fid) and the project by its ID. Only the fields provided in the request will be updated, and any omitted fields will remain unchanged.",
  "request_example": {
    "method": "PATCH",
    "url": "/projects/123/flows/456",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "name": "Updated Flow Name",
      "description": "This is an updated description for the flow."
    }
  },
  "response_example": {
    "status": 200,
    "body": {
      "id": 456,
      "project_id": 123,
      "name": "Updated Flow Name",
      "description": "This is an updated description for the flow.",
      "created_at": "2023-10-01T12:00:00Z",
      "updated_at": "2023-10-02T12:00:00Z"
    }
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/projects/:id/flows/:fid`

**Handler Implementation:**
`flow.UpdateFlow`

---

### {
  "description": "This endpoint updates an existing flow in a project. It allows you to modify the flow's name, description, steps, and edges. The request should include the updated details of the flow, including any new or modified steps and edges. If the flow with the specified ID does not exist, a 404 Not Found error will be returned.",
  "request_example": {
    "name": "Updated Flow Name",
    "description": "This is an updated description of the flow.",
    "steps": [
      {
        "name": "Step 1",
        "sort_order": 1,
        "method": "GET",
        "url": "https://example.com/api/endpoint1",
        "headers": "{\"Content-Type\": \"application/json\"}",
        "body": "",
        "captures": "",
        "asserts": "",
        "position_x": 100,
        "position_y": 100
      },
      {
        "name": "Step 2",
        "sort_order": 2,
        "method": "POST",
        "url": "https://example.com/api/endpoint2",
        "headers": "{\"Content-Type\": \"application/json\"}",
        "body": "{\"key\": \"value\"}",
        "captures": "",
        "asserts": "",
        "position_x": 300,
        "position_y": 100
      }
    ],
    "edges": [
      {
        "source_step_id": 1,
        "target_step_id": 2,
        "variable_mapping": ""
      }
    ]
  },
  "response_example": {
    "id": 1,
    "project_id": 1,
    "name": "Updated Flow Name",
    "description": "This is an updated description of the flow.",
    "steps": [
      {
        "id": 1,
        "flow_id": 1,
        "name": "Step 1",
        "sort_order": 1,
        "method": "GET",
        "url": "https://example.com/api/endpoint1",
        "headers": "{\"Content-Type\": \"application/json\"}",
        "body": "",
        "captures": "",
        "asserts": "",
        "position_x": 100,
        "position_y": 100
      },
      {
        "id": 2,
        "flow_id": 1,
        "name": "Step 2",
        "sort_order": 2,
        "method": "POST",
        "url": "https://example.com/api/endpoint2",
        "headers": "{\"Content-Type\": \"application/json\"}",
        "body": "{\"key\": \"value\"}",
        "captures": "",
        "asserts": "",
        "position_x": 300,
        "position_y": 100
      }
    ],
    "edges": [
      {
        "id": 1,
        "flow_id": 1,
        "source_step_id": 1,
        "target_step_id": 2,
        "variable_mapping": ""
      }
    ]
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/projects/:id/flows/:fid`

**Handler Implementation:**
`flow.SaveFlow`

---

### {
  "description": "This endpoint deletes a specific flow (identified by `fid`) from a project (identified by `id`). If the flow is successfully deleted, it returns a 204 No Content status. If the flow does not exist, it returns a 404 Not Found status. For any other errors, it returns a 500 Internal Server Error status.",
  "request": {
    "method": "DELETE",
    "url": "/projects/123/flows/456",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response": {
    "success": {
      "status_code": 204,
      "body": ""
    },
    "not_found": {
      "status_code": 404,
      "body": {
        "error": "flow not found"
      }
    },
    "internal_server_error": {
      "status_code": 500,
      "body": {
        "error": "An unexpected error occurred"
      }
    }
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id/flows/:fid`

**Handler Implementation:**
`flow.DeleteFlow`

---

### ```json
{
  "description": "This endpoint is used to create a new step within a specific flow of a project. The step can include details such as the name, sort order, HTTP method, URL, headers, body, captures, asserts, and position on the canvas. The request must include the required fields: name, method, and URL.",
  "request_example": {
    "name": "Get User Details",
    "sort_order": 1,
    "method": "GET",
    "url": "https://api.example.com/users/123",
    "headers": "Authorization: Bearer {{token}}\nContent-Type: application/json",
    "body": "{}",
    "captures": "user_id: $.id\nuser_name: $.name",
    "asserts": "status_code == 200\n$.id > 0",
    "position_x": 150,
    "position_y": 100
  },
  "response_example": {
    "id": 1,
    "name": "Get User Details",
    "sort_order": 1,
    "method": "GET",
    "url": "https://api.example.com/users/123",
    "headers": "Authorization: Bearer {{token}}\nContent-Type: application/json",
    "body": "{}",
    "captures": "user_id: $.id\nuser_name: $.name",
    "asserts": "status_code == 200\n$.id > 0",
    "position_x": 150,
    "position_y": 100,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
  }
}
```

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/flows/:fid/steps`

**Handler Implementation:**
`flow.CreateStep`

---

### {
  "description": "This endpoint updates an existing step in a flow within a project. It allows partial updates, meaning only the fields provided in the request will be updated. The step ID is required as a path parameter. If the step is not found, a 404 Not Found error is returned. If the step ID is invalid, a 400 Bad Request error is returned.",
  "request_example": {
    "name": "Updated Step Name",
    "sort_order": 2,
    "method": "POST",
    "url": "https://api.example.com/updated-endpoint",
    "headers": "{\"Content-Type\": \"application/json\"}",
    "body": "{\"key\": \"value\"}",
    "captures": "response.body.data.id",
    "asserts": "response.statusCode == 200",
    "position_x": 150,
    "position_y": 300
  },
  "response_example": {
    "id": 1,
    "name": "Updated Step Name",
    "sort_order": 2,
    "method": "POST",
    "url": "https://api.example.com/updated-endpoint",
    "headers": "{\"Content-Type\": \"application/json\"}",
    "body": "{\"key\": \"value\"}",
    "captures": "response.body.data.id",
    "asserts": "response.statusCode == 200",
    "position_x": 150,
    "position_y": 300,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/projects/:id/flows/:fid/steps/:sid`

**Handler Implementation:**
`flow.UpdateStep`

---

### {
  "description": "This endpoint deletes a specific step from a flow within a project. The step is identified by its ID, which is part of the URL. If the step is successfully deleted, the server responds with a 204 No Content status. If the step ID is invalid or the step does not exist, the server returns a 400 Bad Request or 404 Not Found status, respectively. In case of other errors, a 500 Internal Server Error status is returned.",
  "request_example": {
    "method": "DELETE",
    "url": "/projects/123/flows/456/steps/789"
  },
  "response_example": {
    "status_code": 204,
    "body": ""
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id/flows/:fid/steps/:sid`

**Handler Implementation:**
`flow.DeleteStep`

---

### {
  "description": "This endpoint allows the creation of a new edge between two steps within a specific flow in a project. An edge represents a connection or transition from one step to another, and it can optionally include a variable mapping that defines how variables are passed between the steps.",
  "request_example": {
    "source_step_id": 1,
    "target_step_id": 2,
    "variable_mapping": "{ \"var1\": \"value1\", \"var2\": \"value2\" }"
  },
  "response_example": {
    "id": 3,
    "source_step_id": 1,
    "target_step_id": 2,
    "variable_mapping": "{ \"var1\": \"value1\", \"var2\": \"value2\" }",
    "created_at": "2023-10-05T14:25:30Z",
    "updated_at": "2023-10-05T14:25:30Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/flows/:fid/edges`

**Handler Implementation:**
`flow.CreateEdge`

---

### {
  "description": "This endpoint deletes a specific edge in a flow within a project. The edge is identified by its unique ID. Upon successful deletion, the server responds with a 204 No Content status. If the edge ID is invalid, the server returns a 400 Bad Request. If the edge does not exist, a 404 Not Found error is returned. For any other errors, a 500 Internal Server Error is returned.",
  "request": {
    "method": "DELETE",
    "url": "/projects/123/flows/456/edges/789",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": null
  },
  "response": {
    "status_code": 204,
    "body": null
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id/flows/:fid/edges/:eid`

**Handler Implementation:**
`flow.DeleteEdge`

---

### {
  "description": "This endpoint triggers the execution of a specified flow within a project. The request does not require a body, as it uses the flow ID and user ID to initiate the run. Upon successful execution, it returns the details of the run.",
  "request_example": {},
  "response_example": {
    "id": 1,
    "flow_id": 123,
    "user_id": 456,
    "status": "running",
    "started_at": "2023-10-05T14:25:36Z",
    "completed_at": null
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/flows/:fid/run`

**Handler Implementation:**
`flow.RunFlow`

---

### ```json
{
  "description": "This endpoint retrieves a list of all runs for a specific flow within a project. It returns the runs along with the total count of runs.",
  "request": {
    "method": "GET",
    "url": "/projects/:id/flows/:fid/runs",
    "params": {
      "id": "The ID of the project",
      "fid": "The ID of the flow"
    }
  },
  "response": {
    "200": {
      "description": "A successful response containing a list of runs and the total count of runs.",
      "example": {
        "items": [
          {
            "id": 1,
            "flow_id": 123,
            "status": "completed",
            "started_at": "2023-10-01T12:00:00Z",
            "finished_at": "2023-10-01T12:05:00Z"
          },
          {
            "id": 2,
            "flow_id": 123,
            "status": "failed",
            "started_at": "2023-10-01T12:10:00Z",
            "finished_at": "2023-10-01T12:15:00Z"
          }
        ],
        "total": 2
      }
    },
    "500": {
      "description": "An internal server error occurred while fetching the runs.",
      "example": {
        "error": "Internal Server Error"
      }
    }
  }
}
```

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/flows/:fid/runs`

**Handler Implementation:**
`flow.ListRuns`

---

### {
  "description": "This endpoint retrieves the details of a specific run (identified by its ID) within a flow, which is part of a project. The response includes information about the run, such as its steps and edges. If the run ID is invalid or the run does not exist, appropriate error responses are returned.",
  "request_example": null,
  "response_example": {
    "run": {
      "id": 123,
      "flow_id": 456,
      "project_id": 789,
      "status": "completed",
      "started_at": "2023-10-01T12:00:00Z",
      "ended_at": "2023-10-01T12:15:00Z",
      "steps": [
        {
          "name": "Step 1",
          "sort_order": 1,
          "method": "GET",
          "url": "https://example.com/api/v1/data",
          "headers": "{\"Authorization\": \"Bearer token\"}",
          "body": "",
          "captures": "",
          "asserts": "",
          "position_x": 100,
          "position_y": 100
        },
        {
          "name": "Step 2",
          "sort_order": 2,
          "method": "POST",
          "url": "https://example.com/api/v1/submit",
          "headers": "{\"Content-Type\": \"application/json\"}",
          "body": "{\"key\": \"value\"}",
          "captures": "",
          "asserts": "",
          "position_x": 300,
          "position_y": 100
        }
      ],
      "edges": [
        {
          "source_step_id": 1,
          "target_step_id": 2,
          "variable_mapping": ""
        }
      ]
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/flows/:fid/runs/:rid`

**Handler Implementation:**
`flow.GetRun`

---

### {
  "description": "This endpoint streams Server-Sent Events (SSE) for a specific flow run. It allows clients to receive real-time updates about the execution of a flow, including step events and a final 'done' event when the flow completes. The response is in SSE format, which is a standard for sending automatic updates from the server to the client over a single HTTP connection.",
  "request_example": {
    "method": "GET",
    "url": "/projects/123/flows/456/runs/789/events",
    "headers": {
      "Accept": "text/event-stream"
    }
  },
  "response_example": [
    {
      "event": "step",
      "data": {
        "id": 1,
        "name": "Step 1",
        "status": "success",
        "output": "Step 1 executed successfully"
      }
    },
    {
      "event": "step",
      "data": {
        "id": 2,
        "name": "Step 2",
        "status": "error",
        "output": "Step 2 failed: Invalid input"
      }
    },
    {
      "event": "done",
      "data": {}
    }
  ]
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/flows/:fid/runs/:rid/events`

**Handler Implementation:**
`flow.ExecuteFlowSSE`

---

