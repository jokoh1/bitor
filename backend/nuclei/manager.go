// Package nuclei handles nuclei templates management
package nuclei

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var templatesDir string

// SetTemplatesDir sets the directory where nuclei templates will be stored
func SetTemplatesDir(dir string) {
	templatesDir = dir
}

// GetTemplatesDir returns the current templates directory
func GetTemplatesDir() string {
	return templatesDir
}

// copyDir recursively copies a directory tree
func copyDir(src, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Skip if it's a symlink
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			// Copy the file
			srcFile, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			if _, err = io.Copy(dstFile, srcFile); err != nil {
				return err
			}
		}
	}
	return nil
}

// cleanupUnwantedFiles removes files and directories we don't need
func cleanupUnwantedFiles(dir string) error {
	// List of template directories we want to keep
	wantedDirs := map[string]bool{
		"http":         true,
		"dns":          true,
		"network":      true,
		"file":         true,
		"ssl":          true,
		"headless":     true,
		"websocket":    true,
		"workflows":    true,
		"helpers":      true,
		"technologies": true,
	}

	// Read all entries in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Remove everything except the wanted directories
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)

		// If it's not in our wanted list, remove it
		if !wantedDirs[name] {
			log.Printf("Removing unwanted file/directory: %s", path)
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove %s: %v", path, err)
			}
		}
	}
	return nil
}

// InitializeTemplates ensures nuclei templates are available
func InitializeTemplates() error {
	if templatesDir == "" {
		return fmt.Errorf("templates directory not set")
	}

	// Set up the public templates directory
	publicDir := filepath.Join(templatesDir, "public")
	customDir := filepath.Join(templatesDir, "custom")

	// Create both directories
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		return fmt.Errorf("failed to create public templates directory: %v", err)
	}
	if err := os.MkdirAll(customDir, 0755); err != nil {
		return fmt.Errorf("failed to create custom templates directory: %v", err)
	}

	// Check if public templates already exist
	if _, err := os.Stat(publicDir); err == nil && hasTemplateFiles(publicDir) {
		log.Printf("Public templates directory already exists and contains files at %s, skipping clone", publicDir)
		return nil
	}

	// Create a temporary directory for cloning
	tmpDir, err := os.MkdirTemp("", "nuclei-templates-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed: %v", err)
	}

	log.Printf("Cloning nuclei templates to temporary directory: %s", tmpDir)

	// Create command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Clone with depth 1 to temporary directory
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "https://github.com/projectdiscovery/nuclei-templates.git", tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout while cloning nuclei templates")
		}
		return fmt.Errorf("failed to clone nuclei templates: %v\nOutput: %s", err, output)
	}

	// Clean up unwanted files in temporary directory
	if err := cleanupUnwantedFiles(tmpDir); err != nil {
		return fmt.Errorf("failed to clean up unwanted files: %v", err)
	}

	// Copy the cleaned templates to the public directory
	log.Printf("Copying templates from %s to %s", tmpDir, publicDir)
	if err := copyDir(tmpDir, publicDir); err != nil {
		return fmt.Errorf("failed to copy templates to public directory: %v", err)
	}

	log.Printf("Successfully initialized nuclei templates at %s", templatesDir)
	return nil
}

// hasTemplateFiles checks if a directory contains template files
func hasTemplateFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check for common template directories
			if entry.Name() == "http" || entry.Name() == "dns" || entry.Name() == "network" {
				return true
			}
		}
	}
	return false
}
