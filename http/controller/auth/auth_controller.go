package auth

import (
	"context"
	"net/http"

	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/auth/service"
	"github.com/mccormack-harry/errors/http/errors"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	ForgotPassword(ctx context.Context, email *dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, email *dto.ResetPasswordRequest) error
	SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.Session, error)
	SignOut(ctx context.Context) error
}

type Controller struct {
	goyave.Component
	AuthService Service
}

func NewController(AuthService Service) *Controller {
	return &Controller{
		AuthService: AuthService,
	}
}

func (c *Controller) Init(server *goyave.Server) {
	if c.AuthService == nil {
		c.AuthService = server.Service(service.Auth).(Service)
	}
	c.Component.Init(server)
}

func (c *Controller) RegisterRoutes(router *goyave.Router) {
	authRouter := router.Subrouter("/auth")
	authRouter.Get("/me", c.me)
	authRouter.Post("/forgot-password", c.forgotPassword)
	authRouter.Post("/reset-password", c.resetPassword)
	authRouter.Post("/signin", c.signin)
	authRouter.Post("/signout", c.signout)
}

func (c *Controller) me(response *goyave.Response, request *goyave.Request) {
	user := request.User
	if user == nil {
		response.Status(http.StatusUnauthorized)
		return
	}
	response.JSON(http.StatusOK, user)
}

func (c *Controller) forgotPassword(response *goyave.Response, request *goyave.Request) {
	req := typeutil.MustConvert[*dto.ForgotPasswordRequest](request.Data)
	err := c.AuthService.ForgotPassword(request.Context(), req)
	if errors.WriteError(response, err) {
		return
	}
	response.Status(http.StatusOK)
}

func (c *Controller) resetPassword(response *goyave.Response, request *goyave.Request) {
	req := typeutil.MustConvert[*dto.ResetPasswordRequest](request.Data)
	if err := c.AuthService.ResetPassword(request.Context(), req); err != nil {
		errors.WriteError(response, err)
	}
}

func (c *Controller) signin(response *goyave.Response, request *goyave.Request) {
	req := typeutil.MustConvert[*dto.SignInRequest](request.Data)
	session, err := c.AuthService.SignIn(request.Context(), req)
	if errors.WriteError(response, err) {
		return
	}
	response.Cookie(&http.Cookie{
		Name:     "session_token",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
	})
	response.Status(http.StatusOK)
}

func (c *Controller) signout(response *goyave.Response, request *goyave.Request) {
	err := c.AuthService.SignOut(request.Context())
	if errors.WriteError(response, err) {
		return
	}
	response.Cookie(&http.Cookie{
		Name:     "session_token",
		HttpOnly: true,
		Secure:   true,
		Value:    "",
		MaxAge:   -1, // Delete the cookie
	})
	response.Status(http.StatusOK)
}
