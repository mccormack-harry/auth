package repository

import (
	"context"

	"github.com/mccormack-harry/auth/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/session"
)

type Session struct {
	DB *gorm.DB
}

func NewSession(db *gorm.DB) *Session {
	return &Session{
		DB: db,
	}
}

func (r *Session) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	var s *model.Session
	db := session.DB(ctx, r.DB).Preload("User").Where("token = ?", token).First(&s)
	return s, errors.New(db.Error)
}

func (r *Session) Create(ctx context.Context, s *model.Session) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Omit("id").Create(s)
	return errors.New(db.Error)
}

func (r *Session) Delete(ctx context.Context, id string) error {
	sb := session.DB(ctx, r.DB).Delete(&model.Session{}, "id = ?", id)
	if sb.RowsAffected == 0 {
		return errors.New(gorm.ErrRecordNotFound)
	}
	return errors.New(sb.Error)
}
