// Package types provides core types used across the jone library.
// This package has no internal dependencies to prevent import cycles.
package types

// Column represents a database column definition.
type Column struct {
	Name         string
	DataType     string
	IsPrimaryKey bool
	IsNotNull    bool
	IsUnique     bool
	IsUnsigned   bool
	DefaultValue any
	HasDefault   bool
	RefTable     string
	RefColumn    string
}

// Table represents a database table definition.
type Table struct {
	Name    string
	Columns []*Column
}
