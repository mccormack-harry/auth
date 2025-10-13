package dto

import (
	"github.com/guregu/null/v6"
	t "goyave.dev/goyave/v5/util/typeutil"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type CreateUser struct {
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Phone     string      `json:"phone"`
	Password  null.String `json:"password"`
}

type UpdateUser struct {
	Email     t.Undefined[string] `json:"email"`
	FirstName t.Undefined[string] `json:"firstName"`
	LastName  t.Undefined[string] `json:"lastName"`
	Phone     t.Undefined[string] `json:"phone"`
	Password  t.Undefined[string] `json:"password"`
}
