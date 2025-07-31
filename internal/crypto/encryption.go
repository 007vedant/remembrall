package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32 // AES-256
	saltLength  = 16
	nonceLength = 12 // GCM nonce
	iterations  = 100000
)

// Encryptor handles encryption and decryption operations
type Encryptor struct {
	masterPassword string
}

// NewEncryptor creates a new encryptor with the master password
func NewEncryptor(masterPassword string) *Encryptor {
	return &Encryptor{masterPassword: masterPassword}
}

// deriveKey derives an encryption key from the master password using PBKDF2
func (e *Encryptor) deriveKey(salt []byte) []byte {
	return pbkdf2.Key([]byte(e.masterPassword), salt, iterations, keyLength, sha256.New)
}

// Encrypt encrypts the plaintext using AES-256-GCM
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	// Generate random salt
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from master password and salt
	key := e.deriveKey(salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine salt + nonce + ciphertext
	combined := make([]byte, 0, saltLength+nonceLength+len(ciphertext))
	combined = append(combined, salt...)
	combined = append(combined, nonce...)
	combined = append(combined, ciphertext...)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decrypts the ciphertext using AES-256-GCM
func (e *Encryptor) Decrypt(encodedCiphertext string) (string, error) {
	// Decode from base64
	combined, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check minimum length
	if len(combined) < saltLength+nonceLength {
		return "", fmt.Errorf("invalid ciphertext: too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := combined[:saltLength]
	nonce := combined[saltLength : saltLength+nonceLength]
	ciphertext := combined[saltLength+nonceLength:]

	// Derive key from master password and salt
	key := e.deriveKey(salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: invalid password or corrupted data")
	}

	return string(plaintext), nil
}

// VerifyMasterPassword verifies if the master password is correct by trying to decrypt a test value
func (e *Encryptor) VerifyMasterPassword(testCiphertext string) bool {
	_, err := e.Decrypt(testCiphertext)
	return err == nil
}