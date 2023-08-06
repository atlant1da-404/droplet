package http

import (
	"bytes"
	"fmt"
	"github.com/DataDog/gostackparse"
	"github.com/atlant1da-404/droplet/config"
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/logger"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"runtime/debug"
)

type Options struct {
	Handler  *gin.Engine
	Logger   logger.Logger
	Services service.Services
	Config   *config.Config
}

type RouterOptions struct {
	Handler  *gin.RouterGroup
	Logger   logger.Logger
	Services service.Services
	Config   *config.Config
}

type RouterContext struct {
	logger   logger.Logger
	services service.Services
	config   *config.Config
}

func New(options *Options) {
	options.Handler.Use(
		ginzap.RecoveryWithZap(options.Logger.Named("HTTPController").Unwrap(), true),
		requestIDMiddleware,
		corsMiddleware,
	)

	routerOptions := RouterOptions{
		Handler:  options.Handler.Group("/api/v1"),
		Services: options.Services,
		Logger:   options.Logger.Named("HTTPController"),
		Config:   options.Config,
	}

	// routes
	{
		setupAuthRoutes(routerOptions)
	}
}

// requestIDMiddleware is used to add request id to gin context.
func requestIDMiddleware(c *gin.Context) {
	c.Set("RequestID", uuid.NewString())
}

// corsMiddleware - used to allow incoming cross-origin requests.
func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

// httpResponseError provides a base error type for all errors.
type httpResponseError struct {
	Type          httpErrType `json:"-"`
	Message       string      `json:"message"`
	Code          string      `json:"code,omitempty"`
	Details       interface{} `json:"details,omitempty"`
	InvalidFields interface{} `json:"invalidFields,omitempty"`
}

// httpErrType is used to define error type.
type httpErrType string

// Error is used to convert an error to a string.
func (e httpResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer httpErrType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient httpErrType = "client"
)

// wrapHandler provides unified error handling for all handlers.
func wrapHandler(options RouterOptions, handler func(c *gin.Context) (interface{}, *httpResponseError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		lgr := options.Logger.Named("wrapHandler")
		body, err := handler(c)

		// handle panics
		defer func() {
			if err := recover(); err != nil {
				// get stacktrace
				stacktrace, errors := gostackparse.Parse(bytes.NewReader(debug.Stack()))
				if len(errors) > 0 || len(stacktrace) == 0 {
					lgr.Error("get stacktrace errors", "stacktraceErrors", errors, "stacktrace", "unknown", "err", err)
				} else {
					lgr.Error("unhandled error", "err", err, "stacktrace", stacktrace)
				}

				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					lgr.Error("failed to abort with error", "err", err)
				}
			}
		}()

		// check if middleware
		if body == nil && err == nil {
			return
		}
		lgr = lgr.With("body", body).With("err", err)

		if err != nil {
			if err.Type == ErrorTypeServer {
				lgr.Error("internal server error")

				err := c.AbortWithError(http.StatusInternalServerError, err)
				if err != nil {
					lgr.Error("failed to abort with error", "err", err)
				}
				lgr.Info("aborted with error")

			} else {
				lgr.Info("client error")
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			}
			return
		}
		lgr.Info("request handled")
		c.JSON(http.StatusOK, body)
	}
}
