package domain

import (
	"errors"
	"regexp"
	"strings"
)

// Value Objects are immutable objects that represent a concept by its attributes.
// They have no identity and are compared by their values.

// Email represents a validated email address
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{value: email}, nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// Username represents a validated username
type Username struct {
	value string
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// NewUsername creates a new Username value object
func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return Username{}, errors.New("username cannot be empty")
	}
	if !usernameRegex.MatchString(username) {
		return Username{}, errors.New("username must be 3-32 characters, alphanumeric with _ or -")
	}
	return Username{value: username}, nil
}

// String returns the username as a string
func (u Username) String() string {
	return u.value
}

// Equals checks if two usernames are equal
func (u Username) Equals(other Username) bool {
	return u.value == other.value
}

// Password represents a hashed password
type Password struct {
	hash string
}

// NewPassword creates a new Password from a hash
func NewPassword(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the password hash
func (p Password) Hash() string {
	return p.hash
}

// IsEmpty checks if the password is empty
func (p Password) IsEmpty() bool {
	return p.hash == ""
}

// UserStatus represents the status of a user account
type UserStatus int

const (
	UserStatusDisabled UserStatus = 0
	UserStatusActive   UserStatus = 1
	UserStatusPending  UserStatus = 2
	UserStatusBanned   UserStatus = 3
)

// String returns the string representation of the status
func (s UserStatus) String() string {
	switch s {
	case UserStatusDisabled:
		return "disabled"
	case UserStatusActive:
		return "active"
	case UserStatusPending:
		return "pending"
	case UserStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

// IsActive checks if the status represents an active account
func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

// CanLogin checks if the user can login with this status
func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive
}

// Money represents a monetary value with currency
type Money struct {
	amount   int64  // Amount in smallest unit (cents)
	currency string // ISO 4217 currency code
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency string) (Money, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	if len(currency) != 3 {
		return Money{}, errors.New("currency must be a 3-letter ISO code")
	}
	return Money{amount: amount, currency: currency}, nil
}

// Amount returns the amount in smallest unit
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// Add adds two money values (must be same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add different currencies")
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Subtract subtracts two money values (must be same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract different currencies")
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

// Equals checks if two money values are equal
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// IsZero checks if the amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive checks if the amount is positive
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative checks if the amount is negative
func (m Money) IsNegative() bool {
	return m.amount < 0
}
