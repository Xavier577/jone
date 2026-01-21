package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the latest migration",
	Run:   migrateRollback,
}

func migrateRollback(cmd *cobra.Command, args []string) {
	fmt.Println("Rolling back the latest migration...")
}
