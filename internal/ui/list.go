package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/db"
	"time"

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
		if err := listPasswords(); err != nil {
			exitWithError("Failed to list passwords: %v", err)
		}
	},
}

func listPasswords() error {
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

	fmt.Printf("\nStored applications (%d total):\n", len(entries))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	for i, entry := range entries {
		fmt.Printf("%2d. %-30s (saved: %s)\n", 
			i+1, 
			entry.AppName, 
			entry.CreatedAt.Format("2006-01-02 15:04"))
		
		if entry.UpdatedAt.After(entry.CreatedAt.Add(time.Minute)) {
			fmt.Printf("    %-30s (updated: %s)\n", 
				"", 
				entry.UpdatedAt.Format("2006-01-02 15:04"))
		}
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nUse 'remembrall get <app-name>' to retrieve a password")

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}