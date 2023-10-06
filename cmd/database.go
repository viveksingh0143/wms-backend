package cmd

import (
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "db",
	Short: "database related commands",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
