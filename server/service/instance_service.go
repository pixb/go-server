package service

import (
	"context"

	"connectrpc.com/connect"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/server/auth"
	"github.com/pixb/go-server/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// InstanceStore is an interface that defines the methods needed by InstanceService
type InstanceStore interface {
	ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error)
	Ping(ctx context.Context) error
	Close() error
}

type InstanceService struct {
	Version string
	Demo    bool
	Store   InstanceStore
}

func NewInstanceService(version string, demo bool, store InstanceStore) *InstanceService {
	return &InstanceService{
		Version: version,
		Demo:    demo,
		Store:   store,
	}
}

func (s *InstanceService) GetInstanceProfile(ctx context.Context, req *v1pb.GetInstanceProfileRequest) (*v1pb.InstanceProfile, error) {
	profile := &v1pb.InstanceProfile{
		Version: s.Version,
		Demo:    s.Demo,
	}

	// Try to find the first admin user
	adminRole := store.RoleAdmin
	users, err := s.Store.ListUsers(ctx, &store.FindUser{
		Role: &adminRole,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Use the first admin user if found
	if len(users) > 0 {
		user := users[0]
		profile.Admin = &v1pb.User{
			Id:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			Nickname:          user.Nickname,
			Phone:             user.Phone,
			Role:              auth.StringToRole(user.Role),
			PasswordExpiresAt: timestamppb.New(user.PasswordExpires),
			CreatedAt:         timestamppb.New(user.CreatedAt),
			UpdatedAt:         timestamppb.New(user.UpdatedAt),
		}
	}

	return profile, nil
}
