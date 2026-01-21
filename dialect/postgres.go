package dialect

import (
	"fmt"
	"strings"

	"github.com/Grandbusta/jone/types"
)

// PostgresDialect implements Dialect for PostgreSQL.
type PostgresDialect struct{}

// Name returns "postgresql".
func (d *PostgresDialect) Name() string {
	return "postgresql"
}

// QuoteIdentifier quotes an identifier with double quotes for PostgreSQL.
func (d *PostgresDialect) QuoteIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

// CreateTableSQL generates a CREATE TABLE statement for PostgreSQL.
func (d *PostgresDialect) CreateTableSQL(table *types.Table) string {
	var columns []string
	for _, col := range table.Columns {
		columns = append(columns, d.ColumnDefinitionSQL(col))
	}

	return fmt.Sprintf(
		"CREATE TABLE %s (\n  %s\n);",
		d.QualifyTable(table.Schema, table.Name),
		strings.Join(columns, ",\n  "),
	)
}

// DropTableSQL generates a DROP TABLE statement.
func (d *PostgresDialect) DropTableSQL(schema, name string) string {
	return fmt.Sprintf("DROP TABLE %s;", d.QualifyTable(schema, name))
}

// DropTableIfExistsSQL generates a DROP TABLE IF EXISTS statement.
func (d *PostgresDialect) DropTableIfExistsSQL(schema, name string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", d.QualifyTable(schema, name))
}

// CreateTableIfNotExistsSQL generates a CREATE TABLE IF NOT EXISTS statement.
func (d *PostgresDialect) CreateTableIfNotExistsSQL(table *types.Table) string {
	var columns []string
	for _, col := range table.Columns {
		columns = append(columns, d.ColumnDefinitionSQL(col))
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n  %s\n);",
		d.QualifyTable(table.Schema, table.Name),
		strings.Join(columns, ",\n  "),
	)
}

// ColumnDefinitionSQL generates the column definition SQL.
func (d *PostgresDialect) ColumnDefinitionSQL(col *types.Column) string {
	var parts []string

	parts = append(parts, d.QuoteIdentifier(col.Name))
	parts = append(parts, d.mapDataType(col))

	if col.IsPrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}
	if col.IsNotNull && !col.IsPrimaryKey {
		parts = append(parts, "NOT NULL")
	}
	if col.IsUnique && !col.IsPrimaryKey {
		parts = append(parts, "UNIQUE")
	}
	if col.HasDefault {
		parts = append(parts, fmt.Sprintf("DEFAULT %v", d.formatDefault(col.DefaultValue)))
	}
	if col.RefTable != "" && col.RefColumn != "" {
		parts = append(parts, fmt.Sprintf(
			"REFERENCES %s(%s)",
			d.QuoteIdentifier(col.RefTable),
			d.QuoteIdentifier(col.RefColumn),
		))
	}

	return strings.Join(parts, " ")
}

// mapDataType maps generic types to PostgreSQL-specific types.
func (d *PostgresDialect) mapDataType(col *types.Column) string {
	switch col.DataType {
	case "varchar":
		if col.Length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.Length)
		}
		return "VARCHAR(255)"
	case "char":
		if col.Length > 0 {
			return fmt.Sprintf("CHAR(%d)", col.Length)
		}
		return "CHAR(1)"
	case "int":
		return "INTEGER"
	case "bigint":
		return "BIGINT"
	case "smallint":
		return "SMALLINT"
	case "float":
		if col.Precision > 0 {
			return fmt.Sprintf("FLOAT(%d)", col.Precision)
		}
		return "REAL"
	case "double":
		return "DOUBLE PRECISION"
	case "decimal":
		p := col.Precision
		if p == 0 {
			p = 10
		}
		s := col.Scale
		if s == 0 {
			s = 2
		}
		return fmt.Sprintf("DECIMAL(%d,%d)", p, s)
	case "boolean":
		return "BOOLEAN"
	case "text":
		return "TEXT"
	case "date":
		return "DATE"
	case "time":
		return "TIME"
	case "timestamp":
		return "TIMESTAMP"
	case "uuid":
		return "UUID"
	case "json":
		return "JSON"
	case "jsonb":
		return "JSONB"
	case "binary":
		return "BYTEA"
	case "serial":
		return "SERIAL"
	case "bigserial":
		return "BIGSERIAL"
	default:
		return strings.ToUpper(col.DataType)
	}
}

