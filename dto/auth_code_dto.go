package dto

import (
	"time"

	t "goyave.dev/goyave/v5/util/typeutil"
)

type AuthCode struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`

	User *User `json:"user"`
}

type CreateAuthCode struct {
	Email     string
	ExpiresAt t.Undefined[time.Time]
}
