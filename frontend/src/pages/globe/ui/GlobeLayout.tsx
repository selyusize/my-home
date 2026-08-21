import { GlobeMenuWidget } from '@widgets/globe-menu'
import { Outlet } from 'react-router-dom'

export function GlobeLayout() {
  return (
    <div className="flex min-h-0 flex-1">
      <div className="flex min-h-0 min-w-0 flex-1 flex-col p-6">
        <Outlet />
      </div>
      <GlobeMenuWidget />
    </div>
  )
}
