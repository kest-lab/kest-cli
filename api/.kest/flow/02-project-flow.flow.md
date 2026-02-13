# 🏗️ Project Module Flow

## 1. Create a New Project
```kest
POST /api/v1/projects
{
  "name": "Kest Test Project",
  "description": "A project created by Kest automated flow"
}

[Captures]
projectId: data.id

[Asserts]
status == 201
body.data.name == "Kest Test Project"
duration < 500ms
body.data.id exists
```

## 2. Verify Project Persistence (Side-Effect Check)
```kest
GET /api/v1/projects/{{projectId}}

[Asserts]
status == 200
body.data.id == "{{projectId}}"
body.data.name == "Kest Test Project"
body.data.description exists
```

## 3. List Projects & Verify Inclusion
```kest
GET /api/v1/projects

[Asserts]
status == 200
# Check if our created project exists in the list (if we had a way to search in array)
# For now we just verify the list structure
body.data exists
```

## 4. Update Project
```kest
PATCH /api/v1/projects/{{projectId}}
{
  "name": "Updated Kest Project",
  "description": "New description"
}

[Asserts]
# Some APIs return 200 or 204 on update
status >= 200
status < 300
```

## 5. Verify Update (Consistency Check)
```kest
GET /api/v1/projects/{{projectId}}

[Asserts]
status == 200
body.data.name == "Updated Kest Project"
body.data.description == "New description"
```

## 6. Delete Project
```kest
DELETE /api/v1/projects/{{projectId}}

[Asserts]
status == 200
```

## 7. Verify Cleanup (Negative Test)
```kest
GET /api/v1/projects/{{projectId}}

[Asserts]
status == 404
```
