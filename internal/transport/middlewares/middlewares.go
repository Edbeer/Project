package middlewares

import (
	"context"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
	"github.com/Edbeer/Project/pkg/logger"
	"github.com/google/uuid"
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

// Middleware manager
type MiddlewareManager struct {
	session SessionService
	user    UserService
	config  *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(session SessionService, user UserService, config *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		session: session,
		user:    user,
		config:  config,
		origins: origins,
		logger:  logger,
	}
}
