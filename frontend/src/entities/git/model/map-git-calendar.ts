import type { GitContributionCalendar, GitContributionDay } from './types'

type GitCalendarInput = {
  total: number
  days?: GitContributionDay[] | null
}

export function mapGitCalendar(
  calendar?: GitCalendarInput | null,
): GitContributionCalendar | null {
  const days = calendar?.days
  if (!calendar || !days?.length) {
    return null
  }
  return { total: calendar.total, days }
}
