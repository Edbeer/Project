package redisrepo

import "github.com/go-redis/redis/v9"

// Storage psql
type Storage struct {
	Session *SessionStorage
}

func NewStorage(redis *redis.Client) *Storage {
	return &Storage{
		Session: newSessionStorage(redis),
	}
}