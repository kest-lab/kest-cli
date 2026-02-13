package migration

import (
	"sort"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	return db
}

// setupTestRepository creates a repository with a fresh migrations table
func setupTestRepository(t *testing.T) (Repository, *gorm.DB) {
	db := setupTestDB(t)
	repo := NewDatabaseRepository(db, "migrations")
	err := repo.CreateRepository()
	require.NoError(t, err)
	return repo, db
}

// migrationNameGen generates valid migration names
func migrationNameGen() gopter.Gen {
	return gen.RegexMatch(`[a-z][a-z0-9_]{5,30}`)
}

// batchNumberGen generates valid batch numbers (1-100)
func batchNumberGen() gopter.Gen {
	return gen.IntRange(1, 100)
}

// migrationRecordGen generates a migration record with name and batch
func migrationRecordGen() gopter.Gen {
	return gopter.CombineGens(
		migrationNameGen(),
		batchNumberGen(),
	).Map(func(values []interface{}) struct {
		Name  string
		Batch int
	} {
		return struct {
			Name  string
			Batch int
		}{
			Name:  values[0].(string),
			Batch: values[1].(int),
		}
	})
}

// Feature: migration-system, Property 1: Repository GetRan Ordering
// For any set of migrations logged with different batch numbers and names,
// GetRan() should return them ordered by batch ascending, then by migration name ascending.
func TestProperty1_RepositoryGetRanOrdering(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GetRan returns migrations ordered by batch ASC, migration ASC", prop.ForAll(
		func(records []struct {
			Name  string
			Batch int
		}) bool {
			// Setup fresh repository for each test
			db := setupTestDB(t)
			repo := NewDatabaseRepository(db, "migrations")
			if err := repo.CreateRepository(); err != nil {
				return false
			}

			// Deduplicate migration names
			seen := make(map[string]bool)
			uniqueRecords := make([]struct {
				Name  string
				Batch int
			}, 0)
			for _, r := range records {
				if !seen[r.Name] {
					seen[r.Name] = true
					uniqueRecords = append(uniqueRecords, r)
				}
			}

			// Log all migrations
			for _, r := range uniqueRecords {
				if err := repo.Log(r.Name, r.Batch); err != nil {
					return false
				}
			}

			// Get ran migrations
			ran, err := repo.GetRan()
			if err != nil {
				return false
			}

			// Verify count matches
			if len(ran) != len(uniqueRecords) {
				return false
			}

			// Verify ordering: batch ASC, then migration ASC
			for i := 1; i < len(ran); i++ {
				// Get batch numbers for comparison
				batches, err := repo.GetMigrationBatches()
				if err != nil {
					return false
				}
				prevBatch := batches[ran[i-1]]
				currBatch := batches[ran[i]]

				if prevBatch > currBatch {
					return false // Batch should be ascending
				}
				if prevBatch == currBatch && ran[i-1] > ran[i] {
					return false // Within same batch, name should be ascending
				}
			}

			return true
		},
		gen.SliceOfN(10, migrationRecordGen()),
	))

	properties.TestingRun(t)
}

// Feature: migration-system, Property 2: Repository GetLast Batch Filtering
// For any set of migrations with multiple batches, GetLast() should return
// only migrations from the highest batch number.
func TestProperty2_RepositoryGetLastBatchFiltering(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GetLast returns only migrations from highest batch", prop.ForAll(
		func(records []struct {
			Name  string
			Batch int
		}) bool {
			// Setup fresh repository
			db := setupTestDB(t)
			repo := NewDatabaseRepository(db, "migrations")
			if err := repo.CreateRepository(); err != nil {
				return false
			}

			// Deduplicate migration names
			seen := make(map[string]bool)
			uniqueRecords := make([]struct {
				Name  string
				Batch int
			}, 0)
			for _, r := range records {
				if !seen[r.Name] {
					seen[r.Name] = true
					uniqueRecords = append(uniqueRecords, r)
				}
			}

			if len(uniqueRecords) == 0 {
				// Empty case: GetLast should return empty
				last, err := repo.GetLast()
				return err == nil && len(last) == 0
			}

			// Find max batch
			maxBatch := 0
			for _, r := range uniqueRecords {
				if r.Batch > maxBatch {
					maxBatch = r.Batch
				}
			}

			// Log all migrations
			for _, r := range uniqueRecords {
				if err := repo.Log(r.Name, r.Batch); err != nil {
					return false
				}
			}

			// Get last batch migrations
			last, err := repo.GetLast()
			if err != nil {
				return false
			}

			// Verify all returned migrations are from max batch
			for _, record := range last {
				if record.Batch != maxBatch {
					return false
				}
			}

			// Verify count matches expected
			expectedCount := 0
			for _, r := range uniqueRecords {
				if r.Batch == maxBatch {
					expectedCount++
				}
			}

			return len(last) == expectedCount
		},
		gen.SliceOfN(10, migrationRecordGen()),
	))

	properties.TestingRun(t)
}

