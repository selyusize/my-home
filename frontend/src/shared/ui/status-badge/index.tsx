import { AvatarBadge } from '@/shared/ui/avatar'

type StatusBadgeProps = {
  connected: boolean
}

export function StatusBadge({ connected }: StatusBadgeProps) {
  return (
    <AvatarBadge
      aria-hidden
      className={connected ? 'size-2 bg-emerald-500' : 'size-2 bg-red-500'}
    />
  )
}
