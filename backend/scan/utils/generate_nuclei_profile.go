package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
	"gopkg.in/yaml.v2"
)

func GenerateNucleiProfileYAML(app *pocketbase.PocketBase, scanID string, filePath string) error {
	// Fetch the record from the "nuclei_scans" collection
	record, err := app.Dao().FindRecordById("nuclei_scans", scanID)
	if err != nil {
		return fmt.Errorf("failed to find record: %w", err)
	}

	// Expand the necessary relations
	relations := []string{"nuclei_profile"}
	expandedData, err := ExpandRelations(app, record, relations)
	if err != nil {
		return err
	}

	// Access the expanded "profile" record
	expandedProfile, ok := expandedData["nuclei_profile"]
	if !ok {
		return fmt.Errorf("profile relation not found")
	}

	// Extract profile JSON blob
	var profileJSON types.JsonRaw
	switch profile := expandedProfile.(type) {
	case *models.Record:
		profileJSON = profile.Get("profile").(types.JsonRaw)
	case []*models.Record:
		if len(profile) > 0 {
			profileJSON = profile[0].Get("profile").(types.JsonRaw)
		} else {
			return fmt.Errorf("profile relation is empty")
		}
	default:
		return fmt.Errorf("unexpected type for profile relation")
	}

	// Convert JSON to YAML
	yamlData, err := ConvertJSONToYAML(profileJSON)
	if err != nil {
		return fmt.Errorf("failed to convert JSON to YAML: %w", err)
	}

	// Save profile YAML to the specified file path
	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write nuclei_profile.yaml: %w", err)
	}

	return nil
}

func ConvertJSONToYAML(jsonData []byte) ([]byte, error) {
	var jsonObj interface{}
	if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
		return nil, err
	}
	return yaml.Marshal(jsonObj)
}
