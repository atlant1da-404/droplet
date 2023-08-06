package service

import (
	"context"
	"fmt"
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/pkg/auth"
	"github.com/atlant1da-404/droplet/pkg/hash"
)

type authService struct {
	serviceContext
	hash hash.Hash
	auth auth.Authenticator
}

var _ AuthService = (*authService)(nil)

func NewAuthService(options *Options) AuthService {
	return &authService{
		serviceContext: serviceContext{
			storages: options.Storages,
			config:   options.Config,
			logger:   options.Logger.Named("AuthService"),
		},
		hash: options.Hash,
		auth: options.Auth,
	}
}

func (a authService) SignIn(ctx context.Context, opt *SignInOptions) (*SignInOutput, error) {
	logger := a.logger.
		Named("SignIn").
		WithContext(ctx).
		With("body", opt)

	user, err := a.storages.UserStorage.GetUser(ctx, &GetUserFilter{Email: opt.Email})
	if err != nil {
		logger.Error("failed to get user: ", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user already created")
		return nil, ErrSignInUserNotFound
	}
	logger = logger.With("user", user)

	err = a.hash.CompareHash([]byte(user.Password), []byte(opt.Password))
	if err != nil {
		logger.Info(err.Error())
		return nil, ErrSignInWrongPassword
	}

	accessToken, err := a.auth.GenerateToken(user.Username)
	if err != nil {
		logger.Error("failed to generate token for user: ", err)
		return nil, fmt.Errorf("failed to generate token for user: %w", err)
	}

	logger.Info("successfully signed user")
	return &SignInOutput{AccessToken: accessToken}, nil
}

func (a authService) SignUp(ctx context.Context, opt *SignUpOptions) (*SignUpOutput, error) {
	logger := a.logger.
		Named("SignUp").
		WithContext(ctx).
		With("body", opt)

	user, err := a.storages.UserStorage.GetUser(ctx, &GetUserFilter{Email: opt.Email})
	if err != nil {
		logger.Error("failed to get user: ", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		logger.Info("user already created")
		return nil, ErrSignUpUserAlreadyCreated
	}

	hashedPassword, err := a.hash.GenerateHash(opt.Password)
	if err != nil {
		logger.Error("failed to hash user password: ", err)
		return nil, fmt.Errorf("failed to hash user: %w", err)
	}

	createdUser, err := a.storages.UserStorage.CreateUser(ctx, &entity.User{Email: opt.Email, Password: hashedPassword, Username: opt.Username, MacAddress: opt.MacAddress})
	if err != nil {
		logger.Error("failed to create user: ", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := a.auth.GenerateToken(createdUser.Username)
	if err != nil {
		logger.Error("failed to generate token for user: ", err)
		return nil, fmt.Errorf("failed to generate token for user: %w", err)
	}

	logger.Info("successfully handled sign up")
	return &SignUpOutput{Id: createdUser.Id, Username: createdUser.Username, Email: createdUser.Email, AccessToken: accessToken}, nil
}
