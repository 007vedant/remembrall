package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"remembrall/internal/crypto"
)

const (
	testString = "remembrall-verification-test"
	masterFile = ".remembrall-master"
)

// MasterPasswordManager handles master password operations
type MasterPasswordManager struct {
	masterFilePath string
}

// NewMasterPasswordManager creates a new master password manager
func NewMasterPasswordManager() (*MasterPasswordManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	masterFilePath := filepath.Join(homeDir, masterFile)
	return &MasterPasswordManager{masterFilePath: masterFilePath}, nil
}

// IsFirstTime checks if this is the first time running Remembrall
func (m *MasterPasswordManager) IsFirstTime() bool {
	_, err := os.Stat(m.masterFilePath)
	return os.IsNotExist(err)
}

// SetupMasterPassword sets up the master password for first-time use
func (m *MasterPasswordManager) SetupMasterPassword() (string, error) {
	if !m.IsFirstTime() {
		return "", fmt.Errorf("master password already exists")
	}

	masterPassword, err := PromptNewMasterPassword()
	if err != nil {
		return "", fmt.Errorf("failed to get master password: %w", err)
	}

	// Create encryptor with the master password
	encryptor := crypto.NewEncryptor(masterPassword)

	// Encrypt a test string to verify the password later
	encryptedTest, err := encryptor.Encrypt(testString)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt test string: %w", err)
	}

	// Save the encrypted test string to file
	err = os.WriteFile(m.masterFilePath, []byte(encryptedTest), 0600)
	if err != nil {
		return "", fmt.Errorf("failed to save master password verification: %w", err)
	}

	fmt.Println("Master password has been set up successfully!")
	return masterPassword, nil
}

// VerifyMasterPassword verifies the provided master password
func (m *MasterPasswordManager) VerifyMasterPassword(masterPassword string) error {
	if m.IsFirstTime() {
		return fmt.Errorf("master password not set up. Run any command to set it up")
	}

	// Read the encrypted test string
	encryptedTest, err := os.ReadFile(m.masterFilePath)
	if err != nil {
		return fmt.Errorf("failed to read master password file: %w", err)
	}

	// Create encryptor with the provided password
	encryptor := crypto.NewEncryptor(masterPassword)

	// Try to decrypt the test string
	decryptedTest, err := encryptor.Decrypt(string(encryptedTest))
	if err != nil {
		return fmt.Errorf("invalid master password")
	}

	// Verify it matches our test string
	if decryptedTest != testString {
		return fmt.Errorf("invalid master password")
	}

	return nil
}

// PromptAndVerifyMasterPassword prompts for master password and verifies it
func (m *MasterPasswordManager) PromptAndVerifyMasterPassword() (string, error) {
	// Check if this is first time setup
	if m.IsFirstTime() {
		return m.SetupMasterPassword()
	}

	// Prompt for existing master password
	masterPassword, err := PromptMasterPassword()
	if err != nil {
		return "", fmt.Errorf("failed to get master password: %w", err)
	}

	// Verify the password
	err = m.VerifyMasterPassword(masterPassword)
	if err != nil {
		return "", err
	}

	return masterPassword, nil
}