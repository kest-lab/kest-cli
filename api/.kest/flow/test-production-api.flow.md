# Production API Test - api.kest.dev

测试生产环境 API 接口

---

## Step 1: 健康检查

```kest
GET https://api.kest.dev/v1/health

[Asserts]
status >= 200
status < 300
body.status == "ok"
body.version exists
duration < 2000ms
```

---

## Step 2: 用户注册

```kest
POST https://api.kest.dev/v1/register
Content-Type: application/json

{
  "username": "testuser{{$timestamp}}",
  "email": "test{{$timestamp}}@example.com",
  "password": "Test123456"
}

[Captures]
registered_username: data.username
registered_email: data.email

[Asserts]
status >= 200
status < 300
body.data.username exists
body.data.email exists
duration < 3000ms
```

---

## Step 3: 用户登录

```kest
POST https://api.kest.dev/v1/login
Content-Type: application/json

{
  "username": "{{registered_username}}",
  "password": "Test123456"
}

[Captures]
access_token: data.access_token

[Asserts]
status >= 200
status < 300
body.data.access_token exists
body.data.user.username == "{{registered_username}}"
duration < 3000ms
```

---

## Step 4: 获取用户资料

```kest
GET https://api.kest.dev/v1/users/profile
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
status < 300
body.data.username == "{{registered_username}}"
body.data.email == "{{registered_email}}"
duration < 2000ms
```

---

## Step 5: 创建项目

```kest
POST https://api.kest.dev/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Test Project {{$timestamp}}",
  "platform": "javascript"
}

[Captures]
project_id: data.id
project_name: data.name

[Asserts]
status >= 200
status < 300
body.data.id exists
body.data.name exists
body.data.slug exists
duration < 3000ms
```

---

## Step 6: 获取项目详情

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
status < 300
body.data.id == "{{project_id}}"
body.data.name == "{{project_name}}"
duration < 2000ms
```

---

## Step 7: 获取项目列表

```kest
GET https://api.kest.dev/v1/projects
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
status < 300
body.data.data exists
body.data.meta.total >= 1
duration < 2000ms
```

---

## Step 8: 删除项目

```kest
DELETE https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status >= 200
status < 300
body.code == 0
duration < 2000ms
```

---

## Step 9: 验证项目已删除

```kest
GET https://api.kest.dev/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 404
body.code == 404
duration < 1000ms
```
