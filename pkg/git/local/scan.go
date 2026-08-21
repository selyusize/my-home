package local

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/selyusize/my-home/pkg/git"
)

// Clone is a local git repository, optionally matched to a hosting provider.
type Clone struct {
	Provider string `json:"provider"`
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Path     string `json:"path"`
	HTMLURL  string `json:"htmlUrl"`
}

// Remote is a parsed git hosting remote.
type Remote struct {
	Provider git.Provider
	Owner    string
	Name     string
	FullName string
}

// Scan walks root and its immediate subdirectories for git clones.
func Scan(root string) ([]Clone, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, nil
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("git/local: abs %s: %w", root, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("git/local: stat %s: %w", abs, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("git/local: %s is not a directory", abs)
	}

	var clones []Clone
	if clone, ok := inspect(abs); ok {
		clones = append(clones, clone)
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, fmt.Errorf("git/local: read %s: %w", abs, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		child := filepath.Join(abs, entry.Name())
		clone, ok := inspect(child)
		if !ok {
			continue
		}
		clones = append(clones, clone)
	}
	return clones, nil
}

func inspect(dir string) (Clone, bool) {
	if !isGitRepo(dir) {
		return Clone{}, false
	}
	clone := Clone{
		Name: filepath.Base(dir),
		Path: dir,
	}
	remote, ok := readRemote(dir)
	if !ok {
		return clone, true
	}
	clone.Provider = string(remote.Provider)
	clone.Owner = remote.Owner
	clone.Name = remote.Name
	clone.FullName = remote.FullName
	clone.HTMLURL = HTMLURL(clone.Provider, clone.FullName)
	return clone, true
}

func isGitRepo(dir string) bool {
	info, err := os.Lstat(filepath.Join(dir, ".git"))
	if err != nil {
		return false
	}
	return info.IsDir() || info.Mode().IsRegular()
}

// HTMLURL is the https page for a github.com or gitlab.com clone.
func HTMLURL(provider, fullName string) string {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return ""
	}
	switch git.Provider(provider) {
	case git.ProviderGitHub:
		return "https://github.com/" + fullName
	case git.ProviderGitLab:
		return "https://gitlab.com/" + fullName
	default:
		return ""
	}
}

func readRemote(dir string) (Remote, bool) {
	configPath, ok := gitConfigPath(dir)
	if !ok {
		return Remote{}, false
	}
	urls, err := remoteURLs(configPath)
	if err != nil || len(urls) == 0 {
		return Remote{}, false
	}
	if origin, found := urls["origin"]; found {
		remote, parsed := ParseRemote(origin)
		return remote, parsed
	}
	for _, raw := range urls {
		remote, parsed := ParseRemote(raw)
		if parsed {
			return remote, true
		}
	}
	return Remote{}, false
}

func gitConfigPath(dir string) (string, bool) {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Lstat(gitPath)
	if err != nil {
		return "", false
	}
	if info.IsDir() {
		config := filepath.Join(gitPath, "config")
		if fileExists(config) {
			return config, true
		}
		return "", false
	}
	if !info.Mode().IsRegular() {
		return "", false
	}
	gitDir, ok := parseGitFile(gitPath)
	if !ok {
		return "", false
	}
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(dir, gitDir)
	}
	config := filepath.Join(gitDir, "config")
	if fileExists(config) {
		return config, true
	}
	common := filepath.Join(gitDir, "commondir")
	if !fileExists(common) {
		return "", false
	}
	data, err := os.ReadFile(common)
	if err != nil {
		return "", false
	}
	commonDir := strings.TrimSpace(string(data))
	if commonDir == "" {
		return "", false
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(gitDir, commonDir)
	}
	config = filepath.Join(commonDir, "config")
	if fileExists(config) {
		return config, true
	}
	return "", false
}

func parseGitFile(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	line := strings.TrimSpace(string(data))
	const prefix = "gitdir:"
	if !strings.HasPrefix(strings.ToLower(line), prefix) {
		return "", false
	}
	gitDir := strings.TrimSpace(line[len(prefix):])
	if gitDir == "" {
		return "", false
	}
	return gitDir, true
}

func remoteURLs(configPath string) (map[string]string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urls := make(map[string]string)
	var remoteName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			remoteName = parseRemoteSection(line)
			continue
		}
		if remoteName == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(key) != "url" {
			continue
		}
		urls[remoteName] = strings.TrimSpace(value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func parseRemoteSection(line string) string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
		return ""
	}
	body := strings.TrimSpace(line[1 : len(line)-1])
	name, rest, ok := strings.Cut(body, " ")
	if !ok || name != "remote" {
		return ""
	}
	return strings.Trim(strings.TrimSpace(rest), `"`)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// ParseRemote turns a git remote URL into a hosting provider and full name.
func ParseRemote(raw string) (Remote, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Remote{}, false
	}

	host, path, ok := splitRemote(raw)
	if !ok {
		return Remote{}, false
	}
	provider := providerForHost(host)
	if provider == git.ProviderUnknown {
		return Remote{}, false
	}

	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")
	path = strings.Trim(path, "/")
	if path == "" {
		return Remote{}, false
	}
	owner, name, ok := splitFullName(path)
	if !ok {
		return Remote{}, false
	}
	return Remote{
		Provider: provider,
		Owner:    owner,
		Name:     name,
		FullName: owner + "/" + name,
	}, true
}

func splitRemote(raw string) (host, path string, ok bool) {
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host == "" || parsed.Path == "" {
			return "", "", false
		}
		return parsed.Hostname(), parsed.Path, true
	}

	if strings.HasPrefix(raw, "git@") || strings.Contains(raw, ":") {
		withoutUser := raw
		if at := strings.LastIndex(raw, "@"); at >= 0 {
			withoutUser = raw[at+1:]
		}
		host, path, found := strings.Cut(withoutUser, ":")
		if !found || host == "" || path == "" {
			return "", "", false
		}
		if strings.Contains(host, "/") {
			return "", "", false
		}
		return host, path, true
	}
	return "", "", false
}

func providerForHost(host string) git.Provider {
	host = strings.ToLower(strings.TrimPrefix(host, "www."))
	switch host {
	case "github.com":
		return git.ProviderGitHub
	case "gitlab.com":
		return git.ProviderGitLab
	default:
		return git.ProviderUnknown
	}
}

func splitFullName(path string) (owner, name string, ok bool) {
	path = strings.Trim(path, "/")
	slash := strings.LastIndex(path, "/")
	if slash <= 0 || slash == len(path)-1 {
		return "", "", false
	}
	owner = path[:slash]
	name = path[slash+1:]
	if owner == "" || name == "" {
		return "", "", false
	}
	return owner, name, true
}
