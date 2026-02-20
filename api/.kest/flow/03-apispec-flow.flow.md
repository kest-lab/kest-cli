# 📄 API Specification Flow (Complete)

## 1. Setup: Create Parent Project
```kest
POST /v1/projects
{
  "name": "Project {{$randomInt}}",
  "description": "Project for testing API specs"
}

[Captures]
pid: data.id

[Asserts]
status == 201
```

## 2. Setup: Create Category
```kest
POST /v1/projects/{{pid}}/categories
{
  "name": "User APIs",
  "description": "User related endpoints"
}

[Captures]
cid: data.id

[Asserts]
status == 201
body.data.name == "User APIs"
```

## 3. Create API Spec
```kest
POST /v1/projects/{{pid}}/api-specs
{
  "summary": "Get User Profile",
  "method": "GET",
  "path": "/v1/users/profile",
  "category_id": {{cid}},
  "version": "1.0.0",
  "description": "Retrieves the authenticated user's profile"
}

[Captures]
sid: data.id

[Asserts]
status == 201
body.data.summary == "Get User Profile"
body.data.path == "/v1/users/profile"
```

## 4. Verify Spec Details
```kest
GET /v1/projects/{{pid}}/api-specs/{{sid}}

[Asserts]
status == 200
body.data.id == "{{sid}}"
body.data.method == "GET"
body.data.category_id == "{{cid}}"
duration < 200ms
```

## 5. Update Spec
```kest
PATCH /v1/projects/{{pid}}/api-specs/{{sid}}
{
  "description": "Updated API Spec Description"
}

[Asserts]
status == 200
body.data.description == "Updated API Spec Description"
```

## 6. Delete Spec
```kest
DELETE /v1/projects/{{pid}}/api-specs/{{sid}}

[Asserts]
status == 204
```

## 7. Cleanup Category
```kest
DELETE /v1/projects/{{pid}}/categories/{{cid}}

[Asserts]
status == 204
```

## 8. Cleanup Project
```kest
DELETE /v1/projects/{{pid}}

[Asserts]
status == 200
```
