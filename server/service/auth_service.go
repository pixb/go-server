package service

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthStore is an interface that defines the methods needed by AuthService
type AuthStore interface {
	GetUserByUsername(ctx context.Context, username string) (*store.User, error)
	CreateRefreshToken(ctx context.Context, create *store.CreateRefreshToken) (*store.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, update *store.UpdateRefreshToken) (*store.RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (*store.RefreshToken, error)
	GetUser(ctx context.Context, find *store.FindUser) (*store.User, error)
	Ping(ctx context.Context) error
	Close() error
}

type AuthService struct {
	Secret string
	Store  AuthStore
}

func NewAuthService(secret string, store AuthStore) *AuthService {
	return &AuthService{
		Secret: secret,
		Store:  store,
	}
}

func (s *AuthService) Login(ctx context.Context, req *v1pb.LoginRequest) (*v1pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username and password are required"))
	}

	user, err := s.Store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("username not found"))
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid password"))
	}

	// Check if password has expired
	if time.Now().After(user.PasswordExpires) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("password expired"))
	}

	// Generate access token
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, s.Secret)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate access token"))
	}

	// Generate refresh token
	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate refresh token"))
	}

	// Save refresh token to database
	_, err = s.Store.CreateRefreshToken(ctx, &store.CreateRefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(auth.RefreshTokenDuration),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to save refresh token"))
	}

	// Calculate access token expiration time
	accessTokenExpiresAt := time.Now().Add(auth.AccessTokenDuration)

	return &v1pb.LoginResponse{
		AccessToken:          accessToken,
		RefreshToken:         refreshTokenString,
		AccessTokenExpiresAt: timestamppb.New(accessTokenExpiresAt),
		User: &v1pb.User{
			Id:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			Nickname:          user.Nickname,
			Phone:             user.Phone,
			Role:              auth.StringToRole(user.Role),
			PasswordExpiresAt: timestamppb.New(user.PasswordExpires),
			CreatedAt:         timestamppb.New(user.CreatedAt),
			UpdatedAt:         timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *v1pb.RefreshTokenRequest) (*v1pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("refresh token is required"))
	}

	// Validate refresh token
	refreshToken, err := s.Store.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to validate refresh token"))
	}
	if refreshToken == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid refresh token"))
	}

	// Check if refresh token is expired or revoked
	if refreshToken.Revoked || time.Now().After(refreshToken.ExpiresAt) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("refresh token expired or revoked"))
	}

	// Get user information
	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &refreshToken.UserID})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get user"))
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	// Revoke old refresh token
	_, err = s.Store.UpdateRefreshToken(ctx, &store.UpdateRefreshToken{
		ID:      refreshToken.ID,
		Revoked: &[]bool{true}[0],
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to revoke old refresh token"))
	}

	// Generate new access token
	newAccessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, s.Secret)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate access token"))
	}

	// Generate new refresh token
	newRefreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate refresh token"))
	}

	// Save new refresh token to database
	_, err = s.Store.CreateRefreshToken(ctx, &store.CreateRefreshToken{
		UserID:    user.ID,
		Token:     newRefreshTokenString,
		ExpiresAt: time.Now().Add(auth.RefreshTokenDuration),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to save refresh token"))
	}

	// Calculate access token expiration time
	accessTokenExpiresAt := time.Now().Add(auth.AccessTokenDuration)

	return &v1pb.RefreshTokenResponse{
		AccessToken:          newAccessToken,
		RefreshToken:         newRefreshTokenString,
		AccessTokenExpiresAt: timestamppb.New(accessTokenExpiresAt),
		User: &v1pb.User{
			Id:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			Nickname:          user.Nickname,
			Phone:             user.Phone,
			Role:              auth.StringToRole(user.Role),
			PasswordExpiresAt: timestamppb.New(user.PasswordExpires),
			CreatedAt:         timestamppb.New(user.CreatedAt),
			UpdatedAt:         timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *v1pb.ValidateTokenRequest) (*v1pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("token is required"))
	}

	// Validate token
	_, err := auth.ValidateAccessToken(req.Token, s.Secret)
	if err != nil {
		return &v1pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &v1pb.ValidateTokenResponse{
		Valid: true,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *v1pb.LogoutRequest) (*v1pb.LogoutResponse, error) {
	if req.Token == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("token is required"))
	}

	// Validate token
	_, err := auth.ValidateAccessToken(req.Token, s.Secret)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}

	// For now, just return success (implement token revocation later)
	return &v1pb.LogoutResponse{
		Success: true,
	}, nil
}
