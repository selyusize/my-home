const relativeTime = new Intl.RelativeTimeFormat('ru', { numeric: 'auto' })

const UNITS = [
  { unit: 'year' as const, seconds: 60 * 60 * 24 * 365 },
  { unit: 'month' as const, seconds: 60 * 60 * 24 * 30 },
  { unit: 'day' as const, seconds: 60 * 60 * 24 },
  { unit: 'hour' as const, seconds: 60 * 60 },
  { unit: 'minute' as const, seconds: 60 },
  { unit: 'second' as const, seconds: 1 },
]

export function formatRelativeTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getFullYear() < 1970) {
    return ''
  }

  const diffSeconds = Math.round((date.getTime() - Date.now()) / 1000)
  const unit = UNITS.find(({ seconds }) => Math.abs(diffSeconds) >= seconds) ?? UNITS[UNITS.length - 1]

  return relativeTime.format(Math.round(diffSeconds / unit.seconds), unit.unit)
}
