import { LeftPanelWidget } from '@widgets/left-panel'
import { Outlet } from 'react-router-dom'

export function AppLayout() {
  return (
    <div className="flex h-full flex-col bg-[#06070f] text-white">
      <header className="drag flex h-titlebar shrink-0 items-center border-b border-white/10 pr-4 pl-traffic-lights" />
      <div className="flex min-h-0 flex-1">
        <LeftPanelWidget />
        <main className="min-w-0 flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
