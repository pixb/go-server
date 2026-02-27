package auth

import (
	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
	"github.com/pixb/go-server/store"
)

// StringToRole converts a store Role to protobuf Role enum
func StringToRole(role store.Role) v1pb.Role {
	switch role {
	case store.RoleAdmin:
		return v1pb.Role_ROLE_ADMIN
	case store.RoleUser:
		return v1pb.Role_ROLE_USER
	default:
		return v1pb.Role_ROLE_UNSPECIFIED
	}
}

// RoleToString converts protobuf Role enum to a store Role
func RoleToString(role v1pb.Role) store.Role {
	switch role {
	case v1pb.Role_ROLE_ADMIN:
		return store.RoleAdmin
	case v1pb.Role_ROLE_USER:
		return store.RoleUser
	default:
		return ""
	}
}
