---
name: api-development
description: ZGO API development standards including pagination, error handling, and RESTful design
version: 1.0.0
category: development
tags: [api, rest, pagination, errors, standards]
author: ZGO Team
updated: 2026-01-24
---

# API Development Standards

## 📋 Purpose

This skill provides comprehensive API development standards for the ZGO Go backend project, ensuring consistent and high-quality REST APIs across all modules.

## 🎯 When to Use

- Creating new API endpoints
- Reviewing API code
- Implementing list/collection endpoints
- Handling errors in handlers
- Designing RESTful resources

## ⚙️ Prerequisites

- [ ] Understanding of Go and Gin framework
- [ ] Familiarity with RESTful principles
- [ ] Knowledge of ZGO's module structure (8-file standard)

## 📚 Core Standards

### 1. Pagination (MANDATORY)

**Rule**: All list/collection endpoints **MUST** implement pagination.

#### Why Pagination is Required

1. **Performance**: Prevents loading thousands of records into memory
2. **User Experience**: Faster response times
3. **Scalability**: Database and network efficiency
4. **API Consistency**: Uniform behavior across all list endpoints

#### Implementation

```go
package user

import (
    "github.com/gin-gonic/gin"
    "github.com/zgiai/zgo/pkg/handler"
    "github.com/zgiai/zgo/pkg/response"
    "github.com/zgiai/zgo/pkg/pagination"
    "github.com/zgiai/zgo/internal/domain"
)

// ✅ CORRECT - With pagination
func (h *Handler) List(c *gin.Context) {
    // One-liner pagination (recommended)
    users, paginator, err := pagination.PaginateFromContext[*domain.User](c, h.db)
    if err != nil {
        response.HandleError(c, "Failed to fetch users", err)
        return
    }
    
    // Auto-detects pagination and includes meta + links
    response.Success(c, paginator)
}

// ❌ WRONG - Without pagination
func (h *Handler) List(c *gin.Context) {
    var users []User
    h.db.Find(&users)  // Loads ALL records!
    response.Success(c, users)
}
```

#### Pagination with Filters

```go
func (h *Handler) List(c *gin.Context) {
    // Apply filters before pagination
    query := h.db.Model(&UserPO{})
    
    // Filter by status
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Filter by search
    if search := c.Query("search"); search != "" {
        query = query.Where("username LIKE ? OR email LIKE ?", 
            "%"+search+"%", "%"+search+"%")
    }
    
    // Paginate filtered results
    users, paginator, err := pagination.PaginateFromContext[*domain.User](c, query)
    if err != nil {
        response.HandleError(c, "Failed to fetch users", err)
        return
    }
    
    response.Success(c, paginator)
}
```

#### Pagination Standards

| Setting | Value | Description |
|---------|-------|-------------|
| **Default Page Size** | 20 | Records per page if not specified |
| **Maximum Page Size** | 100 | Upper limit to prevent abuse |
| **Query Parameter** | `page` | Page number (1-indexed) |
| **Query Parameter** | `page_size` | Records per page |

**Example Request**:
```
GET /api/users?page=2&page_size=20&status=active&search=john
```

**Example Response**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 21,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    }
  ],
  "meta": {
    "total": 150,
    "page": 2,
    "page_size": 20,
    "total_pages": 8
  },
  "links": {
    "first": "/api/users?page=1&page_size=20",
    "last": "/api/users?page=8&page_size=20",
    "prev": "/api/users?page=1&page_size=20",
    "next": "/api/users?page=3&page_size=20"
  }
}
```

---

### 2. Unified Error Responses (MANDATORY)

**Rule**: All error responses **MUST** use `pkg/response` package functions.

#### Why Unified Errors

1. **Consistency**: Same format across all APIs
2. **Automatic Mapping**: Framework handles status code logic
3. **User-Friendly**: Clean error messages for clients
4. **Debugging**: Structured errors for logging

#### Available Error Functions

```go
import "github.com/zgiai/zgo/pkg/response"

