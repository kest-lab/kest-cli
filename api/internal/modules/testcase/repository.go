package testcase

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository defines the interface for test case data access
type Repository interface {
	Create(ctx context.Context, tc *TestCasePO) error
	GetByID(ctx context.Context, id uint) (*TestCasePO, error)
	List(ctx context.Context, filter *ListFilter) ([]*TestCasePO, int64, error)
	Update(ctx context.Context, tc *TestCasePO) error
	Delete(ctx context.Context, id uint) error
	CountByAPISpec(ctx context.Context, apiSpecID uint) (int64, error)
}

// ListFilter represents the filter for listing test cases
type ListFilter struct {
	ProjectID *uint
	APISpecID *uint
	Env       *string
	Keyword   *string
	Page      int
	PageSize  int
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new test case repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new test case
func (r *repository) Create(ctx context.Context, tc *TestCasePO) error {
	return r.db.WithContext(ctx).Create(tc).Error
}

// GetByID gets a test case by ID
func (r *repository) GetByID(ctx context.Context, id uint) (*TestCasePO, error) {
	var tc TestCasePO
	err := r.db.WithContext(ctx).First(&tc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tc, nil
}

// List lists test cases with filtering and pagination
func (r *repository) List(ctx context.Context, filter *ListFilter) ([]*TestCasePO, int64, error) {
	query := r.db.WithContext(ctx).Model(&TestCasePO{})

	// Apply filters
	if filter.APISpecID != nil {
		query = query.Where("api_spec_id = ?", *filter.APISpecID)
	}

	if filter.Env != nil && *filter.Env != "" {
		query = query.Where("env = ?", *filter.Env)
	}

	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := "%" + *filter.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	// Get results
	var testCases []*TestCasePO
	err := query.
		Order("id DESC").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&testCases).Error

	if err != nil {
		return nil, 0, err
	}

	return testCases, total, nil
}

// Update updates a test case
func (r *repository) Update(ctx context.Context, tc *TestCasePO) error {
	return r.db.WithContext(ctx).Save(tc).Error
}

// Delete deletes a test case (soft delete)
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&TestCasePO{}, id).Error
}

// CountByAPISpec counts test cases for an API spec
func (r *repository) CountByAPISpec(ctx context.Context, apiSpecID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TestCasePO{}).
		Where("api_spec_id = ?", apiSpecID).
		Count(&count).Error
	return count, err
}
