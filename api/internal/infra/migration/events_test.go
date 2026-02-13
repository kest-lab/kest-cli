package migration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Unit tests for migration events
// Test event name methods return correct strings
// Test event fields are properly set
// Requirements: 7.1-7.6

func TestMigrationsStarted_EventName(t *testing.T) {
	event := NewMigrationsStarted("up")
	assert.Equal(t, EventMigrationsStarted, event.EventName())
	assert.Equal(t, "migration.migrations_started", event.EventName())
}

func TestMigrationsStarted_Fields(t *testing.T) {
	event := NewMigrationsStarted("up")
	assert.Equal(t, "up", event.Direction)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestMigrationsStarted_DirectionDown(t *testing.T) {
	event := NewMigrationsStarted("down")
	assert.Equal(t, "down", event.Direction)
}

func TestMigrationsEnded_EventName(t *testing.T) {
	event := NewMigrationsEnded("up")
	assert.Equal(t, EventMigrationsEnded, event.EventName())
	assert.Equal(t, "migration.migrations_ended", event.EventName())
}

func TestMigrationsEnded_Fields(t *testing.T) {
	event := NewMigrationsEnded("down")
	assert.Equal(t, "down", event.Direction)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestMigrationStarted_EventName(t *testing.T) {
	event := NewMigrationStarted("2024_01_01_000000_create_users_table", "up")
	assert.Equal(t, EventMigrationStarted, event.EventName())
	assert.Equal(t, "migration.migration_started", event.EventName())
}

func TestMigrationStarted_Fields(t *testing.T) {
	event := NewMigrationStarted("2024_01_01_000000_create_users_table", "up")
	assert.Equal(t, "2024_01_01_000000_create_users_table", event.Migration)
	assert.Equal(t, "up", event.Method)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestMigrationStarted_MethodDown(t *testing.T) {
	event := NewMigrationStarted("2024_01_01_000000_create_users_table", "down")
	assert.Equal(t, "down", event.Method)
}

func TestMigrationEnded_EventName(t *testing.T) {
	event := NewMigrationEnded("2024_01_01_000000_create_users_table", "up")
	assert.Equal(t, EventMigrationEnded, event.EventName())
	assert.Equal(t, "migration.migration_ended", event.EventName())
}

func TestMigrationEnded_Fields(t *testing.T) {
	event := NewMigrationEnded("2024_01_01_000000_create_users_table", "down")
	assert.Equal(t, "2024_01_01_000000_create_users_table", event.Migration)
	assert.Equal(t, "down", event.Method)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestMigrationSkipped_EventName(t *testing.T) {
	event := NewMigrationSkipped("2024_01_01_000000_create_users_table")
	assert.Equal(t, EventMigrationSkipped, event.EventName())
	assert.Equal(t, "migration.migration_skipped", event.EventName())
}

func TestMigrationSkipped_Fields(t *testing.T) {
	event := NewMigrationSkipped("2024_01_01_000000_create_users_table")
	assert.Equal(t, "2024_01_01_000000_create_users_table", event.Migration)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestNoPendingMigrations_EventName(t *testing.T) {
	event := NewNoPendingMigrations("up")
	assert.Equal(t, EventNoPendingMigrations, event.EventName())
	assert.Equal(t, "migration.no_pending", event.EventName())
}

func TestNoPendingMigrations_Fields(t *testing.T) {
	event := NewNoPendingMigrations("up")
	assert.Equal(t, "up", event.Direction)
	assert.Equal(t, "migration", event.Metadata().Source)
	assert.NotEmpty(t, event.Metadata().ID)
	assert.False(t, event.OccurredAt().IsZero())
}

func TestNoPendingMigrations_DirectionDown(t *testing.T) {
	event := NewNoPendingMigrations("down")
	assert.Equal(t, "down", event.Direction)
}

// Test that events have unique IDs
func TestEvents_UniqueIDs(t *testing.T) {
	event1 := NewMigrationsStarted("up")
	event2 := NewMigrationsStarted("up")
	assert.NotEqual(t, event1.Metadata().ID, event2.Metadata().ID)
}

// Test that events have reasonable timestamps
func TestEvents_Timestamps(t *testing.T) {
	before := time.Now()
	event := NewMigrationStarted("test_migration", "up")
	after := time.Now()

	assert.True(t, event.OccurredAt().After(before) || event.OccurredAt().Equal(before))
	assert.True(t, event.OccurredAt().Before(after) || event.OccurredAt().Equal(after))
}

// Test all event constants are unique
func TestEventConstants_Unique(t *testing.T) {
	constants := []string{
		EventMigrationsStarted,
		EventMigrationsEnded,
		EventMigrationStarted,
		EventMigrationEnded,
		EventMigrationSkipped,
		EventNoPendingMigrations,
	}

	seen := make(map[string]bool)
	for _, c := range constants {
		assert.False(t, seen[c], "Duplicate event constant: %s", c)
		seen[c] = true
	}
}

// Test event constants have correct prefix
func TestEventConstants_Prefix(t *testing.T) {
	assert.Contains(t, EventMigrationsStarted, "migration.")
	assert.Contains(t, EventMigrationsEnded, "migration.")
	assert.Contains(t, EventMigrationStarted, "migration.")
	assert.Contains(t, EventMigrationEnded, "migration.")
	assert.Contains(t, EventMigrationSkipped, "migration.")
	assert.Contains(t, EventNoPendingMigrations, "migration.")
}
