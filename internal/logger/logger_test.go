package logger

import (
	"os"
	"strings"
	"testing"
)

func TestSessionLogging(t *testing.T) {
	sessionName := "test_session"
	sl, err := StartSession(sessionName)
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer sl.File.Close()

	logMsg := "Hello Kest Log"
	LogToSession("%s", logMsg)

	path := GetSessionPath()
	if path == "" {
		t.Error("Expected session path to be non-empty")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), logMsg) {
		t.Errorf("Log content does not contain expected message: %s", logMsg)
	}

	EndSession()

	if GetSessionPath() != "" {
		t.Error("Expected session path to be empty after EndSession")
	}

	// Cleanup
	os.Remove(path)
}
