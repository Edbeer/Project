package api

import (
	"context"
	"net/http"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/Edbeer/Project/pkg/utils"

	"github.com/Edbeer/Project/config"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	SignUp(ctx context.Context, input *entity.InputUser) (*entity.User, error)
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

func (h *UserHandler) SignUp() echo.HandlerFunc {
	type InputUser struct {
		Name     string `json:"name" db:"name" validate:"required_with,lte=30"`
		Email    string `json:"email" db:"email" validate:"omitempty,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		user := &InputUser{}
		if err := utils.ReadRequest(c, user); err != nil {
			return c.JSON(http.StatusNoContent, "Bad")
		}

		createdUser, err := h.user.SignUp(ctx, &entity.InputUser{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		})
		if err != nil {
			return c.JSON(http.StatusNoContent, "Bad")
		}

		return c.JSON(http.StatusOK, createdUser)
	}
}
