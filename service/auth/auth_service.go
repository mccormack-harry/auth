package auth

import (
	"context"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/service/errors"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dto.User, error)
	GetByEmail(ctx context.Context, email string) (*dto.User, error)
	Update(ctx context.Context, id string, user *dto.UpdateUser) error
}

type SessionService interface {
	GetByToken(ctx context.Context, token string) (*dto.Session, error)
	Create(ctx context.Context, session *dto.CreateSession) (*dto.Session, error)
	Delete(ctx context.Context, session *dto.Session) error
}

type AuthCodeService interface {
	Create(ctx context.Context, create *dto.CreateAuthCode) (*dto.AuthCode, error)
	Validate(ctx context.Context, email, code string) (*dto.AuthCode, error)
}

type PasswordService interface {
	Create(ctx context.Context, create *dto.CreatePassword) (*dto.Password, error)
	Check(ctx context.Context, check *dto.CheckPassword) error
}

type MailService interface {
	SendEmail(recipients []string, subject, html string) error
}

type MailTemplates interface {
	ForgotPassword(*dto.User, *dto.AuthCode) (subject string, html string)
}

type Service struct {
	UserService     UserService
	PasswordService PasswordService
	SessionService  SessionService
	AuthCodeService AuthCodeService
	MailService     MailService
	MailTemplates   MailTemplates
}

func NewService(
	userService UserService,
	passwordService PasswordService,
	sessionService SessionService,
	authCodeService AuthCodeService,
	mailService MailService,
	mailTemplates MailTemplates,
) *Service {
	return &Service{
		UserService:     userService,
		PasswordService: passwordService,
		SessionService:  sessionService,
		AuthCodeService: authCodeService,
		MailService:     mailService,
		MailTemplates:   mailTemplates,
	}
}

func (s *Service) Name() string {
	return service.Auth
}

func (s *Service) Authenticate(ctx context.Context, token string) (*dto.User, *dto.Session, error) {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	if actor.Type != actorctx.ActorTypeSystem {
		return nil, nil, errors.Forbidden("only system can authenticate users")
	}
	session, err := s.SessionService.GetByToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}
	return session.User, session, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	user, err := s.UserService.GetByEmail(actorctx.WithSystem(ctx), req.Email)
	if err != nil {
		return errors.Wrap(errors.NotFound("user not found"), err)
	}
	create := &dto.CreateAuthCode{
		Email: user.Email,
	}
	authCode, err := s.AuthCodeService.Create(actorctx.WithSystem(ctx), create)
	if err != nil {
		return err
	}

	subject, html := s.MailTemplates.ForgotPassword(user, authCode)
	err = s.MailService.SendEmail(
		[]string{authCode.Email},
		subject,
		html,
	)
	return errors.New(err)
}

func (s *Service) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	authCode, err := s.AuthCodeService.Validate(actorctx.WithSystem(ctx), req.Email, req.Code)
	if err != nil {
		return errors.New(err)
	}

	if authCode.User == nil {
		return errors.New("auth code is not associated with a user")
	}

	createPassword := &dto.CreatePassword{
		UserID:   authCode.User.ID,
		Password: req.Password,
	}
	_, err = s.PasswordService.Create(actorctx.WithSystem(ctx), createPassword)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Service) SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.Session, error) {
	actor, _ := actorctx.ActorFromContext(ctx)
	if actor != nil {
		return nil, errors.Forbidden("already signed in")
	}
	user, err := s.UserService.GetByEmail(actorctx.WithSystem(ctx), req.Email)
	if err != nil {
		return nil, errors.Wrap(errors.Unauthorized("invalid email"), errors.New(err))
	}
	check := &dto.CheckPassword{
		UserID:   user.ID,
		Password: req.Password,
	}
	err = s.PasswordService.Check(actorctx.WithSystem(ctx), check)
	if err != nil {
		return nil, errors.Unauthorized("invalid credentials")
	}
	createSession := &dto.CreateSession{
		UserID: user.ID,
	}
	session, err := s.SessionService.Create(actorctx.WithSystem(ctx), createSession)
	if err != nil {
		return nil, errors.New(err)
	}
	return session, nil
}

func (s *Service) SignOut(ctx context.Context) error {
	actor, err := actorctx.ActorFromContext(ctx)
	if err != nil {
		return err
	}
	// TODO admin sign out?
	if actor.Type == actorctx.ActorTypeSystem {
		return errors.Forbidden("system cannot sign out")
	}
	if actor.Session == nil {
		return errors.New("no session found for user")
	}
	return s.SessionService.Delete(actorctx.WithSystem(ctx), actor.Session)
}
