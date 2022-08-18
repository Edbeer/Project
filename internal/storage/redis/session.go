package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/pkg/errors"

	"github.com/go-redis/redis/v9"
)

type Manager interface {
	NewRefreshToken() string
}

// Session redis storage
type SessionStorage struct {
	redis   *redis.Client
	manager Manager
}

// Session storage constructor
func newSessionStorage(redis *redis.Client, manager Manager) *SessionStorage {
	return &SessionStorage{
		redis:   redis,
		manager: manager,
	}
}

// Add refresh token in redis
func (s *SessionStorage) CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error) {

	session.RefreshToken = s.manager.NewRefreshToken()
	key := s.createKey(session.RefreshToken)

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", errors.Wrap(err, "SessionStorage.CreateSession.Marshal")
	}
	if err := s.redis.Set(ctx, key, sessionBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "SessionStorage.CreateSession.Set")
	}

	return session.RefreshToken, nil
}

func (s *SessionStorage) createKey(refreshToken string) string {
	return fmt.Sprintf("session: %s", refreshToken)
}