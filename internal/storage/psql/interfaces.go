//go:generate mockgen -source interfaces.go -destination mock/psql_storage_mock.go -package mock
package psql

import (
	"context"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
)

// User psql storage interface
type UserPsql interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}