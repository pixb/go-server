package auth

import (
	"context"

	"github.com/pixb/go-server/store"
)

type UserClaims struct {
	UserID   int64
	Username string
	Role     string
}

type ContextKey int

const (
	UserIDContextKey ContextKey = iota
	AccessTokenContextKey
	UserClaimsContextKey
)

func GetUserID(ctx context.Context) int64 {
	if v, ok := ctx.Value(UserIDContextKey).(int64); ok {
		return v
	}
	return 0
}

func GetUserClaims(ctx context.Context) *UserClaims {
	if v, ok := ctx.Value(UserClaimsContextKey).(*UserClaims); ok {
		return v
	}
	return nil
}

func SetUserClaimsInContext(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, UserClaimsContextKey, claims)
}

func SetUserInContext(ctx context.Context, user *store.User, accessToken string) context.Context {
	ctx = context.WithValue(ctx, UserIDContextKey, user.ID)
	if accessToken != "" {
		ctx = context.WithValue(ctx, AccessTokenContextKey, accessToken)
	}
	return ctx
}
