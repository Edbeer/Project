package service

import (
	"github.com/Edbeer/Project/internal/storage/psql"
	"github.com/Edbeer/Project/config"
)

// Services
type Services struct {
	User *UserService
}

// Dependencies
type Deps struct {
	Config       *config.Config
	PsqlStorage  *psql.Storage
	Hash         PasswordHasher
}

// New services constructor
func NewServices(deps Deps) *Services {
	userService := NewUserService(deps.Config, deps.PsqlStorage.User, deps.Hash)
	return &Services{
		User: userService,
	}
}
