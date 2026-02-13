package idgen

import (
	"strings"
	"testing"
)

func TestUUID(t *testing.T) {
	id := UUID()
	if len(id) != 36 {
		t.Errorf("UUID length = %d, want 36", len(id))
	}
	if strings.Count(id, "-") != 4 {
		t.Error("UUID should contain 4 dashes")
	}
}

func TestSnowflake(t *testing.T) {
	id1 := Snowflake()
	id2 := Snowflake()

	if id1 == id2 {
		t.Error("Snowflake IDs should be unique")
	}

	t.Logf("Snowflake ID: %s", id1)
}

func TestSnowflakeInt64(t *testing.T) {
	id := SnowflakeInt64()
	if id <= 0 {
		t.Error("Snowflake int64 should be positive")
	}
	t.Logf("Snowflake int64: %d", id)
}

func TestShortID(t *testing.T) {
	id := ShortID()
	if len(id) != 8 {
		t.Errorf("ShortID length = %d, want 8", len(id))
	}
	t.Logf("ShortID: %s", id)
}

func TestNanoID(t *testing.T) {
	id := NanoID()
	if len(id) != 21 {
		t.Errorf("NanoID length = %d, want 21", len(id))
	}
	t.Logf("NanoID: %s", id)
}

func TestRandomHex(t *testing.T) {
	hex := RandomHex(32)
	if len(hex) != 32 {
		t.Errorf("RandomHex length = %d, want 32", len(hex))
	}
	t.Logf("RandomHex: %s", hex)
}

func TestUniqueness(t *testing.T) {
	gen, _ := NewSnowflakeGenerator(1)
	seen := make(map[string]bool)
	count := 10000

	for i := 0; i < count; i++ {
		id := gen.Generate()
		if seen[id] {
			t.Errorf("Duplicate ID found: %s", id)
		}
		seen[id] = true
	}
}
