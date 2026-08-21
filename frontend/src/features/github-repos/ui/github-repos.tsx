import { useMemo, useState } from 'react'
import { getErrorMessage } from '@/shared/lib/error-message'
import { Button } from '@/shared/ui/button'
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from '@/shared/ui/empty'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { filterRepos } from '../model/filter-repos'
import { useGitHubRepos } from '../model/queries'
import { useDebouncedValue } from '../model/use-debounced-value'
import { RepoList, RepoListSkeleton } from './repo-list'

type GitHubReposProps = {
  onOpenPage?: (url: string) => void
  localFullNames?: ReadonlySet<string>
}

export function GitHubRepos({ onOpenPage, localFullNames }: GitHubReposProps) {
  const [query, setQuery] = useState('')
  const search = useDebouncedValue(query)
  const reposQuery = useGitHubRepos()
  const repos = useMemo(
    () => reposQuery.data?.pages.flatMap((page) => page.items) ?? [],
    [reposQuery.data],
  )
  const visibleRepos = useMemo(() => filterRepos(repos, search), [repos, search])

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-3">
      <div className="grid shrink-0 gap-1.5">
        <Label htmlFor="github-repo-search">Поиск</Label>
        <Input
          id="github-repo-search"
          name="repository-search"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          autoComplete="off"
          spellCheck={false}
          placeholder="my-home…"
        />
      </div>
      <RepoListBody
        isPending={reposQuery.isPending}
        isError={reposQuery.isError}
        error={reposQuery.error}
        search={search}
        visibleRepos={visibleRepos}
        localFullNames={localFullNames}
        hasMore={Boolean(reposQuery.hasNextPage)}
        loadingMore={reposQuery.isFetchingNextPage}
        onOpen={onOpenPage}
        onRetry={() => void reposQuery.refetch()}
        onLoadMore={() => void reposQuery.fetchNextPage()}
      />
    </div>
  )
}

function RepoListBody({
  isPending,
  isError,
  error,
  search,
  visibleRepos,
  localFullNames,
  hasMore,
  loadingMore,
  onOpen,
  onRetry,
  onLoadMore,
}: {
  isPending: boolean
  isError: boolean
  error: unknown
  search: string
  visibleRepos: ReturnType<typeof filterRepos>
  localFullNames?: ReadonlySet<string>
  hasMore: boolean
  loadingMore: boolean
  onOpen?: (url: string) => void
  onRetry: () => void
  onLoadMore: () => void
}) {
  if (isPending) {
    return <RepoListSkeleton />
  }
  if (isError) {
    return (
      <Empty className="flex-1">
        <EmptyHeader>
          <EmptyTitle>Не удалось загрузить репозитории</EmptyTitle>
          <EmptyDescription>
            {getErrorMessage(error)}. Попробуйте ещё раз.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button type="button" className="cursor-pointer" onClick={onRetry}>
            Повторить
          </Button>
        </EmptyContent>
      </Empty>
    )
  }
  if (visibleRepos.length === 0) {
    return (
      <Empty className="flex-1">
        <EmptyHeader>
          <EmptyTitle>{search ? 'Ничего не найдено' : 'Репозиториев нет'}</EmptyTitle>
          <EmptyDescription>
            {search
              ? 'Измените запрос или загрузите ещё репозитории.'
              : 'У этой учётной записи пока нет репозиториев.'}
          </EmptyDescription>
        </EmptyHeader>
      </Empty>
    )
  }
  return (
    <RepoList
      repos={visibleRepos}
      localFullNames={localFullNames}
      hasMore={hasMore}
      loadingMore={loadingMore}
      onOpen={onOpen}
      onLoadMore={onLoadMore}
    />
  )
}
