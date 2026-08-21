package dl

import "errors"

var (
	ErrNotInstalled        = errors.New("dl: not installed")
	ErrUnsupportedPlatform = errors.New("dl: unsupported platform")
	ErrMissingReleaseAsset = errors.New("dl: release asset not found")
)
