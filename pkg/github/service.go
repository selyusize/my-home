package github

import (
	"context"
	"fmt"
	"sync"

	gh "github.com/google/go-github/v90/github"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

// Service is a GitHub API client wrapper.
type Service struct {
	mu     sync.RWMutex
	cfg    Config
	oauth  *oauth2.Config
	token  string
	client *gh.Client
}

// New creates a GitHub service. OAuth fields are optional until AuthURL/Exchange are used.
func New(cfg Config) *Service {
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"read:user", "user:email", "repo"}
	}

	s := &Service{cfg: cfg}
	if cfg.ClientID != "" {
		s.oauth = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     oauth2github.Endpoint,
		}
	}
	return s
}

// IsAuthenticated reports whether an access token is set.
func (s *Service) IsAuthenticated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token != ""
}

// AccessToken returns the current access token.
func (s *Service) AccessToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

// Logout clears the current session.
func (s *Service) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
	s.client = nil
}

func (s *Service) api(ctx context.Context) (*gh.Client, error) {
	s.mu.RLock()
	client := s.client
	token := s.token
	s.mu.RUnlock()

	if token == "" {
		return nil, ErrNotAuthenticated
	}
	if client != nil {
		return client, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		return s.client, nil
	}

	client, err := gh.NewClient(gh.WithAuthToken(s.token))
	if err != nil {
		return nil, fmt.Errorf("github: create client: %w", err)
	}
	s.client = client
	_ = ctx
	return s.client, nil
}

func listOpts(opts ListOptions) *gh.ListOptions {
	page := opts.Page
	perPage := opts.PerPage
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 30
	}
	if perPage > 100 {
		perPage = 100
	}
	return &gh.ListOptions{Page: page, PerPage: perPage}
}

func pageInfoFrom(resp *gh.Response) PageInfo {
	if resp == nil {
		return PageInfo{}
	}
	return PageInfo{
		NextPage:  resp.NextPage,
		PrevPage:  resp.PrevPage,
		FirstPage: resp.FirstPage,
		LastPage:  resp.LastPage,
	}
}
