import * as DLService from '@bindings/github.com/selyusize/my-home/internal/dl/dlservice'
import type { Status } from '@bindings/github.com/selyusize/my-home/pkg/dl/models'

export const dlApi = {
  status: DLService.Status,
  install: DLService.Install,
  uninstall: DLService.Uninstall,
  serviceUp: DLService.ServiceUp,
  serviceDown: DLService.ServiceDown,
}

export type DlStatus = Status