// Feature: migration-system, Property 3: Repository Next Batch Number
// For any set of logged migrations, GetNextBatchNumber() should return
// the maximum batch number plus one. For an empty repository, it should return 1.
func TestProperty3_RepositoryNextBatchNumber(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GetNextBatchNumber returns max batch + 1", prop.ForAll(
		func(records []struct {
			Name  string
			Batch int
		}) bool {
			// Setup fresh repository
			db := setupTestDB(t)
			repo := NewDatabaseRepository(db, "migrations")
			if err := repo.CreateRepository(); err != nil {
				return false
			}

			// Deduplicate migration names
			seen := make(map[string]bool)
			uniqueRecords := make([]struct {
				Name  string
				Batch int
			}, 0)
			for _, r := range records {
				if !seen[r.Name] {
					seen[r.Name] = true
					uniqueRecords = append(uniqueRecords, r)
				}
			}

			// Empty case
			if len(uniqueRecords) == 0 {
				next, err := repo.GetNextBatchNumber()
				return err == nil && next == 1
			}

			// Find max batch
			maxBatch := 0
			for _, r := range uniqueRecords {
				if r.Batch > maxBatch {
					maxBatch = r.Batch
				}
			}

			// Log all migrations
			for _, r := range uniqueRecords {
				if err := repo.Log(r.Name, r.Batch); err != nil {
					return false
				}
			}

			// Verify next batch number
			next, err := repo.GetNextBatchNumber()
			if err != nil {
				return false
			}

			return next == maxBatch+1
		},
		gen.SliceOfN(10, migrationRecordGen()),
	))

	properties.TestingRun(t)
}

// Feature: migration-system, Property 4: Migration Log Round-Trip
// For any migration name and batch number, after calling Log(name, batch),
// the migration should appear in GetRan() results and GetMigrationBatches()
// should return the correct batch number for that migration.
func TestProperty4_MigrationLogRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Log then GetRan/GetMigrationBatches returns correct data", prop.ForAll(
		func(name string, batch int) bool {
			// Setup fresh repository
			db := setupTestDB(t)
			repo := NewDatabaseRepository(db, "migrations")
			if err := repo.CreateRepository(); err != nil {
				return false
			}

			// Log the migration
			if err := repo.Log(name, batch); err != nil {
				return false
			}

			// Verify it appears in GetRan
			ran, err := repo.GetRan()
			if err != nil {
				return false
			}
			found := false
			for _, m := range ran {
				if m == name {
					found = true
					break
				}
			}
			if !found {
				return false
			}

			// Verify batch number in GetMigrationBatches
			batches, err := repo.GetMigrationBatches()
			if err != nil {
				return false
			}
			recordedBatch, exists := batches[name]
			if !exists {
				return false
			}

			return recordedBatch == batch
		},
		migrationNameGen(),
		batchNumberGen(),
	))

	properties.TestingRun(t)
}

// Unit tests for edge cases and specific examples

