# Testcase Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `testcase` module provides the following API endpoints:

### ```json
{
  "description": "The GET / endpoint in the testcase module is not provided with specific handler code. Based on the related DTOs, it seems this module is designed to manage and interact with test cases, including creating, updating, and running them. However, without the actual handler code, the exact functionality of the GET / endpoint cannot be determined. The provided DTOs suggest that the module supports operations such as duplicating test cases, running test cases, and creating or updating test cases.",
  "request_example": null,
  "response_example": null
}
```

**Endpoint:**
<kbd>GET</kbd> `/`

**Handler Implementation:**
`testcase.unknown`

---

### ```json
{
  "description": "The POST / endpoint in the 'testcase' module is designed to handle various operations related to test cases. Depending on the request body, it can create, update, duplicate, or run a test case. The specific operation is determined by the structure of the request body.",
  "request_examples": {
    "create_test_case": {
      "description": "Creates a new test case based on the provided API specification and other details.",
      "body": {
        "api_spec_id": 1,
        "name": "New Test Case",
        "description": "This is a new test case for the API",
        "env": "development",
        "headers": {
          "Content-Type": "application/json"
        },
        "query_params": {
          "param1": "value1"
        },
        "path_params": {
          "id": "123"
        },
        "request_body": {
          "key": "value"
        },
        "pre_script": "console.log('Pre-script executed')",
        "post_script": "console.log('Post-script executed')",
        "assertions": [
          {
            "type": "status_code",
            "comparison": "equals",
            "target": 200
          }
        ],
        "extract_vars": [
          {
            "name": "token",
            "source": "response_body",
            "json_path": "$.auth.token"
          }
        ]
      }
    },
    "update_test_case": {
      "description": "Updates an existing test case with new details. Only the fields that need to be updated are included.",
      "body": {
        "name": "Updated Test Case Name",
        "description": "This is an updated description",
        "env": "staging",
        "headers": {
          "Authorization": "Bearer token"
        },
        "query_params": {
          "param2": "value2"
        },
        "path_params": {
          "id": "456"
        },
        "request_body": {
          "new_key": "new_value"
        },
        "pre_script": "console.log('Updated pre-script executed')",
        "post_script": "console.log('Updated post-script executed')",
        "assertions": [
          {
            "type": "response_time",
            "comparison": "less_than",
            "target": 500
          }
        ],
        "extract_vars": [
          {
            "name": "user_id",
            "source": "response_headers",
            "header_name": "X-User-ID"
          }
        ]
      }
    },
    "duplicate_test_case": {
      "description": "Duplicates an existing test case with a new name.",
      "body": {
        "name": "Duplicate Test Case"
      }
    },
    "run_test_case": {
      "description": "Runs a test case with optional environment and variable overrides.",
      "body": {
        "env_id": 2,
        "global_vars": {
          "base_url": "https://api.example.com"
        },
        "variable_keys": {
          "user_id": "12345"
        }
      }
    }
  }
}
```

**Endpoint:**
<kbd>POST</kbd> `/`

**Handler Implementation:**
`testcase.unknown`

---

### ```json
{
  "description": "This endpoint retrieves the details of a specific test case by its unique identifier (tcid). It allows users to fetch the complete information about a particular test case, including its name, description, environment, headers, query parameters, path parameters, request body, pre and post scripts, assertions, and extract variables.",
  "request": {
    "method": "GET",
    "url": "/testcase/:tcid",
    "parameters": {
      "path": {
        "tcid": "123"
      }
    },
    "headers": {},
    "body": null
  },
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "id": 123,
      "api_spec_id": 456,
      "name": "Test Case 1",
      "description": "This is a sample test case for API testing.",
      "env": "production",
      "headers": {
        "Authorization": "Bearer abcdefg123456789"
      },
      "query_params": {
        "page": "1",
        "limit": "10"
      },
      "path_params": {
        "userId": "12345"
      },
      "request_body": {
        "key": "value"
      },
      "pre_script": "console.log('Pre-script executed')",
      "post_script": "console.log('Post-script executed')",
      "assertions": [
        {
          "type": "equal",
          "property": "status",
          "expected": 200
        },
        {
          "type": "contains",
          "property": "response.body.message",
          "expected": "success"
        }
      ],
      "extract_vars": [
        {
          "name": "token",
          "from": "response.headers.Authorization"
        }
      ]
    }
  }
}
```

**Endpoint:**
<kbd>GET</kbd> `/:tcid`

**Handler Implementation:**
`testcase.unknown`

---

### {
  "description": "This endpoint allows for the partial update of a specific test case identified by its unique ID (tcid). It enables the modification of various attributes such as name, description, environment, headers, query parameters, path parameters, request body, pre and post scripts, assertions, and extract variables. Only the fields that need to be updated are required in the request, making it flexible for selective updates.",
  "request_example": {
    "name": "Updated Test Case Name",
    "description": "This is an updated description for the test case.",
    "env": "staging",
    "headers": {
      "Content-Type": "application/json"
    },
    "query_params": {
      "sort": "desc",
      "limit": "10"
    },
    "path_params": {
      "user_id": "12345"
    },
    "request_body": {
      "key": "value"
    },
    "pre_script": "console.log('Pre-script executed')",
    "post_script": "console.log('Post-script executed')",
    "assertions": [
      {
        "type": "equals",
        "target": "response.body.message",
        "expected": "Success"
      }
    ],
    "extract_vars": [
      {
        "source": "response.body.id",
        "var_name": "new_user_id"
      }
    ]
  },
  "response_example": {
    "status": "success",
    "message": "Test case updated successfully",
    "data": {
      "id": 1,
      "name": "Updated Test Case Name",
      "description": "This is an updated description for the test case.",
      "env": "staging",
      "headers": {
        "Content-Type": "application/json"
      },
      "query_params": {
        "sort": "desc",
        "limit": "10"
      },
      "path_params": {
        "user_id": "12345"
      },
      "request_body": {
        "key": "value"
      },
      "pre_script": "console.log('Pre-script executed')",
      "post_script": "console.log('Post-script executed')",
      "assertions": [
        {
          "type": "equals",
          "target": "response.body.message",
          "expected": "Success"
        }
      ],
      "extract_vars": [
        {
          "source": "response.body.id",
          "var_name": "new_user_id"
        }
      ]
    }
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/:tcid`

