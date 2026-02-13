package migration

// MigrationRecord represents a single migration record in the repository.
// It tracks which migrations have been run and in which batch.
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Migration string `gorm:"size:255;not null"`
	Batch     int    `gorm:"not null"`
}

// Repository defines the contract for migration record storage.
// This interface allows swapping implementations without changing the Migrator.
type Repository interface {
	// GetRan returns all completed migration names ordered by batch ascending,
	// then by migration name ascending.
	GetRan() ([]string, error)

	// GetMigrations returns the last N migrations ordered by batch descending,
	// then by migration name descending.
	GetMigrations(steps int) ([]MigrationRecord, error)

	// GetMigrationsByBatch returns all migrations for a specific batch number.
	GetMigrationsByBatch(batch int) ([]MigrationRecord, error)

	// GetLast returns migrations from the last (highest) batch number.
	GetLast() ([]MigrationRecord, error)

	// GetMigrationBatches returns all migrations with their batch numbers
	// as a map where key is migration name and value is batch number.
	GetMigrationBatches() (map[string]int, error)

	// Log records that a migration was run with the given batch number.
	Log(migration string, batch int) error

	// Delete removes a migration record by name.
	Delete(migration string) error

	// GetNextBatchNumber returns the next batch number (max batch + 1).
	// Returns 1 if no migrations have been run.
	GetNextBatchNumber() (int, error)

	// CreateRepository creates the migrations table if it doesn't exist.
	CreateRepository() error

	// RepositoryExists checks if the migrations table exists.
	RepositoryExists() bool

	// DeleteRepository drops the migrations table.
	DeleteRepository() error

	// SetTable sets the migrations table name.
	SetTable(table string)

	// GetTable returns the current migrations table name.
	GetTable() string
}
