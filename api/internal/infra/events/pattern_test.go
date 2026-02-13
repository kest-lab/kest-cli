package events

import "testing"

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		eventName string
		want      bool
	}{
		// Exact matches
		{"exact match", "user.created", "user.created", true},
		{"exact no match", "user.created", "user.deleted", false},

		// Single wildcard (*)
		{"single wildcard suffix", "user.*", "user.created", true},
		{"single wildcard suffix 2", "user.*", "user.deleted", true},
		{"single wildcard prefix", "*.created", "user.created", true},
		{"single wildcard prefix 2", "*.created", "order.created", true},
		{"single wildcard middle", "user.*.done", "user.profile.done", true},
		{"single wildcard no match too many", "user.*", "user.profile.updated", false},
		{"single wildcard no match prefix", "*.created", "user.profile.created", false},

		// Double wildcard (**)
		{"double wildcard all", "**", "user.created", true},
		{"double wildcard all single", "**", "user", true},
		{"double wildcard all multi", "**", "user.profile.updated", true},
		{"double wildcard suffix", "user.**", "user.created", true},
		{"double wildcard suffix multi", "user.**", "user.profile.updated", true},
		{"double wildcard suffix exact", "user.**", "user", true},
		{"double wildcard prefix", "**.created", "user.created", true},
		{"double wildcard prefix multi", "**.created", "user.profile.created", true},
		{"double wildcard middle", "user.**.done", "user.profile.done", true},
		{"double wildcard middle multi", "user.**.done", "user.profile.settings.done", true},

		// Mixed wildcards
		{"mixed wildcards", "user.*.updated", "user.profile.updated", true},
		{"mixed wildcards no match", "user.*.updated", "user.profile.settings.updated", false},

		// Edge cases
		{"empty event", "user.*", "", false},
		{"single segment", "*", "user", true},
		{"single segment no match", "*", "user.created", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPattern(tt.pattern, tt.eventName)
			if got != tt.want {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.eventName, got, tt.want)
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	tests := []struct {
		pattern string
		valid   bool
	}{
		{"user.created", true},
		{"user.*", true},
		{"**", true},
		{"user.**", true},
		{"", false},
		{"user..created", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got := ValidatePattern(tt.pattern)
			if got != tt.valid {
				t.Errorf("ValidatePattern(%q) = %v, want %v", tt.pattern, got, tt.valid)
			}
		})
	}
}
