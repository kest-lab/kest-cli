package events

import (
	"time"
)

// Event represents a system event.
type Event interface {
	// EventName returns the unique identifier for the event
	EventName() string
	// OccurredAt returns when the event happened
	OccurredAt() time.Time
}

// DataEvent is an event that carries a payload
type DataEvent interface {
	Event
	Data() any
}