// formatDefault formats a default value for SQL.
func (d *PostgresDialect) formatDefault(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// AlterTableSQL generates ALTER TABLE statements for all actions.
func (d *PostgresDialect) AlterTableSQL(schema, tableName string, actions []*types.TableAction) []string {
	qualifiedTable := d.QualifyTable(schema, tableName)
	var statements []string
	for _, action := range actions {
		switch action.Type {
		case types.ActionDropColumn:
			statements = append(statements, d.dropColumnSQL(qualifiedTable, action.Name))
		case types.ActionAddColumn:
			statements = append(statements, d.addColumnSQL(qualifiedTable, action.Column))
		case types.ActionRenameColumn:
			statements = append(statements, d.renameColumnSQL(qualifiedTable, action.Name, action.NewName))
		case types.ActionChangeColumnType:
			statements = append(statements, d.changeColumnTypeSQL(qualifiedTable, action.Column))
		case types.ActionSetColumnNotNull:
			statements = append(statements, d.setColumnNotNullSQL(qualifiedTable, action.Name))
		case types.ActionDropColumnNotNull:
			statements = append(statements, d.dropColumnNotNullSQL(qualifiedTable, action.Name))
		case types.ActionSetColumnDefault:
			statements = append(statements, d.setColumnDefaultSQL(qualifiedTable, action.Name, action.DefaultValue))
		case types.ActionDropColumnDefault:
			statements = append(statements, d.dropColumnDefaultSQL(qualifiedTable, action.Name))
		case types.ActionCreateIndex:
			statements = append(statements, d.createIndexSQL(qualifiedTable, action.Index))
		case types.ActionDropIndex:
			statements = append(statements, d.dropIndexSQL(schema, action.Index.Name))
		case types.ActionAddForeignKey:
			statements = append(statements, d.addForeignKeySQL(qualifiedTable, action.ForeignKey))
		case types.ActionDropForeignKey:
			statements = append(statements, d.dropForeignKeySQL(qualifiedTable, action.ForeignKey.Name))
		case types.ActionDropPrimary:
			constraintName := action.Name
			if constraintName == "" {
				constraintName = tableName + "_pkey"
			}
			statements = append(statements, d.dropPrimarySQL(qualifiedTable, constraintName))
		}
	}
	return statements
}

// createIndexSQL generates a CREATE INDEX statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) createIndexSQL(tableName string, idx *types.Index) string {
	unique := ""
	if idx.IsUnique {
		unique = "UNIQUE "
	}

	using := ""
	if idx.Method != "" {
		using = fmt.Sprintf(" USING %s", idx.Method)
	}

	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = d.QuoteIdentifier(c)
	}

	return fmt.Sprintf("CREATE %sINDEX %s ON %s%s (%s);",
		unique,
		d.QuoteIdentifier(idx.Name),
		tableName,
		using,
		strings.Join(cols, ", "))
}

// dropIndexSQL generates a DROP INDEX statement.
// In PostgreSQL, indexes are schema-scoped and need to be qualified.
func (d *PostgresDialect) dropIndexSQL(schema, name string) string {
	if schema == "" {
		return fmt.Sprintf("DROP INDEX %s;", d.QuoteIdentifier(name))
	}
	return fmt.Sprintf("DROP INDEX %s.%s;", d.QuoteIdentifier(schema), d.QuoteIdentifier(name))
}

