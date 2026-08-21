import type { AccountProfile } from '@/shared/ui/profile-dialog'
import type { GitAccount } from './types'

export function mapGitAccountToProfile(account: GitAccount): AccountProfile {
  return {
    name: account.name || account.login,
    handle: account.login,
    email: account.email || undefined,
    bio: account.bio || undefined,
    company: account.company || undefined,
    location: account.location || undefined,
    website: account.blog || undefined,
    avatarUrl: account.avatarUrl || undefined,
    pageUrl: account.htmlUrl || undefined,
    stats: [
      { label: 'Репозитории', value: account.publicRepos },
      { label: 'Подписчики', value: account.followers },
      { label: 'Подписки', value: account.following },
    ],
  }
}
