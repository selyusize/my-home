package git

import "context"

// Client is the common API implemented by every git hosting provider.
type Client interface {
	Provider() Provider

	IsAuthenticated() bool
	AccessToken() string
	Logout()

	AuthURL(state string) (string, error)
	ExchangeCode(ctx context.Context, code string) error
	SetAccessToken(token string) error

	Profile(ctx context.Context) (*User, error)
	User(ctx context.Context, login string) (*User, error)
	ContributionCalendar(ctx context.Context) (*ContributionCalendar, error)

	ListRepos(ctx context.Context, opts ListOptions) ([]Repository, PageInfo, error)
	Repository(ctx context.Context, owner, name string) (*Repository, error)
	ListBranches(ctx context.Context, owner, name string, opts ListOptions) ([]Branch, PageInfo, error)
	ListCommits(ctx context.Context, owner, name, ref string, opts ListOptions) ([]Commit, PageInfo, error)

	ListUserActivity(ctx context.Context, login string, opts ListOptions) ([]ActivityEvent, PageInfo, error)
	ListRepositoryActivity(ctx context.Context, owner, name string, opts ListOptions) ([]ActivityEvent, PageInfo, error)
}

// Constructor creates a Client for a registered provider.
type Constructor func(cfg Config) (Client, error)
