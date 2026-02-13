# Common Code Review Issues

This document lists the most frequently found issues during code reviews and how to fix them.

---

## 🏗️ Architecture Issues

### ❌ Issue 1: Layer Violation

**Problem**:
```go
// handler.go
func (h *Handler) GetUser(c *gin.Context) {
    var user UserPO
    h.db.First(&user, id)  // ❌ Handler accessing DB directly
    c.JSON(200, user)
}
```

**Fix**:
```go
// handler.go
func (h *Handler) GetUser(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return
    }
    
    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.HandleError(c, "User not found", err)
        return
    }
    
    response.Success(c, ToResponse(user))
}
```

**Why**: Handlers should only handle HTTP concerns. Business logic and data access belong in service and repository layers.

---

### ❌ Issue 2: Using PO in Service Layer

**Problem**:
```go
// service.go
func (s *service) Create(ctx context.Context, req *CreateRequest) error {
    po := &UserPO{  // ❌ Service using PO
        Username: req.Username,
        Email:    req.Email,
    }
    return s.repo.Create(ctx, po)
}
```

**Fix**:
```go
// service.go
func (s *service) Create(ctx context.Context, req *CreateRequest) (*domain.User, error) {
    user := &domain.User{  // ✅ Service using domain entity
        Username: req.Username,
        Email:    req.Email,
    }
    
    err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

// repository.go
func (r *repository) Create(ctx context.Context, user *domain.User) error {
    po := ToUserPO(user)  // ✅ Repository handles PO conversion
    return r.db.WithContext(ctx).Create(po).Error
}
```

**Why**: Service layer should work with domain entities for business logic. Repository handles the persistence layer concerns.

---

## 🔒 Security Issues

### ❌ Issue 3: Plaintext Passwords

**Problem**:
```go
// service.go
func (s *service) Create(ctx context.Context, req *CreateUserRequest) error {
    user := &domain.User{
        Password: req.Password,  // ❌ Storing plaintext!
    }
    return s.repo.Create(ctx, user)
}
```

**Fix**:
```go
// service.go
func (s *service) Create(ctx context.Context, req *CreateUserRequest) (*domain.User, error) {
    // Hash password before storing
    hash, err := crypto.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    user := &domain.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hash,  // ✅ Store hash
    }
    
    return user, s.repo.Create(ctx, user)
}

// domain/user.go
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"`  // ✅ Never expose in JSON
}
```

**Why**: Never store passwords in plaintext. Always hash with bcrypt or similar.

---

### ❌ Issue 4: SQL Injection

**Problem**:
```go
// repository.go
func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var po UserPO
    // ❌ SQL injection vulnerability!
    r.db.Raw("SELECT * FROM users WHERE email = '" + email + "'").Scan(&po)
    return FromUserPO(&po), nil
}
```

**Fix**:
```go
// repository.go
func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var po UserPO
    // ✅ Use parameterized query
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error
    if err != nil {
        return nil, err
    }
    return FromUserPO(&po), nil
}
```

**Why**: String concatenation in SQL queries allows SQL injection attacks. Always use parameterized queries.

---

## ❗ Error Handling Issues

### ❌ Issue 5: Ignoring Errors

**Problem**:
```go
func (s *service) Create(ctx context.Context, req *CreateRequest) (*domain.User, error) {
    user := &domain.User{Username: req.Username}
    s.repo.Create(ctx, user)  // ❌ Error ignored!
    return user, nil
}
```

**Fix**:
```go
func (s *service) Create(ctx context.Context, req *CreateRequest) (*domain.User, error) {
    user := &domain.User{Username: req.Username}
    
    err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

**Why**: All errors must be handled. Ignoring errors hides failures and makes debugging impossible.

---

### ❌ Issue 6: No Error Context

**Problem**:
```go
func (s *service) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err  // ❌ No context about what failed
    }
    return user, nil
}
```

**Fix**:
```go
func (s *service) GetByID(ctx context.Context, id uint) (*domain.User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)  // ✅ Clear context
    }
    return user, nil
}
```

**Why**: Error wrapping provides context for debugging. Without it, you don't know where or why the error occurred.

---

## 📊 Performance Issues

### ❌ Issue 7: N+1 Query Problem

**Problem**:
```go
func (s *service) ListWithProfiles(ctx context.Context) ([]*domain.User, error) {
    users, _ := s.repo.List(ctx)
    
    // ❌ N+1 queries: 1 for users, N for profiles
    for _, user := range users {
        profile, _ := s.profileRepo.GetByUserID(ctx, user.ID)
        user.Profile = profile
    }
    
    return users, nil
}
```

**Fix**:
```go
// repository.go
func (r *repository) ListWithProfiles(ctx context.Context) ([]*domain.User, error) {
    var pos []UserPO
    // ✅ Single query with JOIN or Preload
    err := r.db.WithContext(ctx).Preload("Profile").Find(&pos).Error
    if err != nil {
        return nil, err
    }
    
    return ToDomainUsers(pos), nil
}
```

**Why**: Loading related data in loops causes N+1 queries. Use JOINs or eager loading instead.

---

### ❌ Issue 8: Missing Pagination

**Problem**:
```go
// handler.go
func (h *Handler) List(c *gin.Context) {
    var users []User
    h.db.Find(&users)  // ❌ Could return millions of rows!
    response.Success(c, users)
}
```

**Fix**:
```go
// handler.go
func (h *Handler) List(c *gin.Context) {
    users, paginator, err := pagination.PaginateFromContext[*domain.User](c, h.db)
    if err != nil {
        response.HandleError(c, "Failed to list users", err)
        return
    }
    response.Success(c, paginator)  // ✅ Paginated response
}
```

**Why**: Loading all records is slow and uses excessive memory. Always paginate list endpoints.

---

## 📝 Code Quality Issues

### ❌ Issue 9: Poor Naming

**Problem**:
```go
type User struct {}  // ❌ Model missing PO suffix
type CreateDTO struct {}  // ❌ Should be CreateRequest
func Get(id uint) {}  // ❌ Vague: Get what?
```

**Fix**:
```go
type UserPO struct {}  // ✅ Clear: persistence object
type CreateUserRequest struct {}  // ✅ Clear: request DTO
func GetByID(id uint) (*domain.User, error) {}  // ✅ Specific
```

**Why**: Clear naming makes code self-documenting. Follow conventions consistently.

---

### ❌ Issue 10: Missing Validation

**Problem**:
```go
// dto.go
type CreateUserRequest struct {
    Email    string `json:"email"`  // ❌ No validation!
    Password string `json:"password"`
}
```

**Fix**:
```go
// dto.go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`  // ✅ Validated
    Password string `json:"password" binding:"required,min=8,max=72"`
    Username string `json:"username" binding:"required,min=3,max=50"`
}

