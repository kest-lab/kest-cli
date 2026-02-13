package user

import (
	"time"

	"gorm.io/gorm"
)

// UserPO represents the persistent object for the users table.
// It follows all ZGO database naming and column standards.
type UserPO struct {
	// 1. Primary Key
	ID uint `gorm:"primaryKey;autoIncrement"`

	// 2. Core Business Fields
	// Size: 255 is the standard for strings
	// uniqueIndex ensures database-level integrity
	Username string `gorm:"size:255;not null;uniqueIndex:idx_users_username"`
	Email    string `gorm:"size:255;not null;uniqueIndex:idx_users_email"`

	// Password is kept in PO for persistence but hidden from JSON
	Password string `gorm:"size:255;not null"`

	// 3. Optional Fields
	// Use pointers for fields that can be NULL in the DB
	AvatarURL *string `gorm:"size:512"`
	Bio       string  `gorm:"type:text"` // For long content

	// 4. Status and Enums
	// Defaults are set at the DB level for consistency
	Status int `gorm:"index:idx_users_status;default:1"`

	// 5. Audit Timestamps (Managed by GORM)
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time

	// 6. Soft Delete
	// Adding an index to DeletedAt is standard for performance
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName explicitly sets the table name to plural form.
func (UserPO) TableName() string {
	return "users"
}

// =============================================================================
// Indexing Patterns Example
// =============================================================================

// UserLoginLogPO demonstrates a composite index
type UserLoginLogPO struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index:idx_user_login_time;priority:1"`
	LoginIP   string `gorm:"size:45"`
	UserAgent string `gorm:"size:512"`

	// Priority 2 ensures user_id + created_at is efficient
	CreatedAt time.Time `gorm:"index:idx_user_login_time;priority:2"`
}

func (UserLoginLogPO) TableName() string {
	return "user_login_logs"
}
