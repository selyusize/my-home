import type { GitRepository } from '@/entities/git'

export function filterRepos(repos: GitRepository[], query: string) {
  const needle = query.trim().toLowerCase()
  if (!needle) {
    return repos
  }

  return repos.filter((repo) =>
    [repo.name, repo.fullName, repo.description, repo.ownerLogin].some((value) =>
      value.toLowerCase().includes(needle),
    ),
  )
}
