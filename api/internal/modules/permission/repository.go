package permission

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for permission data operations
type Repository interface {
	// Role operations
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uint) error
	FindRoleByID(ctx context.Context, id uint) (*Role, error)
	FindRoleByName(ctx context.Context, name string) (*Role, error)
	FindAllRoles(ctx context.Context) ([]*Role, error)
	FindDefaultRole(ctx context.Context) (*Role, error)

	// Permission operations
	CreatePermission(ctx context.Context, perm *Permission) error
	FindAllPermissions(ctx context.Context) ([]*Permission, error)
	FindPermissionsByModule(ctx context.Context, module string) ([]*Permission, error)

	// Role-Permission operations
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error
	FindPermissionsByRoleID(ctx context.Context, roleID uint) ([]*Permission, error)

	// User-Role operations
	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	FindRolesByUserID(ctx context.Context, userID uint) ([]*Role, error)
	HasPermission(ctx context.Context, userID uint, permissionName string) (bool, error)
}

// RepositoryImpl implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new permission repository
func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// CreateRole creates a new role
func (r *repository) CreateRole(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// UpdateRole updates an existing role
func (r *repository) UpdateRole(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// DeleteRole deletes a role by ID
func (r *repository) DeleteRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}

// FindRoleByID finds a role by ID
func (r *repository) FindRoleByID(ctx context.Context, id uint) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindRoleByName finds a role by name
func (r *repository) FindRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindAllRoles returns all roles
func (r *repository) FindAllRoles(ctx context.Context) ([]*Role, error) {
	var roles []*Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// FindDefaultRole returns the default role
func (r *repository) FindDefaultRole(ctx context.Context) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("is_default = ?", true).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreatePermission creates a new permission
func (r *repository) CreatePermission(ctx context.Context, perm *Permission) error {
	return r.db.WithContext(ctx).Create(perm).Error
}

// FindAllPermissions returns all permissions
func (r *repository) FindAllPermissions(ctx context.Context) ([]*Permission, error) {
	var perms []*Permission
	if err := r.db.WithContext(ctx).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// FindPermissionsByModule returns permissions by module
func (r *repository) FindPermissionsByModule(ctx context.Context, module string) ([]*Permission, error) {
	var perms []*Permission
	if err := r.db.WithContext(ctx).Where("module = ?", module).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *repository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	rp := &RolePermission{RoleID: roleID, PermissionID: permissionID}
	return r.db.WithContext(ctx).FirstOrCreate(rp, rp).Error
}

// RemovePermissionFromRole removes a permission from a role
func (r *repository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{}).Error
}

// FindPermissionsByRoleID returns permissions for a role
func (r *repository) FindPermissionsByRoleID(ctx context.Context, roleID uint) ([]*Permission, error) {
	var perms []*Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

// AssignRoleToUser assigns a role to a user
func (r *repository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	ur := &UserRole{UserID: userID, RoleID: roleID}
	return r.db.WithContext(ctx).FirstOrCreate(ur, ur).Error
}

// RemoveRoleFromUser removes a role from a user
func (r *repository) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

// FindRolesByUserID returns roles for a user
func (r *repository) FindRolesByUserID(ctx context.Context, userID uint) ([]*Role, error) {
	var roles []*Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// HasPermission checks if a user has a specific permission
func (r *repository) HasPermission(ctx context.Context, userID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.name = ?", userID, permissionName).
		Count(&count).Error
	return count > 0, err
}
