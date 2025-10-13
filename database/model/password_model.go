package model

import "time"

type Password struct {
	ID        string    `gorm:"primarykey" json:",omitzero"`
	Password  string    `json:",omitzero"`
	UserID    string    `json:",omitzero"`
	CreatedAt time.Time `json:",omitzero"`
}

func (Password) TableName() string {
	return "passwords"
}
