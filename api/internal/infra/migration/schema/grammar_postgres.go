package schema

import (
	"fmt"
	"strings"
)

// PostgresGrammar generates PostgreSQL-specific SQL statements
type PostgresGrammar struct {
	baseGrammar
}

// Compile generates all SQL statements from a Blueprint
func (g *PostgresGrammar) Compile(blueprint *Blueprint) []string {
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
			if cols, ok := cmd.params["columns"].([]string); ok {
				statements = append(statements, g.CompilePrimary(blueprint.table, cols))
			}
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
			if fk, ok := cmd.params["definition"].(*ForeignKeyDefinition); ok {
				statements = append(statements, g.CompileForeign(blueprint.table, fk))
			}
		case "dropColumn":
			if cols, ok := cmd.params["columns"].([]string); ok {
				statements = append(statements, g.CompileDropColumn(blueprint.table, cols))
			}
		case "renameColumn":
			from, _ := cmd.params["from"].(string)
			to, _ := cmd.params["to"].(string)
			statements = append(statements, g.CompileRenameColumn(blueprint.table, from, to))
		case "dropPrimary":
			statements = append(statements, g.CompileDropPrimary(blueprint.table))
		case "dropUnique":
			if idx, ok := cmd.params["index"].(string); ok {
				statements = append(statements, g.CompileDropUnique(blueprint.table, idx))
			}
		case "dropIndex":
			if idx, ok := cmd.params["index"].(string); ok {
				statements = append(statements, g.CompileDropIndex(blueprint.table, idx))
			}
		case "dropForeign":
			if idx, ok := cmd.params["index"].(string); ok {
				statements = append(statements, g.CompileDropForeign(blueprint.table, idx))
			}
		}
	}

	return statements
}

// CompileCreate generates CREATE TABLE statement for PostgreSQL
func (g *PostgresGrammar) CompileCreate(blueprint *Blueprint) string {
	var columns []string
	var primaryKeys []string

	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		columns = append(columns, colSQL)
		if col.primary {
			primaryKeys = append(primaryKeys, col.name)
		}
	}

	sql := fmt.Sprintf("CREATE TABLE \"%s\" (\n  %s", blueprint.table, strings.Join(columns, ",\n  "))

	if len(primaryKeys) > 0 {
		sql += fmt.Sprintf(",\n  PRIMARY KEY (\"%s\")", strings.Join(primaryKeys, "\", \""))
	}

	sql += "\n)"

	return sql
}

// compileColumn generates column definition SQL
func (g *PostgresGrammar) compileColumn(col *ColumnDefinition) string {
	colType := g.GetColumnType(col)
	modifiers := g.formatPostgresModifiers(col)

	if modifiers != "" {
		return fmt.Sprintf("\"%s\" %s %s", col.name, colType, modifiers)
	}
	return fmt.Sprintf("\"%s\" %s", col.name, colType)
}

// formatPostgresModifiers generates PostgreSQL-specific column modifiers
func (g *PostgresGrammar) formatPostgresModifiers(col *ColumnDefinition) string {
	var parts []string

	if !col.nullable {
		parts = append(parts, "NOT NULL")
	}

	if col.defaultValue != nil {
		parts = append(parts, fmt.Sprintf("DEFAULT %s", g.formatDefault(col.defaultValue)))
	}

	return strings.Join(parts, " ")
}

// CompileAdd generates ALTER TABLE ADD COLUMN statements
func (g *PostgresGrammar) CompileAdd(blueprint *Blueprint) []string {
	var statements []string
	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		statements = append(statements, fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s", blueprint.table, colSQL))
	}
	return statements
}

// CompileChange generates ALTER TABLE ALTER COLUMN statements
func (g *PostgresGrammar) CompileChange(blueprint *Blueprint) []string {
	// For now, return empty - change operations would need additional tracking
	return nil
}

// CompileDrop generates DROP TABLE statement
func (g *PostgresGrammar) CompileDrop(table string) string {
	return fmt.Sprintf("DROP TABLE \"%s\"", table)
}

// CompileDropIfExists generates DROP TABLE IF EXISTS statement
func (g *PostgresGrammar) CompileDropIfExists(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS \"%s\"", table)
}

