package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a jone project",
	Long:  `Initializes a jone project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called frol cli tool")
	},
}
