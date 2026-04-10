package main

import "testing"

func TestBuildSpecSyncEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		projectID string
		want      string
	}{
		{
			name:      "already versioned base path",
			baseURL:   "https://api.kest.dev/v1",
			projectID: "12",
			want:      "https://api.kest.dev/v1/projects/12/cli/spec-sync",
		},
		{
			name:      "already versioned api path",
			baseURL:   "https://api.kest.dev/api/v1/",
			projectID: "12",
			want:      "https://api.kest.dev/api/v1/projects/12/cli/spec-sync",
		},
		{
			name:      "plain origin",
			baseURL:   "https://api.kest.dev",
			projectID: "12",
			want:      "https://api.kest.dev/v1/projects/12/cli/spec-sync",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSpecSyncEndpoint(tt.baseURL, tt.projectID); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
