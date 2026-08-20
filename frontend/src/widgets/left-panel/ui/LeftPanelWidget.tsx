import { cn } from '@/shared/lib/utils'
import { Avatar, AvatarFallback } from '@/shared/ui'
import { NavLink } from 'react-router-dom'
import { navItems } from '../model/nav'

const LeftPanelWidget = () => {
  return (
    <aside className="flex h-full w-20 shrink-0 flex-col items-center border-r border-white/10 py-4">
      <nav className="flex w-full flex-col items-center gap-2 px-2">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === '/'}
            title={item.title}
            aria-label={item.title}
            className={({ isActive }) =>
              cn(
                'flex size-12 items-center justify-center rounded-full transition-[width,height] duration-300 ease-in-out',
                isActive
                  ? 'size-14 bg-white/10 text-white'
                  : 'text-white/60 hover:bg-white/5 hover:text-white',
              )
            }
          >
            {({ isActive }) => (
              <item.icon
                className={cn(
                  'transition-[width,height] duration-300 ease-in-out',
                  isActive ? 'size-7' : 'size-5',
                )}
                strokeWidth={1.75}
              />
            )}
          </NavLink>
        ))}
      </nav>

      <div className="mt-auto flex w-full items-center justify-center px-2 pb-1">
        <button
          type="button"
          title="Профиль"
          aria-label="Профиль"
          className="rounded-full transition-colors hover:bg-white/5"
        >
          <Avatar size="lg" className="size-12">
            <AvatarFallback className="bg-white/10 text-sm font-medium text-white">
              PR
            </AvatarFallback>
          </Avatar>
        </button>
      </div>
    </aside>
  )
}

export { LeftPanelWidget }
