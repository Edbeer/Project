package service

import (
	"context"

	"github.com/Edbeer/Project/pkg/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"

	"github.com/Edbeer/Project/internal/entity"

	"github.com/Edbeer/Project/config"
)

// Token Manager interface
type Manager interface {
	GenerateJWTToken(user *entity.User) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() string
}

// User psql storage interface
type UserPsql interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}

// User service
type UserService struct {
	config       *config.Config
	psql         UserPsql
	tokenManager Manager
}

// New user service constructor
func newUserService(config *config.Config, psql UserPsql, tokenManager Manager) *UserService {
	return &UserService{
		config:       config,
		psql:         psql,
		tokenManager: tokenManager,
	}
}

// Sign-up user
func (u *UserService) SignUp(ctx context.Context, user *entity.User) (*entity.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.SignUp")
	defer span.Finish()

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

	accessToken, err := u.tokenManager.GenerateJWTToken(createdUser)
	if err != nil {
		return nil, err
	}

	return &entity.UserWithToken{
		User:        createdUser,
		AccessToken: accessToken,
	}, nil
}

// Sign-in user
func (u *UserService) SignIn(ctx context.Context, user *entity.User) (*entity.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.SignIn")
	defer span.Finish()
	
	foundUser, err := u.psql.FindUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.tokenManager.GenerateJWTToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &entity.UserWithToken{
		User:        foundUser,
		AccessToken: accessToken,
	}, nil
}

// Get user by id
func (u *UserService) GetUserByID(ctx context.Context, userId uuid.UUID) (*entity.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserByID")
	defer span.Finish()
	
	foundUser, err := u.psql.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.tokenManager.GenerateJWTToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &entity.UserWithToken{
		User:        foundUser,
		AccessToken: accessToken,
	}, nil
}
