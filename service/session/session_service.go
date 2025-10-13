package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"time"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/database/model"
	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/service/errors"
	"goyave.dev/goyave/v5/config"
	"goyave.dev/goyave/v5/util/typeutil"
)

func init() {
	config.Register("session.expiry", config.Entry{
		Value:    int(time.Hour) * 24 * 28, // 1 Month
		Type:     reflect.Int,
		Required: true,
	})
}

type Repository interface {
	GetByToken(ctx context.Context, token string) (*model.Session, error)
	Create(ctx context.Context, session *model.Session) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	Config     *config.Config
	Repository Repository
}

func NewService(config *config.Config, repository Repository) *Service {
	return &Service{
		Config:     config,
		Repository: repository,
	}
}

func (s *Service) Name() string {
	return service.Session
}

func (s *Service) GetByToken(ctx context.Context, token string) (*dto.Session, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can get sessions")
	}
	session, err := s.Repository.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New(err)
	}
	if !session.ExpiresAt.Present || session.ExpiresAt.Val.Before(time.Now()) {
		return nil, errors.Expired("session expired")
	}
	return typeutil.MustConvert[*dto.Session](session), nil
}

func (s *Service) Create(ctx context.Context, session *dto.CreateSession) (*dto.Session, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can create sessions")
	}
	newSession := typeutil.Copy(&model.Session{}, session)
	token, err := generateSessionToken()
	if err != nil {
		return nil, errors.New(err)
	}
	newSession.Token = typeutil.NewUndefined(token)
	newSession.ExpiresAt = typeutil.NewUndefined(s.getExpiry())
	if err := s.Repository.Create(ctx, newSession); err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.Session](newSession), nil
}

func (s *Service) Delete(ctx context.Context, session *dto.Session) error {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return errors.Forbidden("only system can delete sessions")
	}
	err = s.Repository.Delete(ctx, session.ID)
	return errors.New(err)
}

func generateSessionToken() (string, error) {
	const tokenLength = 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New(err)
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func (s *Service) getExpiry() time.Time {
	exp := s.Config.GetInt("session.expiry")
	return time.Now().Add(time.Duration(exp) * time.Millisecond)
}
