package db

import "errors"

var (
	ErrUnknownProvider = errors.New("db: unknown provider")
	ErrMissingConfig   = errors.New("db: host, user, and name or dsn are required")
)
