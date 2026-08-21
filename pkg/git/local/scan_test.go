package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/selyusize/my-home/pkg/git"
)

func TestParseRemote(t *testing.T) {
	t.Parallel()

	cases := []struct {
		raw      string
		provider git.Provider
		fullName string
	}{
		{"https://github.com/selyusize/my-home.git", git.ProviderGitHub, "selyusize/my-home"},
		{"https://github.com/selyusize/my-home", git.ProviderGitHub, "selyusize/my-home"},
		{"git@github.com:selyusize/my-home.git", git.ProviderGitHub, "selyusize/my-home"},
		{"ssh://git@github.com/selyusize/my-home.git", git.ProviderGitHub, "selyusize/my-home"},
		{"https://www.github.com/selyusize/my-home.git", git.ProviderGitHub, "selyusize/my-home"},
		{"https://gitlab.com/group/project.git", git.ProviderGitLab, "group/project"},
		{"git@gitlab.com:group/sub/project.git", git.ProviderGitLab, "group/sub/project"},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			t.Parallel()
			remote, ok := ParseRemote(tc.raw)
			if !ok {
				t.Fatal("expected parse")
			}
			if remote.Provider != tc.provider {
				t.Fatalf("provider=%q", remote.Provider)
			}
			if remote.FullName != tc.fullName {
				t.Fatalf("fullName=%q", remote.FullName)
			}
		})
	}
}

func TestParseRemoteRejectsUnknown(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		"",
		"https://example.com/org/repo.git",
		"not-a-url",
		"git@github.com:onlyname.git",
	} {
		if _, ok := ParseRemote(raw); ok {
			t.Fatalf("parsed %q", raw)
		}
	}
}

func TestScanFindsOriginRemote(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repo := filepath.Join(root, "my-home")
	writeGitConfig(t, filepath.Join(repo, ".git"), `[core]
	bare = false
[remote "origin"]
	url = git@github.com:selyusize/my-home.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`)

	clones, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].Provider != string(git.ProviderGitHub) {
		t.Fatalf("provider=%q", clones[0].Provider)
	}
	if clones[0].FullName != "selyusize/my-home" {
		t.Fatalf("fullName=%q", clones[0].FullName)
	}
	if clones[0].Path != repo {
		t.Fatalf("path=%q", clones[0].Path)
	}
	if clones[0].HTMLURL != "https://github.com/selyusize/my-home" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
}

func TestScanIncludesRootRepo(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeGitConfig(t, filepath.Join(root, ".git"), `[remote "origin"]
	url = https://gitlab.com/acme/app.git
`)

	clones, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].FullName != "acme/app" {
		t.Fatalf("fullName=%q", clones[0].FullName)
	}
	if clones[0].HTMLURL != "https://gitlab.com/acme/app" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
}

func TestScanIncludesGitWithoutRemote(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repo := filepath.Join(root, "notes")
	writeGitConfig(t, filepath.Join(repo, ".git"), "[core]\n\tbare = false\n")

	clones, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].Name != "notes" {
		t.Fatalf("name=%q", clones[0].Name)
	}
	if clones[0].Provider != "" {
		t.Fatalf("provider=%q", clones[0].Provider)
	}
	if clones[0].HTMLURL != "" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
	if clones[0].Path != repo {
		t.Fatalf("path=%q", clones[0].Path)
	}
}

func TestScanUnknownHostHasNoHTMLURL(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repo := filepath.Join(root, "other")
	writeGitConfig(t, filepath.Join(repo, ".git"), `[remote "origin"]
	url = https://example.com/org/repo.git
`)

	clones, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(clones) != 1 {
		t.Fatalf("len=%d", len(clones))
	}
	if clones[0].Name != "other" {
		t.Fatalf("name=%q", clones[0].Name)
	}
	if clones[0].HTMLURL != "" {
		t.Fatalf("htmlUrl=%q", clones[0].HTMLURL)
	}
}

func TestHTMLURL(t *testing.T) {
	t.Parallel()

	if got := HTMLURL(string(git.ProviderGitHub), "acme/app"); got != "https://github.com/acme/app" {
		t.Fatalf("github=%q", got)
	}
	if got := HTMLURL(string(git.ProviderGitLab), "group/sub/app"); got != "https://gitlab.com/group/sub/app" {
		t.Fatalf("gitlab=%q", got)
	}
	if got := HTMLURL("example", "acme/app"); got != "" {
		t.Fatalf("unknown=%q", got)
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
