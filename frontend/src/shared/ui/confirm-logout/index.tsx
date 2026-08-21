import { LogOut } from 'lucide-react'
import { useState } from 'react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/ui/alert-dialog'
import { IconButton } from '@/shared/ui/icon-button'

type ConfirmLogoutButtonProps = {
  label: string
  pending?: boolean
  disabled?: boolean
  onDisconnect: () => void
}

export function ConfirmLogoutButton({
  label,
  pending = false,
  disabled = false,
  onDisconnect,
}: ConfirmLogoutButtonProps) {
  const [confirmOpen, setConfirmOpen] = useState(false)

  return (
    <>
      <IconButton
        label={`Выйти из ${label}`}
        pending={pending}
        disabled={disabled && !pending}
        onClick={() => setConfirmOpen(true)}
      >
        <LogOut />
      </IconButton>
      <AlertDialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <AlertDialogContent className="border-white/20 bg-black text-white">
          <AlertDialogHeader>
            <AlertDialogTitle>Выйти из {label}?</AlertDialogTitle>
            <AlertDialogDescription>
              Вы уверены, что хотите выйти?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel className="cursor-pointer border-white/20 bg-transparent text-white hover:bg-white/10 hover:text-white">
              Отмена
            </AlertDialogCancel>
            <AlertDialogAction
              className="cursor-pointer bg-white text-black hover:bg-white/90"
              onClick={onDisconnect}
            >
              Выйти
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
