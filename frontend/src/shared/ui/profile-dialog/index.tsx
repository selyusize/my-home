import type { ReactNode } from 'react'
import { ExternalLink } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { IconButton } from '@/shared/ui/icon-button'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { getInitials } from './get-initials'
import type { AccountProfile } from './types'

export type { AccountProfile, AccountProfileStat } from './types'

type ProfileDialogProps = {
  label: string
  profile: AccountProfile
  open: boolean
  onOpenChange: (open: boolean) => void
  onOpenPage?: (url: string) => void
  wide?: boolean
  children?: ReactNode
}

export function ProfileDialog({
  label,
  profile,
  open,
  onOpenChange,
  onOpenPage,
  wide = false,
  children,
}: ProfileDialogProps) {
  const details = [
    profile.email ? ['Email', profile.email] : null,
    profile.company ? ['Компания', profile.company] : null,
    profile.location ? ['Локация', profile.location] : null,
    profile.website ? ['Сайт', profile.website] : null,
  ].filter((row): row is [string, string] => row !== null)

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className={cn(
          'border-white/20 bg-black text-white',
          wide ? 'sm:max-w-3xl' : 'sm:max-w-md',
        )}
      >
        <DialogHeader>
          <DialogTitle>{label}</DialogTitle>
          <DialogDescription className="sr-only">
            Профиль {label}
          </DialogDescription>
        </DialogHeader>
        <div className="flex items-start gap-3">
          <Avatar size="lg" className="size-12">
            {profile.avatarUrl ? (
              <AvatarImage src={profile.avatarUrl} alt="" />
            ) : null}
            <AvatarFallback className="bg-white/10 text-sm text-white">
              {getInitials(profile.name)}
            </AvatarFallback>
          </Avatar>
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-1">
              <p className="truncate text-sm font-medium">{profile.name}</p>
              {profile.pageUrl && onOpenPage ? (
                <IconButton
                  label={`Открыть ${label}`}
                  className="size-7"
                  onClick={() => {
                    if (profile.pageUrl) {
                      onOpenPage(profile.pageUrl)
                    }
                  }}
                >
                  <ExternalLink />
                </IconButton>
              ) : null}
            </div>
            {profile.handle ? (
              <p className="truncate text-xs text-white/60">@{profile.handle}</p>
            ) : null}
          </div>
        </div>
        {profile.bio ? (
          <p className="text-sm text-white/80">{profile.bio}</p>
        ) : null}
        {details.length > 0 ? (
          <dl className="grid gap-1 text-sm">
            {details.map(([key, value]) => (
              <div key={key} className="grid grid-cols-[88px_1fr] gap-2">
                <dt className="text-white/50">{key}</dt>
                <dd className="min-w-0 truncate text-white">{value}</dd>
              </div>
            ))}
          </dl>
        ) : null}
        {profile.stats && profile.stats.length > 0 ? (
          <div className="flex gap-4 text-sm">
            {profile.stats.map((stat) => (
              <div key={stat.label}>
                <div className="font-medium">{stat.value}</div>
                <div className="text-xs text-white/50">{stat.label}</div>
              </div>
            ))}
          </div>
        ) : null}
        {children}
      </DialogContent>
    </Dialog>
  )
}
