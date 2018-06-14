package main

import (
	"path/filepath"
	"strings"
)

func getBaseName(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
