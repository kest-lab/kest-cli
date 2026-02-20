# Kest API Complete Test Flow

测试所有新功能：默认值、变量追踪、fail-fast、strict 模式

---

## Test 1: 用户注册（使用默认值）

### Step: Register User
POST https://httpbin.org/post
Content-Type: application/json
```json
{
  "username": "{{username | default: \"testuser\"}}",
  "email": "{{email | default: \"test@example.com\"}}",
  "password": "{{password | default: \"Test@123\"}}",
  "timestamp": "{{$timestamp}}",
  "random_id": "{{$randomInt}}"
}
```

[Asserts]
- status == 200

[Captures]
- user_data = data.json
- request_id = data.headers.X-Request-Id

---

## Test 2: 用户登录（使用捕获的变量）

### Step: Login
POST https://httpbin.org/post
Content-Type: application/json
```json
{
  "username": "{{username | default: \"testuser\"}}",
  "password": "{{password | default: \"Test@123\"}}",
  "device_id": "device_{{$randomInt}}"
}
```

[Asserts]
- status == 200

[Captures]
- auth_token = data.json.username
- login_time = data.json.password

---

## Test 3: 获取用户信息（使用 token）

### Step: Get User Profile
GET https://httpbin.org/get?user={{username | default: "testuser"}}&time={{$timestamp}}
Authorization: Bearer {{auth_token | default: "default_token"}}
Content-Type: application/json

[Asserts]
- status == 200

[Captures]
- profile_data = data.args.user

---

## Test 4: 更新用户信息

### Step: Update Profile
PUT https://httpbin.org/put
Authorization: Bearer {{auth_token | default: "default_token"}}
Content-Type: application/json
```json
{
  "username": "{{username | default: \"testuser\"}}",
  "email": "{{new_email | default: \"updated@example.com\"}}",
  "bio": "{{bio | default: \"Test user bio\"}}",
  "updated_at": "{{$timestamp}}"
}
```

[Asserts]
- status == 200

---

## Test 5: 测试内置变量

### Step: Test Built-in Variables
POST https://httpbin.org/post
Content-Type: application/json
```json
{
  "timestamp1": "{{$timestamp}}",
  "timestamp2": "{{$timestamp}}",
  "random1": "{{$randomInt}}",
  "random2": "{{$randomInt}}",
  "random3": "{{$randomInt}}"
}
```

[Asserts]
- status == 200

---

## Test 6: 混合变量测试

### Step: Mixed Variables
POST https://httpbin.org/post
Content-Type: application/json
```json
{
  "user": "{{username | default: \"testuser\"}}",
  "env": "{{environment | default: \"development\"}}",
  "api_version": "{{api_version | default: \"v1\"}}",
  "timestamp": "{{$timestamp}}",
  "request_id": "req_{{$randomInt}}"
}
```

[Asserts]
- status == 200

[Captures]
- final_data = data.json
