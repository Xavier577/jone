package schema

import (
	"database/sql"
	"fmt"

	"github.com/Grandbusta/jone/config"
	"github.com/Grandbusta/jone/dialect"
)

// Schema provides methods for database schema operations.
type Schema struct {
	dialect dialect.Dialect
	db      *sql.DB
	config  *config.Config
	schema  string // current schema context
}

// New creates a new Schema with the given config.
// It determines the dialect from the config and can optionally connect to the database.
func New(cfg *config.Config) *Schema {
	d := dialect.GetDialect(cfg.Client)
	return &Schema{
		dialect: d,
		config:  cfg,
	}
}

// WithSchema returns a new Schema that operates on the specified schema.
func (s *Schema) WithSchema(schemaName string) *Schema {
	return &Schema{
		dialect: s.dialect,
		db:      s.db,
		config:  s.config,
		schema:  schemaName,
	}
}

// SchemaName returns the current schema name (empty = default).
func (s *Schema) SchemaName() string {
	return s.schema
}

// Dialect returns the current dialect.
func (s *Schema) Dialect() dialect.Dialect {
	return s.dialect
}

// DB returns the database connection, if set.
func (s *Schema) DB() *sql.DB {
	return s.db
}

// SetDB sets the database connection.
func (s *Schema) SetDB(db *sql.DB) {
	s.db = db
}

// Open opens a database connection using the config.
func (s *Schema) Open() error {
	dsn := s.config.Connection.DSN()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}
	s.db = db
	return nil
}

// Close closes the database connection.
func (s *Schema) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Schema) Table(name string, builder func(t *Table)) error {
	t := NewTable(name)
	t.Schema = s.schema // Set schema context
	builder(t)

	// Generate SQL for each action
	statements := s.dialect.AlterTableSQL(s.schema, name, t.Actions)

	for _, sql := range statements {
		fmt.Printf("SQL: %s\n", sql)
		if s.db != nil {
			_, err := s.db.Exec(sql)
			if err != nil {
				return fmt.Errorf("executing ALTER TABLE: %w", err)
			}
		} else {
			fmt.Printf("[DRY RUN] Would execute: %s\n", sql)
		}
	}

	return nil
}

// CreateTable creates a new table with the given name using the builder function.
func (s *Schema) CreateTable(name string, builder func(t *Table)) error {
	t := NewTable(name)
	t.Schema = s.schema // Set schema context
	builder(t)

	sql := s.dialect.CreateTableSQL(t.Table)
	fmt.Printf("SQL: %s\n", sql)

	// Execute if we have a database connection
	if s.db != nil {
		_, err := s.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("executing CREATE TABLE: %w", err)
		}
		fmt.Printf("Created table: %s\n", name)

		// Execute COMMENT ON COLUMN for columns with comments (PostgreSQL needs separate statement)
		qualifiedTable := s.dialect.QualifyTable(s.schema, name)
		for _, col := range t.Columns {
			if col.Comment != "" {
				commentSQL := s.dialect.CommentColumnSQL(qualifiedTable, col.Name, col.Comment)
				fmt.Printf("SQL: %s\n", commentSQL)
				if _, err := s.db.Exec(commentSQL); err != nil {
					return fmt.Errorf("executing COMMENT ON COLUMN: %w", err)
				}
			}
		}
	}
	return nil
}

// CreateTableIfNotExists creates a new table if it doesn't already exist.
func (s *Schema) CreateTableIfNotExists(name string, builder func(t *Table)) error {
	t := NewTable(name)
	t.Schema = s.schema // Set schema context
	builder(t)

	sql := s.dialect.CreateTableIfNotExistsSQL(t.Table)
	fmt.Printf("SQL: %s\n", sql)

	if s.db != nil {
		_, err := s.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("executing CREATE TABLE IF NOT EXISTS: %w", err)
		}
		fmt.Printf("Created table (if not exists): %s\n", name)
	} else {
		fmt.Printf("[DRY RUN] Would create table if not exists: %s with %d columns\n", name, len(t.Columns))
	}

	return nil
}

// DropTable drops a table by name.
func (s *Schema) DropTable(name string) error {
	sql := s.dialect.DropTableSQL(s.schema, name)
	fmt.Printf("SQL: %s\n", sql)

	if s.db != nil {
		_, err := s.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("executing DROP TABLE: %w", err)
		}
		fmt.Printf("Dropped table: %s\n", name)
	} else {
		fmt.Printf("[DRY RUN] Would drop table: %s\n", name)
	}

	return nil
}

// DropTableIfExists drops a table if it exists.
func (s *Schema) DropTableIfExists(name string) error {
	sql := s.dialect.DropTableIfExistsSQL(s.schema, name)
	fmt.Printf("SQL: %s\n", sql)

	if s.db != nil {
		_, err := s.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("executing DROP TABLE IF EXISTS: %w", err)
		}
		fmt.Printf("Dropped table (if existed): %s\n", name)
	} else {
		fmt.Printf("[DRY RUN] Would drop table if exists: %s\n", name)
	}

	return nil
}

// RenameTable renames a table from oldName to newName.
func (s *Schema) RenameTable(oldName, newName string) error {
	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;",
		s.dialect.QualifyTable(s.schema, oldName),
		s.dialect.QuoteIdentifier(newName))
	fmt.Printf("SQL: %s\n", sql)

	if s.db != nil {
		_, err := s.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("executing RENAME TABLE: %w", err)
		}
		fmt.Printf("Renamed table: %s -> %s\n", oldName, newName)
	} else {
		fmt.Printf("[DRY RUN] Would rename table: %s -> %s\n", oldName, newName)
	}

	return nil
}

// HasTable checks if a table exists.
func (s *Schema) HasTable(name string) bool {
	if s.db == nil {
		return false
	}
	sql := s.dialect.HasTableSQL(s.schema, name)
	var count int
	if err := s.db.QueryRow(sql).Scan(&count); err != nil {
		return false
	}
	return count > 0
}

// HasColumn checks if a column exists in a table.
func (s *Schema) HasColumn(table, column string) bool {
	if s.db == nil {
		return false
	}
	sql := s.dialect.HasColumnSQL(s.schema, table, column)
	var count int
	if err := s.db.QueryRow(sql).Scan(&count); err != nil {
		return false
	}
	return count > 0
}
