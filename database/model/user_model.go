package model

import (
	"time"

	t "goyave.dev/goyave/v5/util/typeutil"
)

type User struct {
	Email t.Undefined[string] `json:",omitzero"`

	FirstName t.Undefined[string] `json:",omitzero"`
	LastName  t.Undefined[string] `json:",omitzero"`
	Phone     t.Undefined[string] `json:",omitzero"`

	CreatedAt t.Undefined[time.Time] `json:",omitzero"`
	UpdatedAt t.Undefined[time.Time] `json:",omitzero"`

	ID t.Undefined[string] `gorm:"primarykey" json:",omitzero"`
}

func (User) TableName() string {
	return "users"
}
