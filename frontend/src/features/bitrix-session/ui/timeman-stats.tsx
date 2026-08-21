import {
  formatDayStart,
  formatDurationClock,
  liveTimeMan,
  type BitrixTimeMan,
} from '@/entities/bitrix'
import { useNow } from '../model/use-now'

type TimeManStatsProps = {
  snapshot: BitrixTimeMan | null | undefined
  fetchedAt: number
  active: boolean
}

const statusLabel: Record<string, string> = {
  opened: 'Работаю',
  paused: 'Перерыв',
  closed: 'День закрыт',
  expired: 'День не закрыт',
}

export function TimeManStats({ snapshot, fetchedAt, active }: TimeManStatsProps) {
  const ticking = active && (snapshot?.status === 'opened' || snapshot?.status === 'paused' || snapshot?.status === 'expired')
  const now = useNow(ticking)
  const live = liveTimeMan(snapshot, fetchedAt, ticking ? now : fetchedAt || Date.now())
  if (!live) {
    return null
  }

  const rows = [
    live.startedAt ? ['Начал', formatDayStart(live.startedAt)] : null,
    ['В работе', formatDurationClock(live.workSeconds)],
    ['Перерыв', formatDurationClock(live.breakSeconds)],
  ].filter((row): row is [string, string] => row !== null)

  return (
    <section aria-label="Рабочий день" className="grid gap-2 border-t border-white/10 pt-3">
      <p className="text-xs text-white/50">
        {statusLabel[snapshot?.status ?? ''] ?? 'Рабочий день'}
        {live.isLive ? (
          <span
            aria-hidden="true"
            className={
              snapshot?.status === 'paused'
                ? 'ml-2 inline-block size-1.5 rounded-full bg-amber-400 align-middle'
                : 'ml-2 inline-block size-1.5 rounded-full bg-emerald-400 align-middle'
            }
          />
        ) : null}
      </p>
      <dl className="grid gap-1 text-sm">
        {rows.map(([key, value]) => (
          <div key={key} className="grid grid-cols-[88px_1fr] gap-2">
            <dt className="text-white/50">{key}</dt>
            <dd className="font-medium tabular-nums tracking-tight">{value}</dd>
          </div>
        ))}
      </dl>
    </section>
  )
}
