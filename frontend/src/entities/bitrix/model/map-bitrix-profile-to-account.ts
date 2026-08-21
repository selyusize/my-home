import type { AccountProfile } from '@/shared/ui/profile-dialog'
import type { BitrixAccount } from './types'

export function mapBitrixProfileToAccount(account: BitrixAccount): AccountProfile {
  return {
    name: account.name || account.email || 'Bitrix24',
    handle: account.id ? String(account.id) : undefined,
    email: account.email || undefined,
    company: account.position || undefined,
    avatarUrl: account.avatarUrl || undefined,
    pageUrl: account.pageUrl || account.portalUrl || undefined,
  }
}
