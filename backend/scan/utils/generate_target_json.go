package utils

import (
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func GenerateTargetsJSON(app *pocketbase.PocketBase, scanID string, filePath string) error {
	// Fetch the record from the "nuclei_scans" collection
	record, err := app.Dao().FindRecordById("nuclei_scans", scanID)
	if err != nil {
		return fmt.Errorf("failed to find record: %w", err)
	}

	// Expand the necessary relations
	relations := []string{"nuclei_targets"}
	expandedData, err := ExpandRelations(app, record, relations)
	if err != nil {
		return err
	}

	// Access the expanded "targets" record
	expandedTargets, ok := expandedData["nuclei_targets"]
	if !ok {
		return fmt.Errorf("targets relation not found")
	}

	// Extract targets JSON blob
	var targetsJSON types.JsonRaw
	switch targets := expandedTargets.(type) {
	case *models.Record:
		targetsJSON = targets.Get("targets").(types.JsonRaw)
	case []*models.Record:
		if len(targets) > 0 {
			targetsJSON = targets[0].Get("targets").(types.JsonRaw)
		} else {
			return fmt.Errorf("targets relation is empty")
		}
	default:
		return fmt.Errorf("unexpected type for targets relation")
	}

	// Save targets JSON to the specified file path
	if err := os.WriteFile(filePath, targetsJSON, 0644); err != nil {
		return fmt.Errorf("failed to write targets.json: %w", err)
	}

	return nil
}
