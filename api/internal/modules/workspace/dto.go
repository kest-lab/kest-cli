package workspace

// DTO (Data Transfer Objects) for workspace module

// CreateWorkspaceRequest represents the request to create a workspace
type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Slug        string `json:"slug" binding:"required,max=50,alphanum"`
	Description string `json:"description" binding:"max=500"`
	Type        string `json:"type" binding:"required,oneof=personal team public"`
	Visibility  string `json:"visibility" binding:"oneof=private team public"`
}

// UpdateWorkspaceRequest represents the request to update a workspace
type UpdateWorkspaceRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Visibility  *string `json:"visibility" binding:"omitempty,oneof=private team public"`
}

// AddMemberRequest represents the request to add a member
type AddMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=owner admin editor viewer"`
}

// UpdateMemberRoleRequest represents the request to update member role
type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner admin editor viewer"`
}

// WorkspaceResponse represents the workspace response
type WorkspaceResponse struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	OwnerID     uint                   `json:"owner_id"`
	Visibility  string                 `json:"visibility"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// WorkspaceMemberResponse represents the member response
type WorkspaceMemberResponse struct {
	ID          uint   `json:"id"`
	WorkspaceID uint   `json:"workspace_id"`
	UserID      uint   `json:"user_id"`
	Username    string `json:"username,omitempty"` // Joined from users table
	Email       string `json:"email,omitempty"`    // Joined from users table
	Role        string `json:"role"`
	InvitedBy   uint   `json:"invited_by"`
	JoinedAt    string `json:"joined_at"`
}

// FromWorkspace converts domain Workspace to WorkspaceResponse
func FromWorkspace(w *Workspace) *WorkspaceResponse {
	if w == nil {
		return nil
	}
	return &WorkspaceResponse{
		ID:          w.ID,
		Name:        w.Name,
		Slug:        w.Slug,
		Description: w.Description,
		Type:        w.Type,
		OwnerID:     w.OwnerID,
		Visibility:  w.Visibility,
		Settings:    w.Settings,
		CreatedAt:   w.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   w.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// FromWorkspaceList converts workspace list to response list
func FromWorkspaceList(workspaces []*Workspace) []*WorkspaceResponse {
	result := make([]*WorkspaceResponse, len(workspaces))
	for i, w := range workspaces {
		result[i] = FromWorkspace(w)
	}
	return result
}

// FromMember converts domain WorkspaceMember to response
func FromMember(m *WorkspaceMember) *WorkspaceMemberResponse {
	if m == nil {
		return nil
	}
	return &WorkspaceMemberResponse{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        m.Role,
		InvitedBy:   m.InvitedBy,
		JoinedAt:    m.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// FromMemberList converts member list to response list
func FromMemberList(members []*WorkspaceMember) []*WorkspaceMemberResponse {
	result := make([]*WorkspaceMemberResponse, len(members))
	for i, m := range members {
		result[i] = FromMember(m)
	}
	return result
}
