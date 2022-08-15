package psql

import (
	"context"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
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
	foundUser := &entity.User{}
	query := `SELECT name, email, password, created_at
			FROM users
			WHERE email = $1`
	if err := r.psql.QueryRowxContext(ctx, query, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "UserStoragePsql.FindUserByEmail.StructScan")
	}
	return foundUser, nil
}