// Client Errors (4xx)
response.BadRequest(c, "message", err)           // 400
response.Unauthorized(c)                         // 401  
response.Forbidden(c)                            // 403
response.NotFound(c, "message", err)             // 404
response.Conflict(c, "message", err)             // 409
response.UnprocessableEntity(c, "message", err)  // 422
response.ValidationFailed(c, fieldErrors)        // 422 with fields

// Server Errors (5xx)
response.InternalServerError(c, "message", err)  // 500
response.ServiceUnavailable(c)                   // 503

// Auto-detection (recommended)
response.HandleError(c, "message", err)          // Auto-maps based on error
```

#### Error Handling Pattern

```go
func (h *Handler) Get(c *gin.Context) {
    // Step 1: Parse and validate input
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return  // ✅ 400 already sent by ParseID
    }
    
    // Step 2: Call service
    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        // ✅ Auto-maps error type to status code
        response.HandleError(c, "User not found", err)
        return
    }
    
    // Step 3: Return success
    response.Success(c, ToResponse(user))
}
```

#### Error Response Formats

**Simple Error**:
```json
{
  "code": 404,
  "message": "User not found",
  "error": "record not found"
}
```

**Validation Error**:
```json
{
  "code": 422,
  "message": "Validation failed",
  "errors": {
    "email": ["The email field is required", "Email format is invalid"],
    "password": ["The password must be at least 8 characters"]
  }
}
```

#### Custom Error Types

Define module-specific errors in `service.go`:

```go
package user

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrDuplicateEmail    = errors.New("email already exists")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrAccountSuspended  = errors.New("account is suspended")
)

// Service layer
func (s *service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    user, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrUserNotFound
    }
    
    if user.Status == "suspended" {
        return nil, ErrAccountSuspended
    }
    
    return user, nil
}

// Handler layer
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if !handler.BindJSON(c, &req) {
        return
    }
    
    user, err := h.service.GetByEmail(c.Request.Context(), req.Email)
    if err != nil {
        // Maps custom errors to appropriate status codes
        response.HandleError(c, "Login failed", err)
        return
    }
    
    // ... password verification
}
```

---

### 3. Success Responses (MANDATORY)

**Rule**: Use `pkg/response` for all successful responses.

```go
// 200 OK - Standard success
response.Success(c, data)
response.Success(c, paginator)  // Auto-detects pagination

// 201 Created - After creating a resource
response.Created(c, newUser)

// 204 No Content - After deletion
response.NoContent(c)

// 202 Accepted - For async operations
response.Accepted(c, task)
```

#### Success Response Format

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

---

### 4. HTTP Method Standards (MANDATORY)

| Method | Usage | Response | Body | Example |
|--------|-------|----------|------|---------|
| **GET** | Retrieve resource(s) | 200 + data | No | `GET /api/users/:id` |
| **POST** | Create new resource | 201 + data | Yes | `POST /api/users` |
| **PATCH** | Partial update | 200 + data | Yes | `PATCH /api/users/:id` |
| **PUT** | Full replacement | 200 + data | Yes | `PUT /api/users/:id` |
| **DELETE** | Remove resource | 204 (no body) | No | `DELETE /api/users/:id` |

#### Implementation Examples

```go
// GET - Retrieve single resource
func (h *Handler) Get(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok { return }
    
    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.HandleError(c, "User not found", err)
        return
    }
    
    response.Success(c, ToResponse(user))  // 200
}

// POST - Create new resource
func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    if !handler.BindJSON(c, &req) { return }
    
    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        response.HandleError(c, "Failed to create user", err)
        return
    }
    
    response.Created(c, ToResponse(user))  // 201
}

// PATCH - Partial update
func (h *Handler) Update(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok { return }
    
    var req UpdateUserRequest
    if !handler.BindJSON(c, &req) { return }
    
    user, err := h.service.Update(c.Request.Context(), id, &req)
    if err != nil {
        response.HandleError(c, "Failed to update user", err)
        return
    }
    
    response.Success(c, ToResponse(user))  // 200
}

