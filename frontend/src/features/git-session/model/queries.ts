import { useQuery } from '@tanstack/react-query'
import {
  mapGitCalendar,
  type GitAccount,
  type GitContributionCalendar,
  type GitProvider,
} from '@/entities/git'
import { gitProviderApi } from '../api/provider'

export const gitAccountKeys = {
  all: ['git-account'] as const,
  provider: (provider: GitProvider) => [...gitAccountKeys.all, provider] as const,
  calendar: (provider: GitProvider) =>
    [...gitAccountKeys.provider(provider), 'calendar'] as const,
}

async function loadGitAccount(provider: GitProvider): Promise<GitAccount | null> {
  const api = gitProviderApi[provider]
  const authenticated = await api.isAuthenticated()
  if (!authenticated) {
    return null
  }
  return api.profile()
}

async function loadGitCalendar(
  provider: GitProvider,
): Promise<GitContributionCalendar | null> {
  const api = gitProviderApi[provider]
  const authenticated = await api.isAuthenticated()
  if (!authenticated) {
    return null
  }
  try {
    return mapGitCalendar(await api.contributionCalendar())
  } catch {
    return null
  }
}

export function useGitAccount(provider: GitProvider) {
  return useQuery({
    queryKey: gitAccountKeys.provider(provider),
    queryFn: () => loadGitAccount(provider),
    retry: false,
    staleTime: 30_000,
  })
}

export function useGitCalendar(provider: GitProvider, enabled: boolean) {
  return useQuery({
    queryKey: gitAccountKeys.calendar(provider),
    queryFn: () => loadGitCalendar(provider),
    enabled,
    retry: false,
    staleTime: 60_000,
  })
}
