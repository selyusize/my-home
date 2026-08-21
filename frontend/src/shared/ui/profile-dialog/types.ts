export type AccountProfileStat = {
  label: string
  value: string | number
}

export type AccountProfile = {
  name: string
  handle?: string
  email?: string
  bio?: string
  company?: string
  location?: string
  website?: string
  avatarUrl?: string
  pageUrl?: string
  stats?: AccountProfileStat[]
}
