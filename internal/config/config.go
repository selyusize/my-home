package config

import (
	"os"
	"strconv"

	"github.com/selyusize/my-home/pkg/db"
	"github.com/selyusize/my-home/pkg/docker"
	"github.com/selyusize/my-home/pkg/git"
)

func Database() db.Config {
	return db.Config{
		DSN:      os.Getenv("DATABASE_DSN"),
		Host:     stringOr("DATABASE_HOST", "localhost"),
		Port:     intOr("DATABASE_PORT", 5432),
		User:     stringOr("DATABASE_USER", "myhome"),
		Password: stringOr("DATABASE_PASSWORD", "myhome"),
		Name:     stringOr("DATABASE_NAME", "myhome"),
		SSLMode:  stringOr("DATABASE_SSLMODE", "disable"),
		TimeZone: stringOr("DATABASE_TIMEZONE", "UTC"),
	}
}

func GitHub() git.Config {
	return git.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
	}
}

func GitLab() git.Config {
	return git.Config{
		ClientID:     os.Getenv("GITLAB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITLAB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITLAB_REDIRECT_URL"),
		BaseURL:      os.Getenv("GITLAB_BASE_URL"),
	}
}

func Docker() docker.Config {
	return docker.Config{
		Host:       os.Getenv("DOCKER_HOST"),
		APIVersion: os.Getenv("DOCKER_API_VERSION"),
		CertPath:   os.Getenv("DOCKER_CERT_PATH"),
	}
}

func stringOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func intOr(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}
