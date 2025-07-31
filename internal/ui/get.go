package ui

import (
	"fmt"
	"remembrall/internal/auth"
	"remembrall/internal/crypto"
	"remembrall/internal/db"
	"time"

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
		
		if err := getPassword(appName); err != nil {
			exitWithError("Failed to retrieve password: %v", err)
		}
	},
}

func getPassword(appName string) error {
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

	// Initialize database store
	store, err := db.NewSQLiteStore()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	// Get encrypted password from database
	entry, err := store.Get(appName)
	if err != nil {
		return fmt.Errorf("failed to retrieve from database: %w", err)
	}

	// Initialize encryptor with master password
	encryptor := crypto.NewEncryptor(masterPassword)
	
	// Decrypt the password
	decryptedPassword, err := encryptor.Decrypt(entry.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Display password with timeout
	fmt.Printf("\nPassword for '%s':\n", appName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  %s\n", decryptedPassword)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("\nPassword will be cleared in 5 seconds...")
	
	// Wait for 5 seconds
	time.Sleep(5 * time.Second)
	
	// Clear the screen
	auth.ClearScreen()
	fmt.Println("Password cleared for security.")

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}