# Variable Features Test Flow

Test the new variable features: default values, strict mode, and fail-fast.

---

## Test 1: Default Values

### Step: Test Default Value
POST https://httpbin.org/post
```json
{
  "username": "{{username | default: \"test_user\"}}",
  "password": "{{password | default: \"Test@123\"}}",
  "timestamp": "{{$timestamp}}",
  "random": "{{$randomInt}}"
}
```

[Asserts]
- status == 200

[Captures]
- request_data = data.json

---

## Test 2: Variable Override

### Step: Test CLI Override
POST https://httpbin.org/post
```json
{
  "api_key": "{{api_key | default: \"default_key\"}}",
  "env": "{{env | default: \"dev\"}}"
}
```

[Asserts]
- status == 200

---

## Test 3: Capture and Reuse

### Step: Create Resource
POST https://httpbin.org/post
```json
{
  "name": "Resource_{{$timestamp}}",
  "type": "test"
}
```

[Captures]
- resource_name = data.json.name
- resource_type = data.json.type

### Step: Use Captured Variable
GET https://httpbin.org/get?name={{resource_name}}&type={{resource_type}}

[Asserts]
- status == 200
