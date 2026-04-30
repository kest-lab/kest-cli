package permission

import idpkg "github.com/kest-labs/kest/api/pkg/id"

// CreateRoleRequest is the request for creating a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

// UpdateRoleRequest is the request for updating a role
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// AssignRoleRequest is the request for assigning a role to a user
type AssignRoleRequest struct {
	UserID idpkg.Compatible `json:"user_id" binding:"required" swaggertype:"string"`
	RoleID idpkg.Compatible `json:"role_id" binding:"required" swaggertype:"string"`
}

// AssignPermissionRequest is the request for assigning a permission to a role
type AssignPermissionRequest struct {
	RoleID       idpkg.Compatible `json:"role_id" binding:"required" swaggertype:"string"`
	PermissionID idpkg.Compatible `json:"permission_id" binding:"required" swaggertype:"string"`
}

// RoleResponse is the response for role data
type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	CreatedAt   string `json:"created_at"`
}

// PermissionResponse is the response for permission data
type PermissionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Module      string `json:"module"`
}

// UserRolesResponse is the response for user roles
type UserRolesResponse struct {
	UserID string          `json:"user_id"`
	Roles  []*RoleResponse `json:"roles"`
}
