package password

import (
	"context"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/database/model"
	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/service/errors"
	"golang.org/x/crypto/bcrypt"
	"goyave.dev/goyave/v5/util/typeutil"
)

type PasswordRepository interface {
	GetByUserID(ctx context.Context, userID string) (*model.Password, error)
	Create(ctx context.Context, password *model.Password) error
}

type Service struct {
	PasswordRepository PasswordRepository
}

func NewService(repository PasswordRepository) *Service {
	return &Service{PasswordRepository: repository}
}

func (s *Service) Name() string {
	return service.Password
}

func (s *Service) Check(ctx context.Context, check *dto.CheckPassword) error {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return errors.Forbidden("only system can check passwords")
	}
	password, err := s.PasswordRepository.GetByUserID(ctx, check.UserID)
	if err != nil {
		return errors.New(err)
	}
	if checkPassword(password.Password, check.Password) {
		return nil
	}
	return errors.Unauthorized("invalid credentials")
}

func (s *Service) Create(ctx context.Context, create *dto.CreatePassword) (*dto.Password, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can create passwords")
	}
	create.Password, err = hashPassword(create.Password)
	if err != nil {
		return nil, errors.New(err)
	}
	password := typeutil.Copy(&model.Password{}, create)
	err = s.PasswordRepository.Create(ctx, password)
	if err != nil {
		return nil, err
	}
	return typeutil.MustConvert[*dto.Password](password), nil
}

// TODO hasher interface?

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
