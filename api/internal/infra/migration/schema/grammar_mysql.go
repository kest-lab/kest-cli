package schema

import (
	"fmt"
	"strings"
)

// MySQLGrammar generates MySQL-specific SQL statements
type MySQLGrammar struct {
	baseGrammar
}

// Compile generates all SQL statements from a Blueprint
func (g *MySQLGrammar) Compile(blueprint *Blueprint) []string {
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

// CompileCreate generates CREATE TABLE statement for MySQL
func (g *MySQLGrammar) CompileCreate(blueprint *Blueprint) string {
	var columns []string
	var primaryKeys []string

	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		columns = append(columns, colSQL)
		if col.primary {
			primaryKeys = append(primaryKeys, col.name)
		}
	}

	sql := fmt.Sprintf("CREATE TABLE `%s` (\n  %s", blueprint.table, strings.Join(columns, ",\n  "))

	if len(primaryKeys) > 0 {
		sql += fmt.Sprintf(",\n  PRIMARY KEY (`%s`)", strings.Join(primaryKeys, "`, `"))
	}

	sql += "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

	return sql
}

// compileColumn generates column definition SQL
func (g *MySQLGrammar) compileColumn(col *ColumnDefinition) string {
	colType := g.GetColumnType(col)
	modifiers := g.formatMySQLModifiers(col)

	if modifiers != "" {
		return fmt.Sprintf("`%s` %s %s", col.name, colType, modifiers)
	}
	return fmt.Sprintf("`%s` %s", col.name, colType)
}

// formatMySQLModifiers generates MySQL-specific column modifiers
func (g *MySQLGrammar) formatMySQLModifiers(col *ColumnDefinition) string {
	var parts []string

	if col.unsigned {
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

// CompileAdd generates ALTER TABLE ADD COLUMN statements
func (g *MySQLGrammar) CompileAdd(blueprint *Blueprint) []string {
	var statements []string
	for _, col := range blueprint.columns {
		colSQL := g.compileColumn(col)
		statements = append(statements, fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s", blueprint.table, colSQL))
	}
	return statements
}

// CompileChange generates ALTER TABLE MODIFY COLUMN statements
func (g *MySQLGrammar) CompileChange(blueprint *Blueprint) []string {
	// For now, return empty - change operations would need additional tracking
	return nil
}

// CompileDrop generates DROP TABLE statement
func (g *MySQLGrammar) CompileDrop(table string) string {
	return fmt.Sprintf("DROP TABLE `%s`", table)
}

// CompileDropIfExists generates DROP TABLE IF EXISTS statement
func (g *MySQLGrammar) CompileDropIfExists(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
}

// CompileRename generates RENAME TABLE statement
func (g *MySQLGrammar) CompileRename(from, to string) string {
	return fmt.Sprintf("RENAME TABLE `%s` TO `%s`", from, to)
}

// CompileDropColumn generates ALTER TABLE DROP COLUMN statement
func (g *MySQLGrammar) CompileDropColumn(table string, columns []string) string {
	var drops []string
	for _, col := range columns {
		drops = append(drops, fmt.Sprintf("DROP COLUMN `%s`", col))
	}
	return fmt.Sprintf("ALTER TABLE `%s` %s", table, strings.Join(drops, ", "))
}

// CompileRenameColumn generates ALTER TABLE RENAME COLUMN statement
func (g *MySQLGrammar) CompileRenameColumn(table, from, to string) string {
	return fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`", table, from, to)
}

// CompilePrimary generates ADD PRIMARY KEY statement
func (g *MySQLGrammar) CompilePrimary(table string, columns []string) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD PRIMARY KEY (`%s`)", table, strings.Join(columns, "`, `"))
}

// CompileUnique generates CREATE UNIQUE INDEX statement
func (g *MySQLGrammar) CompileUnique(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE UNIQUE INDEX `%s` ON `%s` (`%s`)", indexName, table, strings.Join(columns, "`, `"))
}

// CompileIndex generates CREATE INDEX statement
func (g *MySQLGrammar) CompileIndex(table string, columns []string, indexName string) string {
	return fmt.Sprintf("CREATE INDEX `%s` ON `%s` (`%s`)", indexName, table, strings.Join(columns, "`, `"))
}

// CompileForeign generates ADD FOREIGN KEY statement
func (g *MySQLGrammar) CompileForeign(table string, fk *ForeignKeyDefinition) string {
	sql := fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `%s_%s_foreign` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)",
		table, table, fk.column, fk.column, fk.on, fk.references)

	if fk.onDelete != "" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.onDelete)
	}
	if fk.onUpdate != "" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.onUpdate)
	}

	return sql
}

// CompileDropPrimary generates DROP PRIMARY KEY statement
func (g *MySQLGrammar) CompileDropPrimary(table string) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY", table)
}

// CompileDropUnique generates DROP INDEX statement for unique index
func (g *MySQLGrammar) CompileDropUnique(table, index string) string {
	return fmt.Sprintf("DROP INDEX `%s` ON `%s`", index, table)
}

// CompileDropIndex generates DROP INDEX statement
func (g *MySQLGrammar) CompileDropIndex(table, index string) string {
	return fmt.Sprintf("DROP INDEX `%s` ON `%s`", index, table)
}

// CompileDropForeign generates DROP FOREIGN KEY statement
func (g *MySQLGrammar) CompileDropForeign(table, index string) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", table, index)
}

// GetColumnType returns the MySQL type for a column definition
func (g *MySQLGrammar) GetColumnType(col *ColumnDefinition) string {
	switch col.colType {
	case "bigInteger":
		return "BIGINT"
	case "integer":
		return "INT"
	case "smallInteger":
		return "SMALLINT"
	case "tinyInteger":
		return "TINYINT"
	case "string":
		if col.length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.length)
		}
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "mediumText":
		return "MEDIUMTEXT"
	case "longText":
		return "LONGTEXT"
	case "boolean":
		return "TINYINT(1)"
	case "timestamp":
		return "TIMESTAMP"
	case "datetime":
		return "DATETIME"
	case "date":
		return "DATE"
	case "time":
		return "TIME"
	case "json":
		return "JSON"
	case "binary":
		return "BLOB"
	case "float":
		return "FLOAT"
	case "double":
		return "DOUBLE"
	case "decimal":
		if col.precision > 0 && col.scale > 0 {
			return fmt.Sprintf("DECIMAL(%d,%d)", col.precision, col.scale)
		}
		return "DECIMAL(8,2)"
	case "uuid":
		return "CHAR(36)"
	default:
		return "VARCHAR(255)"
	}
}
