import type { QueryClient } from '@tanstack/react-query'
import type { LocalClone } from '../api/local-repos'
import { repoSettingsKeys } from './keys'

const gitProviders = ['github', 'gitlab'] as const

export type LocalClonesResult = {
  isConfigured: boolean
  clones: LocalClone[]
}

export function cacheLocalClones(queryClient: QueryClient, clones: LocalClone[]) {
  queryClient.setQueryData<LocalClonesResult>(repoSettingsKeys.explorer(), {
    isConfigured: true,
    clones,
  })
  for (const provider of gitProviders) {
    queryClient.setQueryData(
      repoSettingsKeys.clones(provider),
      new Set(
        clones
          .filter((clone) => clone.provider === provider)
          .map((clone) => clone.fullName.toLowerCase()),
      ),
    )
  }
}
