package repository

import (
	"context"

	"github.com/mccormack-harry/auth/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/session"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (r *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user *model.User
	db := session.DB(ctx, r.DB).Where("id = ?", id).First(&user)
	return user, errors.New(db.Error)
}

func (r *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	db := session.DB(ctx, r.DB).Where("email = ?", email).First(&user)
	return user, errors.New(db.Error)
}

func (r *User) Update(ctx context.Context, user *model.User) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Save(&user)
	return errors.New(db.Error)
}

func (r *User) Create(ctx context.Context, user *model.User) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Omit("id").Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&user)
	return errors.New(db.Error)
}
