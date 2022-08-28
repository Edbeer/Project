//go:generate mockgen -source interfaces.go -destination mock/redis_storage_mock.go -package mock
package redisrepo

import (
	"context"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
)

// Session storage interface
type SessionRedis interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}