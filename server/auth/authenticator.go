package auth

import (
	"context"
	"strings"

	"github.com/pixb/echo-demo/go-server/store"
)

const PersonalAccessTokenPrefix = "pat_"

type AuthResult struct {
	User        *store.User
	Claims      *UserClaims
	AccessToken string
}

type Authenticator struct {
	Store  *store.Store
	Secret string
}

func NewAuthenticator(store *store.Store, secret string) *Authenticator {
	return &Authenticator{Store: store, Secret: secret}
}

func (a *Authenticator) Authenticate(ctx context.Context, authHeader string) *AuthResult {
	token := ExtractBearerToken(authHeader)
	if token == "" {
		return nil
	}

	if !strings.HasPrefix(token, PersonalAccessTokenPrefix) {
		claims, err := a.AuthenticateByAccessTokenV2(token)
		if err == nil && claims != nil {
			return &AuthResult{
				Claims:      claims,
				AccessToken: token,
			}
		}
	}

	return nil
}

func (a *Authenticator) AuthenticateByAccessTokenV2(token string) (*UserClaims, error) {
	claims, err := ValidateAccessToken(token, a.Secret)
	if err != nil {
		return nil, err
	}
	return &UserClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}
