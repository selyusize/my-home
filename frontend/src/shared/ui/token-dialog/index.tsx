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

type TokenDialogProps = {
  label: string
  disabled?: boolean
  pending?: boolean
  onSubmit: (token: string) => Promise<void>
}

export function TokenDialog({
  label,
  disabled,
  pending = false,
  onSubmit,
}: TokenDialogProps) {
  const inputId = useId()
  const [open, setOpen] = useState(false)
  const [token, setToken] = useState('')

  const close = () => {
    setToken('')
    setOpen(false)
  }

  const handleSubmit = async () => {
    try {
      await onSubmit(token.trim())
      close()
    } catch {
      // caller reports the error; keep the dialog open
    }
  }

  return (
    <>
      <IconButton
        label={`Токен ${label}`}
        disabled={disabled}
        onClick={() => setOpen(true)}
      >
        <KeyRound />
      </IconButton>
      <Dialog
        open={open}
        onOpenChange={(next) => {
          setOpen(next)
          if (!next) {
            setToken('')
          }
        }}
      >
        <DialogContent className="border-white/20 bg-black sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Токен {label}</DialogTitle>
            <DialogDescription>Вставьте access token.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-2">
            <Label htmlFor={inputId}>Токен</Label>
            <Input
              id={inputId}
              type="password"
              autoComplete="off"
              value={token}
              onChange={(event) => setToken(event.target.value)}
              disabled={pending}
            />
          </div>
          <DialogFooter className="border-white/20 bg-transparent">
            <Button
              type="button"
              variant="outline"
              className="cursor-pointer border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white"
              disabled={pending}
              onClick={close}
            >
              Отмена
            </Button>
            <Button
              type="button"
              className="cursor-pointer bg-white text-black hover:bg-white/90"
              disabled={pending || token.trim().length === 0}
              onClick={() => void handleSubmit()}
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
