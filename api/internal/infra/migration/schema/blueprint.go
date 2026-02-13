package schema

// Command represents a schema command (like creating indexes)
type Command struct {
	name   string
	params map[string]interface{}
}

// Blueprint defines the structure of a database table.
// It provides a fluent DSL for defining columns, indexes, and constraints.
type Blueprint struct {
	table    string
	columns  []*ColumnDefinition
	commands []*Command
	creating bool
}

// NewBlueprint creates a new blueprint for a table.
func NewBlueprint(table string) *Blueprint {
	return &Blueprint{
		table:    table,
		columns:  make([]*ColumnDefinition, 0),
		commands: make([]*Command, 0),
	}
}

// Create marks this blueprint as creating a new table.
func (b *Blueprint) Create() {
	b.creating = true
	b.addCommand("create", nil)
}

// GetTable returns the table name.
func (b *Blueprint) GetTable() string {
	return b.table
}

// GetColumns returns all column definitions.
func (b *Blueprint) GetColumns() []*ColumnDefinition {
	return b.columns
}

// GetCommands returns all commands.
func (b *Blueprint) GetCommands() []*Command {
	return b.commands
}

// IsCreating returns whether this blueprint is creating a new table.
func (b *Blueprint) IsCreating() bool {
	return b.creating
}

// ID creates an auto-incrementing big integer primary key.
// Default column name is "id".
func (b *Blueprint) ID(column ...string) *ColumnDefinition {
	name := "id"
	if len(column) > 0 {
		name = column[0]
	}
	return b.BigIncrements(name)
}

// BigIncrements creates an auto-incrementing big integer column.
func (b *Blueprint) BigIncrements(column string) *ColumnDefinition {
	return b.addColumn("bigInteger", column).Unsigned().AutoIncrement().Primary()
}

// Increments creates an auto-incrementing integer column.
func (b *Blueprint) Increments(column string) *ColumnDefinition {
	return b.addColumn("integer", column).Unsigned().AutoIncrement().Primary()
}

// String creates a VARCHAR column with optional length (default 255).
func (b *Blueprint) String(column string, length ...int) *ColumnDefinition {
	l := 255
	if len(length) > 0 {
		l = length[0]
	}
	col := b.addColumn("string", column)
	col.length = l
	return col
}

// Text creates a TEXT column.
func (b *Blueprint) Text(column string) *ColumnDefinition {
	return b.addColumn("text", column)
}

// MediumText creates a MEDIUMTEXT column.
func (b *Blueprint) MediumText(column string) *ColumnDefinition {
	return b.addColumn("mediumText", column)
}

// LongText creates a LONGTEXT column.
func (b *Blueprint) LongText(column string) *ColumnDefinition {
	return b.addColumn("longText", column)
}

// Integer creates an INTEGER column.
func (b *Blueprint) Integer(column string) *ColumnDefinition {
	return b.addColumn("integer", column)
}

// TinyInteger creates a TINYINT column.
func (b *Blueprint) TinyInteger(column string) *ColumnDefinition {
	return b.addColumn("tinyInteger", column)
}

// SmallInteger creates a SMALLINT column.
func (b *Blueprint) SmallInteger(column string) *ColumnDefinition {
	return b.addColumn("smallInteger", column)
}

// BigInteger creates a BIGINT column.
func (b *Blueprint) BigInteger(column string) *ColumnDefinition {
	return b.addColumn("bigInteger", column)
}

// UnsignedBigInteger creates an unsigned BIGINT column.
func (b *Blueprint) UnsignedBigInteger(column string) *ColumnDefinition {
	return b.BigInteger(column).Unsigned()
}

// UnsignedInteger creates an unsigned INTEGER column.
func (b *Blueprint) UnsignedInteger(column string) *ColumnDefinition {
	return b.Integer(column).Unsigned()
}

// Boolean creates a BOOLEAN column.
func (b *Blueprint) Boolean(column string) *ColumnDefinition {
	return b.addColumn("boolean", column)
}

// Timestamp creates a TIMESTAMP column.
func (b *Blueprint) Timestamp(column string) *ColumnDefinition {
	return b.addColumn("timestamp", column)
}

// DateTime creates a DATETIME column.
func (b *Blueprint) DateTime(column string) *ColumnDefinition {
	return b.addColumn("datetime", column)
}

// Date creates a DATE column.
func (b *Blueprint) Date(column string) *ColumnDefinition {
	return b.addColumn("date", column)
}

// Time creates a TIME column.
func (b *Blueprint) Time(column string) *ColumnDefinition {
	return b.addColumn("time", column)
}

