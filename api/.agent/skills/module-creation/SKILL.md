---
name: module-creation
description: Complete workflow for creating a new DDD module following ZGO's 8-file standard
version: 1.0.0
category: development
tags: [module, ddd, scaffolding, architecture]
author: ZGO Team
updated: 2026-01-24
---

# Module Creation Skill

## 📋 Purpose

This skill guides you through creating a standardized DDD (Domain-Driven Design) module in the ZGO framework, ensuring adherence to the **8-file structure** and architectural best practices.

ZGO uses a layered architecture where each business module is self-contained with clear separation of concerns:
- **Model Layer** (Database entities)
- **DTO Layer** (Data Transfer Objects + Mappers)
- **Repository Layer** (Data access)
- **Service Layer** (Business logic)
- **Handler Layer** (HTTP controllers)
- **Routes** (Endpoint registration)
- **Provider** (Dependency injection)
- **Tests** (Unit/integration tests)

## 🎯 When to Use

Use this skill when:
- Creating a new business module (e.g., User, Blog, Product, Order)
- Scaffolding a complete CRUD feature
- Ensuring consistency with ZGO's DDD architecture
- Teaching or onboarding team members to ZGO patterns

## ⚙️ Prerequisites

- [ ] Go 1.21+ installed
- [ ] ZGO project cloned and set up
- [ ] Wire tool installed: `go install github.com/google/wire/cmd/wire@latest`
- [ ] Basic understanding of DDD concepts
- [ ] Database connection configured

## 🚀 Workflow Steps

### Step 1: Define Module Scope

Before writing any code, clearly define:

**Business Requirements**:
- Module name (PascalCase): `Blog`, `UserProfile`, `Product`
- Core domain entities and their relationships
- Required operations (CRUD, custom actions)

**Technical Specification**:
```markdown
Module: Blog
Domain Entity: BlogPost
Database Table: blog_posts
API Endpoints:
  - GET    /api/blogs          → List all posts (paginated)
  - POST   /api/blogs          → Create new post
  - GET    /api/blogs/:id      → Get post by ID
  - PATCH  /api/blogs/:id      → Update post
  - DELETE /api/blogs/:id      → Delete post
  
Fields:
  - id (uint, primary key)
  - title (string, required, max 255)
  - content (text)
  - author_id (uint, foreign key)
  - status (enum: draft/published)
  - created_at, updated_at, deleted_at (GORM standard)
```

### Step 2: Create Module Directory

```bash
# Navigate to modules directory
cd internal/modules

# Create module folder (lowercase, singular)
mkdir blog
cd blog
```

**Expected structure**:
```
internal/modules/blog/
├── (8 files will be created below)
```

### Step 3: Create Database Entity (model.go)

**Purpose**: Define the database table structure using GORM.

**File**: `model.go`

```go
package blog

import (
    "time"
    "gorm.io/gorm"
)

// BlogPostPO is the persistent object for blog posts
// Naming: {Entity}PO (Persistent Object)
type BlogPostPO struct {
    ID        uint           `gorm:"primaryKey"`
    Title     string         `gorm:"size:255;not null;index"`
    Content   string         `gorm:"type:text"`
    AuthorID  uint           `gorm:"index;not null"`
    Status    string         `gorm:"size:20;default:'draft';index"` // draft, published
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the default table name
func (BlogPostPO) TableName() string {
    return "blog_posts"
}
```

**Key points**:
- Suffix `PO` for database models
- Use GORM tags for constraints
- Always include soft delete (`DeletedAt`)
- Use `TableName()` for explicit table names

### Step 4: Create Domain Entity and DTOs (dto.go)

**Purpose**: Define domain objects, request/response DTOs, and mapper functions.

**File**: `dto.go`

```go
package blog

import (
    "time"
    "github.com/zgiai/zgo/internal/domain"
)

// ========== Domain Entity (lives in internal/domain) ==========
// Note: In practice, domain.BlogPost should be in internal/domain/blog_post.go
// For this example, we define it here, but prefer domain package

// ========== Request DTOs ==========

type CreateBlogPostRequest struct {
    Title   string `json:"title" binding:"required,max=255"`
    Content string `json:"content" binding:"required"`
    Status  string `json:"status" binding:"omitempty,oneof=draft published"`
}

type UpdateBlogPostRequest struct {
    Title   *string `json:"title" binding:"omitempty,max=255"`
    Content *string `json:"content" binding:"omitempty"`
    Status  *string `json:"status" binding:"omitempty,oneof=draft published"`
}

// ========== Response DTOs ==========

type BlogPostResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    AuthorID  uint      `json:"author_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ========== Mapper Functions ==========