func TestDatabaseRepository_EmptyRepository(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// GetRan should return empty slice
	ran, err := repo.GetRan()
	assert.NoError(t, err)
	assert.Empty(t, ran)

	// GetLast should return empty slice
	last, err := repo.GetLast()
	assert.NoError(t, err)
	assert.Empty(t, last)

	// GetNextBatchNumber should return 1
	next, err := repo.GetNextBatchNumber()
	assert.NoError(t, err)
	assert.Equal(t, 1, next)

	// GetMigrationBatches should return empty map
	batches, err := repo.GetMigrationBatches()
	assert.NoError(t, err)
	assert.Empty(t, batches)
}

func TestDatabaseRepository_DeleteMigration(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Log a migration
	err := repo.Log("test_migration", 1)
	require.NoError(t, err)

	// Verify it exists
	ran, err := repo.GetRan()
	require.NoError(t, err)
	assert.Contains(t, ran, "test_migration")

	// Delete it
	err = repo.Delete("test_migration")
	assert.NoError(t, err)

	// Verify it's gone
	ran, err = repo.GetRan()
	require.NoError(t, err)
	assert.NotContains(t, ran, "test_migration")
}

func TestDatabaseRepository_DeleteNonExistent(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Try to delete non-existent migration
	err := repo.Delete("non_existent")
	assert.ErrorIs(t, err, ErrMigrationNotFound)
}

func TestDatabaseRepository_RepositoryLifecycle(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatabaseRepository(db, "test_migrations")

	// Repository should not exist initially
	assert.False(t, repo.RepositoryExists())

	// Create repository
	err := repo.CreateRepository()
	assert.NoError(t, err)

	// Repository should exist now
	assert.True(t, repo.RepositoryExists())

	// Delete repository
	err = repo.DeleteRepository()
	assert.NoError(t, err)

	// Repository should not exist anymore
	assert.False(t, repo.RepositoryExists())
}

func TestDatabaseRepository_GetMigrations(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Log migrations in different batches
	migrations := []struct {
		name  string
		batch int
	}{
		{"migration_a", 1},
		{"migration_b", 1},
		{"migration_c", 2},
		{"migration_d", 2},
		{"migration_e", 3},
	}

	for _, m := range migrations {
		err := repo.Log(m.name, m.batch)
		require.NoError(t, err)
	}

	// Get last 3 migrations
	records, err := repo.GetMigrations(3)
	require.NoError(t, err)
	assert.Len(t, records, 3)

	// Should be ordered by batch DESC, migration DESC
	// So we expect: migration_e (batch 3), migration_d (batch 2), migration_c (batch 2)
	assert.Equal(t, "migration_e", records[0].Migration)
	assert.Equal(t, 3, records[0].Batch)
}

func TestDatabaseRepository_GetMigrationsByBatch(t *testing.T) {
	repo, _ := setupTestRepository(t)

	// Log migrations in different batches
	err := repo.Log("migration_a", 1)
	require.NoError(t, err)
	err = repo.Log("migration_b", 1)
	require.NoError(t, err)
	err = repo.Log("migration_c", 2)
	require.NoError(t, err)

	// Get batch 1 migrations
	records, err := repo.GetMigrationsByBatch(1)
	require.NoError(t, err)
	assert.Len(t, records, 2)

	names := []string{records[0].Migration, records[1].Migration}
	sort.Strings(names)
	assert.Equal(t, []string{"migration_a", "migration_b"}, names)
}

func TestDatabaseRepository_ConfigurableTableName(t *testing.T) {
	db := setupTestDB(t)

	// Create repository with custom table name
	repo := NewDatabaseRepository(db, "custom_migrations")
	assert.Equal(t, "custom_migrations", repo.GetTable())

	// Create the table
	err := repo.CreateRepository()
	require.NoError(t, err)

	// Verify table exists with custom name
	assert.True(t, db.Migrator().HasTable("custom_migrations"))
	assert.False(t, db.Migrator().HasTable("migrations"))
}

func TestDatabaseRepository_DefaultTableName(t *testing.T) {
	db := setupTestDB(t)

	// Create repository with empty table name (should default to "migrations")
	repo := NewDatabaseRepository(db, "")
	assert.Equal(t, "migrations", repo.GetTable())
}
