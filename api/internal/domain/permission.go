package domain

import (
	"context"
	"time"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uint
	Name        string
	DisplayName string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          uint
	Name        string
	DisplayName string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole represents the association between a user and a role
type UserRole struct {
	UserID    uint
	RoleID    uint
	CreatedAt time.Time
}

// RolePermission represents the association between a role and a permission
type RolePermission struct {
	RoleID       uint
	PermissionID uint
	CreatedAt    time.Time
}

// RoleRepository defines the contract for role data operations
type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	FindAll(ctx context.Context) ([]*Role, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Role, error)
}

// PermissionRepository defines the contract for permission data operations
type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Permission, error)
	FindByName(ctx context.Context, name string) (*Permission, error)
	FindAll(ctx context.Context) ([]*Permission, error)
	FindByRoleID(ctx context.Context, roleID uint) ([]*Permission, error)
	FindByUserID(ctx context.Context, userID uint) ([]*Permission, error)
}

// UserRoleRepository defines the contract for user-role association operations
type UserRoleRepository interface {
	Assign(ctx context.Context, userID, roleID uint) error
	Revoke(ctx context.Context, userID, roleID uint) error
	HasRole(ctx context.Context, userID, roleID uint) (bool, error)
	GetUserRoles(ctx context.Context, userID uint) ([]*Role, error)
}

// RolePermissionRepository defines the contract for role-permission association operations
type RolePermissionRepository interface {
	Grant(ctx context.Context, roleID, permissionID uint) error
	Revoke(ctx context.Context, roleID, permissionID uint) error
	HasPermission(ctx context.Context, roleID, permissionID uint) (bool, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]*Permission, error)
}

// AuthorizationService handles authorization logic
type AuthorizationService struct {
	userRoleRepo       UserRoleRepository
	rolePermissionRepo RolePermissionRepository
	permissionRepo     PermissionRepository
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(
	userRoleRepo UserRoleRepository,
	rolePermissionRepo RolePermissionRepository,
	permissionRepo PermissionRepository,
) *AuthorizationService {
	return &AuthorizationService{
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
		permissionRepo:     permissionRepo,
	}
}

// HasRole checks if a user has a specific role
func (s *AuthorizationService) HasRole(ctx context.Context, userID, roleID uint) (bool, error) {
	return s.userRoleRepo.HasRole(ctx, userID, roleID)
}

// HasPermission checks if a user has a specific permission (through any of their roles)
func (s *AuthorizationService) HasPermission(ctx context.Context, userID uint, permissionName string) (bool, error) {
	// Get user's permissions
	permissions, err := s.permissionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p.Name == permissionName {
			return true, nil
		}
	}
	return false, nil
}

// HasAnyPermission checks if a user has any of the specified permissions
func (s *AuthorizationService) HasAnyPermission(ctx context.Context, userID uint, permissionNames ...string) (bool, error) {
	permissions, err := s.permissionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	permSet := make(map[string]bool)
	for _, p := range permissions {
		permSet[p.Name] = true
	}

	for _, name := range permissionNames {
		if permSet[name] {
			return true, nil
		}
	}
	return false, nil
}

// HasAllPermissions checks if a user has all of the specified permissions
func (s *AuthorizationService) HasAllPermissions(ctx context.Context, userID uint, permissionNames ...string) (bool, error) {
	permissions, err := s.permissionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	permSet := make(map[string]bool)
	for _, p := range permissions {
		permSet[p.Name] = true
	}

	for _, name := range permissionNames {
		if !permSet[name] {
			return false, nil
		}
	}
	return true, nil
}
