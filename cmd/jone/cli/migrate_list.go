package cli

import (
	"fmt"
	"os"

	"github.com/Grandbusta/jone/internal/term"
	"github.com/spf13/cobra"
)

var migrateListCmd = &cobra.Command{
	Use:     "migrate:list",
	Aliases: []string{"migrate:status"},
	Short:   "Lists all migrations with their status (alias: migrate:status)",
	Long:    `Lists all registered migrations showing which are applied and which are pending.`,
	Run:     migrateListJone,
}

func migrateListJone(cmd *cobra.Command, args []string) {
	execParams := RunExecParams{
		Command: "migrate:list",
	}
	if err := runMigrations(execParams); err != nil {
		fmt.Println(term.RedText(fmt.Sprintf("Error listing migrations: %v", err)))
		os.Exit(1)
	}
}
