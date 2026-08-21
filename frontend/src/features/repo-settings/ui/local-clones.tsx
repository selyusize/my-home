import { Open as openPage } from '@bindings/github.com/selyusize/my-home/internal/window/windowservice'
import { GitHubIcon, GitLabIcon, gitProviders } from '@/entities/git'
import { getErrorMessage } from '@/shared/lib/error-message'
import { Button } from '@/shared/ui/button'
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/shared/ui/empty'
import { IconButton } from '@/shared/ui/icon-button'
import { Skeleton } from '@/shared/ui/skeleton'
import { FolderGit2Icon } from 'lucide-react'
import { Link } from 'react-router-dom'
import { toast } from 'sonner'
import type { LocalClone } from '../api/local-repos'
import { useLocalClones } from '../model/use-local-clones'

export function LocalClones() {
  const clonesQuery = useLocalClones()

  if (clonesQuery.isPending) {
    return <LocalClonesSkeleton />
  }

  if (clonesQuery.isError) {
    return (
      <Empty className="flex min-h-0 flex-1 justify-center">
        <EmptyHeader>
          <EmptyTitle>Не удалось просканировать папку</EmptyTitle>
          <EmptyDescription>
            {getErrorMessage(clonesQuery.error)}. Попробуйте ещё раз.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button
            type="button"
            className="cursor-pointer"
            onClick={() => void clonesQuery.refetch()}
          >
            Повторить
          </Button>
        </EmptyContent>
      </Empty>
    )
  }

  if (!clonesQuery.data?.isConfigured) {
    return (
      <Empty className="flex min-h-0 flex-1 justify-center">
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <FolderGit2Icon aria-hidden="true" />
          </EmptyMedia>
          <EmptyTitle>Папка не указана</EmptyTitle>
          <EmptyDescription>
            Укажите папку с git-репозиториями в профиле.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button asChild className="cursor-pointer">
            <Link to="/profile">Открыть профиль</Link>
          </Button>
        </EmptyContent>
      </Empty>
    )
  }

  if (clonesQuery.data.clones.length === 0) {
    return (
      <Empty className="flex min-h-0 flex-1 justify-center">
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <FolderGit2Icon aria-hidden="true" />
          </EmptyMedia>
          <EmptyTitle>Репозиториев нет</EmptyTitle>
          <EmptyDescription>
            В указанной папке не найдено локальных git-репозиториев.
          </EmptyDescription>
        </EmptyHeader>
      </Empty>
    )
  }

  return <LocalCloneList clones={clonesQuery.data.clones} />
}

function LocalCloneList({ clones }: { clones: LocalClone[] }) {
  return (
    <ul
      className="min-h-0 flex-1 overflow-y-auto"
      aria-label="Локальные репозитории"
    >
      {clones.map((clone) => (
        <li
          key={clone.path}
          className="flex min-h-10 items-center gap-2.5 border-b border-white/10 px-1 py-2 last:border-b-0 hover:bg-white/5"
        >
          <span className="min-w-0 flex-1">
            <span
              className="block truncate text-sm font-medium text-white"
              translate="no"
            >
              {clone.name}
            </span>
            <span
              className="mt-0.5 block truncate text-xs text-white/50"
              translate="no"
            >
              {clone.path}
            </span>
          </span>
          <CloudBadge clone={clone} />
        </li>
      ))}
    </ul>
  )
}

function CloudBadge({ clone }: { clone: LocalClone }) {
  if (
    !clone.htmlUrl ||
    (clone.provider !== 'github' && clone.provider !== 'gitlab')
  ) {
    return null
  }

  const meta = gitProviders[clone.provider]
  const Icon = clone.provider === 'gitlab' ? GitLabIcon : GitHubIcon

  const onOpen = async () => {
    try {
      await openPage(clone.fullName || clone.name, clone.htmlUrl)
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <IconButton
      label={`Открыть на ${meta.label}`}
      variant="ghost"
      className="size-8 shrink-0 text-white/70 hover:text-white"
      onClick={() => void onOpen()}
    >
      <Icon className="size-4" />
    </IconButton>
  )
}

function LocalClonesSkeleton() {
  return (
    <div className="space-y-2" aria-hidden="true">
      {Array.from({ length: 8 }, (_, index) => (
        <Skeleton key={index} className="h-10 rounded-md" />
      ))}
    </div>
  )
}
