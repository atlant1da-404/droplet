package http

import (
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/errs"
	"github.com/gin-gonic/gin"
)

type authRouter struct {
	RouterContext
}

func setupAuthRoutes(options RouterOptions) {
	a := &authRouter{
		RouterContext{
			logger:   options.Logger,
			services: options.Services,
			config:   options.Config,
		},
	}

	g := options.Handler.Group("/auth")
	{
		g.POST("/sign-in", wrapHandler(options, a.signIn))
		g.POST("/sign-up", wrapHandler(options, a.signUp))
	}
}

type signInRequestBody struct {
	*service.SignInOptions
} // @name signInRequestBody

type signInResponseBody struct {
	*service.SignInOutput
} // @name signInResponseBody

type signInResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found,wrong_password"`
} // @name signInResponseError

func (e signInResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           SignIn
// @Summary      Login user.
// @Accept       application/json
// @Produce      application/json
// @Param        fields body signInRequestBody true "data"
// @Success      200 {object} signInResponseBody
// @Failure      422,500 {object} signInResponseError
// @Router       /sign-up [POST]
func (a *authRouter) signIn(c *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("signIn").WithContext(c)

	var body signInRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	signed, err := a.services.AuthService.SignIn(c, body.SignInOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, signInResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to sign in", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to sign in", Details: err}
	}

	logger.Info("successfully signed in")
	return &signInResponseBody{signed}, nil
}

type signUpRequestBody struct {
	*service.SignUpOptions
} // @name signUpRequestBody

type signUpResponseBody struct {
	*service.SignUpOutput
} // @name signUpResponseBody

type signUpResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_already_created"`
} // @name signUpResponseError

func (e signUpResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           SignUp
// @Summary      Creates and returns user.
// @Accept       application/json
// @Produce      application/json
// @Param        fields body signUpRequestBody true "data"
// @Success      200 {object} signUpResponseBody
// @Failure      422,500 {object} signUpResponseError
// @Router       /sign-up [POST]
func (a *authRouter) signUp(c *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("signUp").WithContext(c)

	var body signUpRequestBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	createdUser, err := a.services.AuthService.SignUp(c, body.SignUpOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, signUpResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to create and return user", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create and return user", Details: err}
	}
	logger = logger.With("createdUser", createdUser)

	logger.Info("user created and returned")
	return &signUpResponseBody{createdUser}, nil
}