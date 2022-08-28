package psql

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Edbeer/Project/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newUserStorage(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		columns := []string{
			"name",
			"email",
			"password"}
		rows := sqlmock.NewRows(columns).AddRow(
			"PavelV",
			"edbeermtn@gmail.com",
			"12345678",
		)

		user := &entity.User{
			Name:     "PavelV",
			Email:    "edbeermtn@gmail.com",
			Password: "12345678",
		}

		query := `INSERT INTO users (name, email, password, created_at) 
			VALUES ($1, $2, $3, now()) 
			RETURNING *`
		mock.ExpectQuery(query).WithArgs(
			&user.Name, &user.Email, &user.Password,
		).WillReturnRows(rows)

		createdUser, err := userStorage.Create(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func Test_FindUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newUserStorage(sqlxDB)

	t.Run("FindUserByEmail", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"PavelV",
			"edbeermtn@gmail.com",
		)

		testUser := &entity.User{
			ID:    uid,
			Name:  "PavelV",
			Email: "edbeermtn@gmail.com",
		}

		query := `SELECT user_id, name, email, password, created_at
			FROM users
			WHERE email = $1`
		mock.ExpectQuery(query).WithArgs(&testUser.Email).WillReturnRows(rows)

		user, err := userStorage.FindUserByEmail(context.Background(), testUser)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.Name, testUser.Name)
	})
}

func Test_GetUserByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newUserStorage(sqlxDB)

	t.Run("GetUserByID", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"PavelV",
			"edbeermtn@gmail.com",
		)

		testUser := &entity.User{
			ID:    uid,
			Name:  "PavelV",
			Email: "edbeermtn@gmail.com",
		}

		query := `SELECT user_id, name, email, password, created_at
			FROM users
			WHERE user_id = $1`
		mock.ExpectQuery(query).WithArgs(uid).WillReturnRows(rows)

		user, err := userStorage.GetUserByID(context.Background(), uid)
		require.NoError(t, err)
		require.Equal(t, user.Name, testUser.Name)
		fmt.Printf("test user: %s \n", testUser.Name)
		fmt.Printf("user: %s \n", user.Name)
	})
}
