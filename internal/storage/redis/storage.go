package redisrepo

import (
	"github.com/go-redis/redis/v9"
)

type Deps struct {
	Redis *redis.Client
	Manager Manager
}

// Storage redis
type Storage struct {
	Session *SessionStorage
}

func NewStorage(deps Deps) *Storage {
	return &Storage{
		Session: newSessionStorage(deps.Redis, deps.Manager),
	}
}
