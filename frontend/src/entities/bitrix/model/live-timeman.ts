import type { BitrixTimeMan, BitrixTimeManStatus } from './types'

export type LiveTimeMan = {
  workSeconds: number
  breakSeconds: number
  startedAt: Date | null
  isLive: boolean
}

function parseStart(value?: string): Date | null {
  const raw = value?.trim()
  if (!raw) {
    return null
  }
  const parsed = new Date(raw)
  if (Number.isNaN(parsed.getTime())) {
    return null
  }
  return parsed
}

export function liveTimeMan(
  snapshot: BitrixTimeMan | null | undefined,
  fetchedAt: number,
  now: number,
): LiveTimeMan | null {
  if (!snapshot?.status) {
    return null
  }

  const startedAt = parseStart(snapshot.timeStart)
  const duration = Math.max(0, snapshot.duration ?? 0)
  const leaks = Math.max(0, snapshot.timeLeaks ?? 0)
  const sinceFetch = Math.max(0, Math.floor((now - fetchedAt) / 1000))
  const elapsedFromStart = startedAt
    ? Math.max(0, Math.floor((now - startedAt.getTime()) / 1000))
    : 0

  const status: BitrixTimeManStatus = snapshot.status
  if (status === 'closed') {
    return { workSeconds: duration, breakSeconds: leaks, startedAt, isLive: false }
  }

  if (status === 'paused') {
    const breakSeconds = leaks + sinceFetch
    const workSeconds =
      duration > 0 ? duration : Math.max(0, elapsedFromStart - breakSeconds)
    return { workSeconds, breakSeconds, startedAt, isLive: true }
  }

  if (status === 'opened' || status === 'expired') {
    const workSeconds =
      duration > 0 ? duration + sinceFetch : Math.max(0, elapsedFromStart - leaks)
    return { workSeconds, breakSeconds: leaks, startedAt, isLive: true }
  }

  return null
}

export function formatDurationClock(totalSeconds: number): string {
  const seconds = Math.max(0, Math.floor(totalSeconds))
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const rest = seconds % 60
  return [hours, minutes, rest]
    .map((part) => String(part).padStart(2, '0'))
    .join(':')
}

export function formatDayStart(date: Date): string {
  return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}
