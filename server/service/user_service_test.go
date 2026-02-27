package service

import (
	"context"
	"testing"
	"time"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mock implementation of store.Store
type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateUser(ctx context.Context, user *store.User) (*store.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) UpdateUser(ctx context.Context, update *store.UpdateUser) (*store.User, error) {
	args := m.Called(ctx, update)
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error) {
	args := m.Called(ctx, find)
	return args.Get(0).([]*store.User), args.Error(1)
}

func (m *MockStore) DeleteUser(ctx context.Context, delete *store.DeleteUser) error {
	args := m.Called(ctx, delete)
	return args.Error(0)
}

func (m *MockStore) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) CreateRefreshToken(ctx context.Context, create *store.CreateRefreshToken) (*store.RefreshToken, error) {
	args := m.Called(ctx, create)
	return args.Get(0).(*store.RefreshToken), args.Error(1)
}

func (m *MockStore) UpdateRefreshToken(ctx context.Context, update *store.UpdateRefreshToken) (*store.RefreshToken, error) {
	args := m.Called(ctx, update)
	return args.Get(0).(*store.RefreshToken), args.Error(1)
}

func (m *MockStore) ListRefreshTokens(ctx context.Context, find *store.FindRefreshToken) ([]*store.RefreshToken, error) {
	args := m.Called(ctx, find)
	return args.Get(0).([]*store.RefreshToken), args.Error(1)
}

func (m *MockStore) DeleteRefreshToken(ctx context.Context, delete *store.DeleteRefreshToken) error {
	args := m.Called(ctx, delete)
	return args.Error(0)
}

func (m *MockStore) GetRefreshToken(ctx context.Context, token string) (*store.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.RefreshToken), args.Error(1)
}

func (m *MockStore) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStore) GetUser(ctx context.Context, find *store.FindUser) (*store.User, error) {
	args := m.Called(ctx, find)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockStore) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestUserService_RegisterUser(t *testing.T) {
	// Create mock store
	mockStore := new(MockStore)

	// Test data
	req := &v1pb.RegisterUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Nickname: "Test User",
		Phone:    "13800138000",
	}

	// Mock responses
	mockStore.On("GetUserByUsername", mock.Anything, req.Username).Return(nil, nil)
	mockStore.On("GetUserByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockStore.On("CreateUser", mock.Anything, mock.AnythingOfType("*store.User")).Return(&store.User{
		ID:              1,
		Username:        req.Username,
		Email:           req.Email,
		Password:        "hashedpassword",
		Nickname:        req.Nickname,
		Phone:           req.Phone,
		Role:            store.RoleUser,
		PasswordExpires: time.Now().AddDate(0, 0, 90),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil)

	// Create user service
	userService := NewUserService("testsecret", mockStore)

	// Test RegisterUser
	resp, err := userService.RegisterUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Nickname, resp.User.Nickname)
	assert.Equal(t, req.Phone, resp.User.Phone)
	assert.Equal(t, v1pb.Role_ROLE_USER, resp.User.Role)
	assert.NotNil(t, resp.User.PasswordExpiresAt)
	assert.NotNil(t, resp.User.CreatedAt)
	assert.NotNil(t, resp.User.UpdatedAt)

	// Verify mock calls
	mockStore.AssertExpectations(t)
}
