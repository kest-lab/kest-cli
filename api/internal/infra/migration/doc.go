// Package migration provides a Laravel-inspired database migration system
// for the ZGO framework.
//
// The migration system consists of several interconnected components:
//
//   - Repository: Abstraction for migration record storage
//   - Migration: Interface for individual migration files
//   - MigratorOptions: Configuration for migration operations
//   - RollbackOptions: Configuration for rollback operations
//
// # Basic Usage
//
// Create a migration by implementing the Migration interface:
//
//	type CreateUsersTable struct {
//	    migration.BaseMigration
//	}
//
//	func (m *CreateUsersTable) Up(db *gorm.DB) error {
//	    return db.AutoMigrate(&User{})
//	}
//
//	func (m *CreateUsersTable) Down(db *gorm.DB) error {
//	    return db.Migrator().DropTable("users")
//	}
//
// # Repository Pattern
//
// The Repository interface allows swapping storage implementations:
//
//	repo := migration.NewDatabaseRepository(db, "migrations")
//	if !repo.RepositoryExists() {
//	    repo.CreateRepository()
//	}
//
// # Migration Options
//
// Configure migration behavior with options:
//
//	opts := migration.NewMigratorOptions().
//	    WithPretend().  // Show SQL without executing
//	    WithStep()      // Increment batch per migration
//
// # Rollback Options
//
// Configure rollback behavior:
//
//	opts := migration.NewRollbackOptions().
//	    WithSteps(3).   // Rollback last 3 migrations
//	    WithPretend()   // Show SQL without executing
package migration
