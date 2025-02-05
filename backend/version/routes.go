package version

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes registers the version routes
func RegisterRoutes(e *echo.Echo) {
	e.GET("/api/version", HandleGetVersion)
}
