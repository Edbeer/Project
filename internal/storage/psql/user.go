package psql

import (
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
func (u *UserStorage) Create() error {
	return nil
}