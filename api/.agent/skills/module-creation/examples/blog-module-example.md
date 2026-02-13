# Blog Module - Complete Example

This is a complete implementation of a Blog module following ZGO's 8-file DDD standard.

## Module Specification

- **Module Name**: Blog
- **Domain Entity**: BlogPost
- **Table**: `blog_posts`
- **Features**: Full CRUD operations + pagination

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/blogs` | List all posts (paginated) | No |
| POST | `/api/blogs` | Create new post | Yes |
| GET | `/api/blogs/:id` | Get post by ID | No |
| PATCH | `/api/blogs/:id` | Update post | Yes |
| DELETE | `/api/blogs/:id` | Delete post | Yes |

## Database Schema

```sql
CREATE TABLE blog_posts (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    author_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_title (title),
    INDEX idx_author_id (author_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);
```

## File Structure

```
internal/modules/blog/
├── model.go              # BlogPostPO - database entity
├── dto.go                # Request/Response DTOs + mappers
├── repository.go         # Repository interface + implementation
├── service.go            # Service interface + implementation
├── handler.go            # HTTP handlers
├── routes.go             # Route registration
├── provider.go           # Wire DI configuration
└── service_test.go       # Unit tests
```

## Implementation Highlights

### 1. DTOs and Domain Separation

```
Request DTO → Service (Domain) → Repository (PO)
                                 ← Repository (Domain)
            ← Service (Domain)
Response DTO ←
```

### 2. Mapper Functions

Three key mappers in `dto.go`:
- `ToBlogPostPO()` - Domain → PO (for database write)
- `FromBlogPostPO()` - PO → Domain (for database read)
- `ToResponse()` - Domain → Response DTO (for API response)

### 3. Repository Pattern

Repository layer:
- **Only** handles database operations
- **Always** returns domain entities (not PO)
- Uses GORM for ORM operations
- Implements soft delete

### 4. Service Layer

Service layer:
- Contains business logic
- Input validation
- Default values
- Business rule enforcement
- Works entirely with domain entities

### 5. Handler Best Practices

Handlers use ZGO utilities:
- `handler.ParseID()` - Parse URL params
- `handler.GetUserID()` - Get authenticated user
- `handler.BindJSON()` - Bind request body
- `handler.QueryInt()` - Parse query params
- `response.Success()` - Standard success response
- `response.HandleError()` - Error handling

### 6. Wire Dependency Injection

Provider set binds interfaces to implementations:
```go
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*repository)),
    NewService,
    wire.Bind(new(Service), new(*service)),
    NewHandler,
)
```

## Testing

Unit tests use mocks:
- Mock repository interface
- Test service business logic
- Use `testify/mock` and `testify/assert`

```bash
# Run tests
go test ./internal/modules/blog/... -v
```

## Usage Example

### 1. Create a Blog Post

```bash
curl -X POST http://localhost:8080/api/blogs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "status": "draft"
  }'
```

Response:
```json
{
  "code": 201,
  "message": "Created",
  "data": {
    "id": 1,
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "author_id": 123,
    "status": "draft",
    "created_at": "2026-01-24T10:30:00Z",
    "updated_at": "2026-01-24T10:30:00Z"
  }
}
```

### 2. List Blog Posts

```bash
curl http://localhost:8080/api/blogs?page=1&per_page=20
```

Response:
```json
{
  "code": 200,
  "message": "OK",
  "data": [
    {
      "id": 1,
      "title": "My First Post",
      "content": "...",
      "author_id": 123,
      "status": "draft",
      "created_at": "2026-01-24T10:30:00Z",
      "updated_at": "2026-01-24T10:30:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "per_page": 20,
    "total_pages": 1
  }
}
```

### 3. Update a Blog Post

```bash
curl -X PATCH http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "published"
  }'
```

### 4. Delete a Blog Post

```bash
curl -X DELETE http://localhost:8080/api/blogs/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "code": 204,
  "message": "No Content"
}
```

## Key Takeaways

1. **8-file structure** ensures consistency across all modules
2. **Domain layer** is the source of truth for business entities
3. **Repository** abstracts database details
4. **Service** contains all business logic
5. **Handler** is thin, delegates to service
6. **Wire** automates dependency injection
7. **Testing** focuses on business logic with mocks

## Common Pitfalls to Avoid

❌ **Don't**: Put business logic in handlers
✅ **Do**: Keep handlers thin, delegate to service

❌ **Don't**: Return PO from repository
✅ **Do**: Always convert PO → Domain in repository

❌ **Don't**: Skip input validation in service
✅ **Do**: Validate all inputs, enforce business rules

❌ **Don't**: Use global state or singleton
✅ **Do**: Use dependency injection

❌ **Don't**: Forget to add soft delete (DeletedAt)
✅ **Do**: Always include GORM soft delete

## Further Reading

- [ZGO AGENTS.md](../../../AGENTS.md)
- [DDD Layered Architecture](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Wire Dependency Injection](https://github.com/google/wire)
