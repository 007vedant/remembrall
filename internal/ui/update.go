package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/crypto"
	"remembrall/internal/db"

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
		
		if err := updatePassword(appName); err != nil {
			exitWithError("Failed to update password: %v", err)
		}
		
		fmt.Printf("✓ Password for '%s' updated successfully!\n", appName)
	},
}

func updatePassword(appName string) error {
	// Initialize master password manager
	masterMgr, err := auth.NewMasterPasswordManager()
	if err != nil {
		return fmt.Errorf("failed to initialize master password manager: %w", err)
	}

	// Prompt and verify master password
	masterPassword, err := masterMgr.PromptAndVerifyMasterPassword()
	if err != nil {
		return fmt.Errorf("master password verification failed: %w", err)
	}

	// Initialize database store to check if entry exists
	store, err := db.NewSQLiteStore()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	// Check if the entry exists
	_, err = store.Get(appName)
	if err != nil {
		return fmt.Errorf("application '%s' not found. Use 'save' command to add new passwords", appName)
	}

	// Prompt for new application password
	fmt.Printf("Enter new password for '%s'\n", appName)
	newPassword, err := auth.PromptApplicationPassword(appName)
	if err != nil {
		return fmt.Errorf("failed to get new password: %w", err)
	}

	// Initialize encryptor with master password
	encryptor := crypto.NewEncryptor(masterPassword)
	
	// Encrypt the new password
	encryptedPassword, err := encryptor.Encrypt(newPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Update password in database
	err = store.Update(appName, encryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to update in database: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}