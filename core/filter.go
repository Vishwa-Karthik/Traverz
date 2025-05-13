package core

import (
	"regexp"
)

// IsExcluded checks if a given path (basename) matches any of the exclusion patterns.
func IsExcluded(name string, patterns []*regexp.Regexp) bool {
	// Using filepath.Base(path) might be better if full paths are passed,
	// but for now, assuming 'name' is already the basename.
	for _, pattern := range patterns {
		if pattern.MatchString(name) {
			return true
		}
	}
	return false
}
