package request

import (
	"encoding/json"
	"testing"
)

func TestKeyValueJSONRoundTripPreservesDisabledState(t *testing.T) {
	original := KeyValue{
		Key:     "Authorization",
		Value:   "Bearer <token>",
		Enabled: false,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("expected marshal to succeed, got %v", err)
	}

	if string(data) != `{"key":"Authorization","value":"Bearer \u003ctoken\u003e","enabled":false}` {
		t.Fatalf("expected enabled=false to be serialized, got %s", data)
	}

	var decoded KeyValue
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("expected unmarshal to succeed, got %v", err)
	}

	if decoded.Enabled {
		t.Fatal("expected enabled=false to survive round-trip")
	}
}

func TestKeyValueJSONDefaultsEnabledToTrueWhenOmitted(t *testing.T) {
	var decoded KeyValue
	if err := json.Unmarshal([]byte(`{"key":"X-Test","value":"1"}`), &decoded); err != nil {
		t.Fatalf("expected unmarshal to succeed, got %v", err)
	}

	if !decoded.Enabled {
		t.Fatal("expected missing enabled field to default to true")
	}
}
