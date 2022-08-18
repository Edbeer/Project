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
	SignUp(ctx context.Context, input *entity.InputUser) (*entity.UserWithToken, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
}

// User handler
type UserHandler struct {
	config  *config.Config
	user    UserService
	session SessionService
}

// New user handler constructor
func NewUserHandler(config *config.Config, user UserService, session SessionService) *UserHandler {
	return &UserHandler{
		config:  config,
		user:    user,
		session: session,
	}
}

func (h *UserHandler) SignUp() echo.HandlerFunc {
	type inputUser struct {
		Name     string `json:"name" db:"name" validate:"required_with,lte=30"`
		Email    string `json:"email" db:"email" validate:"omitempty,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		user := &inputUser{}
		if err := utils.ReadRequest(c, user); err != nil {
			return c.JSON(http.StatusBadRequest, "400")
		}

		createdUser, err := h.user.SignUp(ctx, &entity.InputUser{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		})
		if err != nil {
			return c.JSON(http.StatusNoContent, "204")
		}

		// TODO 
		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: createdUser.User.ID,
		}, h.config.Session.Expire)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "500")
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (h *UserHandler) SignIn() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			return c.JSON(http.StatusBadRequest, "400")
		}
		userWithToken, err := h.user.SignIn(ctx, &entity.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			return c.JSON(http.StatusNoContent, "204")
		}

		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: userWithToken.User.ID,
		}, h.config.Cookie.MaxAge)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "500")
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))
		return c.JSON(http.StatusOK, userWithToken)
	}
}
