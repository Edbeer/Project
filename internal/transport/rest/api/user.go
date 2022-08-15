package api

import (
	"net/http"

	"github.com/Edbeer/Project/config"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	Create() error
}

// User handler
type UserHandler struct {
	config *config.Config
	user   UserService
}

// New user handler constructor
func NewUserHandler(config *config.Config, user UserService) *UserHandler {
	return &UserHandler{
		config:      config,
		user: user,
	}
}

func (h *UserHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	}
}
