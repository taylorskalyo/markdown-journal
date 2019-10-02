package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var application = &cobra.Command{
	Use:   "markdown-journal",
	Short: "markdown-journal helps you manage a markdown journal",
	Long:  `A markdown journaling system`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute root command.
func Execute() {
	if err := application.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
