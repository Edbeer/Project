package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Edbeer/Project/config"
)
// Dependencies
type Deps struct {
	UserService UserService
	Config      *config.Config
}

// Handlers
type Handlers struct {
	user *UserHandler
}

// New handlers constructor
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		user: NewUserHandler(deps.Config, deps.UserService),
	}
}

func (h *Handlers) Init(e *echo.Echo) error {
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

	h.initApi(e)

	return nil
}

func (h *Handlers) initApi(e *echo.Echo) {
	api := e.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/create", h.user.SignUp())
		}
	}
}