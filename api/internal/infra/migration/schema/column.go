package schema

// ColumnDefinition represents a column in a blueprint.
// It supports method chaining for fluent column configuration.
type ColumnDefinition struct {
	name          string
	colType       string
	length        int
	precision     int
	scale         int
	nullable      bool
	defaultValue  interface{}
	unsigned      bool
	autoIncrement bool
	primary       bool
	comment       string
	after         string
	first         bool
	change        bool
}

// Nullable marks the column as nullable.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Nullable() *ColumnDefinition {
	c.nullable = true
	return c
}

// Default sets the default value for the column.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Default(value interface{}) *ColumnDefinition {
	c.defaultValue = value
	return c
}

// Unsigned marks the column as unsigned (for numeric types).
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Unsigned() *ColumnDefinition {
	c.unsigned = true
	return c
}

// AutoIncrement marks the column as auto-incrementing.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) AutoIncrement() *ColumnDefinition {
	c.autoIncrement = true
	return c
}

// Primary marks the column as a primary key.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Primary() *ColumnDefinition {
	c.primary = true
	return c
}

// Comment adds a comment to the column.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Comment(comment string) *ColumnDefinition {
	c.comment = comment
	return c
}

// After places the column after another column (MySQL only).
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) After(column string) *ColumnDefinition {
	c.after = column
	return c
}

// First places the column first in the table (MySQL only).
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) First() *ColumnDefinition {
	c.first = true
	return c
}

// Change marks the column for modification instead of creation.
// Returns the ColumnDefinition for method chaining.
func (c *ColumnDefinition) Change() *ColumnDefinition {
	c.change = true
	return c
}

// GetName returns the column name.
func (c *ColumnDefinition) GetName() string {
	return c.name
}

// GetType returns the column type.
func (c *ColumnDefinition) GetType() string {
	return c.colType
}

// GetLength returns the column length.
func (c *ColumnDefinition) GetLength() int {
	return c.length
}

// IsNullable returns whether the column is nullable.
func (c *ColumnDefinition) IsNullable() bool {
	return c.nullable
}

// GetDefault returns the default value.
func (c *ColumnDefinition) GetDefault() interface{} {
	return c.defaultValue
}

// IsUnsigned returns whether the column is unsigned.
func (c *ColumnDefinition) IsUnsigned() bool {
	return c.unsigned
}

// IsAutoIncrement returns whether the column is auto-incrementing.
func (c *ColumnDefinition) IsAutoIncrement() bool {
	return c.autoIncrement
}

// IsPrimary returns whether the column is a primary key.
func (c *ColumnDefinition) IsPrimary() bool {
	return c.primary
}

// GetComment returns the column comment.
func (c *ColumnDefinition) GetComment() string {
	return c.comment
}
