package redisrepo

import "github.com/go-redis/redis/v9"

// Session redis storage
type SessionStorage struct {
	redis *redis.Client
}

// Session storage constructor
func newSessionStorage(redis *redis.Client) *SessionStorage {
	return &SessionStorage{redis: redis}
}

// Create session
func (s *SessionStorage) CreateSession() error {
	return nil
}