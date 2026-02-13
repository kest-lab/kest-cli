package domain

import (
	"errors"
	"time"
)

// Domain-specific errors
// These are business errors that can be returned by any layer
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountDisabled    = errors.New("account is disabled")

	// Permission errors
	ErrPermissionDenied = errors.New("permission denied")
	ErrRoleNotFound     = errors.New("role not found")

	// Generic errors
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrInvalidInput = errors.New("invalid input")
)

// Events
const (
	EventUserCreated = "user.created"
)

// UserCreatedEvent is triggered when a new user registers
type UserCreatedEvent struct {
	User       *User
	occurredAt time.Time
}

func NewUserCreatedEvent(user *User) UserCreatedEvent {
	return UserCreatedEvent{
		User:       user,
		occurredAt: time.Now(),
	}
}

func (e UserCreatedEvent) EventName() string {
	return EventUserCreated
}

func (e UserCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e UserCreatedEvent) Data() any {
	return e.User
}
