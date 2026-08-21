package github

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/git"

	gh "github.com/google/go-github/v90/github"
)

// Profile returns the authenticated user's profile.
func (c *Client) Profile(ctx context.Context) (*git.User, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("github: get profile: %w", err)
	}
	return mapUser(user), nil
}

// User returns a public profile by login.
func (c *Client) User(ctx context.Context, login string) (*git.User, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}
	if login == "" {
		return nil, fmt.Errorf("github: empty login")
	}

	user, _, err := client.Users.Get(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("github: get user %q: %w", login, err)
	}
	return mapUser(user), nil
}

func mapUser(u *gh.User) *git.User {
	if u == nil {
		return nil
	}
	return &git.User{
		ID:          deref(u.ID),
		Login:       deref(u.Login),
		Name:        deref(u.Name),
		Email:       deref(u.Email),
		Bio:         deref(u.Bio),
		Company:     deref(u.Company),
		Location:    deref(u.Location),
		Blog:        deref(u.Blog),
		AvatarURL:   deref(u.AvatarURL),
		HTMLURL:     deref(u.HTMLURL),
		PublicRepos: deref(u.PublicRepos),
		Followers:   deref(u.Followers),
		Following:   deref(u.Following),
		CreatedAt:   timeOrZero(u.CreatedAt),
	}
}
