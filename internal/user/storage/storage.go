package storage

import (
	"awesomeProject/internal/user/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	FindAll(ctx context.Context, sortOptions SortOptions) ([]model.User, error)
	FindOne(ctx context.Context, id string) (model.User, error)
	Update(ctx context.Context, user *model.User, id string) error
	Delete(ctx context.Context, id string) error
}

type SortOptions interface {
	GetOrderBy() string
}
