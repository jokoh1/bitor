package models

type SignedURLRequest struct {
	ScanID   string `json:"scan_id" form:"scan_id"`
	FileType string `json:"file_type" form:"file_type"`
}
