package schema

import "strings"

// ForeignKeyDefinition represents a foreign key constraint.
// It supports method chaining for fluent foreign key configuration.
type ForeignKeyDefinition struct {
	column     string
	references string
	on         string
	onDelete   string
	onUpdate   string
}

// References sets the referenced column.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) References(column string) *ForeignKeyDefinition {
	f.references = column
	return f
}

// On sets the referenced table.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) On(table string) *ForeignKeyDefinition {
	f.on = table
	return f
}

// OnDelete sets the ON DELETE action.
// Common values: "CASCADE", "SET NULL", "RESTRICT", "NO ACTION"
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) OnDelete(action string) *ForeignKeyDefinition {
	f.onDelete = strings.ToUpper(action)
	return f
}

// OnUpdate sets the ON UPDATE action.
// Common values: "CASCADE", "SET NULL", "RESTRICT", "NO ACTION"
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) OnUpdate(action string) *ForeignKeyDefinition {
	f.onUpdate = strings.ToUpper(action)
	return f
}

// Cascade sets both ON DELETE and ON UPDATE to CASCADE.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) Cascade() *ForeignKeyDefinition {
	f.onDelete = "CASCADE"
	f.onUpdate = "CASCADE"
	return f
}

// CascadeOnDelete sets ON DELETE to CASCADE.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) CascadeOnDelete() *ForeignKeyDefinition {
	f.onDelete = "CASCADE"
	return f
}

// CascadeOnUpdate sets ON UPDATE to CASCADE.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) CascadeOnUpdate() *ForeignKeyDefinition {
	f.onUpdate = "CASCADE"
	return f
}

// NullOnDelete sets ON DELETE to SET NULL.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) NullOnDelete() *ForeignKeyDefinition {
	f.onDelete = "SET NULL"
	return f
}

// RestrictOnDelete sets ON DELETE to RESTRICT.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) RestrictOnDelete() *ForeignKeyDefinition {
	f.onDelete = "RESTRICT"
	return f
}

// Constrained is a shorthand for common foreign key pattern.
// If table is provided, it sets the referenced table.
// It also sets the referenced column to "id" by default.
// Returns the ForeignKeyDefinition for method chaining.
func (f *ForeignKeyDefinition) Constrained(table ...string) *ForeignKeyDefinition {
	if len(table) > 0 {
		f.on = table[0]
	}
	if f.references == "" {
		f.references = "id"
	}
	return f
}

// GetColumn returns the foreign key column name.
func (f *ForeignKeyDefinition) GetColumn() string {
	return f.column
}

// GetReferences returns the referenced column name.
func (f *ForeignKeyDefinition) GetReferences() string {
	return f.references
}

// GetOn returns the referenced table name.
func (f *ForeignKeyDefinition) GetOn() string {
	return f.on
}

// GetOnDelete returns the ON DELETE action.
func (f *ForeignKeyDefinition) GetOnDelete() string {
	return f.onDelete
}

// GetOnUpdate returns the ON UPDATE action.
func (f *ForeignKeyDefinition) GetOnUpdate() string {
	return f.onUpdate
}
