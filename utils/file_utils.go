package utils

import (
	"os"
)

// PathExists checks if a path exists and if it's a directory.
func PathExists(path string) (exists bool, isDir bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false
	}
	if err != nil {
		// Some other error (e.g., permission denied to stat)
		// Treat as not existing or not a dir for simplicity here,
		// or handle error more granularly if needed.
		return false, false
	}
	return true, info.IsDir()
}

// IsDir checks if a given path is a directory.
// Note: This is somewhat redundant with PathExists but kept for semantic clarity if needed.
func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
