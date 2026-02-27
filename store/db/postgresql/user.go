package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pixb/go-server/store"
)

func (d *Driver) CreateUser(ctx context.Context, create *store.User) (*store.User, error) {
	var id int64
	now := time.Now()
	passwordExpires := now.AddDate(0, 0, 90) // Default 90 days expiration

	err := d.db.QueryRowContext(ctx,
		`INSERT INTO users (username, nickname, password, phone, email, role, password_expires, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		create.Username, create.Nickname, create.Password, create.Phone, create.Email, create.Role, passwordExpires, now, now).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
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
	query := `UPDATE users SET updated_at = $1`
	args := []interface{}{time.Now()}
	argCount := 1

	if update.Username != nil {
		argCount++
		query += fmt.Sprintf(", username = $%d", argCount)
		args = append(args, *update.Username)
	}
	if update.Nickname != nil {
		argCount++
		query += fmt.Sprintf(", nickname = $%d", argCount)
		args = append(args, *update.Nickname)
	}
	if update.Password != nil {
		argCount++
		query += fmt.Sprintf(", password = $%d", argCount)
		args = append(args, *update.Password)
		// Update password expiration when password is changed
		argCount++
		query += fmt.Sprintf(", password_expires = $%d", argCount)
		args = append(args, time.Now().AddDate(0, 0, 90))
	}
	if update.Phone != nil {
		argCount++
		query += fmt.Sprintf(", phone = $%d", argCount)
		args = append(args, *update.Phone)
	}
	if update.Email != nil {
		argCount++
		query += fmt.Sprintf(", email = $%d", argCount)
		args = append(args, *update.Email)
	}
	if update.Role != nil {
		argCount++
		query += fmt.Sprintf(", role = $%d", argCount)
		args = append(args, *update.Role)
	}
	if update.PasswordExpires != nil {
		argCount++
		query += fmt.Sprintf(", password_expires = $%d", argCount)
		args = append(args, *update.PasswordExpires)
	}

	argCount++
	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL RETURNING id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at", argCount)
	args = append(args, update.ID)

	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Password, &user.Phone, &user.Email, &user.Role, &user.PasswordExpires, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	user.DeletedAt = deletedAt

	return &user, nil
}

func (d *Driver) ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error) {
	query := `SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}

	if find.ID != nil {
		args = append(args, *find.ID)
		query += fmt.Sprintf(" AND id = $%d", len(args))
	}
	if find.Username != nil {
		args = append(args, *find.Username)
		query += fmt.Sprintf(" AND username = $%d", len(args))
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
	_, err := d.db.ExecContext(ctx, `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`, time.Now(), delete.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (d *Driver) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		`SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE username = $1 AND deleted_at IS NULL`,
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
		`SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`,
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

func (d *Driver) GetUserByID(ctx context.Context, id int64) (*store.User, error) {
	var user store.User
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		`SELECT id, username, nickname, password, phone, email, role, password_expires, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`,
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
