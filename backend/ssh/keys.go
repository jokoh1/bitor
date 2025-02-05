package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/crypto/ssh"
)

// Assume ansibleBasePath is defined and initialized elsewhere in your code
var ansibleBasePath string

// SetAnsibleBasePath sets the base path for Ansible operations
func SetAnsibleBasePath(path string) {
	ansibleBasePath = path
	log.Printf("Set Ansible base path for SSH keys: %s", path)
}

// getKeysDir returns the path to the keys directory
func getKeysDir() string {
	return filepath.Join(ansibleBasePath, "keys")
}

// ValidateKeyPair checks if a private key and public key form a valid pair
func ValidateKeyPair(privateKeyData, publicKeyData []byte) error {
	// Parse the private key
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return fmt.Errorf("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	// Generate a public key from the private key
	generatedPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to generate public key from private key: %w", err)
	}

	// Compare the fingerprints
	if ssh.FingerprintSHA256(pubKey) != ssh.FingerprintSHA256(generatedPublicKey) {
		return fmt.Errorf("keys do not form a valid pair")
	}

	return nil
}

// VerifySSHKeys checks if the SSH keys in the database match those in the ansible/keys folder
func VerifySSHKeys(app *pocketbase.PocketBase) error {
	// Get keys from database
	records, err := app.Dao().FindRecordsByExpr("ansible")
	if err != nil {
		return fmt.Errorf("failed to query records: %w", err)
	}

	// Find record with SSH keys
	var record *models.Record
	for _, r := range records {
		if r.Get("ssh_public_key") != "" && r.Get("ssh_private_key") != "" {
			record = r
			break
		}
	}

	if record == nil {
		return fmt.Errorf("no valid SSH key pair found in database")
	}

	// First validate the keys in the database form a valid pair
	dbPrivateKey := []byte(record.Get("ssh_private_key").(string))
	dbPublicKey := []byte(record.Get("ssh_public_key").(string))

	if err := ValidateKeyPair(dbPrivateKey, dbPublicKey); err != nil {
		log.Printf("Invalid key pair in database, generating new keys...")
		return InitializeSSHKeys(app)
	}

	// Check if keys exist in ansible/keys directory
	keysDir := getKeysDir()
	publicKeyPath := filepath.Join(keysDir, "key.pub")
	privateKeyPath := filepath.Join(keysDir, "key")

	// Read files if they exist
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Printf("Public key file not found or unreadable: %v", err)
		return GetSSHKeys(app) // Write keys to files
	}

	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("Private key file not found or unreadable: %v", err)
		return GetSSHKeys(app) // Write keys to files
	}

	// Validate the keys in the files form a valid pair
	if err := ValidateKeyPair(privateKeyFile, publicKeyFile); err != nil {
		log.Printf("Invalid key pair in files, updating from database...")
		return GetSSHKeys(app)
	}

	// Compare keys with database
	if string(publicKeyFile) != record.Get("ssh_public_key").(string) ||
		string(privateKeyFile) != record.Get("ssh_private_key").(string) {
		log.Printf("SSH keys in files do not match database, updating files...")
		return GetSSHKeys(app) // Update files with keys from database
	}

	log.Printf("SSH keys verified successfully")
	return nil
}

func InitializeSSHKeys(app *pocketbase.PocketBase) error {
	// Check if keys already exist
	collection, err := app.Dao().FindCollectionByNameOrId("ansible")
	if err != nil {
		return fmt.Errorf("failed to find ansible collection: %w", err)
	}

	// Check for existing public key record
	records, err := app.Dao().FindRecordsByExpr("ansible")
	if err != nil {
		return fmt.Errorf("failed to query records: %w", err)
	}

	// Check if any record has a non-empty ssh_public_key
	for _, record := range records {
		if record.Get("ssh_public_key") != "" {
			log.Printf("SSH keys exist in database, verifying files...")
			return VerifySSHKeys(app)
		}
	}

	log.Printf("No SSH keys found, generating new key pair...")

	// Generate new SSH key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Convert private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyString := string(pem.EncodeToMemory(privateKeyPEM))

	// Generate public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKeyString := string(ssh.MarshalAuthorizedKey(publicKey))

	// Create a single record for both keys
	sshRecord := models.NewRecord(collection)
	sshRecord.Set("type", "ssh_keys")
	sshRecord.Set("name", "default")
	sshRecord.Set("ssh_private_key", privateKeyString)
	sshRecord.Set("ssh_public_key", publicKeyString)

	if err := app.Dao().SaveRecord(sshRecord); err != nil {
		return fmt.Errorf("failed to save SSH keys: %w", err)
	}

	log.Printf("Successfully generated and stored new SSH key pair")

	return nil
}

func GetSSHKeys(app *pocketbase.PocketBase) error {
	// Get the first record from ansible collection
	records, err := app.Dao().FindRecordsByExpr("ansible")
	if err != nil {
		return fmt.Errorf("failed to query records: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no SSH keys found in database")
	}

	// Get the first record that has both keys
	var record *models.Record
	for _, r := range records {
		if r.Get("ssh_public_key") != "" && r.Get("ssh_private_key") != "" {
			record = r
			break
		}
	}

	if record == nil {
		return fmt.Errorf("no valid SSH key pair found")
	}

	// Create ansible/keys directory if it doesn't exist
	keysDir := getKeysDir()
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Write public key
	publicKeyPath := filepath.Join(keysDir, "key.pub")
	if err := os.WriteFile(publicKeyPath, []byte(record.Get("ssh_public_key").(string)), 0600); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	// Write private key
	privateKeyPath := filepath.Join(keysDir, "key")
	if err := os.WriteFile(privateKeyPath, []byte(record.Get("ssh_private_key").(string)), 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	return nil
}
