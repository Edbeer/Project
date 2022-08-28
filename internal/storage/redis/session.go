package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedis.CreateSession")
	defer span.Finish()

	session.RefreshToken = s.manager.NewRefreshToken()

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", errors.Wrap(err, "SessionStorage.CreateSession.Marshal")
	}
	if err := s.redis.Set(ctx, session.RefreshToken, sessionBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "SessionStorage.CreateSession.Set")
	}

	return session.RefreshToken, nil
}

// Get user id from session
func (s *SessionStorage) GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedis.GetUserID")
	defer span.Finish()

	sessionBytes, err := s.redis.Get(ctx, refreshToken).Bytes()
	if err != nil {
		return uuid.Nil , errors.Wrap(err, "SessionStorage.GetUserID.Get")
	}
	session := &entity.Session{}
	if err = json.Unmarshal(sessionBytes, session); err != nil {
		return uuid.Nil, errors.Wrap(err, "SessionStorage.GetSessionByID.Get")
	}

	return session.UserID, nil
}

// Delete session cookie
func (s *SessionStorage) DeleteSession(ctx context.Context, refreshToken string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedis.DeleteSession")
	defer span.Finish()
	if err := s.redis.Del(ctx, refreshToken).Err(); err != nil {
		return errors.Wrap(err, "SessionStorage.DeleteSession.Del")
	}
	return nil
}