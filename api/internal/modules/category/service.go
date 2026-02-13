package category

import (
	"context"
	"fmt"
)

// Service defines the interface for category business logic
type Service interface {
	CreateCategory(ctx context.Context, projectID uint, req *CreateCategoryRequest) (*CategoryResponse, error)
	GetCategory(ctx context.Context, id uint) (*CategoryResponse, error)
	ListCategories(ctx context.Context, projectID uint) ([]*CategoryResponse, error)
	GetCategoryTree(ctx context.Context, projectID uint) ([]*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error
	SortCategories(ctx context.Context, projectID uint, req *SortCategoriesRequest) error
}

type service struct {
	repo Repository
}

// NewService creates a new category service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(ctx context.Context, projectID uint, req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := ToCategoryPO(projectID, req)
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return FromCategoryPO(category), nil
}

func (s *service) GetCategory(ctx context.Context, id uint) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	return FromCategoryPO(category), nil
}

func (s *service) ListCategories(ctx context.Context, projectID uint) ([]*CategoryResponse, error) {
	categories, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return ToResponseList(categories), nil
}

func (s *service) GetCategoryTree(ctx context.Context, projectID uint) ([]*CategoryResponse, error) {
	categories, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories for tree: %w", err)
	}

	resps := ToResponseList(categories)
	return buildTree(resps, nil), nil
}

func buildTree(categories []*CategoryResponse, parentID *uint) []*CategoryResponse {
	var tree []*CategoryResponse
	for _, cat := range categories {
		if (parentID == nil && cat.ParentID == nil) || (parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID) {
			cat.Children = make([]CategoryResponse, 0)
			children := buildTree(categories, &cat.ID)
			for _, child := range children {
				cat.Children = append(cat.Children, *child)
			}
			tree = append(tree, cat)
		}
	}
	return tree
}

func (s *service) UpdateCategory(ctx context.Context, id uint, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category for update: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.ParentID != nil {
		category.ParentID = *req.ParentID
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	return FromCategoryPO(category), nil
}

func (s *service) DeleteCategory(ctx context.Context, id uint) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return fmt.Errorf("category not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) SortCategories(ctx context.Context, projectID uint, req *SortCategoriesRequest) error {
	return s.repo.UpdateSortOrder(ctx, projectID, req.CategoryIDs)
}
