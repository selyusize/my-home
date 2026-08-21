export const repoSettingsKeys = {
  all: ['repo-settings'] as const,
  settings: () => [...repoSettingsKeys.all, 'settings'] as const,
  clones: (provider?: 'github' | 'gitlab') =>
    [...repoSettingsKeys.all, 'clones', provider] as const,
  explorer: () => [...repoSettingsKeys.all, 'explorer'] as const,
}
