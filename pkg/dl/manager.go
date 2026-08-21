package dl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const appDirName = "my-home"

const githubLatestReleaseURL = "https://api.github.com/repos/local-deploy/dl/releases/latest"

// Manager owns an isolated official dl binary and runs it.
type Manager struct {
	root      string
	latestURL string
}

// New creates a manager in the default runtime directory.
func New() (*Manager, error) {
	root, err := DefaultRoot()
	if err != nil {
		return nil, err
	}
	return NewAt(root), nil
}

// Option configures a Manager.
type Option func(*Manager)

// WithReleaseURL overrides the GitHub latest-release endpoint.
func WithReleaseURL(url string) Option {
	return func(m *Manager) {
		if url != "" {
			m.latestURL = url
		}
	}
}

// NewAt creates a manager rooted at the given runtime directory.
func NewAt(root string, opts ...Option) *Manager {
	m := &Manager{
		root:      root,
		latestURL: githubLatestReleaseURL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(m)
		}
	}
	return m
}

// DefaultRoot is os.UserConfigDir()/my-home/runtime/dl.
func DefaultRoot() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("dl: user config dir: %w", err)
	}
	return filepath.Join(dir, appDirName, "runtime", "dl"), nil
}

func (m *Manager) binPath() string {
	return filepath.Join(m.root, "bin", "dl")
}

func (m *Manager) homeDir() string {
	return filepath.Join(m.root, "home")
}

func (m *Manager) ensureDirs() error {
	for _, dir := range []string{filepath.Join(m.root, "bin"), m.homeDir()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("dl: create runtime: %w", err)
		}
	}
	return nil
}

// Status reports the managed binary without requiring a project directory.
func (m *Manager) Status(ctx context.Context) (*Status, error) {
	st := &Status{
		Path:     m.binPath(),
		Services: defaultServices(),
	}
	if _, err := os.Stat(st.Path); err == nil {
		st.IsInstalled = true
		result, err := m.Exec(ctx, m.root, []string{"version"})
		if err == nil {
			st.Version = strings.TrimSpace(result.Stdout)
		}
	}

	latest, err := fetchLatestTag(ctx, m.latestURL)
	if err == nil {
		st.Latest = latest
	}
	st.IsUpdateAvailable = isUpdateAvailable(st.Version, st.Latest)
	return st, nil
}

// ServiceUp starts dl infrastructure containers (traefik, portainer, mail).
func (m *Manager) ServiceUp(ctx context.Context) error {
	return m.service(ctx, "up")
}

// ServiceDown stops and removes dl infrastructure containers.
func (m *Manager) ServiceDown(ctx context.Context) error {
	return m.service(ctx, "down")
}

func (m *Manager) service(ctx context.Context, action string) error {
	result, err := m.Exec(ctx, m.root, []string{"service", action})
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("dl service %s: %s", action, resultMessage(result))
	}
	return nil
}

func resultMessage(result *Result) string {
	if result == nil {
		return "unknown error"
	}
	msg := strings.TrimSpace(result.Stderr)
	if msg == "" {
		msg = strings.TrimSpace(result.Stdout)
	}
	if msg == "" {
		return fmt.Sprintf("exit %d", result.ExitCode)
	}
	return msg
}

// Exec runs the managed dl binary. It never searches PATH for another dl.
func (m *Manager) Exec(ctx context.Context, workdir string, args []string) (*Result, error) {
	bin := m.binPath()
	if _, err := os.Stat(bin); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotInstalled
		}
		return nil, fmt.Errorf("dl: stat binary: %w", err)
	}
	if err := m.ensureDirs(); err != nil {
		return nil, err
	}

	dir := strings.TrimSpace(workdir)
	if dir == "" {
		dir = m.root
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = ""
	}

	parent := os.Environ()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = execEnv(m.homeDir(), filepath.Join(m.root, "bin"), dockerConfigDir(parent, userHome), parent)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			result.ExitCode = exit.ExitCode()
			return result, nil
		}
		return nil, fmt.Errorf("dl: exec: %w", err)
	}
	return result, nil
}

// Uninstall removes the isolated runtime directory. It never touches a system dl.
func (m *Manager) Uninstall() error {
	if err := os.RemoveAll(m.root); err != nil {
		return fmt.Errorf("dl: uninstall: %w", err)
	}
	return nil
}
