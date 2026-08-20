package github

import (
	"time"

	gh "github.com/google/go-github/v90/github"
)

func deref[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

func timeOrZero(t *gh.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.Time
}
