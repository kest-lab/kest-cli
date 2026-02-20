# 路由调试测试

测试 environments 和 categories 路由是否正常工作

---

## Step 1: 登录获取 token

```kest
POST /v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

[Captures]
access_token: data.access_token
user_id: data.user.id
```

---

## Step 2: 创建测试项目

```kest
POST /v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Route Test Project",
  "slug": "route-test-{{$timestamp}}"
}

[Captures]
project_id: data.id
```

---

## Step 3: 测试 Environments 列表（应该返回 200）

```kest
GET /v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
```

---

## Step 4: 创建 Environment

```kest
POST /v1/projects/{{project_id}}/environments
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "test-env",
  "display_name": "Test Environment",
  "base_url": "http://localhost:8080"
}

[Captures]
env_id: data.id

[Asserts]
status == 201
```

---

## Step 5: 测试 Categories 列表（应该返回 200）

```kest
GET /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
```

---

## Step 6: 创建 Category

```kest
POST /v1/projects/{{project_id}}/categories
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "API Tests",
  "description": "API test cases category"
}

[Captures]
category_id: data.id

[Asserts]
status == 201
```

---

## Step 7: 清理 - 删除 Environment

```kest
DELETE /v1/projects/{{project_id}}/environments/{{env_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 8: 清理 - 删除 Category

```kest
DELETE /v1/projects/{{project_id}}/categories/{{category_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```

---

## Step 9: 清理 - 删除项目

```kest
DELETE /v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 204
```
