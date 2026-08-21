import type {
  Branch,
  Commit,
  Repository,
  User,
} from '@bindings/github.com/selyusize/my-home/pkg/git/models'

export type GitProvider = 'github' | 'gitlab'
export type GitAccount = User
export type GitRepository = Repository
export type GitBranch = Branch
export type GitCommit = Commit

export type GitContributionDay = {
  date: string
  count: number
  level: number
}

export type GitContributionCalendar = {
  total: number
  days: GitContributionDay[]
}

export const gitProviders: Record<GitProvider, { label: string; host: string }> = {
  github: { label: 'GitHub', host: 'github.com' },
  gitlab: { label: 'GitLab', host: 'gitlab.com' },
}
