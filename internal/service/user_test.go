package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
	mockstorage "github.com/Edbeer/Project/internal/storage/psql/mock"
	"github.com/Edbeer/Project/pkg/jwt"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.Server{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.Server.JwtSecretKey)
	mockUserStorage := mockstorage.NewMockUserPsql(ctrl)
	userService := newUserService(config, mockUserStorage, manager)

	user := &entity.User{
		Name:     "PavelV",
		Password: "12345678",
		Email:    "edbeermtn@gmail.com",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "UserService.SignUp")
	defer span.Finish()

	mockUserStorage.EXPECT().FindUserByEmail(ctxWithTrace, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockUserStorage.EXPECT().Create(ctxWithTrace, gomock.Eq(user)).Return(user, nil)

	createdUser, err := userService.SignUp(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Nil(t, err)
}

func TestService_SignIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.Server{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.Server.JwtSecretKey)
	mockUserStorage := mockstorage.NewMockUserPsql(ctrl)
	userService := newUserService(config, mockUserStorage, manager)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "UserService.SignUp")
	defer span.Finish()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &entity.User{
		Password: string(hashPassword),
		Email:    "edbeermtn@gmail.com",
	}

	mockUserStorage.EXPECT().FindUserByEmail(ctxWithTrace, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := userService.SignIn(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		Server: config.Server{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.Server.JwtSecretKey)
	mockUserStorage := mockstorage.NewMockUserPsql(ctrl)
	userService := newUserService(config, mockUserStorage, manager)

	user := &entity.User{
		Password: "12345678",
		Email: "edbeermtn@gmail.com",
	}

	ctx := context.Background()
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "UserService.SignUp")
	defer span.Finish()

	mockUserStorage.EXPECT().GetUserByID(ctxWithTrace, gomock.Eq(user.ID)).Return(user, nil)

	u, err := userService.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}