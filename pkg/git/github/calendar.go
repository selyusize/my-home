package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/selyusize/my-home/pkg/git"
)

const contributionCalendarQuery = `
query {
  viewer {
    contributionsCollection {
      contributionCalendar {
        totalContributions
        weeks {
          contributionDays {
            date
            contributionCount
            contributionLevel
          }
        }
      }
    }
  }
}
`

type graphqlCalendarResponse struct {
	Data struct {
		Viewer struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
					Weeks              []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
							ContributionLevel string `json:"contributionLevel"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"viewer"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// ContributionCalendar returns the authenticated user's GitHub contribution heatmap.
func (c *Client) ContributionCalendar(ctx context.Context) (*git.ContributionCalendar, error) {
	client, err := c.api(ctx)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(ctx, http.MethodPost, "graphql", map[string]string{
		"query": contributionCalendarQuery,
	})
	if err != nil {
		return nil, fmt.Errorf("github: contribution calendar: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	var out graphqlCalendarResponse
	if _, err := client.Do(req, &out); err != nil {
		return nil, fmt.Errorf("github: contribution calendar: %w", err)
	}
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("github: contribution calendar: %s", out.Errors[0].Message)
	}

	cal := out.Data.Viewer.ContributionsCollection.ContributionCalendar
	days := make([]git.ContributionDay, 0, 371)
	for _, week := range cal.Weeks {
		for _, day := range week.ContributionDays {
			if day.Date == "" {
				continue
			}
			days = append(days, git.ContributionDay{
				Date:  day.Date,
				Count: day.ContributionCount,
				Level: githubContributionLevel(day.ContributionLevel),
			})
		}
	}

	return &git.ContributionCalendar{
		Total: cal.TotalContributions,
		Days:  days,
	}, nil
}

func githubContributionLevel(level string) int {
	switch level {
	case "FIRST_QUARTILE":
		return 1
	case "SECOND_QUARTILE":
		return 2
	case "THIRD_QUARTILE":
		return 3
	case "FOURTH_QUARTILE":
		return 4
	default:
		return 0
	}
}
