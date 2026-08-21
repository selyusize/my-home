package gitlab

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/selyusize/my-home/pkg/git"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

// Profile returns the authenticated user's profile.
func (c *Client) Profile(ctx context.Context) (*git.User, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}

	user, _, err := client.Users.CurrentUser(gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("gitlab: get profile: %w", err)
	}
	return c.mapUser(ctx, client, user), nil
}

// User returns a public profile by username.
func (c *Client) User(ctx context.Context, login string) (*git.User, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}
	if login == "" {
		return nil, fmt.Errorf("gitlab: empty login")
	}

	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: gitlab.Ptr(login),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("gitlab: get user %q: %w", login, err)
	}

	matched := findUser(users, login)
	if matched == nil {
		users, _, err = client.Users.ListUsers(&gitlab.ListUsersOptions{
			Search: gitlab.Ptr(login),
		}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("gitlab: search user %q: %w", login, err)
		}
		matched = findUser(users, login)
	}
	if matched == nil {
		return nil, fmt.Errorf("gitlab: user %q not found", login)
	}

	return c.mapUser(ctx, client, matched), nil
}

func findUser(users []*gitlab.User, login string) *gitlab.User {
	for _, user := range users {
		if user != nil && strings.EqualFold(user.Username, login) {
			return user
		}
	}
	return nil
}

func (c *Client) mapUser(ctx context.Context, client *gitlab.Client, u *gitlab.User) *git.User {
	if u == nil {
		return nil
	}

	email := u.Email
	if email == "" {
		email = u.PublicEmail
	}

	out := &git.User{
		ID:        u.ID,
		Login:     u.Username,
		Name:      u.Name,
		Email:     email,
		Bio:       u.Bio,
		Company:   u.Organization,
		Location:  u.Location,
		Blog:      u.WebsiteURL,
		AvatarURL: u.AvatarURL,
		HTMLURL:   u.WebURL,
		CreatedAt: timeOrZero(u.CreatedAt),
	}

	if client == nil {
		return out
	}

	counts, _, err := client.Users.GetUserAssociationsCount(u.ID, gitlab.WithContext(ctx))
	if err == nil && counts != nil {
		out.PublicRepos = int(counts.ProjectsCount)
	}
	return out
}

func timeOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
