# 👥 Member Module Complete Flow

Complete CRUD testing for project member management.

---

## Step 1: User Login & Create Project (Setup)

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

[Captures]
access_token: data.access_token
admin_user_id: data.user.id

[Asserts]
status == 200
body.code == 0
```

---

## Step 2: Create Test Project

```kest
POST /api/v1/projects
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "name": "Member Test Project {{$timestamp}}",
  "description": "Project for testing member management"
}

[Captures]
project_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 3: Register Second User (for member testing)

```kest
POST /api/v1/register
Content-Type: application/json

{
  "username": "member{{$timestamp}}",
  "email": "member{{$timestamp}}@kest.io",
  "password": "MemberPass123!",
  "nickname": "Test Member"
}

[Captures]
member_user_id: data.id

[Asserts]
status == 201
body.code == 0
body.data.id exists
```

---

## Step 4: Add Member to Project

```kest
POST /api/v1/projects/{{project_id}}/members
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "user_id": "{{member_user_id}}",
  "role": "write"
}

[Asserts]
status == 201
body.code == 0
body.data.user_id == "{{member_user_id}}"
body.data.role == "write"
duration < 1000ms
```

---

## Step 5: List Project Members

```kest
GET /api/v1/projects/{{project_id}}/members
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
duration < 1000ms
```

---

## Step 6: Update Member Role

```kest
PATCH /api/v1/projects/{{project_id}}/members/{{member_user_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "role": "admin"
}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 7: Verify Member Role Update

```kest
GET /api/v1/projects/{{project_id}}/members
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
body.data exists
```

---

## Step 8: Remove Member from Project

```kest
DELETE /api/v1/projects/{{project_id}}/members/{{member_user_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
duration < 1000ms
```

---

## Step 9: Verify Member Removal

```kest
GET /api/v1/projects/{{project_id}}/members
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

## Step 10: Cleanup - Delete Project

```kest
DELETE /api/v1/projects/{{project_id}}
Authorization: Bearer {{access_token}}

[Asserts]
status == 200
body.code == 0
```

---

**✅ Member Module Complete - 10 Steps**