// DELETE - Remove resource
func (h *Handler) Delete(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok { return }
    
    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        response.HandleError(c, "Failed to delete user", err)
        return
    }
    
    response.NoContent(c)  // 204 No Content
}
```

---

### 5. RESTful URL Naming (MANDATORY)

#### Principles

1. **Use nouns, not verbs**: Resources, not actions
2. **Use plural names**: `/users`, not `/user`
3. **Use HTTP methods**: For actions (GET, POST, etc.)
4. **Use nested resources**: For relationships
5. **Use kebab-case**: For multi-word resources (if needed)

#### Correct Examples

```go
// ✅ CORRECT - RESTful resource design
GET    /api/users                 // List all users
POST   /api/users                 // Create new user
GET    /api/users/:id             // Get specific user
PATCH  /api/users/:id             // Update user
DELETE /api/users/:id             // Delete user

// Nested resources
GET    /api/users/:id/posts       // Get user's posts
POST   /api/users/:id/posts       // Create post for user
GET    /api/posts/:id/comments    // Get post's comments

// With query parameters
GET    /api/users?status=active&role=admin
GET    /api/posts?author_id=123&published=true
```

#### Incorrect Examples

```go
// ❌ WRONG - Using verbs
GET    /api/getUsers
POST   /api/createUser
POST   /api/users/delete/:id
GET    /api/fetchUserById/:id

// ❌ WRONG - Using singular
GET    /api/user
POST   /api/user/:id/post

// ❌ WRONG - Non-RESTful patterns
GET    /api/user_list
POST   /api/user-create
GET    /api/get_user_by_id/:id
```

---

### 6. Request Validation (MANDATORY)

#### DTO Validation

Always use `binding` tags for input validation:

```go
// Create Request
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=72"`
    Age      int    `json:"age" binding:"omitempty,gte=0,lte=150"`
    Role     string `json:"role" binding:"omitempty,oneof=admin user guest"`
}

// Update Request (optional fields use pointers)
type UpdateUserRequest struct {
    Username *string `json:"username" binding:"omitempty,min=3,max=50"`
    Email    *string `json:"email" binding:"omitempty,email"`
    Bio      *string `json:"bio" binding:"omitempty,max=500"`
}
```

#### Common Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field is required | `binding:"required"` |
| `omitempty` | Optional field | `binding:"omitempty,email"` |
| `min=N,max=N` | String length or number range | `binding:"min=3,max=50"` |
| `gte=N,lte=N` | Number greater/less than | `binding:"gte=0,lte=150"` |
| `email` | Valid email format | `binding:"email"` |
| `url` | Valid URL format | `binding:"url"` |
| `oneof=a b c` | Enum values | `binding:"oneof=active inactive"` |
| `uuid` | Valid UUID format | `binding:"uuid"` |
| `datetime` | Valid datetime | `binding:"datetime=2006-01-02"` |

#### Validation in Handler

```go
func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    
    // ✅ Automatic validation
    if !handler.BindJSON(c, &req) {
        return  // 422 with field errors already sent
    }
    
    // Continue with business logic
    user, err := h.service.Create(c.Request.Context(), &req)
    // ...
}
```

---

## 🚀 Complete CRUD Example

See [`examples/complete-crud-handler.go`](./examples/complete-crud-handler.go) for a full implementation.

---

## ✅ Verification Checklist

Use this checklist before submitting API code:

### Pagination
- [ ] All list endpoints use `pagination.PaginateFromContext[T]()`
- [ ] Response includes `meta` (total, page, pageSize, totalPages)
- [ ] Response includes `links` (first, last, prev, next)
- [ ] Default page size is 20, max is 100
- [ ] Query parameters are `page` and `page_size`

### Error Handling
- [ ] All errors use `response.*` functions
- [ ] No manual `c.JSON(statusCode, ...)` for errors
- [ ] Error messages are user-friendly
- [ ] Custom errors defined in service layer
- [ ] No sensitive information in error responses

### Success Responses
- [ ] Use `response.Success()` for 200 OK
- [ ] Use `response.Created()` for 201 Created
- [ ] Use `response.NoContent()` for 204 No Content
- [ ] Use `response.Accepted()` for 202 Accepted

### HTTP Methods
- [ ] GET for retrieval (200)
- [ ] POST for creation (201)
- [ ] PATCH for partial updates (200)
- [ ] PUT for full replacement (200)
- [ ] DELETE for removal (204)

### RESTful URLs
- [ ] URLs use plural nouns (`/users`, not `/user`)
- [ ] No verbs in URLs (use HTTP methods)
- [ ] Nested resources follow `/parent/:id/child` pattern
- [ ] Query parameters for filters and sorting

### Validation
- [ ] All DTOs have `binding` tags
- [ ] Using `handler.BindJSON()` for auto-validation
- [ ] Update DTOs use pointers for optional fields
- [ ] Validation errors return 422 with field details

---

## 🔧 Automated Validation

Run the API standards validation script:

```bash
.agent/skills/api-development/scripts/validate-api.sh user
```

This checks:
- ✅ Pagination usage in list endpoints
- ✅ `response.*` usage (no manual status codes)
- ✅ RESTful URL patterns
- ✅ HTTP method correctness
- ✅ Validation tags in DTOs

---

## 🚫 Common Mistakes

### Mistake 1: No Pagination

```go
// ❌ WRONG - Returns all records
func (h *Handler) List(c *gin.Context) {
    var users []User
    h.db.Find(&users)
    response.Success(c, users)
}

