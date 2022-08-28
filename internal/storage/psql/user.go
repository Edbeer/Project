package psql

import (
	"context"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// User psql storage
type UserStorage struct {
	psql *sqlx.DB
}

// New user storage constructor
func newUserStorage(psql *sqlx.DB) *UserStorage {
	return &UserStorage{psql: psql}
}

// Create user
func (r *UserStorage) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserPsql.Create")
	defer span.Finish()

	u := &entity.User{}
	query := `INSERT INTO users (name, email, password, created_at) 
			VALUES ($1, $2, $3, now()) 
			RETURNING *`
	if err := r.psql.QueryRowxContext(ctx, query, 
		&user.Name, &user.Email, &user.Password,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "UserStoragePsql.Create.StructScan")
	}
	return u, nil
}

// Find user by email
func (r *UserStorage) FindUserByEmail(ctx context.Context, user *entity.User) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserPsql.FindUserByEmail")
	defer span.Finish()
	
	foundUser := &entity.User{}
	query := `SELECT user_id, name, email, password, created_at
			FROM users
			WHERE email = $1`
	if err := r.psql.QueryRowxContext(ctx, query, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "UserStoragePsql.FindUserByEmail.StructScan")
	}
	return foundUser, nil
}

// Get user by id
func (r *UserStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserPsql.GetUserByID")
	defer span.Finish()
	
	u := &entity.User{}
	query := `SELECT user_id, name, email, password, created_at
		FROM users
		WHERE user_id = $1`
	if err := r.psql.QueryRowxContext(ctx, query, userID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "AuthStoragePsql.GetUserByID.StructScan")
	}

	return u, nil
}
