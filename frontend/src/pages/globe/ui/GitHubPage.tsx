import { GitHubIcon } from '@/entities/git'
import { GitHubRepos } from '@/features/github-repos'
import { useGitSession } from '@/features/git-session'
import { useLocalCloneNames } from '@/features/repo-settings'
import { Button } from '@/shared/ui/button'
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/shared/ui/empty'
import { Spinner } from '@/shared/ui/spinner'
import { TokenDialog } from '@/shared/ui/token-dialog'

export function GitHubPage() {
  const session = useGitSession('github')
  const localClones = useLocalCloneNames('github')

  if (session.loading) {
    return (
      <div className="flex min-h-0 flex-1 items-center justify-center">
        <Spinner className="size-6 text-white/60" />
      </div>
    )
  }

  if (!session.connected) {
    return (
      <Empty className="flex min-h-0 flex-1 justify-center">
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <GitHubIcon aria-hidden="true" />
          </EmptyMedia>
          <EmptyTitle>GitHub не подключён</EmptyTitle>
          <EmptyDescription>
            Подключите аккаунт, чтобы увидеть репозитории.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <div className="flex items-center gap-2">
            <Button
              type="button"
              className="cursor-pointer bg-white text-black hover:bg-white/90"
              disabled={session.disabled}
              onClick={() => void session.onConnect()}
            >
              {session.pending ? <Spinner className="text-black" /> : null}
              Подключить GitHub
            </Button>
            {session.onSubmitToken ? (
              <TokenDialog
                label="GitHub"
                disabled={session.disabled}
                pending={session.tokenPending}
                onSubmit={session.onSubmitToken}
              />
            ) : null}
          </div>
        </EmptyContent>
      </Empty>
    )
  }

  return <GitHubRepos onOpenPage={session.onOpenPage} localFullNames={localClones.data} />
}
