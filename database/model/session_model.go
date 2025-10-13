package model

import (
	"time"

	t "goyave.dev/goyave/v5/util/typeutil"
)

type Session struct {
	UserID    t.Undefined[string]    `json:",omitzero"`
	Token     t.Undefined[string]    `json:",omitzero"`
	CreatedAt t.Undefined[time.Time] `json:",omitzero"`
	UpdatedAt t.Undefined[time.Time] `json:",omitzero"`
	ExpiresAt t.Undefined[time.Time] `json:",omitzero"`

	User *User               `gorm:"foreignKey:UserID" json:",omitzero"`
	ID   t.Undefined[string] `gorm:"primarykey" json:",omitzero"`
}

func (Session) TableName() string {
	return "sessions"
}
