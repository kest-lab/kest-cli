package events

import "errors"

// Event system errors
var (
	ErrEventBusClosed   = errors.New("event bus is closed")
	ErrEventNotFound    = errors.New("event not found")
	ErrInvalidEventType = errors.New("invalid event type")
	ErrStoreNotEnabled  = errors.New("event store is not enabled")
)
