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

func TestParseSyncResponse(t *testing.T) {
	t.Run("wrapped api response", func(t *testing.T) {
		body := []byte(`{"code":0,"message":"success","data":{"created":1,"updated":2,"skipped":3,"errors":["x"]}}`)
		got, err := parseSyncResponse(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Created != 1 || got.Updated != 2 || got.Skipped != 3 || len(got.Errors) != 1 {
			t.Fatalf("unexpected response: %+v", got)
		}
	})

	t.Run("flat response", func(t *testing.T) {
		body := []byte(`{"created":4,"updated":5,"skipped":6}`)
		got, err := parseSyncResponse(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Created != 4 || got.Updated != 5 || got.Skipped != 6 {
			t.Fatalf("unexpected response: %+v", got)
		}
	})
}
