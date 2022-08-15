package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID         uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty,uuid"`
	Name       string    `json:"name" db:"name" validate:"required_with,lte=30"`
	Email      string    `json:"email" db:"email" validate:"omitempty,email"`
	Password   string    `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	Created_at time.Time `json:"created_at" db:"created_at"`
}

// User with token
type UserWithToken struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}

// Input User
type InputUser struct {
	Name     string
	Email    string
	Password string
}

// Compare user password and payload
func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Prepare user struct for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	return nil
}

// Sanitize password
func (u *User) SanitizePasswor() {
	u.Password = ""
}
