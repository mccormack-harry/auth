package repository

import (
	"context"

	"github.com/mccormack-harry/auth/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/session"
)

type AuthCode struct {
	DB *gorm.DB
}

func NewAuthCode(db *gorm.DB) *AuthCode {
	return &AuthCode{DB: db}
}

func (r *AuthCode) Find(ctx context.Context, email string, code string) (*model.AuthCode, error) {
	var a *model.AuthCode
	db := session.DB(ctx, r.DB).Where("email = ? AND code = ?", email, code).Preload("User").First(&a)
	return a, errors.New(db.Error)
}

func (r *AuthCode) Create(ctx context.Context, code *model.AuthCode) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Omit("id").Create(&code)
	return errors.New(db.Error)
}

func (r *AuthCode) Delete(ctx context.Context, id string) error {
	db := session.DB(ctx, r.DB).Delete(&model.AuthCode{}, "id = ?", id)
	if db.RowsAffected == 0 {
		return errors.New(gorm.ErrRecordNotFound)
	}
	return errors.New(db.Error)
}
