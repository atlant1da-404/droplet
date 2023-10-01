package http

import (
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/errs"
	"github.com/gin-gonic/gin"
)

type nodeRouter struct {
	RouterContext
}

func setupNodeRoutes(options RouterOptions) {
	router := &authRouter{
		RouterContext{
			logger:   options.Logger,
			services: options.Services,
			config:   options.Config,
		},
	}

	routerGroup := options.Handler.Group("/node")
	{
		routerGroup.POST("", wrapHandler(options, router.signIn))
	}
}

type createNodeRequestBody struct {
	*service.CreateNodeOptions
} // @name createNodeRequestBody

type createNodeResponseBody struct {
	*service.CreateNodeOutput
} // @name createNodeResponseBody

type createNodeResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:""`
} // @name createNodeResponseError

func (e createNodeResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           CreateNode
// @Summary      Create node.
// @Accept       application/json
// @Produce      application/json
// @Param        fields body createNodeRequestBody true "data"
// @Success      200 {object} createNodeResponseError
// @Failure      422,500 {object} createNodeResponseError
// @Router       /node [POST]
func (a *nodeRouter) createNode(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("signIn").WithContext(requestContext)

	var body createNodeRequestBody
	err := requestContext.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	nodeID, err := a.services.NodeService.CreateNode(requestContext, body.CreateNodeOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, createNodeResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to create account", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create node", Details: err}
	}

	logger.Info("successfully created a node")
	return &createNodeResponseBody{nodeID}, nil
}
