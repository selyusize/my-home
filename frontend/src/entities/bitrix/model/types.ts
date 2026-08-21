export type BitrixAccount = {
  id?: number | string
  name?: string
  email?: string
  position?: string
  avatarUrl?: string
  pageUrl?: string
  portalUrl?: string
}

export type BitrixTimeManStatus = 'opened' | 'paused' | 'closed' | 'expired' | ''

export type BitrixTimeMan = {
  status?: BitrixTimeManStatus
  timeStart?: string
  duration?: number
  timeLeaks?: number
}
