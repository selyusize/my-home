import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { getErrorMessage } from '@/shared/lib/error-message'
import {
  emptyRepoSettings,
  localReposApi,
  type LocalRepoSettings,
} from '../api/local-repos'
import { cacheLocalClones } from './cache-local-clones'
import { repoSettingsKeys } from './keys'

export { repoSettingsKeys }

export function useRepoSettings() {
  const queryClient = useQueryClient()
  const settingsQuery = useQuery({
    queryKey: repoSettingsKeys.settings(),
    queryFn: localReposApi.getSettings,
    retry: false,
    staleTime: Infinity,
  })

  const save = useMutation({
    mutationFn: localReposApi.saveSettings,
    onMutate: async (settings) => {
      await queryClient.cancelQueries({ queryKey: repoSettingsKeys.settings() })
      const previous = queryClient.getQueryData<LocalRepoSettings>(
        repoSettingsKeys.settings(),
      )
      queryClient.setQueryData(repoSettingsKeys.settings(), settings)
      return { previous }
    },
    onError: (error, _settings, context) => {
      if (context?.previous) {
        queryClient.setQueryData(repoSettingsKeys.settings(), context.previous)
      }
      toast.error(getErrorMessage(error))
    },
  })

  const check = useMutation({
    mutationFn: localReposApi.check,
    onSuccess: async (report, settings) => {
      queryClient.setQueryData(repoSettingsKeys.settings(), settings)
      const clones = await localReposApi.listClones()
      cacheLocalClones(queryClient, clones)
      const found = [
        report.github && `GitHub: ${report.github}`,
        report.gitlab && `GitLab: ${report.gitlab}`,
      ]
        .filter(Boolean)
        .join(', ')
      toast.success(found ? `Найдены локальные копии · ${found}` : 'Локальных копий не найдено')
    },
    onError: (error) => {
      toast.error(getErrorMessage(error))
    },
  })

  const pickDirectory = async (title: string) => {
    try {
      return await localReposApi.selectDirectory(title)
    } catch (error) {
      toast.error(getErrorMessage(error))
      return ''
    }
  }

  return {
    settings: settingsQuery.data ?? emptyRepoSettings,
    isPending: settingsQuery.isPending,
    save: save.mutate,
    saving: save.isPending,
    check: check.mutate,
    checking: check.isPending,
    pickDirectory,
  }
}

export function useLocalCloneNames(provider: 'github' | 'gitlab') {
  return useQuery({
    queryKey: repoSettingsKeys.clones(provider),
    queryFn: async () => {
      const clones = await localReposApi.listClones(provider)
      return new Set(clones.map((clone) => clone.fullName.toLowerCase()))
    },
    retry: false,
    staleTime: Infinity,
  })
}
