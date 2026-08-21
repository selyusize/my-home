import { toast } from 'sonner'
import { Open as openPage } from '@bindings/github.com/selyusize/my-home/internal/window/windowservice'
import {
  gitProviders,
  mapGitAccountToProfile,
  type GitContributionCalendar,
  type GitProvider,
} from '@/entities/git'
import { getErrorMessage } from '@/shared/lib/error-message'
import type { AccountProfile } from '@/shared/ui/profile-dialog'
import { useGitAccount, useGitCalendar } from './queries'
import { useGitAuth } from './use-git-auth'

export type GitSessionBindProps = {
  connected: boolean
  loading?: boolean
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

export function useGitSession(provider: GitProvider): GitSessionBindProps {
  const { data: account, isPending } = useGitAccount(provider)
  const { data: calendar } = useGitCalendar(provider, Boolean(account))
  const { login, logout, connectWithToken } = useGitAuth(provider)
  const meta = gitProviders[provider]
  const connected = Boolean(account)
  const switching = login.isPending || logout.isPending
  const busy = isPending || switching || connectWithToken.isPending

  const onConnect = async () => {
    try {
      await login.mutateAsync()
      toast.success(`${meta.label} подключён`)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const onDisconnect = async () => {
    try {
      await logout.mutateAsync()
      toast.success(`${meta.label} отключён`)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const onSubmitToken = async (token: string) => {
    try {
      await connectWithToken.mutateAsync(token)
      toast.success(`${meta.label} подключён`)
    } catch (error) {
      toast.error(getErrorMessage(error))
      throw error
    }
  }

  const onOpenPage = async (url: string) => {
    try {
      const title = [meta.label, account?.login].filter(Boolean).join(' · ')
      await openPage(title, url)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return {
    connected,
    loading: isPending,
    pending: switching,
    disabled: busy,
    tokenPending: connectWithToken.isPending,
    profile: account ? mapGitAccountToProfile(account) : null,
    calendar,
    onConnect,
    onDisconnect,
    onSubmitToken,
    onOpenPage,
  }
}
