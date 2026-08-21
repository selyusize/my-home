import { useDlRuntime } from '../model/use-dl-runtime'
import { DlControl } from './dl-control'

export function DlRuntime() {
  const dl = useDlRuntime()
  return <DlControl {...dl} />
}
