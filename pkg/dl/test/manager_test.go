package dl_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/selyusize/my-home/pkg/dl"
)

func TestDefaultRoot(t *testing.T) {
	t.Parallel()

	root, err := dl.DefaultRoot()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(root, filepath.Join("my-home", "runtime", "dl")) {
		t.Fatalf("DefaultRoot=%q", root)
	}
}

func TestServiceUpDown(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mgr := dl.NewAt(root)
	writeFakeDL(t, root, "#!/bin/sh\nprintf '%s\\n' \"$*\"\n")

	if err := mgr.ServiceUp(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := mgr.ServiceDown(context.Background()); err != nil {
		t.Fatal(err)
	}

	result, err := mgr.Exec(context.Background(), root, []string{"service", "up"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Stdout, "service up") {
		t.Fatalf("stdout=%q", result.Stdout)
	}
}

func TestServiceUpNotInstalled(t *testing.T) {
	t.Parallel()

	mgr := dl.NewAt(t.TempDir())
	if err := mgr.ServiceUp(context.Background()); !errors.Is(err, dl.ErrNotInstalled) {
		t.Fatalf("err=%v", err)
	}
}

func TestServiceUpNonZero(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mgr := dl.NewAt(root)
	writeFakeDL(t, root, "#!/bin/sh\necho boom >&2\nexit 2\n")

	err := mgr.ServiceUp(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("err=%v", err)
	}
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mgr := dl.NewAt(root)
	writeFakeDL(t, root, "#!/bin/sh\n")

	if err := mgr.Uninstall(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("runtime still exists: %v", err)
	}
	_, err := mgr.Exec(context.Background(), "", []string{"version"})
	if !errors.Is(err, dl.ErrNotInstalled) {
		t.Fatalf("err=%v", err)
	}
}

func TestExecUsesIsolatedEnv(t *testing.T) {
	root := t.TempDir()
	mgr := dl.NewAt(root)
	writeFakeDL(t, root, "#!/bin/sh\nprintf 'HOME=%s\\nPATH=%s\\nDOCKER_CONFIG=%s\\nPWD=%s\\n' \"$HOME\" \"$PATH\" \"$DOCKER_CONFIG\" \"$PWD\"\n")

	t.Setenv("DOCKER_CONFIG", "/tmp/custom-docker")

	workdir := t.TempDir()
	result, err := mgr.Exec(context.Background(), workdir, []string{"version"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s", result.ExitCode, result.Stderr)
	}

	home := filepath.Join(root, "home")
	if !strings.Contains(result.Stdout, "HOME="+home) {
		t.Fatalf("stdout=%q want HOME=%s", result.Stdout, home)
	}
	if !strings.Contains(result.Stdout, "PATH="+filepath.Join(root, "bin")) {
		t.Fatalf("stdout=%q want PATH prefix", result.Stdout)
	}
	if !strings.Contains(result.Stdout, "DOCKER_CONFIG=/tmp/custom-docker") {
		t.Fatalf("stdout=%q want DOCKER_CONFIG", result.Stdout)
	}
	gotPWD := ""
	for _, line := range strings.Split(result.Stdout, "\n") {
		if rest, ok := strings.CutPrefix(line, "PWD="); ok {
			gotPWD = rest
			break
		}
	}
	wantPWD, err := filepath.EvalSymlinks(workdir)
	if err != nil {
		wantPWD = workdir
	}
	gotPWDEval, err := filepath.EvalSymlinks(gotPWD)
	if err != nil {
		gotPWDEval = gotPWD
	}
	if gotPWDEval != wantPWD {
		t.Fatalf("PWD=%q want %q", gotPWDEval, wantPWD)
	}
	if strings.Contains(result.Stdout, "deploy") {
		t.Fatalf("must not invoke deploy")
	}
}

func TestInstallAndVersion(t *testing.T) {
	t.Parallel()

	goos, arch, ok := supportedPlatform()
	if !ok {
		t.Skip("unsupported platform")
	}

	script := "#!/bin/sh\necho 'dl version test'\n"
	tarball, err := gzipTar("dl", []byte(script))
	if err != nil {
		t.Fatal(err)
	}

	tag := "1.2.3"
	assetName := fmt.Sprintf("dl-%s-%s-%s.tar.gz", tag, goos, arch)
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": tag,
			"assets": []map[string]string{{
				"name":                 assetName,
				"browser_download_url": srv.URL + "/asset",
			}},
		})
	})
	mux.HandleFunc("/asset", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(tarball)
	})

	root := t.TempDir()
	mgr := dl.NewAt(root, dl.WithReleaseURL(srv.URL+"/releases/latest"))
	if err := mgr.Install(context.Background()); err != nil {
		t.Fatal(err)
	}

	st, err := mgr.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !st.IsInstalled {
		t.Fatal("expected installed")
	}
	if st.Latest != tag {
		t.Fatalf("latest=%q", st.Latest)
	}
	if !strings.Contains(st.Version, "dl version test") {
		t.Fatalf("version=%q", st.Version)
	}

	result, err := mgr.Exec(context.Background(), root, []string{"version"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("version exit=%d", result.ExitCode)
	}
}

func TestSmokeVersionAndUp(t *testing.T) {
	if os.Getenv("MYHOME_DL_SMOKE") != "1" {
		t.Skip("set MYHOME_DL_SMOKE=1 to run live dl version/up")
	}
	project := os.Getenv("MYHOME_DL_PROJECT")
	mgr, err := dl.New()
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.Install(context.Background()); err != nil {
		t.Fatal(err)
	}
	ver, err := mgr.Exec(context.Background(), "", []string{"version"})
	if err != nil {
		t.Fatal(err)
	}
	if ver.ExitCode != 0 {
		t.Fatalf("version exit=%d stderr=%s", ver.ExitCode, ver.Stderr)
	}
	if project == "" {
		t.Log("MYHOME_DL_PROJECT unset; skipped up")
		return
	}
	up, err := mgr.Exec(context.Background(), project, []string{"up"})
	if err != nil {
		t.Fatal(err)
	}
	if up.ExitCode != 0 {
		t.Fatalf("up exit=%d stderr=%s stdout=%s", up.ExitCode, up.Stderr, up.Stdout)
	}
}

func writeFakeDL(t *testing.T, root, script string) {
	t.Helper()
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(binDir, "dl"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}

func supportedPlatform() (goos, arch string, ok bool) {
	switch runtime.GOOS {
	case "darwin", "linux":
		goos = runtime.GOOS
	default:
		return "", "", false
	}
	switch runtime.GOARCH {
	case "amd64", "arm64":
		arch = runtime.GOARCH
	default:
		return "", "", false
	}
	return goos, arch, true
}

func gzipTar(name string, body []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	hdr := &tar.Header{
		Name: name,
		Mode: 0o755,
		Size: int64(len(body)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write(body); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
