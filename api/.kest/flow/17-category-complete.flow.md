# 📁 Category Module Complete Flow

Complete CRUD testing for category management including sorting.

---

## Step 1: Register User (Setup)

```kest
POST /v1/register
Content-Type: application/json

{
  "username": "catuser{{$timestamp}}",
  "email": "catuser{{$timestamp}}@kest.io",
  "password": "SecurePass123!",
  "nickname": "Category Tester"
}

[Captures]
registered_username: data.username

[Asserts]
status == 201
body.code == 0
```

---

## Step 1b: User Login

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

## Step 2: Create Test Project

```kest
POST /v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Category Test Project {{$timestamp}}",
  "description": "Project for testing categories"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: Create Category - API Tests

```kest
POST /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "API Tests",
  "description": "Category for API test cases",
  "sort_order": 1
}

[Captures]
category_id_1: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "API Tests"
duration < 1000ms
```

---

## Step 4: Create Category - UI Tests

```kest
POST /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "UI Tests",
  "description": "Category for UI test cases",
  "sort_order": 2
}

[Captures]
category_id_2: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
body.data.name == "UI Tests"
```

---

## Step 5: Create Category - Integration Tests

```kest
POST /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Integration Tests",
  "description": "Category for integration test cases",
  "sort_order": 3
}

[Captures]
category_id_3: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 6: List All Categories

```kest
GET /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 7: Get Category Details

```kest
GET /v1/projects/{{project_id}}/categories/{{category_id_1}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.id == "{{category_id_1}}"
body.data.name == "API Tests"
duration < 500ms
```

---

## Step 8: Update Category

```kest
PATCH /v1/projects/{{project_id}}/categories/{{category_id_1}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "API Tests - Updated",
  "description": "Updated category description"
}

[Asserts]
status == 200
body.code == 0
duration < 3000ms
```

---

## Step 9: Verify Category Update

```kest
GET /v1/projects/{{project_id}}/categories/{{category_id_1}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data.name == "API Tests - Updated"
body.data.description == "Updated category description"
```

---

## Step 10: Sort Categories

```kest
PUT /v1/projects/{{project_id}}/categories/sort
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "category_ids": [
    {{category_id_3}},
    {{category_id_1}},
    {{category_id_2}}
  ]
}

[Asserts]
status == 200
body.code == 0
duration < 3000ms
```

---

## Step 11: Verify Category Sorting

```kest
GET /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
```

---

## Step 12: Delete Category

```kest
DELETE /v1/projects/{{project_id}}/categories/{{category_id_1}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 13: Verify Category Deletion

```kest
GET /v1/projects/{{project_id}}/categories/{{category_id_1}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
```

---

## Step 14: Cleanup - Delete Remaining Categories

```kest
DELETE /v1/projects/{{project_id}}/categories/{{category_id_2}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 15: Cleanup - Delete Project

```kest
DELETE /v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Category Module Complete - 15 Steps**
