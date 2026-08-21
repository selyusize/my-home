export type DlServiceView = {
  name: string
  running: boolean
  present: boolean
}

export type DlSnapshot = {
  installed: boolean
  version: string
  latest: string
  path: string
  dockerOk: boolean
  dockerVersion: string
  dockerOs: string
  updateAvailable: boolean
  serviceUp: boolean
  services: DlServiceView[]
}

export const EMPTY_DL_SERVICES: DlServiceView[] = [
  { name: 'traefik', running: false, present: false },
  { name: 'portainer', running: false, present: false },
  { name: 'mail', running: false, present: false },
]
