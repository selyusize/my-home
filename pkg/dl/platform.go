package dl

import (
	"fmt"
	"runtime"
	"strings"
)

func releasePlatform() (goos, arch string, err error) {
	switch runtime.GOOS {
	case "darwin":
		goos = "darwin"
	case "linux":
		goos = "linux"
	default:
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedPlatform, runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedPlatform, runtime.GOARCH)
	}

	return goos, arch, nil
}

func tarballNames(tag, goos, arch string) []string {
	tag = strings.TrimSpace(tag)
	names := []string{fmt.Sprintf("dl-%s-%s-%s.tar.gz", tag, goos, arch)}
	trimmed := strings.TrimPrefix(tag, "v")
	if trimmed != tag {
		names = append(names, fmt.Sprintf("dl-%s-%s-%s.tar.gz", trimmed, goos, arch))
	}
	return names
}
