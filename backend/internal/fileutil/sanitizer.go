package fileutil

import (
	"regexp"
	"strings"
)

// SanitizePath removes invalid filesystem characters from a path string.
// It replaces invalid characters with hyphens, trims whitespace, and collapses multiple spaces.
// This function is suitable for path components like author, series, and title.
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

// SanitizeFilename removes invalid filesystem characters from a filename.
// Similar to SanitizePath but specifically for filenames (not directory paths).
// Prevents path traversal by replacing directory separators with hyphens.
func SanitizeFilename(filename string) string {
	// Use same sanitization as SanitizePath for consistency
	return SanitizePath(filename)
}
