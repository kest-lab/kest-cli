# 新功能验证测试

测试所有新功能：默认值、内置变量、strict 模式、fail-fast

---

## Test 1: 默认值语法测试

```kest
POST https://httpbin.org/post
Content-Type: application/json

{
  "username": "{{username | default: \"admin\"}}",
  "password": "{{password | default: \"Test@123\"}}",
  "email": "{{email | default: \"test@example.com\"}}"
}

[Asserts]
status == 200

[Captures]
response_data: data.json
```

---

## Test 2: 内置变量测试

```kest
POST https://httpbin.org/post
Content-Type: application/json

{
  "timestamp": "{{$timestamp}}",
  "random1": "{{$randomInt}}",
  "random2": "{{$randomInt}}"
}

[Asserts]
status == 200
```

---

## Test 3: 混合变量测试

```kest
POST https://httpbin.org/post
Content-Type: application/json

{
  "user": "{{username | default: \"testuser\"}}",
  "id": "user_{{$randomInt}}",
  "created_at": "{{$timestamp}}"
}

[Asserts]
status == 200
```
