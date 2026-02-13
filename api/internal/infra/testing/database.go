package testing

import (
	"testing"

	"gorm.io/gorm"
)

// DatabaseTestCase provides database testing utilities
type DatabaseTestCase struct {
	t  *testing.T
	db *gorm.DB
}

// NewDatabaseTestCase creates a new database test case
func NewDatabaseTestCase(t *testing.T, db *gorm.DB) *DatabaseTestCase {
	return &DatabaseTestCase{t: t, db: db}
}

// AssertDatabaseHas asserts a record exists in the database
func (tc *DatabaseTestCase) AssertDatabaseHas(table string, conditions map[string]interface{}) *DatabaseTestCase {
	var count int64
	query := tc.db.Table(table)
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	query.Count(&count)

	if count == 0 {
		tc.t.Errorf("Expected record in table '%s' with conditions %v, but none found", table, conditions)
	}
	return tc
}

// AssertDatabaseMissing asserts a record does not exist
func (tc *DatabaseTestCase) AssertDatabaseMissing(table string, conditions map[string]interface{}) *DatabaseTestCase {
	var count int64
	query := tc.db.Table(table)
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	query.Count(&count)

	if count > 0 {
		tc.t.Errorf("Expected no record in table '%s' with conditions %v, but found %d", table, conditions, count)
	}
	return tc
}

// AssertDatabaseCount asserts the number of records in a table
func (tc *DatabaseTestCase) AssertDatabaseCount(table string, expected int64, conditions ...map[string]interface{}) *DatabaseTestCase {
	var count int64
	query := tc.db.Table(table)
	if len(conditions) > 0 {
		for k, v := range conditions[0] {
			query = query.Where(k+" = ?", v)
		}
	}
	query.Count(&count)

	if count != expected {
		tc.t.Errorf("Expected %d records in table '%s', got %d", expected, table, count)
	}
	return tc
}

// AssertSoftDeleted asserts a record is soft deleted
func (tc *DatabaseTestCase) AssertSoftDeleted(table string, conditions map[string]interface{}) *DatabaseTestCase {
	var count int64
	query := tc.db.Table(table).Unscoped()
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	query.Where("deleted_at IS NOT NULL").Count(&count)

	if count == 0 {
		tc.t.Errorf("Expected soft deleted record in table '%s' with conditions %v", table, conditions)
	}
	return tc
}

// AssertNotSoftDeleted asserts a record is not soft deleted
func (tc *DatabaseTestCase) AssertNotSoftDeleted(table string, conditions map[string]interface{}) *DatabaseTestCase {
	var count int64
	query := tc.db.Table(table)
	for k, v := range conditions {
		query = query.Where(k+" = ?", v)
	}
	query.Where("deleted_at IS NULL").Count(&count)

	if count == 0 {
		tc.t.Errorf("Expected non-soft-deleted record in table '%s' with conditions %v", table, conditions)
	}
	return tc
}

// Seed runs a seeder function
func (tc *DatabaseTestCase) Seed(seeder func(db *gorm.DB) error) *DatabaseTestCase {
	if err := seeder(tc.db); err != nil {
		tc.t.Fatalf("Failed to seed database: %v", err)
	}
	return tc
}

// Truncate truncates the specified tables
func (tc *DatabaseTestCase) Truncate(tables ...string) *DatabaseTestCase {
	for _, table := range tables {
		if err := tc.db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			tc.t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
	return tc
}

// Transaction runs the test in a transaction that's rolled back afterward
func (tc *DatabaseTestCase) Transaction(fn func(tx *gorm.DB)) {
	tx := tc.db.Begin()
	defer tx.Rollback()

	fn(tx)
}

// DB returns the database connection
func (tc *DatabaseTestCase) DB() *gorm.DB {
	return tc.db
}
