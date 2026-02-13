package permission

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a user role
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:255" json:"description"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
}

// TableName specifies the database table name
func (Role) TableName() string {
	return "roles"
}

// Permission represents a granular permission
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"size:100;not null;unique" json:"name"` // e.g., "users:read", "users:write"
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"size:255" json:"description"`
	Module      string         `gorm:"size:50" json:"module"` // e.g., "users", "posts"
}

// TableName specifies the database table name
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission is the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"not null;index" json:"role_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName specifies the database table name
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole is the many-to-many relationship between users and roles
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	RoleID    uint      `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the database table name
func (UserRole) TableName() string {
	return "user_roles"
}
