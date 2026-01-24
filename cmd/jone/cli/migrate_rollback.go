package cli

import (
	"fmt"
	"os"

	"github.com/Grandbusta/jone/internal/term"
	"github.com/spf13/cobra"
)

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the latest migration",
	Run:   migrateRollback,
}

func init() {
	migrateRollbackCmd.Flags().BoolP("all", "a", false, "Rollback all migrations")
	migrateRollbackCmd.Flags().Bool("dry-run", false, "Show SQL without executing")
}

func migrateRollback(cmd *cobra.Command, args []string) {
	allFlag, _ := cmd.Flags().GetBool("all")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	execParams := RunExecParams{
		Command: "migrate:rollback",
		Flags: map[string]any{
			"all":     allFlag,
			"dry-run": dryRun,
		},
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Println(term.RedText(fmt.Sprintf("Error running migrations: %v", err)))
		os.Exit(1)
	}
}
