//go:build production

package main

import (
	"embed"
	"io/fs"
)

//go:embed ansible/**/*
var ansibleEmbedFS embed.FS

func getAnsibleFS() fs.FS {
	return ansibleEmbedFS
}
