package actorctx

import (
	"context"

	"github.com/mccormack-harry/auth/dto"
	"github.com/mccormack-harry/errors/service/errors"
)

type ActorType string

const (
	ActorTypeSystem ActorType = "system"
	ActorTypeAdmin  ActorType = "admin"
	ActorTypeUser   ActorType = "user"
)

type Actor struct {
	Type    ActorType
	User    *dto.User
	Session *dto.Session
}

type actorKey struct{}

func WithSystem(ctx context.Context) context.Context {
	return context.WithValue(ctx, actorKey{}, &Actor{Type: ActorTypeSystem})
}

func WithAdmin(ctx context.Context, user *dto.User, session *dto.Session) context.Context {
	return context.WithValue(ctx, actorKey{}, &Actor{Type: ActorTypeAdmin, User: user, Session: session})
}

func WithUser(ctx context.Context, user *dto.User, session *dto.Session) context.Context {
	return context.WithValue(ctx, actorKey{}, &Actor{Type: ActorTypeUser, User: user, Session: session})
}

func ActorFromContext(ctx context.Context) (*Actor, error) {
	actor, ok := ctx.Value(actorKey{}).(*Actor)
	if !ok {
		return nil, errors.Unauthorized("unauthorized")
	}
	return actor, nil
}
