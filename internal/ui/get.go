package ui

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <app-name>",
	Short: "Retrieve a password for an application",
	Long: `Retrieve a password for an application or website. You will be prompted
to enter your system password for authentication. The password will be displayed
for 5 seconds and then cleared from the terminal.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		fmt.Printf("Retrieving password for: %s\n", appName)
		// TODO: Implement password retrieval logic
		fmt.Println("Get command not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}