// CompileRename generates ALTER TABLE RENAME statement
func (g *PostgresGrammar) CompileRename(from, to string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\"", from, to)
}

// CompileDropColumn generates ALTER TABLE DROP COLUMN statement
func (g *PostgresGrammar) CompileDropColumn(table string, columns []string) string {
	var drops []string
	for _, col := range columns {
		drops = append(drops, fmt.Sprintf("DROP COLUMN \"%s\"", col))
	}
	return fmt.Sprintf("ALTER TABLE \"%s\" %s", table, strings.Join(drops, ", "))
}

// CompileRenameColumn generates ALTER TABLE RENAME COLUMN statement
func (g *PostgresGrammar) CompileRenameColumn(table, from, to string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\"", table, from, to)
}

// CompilePrimary generates ADD PRIMARY KEY statement
func (g *PostgresGrammar) CompilePrimary(table string, columns []string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" ADD PRIMARY KEY (\"%s\")", table, strings.Join(columns, "\", \""))
}

// CompileUnique generates CREATE UNIQUE INDEX statement
func (g *PostgresGrammar) CompileUnique(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\" (\"%s\")", indexName, table, strings.Join(columns, "\", \""))
}

// CompileIndex generates CREATE INDEX statement
func (g *PostgresGrammar) CompileIndex(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (\"%s\")", indexName, table, strings.Join(columns, "\", \""))
}

// CompileForeign generates ADD FOREIGN KEY statement
func (g *PostgresGrammar) CompileForeign(table string, fk *ForeignKeyDefinition) string {
	sql := fmt.Sprintf("ALTER TABLE \"%s\" ADD CONSTRAINT \"%s_%s_foreign\" FOREIGN KEY (\"%s\") REFERENCES \"%s\" (\"%s\")",
		table, table, fk.column, fk.column, fk.on, fk.references)

	if fk.onDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.onDelete)
	}
	if fk.onUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.onUpdate)
	}

	return sql
}

// CompileDropPrimary generates DROP CONSTRAINT statement for primary key
func (g *PostgresGrammar) CompileDropPrimary(table string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP CONSTRAINT \"%s_pkey\"", table, table)
}

// CompileDropUnique generates DROP INDEX statement
func (g *PostgresGrammar) CompileDropUnique(table, index string) string {
	return fmt.Sprintf("DROP INDEX \"%s\"", index)
}

// CompileDropIndex generates DROP INDEX statement
func (g *PostgresGrammar) CompileDropIndex(table, index string) string {
	return fmt.Sprintf("DROP INDEX \"%s\"", index)
}

// CompileDropForeign generates DROP CONSTRAINT statement for foreign key
func (g *PostgresGrammar) CompileDropForeign(table, index string) string {
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP CONSTRAINT \"%s\"", table, index)
}

// GetColumnType returns the PostgreSQL type for a column definition
func (g *PostgresGrammar) GetColumnType(col *ColumnDefinition) string {
	switch col.colType {
	case "bigInteger":
		if col.autoIncrement {
			return "BIGSERIAL"
		}
		return "BIGINT"
	case "integer":
		if col.autoIncrement {
			return "SERIAL"
		}
		return "INTEGER"
	case "smallInteger":
		if col.autoIncrement {
			return "SMALLSERIAL"
		}
		return "SMALLINT"
	case "tinyInteger":
		return "SMALLINT"
	case "string":
		if col.length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.length)
		}
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "mediumText":
		return "TEXT"
	case "longText":
		return "TEXT"
	case "boolean":
		return "BOOLEAN"
	case "timestamp":
		return "TIMESTAMP"
	case "datetime":
		return "TIMESTAMP"
	case "date":
		return "DATE"
	case "time":
		return "TIME"
	case "json":
		return "JSONB"
	case "binary":
		return "BYTEA"
	case "float":
		return "REAL"
	case "double":
		return "DOUBLE PRECISION"
	case "decimal":
		if col.precision > 0 && col.scale > 0 {
			return fmt.Sprintf("DECIMAL(%d,%d)", col.precision, col.scale)
		}
		return "DECIMAL(8,2)"
	case "uuid":
		return "UUID"
	default:
		return "VARCHAR(255)"
	}
}
