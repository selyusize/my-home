package gitlab

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/selyusize/my-home/pkg/git"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
	"golang.org/x/oauth2"
	oauth2gitlab "golang.org/x/oauth2/gitlab"
)

var _ git.Client = (*Client)(nil)

const defaultBaseURL = "https://gitlab.com"

// Client is a GitLab API client wrapper.
type Client struct {
	mu     sync.RWMutex
	cfg    git.Config
	oauth  *oauth2.Config
	token  string
	client *gitlab.Client
}

// Register registers the GitLab constructor on a git.Factory.
func Register() git.Option {
	return git.WithProvider(git.ProviderGitLab, New)
}

// New creates a GitLab client. OAuth fields are optional until AuthURL/Exchange are used.
func New(cfg git.Config) (git.Client, error) {
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"read_user", "read_api"}
	}
	cfg.BaseURL = normalizeInstanceURL(cfg.BaseURL)

	c := &Client{cfg: cfg}
	if cfg.ClientID != "" {
		c.oauth = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     oauthEndpoint(cfg.BaseURL),
		}
	}
	return c, nil
}

// Provider returns the GitLab provider identifier.
func (c *Client) Provider() git.Provider {
	return git.ProviderGitLab
}

// IsAuthenticated reports whether an access token is set.
func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token != ""
}

// AccessToken returns the current access token.
func (c *Client) AccessToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// Logout clears the current session.
func (c *Client) Logout() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
	c.client = nil
}

func (c *Client) api(_ context.Context) (*gitlab.Client, error) {
	c.mu.RLock()
	client := c.client
	token := c.token
	c.mu.RUnlock()

	if token == "" {
		return nil, git.ErrNotAuthenticated
	}
	if client != nil {
		return client, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client != nil {
		return c.client, nil
	}

	opts := make([]gitlab.ClientOptionFunc, 0, 1)
	if baseURL := strings.TrimSpace(c.cfg.BaseURL); baseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(baseURL))
	}

	client, err := gitlab.NewClient(c.token, opts...)
	if err != nil {
		return nil, fmt.Errorf("gitlab: create client: %w", err)
	}
	c.client = client
	return c.client, nil
}

func oauthEndpoint(baseURL string) oauth2.Endpoint {
	base := normalizeInstanceURL(baseURL)
	if base == "" || base == defaultBaseURL {
		return oauth2gitlab.Endpoint
	}
	return oauth2.Endpoint{
		AuthURL:  base + "/oauth/authorize",
		TokenURL: base + "/oauth/token",
	}
}

func normalizeInstanceURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimRight(raw, "/")
	raw = strings.TrimSuffix(raw, "/api/v4")
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	return raw
}

func listOpts(opts git.ListOptions) gitlab.ListOptions {
	opts = git.NormalizeListOptions(opts)
	return gitlab.ListOptions{
		Page:    int64(opts.Page),
		PerPage: int64(opts.PerPage),
	}
}

func pageInfoFrom(resp *gitlab.Response) git.PageInfo {
	if resp == nil {
		return git.PageInfo{}
	}
	firstPage := 0
	if resp.TotalPages > 0 || resp.CurrentPage > 0 || resp.NextPage > 0 {
		firstPage = 1
	}
	return git.PageInfo{
		NextPage:  int(resp.NextPage),
		PrevPage:  int(resp.PreviousPage),
		FirstPage: firstPage,
		LastPage:  int(resp.TotalPages),
	}
}
