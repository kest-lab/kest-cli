# Kest RBAC Verification Scenario

This document defines a full-link test for the Kest RBAC implementation. It uses both `kest` code blocks and standard Markdown.

## 1. Project Creation
First, User 100 creates a new project. They should automatically become the **Owner**.

```kest
POST /api/v1/projects
X-User-ID: 100
Content-Type: application/json

{
  "name": "Markdown RBAC Project",
  "description": "Created from a markdown doc"
}

[Captures]
projectID: data.id

[Asserts]
status == 201
body.name == "Markdown RBAC Project"
```

## 2. Verify Membership
Let's check if User 100 is indeed the owner.

```kest
GET /api/v1/projects/{{projectID}}/members
X-User-ID: 100

[Asserts]
status == 200
body[0].user_id == 100
body[0].role == "owner"
```

## 3. Add a Read-Only Member
The owner (100) invites User 101 with a `read` role.

```kest
POST /api/v1/projects/{{projectID}}/members
X-User-ID: 100
Content-Type: application/json

{
  "user_id": 101,
  "role": "read"
}

[Asserts]
status == 201
```

## 4. Test Permission Enforcement
User 101 (Read) tries to create an environment. This should be **denied**.

```kest
POST /api/v1/projects/{{projectID}}/environments
X-User-ID: 101
Content-Type: application/json

{
  "name": "Illegal Env",
  "base_url": "https://hacking.com"
}

[Asserts]
status == 403
```

## 5. Owner Success
User 100 (Owner) creates the environment successfully.

```kest
POST /api/v1/projects/{{projectID}}/environments
X-User-ID: 100
Content-Type: application/json

{
  "project_id": {{projectID}},
  "name": "Markdown Env",
  "base_url": "https://api.example.com"
}

[Asserts]
status == 201
```

## 6. Project Isolation
An unrelated user (102) tries to access the project. This should be **denied**.

```kest
GET /api/v1/projects/{{projectID}}
X-User-ID: 102

[Asserts]
status == 403
```
