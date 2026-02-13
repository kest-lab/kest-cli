# 🔄 Flow Module Complete Test

Complete CRUD testing for flow management including steps, edges, and runs.

---

## Step 1: Register User

```kest
POST /v1/register
Content-Type: application/json

{
  "username": "flowuser{{$timestamp}}",
  "email": "flowuser{{$timestamp}}@kest.io",
  "password": "SecurePass123!",
  "nickname": "Flow Tester"
}

[Captures]
registered_username: data.username

[Asserts]
status == 201
body.code == 0
```

---

## Step 2: Login

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "SecurePass123!"
}

[Captures]
access_token: data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
```

---

## Step 3: Create Project

```kest
POST /v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Flow Test Project {{$timestamp}}",
  "description": "Project for testing flows"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 4: Create Flow

```kest
POST /v1/projects/{{project_id}}/flows
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Login Test Flow",
  "description": "Tests the login scenario"
}

[Captures]
flow_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "Login Test Flow"
```

---

## Step 5: Create Step 1 - Register

```kest
POST /v1/projects/{{project_id}}/flows/{{flow_id}}/steps
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Register User",
  "sort_order": 0,
  "method": "POST",
  "url": "/v1/register",
  "headers": "{\"Content-Type\":\"application/json\"}",
  "body": "{\"username\":\"test\",\"password\":\"pass\"}",
  "captures": "username: data.username",
  "asserts": "status == 201",
  "position_x": 100,
  "position_y": 100
}

[Captures]
step_id_1: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.method == "POST"
```

---

## Step 6: Create Step 2 - Login

```kest
POST /v1/projects/{{project_id}}/flows/{{flow_id}}/steps
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Login",
  "sort_order": 1,
  "method": "POST",
  "url": "/v1/login",
  "headers": "{\"Content-Type\":\"application/json\"}",
  "body": "{\"username\":\"{{username}}\",\"password\":\"pass\"}",
  "captures": "access_token: data.access_token",
  "asserts": "status == 200",
  "position_x": 100,
  "position_y": 300
}

[Captures]
step_id_2: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 7: Create Edge (Step1 -> Step2)

```kest
POST /v1/projects/{{project_id}}/flows/{{flow_id}}/edges
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "source_step_id": {{step_id_1}},
  "target_step_id": {{step_id_2}},
  "variable_mapping": "{\"username\":\"username\"}"
}

[Captures]
edge_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.source_step_id exists
body.data.target_step_id exists
```

---

## Step 8: Get Flow Detail (with steps and edges)

```kest
GET /v1/projects/{{project_id}}/flows/{{flow_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "Login Test Flow"
body.data.steps exists
body.data.edges exists
body.data.step_count == 2
```

---

## Step 9: List Flows

```kest
GET /v1/projects/{{project_id}}/flows
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.items exists
body.data.total == 1
```

---

## Step 10: Update Flow

```kest
PATCH /v1/projects/{{project_id}}/flows/{{flow_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Login Test Flow - Updated",
  "description": "Updated description"
}

[Asserts]
status == 200
body.code == 0
body.data.name == "Login Test Flow - Updated"
```

---

## Step 11: Update Step

```kest
PATCH /v1/projects/{{project_id}}/flows/{{flow_id}}/steps/{{step_id_1}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Register User - Updated",
  "position_x": 200,
  "position_y": 200
}

[Asserts]
status == 200
body.code == 0
body.data.name == "Register User - Updated"
```

---

## Step 12: Run Flow

```kest
POST /v1/projects/{{project_id}}/flows/{{flow_id}}/run
Authorization: Bearer {{access_token}}

[Captures]
run_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.status == "pending"
```

---

## Step 13: Get Run Detail

```kest
GET /v1/projects/{{project_id}}/flows/{{flow_id}}/runs/{{run_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id exists
body.data.step_results exists
```

---

## Step 14: List Runs

```kest
GET /v1/projects/{{project_id}}/flows/{{flow_id}}/runs
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.items exists
body.data.total == 1
```

---

## Step 15: Delete Edge

```kest
DELETE /v1/projects/{{project_id}}/flows/{{flow_id}}/edges/{{edge_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 16: Delete Step

```kest
DELETE /v1/projects/{{project_id}}/flows/{{flow_id}}/steps/{{step_id_2}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 17: Delete Flow

```kest
DELETE /v1/projects/{{project_id}}/flows/{{flow_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 18: Verify Flow Deleted

```kest
GET /v1/projects/{{project_id}}/flows/{{flow_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
```

---

## Step 19: Cleanup - Delete Project

```kest
DELETE /v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Flow Module Complete - 19 Steps**
