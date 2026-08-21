import type { LocalRepoSettings } from '../api/local-repos'

export function hasRepoRoot(settings: LocalRepoSettings) {
  return Boolean(
    settings.sharedPath.trim() ||
      (settings.githubSeparate && settings.githubPath.trim()) ||
      (settings.gitlabSeparate && settings.gitlabPath.trim()),
  )
}
