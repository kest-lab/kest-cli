package events

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// EventStore interface for persisting and querying events
type EventStore interface {
	// Store persists an event
	Store(ctx context.Context, event Event) error
	// Query retrieves events matching the filter
	Query(ctx context.Context, filter EventFilter) ([]StoredEvent, error)
	// Replay replays events to a handler in chronological order
	Replay(ctx context.Context, filter EventFilter, handler EventHandler) error
}

// EventFilter defines criteria for querying events
type EventFilter struct {
	EventName     string     // Filter by event name (supports wildcards)
	CorrelationID string     // Filter by correlation ID
	Source        string     // Filter by source
	FromTime      *time.Time // Events after this time
	ToTime        *time.Time // Events before this time
	Limit         int        // Maximum number of events to return
	Offset        int        // Offset for pagination
}

// StoredEvent represents a persisted event
type StoredEvent struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EventID       string    `gorm:"uniqueIndex;size:36" json:"event_id"`
	EventName     string    `gorm:"index;size:255" json:"event_name"`
	CorrelationID string    `gorm:"index;size:36" json:"correlation_id"`
	CausationID   string    `gorm:"size:36" json:"causation_id"`
	Source        string    `gorm:"index;size:255" json:"source"`
	Payload       string    `gorm:"type:text" json:"payload"` // JSON encoded payload
	OccurredAt    time.Time `gorm:"index" json:"occurred_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (StoredEvent) TableName() string {
	return "events"
}

// GormEventStore implements EventStore using GORM
type GormEventStore struct {
	db *gorm.DB
}

// NewGormEventStore creates a new GORM-based event store
func NewGormEventStore(db *gorm.DB) *GormEventStore {
	return &GormEventStore{db: db}
}

// AutoMigrate creates the events table
func (s *GormEventStore) AutoMigrate() error {
	return s.db.AutoMigrate(&StoredEvent{})
}

// Store persists an event to the database
func (s *GormEventStore) Store(ctx context.Context, event Event) error {
	meta := event.Metadata()

	// Serialize the event payload
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	stored := StoredEvent{
		EventID:       meta.ID,
		EventName:     event.EventName(),
		CorrelationID: meta.CorrelationID,
		CausationID:   meta.CausationID,
		Source:        meta.Source,
		Payload:       string(payload),
		OccurredAt:    event.OccurredAt(),
	}

	return s.db.WithContext(ctx).Create(&stored).Error
}

// Query retrieves events matching the filter
func (s *GormEventStore) Query(ctx context.Context, filter EventFilter) ([]StoredEvent, error) {
	query := s.db.WithContext(ctx).Model(&StoredEvent{})

	if filter.EventName != "" {
		if containsWildcard(filter.EventName) {
			// Convert glob pattern to SQL LIKE pattern
			likePattern := globToLike(filter.EventName)
			query = query.Where("event_name LIKE ?", likePattern)
		} else {
			query = query.Where("event_name = ?", filter.EventName)
		}
	}

	if filter.CorrelationID != "" {
		query = query.Where("correlation_id = ?", filter.CorrelationID)
	}

	if filter.Source != "" {
		query = query.Where("source = ?", filter.Source)
	}

	if filter.FromTime != nil {
		query = query.Where("occurred_at >= ?", *filter.FromTime)
	}

	if filter.ToTime != nil {
		query = query.Where("occurred_at <= ?", *filter.ToTime)
	}

	// Order by occurred_at for chronological replay
	query = query.Order("occurred_at ASC, id ASC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var events []StoredEvent
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// Replay replays events to a handler in chronological order
func (s *GormEventStore) Replay(ctx context.Context, filter EventFilter, handler EventHandler) error {
	events, err := s.Query(ctx, filter)
	if err != nil {
		return err
	}

	for _, stored := range events {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Create a replay event from stored data
		replayEvent := &ReplayEvent{
			stored: stored,
		}

		if err := handler(ctx, replayEvent); err != nil {
			return err
		}
	}

	return nil
}

// containsWildcard checks if a pattern contains wildcard characters
func containsWildcard(pattern string) bool {
	for _, c := range pattern {
		if c == '*' || c == '?' {
			return true
		}
	}
	return false
}

// globToLike converts a glob pattern to SQL LIKE pattern
func globToLike(pattern string) string {
	result := make([]byte, 0, len(pattern)*2)
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			result = append(result, '%')
		case '?':
			result = append(result, '_')
		case '%', '_':
			// Escape SQL wildcards
			result = append(result, '\\', pattern[i])
		default:
			result = append(result, pattern[i])
		}
	}
	return string(result)
}

// ReplayEvent wraps a stored event for replay
type ReplayEvent struct {
	stored StoredEvent
}

func (e *ReplayEvent) EventName() string {
	return e.stored.EventName
}

func (e *ReplayEvent) OccurredAt() time.Time {
	return e.stored.OccurredAt
}

func (e *ReplayEvent) Metadata() EventMetadata {
	return EventMetadata{
		ID:            e.stored.EventID,
		CorrelationID: e.stored.CorrelationID,
		CausationID:   e.stored.CausationID,
		Source:        e.stored.Source,
		Timestamp:     e.stored.OccurredAt,
	}
}

// Payload returns the raw JSON payload
func (e *ReplayEvent) Payload() string {
	return e.stored.Payload
}

// UnmarshalPayload unmarshals the payload into the target
func (e *ReplayEvent) UnmarshalPayload(target interface{}) error {
	return json.Unmarshal([]byte(e.stored.Payload), target)
}

// StoringMiddleware creates a middleware that persists events to a store
func StoringMiddleware(store EventStore) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			// Store the event first
			if err := store.Store(ctx, event); err != nil {
				// Log but don't fail - event handling should continue
				// In production, you might want to handle this differently
			}
			return next(ctx, event)
		}
	}
}
