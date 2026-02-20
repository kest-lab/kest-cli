# 🔍 Production API Discovery

Discover the actual API structure at https://api.kest.dev

---

## Step 1: Test Root Path

```kest
GET https://api.kest.dev/

[Asserts]
status in [200, 404]
```

---

## Step 2: Test Health without Prefix

```kest
GET https://api.kest.dev/health

[Asserts]
status in [200, 404]
```

---

## Step 3: Test API v1 Group

```kest
GET https://api.kest.dev/v1/health

[Asserts]
status in [200, 404]
```

---

## Step 4: Test Register with Different Paths

```kest
POST https://api.kest.dev/v1/register
Content-Type: application/json

{
  "username": "test_user_{{$randomInt}}",
  "password": "TestPass123",
  "email": "test_{{$timestamp}}@example.com"
}

[Asserts]
status in [200, 201, 400, 404]
```

---

## Step 5: Test Login with Different Paths

```kest
POST https://api.kest.dev/v1/login
Content-Type: application/json

{
  "username": "test",
  "password": "test"
}

[Asserts]
status in [200, 400, 401, 404]
```

---

## Step 6: Test Projects Endpoint

```kest
GET https://api.kest.dev/v1/projects

[Asserts]
status in [200, 401, 404]
```

---

## Step 7: Check API Documentation

```kest
GET https://api.kest.dev/docs

[Asserts]
status in [200, 404]
```

---

## Step 8: Check Swagger

```kest
GET https://api.kest.dev/swagger

[Asserts]
status in [200, 404]
```

---

## Step 9: Check OpenAPI JSON

```kest
GET https://api.kest.dev/swagger.json

[Asserts]
status in [200, 404]
```

---

## Step 10: Test Common API Paths

```kest
GET https://api.kest.dev/api

[Asserts]
status in [200, 404]
```

```kest
GET https://api.kest.dev/v1

[Asserts]
status in [200, 404]
```

```kest
GET https://api.kest.dev/api/health

[Asserts]
status in [200, 404]
```

---

## Step 11: Test CORS Preflight

```kest
OPTIONS https://api.kest.dev/
Origin: https://kest.dev
Access-Control-Request-Method: POST

[Asserts]
status in [200, 204, 404]
```

---

## Step 12: Test with Different Subdomain

```kest
GET https://api.kest.dev/status

[Asserts]
status in [200, 404]
```

---

## Step 13: Test Web App Path

```kest
GET https://api.kest.dev/app

[Asserts]
status in [200, 404]
```

---

## Step 14: Test API Info Endpoint

```kest
GET https://api.kest.dev/info

[Asserts]
status in [200, 404]
```

---

## Step 15: Test Version Endpoint

```kest
GET https://api.kest.dev/version

[Asserts]
status in [200, 404]
```

---

## Discovery Results

This test will help identify:
- The correct base path for the API
- Whether the API is active at this domain
- Available documentation endpoints
- CORS configuration
- API versioning scheme
