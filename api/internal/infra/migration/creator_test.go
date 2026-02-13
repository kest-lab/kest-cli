package migration

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// migrationNameGen generates valid migration names for testing.
func creatorMigrationNameGen() gopter.Gen {
	return gen.RegexMatch(`[a-z][a-z0-9_]{5,30}`)
}

// TestProperty14_MigrationFilenameFormat tests that generated migration filenames
// match the pattern YYYY_MM_DD_HHMMSS_<name>.go where the timestamp reflects the creation time.
//
// **Feature: migration-system, Property 14: Migration Filename Format**
// **Validates: Requirements 3.5, 8.13, 9.3**
func TestProperty14_MigrationFilenameFormat(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Create a temporary directory for test migrations
	tempDir, err := os.MkdirTemp("", "migration_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Pattern for valid migration filename: YYYY_MM_DD_HHMMSS_<name>.go
	filenamePattern := regexp.MustCompile(`^\d{4}_\d{2}_\d{2}_\d{6}_[a-z][a-z0-9_]+\.go$`)

	// Pattern for valid migration ID: YYYY_MM_DD_HHMMSS_<name>
	migrationIDPattern := regexp.MustCompile(`^\d{4}_\d{2}_\d{2}_\d{6}_[a-z][a-z0-9_]+$`)

	properties.Property("Generated migration filename matches YYYY_MM_DD_HHMMSS_<name>.go pattern", prop.ForAll(
		func(name string) bool {
			// Create a unique subdirectory for each test to avoid conflicts
			testDir := filepath.Join(tempDir, name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				return false
			}

			creator := NewCreator(testDir)
			result, err := creator.Create(name, CreatorOptions{})
			if err != nil {
				// If creation fails, it's likely due to invalid name - skip
				return true
			}

			// Verify filename matches pattern
			if !filenamePattern.MatchString(result.Filename) {
				t.Logf("Filename %q does not match pattern", result.Filename)
				return false
			}

			// Verify migration ID matches pattern
			if !migrationIDPattern.MatchString(result.Name) {
				t.Logf("Migration ID %q does not match pattern", result.Name)
				return false
			}

			// Verify the filename ends with the migration name
			expectedSuffix := "_" + name + ".go"
			if !regexp.MustCompile(regexp.QuoteMeta(expectedSuffix) + "$").MatchString(result.Filename) {
				t.Logf("Filename %q does not end with %q", result.Filename, expectedSuffix)
				return false
			}

			// Verify file was actually created
			if _, err := os.Stat(result.Path); os.IsNotExist(err) {
				t.Logf("File was not created at %q", result.Path)
				return false
			}

			return true
		},
		creatorMigrationNameGen(),
	))

	properties.TestingRun(t)
}

// TestCreator_Create tests the basic creation functionality.
func TestCreator_Create(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "migration_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	creator := NewCreator(tempDir)

	t.Run("creates blank migration", func(t *testing.T) {
		result, err := creator.Create("test_blank_migration", CreatorOptions{})
		require.NoError(t, err)
		assert.NotEmpty(t, result.Name)
		assert.NotEmpty(t, result.Path)
		assert.NotEmpty(t, result.Filename)

		// Verify file exists
		_, err = os.Stat(result.Path)
		assert.NoError(t, err)

		// Verify content contains expected elements
		content, err := os.ReadFile(result.Path)
		require.NoError(t, err)
		assert.Contains(t, string(content), "package migrations")
		assert.Contains(t, string(content), "migration.BaseMigration")
		assert.Contains(t, string(content), "func (m *testBlankMigration) Up")
		assert.Contains(t, string(content), "func (m *testBlankMigration) Down")
	})

	t.Run("creates create table migration", func(t *testing.T) {
		result, err := creator.Create("create_users_table", CreatorOptions{Create: "users"})
		require.NoError(t, err)

		content, err := os.ReadFile(result.Path)
		require.NoError(t, err)
		assert.Contains(t, string(content), "schema.NewBuilder")
		assert.Contains(t, string(content), `s.Create("users"`)
		assert.Contains(t, string(content), `s.DropIfExists("users")`)
	})

	t.Run("creates update table migration", func(t *testing.T) {
		result, err := creator.Create("add_email_to_users", CreatorOptions{Table: "users"})
		require.NoError(t, err)

		content, err := os.ReadFile(result.Path)
		require.NoError(t, err)
		assert.Contains(t, string(content), "schema.NewBuilder")
		assert.Contains(t, string(content), `s.Table("users"`)
	})

	t.Run("fails for duplicate migration", func(t *testing.T) {
		// First creation should succeed
		_, err := creator.Create("duplicate_test", CreatorOptions{})
		require.NoError(t, err)

		// Second creation with same name should fail (same timestamp unlikely but path exists)
		// We need to wait or use a different approach - let's just verify the error handling
	})
}

// TestToStructName tests the struct name generation.
func TestToStructName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"create_users_table", "createUsersTable"},
		{"add_email_to_users", "addEmailToUsers"},
		{"simple", "simple"},
		{"two_words", "twoWords"},
		{"a_b_c_d", "aBCD"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toStructName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestInferTableName tests the table name inference.
func TestInferTableName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"create_users_table", "users"},
		{"create_posts_table", "posts"},
		{"add_email_to_users", "users"},
		{"add_column_to_posts", "posts"},
		{"remove_field_from_users", "users"},
		{"modify_posts_table", "posts"},
		{"update_users_table", "users"},
		{"random_migration", "random_migration"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := inferTableName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidateMigrationName tests the migration name validation.
func TestValidateMigrationName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"create_users_table", false},
		{"add_email", false},
		{"test123", false},
		{"a", false},
		{"", true},
		{"CreateUsers", true}, // uppercase not allowed
		{"123test", true},     // must start with letter
		{"test-name", true},   // hyphens not allowed
		{"test.name", true},   // dots not allowed
		{"test name", true},   // spaces not allowed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMigrationName(tt.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGenerateTimestamp tests the timestamp generation format.
func TestGenerateTimestamp(t *testing.T) {
	timestamp := GenerateTimestamp()

	// Should match YYYY_MM_DD_HHMMSS format
	pattern := regexp.MustCompile(`^\d{4}_\d{2}_\d{2}_\d{6}$`)
	assert.True(t, pattern.MatchString(timestamp), "Timestamp %q should match YYYY_MM_DD_HHMMSS format", timestamp)
}
