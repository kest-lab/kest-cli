package member

import (
	"context"
	"fmt"
)

type Service interface {
	AddMember(ctx context.Context, projectID uint, req *AddMemberRequest) (*MemberResponse, error)
	UpdateMember(ctx context.Context, projectID, userID uint, req *UpdateMemberRequest) (*MemberResponse, error)
	RemoveMember(ctx context.Context, projectID, userID uint) error
	ListMembers(ctx context.Context, projectID uint) ([]MemberResponse, error)
	GetMember(ctx context.Context, projectID, userID uint) (*MemberResponse, error)
	CheckPermission(ctx context.Context, projectID, userID uint, requiredRole string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddMember(ctx context.Context, projectID uint, req *AddMemberRequest) (*MemberResponse, error) {
	// Check if already a member
	existing, _ := s.repo.GetMember(ctx, projectID, req.UserID)
	if existing != nil {
		return nil, fmt.Errorf("user is already a member of this project")
	}

	po := &ProjectMemberPO{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}

	if err := s.repo.AddMember(ctx, po); err != nil {
		return nil, err
	}

	return FromMemberPO(po), nil
}

func (s *service) UpdateMember(ctx context.Context, projectID, userID uint, req *UpdateMemberRequest) (*MemberResponse, error) {
	po, err := s.repo.GetMember(ctx, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}

	po.Role = req.Role
	if err := s.repo.UpdateMember(ctx, po); err != nil {
		return nil, err
	}

	return FromMemberPO(po), nil
}

func (s *service) RemoveMember(ctx context.Context, projectID, userID uint) error {
	po, err := s.repo.GetMember(ctx, projectID, userID)
	if err != nil {
		return fmt.Errorf("member not found")
	}

	// Owner protection - optionally check if it's the last owner
	if po.Role == RoleOwner {
		members, _ := s.repo.ListMembers(ctx, projectID)
		ownersCount := 0
		for _, m := range members {
			if m.Role == RoleOwner {
				ownersCount++
			}
		}
		if ownersCount <= 1 {
			return fmt.Errorf("cannot remove the last owner of the project")
		}
	}

	return s.repo.DeleteMember(ctx, projectID, userID)
}

func (s *service) ListMembers(ctx context.Context, projectID uint) ([]MemberResponse, error) {
	pos, err := s.repo.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return FromMemberPOs(pos), nil
}

func (s *service) GetMember(ctx context.Context, projectID, userID uint) (*MemberResponse, error) {
	po, err := s.repo.GetMember(ctx, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}
	return FromMemberPO(po), nil
}

func (s *service) CheckPermission(ctx context.Context, projectID, userID uint, requiredRole string) (bool, error) {
	po, err := s.repo.GetMember(ctx, projectID, userID)
	if err != nil {
		// No membership means no permissions (unless it's a public project, but here we enforce membership)
		return false, nil
	}

	userLevel := RoleLevel[po.Role]
	reqLevel := RoleLevel[requiredRole]

	return userLevel >= reqLevel, nil
}
