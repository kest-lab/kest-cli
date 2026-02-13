package envelope

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseEnvelope(t *testing.T) {
	// Sample envelope from Sentry Go SDK
	// Note: length must match payload exactly
	payload := `{"level":"error","message":"Test error message"}`
	envelope := `{"event_id":"abc123","sent_at":"2026-01-05T20:00:00Z","dsn":"https://key@sentry.io/1"}
{"type":"event","length":` + fmt.Sprintf("%d", len(payload)) + `}
` + payload + `
`

	env, err := ParseBytes([]byte(envelope))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	// Check header
	if env.Header.EventID != "abc123" {
		t.Errorf("Header.EventID = %q, want %q", env.Header.EventID, "abc123")
	}
	if env.Header.Dsn != "https://key@sentry.io/1" {
		t.Errorf("Header.Dsn = %q, want %q", env.Header.Dsn, "https://key@sentry.io/1")
	}

	// Check items
	if len(env.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(env.Items))
	}

	item := env.Items[0]
	if item.Header.Type != ItemTypeEvent {
		t.Errorf("Item.Type = %q, want %q", item.Header.Type, ItemTypeEvent)
	}
}

func TestParseEnvelopeMultipleItems(t *testing.T) {
	payload1 := `{"message":"first event"}`
	payload2 := `{"op":"http","description":"GET"}`
	envelope := `{"event_id":"multi123"}
{"type":"event","length":` + fmt.Sprintf("%d", len(payload1)) + `}
` + payload1 + `
{"type":"transaction","length":` + fmt.Sprintf("%d", len(payload2)) + `}
` + payload2 + `
`

	env, err := ParseBytes([]byte(envelope))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	if len(env.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(env.Items))
	}

	if env.Items[0].Header.Type != ItemTypeEvent {
		t.Errorf("Items[0].Type = %q, want event", env.Items[0].Header.Type)
	}
	if env.Items[1].Header.Type != ItemTypeTransaction {
		t.Errorf("Items[1].Type = %q, want transaction", env.Items[1].Header.Type)
	}
}

func TestParseEnvelopeNoLength(t *testing.T) {
	// Envelope without explicit length (reads until newline)
	envelope := `{"event_id":"nolen123"}
{"type":"event"}
{"message":"no length specified"}
`

	env, err := ParseBytes([]byte(envelope))
	if err != nil {
		t.Fatalf("ParseBytes() error = %v", err)
	}

	if len(env.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(env.Items))
	}
}

func TestExtractPublicKeyFromAuth(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "standard format",
			header: "Sentry sentry_version=7, sentry_client=sentry.go/0.30.0, sentry_key=abc123",
			want:   "abc123",
		},
		{
			name:   "with secret key",
			header: "Sentry sentry_version=7, sentry_key=public, sentry_secret=secret",
			want:   "public",
		},
		{
			name:   "empty header",
			header: "",
			want:   "",
		},
		{
			name:   "no sentry_key",
			header: "Sentry sentry_version=7, sentry_client=test",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPublicKeyFromAuth(tt.header)
			if result != tt.want {
				t.Errorf("ExtractPublicKeyFromAuth() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestParseEvent(t *testing.T) {
	eventJSON := `{
		"event_id": "evt123",
		"timestamp": "2026-01-05T20:00:00Z",
		"platform": "go",
		"level": "error",
		"message": "Something went wrong",
		"server_name": "api-server-1",
		"environment": "production",
		"tags": {"version": "1.0.0"},
		"exception": [{
			"type": "RuntimeError",
			"value": "division by zero"
		}]
	}`

	event, err := ParseEvent([]byte(eventJSON))
	if err != nil {
		t.Fatalf("ParseEvent() error = %v", err)
	}

	if event.EventID != "evt123" {
		t.Errorf("EventID = %q, want %q", event.EventID, "evt123")
	}
	if event.Platform != "go" {
		t.Errorf("Platform = %q, want %q", event.Platform, "go")
	}
	if event.Level != "error" {
		t.Errorf("Level = %q, want %q", event.Level, "error")
	}
	if event.Environment != "production" {
		t.Errorf("Environment = %q, want %q", event.Environment, "production")
	}
	if len(event.Exception) != 1 {
		t.Fatalf("len(Exception) = %d, want 1", len(event.Exception))
	}
	if event.Exception[0].Type != "RuntimeError" {
		t.Errorf("Exception[0].Type = %q, want %q", event.Exception[0].Type, "RuntimeError")
	}
}

func TestEventGetFingerprint(t *testing.T) {
	// Event with explicit fingerprint
	event1 := &Event{
		Fingerprint: []string{"custom", "fingerprint"},
	}
	fp1 := event1.GetFingerprint()
	if len(fp1) != 2 || fp1[0] != "custom" {
		t.Errorf("GetFingerprint() with explicit = %v", fp1)
	}

	// Event with exception (should use exception type/value)
	event2 := &Event{
		Exception: []Exception{{Type: "ValueError", Value: "invalid input"}},
	}
	fp2 := event2.GetFingerprint()
	if len(fp2) != 2 || fp2[0] != "ValueError" {
		t.Errorf("GetFingerprint() with exception = %v", fp2)
	}

	// Event with only message
	event3 := &Event{
		Message: "something happened",
	}
	fp3 := event3.GetFingerprint()
	if len(fp3) != 1 || fp3[0] != "something happened" {
		t.Errorf("GetFingerprint() with message = %v", fp3)
	}
}

func TestEnvelopeSerialize(t *testing.T) {
	env := NewEnvelope("test123")
	env.AddItem(NewItem(ItemTypeEvent, []byte(`{"level":"error"}`)))

	data, err := env.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "test123") {
		t.Errorf("Serialized envelope should contain event_id")
	}
	if !strings.Contains(result, "event") {
		t.Errorf("Serialized envelope should contain item type")
	}
}
