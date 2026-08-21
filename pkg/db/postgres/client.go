package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/selyusize/my-home/pkg/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ db.Client = (*Client)(nil)

const (
	defaultMaxOpenConns = 20
	defaultMaxIdleConns = 5
	defaultConnLifetime = time.Hour
)

// Client is a GORM wrapper for PostgreSQL.
type Client struct {
	gdb *gorm.DB
}

// Register registers the PostgreSQL constructor on a db.Factory.
func Register() db.Option {
	return db.WithProvider(db.ProviderPostgres, New)
}

// New opens a PostgreSQL connection with GORM.
func New(cfg db.Config) (db.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres: sql db: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = defaultMaxOpenConns
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = defaultMaxIdleConns
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(defaultConnLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	return &Client{gdb: gdb}, nil
}

// Provider returns the PostgreSQL provider identifier.
func (c *Client) Provider() db.Provider {
	return db.ProviderPostgres
}

// DB returns the underlying GORM handle.
func (c *Client) DB() *gorm.DB {
	return c.gdb
}

// Ping checks that PostgreSQL is reachable.
func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.gdb.DB()
	if err != nil {
		return fmt.Errorf("postgres: sql db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres: ping: %w", err)
	}
	return nil
}

// AutoMigrate runs GORM auto-migration for the given models.
func (c *Client) AutoMigrate(dst ...any) error {
	if err := c.gdb.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("postgres: auto migrate: %w", err)
	}
	return nil
}

// Close closes the underlying SQL connection pool.
func (c *Client) Close() error {
	sqlDB, err := c.gdb.DB()
	if err != nil {
		return fmt.Errorf("postgres: sql db: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("postgres: close: %w", err)
	}
	return nil
}
