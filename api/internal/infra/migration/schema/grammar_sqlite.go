package schema

import (
	"fmt"
	"strings"
)

// SQLiteGrammar generates SQLite-specific SQL statements
type SQLiteGrammar struct {
	baseGrammar
}

// Compile generates all SQL statements from a Blueprint
func (g *SQLiteGrammar) Compile(blueprint *Blueprint) []string {
	var statements []string

	if blueprint.creating {
		statements = append(statements, g.CompileCreate(blueprint))
	} else {
		// Handle modifications
		statements = append(statements, g.CompileAdd(blueprint)...)
		statements = append(statements, g.CompileChange(blueprint)...)
	}

	// Process commands
	for _, cmd := range blueprint.commands {
		switch cmd.name {
		case "primary":
			// SQLite handles primary keys in CREATE TABLE
			// For ALTER TABLE, this is not directly supported
		case "unique":
			if cols, ok := cmd.params["columns"].([]string); ok {
				indexName := generateIndexName(blueprint.table, cols, "unique")
				statements = append(statements, g.CompileUnique(blueprint.table, cols, indexName))
			}
		case "index":
			if cols, ok := cmd.params["columns"].([]string); ok {
				indexName := generateIndexName(blueprint.table, cols, "index")
				statements = append(statements, g.CompileIndex(blueprint.table, cols, indexName))
			}
		case "foreign":
			// SQLite foreign keys are defined in CREATE TABLE
			// For ALTER TABLE, this requires table recreation
		case "dropColumn":
			// SQLite doesn't support DROP COLUMN directly in older versions
			// Modern SQLite (3.35.0+) supports it
			if cols, ok := cmd.params["columns"].([]string); ok {
				statements = append(statements, g.CompileDropColumn(blueprint.table, cols))
			}
		case "renameColumn":
			from, _ := cmd.params["from"].(string)
			to, _ := cmd.params["to"].(string)
			statements = append(statements, g.CompileRenameColumn(blueprint.table, from, to))
		case "dropUnique":
			if idx, ok := cmd.params["index"].(string); ok {
				statements = append(statements, g.CompileDropUnique(blueprint.table, idx))
			}
		case "dropIndex":
			if idx, ok := cmd.params["index"].(string); ok {
				statements = append(statements, g.CompileDropIndex(blueprint.table, idx))
			}
		}
	}

	return statements
}

// CompileCreate generates CREATE TABLE statement for SQLite
func (g *SQLiteGrammar) CompileCreate(blueprint *Blueprint) string {
	var columns []string
	var primaryKeys []string
	var foreignKeys []string

	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		columns = append(columns, colSQL)
		if col.primary && !col.autoIncrement {
			primaryKeys = append(primaryKeys, col.name)
		}
	}

	// Collect foreign keys from commands
	for _, cmd := range blueprint.commands {
		if cmd.name == "foreign" {
			if fk, ok := cmd.params["definition"].(*ForeignKeyDefinition); ok {
				fkSQL := g.compileForeignKeyInline(fk)
				foreignKeys = append(foreignKeys, fkSQL)
			}
		}
	}

	sql := fmt.Sprintf("CREATE TABLE \"%s\" (\n  %s", blueprint.table, strings.Join(columns, ",\n  "))

	if len(primaryKeys) > 0 {
		sql += fmt.Sprintf(",\n  PRIMARY KEY (\"%s\")", strings.Join(primaryKeys, "\", \""))
	}

	for _, fk := range foreignKeys {
		sql += ",\n  " + fk
	}

	sql += "\n)"

	return sql
}

// compileColumn generates column definition SQL
func (g *SQLiteGrammar) compileColumn(col *ColumnDefinition) string {
	colType := g.GetColumnType(col)
	modifiers := g.formatSQLiteModifiers(col)

	if modifiers != "" {
		return fmt.Sprintf("\"%s\" %s %s", col.name, colType, modifiers)
	}
	return fmt.Sprintf("\"%s\" %s", col.name, colType)
}

// formatSQLiteModifiers generates SQLite-specific column modifiers
func (g *SQLiteGrammar) formatSQLiteModifiers(col *ColumnDefinition) string {
	var parts []string

	// SQLite uses PRIMARY KEY AUTOINCREMENT for auto-increment
	if col.primary && col.autoIncrement {
		parts = append(parts, "PRIMARY KEY AUTOINCREMENT")
	}

	if !col.nullable && !col.autoIncrement {
		parts = append(parts, "NOT NULL")
	}

	if col.defaultValue != nil {
		parts = append(parts, fmt.Sprintf("DEFAULT %s", g.formatDefault(col.defaultValue)))
	}

	return strings.Join(parts, " ")
}

