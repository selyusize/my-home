import { useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import type { LocalRepoSettings } from '../api/local-repos'
import { useRepoSettings } from '../model/use-repo-settings'

export function RepoSettings() {
  const { settings, isPending, save, check, checking, pickDirectory } = useRepoSettings()
  const [draft, setDraft] = useState<LocalRepoSettings | null>(null)
  const form = draft ?? settings

  const commit = (next: LocalRepoSettings) => {
    setDraft(next)
    save(next)
  }

  const canCheck = Boolean(
    form.sharedPath.trim() ||
      (form.githubSeparate && form.githubPath.trim()) ||
      (form.gitlabSeparate && form.gitlabPath.trim()),
  )

  return (
    <section className="w-full max-w-md space-y-4" aria-labelledby="repo-settings-title">
      <div className="space-y-1">
        <h1 id="repo-settings-title" className="text-lg font-medium text-pretty text-white">
          Репозитории
        </h1>
        <p className="text-sm text-white/70">
          Укажите папки с локальными git clone. Проверка не запускается сама — только по кнопке.
        </p>
      </div>

      <FolderField
        id="repo-shared-path"
        label="Общая папка"
        value={form.sharedPath}
        disabled={isPending || checking}
        onChange={(sharedPath) => setDraft({ ...form, sharedPath })}
        onCommit={(sharedPath) => commit({ ...form, sharedPath })}
        onBrowse={async () => {
          const sharedPath = await pickDirectory('Общая папка с репозиториями')
          if (!sharedPath) {
            return
          }
          commit({ ...form, sharedPath })
        }}
      />

      <SeparateFolder
        id="repo-gitlab-path"
        checkboxId="repo-gitlab-separate"
        checkboxLabel="Отдельная папка для GitLab"
        fieldLabel="Папка GitLab"
        separate={form.gitlabSeparate}
        path={form.gitlabPath}
        disabled={isPending || checking}
        onSeparateChange={(gitlabSeparate) => commit({ ...form, gitlabSeparate })}
        onPathChange={(gitlabPath) => setDraft({ ...form, gitlabPath })}
        onPathCommit={(gitlabPath) => commit({ ...form, gitlabPath })}
        onBrowse={async () => {
          const gitlabPath = await pickDirectory('Папка с репозиториями GitLab')
          if (!gitlabPath) {
            return
          }
          commit({ ...form, gitlabSeparate: true, gitlabPath })
        }}
      />

      <SeparateFolder
        id="repo-github-path"
        checkboxId="repo-github-separate"
        checkboxLabel="Отдельная папка для GitHub"
        fieldLabel="Папка GitHub"
        separate={form.githubSeparate}
        path={form.githubPath}
        disabled={isPending || checking}
        onSeparateChange={(githubSeparate) => commit({ ...form, githubSeparate })}
        onPathChange={(githubPath) => setDraft({ ...form, githubPath })}
        onPathCommit={(githubPath) => commit({ ...form, githubPath })}
        onBrowse={async () => {
          const githubPath = await pickDirectory('Папка с репозиториями GitHub')
          if (!githubPath) {
            return
          }
          commit({ ...form, githubSeparate: true, githubPath })
        }}
      />

      <Button
        type="button"
        className="cursor-pointer"
        disabled={!canCheck || checking || isPending}
        aria-busy={checking}
        onClick={() => check(form)}
      >
        Проверить
      </Button>
    </section>
  )
}

function SeparateFolder({
  id,
  checkboxId,
  checkboxLabel,
  fieldLabel,
  separate,
  path,
  disabled,
  onSeparateChange,
  onPathChange,
  onPathCommit,
  onBrowse,
}: {
  id: string
  checkboxId: string
  checkboxLabel: string
  fieldLabel: string
  separate: boolean
  path: string
  disabled: boolean
  onSeparateChange: (value: boolean) => void
  onPathChange: (value: string) => void
  onPathCommit: (value: string) => void
  onBrowse: () => void
}) {
  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <Checkbox
          id={checkboxId}
          checked={separate}
          disabled={disabled}
          onCheckedChange={(value) => {
            if (value === 'indeterminate') {
              return
            }
            onSeparateChange(value)
          }}
        />
        <Label htmlFor={checkboxId} className="text-white/80">
          {checkboxLabel}
        </Label>
      </div>
      {separate ? (
        <FolderField
          id={id}
          label={fieldLabel}
          value={path}
          disabled={disabled}
          onChange={onPathChange}
          onCommit={onPathCommit}
          onBrowse={onBrowse}
        />
      ) : null}
    </div>
  )
}

function FolderField({
  id,
  label,
  value,
  disabled,
  onChange,
  onCommit,
  onBrowse,
}: {
  id: string
  label: string
  value: string
  disabled: boolean
  onChange: (value: string) => void
  onCommit: (value: string) => void
  onBrowse: () => void
}) {
  return (
    <div className="grid gap-1.5">
      <Label htmlFor={id}>{label}</Label>
      <div className="flex gap-2">
        <Input
          id={id}
          value={value}
          disabled={disabled}
          onChange={(event) => onChange(event.target.value)}
          onBlur={() => onCommit(value)}
          autoComplete="off"
          spellCheck={false}
          placeholder="/Users/me/dev"
        />
        <Button
          type="button"
          variant="outline"
          className="cursor-pointer shrink-0"
          disabled={disabled}
          onClick={onBrowse}
        >
          Выбрать
        </Button>
      </div>
    </div>
  )
}
