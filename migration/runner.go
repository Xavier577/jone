package migration

import (
	"fmt"

	"github.com/Grandbusta/jone/config"
	"github.com/Grandbusta/jone/schema"
)

// RunLatest executes all Up migrations in order using the provided schema.
func RunLatest(cfg *config.Config, registrations []Registration, s *schema.Schema) error {
	fmt.Println("Running migrations (migrate:latest)...")

	for _, reg := range registrations {
		fmt.Printf("Running migration: %s (up)\n", reg.Name)
		reg.Up(s)
		fmt.Printf("Completed migration: %s\n", reg.Name)
	}
	fmt.Println("All migrations completed successfully")
	return nil
}

// RunDown executes all Down migrations in reverse order using the provided schema.
func RunDown(cfg *config.Config, registrations []Registration, s *schema.Schema) error {
	fmt.Println("Running rollback (migrate:down)...")

	// Run in reverse order
	for i := len(registrations) - 1; i >= 0; i-- {
		reg := registrations[i]
		fmt.Printf("Rolling back migration: %s (down)\n", reg.Name)
		reg.Down(s)
		fmt.Printf("Completed rollback: %s\n", reg.Name)
	}
	fmt.Println("All rollbacks completed successfully")
	return nil
}

func RunRollback(cfg *config.Config, registrations []Registration, s *schema.Schema) error {
	fmt.Println("Running rollback (migrate:rollback)...")

	return nil
}
