package storage

import (
	"context"
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accountStorage struct {
	*database.PostgreSQL
}

var _ service.AccountStorage = (*accountStorage)(nil)

func NewAccountStorage(postgresql *database.PostgreSQL) service.AccountStorage {
	return &accountStorage{postgresql}
}

func (u *accountStorage) CreateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	err := u.DB.WithContext(ctx).Create(account).Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (u *accountStorage) GetAccount(ctx context.Context, filter *service.GetAccountFilter) (*entity.Account, error) {
	stmt := u.DB.Preload(clause.Associations)

	if filter.AccountId != "" {
		stmt = stmt.Where(entity.Account{Id: filter.AccountId})
	}

	var account entity.Account
	err := stmt.
		WithContext(ctx).
		First(&account).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}
