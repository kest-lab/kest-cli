package project

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *repository) CreateCLIToken(ctx context.Context, token *ProjectCLIToken, tokenHash string) error {
	po := newProjectCLITokenPO(token, tokenHash)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}

	token.ID = po.ID
	token.CreatedAt = po.CreatedAt
	token.UpdatedAt = po.UpdatedAt
	return nil
}

func (r *repository) GetCLITokenByHash(ctx context.Context, tokenHash string) (*ProjectCLIToken, error) {
	var po ProjectCLITokenPO
	if err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.toDomain(), nil
}

func (r *repository) TouchCLIToken(ctx context.Context, id uint, usedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&ProjectCLITokenPO{}).
		Where("id = ?", id).
		Update("last_used_at", usedAt).Error
}
