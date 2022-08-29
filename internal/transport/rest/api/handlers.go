package api

import (
	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/transport/rest/middlewares"
	"github.com/Edbeer/Project/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/Edbeer/Project/docs"
)

// Dependencies
type Deps struct {
	UserService    UserService
	SessionService SessionService
	Config         *config.Config
}

// Handlers
type Handlers struct {
	user *UserHandler
}

// New handlers constructor
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		user: NewUserHandler(deps.Config, deps.UserService, deps.SessionService),
	}
}

func (h *Handlers) Init(e *echo.Echo, logger logger.Logger) error {
	config := config.GetConfig()
	if config.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Use(middleware.RequestID())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	// Request ID middleware generates a unique id for a request.
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	// Middleware Manager
	mw := middlewares.NewMiddlewareManager(
		h.user.session,
		h.user.user,
		h.user.config,
		[]string{"*"},
		logger,
	)
	docs.SwaggerInfo.Title = "Auth JWT example restapi"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	h.initApi(e, mw)

	return nil
}

func (h *Handlers) initApi(e *echo.Echo, mw *middlewares.MiddlewareManager) {
	api := e.Group("/api")
	{
		h.initUserHandlers(api, mw)
	}
}
