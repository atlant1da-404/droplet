package service

import (
	"context"
	"github.com/atlant1da-404/droplet/config"
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/pkg/auth"
	"github.com/atlant1da-404/droplet/pkg/errs"
	"github.com/atlant1da-404/droplet/pkg/hash"
	"github.com/atlant1da-404/droplet/pkg/logger"
)

type Services struct {
	AuthService    AuthService
	AccountService AccountService
}

type Options struct {
	Storages *Storages
	Config   *config.Config
	Logger   logger.Logger
	Hash     hash.Hash
	Auth     auth.Authenticator
}

type serviceContext struct {
	storages *Storages
	config   *config.Config
	logger   logger.Logger
}

type AuthService interface {
	// SignIn provides logic of authentication of clients and returns access and refresh tokens.
	SignIn(ctx context.Context, options *SignInOptions) (*SignInOutput, error)
	// SignUp provides logic of creating the clients and returns access and refresh tokens.
	SignUp(ctx context.Context, options *SignUpOptions) (*SignUpOutput, error)
	// VerifyToken provides logic of validating provided authorization token.
	VerifyToken(ctx context.Context, options *VerifyTokenOptions) (*VerifyTokenOutput, error)
}

type SignInOptions struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken string
}

type SignUpOptions struct {
	Username   string
	Email      string
	Password   string
	MacAddress string
}

type SignUpOutput struct {
	Id          string
	Username    string
	Email       string
	AccessToken string
}

type VerifyTokenOptions struct {
	AccessToken string
}

type VerifyTokenOutput struct {
	Username string
	UserId   string
}

var (
	ErrSignUpUserAlreadyCreated = errs.New("user already created", "user_already_created")
	ErrSignInUserNotFound       = errs.New("user not found", "user_not_found")
	ErrSignInWrongPassword      = errs.New("wrong password", "wrong_password")
)

type AccountService interface {
	// CreateAccount provides logic of creating account for clients.
	CreateAccount(ctx context.Context, options *CreateAccountOptions) (*CreateAccountOutput, error)
	// GetAccount provides logic of getting account via accountId.
	GetAccount(ctx context.Context, options *GetAccountOptions) (*entity.Account, error)
}

type CreateAccountOptions struct {
	UserId           string `json:"userId"`
	DeviceName       string `json:"deviceName"`
	DeviceOS         string `json:"deviceOs"`
	DeviceMacAddress string `json:"deviceMacAddress"`
	Active           bool   `json:"active"`
	AccountLanguage  string `json:"accountLanguage"`
}

type CreateAccountOutput struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type GetAccountOptions struct {
	AccountId string `json:"accountId"`
	UserId    string `json:"userId"`
}

var (
	ErrCreateAccountUserNotFound = errs.New("user not found", "user_not_found")
	ErrGetAccountAccountNotFound = errs.New("account not found", "account_not_found")
)
