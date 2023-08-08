package service

import (
	"context"
	"fmt"
	"github.com/atlant1da-404/droplet/internal/entity"
)

type accountService struct {
	serviceContext
}

var _ AccountService = (*accountService)(nil)

func NewAccountService(options *Options) AccountService {
	return &accountService{
		serviceContext: serviceContext{
			storages: options.Storages,
			config:   options.Config,
			logger:   options.Logger.Named("AuthService"),
		},
	}
}

func (a accountService) CreateAccount(ctx context.Context, options *CreateAccountOptions) (*CreateAccountOutput, error) {
	logger := a.logger.
		Named("CreateAccount").
		WithContext(ctx).
		With("options", options)

	user, err := a.storages.UserStorage.GetUser(ctx, &GetUserFilter{UserId: options.UserId})
	if err != nil {
		logger.Error("failed to get user: ", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return nil, ErrCreateAccountUserNotFound
	}
	logger = logger.With("user", user)

	account := &entity.Account{
		UserId: user.Id,
		AccountDevices: []entity.AccountDevices{
			{
				Name:       options.DeviceName,
				OS:         options.DeviceOS,
				MacAddress: options.DeviceMacAddress,
				Active:     options.Active,
			},
		},
		AccountSettings: &entity.AccountSettings{
			Language: options.AccountLanguage,
		},
	}
	logger = logger.With("account", account)

	createdAccount, err := a.storages.AccountStorage.CreateAccount(ctx, account)
	if err != nil {
		logger.Error("failed to create account: %w", err)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	logger = logger.With("createdAccount", createdAccount)

	logger.Info("account successfully created")
	return &CreateAccountOutput{Id: createdAccount.Id, UserId: createdAccount.UserId}, nil
}

func (a accountService) GetAccount(ctx context.Context, options *GetAccountOptions) (*entity.Account, error) {
	logger := a.logger.
		Named("GetAccount").
		WithContext(ctx).
		With("options", options)

	account, err := a.storages.AccountStorage.GetAccount(ctx, &GetAccountFilter{AccountId: options.AccountId, UserId: options.UserId})
	if err != nil {
		logger.Error("failed to get account: ", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		logger.Info("account not found")
		return nil, ErrGetAccountAccountNotFound
	}
	logger = logger.With("account", account)

	logger.Info("successfully got account")
	return account, nil
}
