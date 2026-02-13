// Package migration provides a Laravel-inspired database migration system
// with repository pattern, event-driven architecture, and fluent schema builder.
package migration

import "gorm.io/gorm"

// Migration defines the contract for individual migration files.
// Each migration must implement Up() for applying changes and Down() for reverting.
type Migration interface {
	// Up applies the migration
	Up(db *gorm.DB) error

	// Down reverts the migration
	Down(db *gorm.DB) error

	// GetConnection returns the database connection name (empty for default)
	GetConnection() string

	// WithinTransaction indicates if the migration should run in a transaction
	WithinTransaction() bool
}

// BaseMigration provides default implementations for the Migration interface.
// Embed this struct in your migrations to get sensible defaults.
type BaseMigration struct {
	// Connection specifies which database connection to use (empty for default)
	Connection string

	// UseTransaction indicates whether to wrap the migration in a transaction
	UseTransaction bool
}

// GetConnection returns the database connection name
func (m *BaseMigration) GetConnection() string {
	return m.Connection
}

// WithinTransaction returns whether the migration should run in a transaction
func (m *BaseMigration) WithinTransaction() bool {
	return m.UseTransaction
}
