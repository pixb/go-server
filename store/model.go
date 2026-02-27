package store

import "time"

// Role is the type of a role.
type Role string

// Role constants
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// IsValidRole checks if the role is valid
func IsValidRole(role Role) bool {
	return role == RoleAdmin || role == RoleUser
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

type User struct {
	ID              int64
	Username        string
	Nickname        string
	Password        string
	Phone           string
	Email           string
	Role            Role
	PasswordExpires time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type UpdateUser struct {
	ID              int64
	Username        *string
	Nickname        *string
	Password        *string
	Phone           *string
	Email           *string
	Role            *Role
	PasswordExpires *time.Time
	UpdatedAt       *time.Time
}

type CreateUser struct {
	Username        string
	Nickname        string
	Password        string
	Phone           string
	Email           string
	Role            Role
	PasswordExpires time.Time
}

type FindUser struct {
	ID       *int64
	Username *string
	Email    *string
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CreateRefreshToken struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

type UpdateRefreshToken struct {
	ID      int64
	Revoked *bool
}

type DeleteUser struct {
	ID int64
}

type FindRefreshToken struct {
	ID     *int64
	UserID *int64
	Token  *string
}

type DeleteRefreshToken struct {
	ID int64
}
