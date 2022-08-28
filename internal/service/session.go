package service

import (
	"context"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionService.CreateSession")
	defer span.Finish()
	return s.session.CreateSession(ctx, session, expire)
}

func (s *SessionService) GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionService.GetUserID")
	defer span.Finish()
	return s.session.GetUserID(ctx, refreshToken)
}

func (s *SessionService) DeleteSession(ctx context.Context, refreshToken string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionService.DeleteSession")
	defer span.Finish()
	return s.session.DeleteSession(ctx, refreshToken)
}