package fileutil

import "strings"

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
