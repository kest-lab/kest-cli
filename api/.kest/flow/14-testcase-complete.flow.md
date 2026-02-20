# 🧪 Test Case Module Complete Flow

Complete CRUD testing for test case management including duplication and execution.

---

## Step 1: User Login & Create Project (Setup)

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
```

---

## Step 2: Create Test Project

```kest
POST /v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "TestCase Test Project {{$timestamp}}",
  "description": "Project for testing test cases"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: Create Test Case

```kest
POST /v1/projects/{{project_id}}/test-cases
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "User Login Test",
  "description": "Test user login functionality",
  "method": "POST",
  "path": "/api/login",
  "request_body": {
    "username": "test",
    "password": "test123"
  },
  "expected_status": 200,
  "expected_response": {
    "success": true
  }
}

[Captures]
testcase_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "User Login Test"
duration < 1000ms
```

---

## Step 4: Get Test Case Details

```kest
GET /v1/projects/{{project_id}}/test-cases/{{testcase_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{testcase_id}}"
body.data.name == "User Login Test"
body.data.method == "POST"
duration < 500ms
```

---

## Step 5: List All Test Cases

```kest
GET /v1/projects/{{project_id}}/test-cases
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 6: Update Test Case

```kest
PATCH /v1/projects/{{project_id}}/test-cases/{{testcase_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "User Login Test - Updated",
  "description": "Updated test for user login",
  "expected_status": 201
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 7: Verify Test Case Update

```kest
GET /v1/projects/{{project_id}}/test-cases/{{testcase_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "User Login Test - Updated"
body.data.expected_status == 201
```

---

## Step 8: Duplicate Test Case

```kest
POST /v1/projects/{{project_id}}/test-cases/{{testcase_id}}/duplicate
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "User Login Test - Copy"
}

[Captures]
duplicated_testcase_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "User Login Test - Copy"
duration < 1000ms
```

---

## Step 9: Run Test Case

```kest
POST /v1/projects/{{project_id}}/test-cases/{{testcase_id}}/run
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "environment_id": "default"
}

[Captures]
run_result_id: data.id

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 3000ms
```

---

## Step 10: Delete Original Test Case

```kest
DELETE /v1/projects/{{project_id}}/test-cases/{{testcase_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 11: Delete Duplicated Test Case

```kest
DELETE /v1/projects/{{project_id}}/test-cases/{{duplicated_testcase_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 12: Cleanup - Delete Project

```kest
DELETE /v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Test Case Module Complete - 12 Steps**
