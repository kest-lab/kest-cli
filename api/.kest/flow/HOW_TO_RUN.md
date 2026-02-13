# 🚀 How to Run Kest Flow Tests

## Quick Start

### Run a Single Flow Test

```bash
kest run flow auth.flow.md
```

### Example Flow File

```markdown
# auth.flow.md

## Step 1: Login
POST /api/auth/login
Content-Type: application/json

{
  "user": "admin"
}

[Captures]
token: body.access_token

[Asserts]
status == 200
duration < 200ms
```

---

## Flow File Structure

### 1. Request Section
```
POST /api/auth/login
Content-Type: application/json

{
  "user": "admin",
  "password": "password123"
}
```

### 2. Captures (Optional)
Save response data for later steps:
```
[Captures]
token: body.access_token
user_id: body.data.id
```

### 3. Asserts
Validate the response:
```
[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 200ms
```

---

## Common Commands

### Run Flow Tests

```bash
# Run specific flow
kest run flow auth.flow.md

# Run all flows in directory
kest run flow .

# Run with environment
kest run flow auth.flow.md --env production
```

### View Test Results

```bash
# View latest test results
kest logs

# View specific test
kest logs --flow auth.flow.md

# View failed tests only
kest logs --failed
```

---

## Assertion Examples

### Status Codes
```
status == 200
status >= 200
status < 300
```

### Response Body
```
body.code == 0
body.data.id exists
body.data.username == "admin"
body.data.items.length > 0
```

### Performance
```
duration < 200ms
duration < 1s
```

### Headers
```
header.Content-Type == "application/json"
header.Authorization exists
```

---

## Using Captured Variables

Capture data from one step and use it in later steps:

```markdown
## Step 1: Login
POST /api/auth/login
{ "user": "admin" }

[Captures]
token: body.access_token

[Asserts]
status == 200

---

## Step 2: Get Profile
GET /api/users/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.data.username == "admin"
```

---

## Built-in Variables

- `{{$timestamp}}` - Current Unix timestamp
- `{{$uuid}}` - Random UUID
- `{{$random}}` - Random number

Example:
```
{
  "email": "user{{$timestamp}}@example.com",
  "id": "{{$uuid}}"
}
```

---

## Multi-Environment Support

### Define Environments

Create `.kest/env/production.env`:
```
BASE_URL=https://api.kest.dev
API_KEY=your-api-key
```

### Use in Flow
```markdown
POST {{BASE_URL}}/api/auth/login
X-API-Key: {{API_KEY}}
```

### Run with Environment
```bash
kest run flow auth.flow.md --env production
```

---

## Best Practices

1. **Descriptive Step Names**: Use clear, action-oriented names
   ```markdown
   ## Step 1: User Login with Valid Credentials
   ```

2. **Comprehensive Assertions**: Test status, body, and performance
   ```
   [Asserts]
   status == 200
   body.code == 0
   body.data.token exists
   duration < 500ms
   ```

3. **Capture Important Data**: Save tokens, IDs for later steps
   ```
   [Captures]
   access_token: body.data.access_token
   user_id: body.data.user.id
   ```

4. **Use Variables**: Avoid hardcoding dynamic values
   ```
   "email": "test{{$timestamp}}@example.com"
   ```

5. **Organize Flows**: Group related tests in separate files
   - `auth.flow.md` - Authentication tests
   - `users.flow.md` - User management
   - `projects.flow.md` - Project operations

---

## Example: Complete Authentication Flow

```markdown
# auth.flow.md

## Step 1: Register New User
POST /api/v1/register
Content-Type: application/json

{
  "username": "user{{$timestamp}}",
  "email": "user{{$timestamp}}@example.com",
  "password": "SecurePass123!"
}

[Captures]
username: body.data.username
email: body.data.email

[Asserts]
status == 201
body.code == 0
body.data.username exists

---

## Step 2: Login
POST /api/v1/login
Content-Type: application/json

{
  "username": "{{username}}",
  "password": "SecurePass123!"
}

[Captures]
token: body.data.access_token

[Asserts]
status == 200
body.code == 0
body.data.access_token exists
duration < 200ms

---

## Step 3: Get User Profile
GET /api/v1/users/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
body.code == 0
body.data.username == "{{username}}"
body.data.email == "{{email}}"
```

---

## Troubleshooting

### View Detailed Logs
```bash
kest logs --verbose
```

### Debug Failed Assertions
```bash
kest run flow auth.flow.md --debug
```

### Check Request/Response
All requests and responses are logged in `.kest/logs/`

---

## Learn More

- [Flow Best Practices](./FLOW_BEST_PRACTICES.md)
- [Multi-Environment Guide](./MULTI_ENV_GUIDE.md)
- [Quick Reference](./QUICK_REFERENCE.md)
