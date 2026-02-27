package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pixb/go-server/store"
)

func (d *Driver) CreateRefreshToken(ctx context.Context, create *store.CreateRefreshToken) (*store.RefreshToken, error) {
	now := time.Now()
	result, err := d.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		create.UserID, create.Token, create.ExpiresAt, false, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &store.RefreshToken{
		ID:        id,
		UserID:    create.UserID,
		Token:     create.Token,
		ExpiresAt: create.ExpiresAt,
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (d *Driver) UpdateRefreshToken(ctx context.Context, update *store.UpdateRefreshToken) (*store.RefreshToken, error) {
	query := "UPDATE refresh_tokens SET updated_at = ?"
	args := []interface{}{time.Now()}

	if update.Revoked != nil {
		query += ", revoked = ?"
		args = append(args, *update.Revoked)
	}

	query += " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, update.ID)

	_, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	// Return updated token by querying
	return d.GetRefreshTokenByID(ctx, update.ID)
}

func (d *Driver) ListRefreshTokens(ctx context.Context, find *store.FindRefreshToken) ([]*store.RefreshToken, error) {
	query := "SELECT id, user_id, token, expires_at, revoked, created_at, updated_at, deleted_at FROM refresh_tokens WHERE deleted_at IS NULL"
	args := []interface{}{}

	if find.ID != nil {
		query += " AND id = ?"
		args = append(args, *find.ID)
	}
	if find.UserID != nil {
		query += " AND user_id = ?"
		args = append(args, *find.UserID)
	}
	if find.Token != nil {
		query += " AND token = ?"
		args = append(args, *find.Token)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list refresh tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*store.RefreshToken
	for rows.Next() {
		var token store.RefreshToken
		var deletedAt *time.Time
		if err := rows.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.Revoked, &token.CreatedAt, &token.UpdatedAt, &deletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan refresh token: %w", err)
		}
		token.DeletedAt = deletedAt
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

func (d *Driver) DeleteRefreshToken(ctx context.Context, delete *store.DeleteRefreshToken) error {
	_, err := d.db.ExecContext(ctx, "UPDATE refresh_tokens SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", time.Now(), delete.ID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (d *Driver) GetRefreshToken(ctx context.Context, token string) (*store.RefreshToken, error) {
	var refreshToken store.RefreshToken
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		"SELECT id, user_id, token, expires_at, revoked, created_at, updated_at, deleted_at FROM refresh_tokens WHERE token = ? AND deleted_at IS NULL",
		token).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.Revoked, &refreshToken.CreatedAt, &refreshToken.UpdatedAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	refreshToken.DeletedAt = deletedAt
	return &refreshToken, nil
}

func (d *Driver) GetRefreshTokenByID(ctx context.Context, id int64) (*store.RefreshToken, error) {
	var refreshToken store.RefreshToken
	var deletedAt *time.Time
	err := d.db.QueryRowContext(ctx,
		"SELECT id, user_id, token, expires_at, revoked, created_at, updated_at, deleted_at FROM refresh_tokens WHERE id = ? AND deleted_at IS NULL",
		id).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.Revoked, &refreshToken.CreatedAt, &refreshToken.UpdatedAt, &deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get refresh token by id: %w", err)
	}
	refreshToken.DeletedAt = deletedAt
	return &refreshToken, nil
}
