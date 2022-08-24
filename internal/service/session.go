package service

import (
	"context"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
)

// Session storage interface
type SessionStorage interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}

// User service
type SessionService struct {
	config  *config.Config
	session SessionStorage
}

// New user service constructor
func NewSessionService(config *config.Config, session SessionStorage) *SessionService {
	return &SessionService{
		config:  config,
		session: session,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error) {
	return s.session.CreateSession(ctx, session, expire)
}

func (s *SessionService) GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	return s.session.GetUserID(ctx, refreshToken)
}

func (s *SessionService) DeleteSession(ctx context.Context, refreshToken string) error {
	return s.session.DeleteSession(ctx, refreshToken)
}