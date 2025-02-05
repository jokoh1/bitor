package version

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Version is the application version injected at build time
var Version string

// HandleGetVersion handles the GET /api/version endpoint
func HandleGetVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"version": Version,
	})
}
