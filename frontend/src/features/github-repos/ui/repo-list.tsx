import type { GitRepository } from '@/entities/git'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/ui/avatar'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/shared/ui/tooltip'
import { FolderGit2Icon, LockIcon } from 'lucide-react'
import { formatRelativeTime } from '../model/format-relative-time'

type RepoListProps = {
  repos: GitRepository[]
  localFullNames?: ReadonlySet<string>
  hasMore: boolean
  loadingMore: boolean
  onOpen?: (url: string) => void
  onLoadMore: () => void
}

export function RepoList({
  repos,
  localFullNames,
  hasMore,
  loadingMore,
  onOpen,
  onLoadMore,
}: RepoListProps) {
  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <ul className="min-h-0 flex-1 overflow-y-auto" aria-label="Репозитории GitHub">
        {repos.map((repo) => {
          const pushed = formatRelativeTime(repo.pushedAt)
          const canOpen = Boolean(repo.htmlUrl && onOpen)
          const isLocal = localFullNames?.has(repo.fullName.toLowerCase()) ?? false
          return (
            <li key={repo.id} className="border-b border-white/10 last:border-b-0">
              {canOpen ? (
                <button
                  type="button"
                  onClick={() => onOpen?.(repo.htmlUrl)}
                  className="flex min-h-10 w-full cursor-pointer items-center gap-2.5 px-1 py-2 text-left outline-none transition-colors hover:bg-white/5 focus-visible:ring-2 focus-visible:ring-ring/50"
                >
                  <RepoRow repo={repo} pushed={pushed} isLocal={isLocal} />
                </button>
              ) : (
                <div className="flex min-h-10 items-center gap-2.5 px-1 py-2">
                  <RepoRow repo={repo} pushed={pushed} isLocal={isLocal} />
                </div>
              )}
            </li>
          )
        })}
      </ul>
      {hasMore ? (
        <Button
          type="button"
          variant="ghost"
          className="mt-2 cursor-pointer text-white/70 hover:bg-white/5 hover:text-white"
          disabled={loadingMore}
          onClick={onLoadMore}
        >
          {loadingMore ? 'Загрузка…' : 'Показать ещё'}
        </Button>
      ) : null}
    </div>
  )
}

function RepoRow({
  repo,
  pushed,
  isLocal,
}: {
  repo: GitRepository
  pushed: string
  isLocal: boolean
}) {
  return (
    <>
      <Avatar size="sm" className="shrink-0">
        <AvatarImage src={repo.ownerAvatar} alt="" width={24} height={24} />
        <AvatarFallback>{repo.ownerLogin.slice(0, 2).toUpperCase() || '?'}</AvatarFallback>
      </Avatar>
      <span className="min-w-0 flex-1">
        <span className="flex min-w-0 items-center gap-1.5">
          <span className="truncate text-sm font-medium text-white" translate="no">
            {repo.name}
          </span>
          {isLocal ? <LocalCloneMark /> : null}
          {repo.private ? (
            <LockIcon className="size-3.5 shrink-0 text-white/50" aria-hidden="true" />
          ) : null}
        </span>
        <span className="mt-0.5 block truncate text-xs tabular-nums text-white/50">
          {[repo.ownerLogin, repo.language, pushed].filter(Boolean).join(' · ')}
        </span>
      </span>
    </>
  )
}

function LocalCloneMark() {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span
          className="inline-flex shrink-0"
          aria-label="Есть локальная копия"
        >
          <FolderGit2Icon className="size-3.5 text-emerald-400" aria-hidden="true" />
        </span>
      </TooltipTrigger>
      <TooltipContent>Есть локальная копия</TooltipContent>
    </Tooltip>
  )
}

export function RepoListSkeleton() {
  return (
    <div className="space-y-2" aria-hidden="true">
      {Array.from({ length: 8 }, (_, index) => (
        <Skeleton key={index} className="h-10 rounded-md" />
      ))}
    </div>
  )
}
