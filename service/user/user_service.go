package user

import (
	"context"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/database/model"
	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/service/errors"
	"goyave.dev/goyave/v5/util/session"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
}

type PasswordService interface {
	Create(ctx context.Context, create *dto.CreatePassword) (*dto.Password, error)
}

type Service struct {
	Session         session.Session
	Repository      Repository
	PasswordService PasswordService
}

func NewService(
	session session.Session,
	repository Repository,
	passwordService PasswordService,
) *Service {
	s := &Service{
		Session:         session,
		Repository:      repository,
		PasswordService: passwordService,
	}
	return s
}

func (s *Service) Name() string {
	return service.User
}

func (s *Service) GetByID(ctx context.Context, id string) (*dto.User, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		// TODO admin get users?
		return nil, errors.Forbidden("only system can get users")
	}
	user, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.User](user), nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*dto.User, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, errors.New(err)
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, errors.Forbidden("only system can get users")
	}
	user, err := s.Repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.User](user), nil
}

func (s *Service) Update(ctx context.Context, id string, updateUser *dto.UpdateUser) error {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return errors.New(err)
	}
	if actor.Type == actorctx.ActorTypeUser && id != actor.User.ID {
		return errors.Forbidden("user can only update self")
	}
	err = s.Session.Transaction(ctx, func(ctx context.Context) error {
		user, err := s.Repository.GetByID(ctx, id)
		if err != nil {
			return errors.New(err)
		}
		user = typeutil.Copy(user, updateUser)
		err = s.Repository.Update(ctx, user)
		return errors.New(err)
	})
	return errors.New(err)
}

func (s *Service) Create(ctx context.Context, create *dto.CreateUser) (*dto.User, error) {
	user := typeutil.Copy(&model.User{}, create)
	err := s.Session.Transaction(ctx, func(ctx context.Context) error {
		err := s.Repository.Create(ctx, user)
		if err != nil {
			return errors.New(err)
		}
		if create.Password.Valid {
			_, err = s.PasswordService.Create(actorctx.WithSystem(ctx), &dto.CreatePassword{
				UserID:   user.ID.Val,
				Password: create.Password.String,
			})
			if err != nil {
				return errors.New(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New(err)
	}

	return typeutil.MustConvert[*dto.User](user), nil
}