// ToBlogPostPO converts domain entity to persistent object
func ToBlogPostPO(post *domain.BlogPost) *BlogPostPO {
    return &BlogPostPO{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        AuthorID:  post.AuthorID,
        Status:    post.Status,
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }
}

// FromBlogPostPO converts persistent object to domain entity
func FromBlogPostPO(po *BlogPostPO) *domain.BlogPost {
    if po == nil {
        return nil
    }
    return &domain.BlogPost{
        ID:        po.ID,
        Title:     po.Title,
        Content:   po.Content,
        AuthorID:  po.AuthorID,
        Status:    po.Status,
        CreatedAt: po.CreatedAt,
        UpdatedAt: po.UpdatedAt,
    }
}

// ToResponse converts domain entity to response DTO
func ToResponse(post *domain.BlogPost) *BlogPostResponse {
    if post == nil {
        return nil
    }
    return &BlogPostResponse{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        AuthorID:  post.AuthorID,
        Status:    post.Status,
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }
}
```

**Data flow**:
```
Handler (DTO) → Service (domain.BlogPost) → Repository (BlogPostPO)
                                           ← Repository (domain.BlogPost)
                ← Service (domain.BlogPost)
Handler (DTO) ←
```

### Step 5: Create Repository Layer (repository.go)

**Purpose**: Handle all database operations, return domain entities.

**File**: `repository.go`

```go
package blog

import (
    "context"
    "github.com/zgiai/zgo/internal/domain"
    "gorm.io/gorm"
)

// Repository defines blog post data access interface
type Repository interface {
    Create(ctx context.Context, post *domain.BlogPost) error
    GetByID(ctx context.Context, id uint) (*domain.BlogPost, error)
    Update(ctx context.Context, post *domain.BlogPost) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, pageSize int) ([]*domain.BlogPost, int64, error)
}

// repository is the private implementation
type repository struct {
    db *gorm.DB
}

// NewRepository creates a new blog repository
func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, post *domain.BlogPost) error {
    po := ToBlogPostPO(post)
    if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
        return err
    }
    *post = *FromBlogPostPO(po) // Update with generated ID
    return nil
}

func (r *repository) GetByID(ctx context.Context, id uint) (*domain.BlogPost, error) {
    var po BlogPostPO
    if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
        return nil, err
    }
    return FromBlogPostPO(&po), nil
}

func (r *repository) Update(ctx context.Context, post *domain.BlogPost) error {
    po := ToBlogPostPO(post)
    return r.db.WithContext(ctx).Save(po).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&BlogPostPO{}, id).Error
}

func (r *repository) List(ctx context.Context, page, pageSize int) ([]*domain.BlogPost, int64, error) {
    var posts []BlogPostPO
    var total int64
    
    offset := (page - 1) * pageSize
    
    if err := r.db.WithContext(ctx).Model(&BlogPostPO{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    if err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(pageSize).
        Find(&posts).Error; err != nil {
        return nil, 0, err
    }
    
    result := make([]*domain.BlogPost, len(posts))
    for i, po := range posts {
        result[i] = FromBlogPostPO(&po)
    }
    
    return result, total, nil
}
```

**Key patterns**:
- Interface-based design
- Private struct implementation
- Constructor returns interface
- Always use `context.Context`
- Convert PO ↔ Domain at repository boundary

### Step 6: Create Service Layer (service.go)

**Purpose**: Implement business logic using domain entities.

**File**: `service.go`

```go
package blog

import (
    "context"
    "errors"
    "github.com/zgiai/zgo/internal/domain"
)

var (
    ErrBlogPostNotFound     = errors.New("blog post not found")
    ErrInvalidBlogPostData  = errors.New("invalid blog post data")
    ErrUnauthorized         = errors.New("unauthorized operation")
)

// Service defines blog post business logic interface
type Service interface {
    Create(ctx context.Context, req *CreateBlogPostRequest, authorID uint) (*domain.BlogPost, error)
    GetByID(ctx context.Context, id uint) (*domain.BlogPost, error)
    Update(ctx context.Context, id uint, req *UpdateBlogPostRequest) (*domain.BlogPost, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, pageSize int) ([]*domain.BlogPost, int64, error)
}

// service is the private implementation
type service struct {
    repo Repository
}

// NewService creates a new blog service
func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateBlogPostRequest, authorID uint) (*domain.BlogPost, error) {
    // Business validation
    if req.Title == "" {
        return nil, ErrInvalidBlogPostData
    }
    
    // Set default status
    status := req.Status
    if status == "" {
        status = "draft"
    }
    
    post := &domain.BlogPost{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: authorID,
        Status:   status,
    }
    
    if err := s.repo.Create(ctx, post); err != nil {
        return nil, err
    }
    
    return post, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*domain.BlogPost, error) {
    post, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrBlogPostNotFound
    }
    return post, nil
}

