package domain

import (
	"context"
	"sync"
)

// AggregateRoot is the base for all aggregate roots.
// It provides domain event collection and dispatch capabilities.
type AggregateRoot struct {
	events []Event
	mu     sync.Mutex
}

// AddEvent adds a domain event to be dispatched later
func (a *AggregateRoot) AddEvent(event Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, event)
}

// GetEvents returns all pending domain events
func (a *AggregateRoot) GetEvents() []Event {
	a.mu.Lock()
	defer a.mu.Unlock()
	events := make([]Event, len(a.events))
	copy(events, a.events)
	return events
}

// ClearEvents clears all pending domain events
func (a *AggregateRoot) ClearEvents() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = nil
}

// DispatchEvents dispatches all pending events and clears them
func (a *AggregateRoot) DispatchEvents(ctx context.Context) error {
	events := a.GetEvents()
	for _, event := range events {
		if err := Dispatch(ctx, event); err != nil {
			return err
		}
	}
	a.ClearEvents()
	return nil
}

// DispatchEventsAsync dispatches all pending events asynchronously and clears them
func (a *AggregateRoot) DispatchEventsAsync(ctx context.Context) {
	events := a.GetEvents()
	for _, event := range events {
		DispatchAsync(ctx, event)
	}
	a.ClearEvents()
}

// UserAggregate is the aggregate root for User domain
type UserAggregate struct {
	AggregateRoot
	User *User
}

// NewUserAggregate creates a new user aggregate
func NewUserAggregate(user *User) *UserAggregate {
	return &UserAggregate{User: user}
}

// Register handles user registration business logic
func (a *UserAggregate) Register(username, email, hashedPassword string) {
	a.User = &User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Status:   int(UserStatusActive),
	}
}

// MarkRegistered marks the user as registered and adds the event
func (a *UserAggregate) MarkRegistered() {
	a.AddEvent(NewUserRegisteredEvent(a.User.ID, a.User.Username, a.User.Email))
}

// MarkLoggedIn marks the user as logged in and adds the event
func (a *UserAggregate) MarkLoggedIn(ipAddress string) {
	a.AddEvent(NewUserLoggedInEvent(a.User.ID, a.User.Username, ipAddress))
}

// ChangePassword changes the user's password
func (a *UserAggregate) ChangePassword(newHashedPassword string) {
	a.User.Password = newHashedPassword
	a.AddEvent(NewUserPasswordChangedEvent(a.User.ID))
}

// Delete marks the user for deletion
func (a *UserAggregate) Delete() {
	a.AddEvent(NewUserDeletedEvent(a.User.ID, a.User.Email))
}

// Disable disables the user account
func (a *UserAggregate) Disable() error {
	if a.User.Status == int(UserStatusDisabled) {
		return ErrAccountDisabled
	}
	a.User.Status = int(UserStatusDisabled)
	return nil
}

// Enable enables the user account
func (a *UserAggregate) Enable() {
	a.User.Status = int(UserStatusActive)
}

// IsActive checks if the user is active
func (a *UserAggregate) IsActive() bool {
	return UserStatus(a.User.Status).IsActive()
}

// CanLogin checks if the user can login
func (a *UserAggregate) CanLogin() bool {
	return UserStatus(a.User.Status).CanLogin()
}
