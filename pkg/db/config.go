package db

import (
	"fmt"
	"net/url"
	"strings"
)

// Config holds connection settings shared by database providers.
type Config struct {
	DSN          string
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	TimeZone     string
	MaxOpenConns int
	MaxIdleConns int
}

// PostgresDSN builds a libpq connection string.
func (c Config) PostgresDSN() string {
	if dsn := strings.TrimSpace(c.DSN); dsn != "" {
		return dsn
	}

	port := c.Port
	if port <= 0 {
		port = 5432
	}
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	timeZone := c.TimeZone
	if timeZone == "" {
		timeZone = "UTC"
	}

	u := url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", c.Host, port),
		Path:   c.Name,
	}
	if c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	} else if c.User != "" {
		u.User = url.User(c.User)
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	q.Set("TimeZone", timeZone)
	u.RawQuery = q.Encode()
	return u.String()
}

// Validate reports whether the config has enough data to connect.
func (c Config) Validate() error {
	if strings.TrimSpace(c.DSN) != "" {
		return nil
	}
	if strings.TrimSpace(c.Host) == "" || strings.TrimSpace(c.User) == "" || strings.TrimSpace(c.Name) == "" {
		return ErrMissingConfig
	}
	return nil
}
