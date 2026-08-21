package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/selyusize/my-home/pkg/git"

	gh "github.com/google/go-github/v90/github"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

var _ git.Client = (*Client)(nil)

// Client is a GitHub API client wrapper.
type Client struct {
	mu     sync.RWMutex
	cfg    git.Config
	oauth  *oauth2.Config
	token  string
	client *gh.Client
}

// Register registers the GitHub constructor on a git.Factory.
func Register() git.Option {
	return git.WithProvider(git.ProviderGitHub, New)
}

// New creates a GitHub client. OAuth fields are optional until AuthURL/Exchange are used.
func New(cfg git.Config) (git.Client, error) {
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"read:user", "user:email", "repo"}
	}

	c := &Client{cfg: cfg}
	if cfg.ClientID != "" {
		c.oauth = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     oauth2github.Endpoint,
		}
	}
	return c, nil
}

// Provider returns the GitHub provider identifier.
func (c *Client) Provider() git.Provider {
	return git.ProviderGitHub
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

func (c *Client) api(ctx context.Context) (*gh.Client, error) {
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

	client, err := gh.NewClient(gh.WithAuthToken(c.token))
	if err != nil {
		return nil, fmt.Errorf("github: create client: %w", err)
	}
	c.client = client
	_ = ctx
	return c.client, nil
}

func listOpts(opts git.ListOptions) *gh.ListOptions {
	opts = git.NormalizeListOptions(opts)
	return &gh.ListOptions{Page: opts.Page, PerPage: opts.PerPage}
}

func pageInfoFrom(resp *gh.Response) git.PageInfo {
	if resp == nil {
		return git.PageInfo{}
	}
	return git.PageInfo{
		NextPage:  resp.NextPage,
		PrevPage:  resp.PrevPage,
		FirstPage: resp.FirstPage,
		LastPage:  resp.LastPage,
	}
}
