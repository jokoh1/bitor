package utils

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func ExpandRelations(app *pocketbase.PocketBase, record *models.Record, relations []string) (map[string]interface{}, error) {
	expandErrors := app.Dao().ExpandRecord(record, relations, nil)
	if len(expandErrors) > 0 {
		for field, err := range expandErrors {
			return nil, fmt.Errorf("failed to expand relation %s: %v", field, err)
		}
	}

	expandedData := make(map[string]interface{})
	for _, relation := range relations {
		expandedField, ok := record.Expand()[relation]
		if !ok {
			return nil, fmt.Errorf("%s relation not found", relation)
		}
		expandedData[relation] = expandedField
	}

	return expandedData, nil
}
