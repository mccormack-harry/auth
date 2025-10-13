package middleware

import (
	"context"
	"net/http"

	"github.com/mccormack-harry/auth/actorctx"
	"github.com/mccormack-harry/auth/dto"
	"github.com/samber/lo"
	"goyave.dev/goyave/v5"
)

type Service interface {
	Authenticate(ctx context.Context, token string) (*dto.User, *dto.Session, error)
}

type Auth struct {
	goyave.Component
	AuthService Service
}

func NewAuth(authService Service) *Auth {
	return &Auth{
		AuthService: authService,
	}
}

func (m *Auth) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		var sessionToken string
		var hasCookie bool
		cookie, ok := lo.Find(request.Cookies(), func(cookie *http.Cookie) bool {
			return cookie.Name == "session_token"
		})
		if ok {
			hasCookie = true
			sessionToken = cookie.Value
		}
		if hasCookie && sessionToken != "" {
			ctx := actorctx.WithSystem(request.Context())
			user, session, err := m.AuthService.Authenticate(ctx, sessionToken)
			if err == nil {
				request.User = user
				request.Extra["session"] = session
				request.WithContext(actorctx.WithUser(request.Context(), user, session))
				next(response, request)
				return
			}
		}

		//request.WithContext(actorctx.WithGuest(request.Context()))
		if hasCookie {
			response.Cookie(&http.Cookie{
				Name:     "session_token",
				HttpOnly: true,
				Secure:   true,
				Value:    "",
				MaxAge:   -1,
			})
		}
		next(response, request)
	}
}
