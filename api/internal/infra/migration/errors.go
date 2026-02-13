package migration

import "errors"

// Common errors returned by the migration system.
var (
	// ErrRepositoryNotExists is returned when operations are attempted
	// before the migrations table has been created.
	ErrRepositoryNotExists = errors.New("migration repository does not exist")

	// ErrMigrationNotFound is returned when trying to delete or access
	// a migration record that doesn't exist.
	ErrMigrationNotFound = errors.New("migration not found")

	// ErrMigrationNotRegistered is returned when a migration file
	// is not found in the registry.
	ErrMigrationNotRegistered = errors.New("migration not registered")

	// ErrNoMigrationsToRun is returned when there are no pending migrations.
	ErrNoMigrationsToRun = errors.New("no migrations to run")

	// ErrNoMigrationsToRollback is returned when there are no migrations
	// to rollback.
	ErrNoMigrationsToRollback = errors.New("no migrations to rollback")
)

// MigrationError wraps an error with migration context.
type MigrationError struct {
	Migration string
	Method    string // "up" or "down"
	Err       error
}

// Error returns the error message with migration context.
func (e *MigrationError) Error() string {
	return "migration " + e.Migration + " " + e.Method + " failed: " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *MigrationError) Unwrap() error {
	return e.Err
}

// NewMigrationError creates a new MigrationError.
func NewMigrationError(migration, method string, err error) *MigrationError {
	return &MigrationError{
		Migration: migration,
		Method:    method,
		Err:       err,
	}
}
