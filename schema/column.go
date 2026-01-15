package schema

import "github.com/Grandbusta/jone/types"

// Column wraps types.Column and provides modifier methods.
type Column struct {
	*types.Column
}

// Primary marks this column as a primary key.
func (c *Column) Primary() *Column {
	c.IsPrimaryKey = true
	return c
}

// NotNullable marks this column as NOT NULL.
func (c *Column) NotNullable() *Column {
	c.IsNotNull = true
	return c
}

// Nullable marks this column as nullable (removes NOT NULL constraint).
func (c *Column) Nullable() *Column {
	c.IsNotNull = false
	return c
}

// Unique adds a UNIQUE constraint to this column.
func (c *Column) Unique() *Column {
	c.IsUnique = true
	return c
}

// Unsigned marks this column as unsigned (for numeric types).
func (c *Column) Unsigned() *Column {
	c.IsUnsigned = true
	return c
}

// Default sets a default value for this column.
func (c *Column) Default(value any) *Column {
	c.HasDefault = true
	c.DefaultValue = value
	return c
}

// References sets up a foreign key reference to another table's column.
func (c *Column) References(table, column string) *Column {
	c.RefTable = table
	c.RefColumn = column
	return c
}
