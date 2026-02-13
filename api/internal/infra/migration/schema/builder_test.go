package schema

import (
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

// tableNameGen generates valid table names
func tableNameGen() gopter.Gen {
	return gen.RegexMatch(`[a-z][a-z0-9_]{2,20}`)
}

// Feature: migration-system, Property 13: Schema HasTable Reflects Existence
// For any table name, HasTable() should return true if and only if the table exists in the database.
func TestProperty13_SchemaHasTableReflectsExistence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("HasTable returns true iff table exists", prop.ForAll(
		func(tableName string) bool {
			// Setup fresh database for each test
			db := setupTestDB(t)
			builder := NewBuilder(db)

			// Initially table should not exist
			if builder.HasTable(tableName) {
				return false
			}

			// Create the table
			err := builder.Create(tableName, func(bp *Blueprint) {
				bp.ID()
				bp.String("name")
			})
			if err != nil {
				return false
			}

			// Now table should exist
			if !builder.HasTable(tableName) {
				return false
			}

			// Drop the table
			err = builder.Drop(tableName)
			if err != nil {
				return false
			}

			// Table should not exist anymore
			if builder.HasTable(tableName) {
				return false
			}

			return true
		},
		tableNameGen(),
	))

	properties.TestingRun(t)
}

// Unit tests for Schema Builder

func TestBuilder_Create(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	err := builder.Create("users", func(bp *Blueprint) {
		bp.ID()
		bp.String("name")
		bp.String("email", 100)
		bp.Text("bio")
		bp.Boolean("active").Default(true)
		bp.Timestamps()
	})
	require.NoError(t, err)

	// Verify table exists
	assert.True(t, builder.HasTable("users"))

	// Verify columns exist
	assert.True(t, builder.HasColumn("users", "id"))
	assert.True(t, builder.HasColumn("users", "name"))
	assert.True(t, builder.HasColumn("users", "email"))
	assert.True(t, builder.HasColumn("users", "bio"))
	assert.True(t, builder.HasColumn("users", "active"))
	assert.True(t, builder.HasColumn("users", "created_at"))
	assert.True(t, builder.HasColumn("users", "updated_at"))
}

func TestBuilder_Table_AddColumn(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	// Create initial table
	err := builder.Create("posts", func(bp *Blueprint) {
		bp.ID()
		bp.String("title")
	})
	require.NoError(t, err)

	// Add a new column
	err = builder.Table("posts", func(bp *Blueprint) {
		bp.Text("content")
	})
	require.NoError(t, err)

	// Verify new column exists
	assert.True(t, builder.HasColumn("posts", "content"))
}

func TestBuilder_Drop(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	// Create table
	err := builder.Create("temp_table", func(bp *Blueprint) {
		bp.ID()
	})
	require.NoError(t, err)
	assert.True(t, builder.HasTable("temp_table"))

	// Drop table
	err = builder.Drop("temp_table")
	require.NoError(t, err)
	assert.False(t, builder.HasTable("temp_table"))
}

func TestBuilder_DropIfExists(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	// Should not error when table doesn't exist
	err := builder.DropIfExists("non_existent_table")
	assert.NoError(t, err)

	// Create and drop
	err = builder.Create("temp_table", func(bp *Blueprint) {
		bp.ID()
	})
	require.NoError(t, err)

	err = builder.DropIfExists("temp_table")
	assert.NoError(t, err)
	assert.False(t, builder.HasTable("temp_table"))
}

func TestBuilder_Rename(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	// Create table
	err := builder.Create("old_name", func(bp *Blueprint) {
		bp.ID()
	})
	require.NoError(t, err)

	// Rename table
	err = builder.Rename("old_name", "new_name")
	require.NoError(t, err)

	assert.False(t, builder.HasTable("old_name"))
	assert.True(t, builder.HasTable("new_name"))
}

func TestBuilder_HasColumn(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	err := builder.Create("test_table", func(bp *Blueprint) {
		bp.ID()
		bp.String("existing_column")
	})
	require.NoError(t, err)

	assert.True(t, builder.HasColumn("test_table", "existing_column"))
	assert.False(t, builder.HasColumn("test_table", "non_existent_column"))
}

func TestBuilder_GetColumnListing(t *testing.T) {
	db := setupTestDB(t)
	builder := NewBuilder(db)

	err := builder.Create("columns_test", func(bp *Blueprint) {
		bp.ID()
		bp.String("name")
		bp.Integer("age")
	})
	require.NoError(t, err)

	columns, err := builder.GetColumnListing("columns_test")
	require.NoError(t, err)

	assert.Contains(t, columns, "id")
	assert.Contains(t, columns, "name")
	assert.Contains(t, columns, "age")
}

