package model

import (
	"time"

	t "goyave.dev/goyave/v5/util/typeutil"
)

type AuthCode struct {
	ID        t.Undefined[string]    `gorm:"primaryKey" json:",omitzero"`
	Email     t.Undefined[string]    `json:",omitzero"`
	Code      t.Undefined[string]    `json:",omitzero"`
	CreatedAt t.Undefined[time.Time] `json:",omitzero"`
	ExpiresAt t.Undefined[time.Time] `json:",omitzero"`

	User *User `gorm:"foreignKey:Email;references:Email" json:",omitempty"`
}
