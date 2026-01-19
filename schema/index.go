package schema

import (
	"strings"

	"github.com/Grandbusta/jone/types"
)

// IndexBuilder provides a fluent interface for creating indexes.
type IndexBuilder struct {
	table   *Table
	columns []string
	name    string
	method  string
	unique  bool
}

// Name sets a custom name for the index.
func (b *IndexBuilder) Name(n string) *IndexBuilder {
	b.name = n
	// Update the action that was already added
	b.updateAction()
	return b
}

// Using sets the index method (e.g., btree, hash, gin, gist).
func (b *IndexBuilder) Using(method string) *IndexBuilder {
	b.method = method
	b.updateAction()
	return b
}

// build creates the Index struct with auto-generated name if needed.
func (b *IndexBuilder) build() *types.Index {
	name := b.name
	if name == "" {
		name = b.generateName()
	}
	return &types.Index{
		Name:      name,
		Columns:   b.columns,
		IsUnique:  b.unique,
		Method:    b.method,
		TableName: b.table.Name,
	}
}

// generateName creates an auto-generated index name.
func (b *IndexBuilder) generateName() string {
	prefix := "idx"
	if b.unique {
		prefix = "uq"
	}
	return prefix + "_" + b.table.Name + "_" + strings.Join(b.columns, "_")
}

// updateAction updates the last action with the current builder state.
func (b *IndexBuilder) updateAction() {
	if len(b.table.Actions) > 0 {
		lastAction := b.table.Actions[len(b.table.Actions)-1]
		if lastAction.Type == types.ActionCreateIndex {
			lastAction.Index = b.build()
		}
	}
}
