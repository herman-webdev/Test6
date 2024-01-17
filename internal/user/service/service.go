package service

import (
	"awesomeProject/internal/user/dto"
	"awesomeProject/internal/user/model"
	"awesomeProject/pkg/api/sort"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, dto dto.CreateUserDto) error
	GetAll(ctx context.Context, sortOptions sort.Options) ([]model.User, error)
	GetOne(ctx context.Context, uuid string) (model.User, error)
	UpdateOne(ctx context.Context, dto dto.UpdateUserDto, uuid string) error
	DeleteOne(ctx context.Context, uuid string) error
}
