package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathBelongsTo(t *testing.T) {
	t.Parallel()

	root := filepath.Join(string(filepath.Separator), "repos")
	child := filepath.Join(root, "app")
	if !pathBelongsTo(child, root) {
		t.Fatal("expected child to belong")
	}
	if !pathBelongsTo(root, root) {
		t.Fatal("expected root to belong to itself")
	}
	if pathBelongsTo(filepath.Join(string(filepath.Separator), "other"), root) {
		t.Fatal("expected outsider to be rejected")
	}
}

func TestCollectClonesKeepsUnknownInSeparateRoot(t *testing.T) {
	t.Parallel()

	githubRoot := t.TempDir()
	gitlabRoot := t.TempDir()
	repo := filepath.Join(gitlabRoot, "notes")
	writeGitConfig(t, filepath.Join(repo, ".git"), "[core]\n\tbare = false\n")

	clones, err := collectClones(LocalRepoSettings{
		GitHubSeparate: true,
		GitHubPath:     githubRoot,
		GitLabSeparate: true,
		GitLabPath:     gitlabRoot,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].Path != repo {
		t.Fatalf("path=%q", clones[0].Path)
	}
	if clones[0].HTMLURL != "" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
}

func TestCollectClonesKeepsGitHubCloneInGitLabRoot(t *testing.T) {
	t.Parallel()

	githubRoot := t.TempDir()
	gitlabRoot := t.TempDir()
	repo := filepath.Join(gitlabRoot, "my-home")
	writeGitConfig(t, filepath.Join(repo, ".git"), `[remote "origin"]
	url = https://github.com/acme/my-home.git
`)

	clones, err := collectClones(LocalRepoSettings{
		GitHubSeparate: true,
		GitHubPath:     githubRoot,
		GitLabSeparate: true,
		GitLabPath:     gitlabRoot,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].Path != repo {
		t.Fatalf("path=%q", clones[0].Path)
	}
	if clones[0].HTMLURL != "https://github.com/acme/my-home" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
}

func writeGitConfig(t *testing.T, gitDir, contents string) {
	t.Helper()
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
