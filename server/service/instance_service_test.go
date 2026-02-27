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

func TestInstanceService_GetInstanceProfile(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		demo          bool
		adminUsers    []*store.User
		expectedAdmin bool
		expectedError bool
	}{
		{
			name:    "instance with admin user",
			version: "1.0.0",
			demo:    false,
			adminUsers: []*store.User{
				{
					ID:              1,
					Username:        "admin",
					Email:           "admin@example.com",
					Password:        "hashedpassword",
					Nickname:        "Admin User",
					Phone:           "13800138000",
					Role:            store.RoleAdmin,
					PasswordExpires: time.Now().AddDate(0, 0, 90),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedAdmin: true,
			expectedError: false,
		},
		{
			name:          "instance with only regular user",
			version:       "1.0.0",
			demo:          false,
			adminUsers:    []*store.User{},
			expectedAdmin: false,
			expectedError: false,
		},
		{
			name:          "empty instance - no users",
			version:       "0.1.0",
			demo:          true,
			adminUsers:    []*store.User{},
			expectedAdmin: false,
			expectedError: false,
		},
		{
			name:    "demo instance with admin",
			version: "0.1.0-demo",
			demo:    true,
			adminUsers: []*store.User{
				{
					ID:              1,
					Username:        "admin",
					Email:           "admin@example.com",
					Password:        "hashedpassword",
					Nickname:        "Admin User",
					Phone:           "13800138000",
					Role:            store.RoleAdmin,
					PasswordExpires: time.Now().AddDate(0, 0, 90),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedAdmin: true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock store
			mockStore := new(MockStore)
			mockStore.On("ListUsers", mock.Anything, mock.AnythingOfType("*store.FindUser")).Return(tt.adminUsers, nil)

			// Create instance service
			instanceService := NewInstanceService(tt.version, tt.demo, mockStore)

			// Test GetInstanceProfile
			resp, err := instanceService.GetInstanceProfile(context.Background(), &v1pb.GetInstanceProfileRequest{})

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.version, resp.Version)
			assert.Equal(t, tt.demo, resp.Demo)

			if tt.expectedAdmin {
				assert.NotNil(t, resp.Admin)
				assert.Equal(t, v1pb.Role_ROLE_ADMIN, resp.Admin.Role)
			} else {
				assert.Nil(t, resp.Admin)
			}

			// Verify mock calls
			mockStore.AssertExpectations(t)
		})
	}
}