func (s *service) Update(ctx context.Context, id uint, req *UpdateBlogPostRequest) (*domain.BlogPost, error) {
    post, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrBlogPostNotFound
    }
    
    // Apply partial updates
    if req.Title != nil {
        post.Title = *req.Title
    }
    if req.Content != nil {
        post.Content = *req.Content
    }
    if req.Status != nil {
        post.Status = *req.Status
    }
    
    if err := s.repo.Update(ctx, post); err != nil {
        return nil, err
    }
    
    return post, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
    if _, err := s.repo.GetByID(ctx, id); err != nil {
        return ErrBlogPostNotFound
    }
    return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*domain.BlogPost, int64, error) {
    return s.repo.List(ctx, page, pageSize)
}
```

**Business logic examples**:
- Input validation
- Default values
- Authorization checks
- Business rules enforcement

### Step 7: Create HTTP Handlers (handler.go)

**Purpose**: Handle HTTP requests, use service layer, return responses.

**File**: `handler.go`

```go
package blog

import (
    "github.com/gin-gonic/gin"
    "github.com/zgiai/zgo/pkg/handler"
    "github.com/zgiai/zgo/pkg/response"
)

// Handler handles blog post HTTP requests
type Handler struct {
    service Service
}

// NewHandler creates a new blog handler
func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

// Create godoc
// @Summary Create blog post
// @Tags blogs
// @Accept json
// @Produce json
// @Param body body CreateBlogPostRequest true "Blog post data"
// @Success 201 {object} BlogPostResponse
// @Router /api/blogs [post]
func (h *Handler) Create(c *gin.Context) {
    // Get authenticated user
    userID, ok := handler.GetUserID(c)
    if !ok {
        return // 401 already sent
    }
    
    // Bind request
    var req CreateBlogPostRequest
    if !handler.BindJSON(c, &req) {
        return // 400 already sent
    }
    
    // Call service
    post, err := h.service.Create(c.Request.Context(), &req, userID)
    if err != nil {
        response.HandleError(c, "Failed to create blog post", err)
        return
    }
    
    // Return response
    response.Created(c, ToResponse(post))
}

// Get godoc
// @Summary Get blog post by ID
// @Tags blogs
// @Produce json
// @Param id path int true "Blog post ID"
// @Success 200 {object} BlogPostResponse
// @Router /api/blogs/{id} [get]
func (h *Handler) Get(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return
    }
    
    post, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.HandleError(c, "Blog post not found", err)
        return
    }
    
    response.Success(c, ToResponse(post))
}

// Update godoc
// @Summary Update blog post
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog post ID"
// @Param body body UpdateBlogPostRequest true "Updated data"
// @Success 200 {object} BlogPostResponse
// @Router /api/blogs/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return
    }
    
    var req UpdateBlogPostRequest
    if !handler.BindJSON(c, &req) {
        return
    }
    
    post, err := h.service.Update(c.Request.Context(), id, &req)
    if err != nil {
        response.HandleError(c, "Failed to update blog post", err)
        return
    }
    
    response.Success(c, ToResponse(post))
}

// Delete godoc
// @Summary Delete blog post
// @Tags blogs
// @Param id path int true "Blog post ID"
// @Success 204
// @Router /api/blogs/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return
    }
    
    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        response.HandleError(c, "Failed to delete blog post", err)
        return
    }
    
    response.NoContent(c)
}

