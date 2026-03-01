package member

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	AddMember(ctx context.Context, member *ProjectMemberPO) error
	UpdateMember(ctx context.Context, member *ProjectMemberPO) error
	DeleteMember(ctx context.Context, projectID, userID uint) error
	GetMember(ctx context.Context, projectID, userID uint) (*ProjectMemberPO, error)
	ListMembers(ctx context.Context, projectID uint) ([]ProjectMemberPO, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AddMember(ctx context.Context, member *ProjectMemberPO) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *repository) UpdateMember(ctx context.Context, member *ProjectMemberPO) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *repository) DeleteMember(ctx context.Context, projectID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&ProjectMemberPO{}).Error
}

func (r *repository) GetMember(ctx context.Context, projectID, userID uint) (*ProjectMemberPO, error) {
	var member ProjectMemberPO
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Preload("User").
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *repository) ListMembers(ctx context.Context, projectID uint) ([]ProjectMemberPO, error) {
	var members []ProjectMemberPO
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("User").
		Find(&members).Error
	return members, err
}
