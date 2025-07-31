package ui

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored applications",
	Long: `List all applications for which passwords are stored. You will be prompted
to enter your system password for authentication. Only application names are shown,
not the actual passwords.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all stored applications:")
		// TODO: Implement list logic
		fmt.Println("List command not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}