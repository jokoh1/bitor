package findings

import (
	"orbit/auth" // Adjust the import path as necessary

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// Ensure you have imported Echo's middleware package

// RegisterRoutes registers the scan routes with the authentication middleware.
func RegisterRoutes(app *pocketbase.PocketBase, e *core.ServeEvent) {
	// Create a group for the scan routes with the authentication middleware
	findingsGroup := e.Router.Group("/api/findings",
		apis.LoadAuthContext(app),     // Apply LoadAuthContext middleware first
		auth.RequireAuthOrAPIKey(app), // Use the custom middleware from the auth package
		apis.ActivityLogger(app),      // Optional: log activities
	)

	// Register scan routes within the group
	findingsGroup.GET("/grouped", HandleGroupedFindings(app))
	findingsGroup.GET("", HandleFindings(app))
	findingsGroup.POST("/bulk-update", HandleBulkUpdateFindings(app))
	findingsGroup.GET("/by-client", HandleVulnerabilitiesByClient(app))
	findingsGroup.GET("/recent", HandleRecentFindings(app))
}
