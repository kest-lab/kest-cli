package schema

import (
	"gorm.io/gorm"
)

// Builder provides a fluent interface for schema operations.
// It wraps GORM's database connection and uses Grammar implementations
// to generate dialect-specific SQL.
type Builder struct {
	db      *gorm.DB
	grammar Grammar
}

// NewBuilder creates a new schema builder for the given database connection.
// It automatically selects the appropriate grammar based on the database dialect.
func NewBuilder(db *gorm.DB) *Builder {
	return &Builder{
		db:      db,
		grammar: NewGrammar(db.Dialector.Name()),
	}
}

// NewBuilderWithGrammar creates a new schema builder with a specific grammar.
// Use this when you need to override the default grammar selection.
func NewBuilderWithGrammar(db *gorm.DB, grammar Grammar) *Builder {
	return &Builder{
		db:      db,
		grammar: grammar,
	}
}

// Create creates a new table using the blueprint callback.
// The callback receives a Blueprint that can be used to define columns and indexes.
func (b *Builder) Create(table string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(table)
	blueprint.Create()
	callback(blueprint)

	statements := blueprint.ToSQL(b.grammar)
	for _, sql := range statements {
		if sql == "" {
			continue
		}
		if err := b.db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// Table modifies an existing table using the blueprint callback.
// The callback receives a Blueprint that can be used to add/modify columns and indexes.
func (b *Builder) Table(table string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(table)
	callback(blueprint)

	statements := blueprint.ToSQL(b.grammar)
	for _, sql := range statements {
		if sql == "" {
			continue
		}
		if err := b.db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// Drop drops a table.
func (b *Builder) Drop(table string) error {
	sql := b.grammar.CompileDrop(table)
	return b.db.Exec(sql).Error
}

// DropIfExists drops a table if it exists.
func (b *Builder) DropIfExists(table string) error {
	sql := b.grammar.CompileDropIfExists(table)
	return b.db.Exec(sql).Error
}

// Rename renames a table from one name to another.
func (b *Builder) Rename(from, to string) error {
	sql := b.grammar.CompileRename(from, to)
	return b.db.Exec(sql).Error
}

// HasTable checks if a table exists in the database.
func (b *Builder) HasTable(table string) bool {
	return b.db.Migrator().HasTable(table)
}

// HasColumn checks if a column exists in a table.
func (b *Builder) HasColumn(table, column string) bool {
	return b.db.Migrator().HasColumn(table, column)
}

// GetColumnListing returns all column names for a table.
func (b *Builder) GetColumnListing(table string) ([]string, error) {
	var columns []string

	// Use GORM's migrator to get column types
	columnTypes, err := b.db.Migrator().ColumnTypes(table)
	if err != nil {
		return nil, err
	}

	for _, col := range columnTypes {
		columns = append(columns, col.Name())
	}

	return columns, nil
}

// DropAllTables drops all tables in the database.
// Use with caution - this is destructive!
func (b *Builder) DropAllTables() error {
	// Get all tables
	var tables []string

	switch b.db.Dialector.Name() {
	case "mysql":
		if err := b.db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
			return err
		}
	case "postgres":
		if err := b.db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
			return err
		}
	case "sqlite", "sqlite3":
		if err := b.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
			return err
		}
	}

	// Disable foreign key checks for MySQL
	if b.db.Dialector.Name() == "mysql" {
		b.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer b.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}

	// Drop each table
	for _, table := range tables {
		if err := b.Drop(table); err != nil {
			return err
		}
	}

	return nil
}

// GetConnection returns the underlying database connection.
func (b *Builder) GetConnection() *gorm.DB {
	return b.db
}

// GetGrammar returns the grammar being used.
func (b *Builder) GetGrammar() Grammar {
	return b.grammar
}

// SetGrammar sets a new grammar for the builder.
func (b *Builder) SetGrammar(grammar Grammar) {
	b.grammar = grammar
}
