import { useState } from 'react'
import { Pause, Play, Square } from 'lucide-react'
import { BitrixIcon, type BitrixTimeMan, type BitrixTimeManStatus } from '@/entities/bitrix'
import { ConfirmLogoutButton } from '@/shared/ui/confirm-logout'
import { IconButton } from '@/shared/ui/icon-button'
import { ProfileDialog, type AccountProfile } from '@/shared/ui/profile-dialog'
import { StatusBadge } from '@/shared/ui/status-badge'
import { BitrixWebhookDialog } from './bitrix-webhook-dialog'
import { TimeManStats } from './timeman-stats'

export type BitrixControlProps = {
  connected: boolean
  pending?: boolean
  disabled?: boolean
  webhookPending?: boolean
  timeManOpenPending?: boolean
  timeManPausePending?: boolean
  timeManClosePending?: boolean
  timeManStatus?: BitrixTimeManStatus
  timeMan?: BitrixTimeMan | null
  timeManFetchedAt?: number
  profile?: AccountProfile | null
  onDisconnect: () => void
  onSubmitWebhook: (domain: string, webhook: string) => Promise<void>
  onOpenPage?: (url: string) => void
  onOpenProfile?: () => void
  onTimeManOpen: () => void
  onTimeManPause: () => void
  onTimeManClose: () => void
}

export function BitrixControl({
  connected,
  pending = false,
  disabled = false,
  webhookPending = false,
  timeManOpenPending = false,
  timeManPausePending = false,
  timeManClosePending = false,
  timeManStatus = '',
  timeMan,
  timeManFetchedAt = 0,
  profile,
  onDisconnect,
  onSubmitWebhook,
  onOpenPage,
  onOpenProfile,
  onTimeManOpen,
  onTimeManPause,
  onTimeManClose,
}: BitrixControlProps) {
  const [webhookOpen, setWebhookOpen] = useState(false)
  const [profileOpen, setProfileOpen] = useState(false)

  const timeManBusy = timeManOpenPending || timeManPausePending || timeManClosePending
  const busy = disabled || pending || webhookPending || timeManBusy
  const canStart =
    connected &&
    (timeManStatus === 'closed' || timeManStatus === 'paused' || timeManStatus === 'expired')
  const canPause = connected && timeManStatus === 'opened'
  const canClose =
    connected &&
    (timeManStatus === 'opened' ||
      timeManStatus === 'paused' ||
      timeManStatus === 'expired')

  const handleIconClick = () => {
    if (!connected) {
      setWebhookOpen(true)
      return
    }
    if (profile) {
      onOpenProfile?.()
      setProfileOpen(true)
    }
  }

  return (
    <div className="flex items-center gap-1">
      <IconButton
        label={connected ? 'Профиль Bitrix24' : 'Подключить Bitrix24'}
        pending={pending}
        disabled={busy && !pending}
        onClick={handleIconClick}
        badge={<StatusBadge connected={connected} />}
      >
        <BitrixIcon />
      </IconButton>
      {connected ? (
        <>
          <IconButton
            label="Начать"
            pending={timeManOpenPending}
            disabled={!canStart || busy}
            className={
              canStart && !busy
                ? 'text-emerald-500 hover:bg-emerald-500/10 hover:text-emerald-400'
                : undefined
            }
            onClick={onTimeManOpen}
          >
            <Play />
          </IconButton>
          <IconButton
            label="Остановить"
            pending={timeManPausePending}
            disabled={!canPause || busy}
            className={
              canPause && !busy
                ? 'text-amber-400 hover:bg-amber-500/10 hover:text-amber-300'
                : undefined
            }
            onClick={onTimeManPause}
          >
            <Pause />
          </IconButton>
          <IconButton
            label="Закончить день"
            pending={timeManClosePending}
            disabled={!canClose || busy}
            className={
              canClose && !busy
                ? 'text-red-500 hover:bg-red-500/10 hover:text-red-400'
                : undefined
            }
            onClick={onTimeManClose}
          >
            <Square />
          </IconButton>
          <ConfirmLogoutButton
            label="Bitrix24"
            pending={pending}
            disabled={busy && !pending}
            onDisconnect={onDisconnect}
          />
        </>
      ) : (
        <BitrixWebhookDialog
          showTrigger
          triggerDisabled={busy}
          open={webhookOpen}
          onOpenChange={setWebhookOpen}
          pending={webhookPending}
          onSubmit={onSubmitWebhook}
        />
      )}
      {profile ? (
        <ProfileDialog
          label="Bitrix24"
          profile={profile}
          open={profileOpen}
          onOpenChange={setProfileOpen}
          onOpenPage={onOpenPage}
        >
          <TimeManStats snapshot={timeMan} fetchedAt={timeManFetchedAt} active={profileOpen} />
        </ProfileDialog>
      ) : null}
    </div>
  )
}
