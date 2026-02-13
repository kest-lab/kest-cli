package workspace

import (
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

// Service defines the workspace business logic interface
type Service interface {
	// Workspace operations
	CreateWorkspace(req *CreateWorkspaceRequest, ownerID uint) (*Workspace, error)
	UpdateWorkspace(id uint, req *UpdateWorkspaceRequest, userID uint, isSuperAdmin bool) (*Workspace, error)
	DeleteWorkspace(id uint, userID uint, isSuperAdmin bool) error
	GetWorkspace(id uint, userID uint, isSuperAdmin bool) (*Workspace, error)
	ListWorkspaces(userID uint, isSuperAdmin bool) ([]*Workspace, error)

	// Member operations
	AddMember(workspaceID uint, req *AddMemberRequest, inviterID uint, isSuperAdmin bool) error
	RemoveMember(workspaceID, targetUserID, requesterID uint, isSuperAdmin bool) error
	UpdateMemberRole(workspaceID, targetUserID uint, role string, requesterID uint, isSuperAdmin bool) error
	ListMembers(workspaceID, userID uint, isSuperAdmin bool) ([]*WorkspaceMember, error)

	// Permission checks
	CheckUserRole(workspaceID, userID uint, isSuperAdmin bool) (string, error)
	HasPermission(workspaceID, userID uint, requiredRole string, isSuperAdmin bool) (bool, error)
}

// service implements Service interface
type service struct {
	repo Repository
}

// NewService creates a new workspace service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateWorkspace creates a new workspace
func (s *service) CreateWorkspace(req *CreateWorkspaceRequest, ownerID uint) (*Workspace, error) {
	// Generate slug if not provided or sanitize it
	workspaceSlug := req.Slug
	if workspaceSlug == "" {
		workspaceSlug = slug.Make(req.Name)
	} else {
		workspaceSlug = slug.Make(workspaceSlug)
	}

	// Check if slug already exists
	existing, _ := s.repo.FindBySlug(workspaceSlug)
	if existing != nil {
		return nil, errors.New("workspace slug already exists")
	}

	// Set default visibility based on type
	visibility := req.Visibility
	if visibility == "" {
		if req.Type == TypePersonal {
			visibility = VisibilityPrivate
		} else {
			visibility = VisibilityTeam
		}
	}

	// Create workspace
	workspace := &WorkspacePO{
		Name:        req.Name,
		Slug:        workspaceSlug,
		Description: req.Description,
		Type:        req.Type,
		OwnerID:     ownerID,
		Visibility:  visibility,
	}

	if err := s.repo.Create(workspace); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Add owner as first member
	member := &WorkspaceMemberPO{
		WorkspaceID: workspace.ID,
		UserID:      ownerID,
		Role:        RoleOwner,
		InvitedBy:   ownerID,
		JoinedAt:    time.Now(),
	}

	if err := s.repo.AddMember(member); err != nil {
		// Rollback workspace creation if adding member fails
		s.repo.Delete(workspace.ID)
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	return workspace.toDomain(), nil
}

// UpdateWorkspace updates a workspace
func (s *service) UpdateWorkspace(id uint, req *UpdateWorkspaceRequest, userID uint, isSuperAdmin bool) (*Workspace, error) {
	// Check permission (only owner and admin can update, or super admin)
	hasPermission, err := s.repo.HasPermission(id, userID, RoleAdmin, isSuperAdmin)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("insufficient permissions to update workspace")
	}

	// Get existing workspace
	workspace, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		workspace.Name = *req.Name
	}
	if req.Description != nil {
		workspace.Description = *req.Description
	}
	if req.Visibility != nil {
		workspace.Visibility = *req.Visibility
	}

	if err := s.repo.Update(workspace); err != nil {
		return nil, err
	}

	return workspace.toDomain(), nil
}

// DeleteWorkspace deletes a workspace (only owner or super admin can delete)
func (s *service) DeleteWorkspace(id uint, userID uint, isSuperAdmin bool) error {
	workspace, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Only owner or super admin can delete
	if !isSuperAdmin && workspace.OwnerID != userID {
		return errors.New("only workspace owner can delete workspace")
	}

	return s.repo.Delete(id)
}

