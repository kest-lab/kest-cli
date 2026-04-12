package category

import (
	"context"
	"errors"
	"testing"
)

func TestServiceGetCategoryRespectsProjectScope(t *testing.T) {
	repo := &stubCategoryRepository{
		categories: map[uint]*CategoryPO{
			7: {
				ID:        7,
				ProjectID: 2,
				Name:      "Foreign",
			},
		},
	}

	service := NewService(repo)
	_, err := service.GetCategory(context.Background(), 1, 7)
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestServiceUpdateCategoryRejectsCrossProjectParent(t *testing.T) {
	projectID := uint(1)
	categoryID := uint(10)
	foreignParentID := uint(99)

	repo := &stubCategoryRepository{
		categories: map[uint]*CategoryPO{
			categoryID: {
				ID:        categoryID,
				ProjectID: projectID,
				Name:      "Orders",
			},
			foreignParentID: {
				ID:        foreignParentID,
				ProjectID: 2,
				Name:      "Foreign Parent",
			},
		},
	}

	service := NewService(repo)
	req := &UpdateCategoryRequest{
		ParentID: ptrToPtr(foreignParentID),
	}

	_, err := service.UpdateCategory(context.Background(), projectID, categoryID, req)
	if !errors.Is(err, ErrInvalidParentCategory) {
		t.Fatalf("expected ErrInvalidParentCategory, got %v", err)
	}
}

func TestServiceDeleteCategoryRespectsProjectScope(t *testing.T) {
	repo := &stubCategoryRepository{
		categories: map[uint]*CategoryPO{
			8: {
				ID:        8,
				ProjectID: 2,
				Name:      "Foreign",
			},
		},
	}

	service := NewService(repo)
	err := service.DeleteCategory(context.Background(), 1, 8)
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

type stubCategoryRepository struct {
	categories map[uint]*CategoryPO
}

func (r *stubCategoryRepository) Create(_ context.Context, category *CategoryPO) error {
	r.categories[category.ID] = cloneCategory(category)
	return nil
}

func (r *stubCategoryRepository) GetByID(_ context.Context, id uint) (*CategoryPO, error) {
	return cloneCategory(r.categories[id]), nil
}

func (r *stubCategoryRepository) GetByIDAndProject(_ context.Context, id, projectID uint) (*CategoryPO, error) {
	category := r.categories[id]
	if category == nil || category.ProjectID != projectID {
		return nil, nil
	}

	return cloneCategory(category), nil
}

func (r *stubCategoryRepository) ListByProject(_ context.Context, projectID uint) ([]*CategoryPO, error) {
	var categories []*CategoryPO
	for _, category := range r.categories {
		if category.ProjectID == projectID {
			categories = append(categories, cloneCategory(category))
		}
	}
	return categories, nil
}

func (r *stubCategoryRepository) Update(_ context.Context, category *CategoryPO) error {
	r.categories[category.ID] = cloneCategory(category)
	return nil
}

func (r *stubCategoryRepository) Delete(_ context.Context, id uint) error {
	delete(r.categories, id)
	return nil
}

func (r *stubCategoryRepository) UpdateSortOrder(_ context.Context, projectID uint, categoryIDs []uint) error {
	for i, id := range categoryIDs {
		category := r.categories[id]
		if category == nil || category.ProjectID != projectID {
			continue
		}
		category.SortOrder = i
	}
	return nil
}

func cloneCategory(category *CategoryPO) *CategoryPO {
	if category == nil {
		return nil
	}

	cloned := *category
	if category.ParentID != nil {
		parentID := *category.ParentID
		cloned.ParentID = &parentID
	}

	return &cloned
}

func ptrToPtr(value uint) **uint {
	ptr := &value
	return &ptr
}
