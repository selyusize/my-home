import { useId, useState } from 'react'
import { KeyRound } from 'lucide-react'
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
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Spinner } from '@/shared/ui/spinner'

export type BitrixWebhookDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  pending?: boolean
  showTrigger?: boolean
  triggerDisabled?: boolean
  onSubmit: (domain: string, webhook: string) => Promise<void>
}

export function BitrixWebhookDialog({
  open,
  onOpenChange,
  pending = false,
  showTrigger = false,
  triggerDisabled = false,
  onSubmit,
}: BitrixWebhookDialogProps) {
  const domainId = useId()
  const webhookId = useId()
  const [domain, setDomain] = useState('')
  const [webhook, setWebhook] = useState('')

  const closeDialog = () => {
    setDomain('')
    setWebhook('')
    onOpenChange(false)
  }

  const submitWebhook = async () => {
    try {
      await onSubmit(domain.trim(), webhook.trim())
      closeDialog()
    } catch {
      // caller reports the error; keep the dialog open
    }
  }

  return (
    <>
      {showTrigger ? (
        <IconButton
          label="Вебхук Bitrix24"
          disabled={triggerDisabled}
          onClick={() => onOpenChange(true)}
        >
          <KeyRound />
        </IconButton>
      ) : null}
      <Dialog
        open={open}
        onOpenChange={(next) => {
          onOpenChange(next)
          if (!next) {
            setDomain('')
            setWebhook('')
          }
        }}
      >
        <DialogContent className="border-white/20 bg-black sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Вебхук Bitrix24</DialogTitle>
            <DialogDescription>
              Портал и входящий вебхук. Нужны права user и timeman.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-3">
            <div className="grid gap-2">
              <Label htmlFor={domainId}>Портал</Label>
              <Input
                id={domainId}
                autoComplete="off"
                placeholder="company.bitrix24.ru"
                value={domain}
                onChange={(event) => setDomain(event.target.value)}
                disabled={pending}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor={webhookId}>Вебхук</Label>
              <Input
                id={webhookId}
                type="password"
                autoComplete="off"
                placeholder="userId/code или URL"
                value={webhook}
                onChange={(event) => setWebhook(event.target.value)}
                disabled={pending}
              />
            </div>
          </div>
          <DialogFooter className="border-white/20 bg-transparent">
            <Button
              type="button"
              variant="outline"
              className="cursor-pointer border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white"
              disabled={pending}
              onClick={closeDialog}
            >
              Отмена
            </Button>
            <Button
              type="button"
              className="cursor-pointer bg-white text-black hover:bg-white/90"
              disabled={pending || webhook.trim().length === 0 || (domain.trim().length === 0 && !webhook.includes('://'))}
              onClick={() => void submitWebhook()}
            >
              {pending ? <Spinner className="text-black" /> : null}
              Подключить
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
