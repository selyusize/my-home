import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Open as openPage } from '@bindings/github.com/selyusize/my-home/internal/window/windowservice'
import { EMPTY_DL_SERVICES, mapDlStatusToView, type DlSnapshot } from '@/entities/dl'
import { getErrorMessage } from '@/shared/lib/error-message'
import { dlApi, type DlStatus } from '../api/dl'

export const dlKeys = {
  all: ['dl'] as const,
  status: () => [...dlKeys.all, 'status'] as const,
}

const DL_REPO_URL = 'https://github.com/local-deploy/dl'

function buildClearedDlStatus(prev?: DlStatus | null): DlStatus {
  return {
    path: prev?.path ?? '',
    installed: false,
    version: '',
    latest: prev?.latest ?? '',
    updateAvailable: false,
    dockerOk: prev?.dockerOk ?? false,
    dockerVersion: prev?.dockerVersion ?? '',
    dockerOs: prev?.dockerOs ?? '',
    serviceUp: false,
    services: EMPTY_DL_SERVICES,
  }
}

export type DlRuntimeBindProps = {
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

export function useDlRuntime(): DlRuntimeBindProps {
  const queryClient = useQueryClient()
  const { data } = useQuery({
    queryKey: dlKeys.status(),
    queryFn: dlApi.status,
    retry: false,
    staleTime: 15_000,
  })

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: dlKeys.all })

  const install = useMutation({
    mutationFn: dlApi.install,
    onSuccess: () => {
      toast.success('dl установлен')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const update = useMutation({
    mutationFn: dlApi.install,
    onSuccess: () => {
      toast.success('dl обновлён')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const uninstall = useMutation({
    mutationFn: dlApi.uninstall,
    onSuccess: () => {
      queryClient.setQueryData(dlKeys.status(), buildClearedDlStatus(data))
      toast.success('dl удалён')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const serviceUp = useMutation({
    mutationFn: dlApi.serviceUp,
    onSuccess: () => {
      toast.success('dl service запущен')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  const serviceDown = useMutation({
    mutationFn: dlApi.serviceDown,
    onSuccess: () => {
      toast.success('dl service остановлен')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
    onSettled: invalidate,
  })

  return {
    status: mapDlStatusToView(data),
    pending: install.isPending || uninstall.isPending,
    serviceUpPending: serviceUp.isPending,
    serviceDownPending: serviceDown.isPending,
    updatePending: update.isPending,
    onInstall: () => {
      install.mutate()
    },
    onUninstall: () => {
      uninstall.mutate()
    },
    onServiceUp: () => {
      serviceUp.mutate()
    },
    onServiceDown: () => {
      serviceDown.mutate()
    },
    onUpdate: () => {
      update.mutate()
    },
    repoUrl: DL_REPO_URL,
    onOpenRepo: () => {
      void openPage('dl', DL_REPO_URL).catch((error) => {
        toast.error(getErrorMessage(error))
      })
    },
  }
}
