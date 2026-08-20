package services

import (
	"context"
	"os"

	ghsvc "github.com/selyusize/my-home/pkg/github"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GitHubService exposes the GitHub API to the Wails frontend.
type GitHubService struct {
	*ghsvc.Service
}

// NewGitHubService creates a GitHub service configured from environment variables:
//   - GITHUB_CLIENT_ID
//   - GITHUB_CLIENT_SECRET
//   - GITHUB_REDIRECT_URL
//   - GITHUB_TOKEN (optional, applied on startup)
func NewGitHubService() *GitHubService {
	return &GitHubService{
		Service: ghsvc.New(ghsvc.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		}),
	}
}

// ServiceStartup optionally authenticates with GITHUB_TOKEN.
func (s *GitHubService) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil
	}
	return s.SetAccessToken(token)
}
