package git

import "time"

// ContributionLevel maps a daily count onto GitHub's 0–4 intensity scale.
func ContributionLevel(count int) int {
	switch {
	case count <= 0:
		return 0
	case count <= 2:
		return 1
	case count <= 5:
		return 2
	case count <= 8:
		return 3
	default:
		return 4
	}
}

// BuildContributionCalendar expands sparse date→count data into a Sunday-aligned year grid.
func BuildContributionCalendar(counts map[string]int, now time.Time) ContributionCalendar {
	now = now.UTC().Truncate(24 * time.Hour)
	start := now.AddDate(0, 0, -364)
	for start.Weekday() != time.Sunday {
		start = start.AddDate(0, 0, -1)
	}

	days := make([]ContributionDay, 0, 371)
	total := 0
	for day := start; !day.After(now); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		count := counts[key]
		total += count
		days = append(days, ContributionDay{
			Date:  key,
			Count: count,
			Level: ContributionLevel(count),
		})
	}

	return ContributionCalendar{Total: total, Days: days}
}