// List godoc
// @Summary List blog posts (paginated)
// @Tags blogs
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.PaginatedResponse
// @Router /api/blogs [get]
func (h *Handler) List(c *gin.Context) {
    page := handler.QueryInt(c, "page", 1)
    perPage := handler.QueryInt(c, "per_page", 20)
    
    posts, total, err := h.service.List(c.Request.Context(), page, perPage)
    if err != nil {
        response.HandleError(c, "Failed to list blog posts", err)
        return
    }
    
    // Convert to responses
    responses := make([]*BlogPostResponse, len(posts))
    for i, post := range posts {
        responses[i] = ToResponse(post)
    }
    
    // Return paginated response
    // Note: Use pkg/pagination for automatic pagination
    response.Success(c, map[string]any{
        "data": responses,
        "meta": map[string]any{
            "total":       total,
            "page":        page,
            "per_page":    perPage,
            "total_pages": (total + int64(perPage) - 1) / int64(perPage),
        },
    })
}
```

**Handler best practices**:
- Use `pkg/handler` utilities
- Use `pkg/response` for consistent responses
- Add Swagger annotations
- Keep handlers thin (delegate to service)

### Step 8: Register Routes (routes.go)

**Purpose**: Define API endpoints and HTTP methods.

**File**: `routes.go`

```go
package blog

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes registers blog post routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
    blogs := router.Group("/blogs")
    {
        // Public routes
        blogs.GET("", handler.List)
        blogs.GET("/:id", handler.Get)
        
        // Protected routes (require authentication)
        blogs.POST("", authMiddleware, handler.Create)
        blogs.PATCH("/:id", authMiddleware, handler.Update)
        blogs.DELETE("/:id", authMiddleware, handler.Delete)
    }
}
```

**Route patterns**:
- Group related endpoints
- Apply middleware selectively
- Follow RESTful conventions

### Step 9: Wire Dependency Injection (provider.go)

**Purpose**: Configure Wire to auto-generate DI code.

**File**: `provider.go`

```go
package blog

import "github.com/google/wire"

// ProviderSet is the Wire provider set for blog module
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*repository)),
    NewService,
    wire.Bind(new(Service), new(*service)),
    NewHandler,
)
```

**Wire pattern**:
- Export `ProviderSet`
- Bind interfaces to implementations
- Constructors must match signatures

### Step 10: Create Unit Tests (service_test.go)

**Purpose**: Test business logic in isolation.

**File**: `service_test.go`

```go
package blog

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/zgiai/zgo/internal/domain"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, post *domain.BlogPost) error {
    args := m.Called(ctx, post)
    return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uint) (*domain.BlogPost, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.BlogPost), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, post *domain.BlogPost) error {
    args := m.Called(ctx, post)
    return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, page, pageSize int) ([]*domain.BlogPost, int64, error) {
    args := m.Called(ctx, page, pageSize)
    return args.Get(0).([]*domain.BlogPost), args.Get(1).(int64), args.Error(2)
}

// TestCreate tests the Create service method
func TestCreate(t *testing.T) {
    mockRepo := new(MockRepository)
    svc := NewService(mockRepo)
    ctx := context.Background()
    
    req := &CreateBlogPostRequest{
        Title:   "Test Post",
        Content: "Test content",
        Status:  "draft",
    }
    
    mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.BlogPost")).Return(nil)
    
    post, err := svc.Create(ctx, req, 1)
    
    assert.NoError(t, err)
    assert.NotNil(t, post)
    assert.Equal(t, "Test Post", post.Title)
    assert.Equal(t, uint(1), post.AuthorID)
    mockRepo.AssertExpectations(t)
}
```

**Testing best practices**:
- Mock repository dependencies
- Test business logic, not infrastructure
- Use table-driven tests for multiple cases

### Step 11: Integrate with Wire DI System

**Add module to main wire config**:

Edit `internal/wiring/wire.go`:

```go
// Add blog import
import (
    // ... existing imports ...
    "github.com/zgiai/zgo/internal/modules/blog"
)

// Add to initializeApplication function
func initializeApplication(cfg *config.Config, db *gorm.DB, ...) (*bootstrap.Application, func(), error) {
    wire.Build(
        // ... existing providers ...
        blog.ProviderSet,  // Add this line
        // ... rest of the providers ...
    )
    return nil, nil, nil
}
```

**Generate Wire code**:

```bash
cd internal/wiring
wire
```

Expected output:
```
wire: blog: wrote /path/to/zgo/internal/wiring/wire_gen.go
```

### Step 12: Register Routes in Application

Edit `routes/api.go` (or wherever routes are registered):

```go
import (
    // ... existing imports ...
    "github.com/zgiai/zgo/internal/modules/blog"
)

func RegisterAPIRoutes(app *bootstrap.Application) {
    api := app.Router.Group("/api")
    
    // ... existing routes ...
    
    // Blog routes
    blog.RegisterRoutes(api, app.BlogHandler, app.AuthMiddleware)
}
```

### Step 13: Create Database Migration

Create migration file:

```bash
# Create migration file (manual for now, or use migration tool)
cat > database/migrations/013_create_blog_posts_table.go << 'EOF'
package migrations

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"
)

