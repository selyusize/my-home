package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v90/github"
)

// ListRepos returns repositories for the authenticated user.
func (s *Service) ListRepos(ctx context.Context, opts ListOptions) ([]Repository, PageInfo, error) {
	client, err := s.api(ctx)
	if err != nil {
		return nil, PageInfo{}, err
	}

	repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: *listOpts(opts),
		Sort:        "updated",
		Direction:   "desc",
	})
	if err != nil {
		return nil, PageInfo{}, fmt.Errorf("github: list repositories: %w", err)
	}

	out := make([]Repository, 0, len(repos))
	for _, repo := range repos {
		if mapped := mapRepository(repo); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

// GetRepository returns repository details by owner and name.
func (s *Service) GetRepository(ctx context.Context, owner, name string) (*Repository, error) {
	client, err := s.api(ctx)
	if err != nil {
		return nil, err
	}
	if owner == "" || name == "" {
		return nil, fmt.Errorf("github: owner and name are required")
	}

	repo, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, fmt.Errorf("github: get repository %s/%s: %w", owner, name, err)
	}
	return mapRepository(repo), nil
}

func mapRepository(r *gh.Repository) *Repository {
	if r == nil {
		return nil
	}

	ownerLogin := ""
	ownerAvatar := ""
	if r.Owner != nil {
		ownerLogin = deref(r.Owner.Login)
		ownerAvatar = deref(r.Owner.AvatarURL)
	}

	return &Repository{
		ID:            deref(r.ID),
		Name:          deref(r.Name),
		FullName:      deref(r.FullName),
		Description:   deref(r.Description),
		HTMLURL:       deref(r.HTMLURL),
		CloneURL:      deref(r.CloneURL),
		SSHURL:        deref(r.SSHURL),
		DefaultBranch: deref(r.DefaultBranch),
		Language:      deref(r.Language),
		Private:       deref(r.Private),
		Fork:          deref(r.Fork),
		Archived:      deref(r.Archived),
		Stars:         deref(r.StargazersCount),
		Forks:         deref(r.ForksCount),
		OpenIssues:    deref(r.OpenIssuesCount),
		OwnerLogin:    ownerLogin,
		OwnerAvatar:   ownerAvatar,
		PushedAt:      timeOrZero(r.PushedAt),
		UpdatedAt:     timeOrZero(r.UpdatedAt),
		CreatedAt:     timeOrZero(r.CreatedAt),
	}
}