// ✅ CORRECT - With pagination
func (h *Handler) List(c *gin.Context) {
    users, paginator, _ := pagination.PaginateFromContext[*domain.User](c, h.db)
    response.Success(c, paginator)
}
```

### Mistake 2: Manual Status Codes

```go
// ❌ WRONG
c.JSON(404, gin.H{"error": "not found"})
c.AbortWithStatusJSON(400, map[string]any{"error": "bad request"})

// ✅ CORRECT
response.NotFound(c, "User not found", err)
response.BadRequest(c, "Invalid input", err)
```

### Mistake 3: Verbs in URLs

```go
// ❌ WRONG
GET /api/getUsers
POST /api/createUser
POST /api/users/delete/:id

// ✅ CORRECT
GET /api/users
POST /api/users
DELETE /api/users/:id
```

### Mistake 4: Wrong HTTP Methods

```go
// ❌ WRONG - GET for updates
GET /api/users/:id/update

// ✅ CORRECT - PATCH for updates
PATCH /api/users/:id
```

### Mistake 5: Returning 200 on DELETE

```go
// ❌ WRONG - DELETE returns data
func (h *Handler) Delete(c *gin.Context) {
    h.service.Delete(c.Request.Context(), id)
    response.Success(c, gin.H{"message": "deleted"})  // 200
}

// ✅ CORRECT - DELETE returns 204
func (h *Handler) Delete(c *gin.Context) {
    h.service.Delete(c.Request.Context(), id)
    response.NoContent(c)  // 204
}
```

---

## 📚 Additional Resources

- [RESTful API Design Best Practices](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
- [HTTP Status Codes](https://httpstatuses.com/)
- [Pagination Best Practices](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)

---

## 🔗 Related Skills

- [`module-creation`](../module-creation/): For creating new modules with handlers
- [`coding-standards`](../coding-standards/): For general code quality standards
- [`testing-strategy`](../testing-strategy/): For API testing patterns

---

## 📝 Quick Reference

```go
// Pagination (MUST)
users, paginator, _ := pagination.PaginateFromContext[T](c, db)
response.Success(c, paginator)

// Errors (MUST)
response.HandleError(c, "message", err)    // Auto-map
response.NotFound(c, "message", err)       // 404
response.BadRequest(c, "message", err)     // 400

// Success (MUST)
response.Success(c, data)                  // 200
response.Created(c, resource)              // 201
response.NoContent(c)                      // 204

// Validation (MUST)
if !handler.BindJSON(c, &req) { return }
```

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team