// Unit tests for Blueprint

func TestBlueprint_ColumnTypes(t *testing.T) {
	bp := NewBlueprint("test")
	bp.Create()

	// Test various column types
	bp.ID()
	bp.String("name")
	bp.Text("description")
	bp.Integer("count")
	bp.BigInteger("big_count")
	bp.Boolean("active")
	bp.Timestamp("created_at")
	bp.JSON("metadata")

	assert.Len(t, bp.GetColumns(), 8)
	assert.True(t, bp.IsCreating())
}

func TestBlueprint_Timestamps(t *testing.T) {
	bp := NewBlueprint("test")
	bp.Create()
	bp.Timestamps()

	columns := bp.GetColumns()
	assert.Len(t, columns, 2)

	// Find created_at and updated_at
	var createdAt, updatedAt *ColumnDefinition
	for _, col := range columns {
		if col.GetName() == "created_at" {
			createdAt = col
		}
		if col.GetName() == "updated_at" {
			updatedAt = col
		}
	}

	assert.NotNil(t, createdAt)
	assert.NotNil(t, updatedAt)
	assert.True(t, createdAt.IsNullable())
	assert.True(t, updatedAt.IsNullable())
}

func TestBlueprint_SoftDeletes(t *testing.T) {
	bp := NewBlueprint("test")
	bp.Create()
	bp.SoftDeletes()

	columns := bp.GetColumns()
	assert.Len(t, columns, 1)
	assert.Equal(t, "deleted_at", columns[0].GetName())
	assert.True(t, columns[0].IsNullable())
}

func TestBlueprint_SoftDeletesCustomColumn(t *testing.T) {
	bp := NewBlueprint("test")
	bp.Create()
	bp.SoftDeletes("removed_at")

	columns := bp.GetColumns()
	assert.Len(t, columns, 1)
	assert.Equal(t, "removed_at", columns[0].GetName())
}

func TestBlueprint_Indexes(t *testing.T) {
	bp := NewBlueprint("test")
	bp.Create()
	bp.String("email")
	bp.Unique("email")
	bp.Index("email")

	commands := bp.GetCommands()
	// Should have: create, unique, index
	assert.Len(t, commands, 3)
}

func TestBlueprint_ForeignKey(t *testing.T) {
	bp := NewBlueprint("posts")
	bp.Create()
	bp.UnsignedBigInteger("user_id")
	bp.Foreign("user_id").References("id").On("users").Cascade()

	commands := bp.GetCommands()
	// Should have: create, foreign
	assert.Len(t, commands, 2)
}

// Unit tests for ColumnDefinition

func TestColumnDefinition_Chaining(t *testing.T) {
	bp := NewBlueprint("test")
	col := bp.String("name").Nullable().Default("unknown").Comment("User name")

	assert.True(t, col.IsNullable())
	assert.Equal(t, "unknown", col.GetDefault())
	assert.Equal(t, "User name", col.GetComment())
}

func TestColumnDefinition_Unsigned(t *testing.T) {
	bp := NewBlueprint("test")
	col := bp.Integer("count").Unsigned()

	assert.True(t, col.IsUnsigned())
}

func TestColumnDefinition_AutoIncrement(t *testing.T) {
	bp := NewBlueprint("test")
	col := bp.BigInteger("id").AutoIncrement().Primary()

	assert.True(t, col.IsAutoIncrement())
	assert.True(t, col.IsPrimary())
}

// Unit tests for ForeignKeyDefinition

func TestForeignKeyDefinition_Chaining(t *testing.T) {
	fk := &ForeignKeyDefinition{column: "user_id"}
	fk.References("id").On("users").OnDelete("CASCADE").OnUpdate("SET NULL")

	assert.Equal(t, "user_id", fk.GetColumn())
	assert.Equal(t, "id", fk.GetReferences())
	assert.Equal(t, "users", fk.GetOn())
	assert.Equal(t, "CASCADE", fk.GetOnDelete())
	assert.Equal(t, "SET NULL", fk.GetOnUpdate())
}

func TestForeignKeyDefinition_Cascade(t *testing.T) {
	fk := &ForeignKeyDefinition{column: "user_id"}
	fk.Cascade()

	assert.Equal(t, "CASCADE", fk.GetOnDelete())
	assert.Equal(t, "CASCADE", fk.GetOnUpdate())
}

