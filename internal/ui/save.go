package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/crypto"
	"remembrall/internal/db"

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
		
		if err := savePassword(appName); err != nil {
			exitWithError("Failed to save password: %v", err)
		}
		
		fmt.Printf("✓ Password for '%s' saved successfully!\n", appName)
	},
}

func savePassword(appName string) error {
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

	// Prompt for application password
	appPassword, err := auth.PromptApplicationPassword(appName)
	if err != nil {
		return fmt.Errorf("failed to get application password: %w", err)
	}

	// Initialize encryptor with master password
	encryptor := crypto.NewEncryptor(masterPassword)
	
	// Encrypt the application password
	encryptedPassword, err := encryptor.Encrypt(appPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Initialize database store
	store, err := db.NewSQLiteStore()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	// Save encrypted password to database
	err = store.Save(appName, encryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to save to database: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(saveCmd)
}