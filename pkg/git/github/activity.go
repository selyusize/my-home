package github

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/git"

	gh "github.com/google/go-github/v90/github"
)

// ListUserActivity returns events performed by the authenticated user.
// If login is empty, the authenticated profile login is used.
func (c *Client) ListUserActivity(ctx context.Context, login string, opts git.ListOptions) ([]git.ActivityEvent, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}

	if login == "" {
		profile, err := c.Profile(ctx)
		if err != nil {
			return nil, git.PageInfo{}, err
		}
		login = profile.Login
	}

	events, resp, err := client.Activity.ListEventsPerformedByUser(ctx, login, false, listOpts(opts))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("github: list user activity: %w", err)
	}

	return mapEvents(events), pageInfoFrom(resp), nil
}

// ListRepositoryActivity returns events for a repository.
func (c *Client) ListRepositoryActivity(ctx context.Context, owner, name string, opts git.ListOptions) ([]git.ActivityEvent, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("github: owner and name are required")
	}

	events, resp, err := client.Activity.ListRepositoryEvents(ctx, owner, name, listOpts(opts))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("github: list repository activity: %w", err)
	}

	return mapEvents(events), pageInfoFrom(resp), nil
}

func mapEvents(events []*gh.Event) []git.ActivityEvent {
	out := make([]git.ActivityEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}

		actor := ""
		if event.Actor != nil {
			actor = deref(event.Actor.Login)
		}

		repoName := ""
		repoURL := ""
		if event.Repo != nil {
			repoName = deref(event.Repo.Name)
			repoURL = deref(event.Repo.URL)
		}

		out = append(out, git.ActivityEvent{
			ID:        deref(event.ID),
			Type:      deref(event.Type),
			Actor:     actor,
			RepoName:  repoName,
			RepoURL:   repoURL,
			IsPublic:  deref(event.Public),
			CreatedAt: timeOrZero(event.CreatedAt),
		})
	}
	return out
}