// Timestamps creates created_at and updated_at timestamp columns.
func (b *Blueprint) Timestamps() {
	b.Timestamp("created_at").Nullable()
	b.Timestamp("updated_at").Nullable()
}

// SoftDeletes creates a deleted_at column for soft deletes.
// Default column name is "deleted_at".
func (b *Blueprint) SoftDeletes(column ...string) *ColumnDefinition {
	name := "deleted_at"
	if len(column) > 0 {
		name = column[0]
	}
	return b.Timestamp(name).Nullable()
}

// JSON creates a JSON column.
func (b *Blueprint) JSON(column string) *ColumnDefinition {
	return b.addColumn("json", column)
}

// Binary creates a BINARY/BLOB column.
func (b *Blueprint) Binary(column string) *ColumnDefinition {
	return b.addColumn("binary", column)
}

// Float creates a FLOAT column.
func (b *Blueprint) Float(column string) *ColumnDefinition {
	return b.addColumn("float", column)
}

// Double creates a DOUBLE column.
func (b *Blueprint) Double(column string) *ColumnDefinition {
	return b.addColumn("double", column)
}

// Decimal creates a DECIMAL column with precision and scale.
func (b *Blueprint) Decimal(column string, precision, scale int) *ColumnDefinition {
	col := b.addColumn("decimal", column)
	col.precision = precision
	col.scale = scale
	return col
}

// UUID creates a UUID column.
func (b *Blueprint) UUID(column string) *ColumnDefinition {
	return b.addColumn("uuid", column)
}

// addColumn adds a column to the blueprint.
func (b *Blueprint) addColumn(colType, name string) *ColumnDefinition {
	col := &ColumnDefinition{
		name:    name,
		colType: colType,
	}
	b.columns = append(b.columns, col)
	return col
}

// addCommand adds a command to the blueprint.
func (b *Blueprint) addCommand(name string, params map[string]interface{}) {
	b.commands = append(b.commands, &Command{name: name, params: params})
}

// ToSQL generates SQL statements for the blueprint using the given grammar.
func (b *Blueprint) ToSQL(grammar Grammar) []string {
	return grammar.Compile(b)
}

// Primary creates a primary key on the specified columns.
func (b *Blueprint) Primary(columns ...string) {
	b.addCommand("primary", map[string]interface{}{"columns": columns})
}

// Unique creates a unique index on the specified columns.
func (b *Blueprint) Unique(columns ...string) {
	b.addCommand("unique", map[string]interface{}{"columns": columns})
}

// Index creates an index on the specified columns.
func (b *Blueprint) Index(columns ...string) {
	b.addCommand("index", map[string]interface{}{"columns": columns})
}

// Foreign creates a foreign key constraint on the specified column.
// Returns a ForeignKeyDefinition for fluent configuration.
func (b *Blueprint) Foreign(column string) *ForeignKeyDefinition {
	fk := &ForeignKeyDefinition{column: column}
	b.addCommand("foreign", map[string]interface{}{"definition": fk})
	return fk
}

// DropColumn drops the specified columns from the table.
func (b *Blueprint) DropColumn(columns ...string) {
	b.addCommand("dropColumn", map[string]interface{}{"columns": columns})
}

// RenameColumn renames a column from one name to another.
func (b *Blueprint) RenameColumn(from, to string) {
	b.addCommand("renameColumn", map[string]interface{}{"from": from, "to": to})
}

// DropPrimary drops the primary key from the table.
func (b *Blueprint) DropPrimary() {
	b.addCommand("dropPrimary", nil)
}

// DropUnique drops a unique index by name.
func (b *Blueprint) DropUnique(index string) {
	b.addCommand("dropUnique", map[string]interface{}{"index": index})
}

// DropIndex drops an index by name.
func (b *Blueprint) DropIndex(index string) {
	b.addCommand("dropIndex", map[string]interface{}{"index": index})
}

// DropForeign drops a foreign key constraint by name.
func (b *Blueprint) DropForeign(index string) {
	b.addCommand("dropForeign", map[string]interface{}{"index": index})
}

// ForeignID creates an unsigned big integer column and a foreign key constraint.
// This is a convenience method for the common pattern of foreign key columns.
func (b *Blueprint) ForeignID(column string) *ForeignKeyDefinition {
	b.UnsignedBigInteger(column)
	return b.Foreign(column)
}

// DropTimestamps drops the created_at and updated_at columns.
func (b *Blueprint) DropTimestamps() {
	b.DropColumn("created_at", "updated_at")
}

// DropSoftDeletes drops the deleted_at column.
func (b *Blueprint) DropSoftDeletes(column ...string) {
	name := "deleted_at"
	if len(column) > 0 {
		name = column[0]
	}
	b.DropColumn(name)
}
