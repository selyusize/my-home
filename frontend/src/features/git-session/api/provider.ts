import * as GitHubService from '@bindings/github.com/selyusize/my-home/internal/git/githubservice'
import * as GitLabService from '@bindings/github.com/selyusize/my-home/internal/git/gitlabservice'
import type { ContributionCalendar } from '@bindings/github.com/selyusize/my-home/pkg/git/models'
import type { GitAccount, GitProvider } from '@/entities/git'

type GitProviderApi = {
  isAuthenticated: () => Promise<boolean>
  profile: () => Promise<GitAccount | null>
  contributionCalendar: () => Promise<ContributionCalendar | null>
  login: () => Promise<GitAccount | null>
  logout: () => Promise<void>
  setAccessToken: (token: string) => Promise<void>
}

export const gitProviderApi: Record<GitProvider, GitProviderApi> = {
  github: {
    isAuthenticated: GitHubService.IsAuthenticated,
    profile: GitHubService.Profile,
    contributionCalendar: GitHubService.ContributionCalendar,
    login: GitHubService.Login,
    logout: GitHubService.Logout,
    setAccessToken: GitHubService.SetAccessToken,
  },
  gitlab: {
    isAuthenticated: GitLabService.IsAuthenticated,
    profile: GitLabService.Profile,
    contributionCalendar: GitLabService.ContributionCalendar,
    login: GitLabService.Login,
    logout: GitLabService.Logout,
    setAccessToken: GitLabService.SetAccessToken,
  },
}
