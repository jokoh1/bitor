package models

type ScanRequest struct {
	ScanID         string `json:"scan_id"`
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
	NucleiTargets  string `json:"nuclei_targets"`
	NucleiProfile  string `json:"nuclei_profile"`
	NucleiInteract string `json:"nuclei_interact"`

	// Provider references
	ComputeProvider     string `json:"compute_provider"`      // Provider with "compute" use
	StateBucketProvider string `json:"state_bucket_provider"` // Provider with "terraform_storage" use
	ScanBucketProvider  string `json:"scan_bucket_provider"`  // Provider with "scan_storage" use

	// Optional settings overrides
	StatefilePath string `json:"statefile_path,omitempty"` // Override default statefile path
	ScansPath     string `json:"scans_path,omitempty"`     // Override default scans path
}
