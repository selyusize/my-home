import * as LocalReposService from '@bindings/github.com/selyusize/my-home/internal/git/localreposservice'
import type { LocalRepoSettings } from '@bindings/github.com/selyusize/my-home/internal/git/models'
import type { Clone } from '@bindings/github.com/selyusize/my-home/pkg/git/local/models'

export type { LocalRepoSettings }
export type LocalClone = Clone

export const emptyRepoSettings: LocalRepoSettings = {
  sharedPath: '',
  githubSeparate: false,
  githubPath: '',
  gitlabSeparate: false,
  gitlabPath: '',
}

export const localReposApi = {
  getSettings: async (): Promise<LocalRepoSettings> => {
    const settings = await LocalReposService.GetSettings()
    return {
      ...emptyRepoSettings,
      ...settings,
    }
  },
  saveSettings: (settings: LocalRepoSettings) => LocalReposService.SaveSettings(settings),
  selectDirectory: (title: string) => LocalReposService.SelectDirectory(title),
  check: (settings: LocalRepoSettings) => LocalReposService.Check(settings),
  listClones: async (provider: 'github' | 'gitlab' | '' = '') =>
    (await LocalReposService.ListClones(provider)) ?? [],
}
