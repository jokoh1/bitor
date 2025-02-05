package setup

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// AddPasswordChangeMiddleware adds middleware to check if a user needs to change their password
func AddPasswordChangeMiddleware(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				path := c.Path()

				// Skip middleware for:
				// 1. Authentication endpoints
				// 2. Password change endpoints
				// 3. Static files
				// 4. Setup page
				if path == "/api/collections/users/auth-with-password" ||
					path == "/api/collections/users/records/:id" ||
					strings.HasPrefix(path, "/_") ||
					!strings.HasPrefix(path, "/api/") ||
					strings.Contains(path, "change-password") {
					return next(c)
				}

				// Get the authenticated record
				authRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
				if authRecord == nil {
					return next(c)
				}

				// Check if password change is required
				requireChange, _ := authRecord.Get("requirePasswordChange").(bool)
				if requireChange {
					return apis.NewForbiddenError("Password change required", nil)
				}

				return next(c)
			}
		})
		return nil
	})
}
