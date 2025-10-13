package dto

import "time"

type CreatePassword struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

type Password struct {
	ID        string    `json:"id"`
	Password  string    `json:"password"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type CheckPassword struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}
