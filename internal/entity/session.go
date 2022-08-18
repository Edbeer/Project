package entity

import (

	"github.com/google/uuid"
)

type Session struct {
	RefreshToken string  `json:"refresh_token" redis:"refresh_token"`
	UserID    uuid.UUID `json:"user_id" redis:"user_id"`
}
