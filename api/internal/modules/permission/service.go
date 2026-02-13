package permission

import (
	"context"
	"errors"
)

// Service defines the interface for permission operations
type Service interface {
	// Role management
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req *UpdateRoleRequest) (*RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error
	GetRole(ctx context.Context, id uint) (*RoleResponse, error)
	ListRoles(ctx context.Context) ([]*RoleResponse, error)

	// User role management
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*RoleResponse, error)

	// Permission checking
	HasPermission(ctx context.Context, userID uint, permission string) (bool, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]*PermissionResponse, error)

	// Permission management
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error
	ListPermissions(ctx context.Context) ([]*PermissionResponse, error)
}

// ServiceImpl implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new permission service
func NewService(repo Repository) *service {
	return &service{repo: repo}
}

// CreateRole creates a new role
func (s *service) CreateRole(ctx context.Context, req *CreateRoleRequest) (*RoleResponse, error) {
	role := &Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsDefault:   req.IsDefault,
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

// UpdateRole updates an existing role
func (s *service) UpdateRole(ctx context.Context, id uint, req *UpdateRoleRequest) (*RoleResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

// DeleteRole deletes a role
func (s *service) DeleteRole(ctx context.Context, id uint) error {
	return s.repo.DeleteRole(ctx, id)
}

// GetRole gets a role by ID
func (s *service) GetRole(ctx context.Context, id uint) (*RoleResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toRoleResponse(role), nil
}

// ListRoles lists all roles
func (s *service) ListRoles(ctx context.Context) ([]*RoleResponse, error) {
	roles, err := s.repo.FindAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*RoleResponse
	for _, r := range roles {
		responses = append(responses, toRoleResponse(r))
	}
	return responses, nil
}

// AssignRoleToUser assigns a role to a user
func (s *service) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

// RemoveRoleFromUser removes a role from a user
func (s *service) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return s.repo.RemoveRoleFromUser(ctx, userID, roleID)
}

// GetUserRoles gets all roles for a user
func (s *service) GetUserRoles(ctx context.Context, userID uint) ([]*RoleResponse, error) {
	roles, err := s.repo.FindRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*RoleResponse
	for _, r := range roles {
		responses = append(responses, toRoleResponse(r))
	}
	return responses, nil
}

// HasPermission checks if a user has a specific permission
func (s *service) HasPermission(ctx context.Context, userID uint, permission string) (bool, error) {
	return s.repo.HasPermission(ctx, userID, permission)
}

// GetRolePermissions gets all permissions for a role
func (s *service) GetRolePermissions(ctx context.Context, roleID uint) ([]*PermissionResponse, error) {
	perms, err := s.repo.FindPermissionsByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	var responses []*PermissionResponse
	for _, p := range perms {
		responses = append(responses, toPermissionResponse(p))
	}
	return responses, nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *service) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	return s.repo.AssignPermissionToRole(ctx, roleID, permissionID)
}

// RemovePermissionFromRole removes a permission from a role
func (s *service) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	return s.repo.RemovePermissionFromRole(ctx, roleID, permissionID)
}

// ListPermissions lists all permissions
func (s *service) ListPermissions(ctx context.Context) ([]*PermissionResponse, error) {
	perms, err := s.repo.FindAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*PermissionResponse
	for _, p := range perms {
		responses = append(responses, toPermissionResponse(p))
	}
	return responses, nil
}

// Helper functions
func toRoleResponse(r *Role) *RoleResponse {
	return &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsDefault:   r.IsDefault,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPermissionResponse(p *Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		Description: p.Description,
		Module:      p.Module,
	}
}
