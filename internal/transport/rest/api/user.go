package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/Edbeer/Project/internal/transport/rest/middlewares"
	"github.com/Edbeer/Project/pkg/httpe"
	"github.com/Edbeer/Project/pkg/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"

	"github.com/Edbeer/Project/config"
	"github.com/labstack/echo/v4"
)

// User service interface
type UserService interface {
	SignUp(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.UserWithToken, error)
}

// Session service interface
type SessionService interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}

// init user handlers
func (h *Handlers) initUserHandlers(api *echo.Group, mw *middlewares.MiddlewareManager) {
	user := api.Group("/user")
	{
		user.POST("/sign-up", h.user.SignUp())
		user.POST("/sign-in", h.user.SignIn())
		user.POST("/auth/refresh", h.user.RefreshTokens())
		user.Use(mw.AuthJWTMiddleware())
		user.POST("/sign-out", h.user.SignOut())
		user.GET("/me", h.user.GetMe())
	}
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

type inputUser struct {
	Name     string `json:"name" validate:"required_with,lte=30"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"required,gte=6"`
}

// SignUp godoc
// @Summary Register new user
// @Description register new user, returns user and access token
// @Tags User
// @Accept json
// @Produce json
// @Param input body inputUser true "sign up info"
// @Success 201 {object} entity.User
// @Router /user/sign-up [post]
func (h *UserHandler) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "UserHandler.SignUp")
		defer span.Finish()

		user := &inputUser{}
		if err := utils.ReadRequest(c, user); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		createdUser, err := h.user.SignUp(ctx, &entity.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		})
		if err != nil {
			return c.JSON(httpe.ParseErrors(err).Status(), httpe.ParseErrors(err))
		}

		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: createdUser.User.ID,
		}, h.config.Cookie.MaxAge)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

type Login struct {
	Email    string `json:"email" validate:"omitempty,lte=60,email"`
	Password string `json:"password,omitempty" validate:"required,gte=6"`
}

// SignIn godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags User
// @Accept json
// @Produce json
// @Param input body Login true "sign up info"
// @Success 200 {object} entity.User
// @Router /user/sign-in [post]
func (h *UserHandler) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "UserHandler.SignIn")
		defer span.Finish()

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		userWithToken, err := h.user.SignIn(ctx, &entity.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: userWithToken.User.ID,
		}, h.config.Cookie.MaxAge)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))
		return c.JSON(http.StatusOK, userWithToken)
	}
}

type RefreshToken struct {
	Token string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary Refresh Tokens
// @Tags User
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body RefreshToken true "sign up info"
// @Success 200 {object} TokenResponse
// @Router /user/auth/refresh [post]
func (h *UserHandler) RefreshTokens() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "UserHandler.RefreshTokens")
		defer span.Finish()

		token := &RefreshToken{}
		if err := utils.ReadRequest(c, token); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		uuid, err := h.session.GetUserID(ctx, token.Token)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		
		user, err := h.user.GetUserByID(ctx, uuid)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		refreshToken, err := h.session.CreateSession(ctx, &entity.Session{
			UserID: user.User.ID,
		}, h.config.Cookie.MaxAge)
		if err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}

		c.SetCookie(utils.ConfigureJWTCookie(h.config, refreshToken))
		return c.JSON(http.StatusOK, TokenResponse{
			AccessToken: user.AccessToken,
			RefreshToken: refreshToken,
		})
	}
}

// SignOut godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /user/sign-out [post]
func (u *UserHandler) SignOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "UserHandler.SignOut")
		defer span.Finish()

		cookie, err := c.Cookie("jwt-token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusInternalServerError, httpe.NewInternalServerError(err))
		}
		if err = u.session.DeleteSession(ctx, cookie.Value); err != nil {
			return c.JSON(httpe.ErrorResponse(err))
		}
		utils.DeleteCookie(c, u.config.Cookie.Name)

		return c.NoContent(http.StatusOK)
	}
}

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} entity.User
// @Failure 500 {object} httpe.RestError
// @Router /user/me [get]
func (u *UserHandler) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*entity.User)
		if !ok {
			httpe.NewUnauthorizedError(httpe.Unauthorized)
		}

		return c.JSON(http.StatusOK, user)
	}
}