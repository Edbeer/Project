package service

import "github.com/Edbeer/Project/config"

type UserPsql interface {
	Create() error
}

// PasswordHasher provides hashing logic to securely store passwords
type PasswordHasher interface {
	Hash(password string) string
}

// User service
type UserService struct {
	config *config.Config
	psql   UserPsql
	hash PasswordHasher
}

// New user service constructor
func NewUserService(config *config.Config, psql UserPsql, hash PasswordHasher) *UserService {
	return &UserService{
		config: config,
		psql:   psql,
		hash: hash,
	}
}

func (u *UserService) Create() error {
	return nil
}