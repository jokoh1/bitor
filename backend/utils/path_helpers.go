package utils

import "path/filepath"

func GetScanPath(basePath, scanID, filename string) string {
	return filepath.Join(basePath, "scans", scanID, filename)
}
