package service

import (
	"context"
	"github.com/atlant1da-404/droplet/config"
	"github.com/atlant1da-404/droplet/pkg/auth"
	"github.com/atlant1da-404/droplet/pkg/errs"
	"github.com/atlant1da-404/droplet/pkg/hash"
	"github.com/atlant1da-404/droplet/pkg/logger"
)

type Services struct {
	AuthService AuthService
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
	SignIn(ctx context.Context, opt *SignInOptions) (*SignInOutput, error)
	// SignUp provides logic of creating the clients and returns access and refresh tokens.
	SignUp(ctx context.Context, opt *SignUpOptions) (*SignUpOutput, error)
}

type SignInOptions struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken string
}

type SignUpOptions struct {
	Username string
	Email    string
	Password string
}

type SignUpOutput struct {
	Id          string
	Username    string
	Email       string
	AccessToken string
}

var (
	ErrSignUpUserAlreadyCreated = errs.New("user already created", "user_already_created")
	ErrSignInUserNotFound       = errs.New("user not found", "user_not_found")
	ErrSignInWrongPassword      = errs.New("wrong password", "wrong_password")
)