func TestForeignKeyDefinition_Constrained(t *testing.T) {
	fk := &ForeignKeyDefinition{column: "user_id"}
	fk.Constrained("users")

	assert.Equal(t, "users", fk.GetOn())
	assert.Equal(t, "id", fk.GetReferences())
}

// Unit tests for Grammar implementations

func TestMySQLGrammar_CompileCreate(t *testing.T) {
	grammar := &MySQLGrammar{}
	bp := NewBlueprint("users")
	bp.Create()
	bp.ID()
	bp.String("name")

	sql := grammar.CompileCreate(bp)

	assert.Contains(t, sql, "CREATE TABLE `users`")
	assert.Contains(t, sql, "`id` BIGINT")
	assert.Contains(t, sql, "`name` VARCHAR(255)")
	assert.Contains(t, sql, "PRIMARY KEY")
}

func TestPostgresGrammar_CompileCreate(t *testing.T) {
	grammar := &PostgresGrammar{}
	bp := NewBlueprint("users")
	bp.Create()
	bp.ID()
	bp.String("name")

	sql := grammar.CompileCreate(bp)

	assert.Contains(t, sql, "CREATE TABLE \"users\"")
	assert.Contains(t, sql, "\"id\" BIGSERIAL")
	assert.Contains(t, sql, "\"name\" VARCHAR(255)")
	assert.Contains(t, sql, "PRIMARY KEY")
}

func TestSQLiteGrammar_CompileCreate(t *testing.T) {
	grammar := &SQLiteGrammar{}
	bp := NewBlueprint("users")
	bp.Create()
	bp.ID()
	bp.String("name")

	sql := grammar.CompileCreate(bp)

	assert.Contains(t, sql, "CREATE TABLE \"users\"")
	assert.Contains(t, sql, "\"id\" INTEGER")
	assert.Contains(t, sql, "PRIMARY KEY AUTOINCREMENT")
	assert.Contains(t, sql, "\"name\" VARCHAR(255)")
}

func TestNewGrammar_Factory(t *testing.T) {
	tests := []struct {
		dialect  string
		expected string
	}{
		{"mysql", "*schema.MySQLGrammar"},
		{"postgres", "*schema.PostgresGrammar"},
		{"postgresql", "*schema.PostgresGrammar"},
		{"sqlite", "*schema.SQLiteGrammar"},
		{"sqlite3", "*schema.SQLiteGrammar"},
		{"unknown", "*schema.SQLiteGrammar"}, // defaults to SQLite
	}

	for _, tt := range tests {
		t.Run(tt.dialect, func(t *testing.T) {
			grammar := NewGrammar(tt.dialect)
			assert.NotNil(t, grammar)
		})
	}
}

func TestGrammar_CompileDrop(t *testing.T) {
	tests := []struct {
		name     string
		grammar  Grammar
		expected string
	}{
		{"MySQL", &MySQLGrammar{}, "DROP TABLE `users`"},
		{"Postgres", &PostgresGrammar{}, "DROP TABLE \"users\""},
		{"SQLite", &SQLiteGrammar{}, "DROP TABLE \"users\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.grammar.CompileDrop("users")
			assert.Equal(t, tt.expected, sql)
		})
	}
}

func TestGrammar_CompileDropIfExists(t *testing.T) {
	tests := []struct {
		name     string
		grammar  Grammar
		expected string
	}{
		{"MySQL", &MySQLGrammar{}, "DROP TABLE IF EXISTS `users`"},
		{"Postgres", &PostgresGrammar{}, "DROP TABLE IF EXISTS \"users\""},
		{"SQLite", &SQLiteGrammar{}, "DROP TABLE IF EXISTS \"users\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.grammar.CompileDropIfExists("users")
			assert.Equal(t, tt.expected, sql)
		})
	}
}

func TestGrammar_CompileRename(t *testing.T) {
	tests := []struct {
		name     string
		grammar  Grammar
		expected string
	}{
		{"MySQL", &MySQLGrammar{}, "RENAME TABLE `old` TO `new`"},
		{"Postgres", &PostgresGrammar{}, "ALTER TABLE \"old\" RENAME TO \"new\""},
		{"SQLite", &SQLiteGrammar{}, "ALTER TABLE \"old\" RENAME TO \"new\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.grammar.CompileRename("old", "new")
			assert.Equal(t, tt.expected, sql)
		})
	}
}
