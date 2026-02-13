package schema

import (
	"fmt"
	"strings"
)

// Grammar defines the contract for compiling Blueprint to SQL statements.
// Each database dialect implements this interface to generate dialect-specific SQL.
type Grammar interface {
	// Compile generates SQL statements from a Blueprint
	Compile(blueprint *Blueprint) []string

	// CompileCreate generates CREATE TABLE statement
	CompileCreate(blueprint *Blueprint) string

	// CompileAdd generates ALTER TABLE ADD COLUMN statements
	CompileAdd(blueprint *Blueprint) []string

	// CompileChange generates ALTER TABLE MODIFY/ALTER COLUMN statements
	CompileChange(blueprint *Blueprint) []string

	// CompileDrop generates DROP TABLE statement
	CompileDrop(table string) string

	// CompileDropIfExists generates DROP TABLE IF EXISTS statement
	CompileDropIfExists(table string) string

	// CompileRename generates RENAME TABLE statement
	CompileRename(from, to string) string

	// CompileDropColumn generates ALTER TABLE DROP COLUMN statement
	CompileDropColumn(table string, columns []string) string

	// CompileRenameColumn generates ALTER TABLE RENAME COLUMN statement
	CompileRenameColumn(table, from, to string) string

	// CompilePrimary generates ADD PRIMARY KEY statement
	CompilePrimary(table string, columns []string) string

	// CompileUnique generates CREATE UNIQUE INDEX statement
	CompileUnique(table string, columns []string, indexName string) string

	// CompileIndex generates CREATE INDEX statement
	CompileIndex(table string, columns []string, indexName string) string

	// CompileForeign generates ADD FOREIGN KEY statement
	CompileForeign(table string, fk *ForeignKeyDefinition) string

	// CompileDropPrimary generates DROP PRIMARY KEY statement
	CompileDropPrimary(table string) string

	// CompileDropUnique generates DROP UNIQUE INDEX statement
	CompileDropUnique(table, index string) string

	// CompileDropIndex generates DROP INDEX statement
	CompileDropIndex(table, index string) string

	// CompileDropForeign generates DROP FOREIGN KEY statement
	CompileDropForeign(table, index string) string

	// GetColumnType returns the SQL type for a column definition
	GetColumnType(column *ColumnDefinition) string
}

// NewGrammar creates a Grammar implementation for the specified dialect.
// Supported dialects: "mysql", "postgres", "sqlite"
func NewGrammar(dialect string) Grammar {
	switch strings.ToLower(dialect) {
	case "mysql":
		return &MySQLGrammar{}
	case "postgres", "postgresql":
		return &PostgresGrammar{}
	case "sqlite", "sqlite3":
		return &SQLiteGrammar{}
	default:
		// Default to SQLite for unknown dialects
		return &SQLiteGrammar{}
	}
}

// baseGrammar provides common functionality for all grammar implementations
type baseGrammar struct{}

// formatColumnModifiers generates common column modifiers
func (g *baseGrammar) formatColumnModifiers(col *ColumnDefinition, unsigned bool) string {
	var parts []string

	if unsigned && col.unsigned {
		parts = append(parts, "UNSIGNED")
	}

	if !col.nullable {
		parts = append(parts, "NOT NULL")
	} else {
		parts = append(parts, "NULL")
	}

	if col.defaultValue != nil {
		parts = append(parts, fmt.Sprintf("DEFAULT %s", g.formatDefault(col.defaultValue)))
	}

	if col.autoIncrement {
		parts = append(parts, "AUTO_INCREMENT")
	}

	if col.comment != "" {
		parts = append(parts, fmt.Sprintf("COMMENT '%s'", col.comment))
	}

	return strings.Join(parts, " ")
}

// formatDefault formats a default value for SQL
func (g *baseGrammar) formatDefault(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// generateIndexName creates a standard index name
func generateIndexName(table string, columns []string, indexType string) string {
	return fmt.Sprintf("%s_%s_%s", table, strings.Join(columns, "_"), indexType)
}
