package models

type NucleiScanArchive struct {
	ID           string `db:"id" json:"id"`
	ScanID       string `db:"scan_id" json:"scan_id"`
	ClientID     string `db:"client_id" json:"client_id"`
	S3ProviderID string `db:"s3_provider_id" json:"s3_provider_id"`
	S3FullPath   string `db:"s3_full_path" json:"s3_full_path"`
	S3SmallPath  string `db:"s3_small_path" json:"s3_small_path"`
	Created      string `db:"created" json:"created"`
	Updated      string `db:"updated" json:"updated"`
}
