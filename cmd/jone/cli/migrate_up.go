package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "migrate:up [migration_name]",
	Short: "Runs the next pending migration or a specific one",
	Long:  `Runs the next pending migration. If a migration name is provided, runs that specific migration.`,
	Run:   migrateUpJone,
}

func init() {
	migrateUpCmd.Flags().Bool("dry-run", false, "Show SQL without executing")
}

func migrateUpJone(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	execParams := RunExecParams{
		Command: "migrate:up",
		Args:    args,
		Flags: map[string]any{
			"dry-run": dryRun,
		},
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Printf("Error running migration: %v\n", err)
		os.Exit(1)
	}
}
