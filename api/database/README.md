# Database Migrations

## Overview

ZGO uses a Laravel-style migration system with:
- **Timestamped filenames** for automatic ordering
- **Auto-registration** using `init()` functions
- **CLI generator** for creating new migrations
- **Rollback support** for database changes

## Directory Structure

```
database/
└── migrations/
    ├── migrations.go                              # Registry (auto-populated)
    ├── 2025_06_18_000000_create_users_table.go
    ├── 2025_06_18_000001_seed_default_users.go
    ├── 2025_12_26_000000_create_roles_table.go
    ├── 2025_12_26_000001_create_permissions_table.go
    ├── 2025_12_26_000002_create_role_permissions_table.go
    ├── 2025_12_26_000003_create_user_roles_table.go
    └── 2025_12_26_000004_seed_default_roles.go
```

## Creating Migrations

### Generate Migration File

```bash
./zgo make:migration create_posts_table
```

**Output:**
```
✓ Migration created: database/migrations/2025_12_26_012920_create_posts_table.go
ℹ Migration ID: 2025_12_26_012920_create_posts_table
```

### Generated File Structure

```go
package migrations

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"
)

func init() {
    register(&gormigrate.Migration{
        ID: "2025_12_26_012920_create_posts_table",
        Migrate: func(db *gorm.DB) error {
            // TODO: Implement migration logic
            // Example: return db.AutoMigrate(&YourModel{})
            return nil
        },
        Rollback: func(db *gorm.DB) error {
            // TODO: Implement rollback logic
            // Example: return db.Migrator().DropTable("your_table")
            return nil
        },
    })
}
```

### Implement Migration Logic

```go
func init() {
    register(&gormigrate.Migration{
        ID: "2025_12_26_012920_create_posts_table",
        Migrate: func(db *gorm.DB) error {
            return db.AutoMigrate(&blog.Post{})
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropTable("posts")
        },
    })
}
```

## Running Migrations

```bash
./zgo migrate
```

**Output:**
```
ℹ Running migrations...
✓ Migrations completed
✓ Done in 0s
```

## Migration Naming Convention

**Format**: `YYYY_MM_DD_HHMMSS_description.go`

**Examples:**
- `2025_06_18_000000_create_users_table.go` - Create users table
- `2025_06_18_000001_seed_default_users.go` - Seed default users
- `2025_12_26_012920_add_status_to_posts.go` - Add column to existing table

**Benefits:**
- ✅ Automatic chronological ordering
- ✅ No naming conflicts (timestamp ensures uniqueness)
- ✅ Clear migration history

## Auto-Registration

Migrations are automatically registered using `init()` functions:

```go
func init() {
    register(&gormigrate.Migration{
        ID: "2025_12_26_012920_create_posts_table",
        // ...
    })
}
```

**No manual registration needed!** The `register()` function adds migrations to the global registry automatically when the package is imported.

## Migration Tracking

All executed migrations are tracked in the `migrations` table:

```sql
SELECT * FROM migrations;
```

**Output:**
```
| id                                      |
|-----------------------------------------|
| 2025_06_18_000000_create_users_table    |
| 2025_06_18_000001_seed_default_users    |
| 2025_12_26_000000_create_roles_table    |
| ...                                     |
```

## Best Practices

1. **One Migration, One Purpose**
   - Each migration should do one thing (create table, add column, etc.)

2. **Always Provide Rollback**
   - Every migration should have a working `Rollback` function

3. **Use Idempotent Operations**
   - Use `FirstOrCreate` for seeds to allow re-running

4. **Test Migrations**
   - Test both `Migrate` and `Rollback` functions

5. **Never Edit Existing Migrations**
   - Once deployed, create a new migration instead

## Common Patterns

### Create Table
```go
Migrate: func(db *gorm.DB) error {
    return db.AutoMigrate(&YourModel{})
}
```

### Add Column
```go
Migrate: func(db *gorm.DB) error {
    return db.Migrator().AddColumn(&User{}, "Status")
}
```

### Seed Data
```go
Migrate: func(db *gorm.DB) error {
    users := []User{
        {Username: "admin", Email: "admin@example.com"},
    }
    for _, user := range users {
        db.FirstOrCreate(&user, User{Username: user.Username})
    }
    return nil
}
```

### Drop Table
```go
Rollback: func(db *gorm.DB) error {
    return db.Migrator().DropTable("your_table")
}
```

## Troubleshooting

### Migration Already Exists
If you see "migration file already exists", the timestamp collision is extremely rare. Wait a second and try again.

### Migration Not Running
Ensure the migration file is in `database/migrations/` and has an `init()` function that calls `register()`.

### Rollback Not Working
Check that your `Rollback` function properly reverses the `Migrate` operation.
