package auth

import v1pb "github.com/pixb/go-server/proto/gen/api/v1"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleUser
}

func StringToRole(role string) v1pb.Role {
	switch role {
	case RoleAdmin:
		return v1pb.Role_ROLE_ADMIN
	case RoleUser:
		return v1pb.Role_ROLE_USER
	default:
		return v1pb.Role_ROLE_UNSPECIFIED
	}
}

func RoleToString(role v1pb.Role) string {
	switch role {
	case v1pb.Role_ROLE_ADMIN:
		return RoleAdmin
	case v1pb.Role_ROLE_USER:
		return RoleUser
	default:
		return ""
	}
}
