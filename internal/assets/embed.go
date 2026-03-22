// Package assets provides access to all embedded project assets
//go:build !noembed

package assets

import "embed"

// All embeddable files from the configured asset roots
//
//go:embed embed/*
var FS embed.FS

// Example accessors

// GetFile returns content of a file preserving original path
func GetFile(subPath string) ([]byte, error) {
	return FS.ReadFile("embed/" + subPath)
}

// MustGetFile panics if file is missing (useful for init/config)
func MustGetFile(subPath string) []byte {
	data, err := GetFile(subPath)
	if err != nil {
		panic("missing embedded asset: embed/" + subPath + " → " + err.Error())
	}
	return data
}
