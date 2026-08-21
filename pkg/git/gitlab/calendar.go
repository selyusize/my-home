package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/selyusize/my-home/pkg/git"
)

// ContributionCalendar returns the authenticated user's GitLab contribution heatmap.
func (c *Client) ContributionCalendar(ctx context.Context) (*git.ContributionCalendar, error) {
	profile, err := c.Profile(ctx)
	if err != nil {
		return nil, err
	}
	if profile.Login == "" {
		return &git.ContributionCalendar{}, nil
	}

	base := c.cfg.BaseURL
	if base == "" {
		base = defaultBaseURL
	}

	calendarURL, err := url.JoinPath(base, "users", profile.Login, "calendar.json")
	if err != nil {
		return nil, fmt.Errorf("gitlab: contribution calendar: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, calendarURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gitlab: contribution calendar: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken())
	req.Header.Set("Private-Token", c.AccessToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab: contribution calendar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return &git.ContributionCalendar{}, nil
	}

	var raw map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("gitlab: contribution calendar: %w", err)
	}

	counts := make(map[string]int, len(raw))
	for date, count := range raw {
		counts[date] = int(count)
	}

	cal := git.BuildContributionCalendar(counts, time.Now())
	return &cal, nil
}
