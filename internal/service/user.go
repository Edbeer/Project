package service

import (
	"github.com/Edbeer/Project/pkg/utils"
	"context"

	"github.com/Edbeer/Project/internal/entity"

	"github.com/Edbeer/Project/config"
)

// User psql storage interface
type UserPsql interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
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

func (u *UserService) SignUp(ctx context.Context, input *entity.InputUser) (*entity.User, error) {	
	user := &entity.User{
		Name: input.Name,
		Password: u.hash.Hash(input.Password),
		Email: input.Email,
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, err
	}

	existsUser, err := u.psql.FindUserByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, err
	}

	if err := utils.ValidateStruct(ctx, user); err != nil {
		return nil, err
	}

	createdUser, err := u.psql.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserService) SignIn() error {
	return nil
}