**Handler Implementation:**
`testcase.unknown`

---

### ```json
{
  "description": "This endpoint is used to delete a specific test case identified by its unique ID (tcid). When this DELETE request is made, the system will remove the test case and all associated data from the database. No request body is needed for this operation, as the test case ID is provided in the URL.",
  "request_example": {
    "method": "DELETE",
    "url": "/testcase/123",
    "headers": {
      "Content-Type": "application/json",
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
<kbd>DELETE</kbd> `/:tcid`

**Handler Implementation:**
`testcase.unknown`

---

### ```json
{
  "description": "This endpoint allows users to duplicate an existing test case by providing a new name for the duplicated test case. The original test case is identified by its unique ID (tcid). The duplicated test case will have the same configuration as the original, except for the name which will be set to the new name provided in the request.",
  "request_example": {
    "method": "POST",
    "url": "/testcase/123/duplicate",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "name": "New Test Case Name"
    }
  },
  "response_example": {
    "status": 201,
    "body": {
      "id": 456,
      "name": "New Test Case Name",
      "api_spec_id": 789,
      "description": "This is a duplicated test case.",
      "env": "production",
      "headers": {
        "Authorization": "Bearer token"
      },
      "query_params": {
        "page": "1",
        "limit": "10"
      },
      "path_params": {
        "user_id": "123"
      },
      "request_body": {
        "key": "value"
      },
      "pre_script": "console.log('Running pre-script...');",
      "post_script": "console.log('Running post-script...');",
      "assertions": [
        {
          "type": "equals",
          "property": "status",
          "value": 200
        }
      ],
      "extract_vars": [
        {
          "name": "token",
          "source": "response.headers.Authorization"
        }
      ]
    }
  }
}
```

**Endpoint:**
<kbd>POST</kbd> `/:tcid/duplicate`

**Handler Implementation:**
`testcase.unknown`

---

### {
  "description": "The POST /from-spec endpoint is used to create a new test case based on an existing API specification. The request includes the ID of the API specification, the name for the new test case, and optionally, the environment and an example ID if the user wants to use a specific example from the API specification as a base for the test case.",
  "request_example": {
    "api_spec_id": 1,
    "name": "New Test Case",
    "env": "development",
    "use_example": true,
    "example_id": 2
  },
  "response_example": {
    "id": 10,
    "name": "New Test Case",
    "api_spec_id": 1,
    "description": null,
    "env": "development",
    "headers": {},
    "query_params": {},
    "path_params": {},
    "request_body": null,
    "pre_script": "",
    "post_script": "",
    "assertions": [],
    "extract_vars": []
  }
}

**Endpoint:**
<kbd>POST</kbd> `/from-spec`

**Handler Implementation:**
`testcase.unknown`

---

### ```json
{
  "description": "This endpoint is used to run a specific test case identified by its unique ID (tcid). The request allows for optional parameters such as environment ID, global variables, and variable keys to be provided. These can be used to override the default settings of the test case, making it flexible for different testing scenarios.",
  "request_example": {
    "env_id": 1,
    "global_vars": {
      "user_id": 12345,
      "access_token": "abc123xyz"
    },
    "variable_keys": {
      "base_url": "https://api.example.com",
      "version": "v2"
    }
  },
  "response_example": {
    "status": "success",
    "message": "Test case executed successfully.",
    "execution_details": {
      "test_case_id": 101,
      "start_time": "2023-10-05T14:25:36Z",
      "end_time": "2023-10-05T14:25:40Z",
      "duration": 4.2,
      "result": "passed",
      "logs": [
        {
          "timestamp": "2023-10-05T14:25:37Z",
          "level": "INFO",
          "message": "Starting test case execution."
        },
        {
          "timestamp": "2023-10-05T14:25:38Z",
          "level": "DEBUG",
          "message": "Executing pre-script."
        },
        {
          "timestamp": "2023-10-05T14:25:39Z",
          "level": "INFO",
          "message": "Sending request to https://api.example.com/v2/data."
        },
        {
          "timestamp": "2023-10-05T14:25:40Z",
          "level": "INFO",
          "message": "Test case passed all assertions."
        }
      ]
    }
  }
}
```

**Endpoint:**
<kbd>POST</kbd> `/:tcid/run`

**Handler Implementation:**
`testcase.unknown`

---

