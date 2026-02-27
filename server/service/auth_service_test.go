package service

import (
	"context"
	"testing"
	"time"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login(t *testing.T) {
	// Create mock store
	mockStore := new(MockStore)

	// Test data
	req := &v1pb.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	// Mock responses
	passwordHash, _ := auth.HashPassword(req.Password)
	mockStore.On("GetUserByUsername", mock.Anything, req.Username).Return(&store.User{
		ID:              1,
		Username:        req.Username,
		Email:           "test@example.com",
		Password:        passwordHash,
		Nickname:        "Test User",
		Phone:           "13800138000",
		Role:            "user",
		PasswordExpires: time.Now().AddDate(0, 0, 90),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil)
	mockStore.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*store.CreateRefreshToken")).Return(&store.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     "testrefreshtoken",
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		CreatedAt: time.Now(),
	}, nil)

	// Create auth service
	authService := NewAuthService("testsecret", mockStore)

	// Test Login
	resp, err := authService.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.NotNil(t, resp.AccessTokenExpiresAt)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.NotNil(t, resp.User.PasswordExpiresAt)
	assert.NotNil(t, resp.User.CreatedAt)
	assert.NotNil(t, resp.User.UpdatedAt)

	// Verify mock calls
	mockStore.AssertExpectations(t)
}

func TestAuthService_RefreshToken(t *testing.T) {
	// Create mock store
	mockStore := new(MockStore)

	// Test data
	req := &v1pb.RefreshTokenRequest{
		RefreshToken: "testrefreshtoken",
	}

	// Mock responses
	mockStore.On("GetRefreshToken", mock.Anything, req.RefreshToken).Return(&store.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     req.RefreshToken,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		CreatedAt: time.Now(),
	}, nil)
	mockStore.On("GetUser", mock.Anything, mock.AnythingOfType("*store.FindUser")).Return(&store.User{
		ID:              1,
		Username:        "testuser",
		Email:           "test@example.com",
		Password:        "hashedpassword",
		Nickname:        "Test User",
		Phone:           "13800138000",
		Role:            "user",
		PasswordExpires: time.Now().AddDate(0, 0, 90),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil)
	mockStore.On("UpdateRefreshToken", mock.Anything, mock.AnythingOfType("*store.UpdateRefreshToken")).Return(&store.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     req.RefreshToken,
		Revoked:   true,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		CreatedAt: time.Now(),
	}, nil)
	mockStore.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("*store.CreateRefreshToken")).Return(&store.RefreshToken{
		ID:        2,
		UserID:    1,
		Token:     "newrefreshtoken",
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		CreatedAt: time.Now(),
	}, nil)

	// Create auth service
	authService := NewAuthService("testsecret", mockStore)

	// Test RefreshToken
	resp, err := authService.RefreshToken(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.NotNil(t, resp.AccessTokenExpiresAt)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.NotNil(t, resp.User.PasswordExpiresAt)
	assert.NotNil(t, resp.User.CreatedAt)
	assert.NotNil(t, resp.User.UpdatedAt)

	// Verify mock calls
	mockStore.AssertExpectations(t)
}
