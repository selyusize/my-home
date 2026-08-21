package git

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/internal/config"
	pkggit "github.com/selyusize/my-home/pkg/git"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GitHubService struct {
	pkggit.Client
	creds       *CredentialStore
	redirectURL string
}

func NewGitHubService(creds *CredentialStore) (*GitHubService, error) {
	cfg := config.GitHub()
	cfg.RedirectURL = loopbackRedirect(cfg.RedirectURL)
	client, err := factory.New(pkggit.ProviderGitHub, cfg)
	if err != nil {
		return nil, err
	}
	return &GitHubService{Client: client, creds: creds, redirectURL: cfg.RedirectURL}, nil
}

func (s *GitHubService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	return restoreGitSession(ctx, s.Client, s.creds, pkggit.ProviderGitHub)
}

func (s *GitHubService) Login(ctx context.Context) (*pkggit.User, error) {
	user, err := loginWithBrowser(ctx, s.Client, s.redirectURL, "GitHub")
	if err != nil {
		return nil, err
	}
	if err := s.creds.Save(pkggit.ProviderGitHub, s.AccessToken()); err != nil {
		return nil, fmt.Errorf("github: persist session: %w", err)
	}
	return user, nil
}

func (s *GitHubService) SetAccessToken(token string) error {
	return persistGitToken(s.Client, s.creds, pkggit.ProviderGitHub, token)
}

func (s *GitHubService) Logout() error {
	return clearGitToken(s.Client, s.creds, pkggit.ProviderGitHub)
}
