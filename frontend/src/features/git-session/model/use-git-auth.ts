import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { GitProvider } from '@/entities/git'
import { gitProviderApi } from '../api/provider'
import { gitAccountKeys } from './queries'

export function useGitAuth(provider: GitProvider) {
  const queryClient = useQueryClient()
  const queryKey = gitAccountKeys.provider(provider)
  const api = gitProviderApi[provider]

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: gitAccountKeys.all })

  const login = useMutation({
    mutationFn: () => api.login(),
    onSuccess: (account) => {
      queryClient.setQueryData(queryKey, account)
    },
    onSettled: invalidate,
  })

  const logout = useMutation({
    mutationFn: () => api.logout(),
    onSuccess: () => {
      queryClient.setQueryData(queryKey, null)
    },
    onSettled: invalidate,
  })

  const connectWithToken = useMutation({
    mutationFn: async (token: string) => {
      await api.setAccessToken(token)
      return api.profile()
    },
    onSuccess: (account) => {
      queryClient.setQueryData(queryKey, account)
    },
    onSettled: invalidate,
  })

  return { login, logout, connectWithToken }
}
