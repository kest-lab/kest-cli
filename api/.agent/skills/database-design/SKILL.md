---
name: database-design
description: Database design best practices, indexing strategies, and migration patterns for ZGO
version: 1.0.0
category: architecture
tags: [database, sql, gorm, migration, optimization]
author: ZGO Team
updated: 2026-01-24
---

# Database Design Standards

## 📋 Purpose

This skill provides the definitive standards for database design in the ZGO project. It ensures data integrity, query performance, and consistent migration workflows across all modules.

## 🎯 When to Use

- Designing tables for a new module
- Adding or modifying columns in existing tables
- Optimizing slow queries using indexes
- Creating migration scripts
- Reviewing database-related code in PRs

---

## 🏗️ Table Design Best Practices

### 1. Naming Conventions
- **Tables**: `snake_case`, plural (e.g., `users`, `user_orders`).
- **Columns**: `snake_case` (e.g., `created_at`, `is_active`).
- **Persistence Objects (PO)**: Must have a `PO` suffix and a `TableName()` method.

### 2. Mandatory Columns
Every table must include these core audit and identification fields:
- `id`: BigInt / Serial (Primary Key)
- `created_at`: Time with zone
- `updated_at`: Time with zone
- `deleted_at`: Indexable time (for GORM soft deletes)

```go
// Example PO Structure
type UserPO struct {
    ID        uint           `gorm:"primaryKey"`
    Username  string         `gorm:"size:255;not null;uniqueIndex"`
    Email     string         `gorm:"size:255;not null;uniqueIndex"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserPO) TableName() string {
    return "users"
}
```

### 3. Data Types
- **Strings**: Use `string` in Go. Specify `size:255` or `text` in GORM tags.
- **Booleans**: Use `bool`. Defaults to `false`.
- **Enums**: Use `int` or `string` with validation in the Service layer. Do not rely solely on DB-level enums.
- **JSON**: Use `datatypes.JSON` or `string` for flexible data blobs.

---

## 🚀 Indexing Strategies

### 1. Primary Rules
- **Unique Indexes**: Use for fields like `email`, `username`, or `slug` to prevent logical duplicates.
- **Query Alignment**: Index columns used frequently in `WHERE` clauses (e.g., `status`, `user_id`).
- **Composite Indexes**: Use for queries that filter by multiple columns together. Order matters: high cardinality fields go first.

### 2. GORM Index Tags
```go
Username  string `gorm:"uniqueIndex:idx_users_username"`
Status    string `gorm:"index:idx_users_status"`
// Composite Index
TenantID  uint   `gorm:"index:idx_tenant_created;priority:1"`
CreatedAt time.Time `gorm:"index:idx_tenant_created;priority:2"`
```

---

## 🛠️ Migration Patterns

### 1. Principle: Forward Only
We use GORM's `AutoMigrate` for simple development, but **production changes must use explicit migration files** or controlled scripts to prevent data loss.

### 2. Safe Conversions
- **Adding columns**: Safe. Always provide a default or allow NULL.
- **Renaming columns**: Unsafe. Use a new column, copy data, then drop the old one.
- **Adding indexes**: Safe, but can be slow on large tables (Concurrently in Postgres).

---

## ⚡ SQL Optimization

### 1. Avoid `SELECT *`
In repositories, only fetch what you need if performance is critical. However, for standard CRUD, GORM's default behavior is acceptable.

### 2. N+1 Query Problem
Never perform a database query inside a loop. Use Eager Loading.

```go
// ❌ WRONG
for _, user := range users {
    var profile Profile
    db.Where("user_id = ?", user.ID).First(&profile) // N queries
}

// ✅ CORRECT (GORM Preload)
db.Preload("Profile").Find(&users) // 2 queries total
```

### 3. Explain Analyze
Whenever a query feels slow, use `EXPLAIN ANALYZE` in your DB console to check for sequential scans vs. index hits.

---

## ✅ Verification Checklist

- [ ] All tables have `id`, `created_at`, `updated_at`, `deleted_at`.
- [ ] Table names are plural `snake_case`.
- [ ] Index tags are added for all search/filter criteria.
- [ ] Unique constraints are placed on unique identifiers.
- [ ] Pointer types are used for optional/nullable fields in POs.
- [ ] `TableName()` is explicitly defined in `model.go`.
- [ ] No DB-level logic (triggers/stored procs) - keep logic in Go.

---

## 📚 Complete Examples

- [**Module model.go Example**](./examples/module_model_example.go)
- [**Complex Indexing Pattern**](./examples/indexing_patterns.go)
- [**Safe Migration Script**](./examples/migration_script.go)

---

## 🔗 Related Skills

- [`module-creation`](../module-creation/): Where `model.go` lives.
- [`coding-standards`](../coding-standards/): General naming and error handling.

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team
