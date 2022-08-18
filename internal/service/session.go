package service

import (
	"context"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
)


type SessionStorage interface {
	CreateSession(ctx context.Context, session *entity.Session, expire int) (string, error)
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
