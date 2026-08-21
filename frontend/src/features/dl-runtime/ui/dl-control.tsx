import { useState } from 'react'
import { Download, Trash2 } from 'lucide-react'
import { DlIcon, type DlSnapshot } from '@/entities/dl'
import { IconButton } from '@/shared/ui/icon-button'
import { StatusBadge } from '@/shared/ui/status-badge'
import { DlSettingsDialog } from './dl-settings-dialog'

export type DlControlProps = {
  status: DlSnapshot
  pending?: boolean
  serviceUpPending?: boolean
  serviceDownPending?: boolean
  updatePending?: boolean
  onInstall: () => void
  onUninstall: () => void
  onServiceUp: () => void
  onServiceDown: () => void
  onUpdate: () => void
  repoUrl?: string
  onOpenRepo?: () => void
}

export function DlControl({
  status,
  pending = false,
  serviceUpPending = false,
  serviceDownPending = false,
  updatePending = false,
  onInstall,
  onUninstall,
  onServiceUp,
  onServiceDown,
  onUpdate,
  repoUrl,
  onOpenRepo,
}: DlControlProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="flex items-center gap-1">
      <IconButton
        label="Настройки dl"
        onClick={() => setOpen(true)}
        badge={<StatusBadge connected={status.installed} />}
      >
        <DlIcon />
      </IconButton>
      <IconButton
        label={status.installed ? 'Удалить dl' : 'Установить dl'}
        pending={pending}
        onClick={status.installed ? onUninstall : onInstall}
      >
        {status.installed ? <Trash2 /> : <Download />}
      </IconButton>
      <DlSettingsDialog
        open={open}
        onOpenChange={setOpen}
        status={status}
        serviceUpPending={serviceUpPending}
        serviceDownPending={serviceDownPending}
        updatePending={updatePending}
        onServiceUp={onServiceUp}
        onServiceDown={onServiceDown}
        onUpdate={onUpdate}
        repoUrl={repoUrl}
        onOpenRepo={onOpenRepo}
      />
    </div>
  )
}
