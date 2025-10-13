package dto

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`

	User *User `json:"user,omitempty"`
}

type CreateSession struct {
	UserID string `json:"userId"`
}
