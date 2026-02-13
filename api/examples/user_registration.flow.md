# User Registration Flow Test

This flow demonstrates how to test the user registration API endpoint.

## Test Scenario

1. Register a new user with random credentials
2. Verify the response contains user data
3. Verify password is not exposed in response

---

## Registration Test

```kest
POST /register
Content-Type: application/json

{
  "username": "test_user_{{$randomInt}}",
  "password": "TestPassword123",
  "email": "test{{$randomInt}}@example.com",
  "nickname": "Test User",
  "phone": "+1234567890"
}

[Captures]
user_id: data.id
username: data.username
email: data.email

[Asserts]
status == 201
body.data.id exists
body.data.username exists
body.data.email exists
body.data.nickname exists
body.data.password !exists
duration < 1000ms
```

---

## Validation Tests

### Test: Username Too Short

```kest
POST /register
Content-Type: application/json

{
  "username": "ab",
  "password": "TestPassword123",
  "email": "test@example.com"
}

[Asserts]
status == 400
body.message exists
```

### Test: Invalid Email Format

```kest
POST /register
Content-Type: application/json

{
  "username": "testuser",
  "password": "TestPassword123",
  "email": "invalid-email"
}

[Asserts]
status == 400
body.message exists
```

### Test: Password Too Short

```kest
POST /register
Content-Type: application/json

{
  "username": "testuser",
  "password": "12345",
  "email": "test@example.com"
}

[Asserts]
status == 400
body.message exists
```

---

## Notes

- Uses `{{$randomInt}}` to generate unique usernames and emails
- Captures user_id for use in subsequent tests
- Validates that password is NOT returned in response
- Tests validation rules for username, email, and password
