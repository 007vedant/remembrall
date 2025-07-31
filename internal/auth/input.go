package auth

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ReadPassword reads a password from stdin without echoing it to the terminal
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	
	// Get the file descriptor for stdin
	fd := int(syscall.Stdin)
	
	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("not running in a terminal")
	}
	
	// Read password without echo
	bytePassword, err := term.ReadPassword(fd)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	
	// Print newline since ReadPassword doesn't echo the Enter key
	fmt.Println()
	
	password := string(bytePassword)
	password = strings.TrimSpace(password)
	
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	
	return password, nil
}

// ReadPasswordWithConfirmation reads a password twice and ensures they match
func ReadPasswordWithConfirmation(prompt, confirmPrompt string) (string, error) {
	password, err := ReadPassword(prompt)
	if err != nil {
		return "", err
	}
	
	confirmation, err := ReadPassword(confirmPrompt)
	if err != nil {
		return "", err
	}
	
	if password != confirmation {
		return "", fmt.Errorf("passwords do not match")
	}
	
	return password, nil
}

// PromptMasterPassword prompts for the master password
func PromptMasterPassword() (string, error) {
	return ReadPassword("Enter your master password: ")
}

// PromptApplicationPassword prompts for an application password
func PromptApplicationPassword(appName string) (string, error) {
	prompt := fmt.Sprintf("Enter password for %s: ", appName)
	return ReadPassword(prompt)
}

// PromptNewMasterPassword prompts for a new master password with confirmation
func PromptNewMasterPassword() (string, error) {
	fmt.Println("Setting up master password for Remembrall...")
	return ReadPasswordWithConfirmation(
		"Enter your new master password: ",
		"Confirm your master password: ",
	)
}

// ClearScreen clears the terminal screen (for security after displaying passwords)
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// PrintAndClear prints a message and clears it after a delay
func PrintAndClear(message string) {
	fmt.Print(message)
	// Note: The actual clearing with timer will be implemented in the UI layer
}