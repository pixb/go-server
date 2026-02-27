package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pixb/go-server/store"
)

func (d *Driver) CreateUser(ctx context.Context, create *store.User) (*store.User, error) {
	now := time.Now()
	passwordExpires := now.AddDate(0, 0, 90) // Default 90 days expiration

	result, err := d.db.ExecContext(ctx,
		`INSERT INTO users (username, nickname, password, phone, email, role, password_expires, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		create.Username, create.Nickname, create.Password, create.Phone, create.Email, create.Role, passwordExpires, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &store.User{
		ID:              id,
		Username:        create.Username,
		Nickname:        create.Nickname,
		Password:        create.Password,
		Phone:           create.Phone,
		Email:           create.Email,
		Role:            create.Role,
		PasswordExpires: passwordExpires,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (d *Driver) UpdateUser(ctx context.Context, update *store.UpdateUser) (*store.User, error) {
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{time.Now()}

	if update.Username != nil {
		query += ", username = ?"
		args = append(args, *update.Username)
	}
	if update.Nickname != nil {
		query += ", nickname = ?"
		args = append(args, *update.Nickname)
	}
	if update.Password != nil {
		query += ", password = ?"
		args = append(args, *update.Password)
		// Update password expiration when password is changed
		query += ", password_expires = ?"
		args = append(args, time.Now().AddDate(0, 0, 90))
	}
	if update.Phone != nil {
		query += ", phone = ?"
		args = append(args, *update.Phone)
	}
	if update.Email != nil {
		query += ", email = ?"
		args = append(args, *update.Email)
	}
	if update.Role != nil {
		query += ", role = ?"
		args = append(args, *update.Role)
	}
	if update.PasswordExpires != nil {
		query += ", password_expires = ?"
		args = append(args, *update.PasswordExpires)
	}

	query += " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, update.ID)

	_, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return updated user by querying
	return d.GetUserByID(ctx, update.ID)
}

func (d *Driver) ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error) {
	query := "SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL"
	args := []interface{}{}

	if find.ID != nil {
		query += " AND id = ?"
		args = append(args, *find.ID)
	}
	if find.Username != nil {
		query += " AND username = ?"
		args = append(args, *find.Username)
	}
	if find.Email != nil {
		query += " AND email = ?"
		args = append(args, *find.Email)
	}
	if find.Role != nil {
		query += " AND role = ?"
		args = append(args, *find.Role)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*store.User
	for rows.Next() {
		var user store.User
		var deletedAt *time.Time
		if err := rows.Scan(&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Phone, &user.Email, &user.Role, &user.PasswordExpires, &user.CreatedAt, &user.UpdatedAt, &deletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.DeletedAt = deletedAt
		users = append(users, &user)
	}

	return users, nil
}

func (d *Driver) DeleteUser(ctx context.Context, delete *store.DeleteUser) error {
	_, err := d.db.ExecContext(ctx, "UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", time.Now(), delete.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (d *Driver) GetUserByID(ctx context.Context, id int64) (*store.User, error) {
	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		"SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE id = ? AND deleted_at IS NULL",
		id).Scan(&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Phone, &user.Email, &user.Role, &user.PasswordExpires, &user.CreatedAt, &user.UpdatedAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	user.DeletedAt = deletedAt
	return &user, nil
}

func (d *Driver) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		"SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE username = ? AND deleted_at IS NULL",
		username).Scan(&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Phone, &user.Email, &user.Role, &user.PasswordExpires, &user.CreatedAt, &user.UpdatedAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	user.DeletedAt = deletedAt
	return &user, nil
}

func (d *Driver) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		"SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE email = ? AND deleted_at IS NULL",
		email).Scan(&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Phone, &user.Email, &user.Role, &user.PasswordExpires, &user.CreatedAt, &user.UpdatedAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	user.DeletedAt = deletedAt
	return &user, nil
}
