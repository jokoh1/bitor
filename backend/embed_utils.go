package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func writeAnsibleFiles(ansibleFS fs.FS, ansibleBasePath string) error {
	// Use fs.WalkDir to traverse the ansibleFS and copy files to ansibleBasePath
	return fs.WalkDir(ansibleFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(ansibleBasePath, path)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			fileData, err := fs.ReadFile(ansibleFS, path)
			if err != nil {
				return err
			}
			if err := os.WriteFile(destPath, fileData, os.ModePerm); err != nil {
				return err
			}
		}
		return nil
	})
}
