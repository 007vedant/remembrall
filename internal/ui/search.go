package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/db"
	"remembrall/internal/search"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for applications using fuzzy matching",
	Long: `Search for stored applications using fuzzy matching. This allows you to find 
applications even with partial names, typos, or approximate matches. You will be 
prompted to enter your system password for authentication.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		
		if err := searchPasswords(query); err != nil {
			exitWithError("Failed to search passwords: %v", err)
		}
	},
}

func searchPasswords(query string) error {
	// Initialize master password manager
	masterMgr, err := auth.NewMasterPasswordManager()
	if err != nil {
		return fmt.Errorf("failed to initialize master password manager: %w", err)
	}

	// Prompt and verify master password
	_, err = masterMgr.PromptAndVerifyMasterPassword()
	if err != nil {
		return fmt.Errorf("master password verification failed: %w", err)
	}

	// Initialize database store
	store, err := db.NewSQLiteStore()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	// Get all password entries
	entries, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to retrieve from database: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No passwords stored yet.")
		fmt.Println("Use 'remembrall save <app-name>' to add your first password.")
		return nil
	}

	// Perform fuzzy search
	results := search.FuzzySearch(entries, query)
	
	if len(results) == 0 {
		fmt.Printf("No matches found for '%s'.\n", query)
		fmt.Println("Use 'remembrall list' to see all stored applications.")
		return nil
	}

	fmt.Printf("\nSearch results for '%s' (%d matches):\n", query, len(results))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, result := range results {
		fmt.Printf("%2d. %s\n", i+1, result.Entry.AppName)
		
		if i >= 9 { // Show max 10 results
			fmt.Printf("    ... and %d more matches\n", len(results)-10)
			break
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nUse 'remembrall get <app-name>' to retrieve a password")

	return nil
}

func init() {
	rootCmd.AddCommand(searchCmd)
}