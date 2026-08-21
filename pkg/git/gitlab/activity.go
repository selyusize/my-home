package gitlab

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/selyusize/my-home/pkg/git"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

// ListUserActivity returns contribution events for a user.
// If login is empty, events of the authenticated user are returned.
func (c *Client) ListUserActivity(ctx context.Context, login string, opts git.ListOptions) ([]git.ActivityEvent, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}

	list := &gitlab.ListContributionEventsOptions{ListOptions: listOpts(opts)}
	var (
		events []*gitlab.ContributionEvent
		resp   *gitlab.Response
	)

	if login == "" {
		events, resp, err = client.Events.ListCurrentUserContributionEvents(list, gitlab.WithContext(ctx))
	} else {
		events, resp, err = client.Users.ListUserContributionEvents(login, list, gitlab.WithContext(ctx))
	}
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: list user activity: %w", err)
	}

	return mapContributionEvents(events), pageInfoFrom(resp), nil
}

// ListRepositoryActivity returns visible events for a project.
func (c *Client) ListRepositoryActivity(ctx context.Context, owner, name string, opts git.ListOptions) ([]git.ActivityEvent, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: owner and name are required")
	}

	pid := projectPath(owner, name)
	events, resp, err := client.Events.ListProjectVisibleEvents(pid, &gitlab.ListProjectVisibleEventsOptions{
		ListOptions: listOpts(opts),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: list repository activity: %w", err)
	}

	fullName := pid
	return mapProjectEvents(events, fullName), pageInfoFrom(resp), nil
}

func mapContributionEvents(events []*gitlab.ContributionEvent) []git.ActivityEvent {
	out := make([]git.ActivityEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}

		out = append(out, git.ActivityEvent{
			ID:        strconv.FormatInt(event.ID, 10),
			Type:      eventType(event.ActionName, event.TargetType),
			Actor:     event.AuthorUsername,
			RepoName:  strconv.FormatInt(event.ProjectID, 10),
			IsPublic:  true,
			CreatedAt: timeOrZero(event.CreatedAt),
		})
	}
	return out
}

func mapProjectEvents(events []*gitlab.ProjectEvent, fullName string) []git.ActivityEvent {
	out := make([]git.ActivityEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}

		repoURL := ""
		if event.Data.Repository != nil {
			repoURL = event.Data.Repository.WebURL
			if fullName == "" {
				fullName = event.Data.Repository.PathWithNamespace
			}
		}

		out = append(out, git.ActivityEvent{
			ID:        strconv.FormatInt(event.ID, 10),
			Type:      eventType(event.ActionName, event.TargetType),
			Actor:     event.AuthorUsername,
			RepoName:  fullName,
			RepoURL:   repoURL,
			IsPublic:  true,
			CreatedAt: parseTime(event.CreatedAt),
		})
	}
	return out
}

func eventType(action, target string) string {
	action = strings.TrimSpace(action)
	target = strings.TrimSpace(target)
	switch {
	case action != "" && target != "":
		return action + ":" + target
	case action != "":
		return action
	default:
		return target
	}
}

func parseTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}

	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05 MST", "2006-01-02 15:04:05 Z07:00"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}
	return time.Time{}
}
