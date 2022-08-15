package entity

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token" redis:"refresh_token"`
	Expire       time.Time `json:"expire" redis:"expire"`
}
