package store

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/pixb/go-server/internal/profile"
	"github.com/pixb/go-server/store/cache"
)

type Driver interface {
	GetDB() *sql.DB
	Close() error
	IsInitialized(ctx context.Context) (bool, error)
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, create *User) (*User, error)
	UpdateUser(ctx context.Context, update *UpdateUser) (*User, error)
	ListUsers(ctx context.Context, find *FindUser) ([]*User, error)
	DeleteUser(ctx context.Context, delete *DeleteUser) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateRefreshToken(ctx context.Context, create *CreateRefreshToken) (*RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, update *UpdateRefreshToken) (*RefreshToken, error)
	ListRefreshTokens(ctx context.Context, find *FindRefreshToken) ([]*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, delete *DeleteRefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)

	// InstanceSetting model related methods.
	UpsertInstanceSetting(ctx context.Context, upsert *InstanceSetting) (*InstanceSetting, error)
	ListInstanceSettings(ctx context.Context, find *FindInstanceSetting) ([]*InstanceSetting, error)
	DeleteInstanceSetting(ctx context.Context, delete *DeleteInstanceSetting) error
}

type Store struct {
	driver  Driver
	profile *profile.Profile

	cacheConfig          *cache.Config
	userCache            *cache.Cache
	instanceSettingCache *cache.Cache
}

func New(driver Driver, profile *profile.Profile) *Store {
	cacheConfig := &cache.Config{
		DefaultTTL:      10 * time.Minute,
		CleanupInterval: 5 * time.Minute,
		MaxItems:        1000,
	}
	return &Store{
		driver:               driver,
		profile:              profile,
		cacheConfig:          cacheConfig,
		userCache:            cache.New(*cacheConfig),
		instanceSettingCache: cache.New(*cacheConfig),
	}
}

func (s *Store) GetDriver() Driver { return s.driver }

func (s *Store) Close() error {
	s.userCache.Close()
	return s.driver.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.driver.Ping(ctx)
}

func (s *Store) CreateUser(ctx context.Context, create *User) (*User, error) {
	user, err := s.driver.CreateUser(ctx, create)
	if err != nil {
		return nil, err
	}
	s.userCache.Set(ctx, strconv.FormatInt(user.ID, 10), user)
	return user, nil
}

func (s *Store) GetUser(ctx context.Context, find *FindUser) (*User, error) {
	if find.ID != nil {
		if cached, ok := s.userCache.Get(ctx, strconv.FormatInt(*find.ID, 10)); ok {
			if user, ok := cached.(*User); ok {
				return user, nil
			}
		}
	}
	list, err := s.ListUsers(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, sql.ErrNoRows
	}
	return list[0], nil
}

func (s *Store) ListUsers(ctx context.Context, find *FindUser) ([]*User, error) {
	return s.driver.ListUsers(ctx, find)
}

func (s *Store) UpdateUser(ctx context.Context, update *UpdateUser) (*User, error) {
	user, err := s.driver.UpdateUser(ctx, update)
	if err != nil {
		return nil, err
	}
	s.userCache.Delete(ctx, strconv.FormatInt(user.ID, 10))
	return user, nil
}

func (s *Store) DeleteUser(ctx context.Context, delete *DeleteUser) error {
	if err := s.driver.DeleteUser(ctx, delete); err != nil {
		return err
	}
	s.userCache.Delete(ctx, strconv.FormatInt(delete.ID, 10))
	return nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return s.driver.GetUserByUsername(ctx, username)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.driver.GetUserByEmail(ctx, email)
}

func (s *Store) CreateRefreshToken(ctx context.Context, create *CreateRefreshToken) (*RefreshToken, error) {
	return s.driver.CreateRefreshToken(ctx, create)
}

func (s *Store) UpdateRefreshToken(ctx context.Context, update *UpdateRefreshToken) (*RefreshToken, error) {
	return s.driver.UpdateRefreshToken(ctx, update)
}

func (s *Store) ListRefreshTokens(ctx context.Context, find *FindRefreshToken) ([]*RefreshToken, error) {
	return s.driver.ListRefreshTokens(ctx, find)
}

func (s *Store) DeleteRefreshToken(ctx context.Context, delete *DeleteRefreshToken) error {
	return s.driver.DeleteRefreshToken(ctx, delete)
}

func (s *Store) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	return s.driver.GetRefreshToken(ctx, token)
}
