package gitlab

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/git"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"
)

// ListRepos returns projects visible to the authenticated user.
func (c *Client) ListRepos(ctx context.Context, opts git.ListOptions) ([]git.Repository, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}

	projects, resp, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		ListOptions: listOpts(opts),
		Membership:  gitlab.Ptr(true),
		OrderBy:     gitlab.Ptr("updated_at"),
		Sort:        gitlab.Ptr("desc"),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: list repositories: %w", err)
	}

	out := make([]git.Repository, 0, len(projects))
	for _, project := range projects {
		if mapped := mapProject(project, ""); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

// Repository returns project details by namespace and path.
func (c *Client) Repository(ctx context.Context, owner, name string) (*git.Repository, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}
	if owner == "" || name == "" {
		return nil, fmt.Errorf("gitlab: owner and name are required")
	}

	pid := projectPath(owner, name)
	project, _, err := client.Projects.GetProject(pid, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("gitlab: get repository %s: %w", pid, err)
	}

	mapped := mapProject(project, "")
	if mapped == nil {
		return nil, fmt.Errorf("gitlab: empty repository %s", pid)
	}

	langs, _, err := client.Projects.GetProjectLanguages(pid, gitlab.WithContext(ctx))
	if err == nil {
		mapped.Language = primaryLanguage(langs)
	}
	return mapped, nil
}

// ListBranches returns branches for a project.
func (c *Client) ListBranches(ctx context.Context, owner, name string, opts git.ListOptions) ([]git.Branch, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: owner and name are required")
	}

	pid := projectPath(owner, name)
	branches, resp, err := client.Branches.ListBranches(pid, &gitlab.ListBranchesOptions{
		ListOptions: listOpts(opts),
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: list branches %s: %w", pid, err)
	}

	out := make([]git.Branch, 0, len(branches))
	for _, branch := range branches {
		if mapped := mapBranch(branch); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

// ListCommits returns commits for a project. An empty ref uses the default branch.
func (c *Client) ListCommits(ctx context.Context, owner, name, ref string, opts git.ListOptions) ([]git.Commit, git.PageInfo, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, git.PageInfo{}, err
	}
	if owner == "" || name == "" {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: owner and name are required")
	}

	pid := projectPath(owner, name)
	commitOpts := &gitlab.ListCommitsOptions{ListOptions: listOpts(opts)}
	if ref != "" {
		commitOpts.RefName = gitlab.Ptr(ref)
	}

	commits, resp, err := client.Commits.ListCommits(pid, commitOpts, gitlab.WithContext(ctx))
	if err != nil {
		return nil, git.PageInfo{}, fmt.Errorf("gitlab: list commits %s: %w", pid, err)
	}

	out := make([]git.Commit, 0, len(commits))
	for _, commit := range commits {
		if mapped := mapCommit(commit); mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out, pageInfoFrom(resp), nil
}

func projectPath(owner, name string) string {
	return owner + "/" + name
}

func mapProject(p *gitlab.Project, language string) *git.Repository {
	if p == nil {
		return nil
	}

	ownerLogin := ""
	ownerAvatar := ""
	if p.Namespace != nil {
		ownerLogin = p.Namespace.Path
		ownerAvatar = p.Namespace.AvatarURL
	}
	if p.Owner != nil {
		if ownerLogin == "" {
			ownerLogin = p.Owner.Username
		}
		if ownerAvatar == "" {
			ownerAvatar = p.Owner.AvatarURL
		}
	}

	return &git.Repository{
		ID:            p.ID,
		Name:          p.Name,
		FullName:      p.PathWithNamespace,
		Description:   p.Description,
		HTMLURL:       p.WebURL,
		CloneURL:      p.HTTPURLToRepo,
		SSHURL:        p.SSHURLToRepo,
		DefaultBranch: p.DefaultBranch,
		Language:      language,
		IsPrivate:     p.Visibility == gitlab.PrivateVisibility,
		IsFork:        p.ForkedFromProject != nil,
		IsArchived:    p.Archived,
		Stars:         int(p.StarCount),
		Forks:         int(p.ForksCount),
		OpenIssues:    int(p.OpenIssuesCount),
		OwnerLogin:    ownerLogin,
		OwnerAvatar:   ownerAvatar,
		PushedAt:      timeOrZero(p.LastActivityAt),
		UpdatedAt:     timeOrZero(p.UpdatedAt),
		CreatedAt:     timeOrZero(p.CreatedAt),
	}
}

func primaryLanguage(langs *gitlab.ProjectLanguages) string {
	if langs == nil || len(*langs) == 0 {
		return ""
	}

	name := ""
	var best float32
	for lang, share := range *langs {
		if share >= best {
			best = share
			name = lang
		}
	}
	return name
}

func mapBranch(b *gitlab.Branch) *git.Branch {
	if b == nil {
		return nil
	}

	sha := ""
	if b.Commit != nil {
		sha = b.Commit.ID
	}

	return &git.Branch{
		Name:        b.Name,
		SHA:         sha,
		IsProtected: b.Protected,
	}
}

func mapCommit(c *gitlab.Commit) *git.Commit {
	if c == nil {
		return nil
	}

	return &git.Commit{
		SHA:       c.ID,
		Message:   c.Message,
		Author:    c.AuthorName,
		HTMLURL:   c.WebURL,
		Date:      timeOrZero(c.CommittedDate),
	}
}
