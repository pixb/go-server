package store

import "time"

type User struct {
	ID              int64
	Username        string
	Nickname        string
	Password        string
	Phone           string
	Email           string
	Role            string
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
	Role            *string
	PasswordExpires *time.Time
	UpdatedAt       *time.Time
}

type FindUser struct {
	ID       *int64
	Username *string
	Limit    *int
}

type DeleteUser struct {
	ID int64
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
	ID        int64
	Revoked   *bool
	UpdatedAt *time.Time
}

type FindRefreshToken struct {
	ID     *int64
	UserID *int64
	Token  *string
}

type DeleteRefreshToken struct {
	ID int64
}
