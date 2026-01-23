package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the latest migration",
	Run:   migrateRollback,
}

func init() {
	migrateRollbackCmd.Flags().BoolP("all", "a", false, "Rollback all migrations")
}

func migrateRollback(cmd *cobra.Command, args []string) {
	allFlag, _ := cmd.Flags().GetBool("all")
	execParams := RunExecParams{
		Command: "migrate:rollback",
		Flags: map[string]any{
			"all": allFlag,
		},
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}
}