// handler.go
func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    if !handler.BindJSON(c, &req) {
        return  // ✅ Validation errors automatically returned
    }
    // ...
}
```

**Why**: Input validation prevents invalid data and provides clear error messages.

---

## 🧪 Testing Issues

### ❌ Issue 11: No Tests

**Problem**:
```go
// service.go exists
// service_test.go doesn't exist ❌
```

**Fix**:
```go
// service_test.go
func TestService_Create_Success(t *testing.T) {
    // Setup
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    req := &CreateUserRequest{
        Email: "test@example.com",
        Username: "testuser",
    }
    
    // Expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Execute
    user, err := service.Create(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockRepo.AssertExpectations(t)
}
```

**Why**: Untested code is broken code. Write tests for all new functionality.

---

### ❌ Issue 12: Testing Implementation, Not Interface

**Problem**:
```go
func TestService_Create(t *testing.T) {
    db := setupTestDB()  // ❌ Testing with real DB
    repo := NewRepository(db)
    service := NewService(repo)
    
    service.Create(context.Background(), &req)
    // Asserting DB state...
}
```

**Fix**:
```go
func TestService_Create(t *testing.T) {
    mockRepo := new(MockRepository)  // ✅ Mock dependency
    service := NewService(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    user, err := service.Create(context.Background(), &req)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)  // ✅ Test behavior, not implementation
}
```

**Why**: Unit tests should test behavior, not implementation details. Use mocks for dependencies.

---

## 📚 Documentation Issues

### ❌ Issue 13: Missing Swagger Docs

**Problem**:
```go
// handler.go
func (h *Handler) Create(c *gin.Context) {  // ❌ No Swagger comments
    // ...
}
```

**Fix**:
```go
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/users [post]
func (h *Handler) Create(c *gin.Context) {
    // ...
}
```

**Why**: API documentation helps consumers understand how to use endpoints.

---

## 🎯 Quick Fix Checklist

Before requesting re-review, verify:

- [ ] All errors handled and wrapped
- [ ] Passwords hashed, not plaintext
- [ ] SQL queries use parameters, not concatenation
- [ ] Layers properly separated (Handler→Service→Repository)
- [ ] Service uses domain entities, not POs
- [ ] List endpoints use pagination
- [ ] Input validation with `binding` tags
- [ ] Tests written and passing
- [ ] Swagger documentation added
- [ ] Naming follows conventions

---

**See Also**:
- [coding-standards](../coding-standards/) - Full coding standards
- [api-development](../api-development/) - API best practices
- [good-review-example](./good-review-example.md) - Example of good review