// addForeignKeySQL generates an ALTER TABLE ADD CONSTRAINT FOREIGN KEY statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) addForeignKeySQL(tableName string, fk *types.ForeignKey) string {
	sql := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)",
		tableName,
		d.QuoteIdentifier(fk.Name),
		d.QuoteIdentifier(fk.Column),
		d.QuoteIdentifier(fk.RefTable),
		d.QuoteIdentifier(fk.RefColumn))

	if fk.OnDelete != "" {
		sql += " ON DELETE " + fk.OnDelete
	}
	if fk.OnUpdate != "" {
		sql += " ON UPDATE " + fk.OnUpdate
	}
	return sql + ";"
}

// dropForeignKeySQL generates an ALTER TABLE DROP CONSTRAINT statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) dropForeignKeySQL(tableName, fkName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;",
		tableName,
		d.QuoteIdentifier(fkName))
}

// dropColumnSQL generates an ALTER TABLE DROP COLUMN statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) dropColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;",
		tableName,
		d.QuoteIdentifier(columnName))
}

// addColumnSQL generates an ALTER TABLE ADD COLUMN statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) addColumnSQL(tableName string, column *types.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;",
		tableName,
		d.ColumnDefinitionSQL(column))
}

// renameColumnSQL generates an ALTER TABLE RENAME COLUMN statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) renameColumnSQL(tableName, oldName, newName string) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;",
		tableName,
		d.QuoteIdentifier(oldName),
		d.QuoteIdentifier(newName))
}

// changeColumnTypeSQL generates an ALTER TABLE ALTER COLUMN TYPE statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) changeColumnTypeSQL(tableName string, column *types.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;",
		tableName,
		d.QuoteIdentifier(column.Name),
		d.mapDataType(column))
}

// setColumnNotNullSQL generates an ALTER TABLE ALTER COLUMN SET NOT NULL statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) setColumnNotNullSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;",
		tableName,
		d.QuoteIdentifier(columnName))
}

// dropColumnNotNullSQL generates an ALTER TABLE ALTER COLUMN DROP NOT NULL statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) dropColumnNotNullSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL;",
		tableName,
		d.QuoteIdentifier(columnName))
}

// setColumnDefaultSQL generates an ALTER TABLE ALTER COLUMN SET DEFAULT statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) setColumnDefaultSQL(tableName, columnName string, defaultValue any) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;",
		tableName,
		d.QuoteIdentifier(columnName),
		d.formatDefault(defaultValue))
}

// dropColumnDefaultSQL generates an ALTER TABLE ALTER COLUMN DROP DEFAULT statement.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) dropColumnDefaultSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT;",
		tableName,
		d.QuoteIdentifier(columnName))
}

// HasTableSQL returns SQL to check if a table exists in PostgreSQL.
func (d *PostgresDialect) HasTableSQL(schema, tableName string) string {
	schemaName := schema
	if schemaName == "" {
		schemaName = "public"
	}
	return fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s'`, schemaName, tableName)
}

// HasColumnSQL returns SQL to check if a column exists in PostgreSQL.
func (d *PostgresDialect) HasColumnSQL(schema, tableName, columnName string) string {
	schemaName := schema
	if schemaName == "" {
		schemaName = "public"
	}
	return fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = '%s' AND table_name = '%s' AND column_name = '%s'`, schemaName, tableName, columnName)
}

// CommentColumnSQL returns SQL to add a comment to a column in PostgreSQL.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) CommentColumnSQL(tableName, columnName, comment string) string {
	return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';",
		tableName,
		d.QuoteIdentifier(columnName),
		comment)
}

// dropPrimarySQL returns SQL to drop the primary key constraint in PostgreSQL.
// tableName should be pre-qualified (e.g., from QualifyTable).
func (d *PostgresDialect) dropPrimarySQL(tableName, constraintName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;",
		tableName,
		d.QuoteIdentifier(constraintName))
}

// QualifyTable returns a schema-qualified table name.
func (d *PostgresDialect) QualifyTable(schema, tableName string) string {
	if schema == "" {
		return d.QuoteIdentifier(tableName)
	}
	return fmt.Sprintf("%s.%s", d.QuoteIdentifier(schema), d.QuoteIdentifier(tableName))
}
