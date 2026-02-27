package service

import (
	"context"
	"errors"
	"regexp"

	"connectrpc.com/connect"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserStore is an interface that defines the methods needed by UserService
type UserStore interface {
	CreateUser(ctx context.Context, create *store.User) (*store.User, error)
	UpdateUser(ctx context.Context, update *store.UpdateUser) (*store.User, error)
	ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error)
	DeleteUser(ctx context.Context, delete *store.DeleteUser) error
	GetUserByUsername(ctx context.Context, username string) (*store.User, error)
	GetUserByEmail(ctx context.Context, email string) (*store.User, error)
	GetUser(ctx context.Context, find *store.FindUser) (*store.User, error)
	Ping(ctx context.Context) error
	Close() error
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

type UserService struct {
	Secret string
	Store  UserStore
}

func NewUserService(secret string, store UserStore) *UserService {
	return &UserService{
		Secret: secret,
		Store:  store,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *v1pb.RegisterUserRequest) (*v1pb.RegisterUserResponse, error) {
	// Validate username
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username must be between 3 and 50 characters"))
	}

	// Validate nickname
	if req.Nickname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("nickname is required"))
	}
	if len(req.Nickname) > 50 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("nickname must be at most 50 characters"))
	}

	// Validate password
	if req.Password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password is required"))
	}
	if len(req.Password) < 6 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password must be at least 6 characters"))
	}

	// Validate phone
	if req.Phone == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("phone is required"))
	}
	if len(req.Phone) != 11 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("phone must be 11 digits"))
	}
	for _, r := range req.Phone {
		if r < '0' || r > '9' {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("phone must be only digits"))
		}
	}

	// Validate email
	if req.Email == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email is required"))
	}
	if len(req.Email) > 100 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email must be at most 100 characters"))
	}
	// Simple email format validation
	if !isValidEmail(req.Email) {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email must be a valid email address"))
	}

	// Check if username already exists
	existingUser, err := s.Store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("username already exists"))
	}

	// Check if email already exists
	existingUserByEmail, err := s.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUserByEmail != nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("email already exists"))
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to hash password"))
	}

	newUser, err := s.Store.CreateUser(ctx, &store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: passwordHash,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     store.RoleUser, // Default role
	})
	if err != nil {
		return nil, err
	}

	return &v1pb.RegisterUserResponse{
		User: &v1pb.User{
			Id:                newUser.ID,
			Username:          newUser.Username,
			Email:             newUser.Email,
			Nickname:          newUser.Nickname,
			Phone:             newUser.Phone,
			Role:              auth.StringToRole(newUser.Role),
			PasswordExpiresAt: timestamppb.New(newUser.PasswordExpires),
			CreatedAt:         timestamppb.New(newUser.CreatedAt),
			UpdatedAt:         timestamppb.New(newUser.UpdatedAt),
		},
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *v1pb.GetUserProfileRequest) (*v1pb.GetUserProfileResponse, error) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	return &v1pb.GetUserProfileResponse{
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

func (s *UserService) UpdateUserProfile(ctx context.Context, req *v1pb.UpdateUserProfileRequest) (*v1pb.UpdateUserProfileResponse, error) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	// Validate parameters
	if req.Nickname != "" && len(req.Nickname) > 50 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("nickname must be at most 50 characters"))
	}

	if req.Phone != "" {
		if len(req.Phone) != 11 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("phone must be 11 digits"))
		}
		for _, r := range req.Phone {
			if r < '0' || r > '9' {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("phone must be only digits"))
			}
		}
	}

	if req.Email != "" {
		if len(req.Email) > 100 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email must be at most 100 characters"))
		}
		if !isValidEmail(req.Email) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email must be a valid email address"))
		}
	}

	// Build UpdateUser struct
	update := &store.UpdateUser{
		ID: userID,
	}

	if req.Nickname != "" {
		update.Nickname = &req.Nickname
	}

	if req.Phone != "" {
		update.Phone = &req.Phone
	}

	if req.Email != "" {
		update.Email = &req.Email
	}

	// Update user
	updatedUser, err := s.Store.UpdateUser(ctx, update)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to update user profile"))
	}

	return &v1pb.UpdateUserProfileResponse{
		User: &v1pb.User{
			Id:                updatedUser.ID,
			Username:          updatedUser.Username,
			Email:             updatedUser.Email,
			Nickname:          updatedUser.Nickname,
			Phone:             updatedUser.Phone,
			Role:              auth.StringToRole(updatedUser.Role),
			PasswordExpiresAt: timestamppb.New(updatedUser.PasswordExpires),
			CreatedAt:         timestamppb.New(updatedUser.CreatedAt),
			UpdatedAt:         timestamppb.New(updatedUser.UpdatedAt),
		},
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *v1pb.ChangePasswordRequest) (*v1pb.ChangePasswordResponse, error) {
	userID := auth.GetUserID(ctx)
	if userID == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("authentication required"))
	}

	// Validate parameters
	if req.OldPassword == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("old password is required"))
	}

	if req.NewPassword == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("new password is required"))
	}

	if len(req.NewPassword) < 6 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("new password must be at least 6 characters"))
	}

	// Get user
	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get user"))
	}
	if user == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
	}

	// Verify old password
	if !auth.CheckPassword(req.OldPassword, user.Password) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("old password is incorrect"))
	}

	// Hash new password
	newPasswordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to hash new password"))
	}

	// Update password
	update := &store.UpdateUser{
		ID:       userID,
		Password: &newPasswordHash,
	}

	updatedUser, err := s.Store.UpdateUser(ctx, update)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to update password"))
	}

	return &v1pb.ChangePasswordResponse{
		User: &v1pb.User{
			Id:                updatedUser.ID,
			Username:          updatedUser.Username,
			Email:             updatedUser.Email,
			Nickname:          updatedUser.Nickname,
			Phone:             updatedUser.Phone,
			Role:              auth.StringToRole(updatedUser.Role),
			PasswordExpiresAt: timestamppb.New(updatedUser.PasswordExpires),
			CreatedAt:         timestamppb.New(updatedUser.CreatedAt),
			UpdatedAt:         timestamppb.New(updatedUser.UpdatedAt),
		},
	}, nil
}
