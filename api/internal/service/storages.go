package service

import (
	"context"
	"github.com/atlant1da-404/droplet/internal/entity"
)

type Storages struct {
	UserStorage    UserStorage
	AccountStorage AccountStorage
}

type UserStorage interface {
	// GetUser provides getting user from storage via requested filters.
	GetUser(ctx context.Context, filter *GetUserFilter) (*entity.User, error)
	// CreateUser provides creating user in the system.
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type GetUserFilter struct {
	Email  string
	UserId string
}

type AccountStorage interface {
	// CreateAccount provides creating account in the system.
	CreateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error)
	// GetAccount provides logic of getting account from storage.
	GetAccount(ctx context.Context, filter *GetAccountFilter) (*entity.Account, error)
}

type GetAccountFilter struct {
	AccountId string
	UserId    string
}
