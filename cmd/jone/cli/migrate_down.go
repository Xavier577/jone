package cli

import (
	"fmt"
	"os"

	"github.com/Grandbusta/jone/internal/term"
	"github.com/spf13/cobra"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down",
	Short: "Rolls back migrations",
	Long:  `Rolls back migrations by generating and executing a runner`,
	Run:   migrateDownJone,
}

func init() {
	migrateDownCmd.Flags().Bool("dry-run", false, "Show SQL without executing")
}

func migrateDownJone(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	execParams := RunExecParams{
		Command: "migrate:down",
		Args:    args,
		Flags: map[string]any{
			"dry-run": dryRun,
		},
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Println(term.RedText(fmt.Sprintf("Error rolling back migrations: %v", err)))
		os.Exit(1)
	}
}
