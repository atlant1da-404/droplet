package http

import (
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type accountRouter struct {
	RouterContext
}

func setupAccountRoutes(options RouterOptions) {
	router := &accountRouter{
		RouterContext{
			logger:   options.Logger,
			services: options.Services,
			config:   options.Config,
		},
	}

	routerGroup := options.Handler.Group("/account")
	{
		routerGroup.POST("", wrapHandler(options, router.createAccount))
		routerGroup.GET("/:id", wrapHandler(options, router.getAccount))
		routerGroup.PATCH("/:id", wrapHandler(options, router.updateAccount))
	}
}

type createAccountRequestBody struct {
	*service.CreateAccountOptions
} // @name createAccountRequestBody

type createAccountResponseBody struct {
	*service.CreateAccountOutput
} // @name createAccountResponseBody

type createAccountResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name createAccountResponseError

func (e createAccountResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           CreateAccount
// @Summary      Creates account.
// @Accept       application/json
// @Produce      application/json
// @Param        fields body createAccountRequestBody true "data"
// @Success      200 {object} createAccountResponseBody
// @Failure      422,500 {object} createAccountResponseError
// @Router       /account [POST]
func (a *accountRouter) createAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("createAccount").WithContext(requestContext)

	var body createAccountRequestBody
	err := requestContext.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	createdAccount, err := a.services.AccountService.CreateAccount(requestContext, body.CreateAccountOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, createAccountResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to create account", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create account", Details: err}
	}
	logger = logger.With("createdAccount", createdAccount)

	logger.Info("account created successfully")
	return &createAccountResponseBody{createdAccount}, nil
}

type getAccountResponseBody struct {
	*entity.Account
} // @name getAccountResponseBody

type getAccountResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name getAccountResponseError

func (e getAccountResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           GetAccount
// @Summary      Gets account.
// @Accept       application/json
// @Produce      application/json
// @Param        id path string true "Account ID"
// @Success      200 {object} getAccountResponseBody
// @Failure      422,500 {object} getAccountResponseError
// @Router       /account/{id} [PUT]
func (a *accountRouter) getAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("getAccount").WithContext(requestContext)

	accountId := requestContext.Param("id")
	if _, ok := uuid.Parse(accountId); ok != nil {
		logger.Info("invalid account id parameter", "param", accountId)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid account id parameter"}
	}
	logger = logger.With("accountId", accountId)
	logger.Debug("parsed params")

	account, err := a.services.AccountService.GetAccount(requestContext, &service.GetAccountOptions{AccountId: accountId})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, getAccountResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
	}
	logger = logger.With("account", account)

	logger.Info("successfully got account")
	return getAccountResponseBody{account}, nil
}

type updateAccountRequestBody struct {
	*entity.Account
} // @name updateAccountRequestBody

type updateAccountResponseBody struct {
	*entity.Account
} // @name updateAccountResponseBody

type updateAccountResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name updateAccountResponseError

func (e updateAccountResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

func (a *accountRouter) updateAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("updateAccount").WithContext(requestContext)

	logger.Info("successfully updated account")
	return nil, nil
}
