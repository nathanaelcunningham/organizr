package fileutil

import (
	"regexp"
	"strings"
)

// SanitizePath removes invalid filesystem characters from a path string.
// It replaces invalid characters with hyphens, trims whitespace, and collapses multiple spaces.
func SanitizePath(path string) string {
	// Remove invalid filesystem characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := path
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "-")
	}

	// Trim whitespace
	result = strings.TrimSpace(result)

	// Collapse multiple spaces to single space
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")

	return result
}
