package web

import (
	"embed"
	"io/fs"
)

//go:embed build/*
var static embed.FS

// NewFS returns a new fs.FS containing the embedded static files.
func NewFS() fs.FS {
	fsys, err := fs.Sub(static, "build")
	if err != nil {
		panic(err)
	}

	return fsys
}
