package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/crypto"
	"remembrall/internal/db"
	"remembrall/internal/search"

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
		
		targetAppName, err := updatePassword(appName)
		if err != nil {
			exitWithError("Failed to update password: %v", err)
		}
		
		fmt.Printf("✓ Password for '%s' updated successfully!\n", targetAppName)
	},
}

func updatePassword(appName string) (string, error) {
	// Initialize master password manager
	masterMgr, err := auth.NewMasterPasswordManager()
	if err != nil {
		return "", fmt.Errorf("failed to initialize master password manager: %w", err)
	}

	// Prompt and verify master password
	masterPassword, err := masterMgr.PromptAndVerifyMasterPassword()
	if err != nil {
		return "", fmt.Errorf("master password verification failed: %w", err)
	}

	// Initialize database store to check if entry exists
	store, err := db.NewSQLiteStore()
	if err != nil {
		return "", fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	// Try exact match first
	targetAppName := appName
	_, err = store.Get(appName)
	if err != nil {
		// If exact match fails, try fuzzy search
		allEntries, listErr := store.List()
		if listErr != nil {
			return "", fmt.Errorf("application '%s' not found. Use 'save' command to add new passwords", appName)
		}

		bestMatch := search.FindBestMatch(allEntries, appName)
		if bestMatch == nil {
			// Show similar matches if available
			results := search.FuzzySearch(allEntries, appName)
			if len(results) > 0 {
				fmt.Printf("No exact match found for '%s'. Did you mean:\n", appName)
				for i, result := range results {
					if i >= 3 { // Show max 3 suggestions
						break
					}
					fmt.Printf("  • %s\n", result.Entry.AppName)
				}
				return "", fmt.Errorf("use exact application name or try 'remembrall list' to see all stored passwords")
			}
			return "", fmt.Errorf("application '%s' not found. Use 'save' command to add new passwords", appName)
		}

		// Found a good match, use it
		fmt.Printf("No exact match found for '%s'.\n", appName)
		fmt.Printf("Updating password for '%s'...\n", bestMatch.AppName)
		targetAppName = bestMatch.AppName
	}

	// Prompt for new application password
	fmt.Printf("Enter new password for '%s'\n", targetAppName)
	newPassword, err := auth.PromptApplicationPassword(targetAppName)
	if err != nil {
		return "", fmt.Errorf("failed to get new password: %w", err)
	}

	// Initialize encryptor with master password
	encryptor := crypto.NewEncryptor(masterPassword)
	
	// Encrypt the new password
	encryptedPassword, err := encryptor.Encrypt(newPassword)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Update password in database
	err = store.Update(targetAppName, encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to update in database: %w", err)
	}

	return targetAppName, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}