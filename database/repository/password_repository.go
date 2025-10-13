package repository

import (
	"context"

	"github.com/mccormack-harry/auth/database/model"
	"github.com/mccormack-harry/errors/service/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/goyave/v5/util/session"
)

type Password struct {
	DB *gorm.DB
}

func NewPassword(db *gorm.DB) *Password {
	return &Password{
		DB: db,
	}
}

func (r Password) GetByUserID(ctx context.Context, userID string) (*model.Password, error) {
	var password *model.Password
	db := session.DB(ctx, r.DB).Where("user_id = ?", userID).First(&password)
	return password, errors.New(db.Error)
}

func (r Password) Create(ctx context.Context, password *model.Password) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Omit("id").Create(&password)
	return errors.New(db.Error)
}
