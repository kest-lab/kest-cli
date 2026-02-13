package workspace

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Workspace types
const (
	TypePersonal = "personal"
	TypeTeam     = "team"
	TypePublic   = "public"
)

// Visibility types
const (
	VisibilityPrivate = "private"
	VisibilityTeam    = "team"
	VisibilityPublic  = "public"
)

// WorkspacePO is the persistent object for database operations
type WorkspacePO struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;not null"`
	Slug        string `gorm:"size:50;uniqueIndex"`
	Description string `gorm:"size:500"`
	Type        string `gorm:"size:20;not null;default:'personal'"` // personal|team|public
	OwnerID     uint   `gorm:"not null;index"`                      // Creator
	Visibility  string `gorm:"size:20;default:'private'"`           // private|team|public
	Settings    string `gorm:"type:text"`                           // JSON settings as string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the database table name
func (WorkspacePO) TableName() string {
	return "workspaces"
}

// WorkspaceMemberPO represents a membership of a user in a workspace
type WorkspaceMemberPO struct {
	ID          uint   `gorm:"primaryKey"`
	WorkspaceID uint   `gorm:"index;uniqueIndex:idx_workspace_user;not null"`
	UserID      uint   `gorm:"index;uniqueIndex:idx_workspace_user;not null"`
	Role        string `gorm:"size:20;not null"` // owner|admin|editor|viewer
	InvitedBy   uint   `gorm:"index"`            // Inviter user ID
	JoinedAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the database table name
func (WorkspaceMemberPO) TableName() string {
	return "workspace_members"
}

// Role constants
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// RoleLevel defines the hierarchy of roles
var RoleLevel = map[string]int{
	RoleOwner:  40,
	RoleAdmin:  30,
	RoleEditor: 20,
	RoleViewer: 10,
}

// Workspace is the domain entity
type Workspace struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	OwnerID     uint                   `json:"owner_id"`
	Visibility  string                 `json:"visibility"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// WorkspaceMember is the domain entity for membership
type WorkspaceMember struct {
	ID          uint      `json:"id"`
	WorkspaceID uint      `json:"workspace_id"`
	UserID      uint      `json:"user_id"`
	Role        string    `json:"role"`
	InvitedBy   uint      `json:"invited_by"`
	JoinedAt    time.Time `json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// toDomain converts WorkspacePO to Workspace domain entity
func (po *WorkspacePO) toDomain() *Workspace {
	if po == nil {
		return nil
	}

	var settings map[string]interface{}
	if po.Settings != "" {
		// Parse JSON string to map (ignore errors, use empty map)
		_ = json.Unmarshal([]byte(po.Settings), &settings)
	}

	return &Workspace{
		ID:          po.ID,
		Name:        po.Name,
		Slug:        po.Slug,
		Description: po.Description,
		Type:        po.Type,
		OwnerID:     po.OwnerID,
		Visibility:  po.Visibility,
		Settings:    settings,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// toDomainList converts a slice of WorkspacePO to Workspace slice
func toDomainList(poList []*WorkspacePO) []*Workspace {
	result := make([]*Workspace, len(poList))
	for i, po := range poList {
		result[i] = po.toDomain()
	}
	return result
}

// toMemberDomain converts WorkspaceMemberPO to WorkspaceMember
func (po *WorkspaceMemberPO) toMemberDomain() *WorkspaceMember {
	if po == nil {
		return nil
	}
	return &WorkspaceMember{
		ID:          po.ID,
		WorkspaceID: po.WorkspaceID,
		UserID:      po.UserID,
		Role:        po.Role,
		InvitedBy:   po.InvitedBy,
		JoinedAt:    po.JoinedAt,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// toMemberDomainList converts slice of members
func toMemberDomainList(poList []*WorkspaceMemberPO) []*WorkspaceMember {
	result := make([]*WorkspaceMember, len(poList))
	for i, po := range poList {
		result[i] = po.toMemberDomain()
	}
	return result
}
