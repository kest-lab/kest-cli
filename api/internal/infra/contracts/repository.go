package contracts

import "context"

// Repository defines the base repository interface
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	Find(ctx context.Context, id uint) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
	FindBy(ctx context.Context, field string, value any) (*T, error)
	FindAllBy(ctx context.Context, field string, value any) ([]*T, error)
}

// PaginatedRepository adds pagination support
type PaginatedRepository[T any] interface {
	Repository[T]
	Paginate(ctx context.Context, page, perPage int) ([]*T, int64, error)
}
