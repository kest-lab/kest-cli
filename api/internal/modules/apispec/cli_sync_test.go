package apispec

import "testing"

func TestSanitizeExampleHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer secret",
		"X-API-Key":     "abc123",
		"Content-Type":  "application/json",
	}

	sanitized := sanitizeExampleHeaders(headers)

	if sanitized["Authorization"] != "[REDACTED]" {
		t.Fatalf("expected Authorization to be redacted, got %q", sanitized["Authorization"])
	}
	if sanitized["X-API-Key"] != "[REDACTED]" {
		t.Fatalf("expected X-API-Key to be redacted, got %q", sanitized["X-API-Key"])
	}
	if sanitized["Content-Type"] != "application/json" {
		t.Fatalf("expected Content-Type to be preserved, got %q", sanitized["Content-Type"])
	}
}

func TestSanitizeJSONStringRedactsSensitiveFields(t *testing.T) {
	input := `{"email":"user@example.com","password":"secret","nested":{"token":"abc","ok":true}}`

	sanitized := sanitizeJSONString(input)
	expected := `{"email":"user@example.com","nested":{"ok":true,"token":"[REDACTED]"},"password":"[REDACTED]"}`

	if sanitized != expected {
		t.Fatalf("expected sanitized JSON %s, got %s", expected, sanitized)
	}
}

func TestSanitizeJSONStringKeepsInvalidJSONUntouched(t *testing.T) {
	input := `{"password":`
	if got := sanitizeJSONString(input); got != input {
		t.Fatalf("expected invalid JSON to remain unchanged, got %q", got)
	}
}
