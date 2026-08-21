package git

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/internal/config"
	pkggit "github.com/selyusize/my-home/pkg/git"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GitLabService struct {
	pkggit.Client
	creds       *CredentialStore
	redirectURL string
}

func NewGitLabService(creds *CredentialStore) (*GitLabService, error) {
	cfg := config.GitLab()
	cfg.RedirectURL = loopbackRedirect(cfg.RedirectURL)
	client, err := factory.New(pkggit.ProviderGitLab, cfg)
	if err != nil {
		return nil, err
	}
	return &GitLabService{Client: client, creds: creds, redirectURL: cfg.RedirectURL}, nil
}

func (s *GitLabService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	return restoreGitSession(ctx, s.Client, s.creds, pkggit.ProviderGitLab)
}

func (s *GitLabService) Login(ctx context.Context) (*pkggit.User, error) {
	user, err := loginWithBrowser(ctx, s.Client, s.redirectURL, "GitLab")
	if err != nil {
		return nil, err
	}
	if err := s.creds.Save(pkggit.ProviderGitLab, s.AccessToken()); err != nil {
		return nil, fmt.Errorf("gitlab: persist session: %w", err)
	}
	return user, nil
}

func (s *GitLabService) SetAccessToken(token string) error {
	return persistGitToken(s.Client, s.creds, pkggit.ProviderGitLab, token)
}

func (s *GitLabService) Logout() error {
	return clearGitToken(s.Client, s.creds, pkggit.ProviderGitLab)
}
