package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateLatestCmd = &cobra.Command{
	Use:   "migrate:latest",
	Short: "Runs all pending migrations",
	Long:  `Runs all pending migrations by generating and executing a runner`,
	Run:   migrateLatestJone,
}

func init() {
	migrateLatestCmd.Flags().Bool("dry-run", false, "Show SQL without executing")
}

func migrateLatestJone(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	execParams := RunExecParams{
		Command: "migrate:latest",
		Flags: map[string]any{
			"dry-run": dryRun,
		},
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}
}
