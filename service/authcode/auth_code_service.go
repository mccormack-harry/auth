package authcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/database/model"
	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/service/errors"
	"goyave.dev/goyave/v5/util/session"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	Find(ctx context.Context, email string, code string) (*model.AuthCode, error)
	Create(ctx context.Context, code *model.AuthCode) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	Session    session.Session
	Repository Repository
}

func NewService(session session.Session, repository Repository) *Service {
	return &Service{
		Repository: repository,
		Session:    session,
	}
}

func (s *Service) Name() string {
	return service.AuthCode
}

func (s *Service) Create(ctx context.Context, create *dto.CreateAuthCode) (*dto.AuthCode, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can create auth codes")
	}
	if !create.ExpiresAt.Present {
		// TODO config auth_code.expiry
		create.ExpiresAt = typeutil.NewUndefined(time.Now().Add(15 * time.Minute))
	}
	authCode := typeutil.MustConvert[*model.AuthCode](create)
	code, err := generateAuthCode()
	if err != nil {
		return nil, errors.New(err)
	}
	authCode.Code = typeutil.NewUndefined(code)
	if err := s.Repository.Create(ctx, authCode); err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.AuthCode](authCode), nil
}

func (s *Service) Validate(ctx context.Context, email, code string) (*dto.AuthCode, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can validate auth codes")
	}
	var res *dto.AuthCode
	err = s.Session.Transaction(ctx, func(ctx context.Context) error {
		authCode, err := s.Repository.Find(ctx, email, code)
		if err != nil {
			return errors.New(err)
		}
		if !authCode.ExpiresAt.Present || time.Now().After(authCode.ExpiresAt.Val) {
			return errors.Unauthorized("auth code expired")
		}
		res = typeutil.MustConvert[*dto.AuthCode](authCode)
		if err := s.Repository.Delete(ctx, authCode.ID.Val); err != nil {
			// TODO logging?
			fmt.Printf("error occurred deleting login code: %v", err)
		}
		return nil
	})
	return res, errors.New(err)
}

func generateAuthCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}
