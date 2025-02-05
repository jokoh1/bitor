package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase"
)

const (
	TestValue = "orbit"
	KeySize   = 32 // AES-256 requires 32 bytes
)

// padKey ensures the key is exactly 32 bytes
func padKey(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("encryption key cannot be empty")
	}

	// If key is shorter than 32 bytes, pad it with zeros
	if len(key) < KeySize {
		paddedKey := make([]byte, KeySize)
		copy(paddedKey, []byte(key))
		return paddedKey, nil
	}

	// If key is longer than 32 bytes, truncate it
	if len(key) > KeySize {
		return []byte(key[:KeySize]), nil
	}

	// Key is exactly 32 bytes
	return []byte(key), nil
}

// Encrypt encrypts data using AES-GCM with proper key padding
func Encrypt(data []byte, key string) (string, error) {
	// Always use environment variable for encryption key
	envKey := os.Getenv("API_ENCRYPTION_KEY")
	if envKey == "" {
		return "", fmt.Errorf("API_ENCRYPTION_KEY environment variable is not set")
	}

	paddedKey, err := padKey(envKey)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM with proper key padding
func Decrypt(encryptedData string, key string) ([]byte, error) {
	// Always use environment variable for decryption key
	envKey := os.Getenv("API_ENCRYPTION_KEY")
	if envKey == "" {
		return nil, fmt.Errorf("API_ENCRYPTION_KEY environment variable is not set")
	}

	paddedKey, err := padKey(envKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	data, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plaintext, nil
}

// ValidateEncryptionKey checks if the encryption key is valid by testing it against a known value
func ValidateEncryptionKey(app *pocketbase.PocketBase) error {
	// Get the encryption key from environment
	encryptionKey := os.Getenv("API_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return fmt.Errorf("API_ENCRYPTION_KEY environment variable is not set")
	}

	// Validate key length
	if len(encryptionKey) != KeySize {
		return fmt.Errorf("API_ENCRYPTION_KEY must be exactly %d bytes (got %d bytes)", KeySize, len(encryptionKey))
	}

	// Find the PocketBase collection "system_settings"
	collection, err := app.Dao().FindCollectionByNameOrId("system_settings")
	if err != nil {
		return fmt.Errorf("failed to find system_settings collection: %w", err)
	}

	// Get the first system settings record
	record, err := app.Dao().FindFirstRecordByFilter(collection.Id, "id != ''")
	if err != nil {
		return fmt.Errorf("failed to find system settings record: %w", err)
	}

	encryptedTest := record.GetString("encyption_test")
	if encryptedTest == "" {
		// No test value exists, create one
		encryptedTest, err = Encrypt([]byte(TestValue), encryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt test value: %w", err)
		}

		// Save the encrypted test value
		record.Set("encyption_test", encryptedTest)
		if err := app.Dao().SaveRecord(record); err != nil {
			return fmt.Errorf("failed to save encryption test: %w", err)
		}
		return nil
	}

	// Test value exists, try to decrypt it
	decrypted, err := Decrypt(encryptedTest, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt test value - encryption key is invalid: %w", err)
	}

	if string(decrypted) != TestValue {
		return fmt.Errorf("decrypted test value does not match expected value - encryption key is invalid")
	}

	return nil
}
