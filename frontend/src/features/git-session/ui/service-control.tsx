import { useState, type ComponentType, type ReactNode, type SVGProps } from 'react'
import {
  ContributionCalendar,
  type GitContributionCalendar,
} from '@/entities/git'
import { ConfirmLogoutButton } from '@/shared/ui/confirm-logout'
import { IconButton } from '@/shared/ui/icon-button'
import { ProfileDialog, type AccountProfile } from '@/shared/ui/profile-dialog'
import { StatusBadge } from '@/shared/ui/status-badge'
import { TokenDialog } from '@/shared/ui/token-dialog'

export type ServiceControlIcon = ComponentType<SVGProps<SVGSVGElement>>

export type ServiceControlProps = {
  label: string
  icon: ServiceControlIcon
  connected: boolean
  pending?: boolean
  disabled?: boolean
  tokenPending?: boolean
  profile?: AccountProfile | null
  calendar?: GitContributionCalendar | null
  onConnect: () => void
  onDisconnect: () => void
  onSubmitToken?: (token: string) => Promise<void>
  onOpenPage?: (url: string) => void
}

export function ServiceControl({
  label,
  icon: Icon,
  connected,
  pending = false,
  disabled = false,
  tokenPending = false,
  profile,
  calendar,
  onConnect,
  onDisconnect,
  onSubmitToken,
  onOpenPage,
}: ServiceControlProps) {
  const [profileOpen, setProfileOpen] = useState(false)

  const handleClick = () => {
    if (!connected) {
      onConnect()
      return
    }
    if (profile) {
      setProfileOpen(true)
    }
  }

  let action: ReactNode = null
  if (connected) {
    action = (
      <ConfirmLogoutButton
        label={label}
        pending={pending}
        disabled={disabled}
        onDisconnect={onDisconnect}
      />
    )
  } else if (onSubmitToken) {
    action = (
      <TokenDialog
        label={label}
        disabled={disabled}
        pending={tokenPending}
        onSubmit={onSubmitToken}
      />
    )
  }

  return (
    <div className="flex items-center gap-1">
      <IconButton
        label={connected ? `Профиль ${label}` : `Подключить ${label}`}
        pending={pending}
        disabled={disabled && !pending}
        onClick={handleClick}
        badge={<StatusBadge connected={connected} />}
      >
        <Icon />
      </IconButton>
      {action}
      {profile ? (
        <ProfileDialog
          label={label}
          profile={profile}
          open={profileOpen}
          onOpenChange={setProfileOpen}
          onOpenPage={onOpenPage}
          wide={Boolean(calendar?.days?.length)}
        >
          {calendar?.days?.length ? (
            <ContributionCalendar calendar={calendar} />
          ) : null}
        </ProfileDialog>
      ) : null}
    </div>
  )
}
