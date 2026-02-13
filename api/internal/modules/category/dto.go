package category

import "time"

// ========== Request DTOs ==========

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	ParentID    *uint  `json:"parent_id"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	ParentID    **uint  `json:"parent_id"` // Double pointer to allow setting to nil
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

type SortCategoriesRequest struct {
	CategoryIDs []uint `json:"category_ids" binding:"required"`
}

// ========== Response DTOs ==========

type CategoryResponse struct {
	ID          uint               `json:"id"`
	ProjectID   uint               `json:"project_id"`
	Name        string             `json:"name"`
	ParentID    *uint              `json:"parent_id"`
	Description string             `json:"description,omitempty"`
	SortOrder   int                `json:"sort_order"`
	Children    []CategoryResponse `json:"children,omitempty"` // For tree structure if needed
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToCategoryPO converts CreateCategoryRequest to CategoryPO
func ToCategoryPO(projectID uint, req *CreateCategoryRequest) *CategoryPO {
	return &CategoryPO{
		ProjectID:   projectID,
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
}

// FromCategoryPO converts CategoryPO to CategoryResponse
func FromCategoryPO(po *CategoryPO) *CategoryResponse {
	if po == nil {
		return nil
	}
	return &CategoryResponse{
		ID:          po.ID,
		ProjectID:   po.ProjectID,
		Name:        po.Name,
		ParentID:    po.ParentID,
		Description: po.Description,
		SortOrder:   po.SortOrder,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

// ToResponseList converts a slice of CategoryPO to a slice of CategoryResponse
func ToResponseList(pos []*CategoryPO) []*CategoryResponse {
	resps := make([]*CategoryResponse, 0, len(pos))
	for _, po := range pos {
		resps = append(resps, FromCategoryPO(po))
	}
	return resps
}
