//go:generate mockgen -source interfaces.go -destination mock/service_mock.go -package mock
package service

import (
	"context"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
)

// User service interface
type User interface {
	SignUp(ctx context.Context, input *entity.User) (*entity.UserWithToken, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.UserWithToken, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.UserWithToken, error)
}

// Session service interface
type Session interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}