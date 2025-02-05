//go:build !production

package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func getAnsibleFS() fs.FS {
	ansiblePath := filepath.Join("..", "ansible")
	return os.DirFS(ansiblePath)
}
