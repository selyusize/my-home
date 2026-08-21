import { useEffect, useState } from 'react'
import { Check, CircleOff, Copy, ExternalLink, Play, Square } from 'lucide-react'
import { DlIcon, type DlServiceView, type DlSnapshot } from '@/entities/dl'
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { IconButton } from '@/shared/ui/icon-button'
import { Spinner } from '@/shared/ui/spinner'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/shared/ui/tooltip'

export type DlSettingsDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  status: DlSnapshot
  serviceUpPending?: boolean
  serviceDownPending?: boolean
  updatePending?: boolean
  onServiceUp: () => void
  onServiceDown: () => void
  onUpdate: () => void
  repoUrl?: string
  onOpenRepo?: () => void
}

const SERVICE_LABELS: Record<string, string> = {
  traefik: 'Traefik',
  portainer: 'Portainer',
  mail: 'Mail',
}

function getServiceStatusLabel(service: DlServiceView) {
  if (service.running) {
    return 'запущен'
  }
  if (service.present) {
    return 'остановлен'
  }
  return 'нет'
}

function ServiceStatusIcon({ service }: { service: DlServiceView }) {
  const label = getServiceStatusLabel(service)
  const Icon = service.running ? Play : service.present ? Square : CircleOff

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span
          aria-label={label}
          className={cn(
            'inline-flex size-7 items-center justify-center',
            service.running && 'text-emerald-400',
            !service.running && service.present && 'text-amber-300',
            !service.present && 'text-white/40',
          )}
        >
          <Icon className="size-4" />
        </span>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="grid grid-cols-[88px_1fr] gap-2">
      <dt className="text-white/50">{label}</dt>
      <dd className="min-w-0 text-white">{value}</dd>
    </div>
  )
}

function PathDetail({ path }: { path: string }) {
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    if (!copied) {
      return
    }
    const id = window.setTimeout(() => setCopied(false), 1500)
    return () => window.clearTimeout(id)
  }, [copied])

  useEffect(() => {
    setCopied(false)
  }, [path])

  const copyPath = async () => {
    try {
      await navigator.clipboard.writeText(path)
      setCopied(true)
    } catch {
      setCopied(false)
    }
  }

  return (
    <div className="grid gap-1">
      <dt className="text-white/50">Путь</dt>
      <dd className="flex items-start gap-1">
        <span className="min-w-0 flex-1 break-all text-white">{path}</span>
        <IconButton
          label={copied ? 'Скопировано' : 'Копировать путь'}
          className="size-7 shrink-0"
          onClick={() => {
            void copyPath()
          }}
        >
          {copied ? <Check /> : <Copy />}
        </IconButton>
      </dd>
    </div>
  )
}

export function DlSettingsDialog({
  open,
  onOpenChange,
  status,
  serviceUpPending = false,
  serviceDownPending = false,
  updatePending = false,
  onServiceUp,
  onServiceDown,
  onUpdate,
  repoUrl,
  onOpenRepo,
}: DlSettingsDialogProps) {
  const serviceBusy = serviceUpPending || serviceDownPending
  const canControl = status.installed && status.dockerOk && !serviceBusy
  const version = status.version || '—'
  const latest = status.latest ? `доступно ${status.latest}` : ''
  const docker =
    status.dockerOk
      ? [status.dockerVersion, status.dockerOs].filter(Boolean).join(' · ') || 'доступен'
      : 'недоступен'

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="border-white/20 bg-black text-white sm:max-w-md">
        <DialogHeader>
          <div className="flex items-center gap-2">
            <DlIcon className="size-8 rounded-md" />
            <div className="flex min-w-0 items-center gap-1">
              <DialogTitle>dl</DialogTitle>
              {onOpenRepo ? (
                <IconButton
                  label="Открыть репозиторий"
                  className="size-7"
                  onClick={onOpenRepo}
                >
                  <ExternalLink />
                </IconButton>
              ) : null}
            </div>
          </div>
          <DialogDescription className="sr-only">
            Состояние изолированного dl и инфраструктурных сервисов
          </DialogDescription>
        </DialogHeader>
        <dl className="grid gap-1 text-sm">
          <Detail label="Статус" value={status.installed ? 'установлен' : 'не установлен'} />
          <Detail
            label="Версия"
            value={status.updateAvailable && latest ? `${version} → ${status.latest}` : version}
          />
          <Detail label="Docker" value={docker} />
          {repoUrl ? (
            <Detail label="Репозиторий" value={repoUrl.replace(/^https?:\/\//, '')} />
          ) : null}
          {status.path ? <PathDetail path={status.path} /> : null}
        </dl>
        <div className="grid gap-2">
          <p className="text-xs text-white/50">Сервисы</p>
          <ul className="grid gap-1.5">
            {status.services.map((service) => (
              <li key={service.name} className="flex items-center justify-between gap-2 text-sm">
                <span>{SERVICE_LABELS[service.name] ?? service.name}</span>
                <ServiceStatusIcon service={service} />
              </li>
            ))}
          </ul>
        </div>
        <DialogFooter className="border-white/20 bg-transparent sm:justify-end">
          <IconButton
            label="dl down"
            pending={serviceDownPending}
            disabled={!canControl || !status.serviceUp}
            className={
              canControl && status.serviceUp
                ? 'text-red-500 hover:bg-red-500/10 hover:text-red-400'
                : undefined
            }
            onClick={onServiceDown}
          >
            <Square />
          </IconButton>
          {status.updateAvailable ? (
            <Button
              type="button"
              variant="outline"
              className="cursor-pointer border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white"
              disabled={!status.installed || updatePending}
              onClick={onUpdate}
            >
              {updatePending ? <Spinner /> : null}
              Обновить
            </Button>
          ) : null}
          <IconButton
            label="dl up"
            pending={serviceUpPending}
            disabled={!canControl || status.serviceUp}
            className={
              canControl && !status.serviceUp
                ? 'text-emerald-500 hover:bg-emerald-500/10 hover:text-emerald-400'
                : undefined
            }
            onClick={onServiceUp}
          >
            <Play />
          </IconButton>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
