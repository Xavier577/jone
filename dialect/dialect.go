// Package dialect provides database-specific SQL generation.
package dialect

import "github.com/Grandbusta/jone/types"

// Dialect defines the interface for database-specific SQL generation.
type Dialect interface {
	// Name returns the dialect name (e.g., "postgresql", "mysql").
	Name() string

	// CreateTableSQL generates a CREATE TABLE statement.
	CreateTableSQL(table *types.Table) string

	// CreateTableIfNotExistsSQL generates a CREATE TABLE IF NOT EXISTS statement.
	CreateTableIfNotExistsSQL(table *types.Table) string

	// DropTableSQL generates a DROP TABLE statement.
	DropTableSQL(schema, name string) string

	// DropTableIfExistsSQL generates a DROP TABLE IF EXISTS statement.
	DropTableIfExistsSQL(schema, name string) string

	// AlterTableSQL generates ALTER TABLE statements for all actions.
	AlterTableSQL(schema, tableName string, actions []*types.TableAction) []string

	// HasTableSQL returns SQL to check if a table exists (returns count).
	HasTableSQL(schema, tableName string) string

	// HasColumnSQL returns SQL to check if a column exists (returns count).
	HasColumnSQL(schema, tableName, columnName string) string

	// ColumnDefinitionSQL generates the column definition for use in CREATE TABLE.
	ColumnDefinitionSQL(col *types.Column) string

	// CommentColumnSQL returns SQL to add a comment to a column.
	CommentColumnSQL(tableName, columnName, comment string) string

	// QuoteIdentifier quotes an identifier (table/column name) for this dialect.
	QuoteIdentifier(name string) string

	// QualifyTable returns a schema-qualified table name.
	// If schema is empty, returns just the quoted table name.
	QualifyTable(schema, tableName string) string
}

// GetDialect returns a dialect implementation by name.
func GetDialect(name string) Dialect {
	switch name {
	case "postgresql", "postgres", "pg":
		return &PostgresDialect{}
	// case "mysql":
	// 	return &MySQLDialect{}
	default:
		// Default to PostgreSQL
		return &PostgresDialect{}
	}
}
