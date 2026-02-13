# Database Seeders

## Overview

ZGO provides a Laravel-style seeder system for populating your database with test or default data.

## Directory Structure

```
database/
└── seeders/
    ├── seeders.go        # Registry (auto-populated)
    ├── user_seeder.go    # User data seeder
    └── role_seeder.go    # Role data seeder
```

## Creating Seeders

### Generate Seeder File

```bash
./zgo make:seeder Product
```

**Output:**
```
✓ Seeder created: database/seeders/product_seeder.go
ℹ Run with: ./zgo db:seed
```

### Generated File Structure

```go
package seeders

import (
    "github.com/zgiai/zgo/internal/infra/database"
)

type ProductSeeder struct{}

func (s *ProductSeeder) Run() error {
    db := database.GetDB()

    // TODO: Implement seeder logic
    // Example:
    // items := []YourModel{
    //     {Field: "value"},
    // }
    // for _, item := range items {
    //     db.FirstOrCreate(&item, YourModel{Field: item.Field})
    // }

    return nil
}

func init() {
    register(&ProductSeeder{})
}
```

### Implement Seeder Logic

```go
package seeders

import (
    "github.com/zgiai/zgo/internal/modules/blog"
    "github.com/zgiai/zgo/internal/infra/database"
)

type PostSeeder struct{}

func (s *PostSeeder) Run() error {
    db := database.GetDB()

    posts := []blog.Post{
        {Title: "First Post", Content: "Hello World", Status: "published"},
        {Title: "Second Post", Content: "Another post", Status: "draft"},
    }

    for _, post := range posts {
        if err := db.FirstOrCreate(&post, blog.Post{Title: post.Title}).Error; err != nil {
            return err
        }
    }

    return nil
}

func init() {
    register(&PostSeeder{})
}
```

## Running Seeders

### Run All Seeders

```bash
./zgo db:seed
```

Or use the alias:

```bash
./zgo seed
```

**Output:**
```
ℹ Running database seeders...
✓ Seeders completed
```

### Run with Migrations

```bash
./zgo migrate --seed
```

## Auto-Registration

Seeders are automatically registered using `init()` functions:

```go
func init() {
    register(&YourSeeder{})
}
```

**No manual registration needed!** The `register()` function adds seeders to the global registry automatically.

## Seeder Interface

All seeders must implement the `Seeder` interface:

```go
type Seeder interface {
    Run() error
}
```

## Best Practices

1. **Use FirstOrCreate for Idempotency**
   - Allows seeders to be run multiple times safely
   ```go
   db.FirstOrCreate(&item, Model{UniqueField: item.UniqueField})
   ```

2. **One Seeder Per Model**
   - Keep seeders focused on a single model or related data

3. **Order Matters**
   - Seeders run in the order they're registered (file alphabetical order)
   - Name files to control execution order if needed

4. **Handle Errors**
   - Always return errors for proper error handling
   ```go
   if err := db.FirstOrCreate(&item).Error; err != nil {
       return err
   }
   ```

5. **Use Realistic Data**
   - Seed data should be realistic for testing
   - Consider using faker libraries for large datasets

## Example Seeders

### User Seeder
```go
type UserSeeder struct{}

func (s *UserSeeder) Run() error {
    db := database.GetDB()

    users := []user.User{
        {
            Username: "admin",
            Email:    "admin@example.com",
            Password: "$2a$10$...", // bcrypt hash
            Status:   1,
        },
    }

    for _, u := range users {
        db.FirstOrCreate(&u, user.User{Email: u.Email})
    }

    return nil
}
```

### Role Seeder
```go
type RoleSeeder struct{}

func (s *RoleSeeder) Run() error {
    db := database.GetDB()

    roles := []permission.Role{
        {Name: "admin", DisplayName: "Administrator"},
        {Name: "user", DisplayName: "User"},
    }

    for _, role := range roles {
        db.FirstOrCreate(&role, permission.Role{Name: role.Name})
    }

    return nil
}
```

## Seeder vs Migration

| Feature | Migration | Seeder |
|---------|-----------|--------|
| Purpose | Schema changes | Data population |
| When | Always run | Optional |
| Idempotent | Should be | Must be |
| Production | Yes | Usually dev/test only |
| Tracked | Yes (migrations table) | No |

## Common Patterns

### Bulk Insert
```go
func (s *ProductSeeder) Run() error {
    db := database.GetDB()
    
    products := []Product{
        {Name: "Product 1", Price: 100},
        {Name: "Product 2", Price: 200},
        // ... many more
    }
    
    return db.Create(&products).Error
}
```

### Relationships
```go
func (s *PostSeeder) Run() error {
    db := database.GetDB()
    
    // Find or create user first
    var user User
    db.FirstOrCreate(&user, User{Email: "author@example.com"})
    
    // Create posts with user relationship
    posts := []Post{
        {Title: "Post 1", UserID: user.ID},
    }
    
    return db.Create(&posts).Error
}
```

### Conditional Seeding
```go
func (s *DemoSeeder) Run() error {
    db := database.GetDB()
    
    // Only seed in development
    if os.Getenv("APP_ENV") != "production" {
        // Seed demo data
    }
    
    return nil
}
```

## Troubleshooting

### Seeder Not Running
- Ensure the seeder file is in `database/seeders/`
- Check that `init()` function calls `register()`
- Verify the seeder implements the `Seeder` interface

### Duplicate Data
- Use `FirstOrCreate` instead of `Create`
- Specify unique fields in the condition

### Import Cycle
- Seeders should import models, not vice versa
- Keep seeder logic simple and focused
