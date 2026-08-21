package github

import (
	"context"
	"fmt"
	"time"

	"github.com/selyusize/my-home/pkg/git"

	gh "github.com/google/go-github/v90/github"
)

// ListRepos returns repositories for the authenticated user.
func (c *Client) ListRepos(ctx context.Context, opts git.ListOptions) ([]git.Repository, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}

	repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: *listOpts(opts),
		Affiliation: "owner,organization_member",
		Sort:        "updated",
		Direction:   "desc",
	})
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("github: list repositories: %w", err)
	}

	out := make([]git.Repository, 0, len(repos))
	for _, repo := range repos {
		if mapped := mapRepository(repo); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

// Repository returns repository details by owner and name.
func (c *Client) Repository(ctx context.Context, owner, name string) (*git.Repository, error) {
	client, err := c.api(ctx)
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

// ListBranches returns branches for a repository.
func (c *Client) ListBranches(ctx context.Context, owner, name string, opts git.ListOptions) ([]git.Branch, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("github: owner and name are required")
	}

	branches, resp, err := client.Repositories.ListBranches(ctx, owner, name, &gh.BranchListOptions{
		ListOptions: *listOpts(opts),
	})
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("github: list branches %s/%s: %w", owner, name, err)
	}

	out := make([]git.Branch, 0, len(branches))
	for _, branch := range branches {
		if mapped := mapBranch(branch); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

// ListCommits returns commits for a repository. An empty ref uses the default branch.
func (c *Client) ListCommits(ctx context.Context, owner, name, ref string, opts git.ListOptions) ([]git.Commit, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("github: owner and name are required")
	}

	commitOpts := &gh.CommitsListOptions{ListOptions: *listOpts(opts)}
	if ref != "" {
		commitOpts.SHA = ref
	}

	commits, resp, err := client.Repositories.ListCommits(ctx, owner, name, commitOpts)
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("github: list commits %s/%s: %w", owner, name, err)
	}

	out := make([]git.Commit, 0, len(commits))
	for _, commit := range commits {
		if mapped := mapCommit(commit); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

func mapRepository(r *gh.Repository) *git.Repository {
	if r == nil {
		return nil
	}

	ownerLogin := ""
	ownerAvatar := ""
	if r.Owner != nil {
		ownerLogin = deref(r.Owner.Login)
		ownerAvatar = deref(r.Owner.AvatarURL)
	}

	return &git.Repository{
		ID:            deref(r.ID),
		Name:          deref(r.Name),
		FullName:      deref(r.FullName),
		Description:   deref(r.Description),
		HTMLURL:       deref(r.HTMLURL),
		CloneURL:      deref(r.CloneURL),
		SSHURL:        deref(r.SSHURL),
		DefaultBranch: deref(r.DefaultBranch),
		Language:      deref(r.Language),
		IsPrivate:     deref(r.Private),
		IsFork:        deref(r.Fork),
		IsArchived:    deref(r.Archived),
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

func mapBranch(b *gh.Branch) *git.Branch {
	if b == nil {
		return nil
	}

	sha := ""
	if b.Commit != nil {
		sha = deref(b.Commit.SHA)
	}

	return &git.Branch{
		Name:        deref(b.Name),
		SHA:         sha,
		IsProtected: deref(b.Protected),
	}
}

func mapCommit(c *gh.RepositoryCommit) *git.Commit {
	if c == nil {
		return nil
	}

	message := ""
	author := ""
	var date time.Time
	if c.Commit != nil {
		message = deref(c.Commit.Message)
		if c.Commit.Author != nil {
			author = deref(c.Commit.Author.Name)
			date = timeOrZero(c.Commit.Author.Date)
		}
	}

	avatar := ""
	if c.Author != nil {
		avatar = deref(c.Author.AvatarURL)
		if author == "" {
			author = deref(c.Author.Login)
		}
	}

	return &git.Commit{
		SHA:       deref(c.SHA),
		Message:   message,
		Author:    author,
		AvatarURL: avatar,
		HTMLURL:   deref(c.HTMLURL),
		Date:      date,
	}
}
