package git

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	pkggit "github.com/selyusize/my-home/pkg/git"
)

const (
	oauthTimeout = 3 * time.Minute
	// defaultLoopbackRedirect is registered on the GitHub/GitLab OAuth app
	// when GITHUB_REDIRECT_URL / GITLAB_REDIRECT_URL are unset.
	defaultLoopbackRedirect = "http://127.0.0.1:8741/oauth/callback"
)

func loopbackRedirect(raw string) string {
	if v := strings.TrimSpace(raw); v != "" {
		return v
	}
	return defaultLoopbackRedirect
}

func loginWithBrowser(ctx context.Context, client pkggit.Client, redirectURL, providerName string) (*pkggit.User, error) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, oauthTimeout)
		defer cancel()
	}

	callback, err := parseLoopbackRedirect(redirectURL)
	if err != nil {
		return nil, err
	}

	state, err := randomOAuthState()
	if err != nil {
		return nil, err
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(callback.Path, func(w http.ResponseWriter, r *http.Request) {
		handleOAuthCallback(w, r, state, providerName, codeCh, errCh)
	})

	listener, err := net.Listen("tcp", callback.Host)
	if err != nil {
		return nil, fmt.Errorf("oauth: listen on %s: %w", callback.Host, err)
	}

	server := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	go func() {
		if serveErr := server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			select {
			case errCh <- fmt.Errorf("oauth: callback server: %w", serveErr):
			default:
			}
		}
	}()

	authURL, err := client.AuthURL(state)
	if err != nil {
		return nil, err
	}
	if err := openBrowser(authURL); err != nil {
		return nil, fmt.Errorf("oauth: open browser: %w", err)
	}

	select {
	case code := <-codeCh:
		if err := client.ExchangeCode(ctx, code); err != nil {
			return nil, err
		}
		return client.Profile(ctx)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("oauth: authorization timed out or was cancelled")
	}
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request, state, providerName string, codeCh chan<- string, errCh chan<- error) {
	query := r.URL.Query()
	if got := query.Get("state"); got != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		sendOAuthErr(errCh, fmt.Errorf("oauth: invalid state"))
		return
	}
	if errCode := query.Get("error"); errCode != "" {
		desc := query.Get("error_description")
		if desc == "" {
			desc = errCode
		}
		writeOAuthPage(w, http.StatusBadRequest, providerName, false)
		sendOAuthErr(errCh, fmt.Errorf("oauth: %s", desc))
		return
	}

	code := query.Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		sendOAuthErr(errCh, fmt.Errorf("oauth: missing authorization code"))
		return
	}

	writeOAuthPage(w, http.StatusOK, providerName, true)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	select {
	case codeCh <- code:
	default:
	}
}

func parseLoopbackRedirect(raw string) (*url.URL, error) {
	raw = loopbackRedirect(raw)

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("oauth: invalid redirect url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("oauth: redirect url must be http or https")
	}

	host := strings.ToLower(parsed.Hostname())
	if host != "127.0.0.1" && host != "localhost" {
		return nil, fmt.Errorf("oauth: redirect url must use localhost")
	}
	if parsed.Port() == "" {
		return nil, fmt.Errorf("oauth: redirect url must include a port")
	}
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed, nil
}

func randomOAuthState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("oauth: generate state: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}

func sendOAuthErr(errCh chan<- error, err error) {
	select {
	case errCh <- err:
	default:
	}
}

func writeOAuthPage(w http.ResponseWriter, status int, providerName string, ok bool) {
	title := "Авторизация не завершена"
	body := "Можно закрыть это окно и вернуться в приложение."
	if ok {
		title = html.EscapeString(providerName) + " подключён"
		body = "Можно закрыть это окно и вернуться в приложение."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <style>
    html,body{height:100%%;margin:0;background:#141414;color:#f0f0f0;font:16px/1.5 -apple-system,BlinkMacSystemFont,sans-serif}
    main{min-height:100%%;display:grid;place-items:center}
    section{padding:32px 36px;border:1px solid rgba(255,255,255,.12);border-radius:28px;background:rgba(255,255,255,.04)}
    h1{margin:0 0 8px;font-size:22px;font-weight:600}
    p{margin:0;color:rgba(240,240,240,.7)}
  </style>
</head>
<body>
  <main>
    <section>
      <h1>%s</h1>
      <p>%s</p>
    </section>
  </main>
</body>
</html>`, title, title, body)
}
