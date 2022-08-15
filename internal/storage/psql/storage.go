package psql

import "github.com/jmoiron/sqlx"

// Storage psql
type Storage struct {
	User *UserStorage
}

func NewStorage(psql *sqlx.DB) *Storage {
	return &Storage{
		User: newUserStorage(psql),
	}
}