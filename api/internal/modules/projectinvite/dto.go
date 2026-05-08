package projectinvite

import (
	"fmt"
	"net/url"
	"time"
)

type CreateProjectInvitationRequest struct {
	Role         string     `json:"role" binding:"required,oneof=admin write read"`
	ExpiresAt    *time.Time `json:"expires_at"`
	MaxUses      *int       `json:"max_uses"`
	TargetUserID string     `json:"target_user_id"`
}

type ProjectInvitationResponse struct {
	ID            string     `json:"id"`
	ProjectID     string     `json:"project_id"`
	TokenPrefix   string     `json:"token_prefix"`
	Slug          string     `json:"slug"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	InviteURL     string     `json:"invite_url"`
	MaxUses       int        `json:"max_uses"`
	UsedCount     int        `json:"used_count"`
	RemainingUses *int       `json:"remaining_uses"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
	CreatedBy     string     `json:"created_by"`
	TargetUserID  string     `json:"target_user_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type PublicProjectInvitationResponse struct {
	ProjectID     string     `json:"project_id"`
	ProjectName   string     `json:"project_name"`
	ProjectSlug   string     `json:"project_slug"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	RemainingUses *int       `json:"remaining_uses"`
	RequiresAuth  bool       `json:"requires_auth"`
}

type AcceptedProjectInvitationMember struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type AcceptProjectInvitationResponse struct {
	ProjectID  string                          `json:"project_id"`
	Member     AcceptedProjectInvitationMember `json:"member"`
	RedirectTo string                          `json:"redirect_to"`
}

type RejectProjectInvitationResponse struct {
	Status string `json:"status"`
}

type PendingProjectInvitationResponse struct {
	ID            string     `json:"id"`
	Slug          string     `json:"slug"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	InviteURL     string     `json:"invite_url"`
	ProjectID     string     `json:"project_id"`
	ProjectName   string     `json:"project_name"`
	ProjectSlug   string     `json:"project_slug"`
	InviterID     string     `json:"inviter_id"`
	InviterName   string     `json:"inviter_name"`
	InviterEmail  string     `json:"inviter_email"`
	InviterAvatar string     `json:"inviter_avatar,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	RemainingUses *int       `json:"remaining_uses"`
}

func toProjectInvitationResponse(invitation *ProjectInvitation, now time.Time) *ProjectInvitationResponse {
	if invitation == nil {
		return nil
	}

	return &ProjectInvitationResponse{
		ID:            invitation.ID,
		ProjectID:     invitation.ProjectID,
		TokenPrefix:   invitation.TokenPrefix,
		Slug:          invitation.Slug,
		Role:          invitation.Role,
		Status:        resolveInvitationStatus(invitation, now),
		InviteURL:     buildProjectInvitationURL(invitation.Slug),
		MaxUses:       invitation.MaxUses,
		UsedCount:     invitation.UsedCount,
		RemainingUses: remainingInvitationUses(invitation),
		ExpiresAt:     invitation.ExpiresAt,
		LastUsedAt:    invitation.LastUsedAt,
		CreatedBy:     invitation.CreatedBy,
		TargetUserID:  invitation.TargetUserID,
		CreatedAt:     invitation.CreatedAt,
		UpdatedAt:     invitation.UpdatedAt,
	}
}

func (r *ProjectInvitationResponse) withBaseURL(baseURL string) *ProjectInvitationResponse {
	if r == nil {
		return nil
	}

	r.InviteURL = buildProjectInvitationURLForBase(r.Slug, baseURL)
	return r
}

func toPublicProjectInvitationResponse(
	invitation *ProjectInvitation,
	project *ProjectSummary,
	now time.Time,
) *PublicProjectInvitationResponse {
	if invitation == nil || project == nil {
		return nil
	}

	return &PublicProjectInvitationResponse{
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		ProjectSlug:   project.Slug,
		Role:          invitation.Role,
		Status:        resolveInvitationStatus(invitation, now),
		ExpiresAt:     invitation.ExpiresAt,
		RemainingUses: remainingInvitationUses(invitation),
		RequiresAuth:  true,
	}
}

func toPendingProjectInvitationResponse(
	pending *PendingProjectInvitation,
	now time.Time,
) *PendingProjectInvitationResponse {
	if pending == nil || pending.Invitation == nil || pending.Project == nil {
		return nil
	}

	inviterName := pending.Invitation.CreatedBy
	inviterEmail := ""
	inviterAvatar := ""
	if pending.Inviter != nil {
		inviterName = resolveInvitationUserDisplayName(pending.Inviter)
		inviterEmail = pending.Inviter.Email
		inviterAvatar = pending.Inviter.Avatar
	}

	return &PendingProjectInvitationResponse{
		ID:            pending.Invitation.ID,
		Slug:          pending.Invitation.Slug,
		Role:          pending.Invitation.Role,
		Status:        resolveInvitationStatus(pending.Invitation, now),
		InviteURL:     buildProjectInvitationURL(pending.Invitation.Slug),
		ProjectID:     pending.Project.ID,
		ProjectName:   pending.Project.Name,
		ProjectSlug:   pending.Project.Slug,
		InviterID:     pending.Invitation.CreatedBy,
		InviterName:   inviterName,
		InviterEmail:  inviterEmail,
		InviterAvatar: inviterAvatar,
		ExpiresAt:     pending.Invitation.ExpiresAt,
		CreatedAt:     pending.Invitation.CreatedAt,
		RemainingUses: remainingInvitationUses(pending.Invitation),
	}
}

func (r *PendingProjectInvitationResponse) withBaseURL(baseURL string) *PendingProjectInvitationResponse {
	if r == nil {
		return nil
	}

	r.InviteURL = buildProjectInvitationURLForBase(r.Slug, baseURL)
	return r
}

func resolveInvitationUserDisplayName(user *InvitationUserSummary) string {
	if user == nil {
		return ""
	}
	if user.Name != "" {
		return user.Name
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.Username != "" {
		return user.Username
	}
	if user.Email != "" {
		return user.Email
	}
	return user.ID
}

func buildProjectInvitationURL(slug string) string {
	return buildProjectInvitationURLForBase(slug, "")
}

func buildProjectInvitationURLForBase(slug, baseURL string) string {
	path := fmt.Sprintf("/invite/project/%s", url.PathEscape(slug))
	base := resolveConfiguredInvitationBaseURL()
	if base != "" {
		return base + path
	}
	base = normalizeInvitationBaseURL(baseURL, false)
	if base != "" {
		return base + path
	}
	return path
}
