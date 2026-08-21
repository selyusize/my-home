package window

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const maxPathLength = 4096

type WindowService struct{}

func NewWindowService() *WindowService {
	return &WindowService{}
}

func (s *WindowService) Open(title, rawURL string) error {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("window: invalid url")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("window: url must be http or https")
	}

	app := application.Get()
	if app == nil {
		return fmt.Errorf("window: app is not running")
	}

	pageURL := parsed.String()
	name := "page:" + parsed.Host + parsed.EscapedPath()
	if title == "" {
		title = parsed.Host
	}

	if win, ok := app.Window.GetByName(name); ok {
		win.SetTitle(title)
		win.SetURL(pageURL)
		win.Show()
		win.Focus()
		return nil
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             name,
		Title:            title,
		Width:            1080,
		Height:           760,
		URL:              pageURL,
		BackgroundColour: application.NewRGB(20, 20, 20),
	})
	return nil
}

func (s *WindowService) OpenInCursor(root string) error {
	dir, err := folderPath(root)
	if err != nil {
		return err
	}
	cmd, err := cursorCommand(dir)
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("window: open cursor: %w", err)
	}
	_ = cmd.Process.Release()
	return nil
}

func folderPath(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" || strings.ContainsRune(root, 0) {
		return "", fmt.Errorf("window: invalid path")
	}
	if utf8.RuneCountInString(root) > maxPathLength {
		return "", fmt.Errorf("window: path is too long")
	}
	if !filepath.IsAbs(root) {
		return "", fmt.Errorf("window: path must be absolute")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("window: abs: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("window: folder not found")
	}
	if !info.IsDir() {
		return "", fmt.Errorf("window: not a directory")
	}
	return abs, nil
}

func cursorCommand(dir string) (*exec.Cmd, error) {
	if bin, err := exec.LookPath("cursor"); err == nil {
		return exec.Command(bin, dir), nil
	}
	if runtime.GOOS == "darwin" {
		return exec.Command("open", "-a", "Cursor", dir), nil
	}
	return nil, fmt.Errorf("window: cursor is not installed")
}
