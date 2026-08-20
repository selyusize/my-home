import { LeftPanelWidget } from '@widgets/left-panel'
import type { ReactNode } from 'react'

type AppLayoutProps = {
  children: ReactNode
}

export function AppLayout({ children }: AppLayoutProps) {
  return (
    <div className="flex h-full flex-col bg-[#06070f] text-white">
      <header className="drag flex h-titlebar shrink-0 items-center border-white/10 pr-4 pl-traffic-lights"></header>
      <div className="flex flex-1">
        <LeftPanelWidget />
        <main className="border w-full">{children}</main>
      </div>
    </div>
  )
}
