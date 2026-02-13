package apispec

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines API specification data access interface
type Repository interface {
	// API Spec operations
	CreateSpec(ctx context.Context, spec *APISpecPO) error
	GetSpecByID(ctx context.Context, id uint) (*APISpecPO, error)
	GetSpecByMethodAndPath(ctx context.Context, projectID uint, method, path string) (*APISpecPO, error)
	UpdateSpec(ctx context.Context, spec *APISpecPO) error
	DeleteSpec(ctx context.Context, id uint) error
	ListSpecs(ctx context.Context, projectID uint, version string, page, pageSize int) ([]*APISpecPO, int64, error)
	ListAllSpecs(ctx context.Context, projectID uint) ([]*APISpecPO, error)

	// API Example operations
	CreateExample(ctx context.Context, example *APIExamplePO) error
	GetExamplesBySpecID(ctx context.Context, specID uint) ([]*APIExamplePO, error)
	GetExampleByID(ctx context.Context, id uint) (*APIExamplePO, error)
	DeleteExample(ctx context.Context, id uint) error
}

// repository is the private implementation
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new API spec repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ========== API Spec Operations ==========

func (r *repository) CreateSpec(ctx context.Context, spec *APISpecPO) error {
	return r.db.WithContext(ctx).Create(spec).Error
}

func (r *repository) GetSpecByID(ctx context.Context, id uint) (*APISpecPO, error) {
	var spec APISpecPO
	if err := r.db.WithContext(ctx).First(&spec, id).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

func (r *repository) GetSpecByMethodAndPath(ctx context.Context, projectID uint, method, path string) (*APISpecPO, error) {
	var spec APISpecPO
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND method = ? AND path = ?", projectID, method, path).
		First(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

func (r *repository) UpdateSpec(ctx context.Context, spec *APISpecPO) error {
	return r.db.WithContext(ctx).Save(spec).Error
}

func (r *repository) DeleteSpec(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&APISpecPO{}, id).Error
}

func (r *repository) ListSpecs(ctx context.Context, projectID uint, version string, page, pageSize int) ([]*APISpecPO, int64, error) {
	var specs []*APISpecPO
	var total int64

	query := r.db.WithContext(ctx).Model(&APISpecPO{}).Where("project_id = ?", projectID)

	if version != "" {
		query = query.Where("version = ?", version)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&specs).Error; err != nil {
		return nil, 0, err
	}

	return specs, total, nil
}

func (r *repository) ListAllSpecs(ctx context.Context, projectID uint) ([]*APISpecPO, error) {
	var specs []*APISpecPO
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("method ASC, path ASC").
		Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

// ========== API Example Operations ==========

func (r *repository) CreateExample(ctx context.Context, example *APIExamplePO) error {
	return r.db.WithContext(ctx).Create(example).Error
}

func (r *repository) GetExamplesBySpecID(ctx context.Context, specID uint) ([]*APIExamplePO, error) {
	var examples []*APIExamplePO
	if err := r.db.WithContext(ctx).
		Where("api_spec_id = ?", specID).
		Order("id DESC").
		Find(&examples).Error; err != nil {
		return nil, err
	}
	return examples, nil
}

func (r *repository) GetExampleByID(ctx context.Context, id uint) (*APIExamplePO, error) {
	var example APIExamplePO
	if err := r.db.WithContext(ctx).First(&example, id).Error; err != nil {
		return nil, err
	}
	return &example, nil
}

func (r *repository) DeleteExample(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&APIExamplePO{}, id).Error
}
