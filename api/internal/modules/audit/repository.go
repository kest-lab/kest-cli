package audit

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the interface for audit logging
type Repository interface {
	Create(ctx context.Context, log *AuditLogPO) error
	ListByProject(ctx context.Context, projectID uint, page, pageSize int) ([]AuditLogPO, int64, error)
	ListAll(ctx context.Context, page, pageSize int) ([]AuditLogPO, int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new audit repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new audit log
func (r *repository) Create(ctx context.Context, log *AuditLogPO) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ListByProject retrieves audit logs for a specific project
func (r *repository) ListByProject(ctx context.Context, projectID uint, page, pageSize int) ([]AuditLogPO, int64, error) {
	var logs []AuditLogPO
	var total int64
	offset := (page - 1) * pageSize

	q := r.db.WithContext(ctx).Model(&AuditLogPO{}).Where("project_id = ?", projectID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ListAll retrieves all audit logs with pagination
func (r *repository) ListAll(ctx context.Context, page, pageSize int) ([]AuditLogPO, int64, error) {
	var logs []AuditLogPO
	var total int64
	offset := (page - 1) * pageSize

	q := r.db.WithContext(ctx).Model(&AuditLogPO{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
