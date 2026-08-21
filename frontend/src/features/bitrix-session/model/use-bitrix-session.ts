import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Open as openPage } from '@bindings/github.com/selyusize/my-home/internal/window/windowservice'
import { TimeManStatus } from '@bindings/github.com/selyusize/my-home/pkg/bitrix/models'
import {
  mapBitrixProfileToAccount,
  type BitrixTimeMan,
  type BitrixTimeManStatus,
} from '@/entities/bitrix'
import { getErrorMessage } from '@/shared/lib/error-message'
import type { AccountProfile } from '@/shared/ui/profile-dialog'
import { bitrixApi, type BitrixProfile } from '../api/bitrix'

export const bitrixKeys = {
  all: ['bitrix'] as const,
  profile: () => [...bitrixKeys.all, 'profile'] as const,
  timeMan: () => [...bitrixKeys.all, 'timeman'] as const,
}

async function loadBitrixProfile(): Promise<BitrixProfile | null> {
  const authenticated = await bitrixApi.isAuthenticated()
  if (!authenticated) {
    return null
  }
  return bitrixApi.profile()
}

export type BitrixSessionBindProps = {
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

export function useBitrixSession(): BitrixSessionBindProps {
  const queryClient = useQueryClient()
  const { data: account, isPending } = useQuery({
    queryKey: bitrixKeys.profile(),
    queryFn: loadBitrixProfile,
    retry: false,
    staleTime: 30_000,
  })
  const connected = Boolean(account)
  const { data: timeMan, dataUpdatedAt: timeManFetchedAt, refetch: refetchTimeMan } = useQuery({
    queryKey: bitrixKeys.timeMan(),
    queryFn: bitrixApi.timeMan,
    enabled: connected,
    retry: false,
    staleTime: 15_000,
  })

  const invalidate = () => queryClient.invalidateQueries({ queryKey: bitrixKeys.all })

  const connect = useMutation({
    mutationFn: ({ domain, webhook }: { domain: string; webhook: string }) =>
      bitrixApi.setWebhook(domain, webhook),
    onSuccess: (profile) => {
      queryClient.setQueryData(bitrixKeys.profile(), profile)
      toast.success('Bitrix24 подключён')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const logout = useMutation({
    mutationFn: bitrixApi.logout,
    onSuccess: () => {
      queryClient.setQueryData(bitrixKeys.profile(), null)
      queryClient.setQueryData(bitrixKeys.timeMan(), null)
      toast.success('Bitrix24 отключён')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const openDay = useMutation({
    mutationFn: bitrixApi.timeManOpen,
    onSuccess: () => {
      toast.success('Рабочий день начат')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: bitrixKeys.timeMan() }),
  })

  const pauseDay = useMutation({
    mutationFn: bitrixApi.timeManPause,
    onSuccess: () => {
      toast.success('Рабочее время остановлено')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: bitrixKeys.timeMan() }),
  })

  const closeDay = useMutation({
    mutationFn: bitrixApi.timeManClose,
    onSuccess: () => {
      toast.success('Рабочий день закончен')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey: bitrixKeys.timeMan() }),
  })

  const switching = logout.isPending

  return {
    connected,
    pending: switching,
    disabled: isPending,
    webhookPending: connect.isPending,
    timeManOpenPending: openDay.isPending,
    timeManPausePending: pauseDay.isPending,
    timeManClosePending: closeDay.isPending,
    timeManStatus: timeMan?.status ?? TimeManStatus.$zero,
    timeMan,
    timeManFetchedAt,
    profile: account ? mapBitrixProfileToAccount(account) : null,
    onDisconnect: () => {
      logout.mutate()
    },
    onSubmitWebhook: async (domain, webhook) => {
      await connect.mutateAsync({ domain, webhook })
    },
    onOpenPage: async (url) => {
      try {
        const title = ['Bitrix24', account?.name].filter(Boolean).join(' · ')
        await openPage(title, url)
      } catch (error) {
        toast.error(getErrorMessage(error))
      }
    },
    onOpenProfile: () => {
      void refetchTimeMan()
    },
    onTimeManOpen: () => {
      openDay.mutate()
    },
    onTimeManPause: () => {
      pauseDay.mutate()
    },
    onTimeManClose: () => {
      closeDay.mutate()
    },
  }
}