// GetWorkspace gets a workspace by ID
func (s *service) GetWorkspace(id uint, userID uint, isSuperAdmin bool) (*Workspace, error) {
	// Check if user has access
	hasPermission, err := s.repo.HasPermission(id, userID, RoleViewer, isSuperAdmin)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("workspace not found or access denied")
	}

	workspace, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return workspace.toDomain(), nil
}

// ListWorkspaces lists all workspaces accessible to a user
func (s *service) ListWorkspaces(userID uint, isSuperAdmin bool) ([]*Workspace, error) {
	workspaces, err := s.repo.ListByUserID(userID, isSuperAdmin)
	if err != nil {
		return nil, err
	}

	return toDomainList(workspaces), nil
}

// AddMember adds a member to a workspace
func (s *service) AddMember(workspaceID uint, req *AddMemberRequest, inviterID uint, isSuperAdmin bool) error {
	// Only admin and owner can add members (or super admin)
	hasPermission, err := s.repo.HasPermission(workspaceID, inviterID, RoleAdmin, isSuperAdmin)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions to add members")
	}

	// Check if user is already a member
	existing, _ := s.repo.FindMember(workspaceID, req.UserID)
	if existing != nil {
		return errors.New("user is already a member of this workspace")
	}

	member := &WorkspaceMemberPO{
		WorkspaceID: workspaceID,
		UserID:      req.UserID,
		Role:        req.Role,
		InvitedBy:   inviterID,
		JoinedAt:    time.Now(),
	}

	return s.repo.AddMember(member)
}

// RemoveMember removes a member from a workspace
func (s *service) RemoveMember(workspaceID, targetUserID, requesterID uint, isSuperAdmin bool) error {
	// Get workspace to check owner
	workspace, err := s.repo.FindByID(workspaceID)
	if err != nil {
		return err
	}

	// Cannot remove owner
	if workspace.OwnerID == targetUserID {
		return errors.New("cannot remove workspace owner")
	}

	// Super admin can remove anyone (except owner)
	if isSuperAdmin {
		return s.repo.RemoveMember(workspaceID, targetUserID)
	}

	// Regular users: owner and admin can remove members
	hasPermission, err := s.repo.HasPermission(workspaceID, requesterID, RoleAdmin, false)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions to remove members")
	}

	return s.repo.RemoveMember(workspaceID, targetUserID)
}

// UpdateMemberRole updates a member's role
func (s *service) UpdateMemberRole(workspaceID, targetUserID uint, role string, requesterID uint, isSuperAdmin bool) error {
	// Get workspace
	workspace, err := s.repo.FindByID(workspaceID)
	if err != nil {
		return err
	}

	// Cannot change owner's role
	if workspace.OwnerID == targetUserID {
		return errors.New("cannot change workspace owner's role")
	}

	// Super admin can change any role (except owner)
	if isSuperAdmin {
		return s.repo.UpdateMemberRole(workspaceID, targetUserID, role)
	}

	// Regular users: only owner and admin can change roles
	hasPermission, err := s.repo.HasPermission(workspaceID, requesterID, RoleAdmin, false)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions to update member roles")
	}

	return s.repo.UpdateMemberRole(workspaceID, targetUserID, role)
}

// ListMembers lists all members of a workspace
func (s *service) ListMembers(workspaceID, userID uint, isSuperAdmin bool) ([]*WorkspaceMember, error) {
	// Check if user has access to workspace
	hasPermission, err := s.repo.HasPermission(workspaceID, userID, RoleViewer, isSuperAdmin)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("workspace not found or access denied")
	}

	members, err := s.repo.ListMembers(workspaceID)
	if err != nil {
		return nil, err
	}

	return toMemberDomainList(members), nil
}

// CheckUserRole returns the user's role in a workspace
func (s *service) CheckUserRole(workspaceID, userID uint, isSuperAdmin bool) (string, error) {
	return s.repo.CheckPermission(workspaceID, userID, isSuperAdmin)
}

// HasPermission checks if a user has at least the required role level
func (s *service) HasPermission(workspaceID, userID uint, requiredRole string, isSuperAdmin bool) (bool, error) {
	return s.repo.HasPermission(workspaceID, userID, requiredRole, isSuperAdmin)
}
