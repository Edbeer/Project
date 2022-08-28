package service

import (
	"context"
	"testing"

	"github.com/Edbeer/Project/internal/entity"
	mockredis "github.com/Edbeer/Project/internal/storage/redis/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestService_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionRedis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis)

	ctx := context.Background()
	session := &entity.Session{}
	rT := "refresh token"

	mockSessionRedis.EXPECT().CreateSession(gomock.Any(), gomock.Eq(session), 10).Return(rT, nil)

	createdSession, err := sessionService.CreateSession(ctx, session, 10)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotEqual(t, createdSession, "")
}

func TestService_GetUserID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionRedis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis)

	ctx := context.Background()
	session := &entity.Session{
		UserID: uuid.New(),
	}
	uid := session.UserID
	rT := "refresh token"

	mockSessionRedis.EXPECT().GetUserID(gomock.Any(), gomock.Eq(rT)).Return(uid, nil)

	uid, err := sessionService.GetUserID(ctx, rT)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, uid)
}

func TestService_DeleteSessionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionRedis(ctrl)
	sessionService := NewSessionService(nil, mockSessionRedis)

	ctx := context.Background()
	rT := "refresh token"

	mockSessionRedis.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(rT)).Return(nil)

	err := sessionService.DeleteSession(ctx, rT)
	require.NoError(t, err)
	require.Nil(t, err)
}