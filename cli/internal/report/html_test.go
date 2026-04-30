package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kest-labs/kest/cli/internal/storage"
	"github.com/kest-labs/kest/cli/internal/summary"
)

func TestWriteRecordHTML(t *testing.T) {
	t.Parallel()

	requestHeaders, err := json.Marshal(map[string]string{
		"Authorization": "Bearer demo-token",
		"Accept":        "application/json",
	})
	if err != nil {
		t.Fatalf("marshal request headers: %v", err)
	}

	responseHeaders, err := json.Marshal(map[string][]string{
		"Content-Type": {"application/json"},
	})
	if err != nil {
		t.Fatalf("marshal response headers: %v", err)
	}

	record := &storage.Record{
		ID:              42,
		Method:          "GET",
		URL:             "https://api.example.com/users?id=1",
		QueryParams:     json.RawMessage(`{"id":["1"]}`),
		RequestHeaders:  requestHeaders,
		RequestBody:     "",
		ResponseStatus:  200,
		ResponseHeaders: responseHeaders,
		ResponseBody:    `{"ok":true,"user":{"id":1}}`,
		DurationMs:      123,
		Environment:     "staging",
		Project:         "demo-project",
		CreatedAt:       time.Date(2026, time.April, 30, 9, 30, 0, 0, time.UTC),
	}

	outputPath := filepath.Join(t.TempDir(), "record.html")
	writtenPath, err := WriteRecordHTML(record, RecordHTMLOptions{OutputPath: outputPath})
	if err != nil {
		t.Fatalf("WriteRecordHTML returned error: %v", err)
	}
	if writtenPath != outputPath {
		t.Fatalf("expected path %q, got %q", outputPath, writtenPath)
	}

	html, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read HTML: %v", err)
	}

	content := string(html)
	assertContains(t, content, "Record #42")
	assertContains(t, content, "https://api.example.com/users?id=1")
	assertContains(t, content, "demo-project")
	assertContains(t, content, "ok")
	assertContains(t, content, "user")
	assertContains(t, content, "Authorization")
}

func TestWriteRunHTML(t *testing.T) {
	t.Parallel()

	summ := summary.NewSummary()
	summ.StartTime = time.Date(2026, time.April, 30, 10, 0, 0, 0, time.UTC)

	summ.AddResult(summary.TestResult{
		Name:   "Login",
		Method: "POST",
		URL:    "https://api.example.com/login",
		RequestHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		RequestBody: `{"email":"demo@example.com"}`,
		Status:      200,
		ResponseHeaders: map[string][]string{
			"Content-Type": {"application/json"},
		},
		ResponseBody: `{"token":"abc123"}`,
		Duration:     120 * time.Millisecond,
		StartTime:    time.Date(2026, time.April, 30, 10, 0, 1, 0, time.UTC),
		RecordID:     7,
		Success:      true,
	})

	summ.AddResult(summary.TestResult{
		Name:         "Generate Signature",
		Method:       "EXEC",
		Command:      `printf '{"signature":"ok"}'`,
		ResponseBody: `{"signature":"ok"}`,
		Duration:     20 * time.Millisecond,
		StartTime:    time.Date(2026, time.April, 30, 10, 0, 2, 0, time.UTC),
		Success:      true,
	})

	outputPath := filepath.Join(t.TempDir(), "run.html")
	writtenPath, err := WriteRunHTML(summ, RunHTMLOptions{
		OutputPath: outputPath,
		SourcePath: "login.flow.md",
		LogPath:    "/tmp/kest-session.log",
	})
	if err != nil {
		t.Fatalf("WriteRunHTML returned error: %v", err)
	}
	if writtenPath != outputPath {
		t.Fatalf("expected path %q, got %q", outputPath, writtenPath)
	}

	html, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read HTML: %v", err)
	}

	content := string(html)
	assertContains(t, content, "login.flow.md")
	assertContains(t, content, "Step Index")
	assertContains(t, content, "Generate Signature")
	assertContains(t, content, "printf")
	assertContains(t, content, "signature")
	assertContains(t, content, "Recorded as #7")
	assertContains(t, content, "/tmp/kest-session.log")
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(content, want) {
		t.Fatalf("expected HTML to contain %q", want)
	}
}
