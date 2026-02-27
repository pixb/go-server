package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/pixb/go-server/internal/profile"
)

type Driver struct {
	db      *sql.DB
	profile *profile.Profile
}

func NewDriver(profile *profile.Profile) (*Driver, error) {
	db, err := sql.Open("postgres", profile.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Driver{
		db:      db,
		profile: profile,
	}, nil
}

func (d *Driver) GetDB() *sql.DB { return d.db }

func (d *Driver) Close() error { return d.db.Close() }

func (d *Driver) IsInitialized(ctx context.Context) (bool, error) {
	var count int
	err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Driver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
