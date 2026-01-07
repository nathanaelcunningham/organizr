package fileutil

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseTemplate replaces placeholders in a template string with values from the vars map.
// Placeholders are in the format {key}.
func ParseTemplate(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ValidateTemplate checks if a template string only uses placeholders from the allowed list.
// Placeholders are in the format {key}. Returns an error if any placeholder is not in allowedVars.
func ValidateTemplate(template string, allowedVars []string) error {
	// Find all placeholders using regex
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	// Empty template or no placeholders is valid
	if len(matches) == 0 {
		return nil
	}

	// Create a set of allowed variables for fast lookup
	allowedSet := make(map[string]bool)
	for _, v := range allowedVars {
		allowedSet[v] = true
	}

	// Check each placeholder
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		placeholder := match[1]
		if !allowedSet[placeholder] {
			return fmt.Errorf("invalid placeholder: {%s}, allowed: %v", placeholder, allowedVars)
		}
	}

	return nil
}
