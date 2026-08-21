import { EMPTY_DL_SERVICES, type DlServiceView, type DlSnapshot } from './types'

type DlStatusInput = {
  installed?: boolean
  version?: string
  latest?: string
  path?: string
  dockerOk?: boolean
  dockerVersion?: string
  dockerOs?: string
  updateAvailable?: boolean
  serviceUp?: boolean
  services?: DlServiceView[] | null
}

export function mapDlStatusToView(data?: DlStatusInput | null): DlSnapshot {
  return {
    installed: Boolean(data?.installed),
    version: data?.version ?? '',
    latest: data?.latest ?? '',
    path: data?.path ?? '',
    dockerOk: Boolean(data?.dockerOk),
    dockerVersion: data?.dockerVersion ?? '',
    dockerOs: data?.dockerOs ?? '',
    updateAvailable: Boolean(data?.updateAvailable),
    serviceUp: Boolean(data?.serviceUp),
    services: data?.services?.length ? data.services : EMPTY_DL_SERVICES,
  }
}
