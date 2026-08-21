import { BitrixSession } from '@/features/bitrix-session'
import { DlRuntime } from '@/features/dl-runtime'
import { GitHubSession, GitLabSession } from '@/features/git-session'
import { RepoSettings } from '@/features/repo-settings'

export function UserPage() {
  return (
    <div className="flex min-h-0 flex-1 gap-6 p-6">
      <div className="min-w-0 flex-1">
        <RepoSettings />
      </div>
      <div className="flex shrink-0 flex-col items-end gap-1">
        <GitHubSession />
        <GitLabSession />
        <BitrixSession />
        <DlRuntime />
      </div>
    </div>
  )
}
