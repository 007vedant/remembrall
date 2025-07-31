package ui

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <app-name>",
	Short: "Update a password for an application",
	Long: `Update an existing password for an application or website. You will be prompted
to enter your system password for authentication, and then the new password
to store. The password input will be hidden from the terminal.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Updating password for: %s\n", appName)
		// TODO: Implement password update logic
		fmt.Println("Update command not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}