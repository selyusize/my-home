package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v90/github"
)

// ListUserActivity returns events performed by the authenticated user.
// If login is empty, the authenticated profile login is used.
func (s *Service) ListUserActivity(ctx context.Context, login string, opts ListOptions) ([]ActivityEvent, PageInfo, error) {
	client, err := s.api(ctx)
	if err != nil {
		return nil, PageInfo{}, err
	}

	if login == "" {
		profile, err := s.GetProfile(ctx)
		if err != nil {
			return nil, PageInfo{}, err
		}
		login = profile.Login
	}

	events, resp, err := client.Activity.ListEventsPerformedByUser(ctx, login, false, listOpts(opts))
	if err != nil {
		return nil, PageInfo{}, fmt.Errorf("github: list user activity: %w", err)
	}

	return mapEvents(events), pageInfoFrom(resp), nil
}

// ListRepositoryActivity returns events for a repository.
func (s *Service) ListRepositoryActivity(ctx context.Context, owner, name string, opts ListOptions) ([]ActivityEvent, PageInfo, error) {
	client, err := s.api(ctx)
	if err != nil {
		return nil, PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, PageInfo{}, fmt.Errorf("github: owner and name are required")
	}

	events, resp, err := client.Activity.ListRepositoryEvents(ctx, owner, name, listOpts(opts))
	if err != nil {
		return nil, PageInfo{}, fmt.Errorf("github: list repository activity: %w", err)
	}

	return mapEvents(events), pageInfoFrom(resp), nil
}

func mapEvents(events []*gh.Event) []ActivityEvent {
	out := make([]ActivityEvent, 0, len(events))
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

		out = append(out, ActivityEvent{
			ID:        deref(event.ID),
			Type:      deref(event.Type),
			Actor:     actor,
			RepoName:  repoName,
			RepoURL:   repoURL,
			Public:    deref(event.Public),
			CreatedAt: timeOrZero(event.CreatedAt),
		})
	}
	return out
}
