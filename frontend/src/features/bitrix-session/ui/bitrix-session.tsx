import { useBitrixSession } from '../model/use-bitrix-session'
import { BitrixControl } from './bitrix-control'

export function BitrixSession() {
  const session = useBitrixSession()
  return <BitrixControl {...session} />
}
