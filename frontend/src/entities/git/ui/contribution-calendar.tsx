import { cn } from '@/shared/lib/utils'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/shared/ui/tooltip'
import type { GitContributionCalendar, GitContributionDay } from '../model/types'

const LEVEL_CLASS = [
  'bg-white/10',
  'bg-[#0e4429]',
  'bg-[#006d32]',
  'bg-[#26a641]',
  'bg-[#39d353]',
] as const

const WEEKDAY_LABELS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'] as const
const MONTHS = [
  'янв',
  'фев',
  'мар',
  'апр',
  'май',
  'июн',
  'июл',
  'авг',
  'сен',
  'окт',
  'ноя',
  'дек',
] as const

const CELL = 10
const GAP = 3

type Week = {
  start: Date
  days: (GitContributionDay | undefined)[]
}

type ContributionCalendarProps = {
  calendar: GitContributionCalendar
}

function parseIsoDate(value: string): Date {
  const [year, month, day] = value.split('-').map(Number)
  return new Date(Date.UTC(year, month - 1, day))
}

function formatIsoDate(date: Date): string {
  return date.toISOString().slice(0, 10)
}

function addUtcDays(date: Date, days: number): Date {
  const next = new Date(date)
  next.setUTCDate(next.getUTCDate() + days)
  return next
}

function getMondayIndex(date: Date): number {
  return (date.getUTCDay() + 6) % 7
}

function getStartOfMonday(date: Date): Date {
  return addUtcDays(date, -getMondayIndex(date))
}

function getMonthName(date: Date): string {
  return MONTHS[date.getUTCMonth()]
}

function getContributionWord(count: number): string {
  const n = Math.abs(count) % 100
  const digit = n % 10
  if (n > 10 && n < 20) {
    return 'вкладов'
  }
  if (digit === 1) {
    return 'вклад'
  }
  if (digit >= 2 && digit <= 4) {
    return 'вклада'
  }
  return 'вкладов'
}

function formatContributionDayLabel(day: GitContributionDay): string {
  const date = new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(parseIsoDate(day.date))
  if (day.count === 0) {
    return `Нет вкладов · ${date}`
  }
  const count = new Intl.NumberFormat('ru-RU').format(day.count)
  return `${count} ${getContributionWord(day.count)} · ${date}`
}

function groupDaysIntoWeeks(days: GitContributionDay[]): Week[] {
  if (days.length === 0) {
    return []
  }

  const byDate = new Map(days.map((day) => [day.date, day]))
  const first = parseIsoDate(days[0].date)
  const last = parseIsoDate(days[days.length - 1].date)
  const start = getStartOfMonday(first)

  const weeks: Week[] = []
  for (let cursor = start; cursor.getTime() <= last.getTime(); ) {
    const weekStart = cursor
    const weekDays: (GitContributionDay | undefined)[] = []
    for (let weekday = 0; weekday < 7; weekday += 1) {
      weekDays.push(byDate.get(formatIsoDate(cursor)))
      cursor = addUtcDays(cursor, 1)
    }
    weeks.push({ start: weekStart, days: weekDays })
  }
  return weeks
}

function groupWeeksByMonth(weeks: Week[]) {
  const groups: { label: string; span: number }[] = []
  for (const week of weeks) {
    const label = getMonthName(week.start)
    const last = groups[groups.length - 1]
    if (last?.label === label) {
      last.span += 1
      continue
    }
    groups.push({ label, span: 1 })
  }
  return groups
}

function getContributionLevelClass(level: number): string {
  return LEVEL_CLASS[Math.min(Math.max(level, 0), LEVEL_CLASS.length - 1)] ?? LEVEL_CLASS[0]
}

function DayCell({ day }: { day: GitContributionDay }) {
  const squareClass = cn('size-full rounded-[2px]', getContributionLevelClass(day.level))

  if (day.count <= 0) {
    return (
      <td className="p-0" style={{ width: CELL, height: CELL }}>
        <div className={squareClass} />
      </td>
    )
  }

  const label = formatContributionDayLabel(day)

  return (
    <td className="p-0" style={{ width: CELL, height: CELL }}>
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type="button"
            aria-label={label}
            className={cn(squareClass, 'block cursor-default border-0 p-0')}
          />
        </TooltipTrigger>
        <TooltipContent className="z-[100]">{label}</TooltipContent>
      </Tooltip>
    </td>
  )
}

export function ContributionCalendar({ calendar }: ContributionCalendarProps) {
  const weeks = groupDaysIntoWeeks(calendar.days)
  if (weeks.length === 0) {
    return null
  }

  const total = new Intl.NumberFormat('ru-RU').format(calendar.total)
  const summary = `${total} ${getContributionWord(calendar.total)} за последний год`
  const months = groupWeeksByMonth(weeks)

  return (
    <div className="grid gap-2 mt-2">
      <p className="text-xs text-white/60">{summary}</p>
      <div className="overflow-x-auto">
        <table
          role="img"
          aria-label={summary}
          className="border-separate text-[10px] leading-none text-white/50"
          style={{ borderSpacing: GAP }}
        >
          <thead>
            <tr>
              <td className="w-7 p-0" />
              {months.map((month, index) => (
                <td
                  key={`${month.label}-${index}`}
                  colSpan={month.span}
                  className="relative h-3 overflow-hidden p-0 align-bottom"
                >
                  <span className="absolute bottom-0 left-0 whitespace-nowrap">
                    {month.label}
                  </span>
                </td>
              ))}
            </tr>
          </thead>
          <tbody>
            {WEEKDAY_LABELS.map((weekday, dayIndex) => (
              <tr key={weekday} className="h-2.5">
                <td className="w-7 p-0 pr-1 text-right align-middle">{weekday}</td>
                {weeks.map((week) => {
                  const day = week.days[dayIndex]
                  const weekKey = formatIsoDate(week.start)
                  if (!day) {
                    return (
                      <td
                        key={`${weekKey}-${dayIndex}`}
                        className="p-0"
                        style={{ width: CELL, height: CELL }}
                      />
                    )
                  }
                  return <DayCell key={day.date} day={day} />
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="flex items-center justify-end gap-1 text-[10px] text-white/50">
        <span>Меньше</span>
        {LEVEL_CLASS.map((cls) => (
          <span key={cls} className={cn('size-2.5 rounded-[2px]', cls)} />
        ))}
        <span>Больше</span>
      </div>
    </div>
  )
}
