package parser

import (
	"strings"
	"testing"
)

func TestGenerateCurlExampleAvoidsDuplicateVersionPrefix(t *testing.T) {
	curl := generateCurlExample(Endpoint{
		Route: Route{
			Method:   "GET",
			Path:     "/v1/projects/:id/categories/:cid",
			IsPublic: true,
		},
	}, APIConfig{
		LocalURL: "http://localhost:8025/api/v1",
	})

	if strings.Contains(curl, "/api/v1/v1/") {
		t.Fatalf("expected curl example to avoid duplicate version prefix, got %q", curl)
	}

	expectedURL := "http://localhost:8025/api/v1/projects/1/categories/1"
	if !strings.Contains(curl, expectedURL) {
		t.Fatalf("expected curl example to contain %q, got %q", expectedURL, curl)
	}
}
