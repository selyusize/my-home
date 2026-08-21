package dl

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 2 * time.Minute}
}

func githubGet(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "my-home")
	req.Header.Set("Accept", "application/vnd.github+json")
	return httpClient().Do(req)
}

func fetchLatestTag(ctx context.Context, latestURL string) (string, error) {
	release, err := fetchLatestRelease(ctx, latestURL)
	if err != nil {
		return "", err
	}
	return release.TagName, nil
}

func fetchLatestRelease(ctx context.Context, latestURL string) (*githubRelease, error) {
	resp, err := githubGet(ctx, latestURL)
	if err != nil {
		return nil, fmt.Errorf("dl: latest release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dl: latest release: unexpected status %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("dl: latest release: %w", err)
	}
	if strings.TrimSpace(release.TagName) == "" {
		return nil, fmt.Errorf("dl: latest release: empty tag")
	}
	return &release, nil
}

func assetURL(release *githubRelease, goos, arch string) (string, error) {
	wanted := make(map[string]struct{})
	for _, name := range tarballNames(release.TagName, goos, arch) {
		wanted[name] = struct{}{}
	}
	for _, asset := range release.Assets {
		if _, ok := wanted[asset.Name]; ok && asset.BrowserDownloadURL != "" {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrMissingReleaseAsset, release.TagName)
}

// Install downloads the latest official dl release into the isolated runtime.
func (m *Manager) Install(ctx context.Context) error {
	goos, arch, err := releasePlatform()
	if err != nil {
		return err
	}
	if err := m.ensureDirs(); err != nil {
		return err
	}

	release, err := fetchLatestRelease(ctx, m.latestURL)
	if err != nil {
		return err
	}
	url, err := assetURL(release, goos, arch)
	if err != nil {
		return err
	}

	resp, err := githubGet(ctx, url)
	if err != nil {
		return fmt.Errorf("dl: download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dl: download: unexpected status %s", resp.Status)
	}

	tmp, err := os.CreateTemp(filepath.Join(m.root, "bin"), "dl-*.tmp")
	if err != nil {
		return fmt.Errorf("dl: install temp: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if err := extractBinary(resp.Body, tmp); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o755); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("dl: chmod: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("dl: close temp: %w", err)
	}

	dest := m.binPath()
	if err := os.Rename(tmpName, dest); err != nil {
		return fmt.Errorf("dl: replace binary: %w", err)
	}
	clearQuarantine(dest)
	return nil
}

func extractBinary(src io.Reader, dest *os.File) error {
	gr, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("dl: gzip: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("dl: tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if filepath.Base(filepath.ToSlash(hdr.Name)) != "dl" {
			continue
		}
		if _, err := io.Copy(dest, tr); err != nil {
			return fmt.Errorf("dl: extract: %w", err)
		}
		return nil
	}
	return fmt.Errorf("dl: archive has no dl binary")
}

func clearQuarantine(path string) {
	if runtime.GOOS != "darwin" {
		return
	}
	_ = exec.Command("xattr", "-d", "com.apple.quarantine", path).Run()
}
