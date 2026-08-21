import { useQuery, useQueryClient } from '@tanstack/react-query'
import { localReposApi } from '../api/local-repos'
import { cacheLocalClones, type LocalClonesResult } from './cache-local-clones'
import { hasRepoRoot } from './has-repo-root'
import { repoSettingsKeys } from './keys'

export type { LocalClonesResult }

export function useLocalClones() {
  const queryClient = useQueryClient()

  return useQuery({
    queryKey: repoSettingsKeys.explorer(),
    queryFn: async (): Promise<LocalClonesResult> => {
      const settings = await localReposApi.getSettings()
      if (!hasRepoRoot(settings)) {
        return { isConfigured: false, clones: [] }
      }

      await localReposApi.check(settings)
      const clones = await localReposApi.listClones()
      cacheLocalClones(queryClient, clones)
      return { isConfigured: true, clones }
    },
    retry: false,
  })
}
