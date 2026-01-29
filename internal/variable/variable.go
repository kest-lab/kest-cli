package variable

import (
	"regexp"
	"strings"
)

var varRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Interpolate replaces {{var}} with values from the map
func Interpolate(text string, vars map[string]string) string {
	return varRegex.ReplaceAllStringFunc(text, func(match string) string {
		name := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := vars[name]; ok {
			return val
		}
		return match
	})
}