func init() {
    Migrations = append(Migrations, &gormigrate.Migration{
        ID: "013_create_blog_posts_table",
        Migrate: func(tx *gorm.DB) error {
            type BlogPost struct {
                ID        uint   `gorm:"primaryKey"`
                Title     string `gorm:"size:255;not null;index"`
                Content   string `gorm:"type:text"`
                AuthorID  uint   `gorm:"index;not null"`
                Status    string `gorm:"size:20;default:'draft';index"`
                CreatedAt int64  `gorm:"autoCreateTime"`
                UpdatedAt int64  `gorm:"autoUpdateTime"`
                DeletedAt *int64 `gorm:"index"`
            }
            return tx.AutoMigrate(&BlogPost{})
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable("blog_posts")
        },
    })
}
EOF
```

**Run migration**:

```bash
make migrate
# Or: ./zgo migrate
```

### Step 14: Create Domain Entity

Create `internal/domain/blog_post.go`:

```go
package domain

import "time"

// BlogPost represents a blog post in the domain layer
type BlogPost struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    AuthorID  uint      `json:"author_id"`
    Status    string    `json:"status"` // draft, published
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Step 15: Validation

**Automated validation script**:

Run the validation script (see `scripts/` folder):

```bash
./.agent/skills/module-creation/scripts/validate-module.sh blog
```

**Manual checklist**:

- [ ] All 8 files created and properly structured
- [ ] Wire generation successful (`cd internal/wiring && wire`)
- [ ] Routes registered in `routes/api.go`
- [ ] Migration created and applied
- [ ] Unit tests passing (`go test ./internal/modules/blog/...`)
- [ ] Domain entity created in `internal/domain/`
- [ ] Swagger annotations added
- [ ] Handler utilities used (`pkg/handler`, `pkg/response`)

**Test the API**:

```bash
# Start server
make air

# Test endpoints (in another terminal)
# Create
curl -X POST http://localhost:8080/api/blogs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Post", "content": "Hello world"}'

# List
curl http://localhost:8080/api/blogs

# Get by ID
curl http://localhost:8080/api/blogs/1

# Update
curl -X PATCH http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "published"}'

# Delete
curl -X DELETE http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🔍 Troubleshooting

### Common Error 1: Wire Generation Fails

**Symptom**:
```
wire: blog: wire_gen.go:XX:YY: no provider found for ...
```

**Cause**: Missing or incorrect provider bindings

**Solution**:
1. Check `provider.go` includes all constructors:
   ```go
   var ProviderSet = wire.NewSet(
       NewRepository,
       wire.Bind(new(Repository), new(*repository)),
       NewService,
       wire.Bind(new(Service), new(*service)),
       NewHandler,
   )
   ```
2. Verify interface and implementation match
3. Ensure `ProviderSet` is added to `wiring/wire.go`

### Common Error 2: Routes Not Working

**Symptom**: 404 Not Found for `/api/blogs`

**Cause**: Routes not registered

**Solution**:
1. Verify `RegisterRoutes` is called in `routes/api.go`
2. Check middleware order (auth middleware might block)
3. Print registered routes:
   ```go
   for _, route := range app.Router.Routes() {
       fmt.Printf("%s %s\n", route.Method, route.Path)
   }
   ```

### Common Error 3: Database Table Not Found

**Symptom**: `Error 1146: Table 'zgo.blog_posts' doesn't exist`

**Cause**: Migration not run

**Solution**:
```bash
# Run migrations
make migrate

# Or manually
./zgo migrate
```

### Common Error 4: JSON Binding Fails

**Symptom**: 400 Bad Request, validation errors

**Cause**: Incorrect struct tags or request body

**Solution**:
1. Check binding tags in DTO:
   ```go
   type CreateBlogPostRequest struct {
       Title string `json:"title" binding:"required,max=255"`
   }
   ```
2. Verify request JSON matches field names (snake_case)
3. Test with `curl -v` to see actual error message

## 📚 Examples

See [`examples/blog-module-complete/`](./examples/blog-module-complete/) for the full implementation of a Blog module.

## 🔗 Related Skills

- [`api-development`](../api-development/): For handler patterns and pagination
- [`testing-strategy`](../testing-strategy/): For comprehensive test coverage
- [`wire-di`](../wire-di/): For advanced dependency injection scenarios
- [`database-migration`](../database-migration/): For migration best practices
- [`swagger-docs`](../swagger-docs/): For API documentation

## 📖 References

- [ZGO AGENTS.md](../../AGENTS.md) - Project development guidelines
- [DDD Layered Architecture](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Wire User Guide](https://github.com/google/wire/blob/main/docs/guide.md)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Web Framework](https://gin-gonic.com/docs/)

---

## 🎉 Success!

You've successfully created a complete DDD module following ZGO's 8-file standard!

**What's next?**
1. Add more business logic and validation
2. Write integration tests
3. Generate Swagger documentation
4. Deploy and test in staging environment
