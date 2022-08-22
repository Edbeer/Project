package api

import (
	"context"
	"net/http"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/Edbeer/Project/pkg/utils"
	"github.com/google/uuid"

	"github.com/Edbeer/Project/config"
	"github.com/labstack/echo/v4"
)

// User service interface
type UserService interface {
	SignUp(ctx context.Context, input *entity.InputUser) (*entity.UserWithToken, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.UserWithToken, error)
}

// Session service interface
type SessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
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

// Update tokens
func (h *UserHandler) RefreshTokens() echo.HandlerFunc {
	type RefreshToken struct {
		Token string `json:"refresh_token" redis:"refresh_token" binding:"required"`
	}
	type tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		token := &RefreshToken{}
		if err := utils.ReadRequest(c, token); err != nil {
			return c.JSON(http.StatusBadRequest, "400")
		}
		uuid, err := h.session.GetUserID(ctx, token.Token)
		if err != nil {
			return c.JSON(http.StatusNoContent, "uid")
		}
		
		user, err := h.user.GetUserByID(ctx, uuid)
		if err != nil {
			return c.JSON(http.StatusNoContent, "204")
		}

		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: user.User.ID,
		}, h.config.Cookie.MaxAge)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "500")
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))
		return c.JSON(http.StatusOK, tokenResponse{
			AccessToken: user.AccessToken,
			RefreshToken: refreshToken,
		})
	}
}
