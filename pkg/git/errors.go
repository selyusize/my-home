package git

import "errors"

var (
	ErrNotAuthenticated   = errors.New("git: not authenticated")
	ErrMissingOAuthConfig = errors.New("git: oauth client id and secret are required")
	ErrUnknownProvider    = errors.New("git: unknown provider")
)
