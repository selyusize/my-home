package db

import (
	"context"

	"gorm.io/gorm"
)

// Client is the common API implemented by every database provider.
type Client interface {
	Provider() Provider
	DB() *gorm.DB
	Ping(ctx context.Context) error
	AutoMigrate(dst ...any) error
	Close() error
}

// Constructor creates a Client for a registered provider.
type Constructor func(cfg Config) (Client, error)