// compileForeignKeyInline generates inline foreign key constraint for CREATE TABLE
func (g *SQLiteGrammar) compileForeignKeyInline(fk *ForeignKeyDefinition) string {
	sql := fmt.Sprintf("FOREIGN KEY (\"%s\") REFERENCES \"%s\" (\"%s\")",
		fk.column, fk.on, fk.references)

	if fk.onDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.onDelete)
	}
	if fk.onUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.onUpdate)
	}

	return sql
}

// CompileAdd generates ALTER TABLE ADD COLUMN statements
func (g *SQLiteGrammar) CompileAdd(blueprint *Blueprint) []string {
	var statements []string
	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		statements = append(statements, fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s", blueprint.table, colSQL))
	}
	return statements
}

// CompileChange generates ALTER TABLE statements for column changes
// Note: SQLite has limited ALTER TABLE support
func (g *SQLiteGrammar) CompileChange(blueprint *Blueprint) []string {
	// SQLite doesn't support modifying columns directly
	// Would require table recreation
	return nil
}

// CompileDrop generates DROP TABLE statement
func (g *SQLiteGrammar) CompileDrop(table string) string {
	return fmt.Sprintf("DROP TABLE \"%s\"", table)
}

// CompileDropIfExists generates DROP TABLE IF EXISTS statement
func (g *SQLiteGrammar) CompileDropIfExists(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS \"%s\"", table)
}

// CompileRename generates ALTER TABLE RENAME statement
func (g *SQLiteGrammar) CompileRename(from, to string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\"", from, to)
}

// CompileDropColumn generates ALTER TABLE DROP COLUMN statement
// Note: Requires SQLite 3.35.0+
func (g *SQLiteGrammar) CompileDropColumn(table string, columns []string) string {
	var statements []string
	for _, col := range columns {
		statements = append(statements, fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\"", table, col))
	}
	// Return first statement; multiple drops need separate calls
	if len(statements) > 0 {
		return statements[0]
	}
	return ""
}

// CompileRenameColumn generates ALTER TABLE RENAME COLUMN statement
// Note: Requires SQLite 3.25.0+
func (g *SQLiteGrammar) CompileRenameColumn(table, from, to string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\"", table, from, to)
}

// CompilePrimary generates ADD PRIMARY KEY statement
// Note: SQLite doesn't support adding primary key after table creation
func (g *SQLiteGrammar) CompilePrimary(table string, columns []string) string {
	// SQLite doesn't support ALTER TABLE ADD PRIMARY KEY
	// Would require table recreation
	return ""
}

// CompileUnique generates CREATE UNIQUE INDEX statement
func (g *SQLiteGrammar) CompileUnique(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\" (\"%s\")", indexName, table, strings.Join(columns, "\", \""))
}

// CompileIndex generates CREATE INDEX statement
func (g *SQLiteGrammar) CompileIndex(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (\"%s\")", indexName, table, strings.Join(columns, "\", \""))
}

// CompileForeign generates foreign key constraint
// Note: SQLite foreign keys must be defined in CREATE TABLE
func (g *SQLiteGrammar) CompileForeign(table string, fk *ForeignKeyDefinition) string {
	// SQLite doesn't support ALTER TABLE ADD FOREIGN KEY
	// Foreign keys must be defined in CREATE TABLE
	return ""
}

// CompileDropPrimary generates statement to drop primary key
// Note: SQLite doesn't support this directly
func (g *SQLiteGrammar) CompileDropPrimary(table string) string {
	// SQLite doesn't support dropping primary key
	// Would require table recreation
	return ""
}

// CompileDropUnique generates DROP INDEX statement
func (g *SQLiteGrammar) CompileDropUnique(table, index string) string {
	return fmt.Sprintf("DROP INDEX \"%s\"", index)
}

// CompileDropIndex generates DROP INDEX statement
func (g *SQLiteGrammar) CompileDropIndex(table, index string) string {
	return fmt.Sprintf("DROP INDEX \"%s\"", index)
}

// CompileDropForeign generates statement to drop foreign key
// Note: SQLite doesn't support this directly
func (g *SQLiteGrammar) CompileDropForeign(table, index string) string {
	// SQLite doesn't support dropping foreign keys
	// Would require table recreation
	return ""
}

// GetColumnType returns the SQLite type for a column definition
func (g *SQLiteGrammar) GetColumnType(col *ColumnDefinition) string {
	switch col.colType {
	case "bigInteger", "integer", "smallInteger", "tinyInteger":
		return "INTEGER"
	case "string":
		if col.length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.length)
		}
		return "VARCHAR(255)"
	case "text", "mediumText", "longText":
		return "TEXT"
	case "boolean":
		return "INTEGER"
	case "timestamp", "datetime":
		return "DATETIME"
	case "date":
		return "DATE"
	case "time":
		return "TIME"
	case "json":
		return "TEXT"
	case "binary":
		return "BLOB"
	case "float", "double":
		return "REAL"
	case "decimal":
		return "NUMERIC"
	case "uuid":
		return "VARCHAR(36)"
	default:
		return "TEXT"
	}
}
