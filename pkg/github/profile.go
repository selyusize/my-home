package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v90/github"
)

// GetProfile returns the authenticated user's profile.
func (s *Service) GetProfile(ctx context.Context) (*User, error) {
	client, err := s.api(ctx)
	if err != nil {
		return nil, err
	}

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("github: get profile: %w", err)
	}
	return mapUser(user), nil
}

// GetUser returns a public profile by login.
func (s *Service) GetUser(ctx context.Context, login string) (*User, error) {
	client, err := s.api(ctx)
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

func mapUser(u *gh.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:            deref(u.ID),
		Login:         deref(u.Login),
		Name:          deref(u.Name),
		Email:         deref(u.Email),
		Bio:           deref(u.Bio),
		Company:       deref(u.Company),
		Location:      deref(u.Location),
		Blog:          deref(u.Blog),
		AvatarURL:     deref(u.AvatarURL),
		HTMLURL:       deref(u.HTMLURL),
		PublicRepos: deref(u.PublicRepos),
		Followers:     deref(u.Followers),
		Following:     deref(u.Following),
		CreatedAt:     timeOrZero(u.CreatedAt),
	}
}
