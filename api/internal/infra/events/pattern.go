package events

import "strings"

// matchPattern checks if an event name matches a subscription pattern.
// Supports glob-style patterns with dot-separated segments:
//   - "user.*" matches "user.created", "user.deleted" (single segment wildcard)
//   - "*.created" matches "user.created", "order.created"
//   - "*" matches any single segment like "user" or "created"
//   - "**" matches any number of segments (including zero)
//   - "user.**" matches "user", "user.created", "user.profile.updated"
//   - "user.created" matches exactly "user.created"
//
// Pattern rules:
//   - Patterns and event names are dot-separated
//   - "*" matches exactly one segment
//   - "**" matches zero or more segments
func matchPattern(pattern, eventName string) bool {
	// Exact match fast path
	if pattern == eventName {
		return true
	}

	// Handle "**" (match everything)
	if pattern == "**" {
		return true
	}

	patternParts := strings.Split(pattern, ".")
	eventParts := strings.Split(eventName, ".")

	return matchParts(patternParts, eventParts)
}

// matchParts recursively matches pattern parts against event parts
func matchParts(pattern, event []string) bool {
	pi, ei := 0, 0

	for pi < len(pattern) {
		if pi >= len(pattern) {
			return ei >= len(event)
		}

		p := pattern[pi]

		switch p {
		case "**":
			// "**" at the end matches everything remaining
			if pi == len(pattern)-1 {
				return true
			}

			// Try matching "**" with 0, 1, 2, ... segments
			for i := ei; i <= len(event); i++ {
				if matchParts(pattern[pi+1:], event[i:]) {
					return true
				}
			}
			return false

		case "*":
			// "*" matches exactly one segment
			if ei >= len(event) {
				return false
			}
			pi++
			ei++

		default:
			// Literal match
			if ei >= len(event) || p != event[ei] {
				return false
			}
			pi++
			ei++
		}
	}

	// Pattern exhausted, check if event is also exhausted
	return ei >= len(event)
}

// ValidatePattern checks if a pattern is valid
func ValidatePattern(pattern string) bool {
	if pattern == "" {
		return false
	}

	parts := strings.Split(pattern, ".")
	for _, part := range parts {
		if part == "" {
			return false // Empty segment (e.g., "user..created")
		}
	}
	return true
}

// NormalizeEventName ensures an event name follows conventions
func NormalizeEventName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
