package service

import (
	"context"
	"github.com/atlant1da-404/droplet/internal/entity"
)

type Storages struct {
	UserStorage UserStorage
}

type UserStorage interface {
	// GetUser provides getting user from storage via requested filters.
	GetUser(ctx context.Context, filter *GetUserFilter) (*entity.User, error)
	// CreateUser provides creating user in the system.
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type GetUserFilter struct {
	Email string
}
