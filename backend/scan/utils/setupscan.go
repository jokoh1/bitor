package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
)

// Helper function to set up scan files.
func SetupScanFiles(yamlContent, yamlFile, targetsFile, profileFile, logDir string, app *pocketbase.PocketBase, scanID string) error {
	// Ensure the log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Ensure the directory for yamlFile exists
	yamlDir := filepath.Dir(yamlFile)
	if err := os.MkdirAll(yamlDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for YAML file: %v", err)
	}

	// Ensure the directory for targetsFile exists
	targetsDir := filepath.Dir(targetsFile)
	if err := os.MkdirAll(targetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for targets file: %v", err)
	}

	// Ensure the directory for profileFile exists
	profileDir := filepath.Dir(profileFile)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for profile file: %v", err)
	}

	// Save YAML to a file
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %v", err)
	}

	// Save targets JSON to a file
	if err := GenerateTargetsJSON(app, scanID, targetsFile); err != nil {
		return fmt.Errorf("failed to generate targets JSON: %v", err)
	}

	// Save nuclei profile YAML to a file
	if err := GenerateNucleiProfileYAML(app, scanID, profileFile); err != nil {
		return fmt.Errorf("failed to generate nuclei profile YAML: %v", err)
	}

	return nil
}
