import { useOrderedItems } from '@/shared/lib/use-ordered-items'
import { Navigate } from 'react-router-dom'
import {
  getGlobeMenuItemId,
  globeMenuItems,
  globeMenuOrderKey,
} from '@widgets/globe-menu'

export function GlobeIndexRedirect() {
  const [items, , isReady] = useOrderedItems(
    globeMenuOrderKey,
    globeMenuItems,
    getGlobeMenuItemId,
  )

  if (!isReady) {
    return null
  }

  return <Navigate to={items[0]?.to ?? '/globe/bitrix'} replace />
}
