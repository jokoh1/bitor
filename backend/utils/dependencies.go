package utils

import (
	"fmt"
	"log"
	"os/exec"
)

// CheckDependencies checks if the required dependencies are installed.
func CheckDependencies() error {
	dependencies := []string{"terraform", "ansible"}

	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			return fmt.Errorf("%s is not installed", dep)
		}
	}

	log.Println("All dependencies are installed.")
	return nil
}
