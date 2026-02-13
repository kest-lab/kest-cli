package migration

import (
	"gorm.io/gorm"
)

// databaseRepository implements Repository using a database table.
// It stores migration records in a configurable table (default: "migrations").
type databaseRepository struct {
	db        *gorm.DB
	tableName string
}

// NewDatabaseRepository creates a new database-backed migration repository.
// If tableName is empty, it defaults to "migrations".
func NewDatabaseRepository(db *gorm.DB, tableName string) Repository {
	if tableName == "" {
		tableName = "migrations"
	}
	return &databaseRepository{
		db:        db,
		tableName: tableName,
	}
}

// SetTable sets the migrations table name.
func (r *databaseRepository) SetTable(table string) {
	r.tableName = table
}

// GetTable returns the current migrations table name.
func (r *databaseRepository) GetTable() string {
	return r.tableName
}

// GetRan returns all completed migration names ordered by batch ascending,
// then by migration name ascending.
func (r *databaseRepository) GetRan() ([]string, error) {
	var migrations []string
	err := r.db.Table(r.tableName).
		Order("batch ASC, migration ASC").
		Pluck("migration", &migrations).Error
	return migrations, err
}

// GetMigrations returns the last N migrations ordered by batch descending,
// then by migration name descending.
func (r *databaseRepository) GetMigrations(steps int) ([]MigrationRecord, error) {
	var records []MigrationRecord
	err := r.db.Table(r.tableName).
		Order("batch DESC, migration DESC").
		Limit(steps).
		Find(&records).Error
	return records, err
}

// GetMigrationsByBatch returns all migrations for a specific batch number.
func (r *databaseRepository) GetMigrationsByBatch(batch int) ([]MigrationRecord, error) {
	var records []MigrationRecord
	err := r.db.Table(r.tableName).
		Where("batch = ?", batch).
		Order("migration DESC").
		Find(&records).Error
	return records, err
}

// GetLast returns migrations from the last (highest) batch number.
func (r *databaseRepository) GetLast() ([]MigrationRecord, error) {
	lastBatch, err := r.getLastBatchNumber()
	if err != nil {
		return nil, err
	}
	if lastBatch == 0 {
		return []MigrationRecord{}, nil
	}
	return r.GetMigrationsByBatch(lastBatch)
}

// getLastBatchNumber returns the highest batch number, or 0 if no migrations exist.
func (r *databaseRepository) getLastBatchNumber() (int, error) {
	var max *int
	err := r.db.Table(r.tableName).
		Select("MAX(batch)").
		Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if max == nil {
		return 0, nil
	}
	return *max, nil
}

// Log records that a migration was run with the given batch number.
func (r *databaseRepository) Log(migration string, batch int) error {
	record := MigrationRecord{
		Migration: migration,
		Batch:     batch,
	}
	return r.db.Table(r.tableName).Create(&record).Error
}

// Delete removes a migration record by name.
func (r *databaseRepository) Delete(migration string) error {
	result := r.db.Table(r.tableName).Where("migration = ?", migration).Delete(&MigrationRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMigrationNotFound
	}
	return nil
}

// GetNextBatchNumber returns the next batch number (max batch + 1).
// Returns 1 if no migrations have been run.
func (r *databaseRepository) GetNextBatchNumber() (int, error) {
	last, err := r.getLastBatchNumber()
	if err != nil {
		return 1, err
	}
	return last + 1, nil
}

// GetMigrationBatches returns all migrations with their batch numbers
// as a map where key is migration name and value is batch number.
func (r *databaseRepository) GetMigrationBatches() (map[string]int, error) {
	var records []MigrationRecord
	err := r.db.Table(r.tableName).
		Select("migration, batch").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	batches := make(map[string]int, len(records))
	for _, record := range records {
		batches[record.Migration] = record.Batch
	}
	return batches, nil
}

// CreateRepository creates the migrations table if it doesn't exist.
// Supports MySQL, PostgreSQL, and SQLite dialects.
func (r *databaseRepository) CreateRepository() error {
	dialect := r.db.Dialector.Name()

	var sql string
	switch dialect {
	case "mysql":
		sql = `CREATE TABLE IF NOT EXISTS ` + r.tableName + ` (
			id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			migration VARCHAR(255) NOT NULL,
			batch INT NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
	case "postgres":
		sql = `CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)`
	case "sqlite":
		sql = `CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)`
	default:
		// Fallback to SQLite-compatible syntax
		sql = `CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)`
	}

	return r.db.Exec(sql).Error
}

// RepositoryExists checks if the migrations table exists.
func (r *databaseRepository) RepositoryExists() bool {
	return r.db.Migrator().HasTable(r.tableName)
}

// DeleteRepository drops the migrations table.
func (r *databaseRepository) DeleteRepository() error {
	dialect := r.db.Dialector.Name()

	var sql string
	switch dialect {
	case "postgres":
		sql = `DROP TABLE IF EXISTS "` + r.tableName + `" CASCADE`
	default:
		sql = `DROP TABLE IF EXISTS ` + r.tableName
	}

	return r.db.Exec(sql).Error
}
