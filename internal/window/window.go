package window

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

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
