package domain

import (
	"context"
	"time"
)

// User represents the core domain entity.
// JSON tags control API output - Password is hidden by default.
//
// For most cases, just use response.Success(c, user) directly.
// The Password field is automatically excluded from JSON output.
type User struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Password     string     `json:"-"` // Always hidden in JSON output
	Nickname     string     `json:"nickname,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	Bio          string     `json:"bio,omitempty"`
	Status       int        `json:"status"`
	IsSuperAdmin bool       `json:"is_super_admin"` // System-level administrator
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// IsActive returns whether the user account is active
func (u *User) IsActive() bool {
	return u.Status == 1
}

// UserRepository defines the contract for user data operations
// Implementations live in modules/user/repository.go
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*User, int64, error)
	Search(ctx context.Context, query string, limit int) ([]*User, error)
}
