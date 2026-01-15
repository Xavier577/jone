// Package migration provides migration registration and execution.
package migration

import "github.com/Grandbusta/jone/schema"

// Registration represents a single migration with its metadata and operations.
type Registration struct {
	Name string
	Up   func(*schema.Schema)
	Down func(*schema.Schema)
}
