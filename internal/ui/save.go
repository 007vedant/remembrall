package ui

import (
	"fmt"

	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save <app-name>",
	Short: "Save a password for an application",
	Long: `Save a password for an application or website. You will be prompted
to enter your system password for authentication, and then the password
to store. The password input will be hidden from the terminal.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Saving password for: %s\n", appName)
		// TODO: Implement password saving logic
		fmt.Println("Save command not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}