import { cn } from '@/shared/lib/utils'
import { useOrderedItems } from '@/shared/lib/use-ordered-items'
import { Avatar, AvatarFallback } from '@/shared/ui/avatar'
import { SortableList } from '@/shared/ui/sortable-list'
import { NavLink } from 'react-router-dom'
import { getNavItemId, getNavItemLabel, navItems } from '../model/nav'

const itemClass = ({ isActive }: { isActive: boolean }) =>
  cn(
    'flex size-10 items-center justify-center rounded-lg outline-none transition-colors hover:bg-white/5 hover:text-white focus-visible:ring-2 focus-visible:ring-ring/50',
    isActive ? 'bg-white/10 text-white' : 'text-white/60',
  )

export function LeftPanelWidget() {
  const [items, reorder] = useOrderedItems('left-nav', navItems, getNavItemId)

  return (
    <aside className="flex h-full w-14 shrink-0 flex-col items-center border-r border-white/10 py-3">
      <nav className="flex w-full flex-col items-center px-2">
        <SortableList
          items={items}
          getId={getNavItemId}
          getLabel={getNavItemLabel}
          onReorder={reorder}
          label="Основное меню"
          className="w-full items-center gap-1"
        >
          {(item) => (
            <NavLink
              to={item.to}
              end={item.to === '/'}
              aria-label={item.title}
              className={itemClass}
              draggable={false}
            >
              <item.icon className="size-5" strokeWidth={1.75} aria-hidden="true" />
            </NavLink>
          )}
        </SortableList>
      </nav>

      <div className="mt-auto flex w-full items-center justify-center px-2">
        <NavLink to="/profile" aria-label="Профиль" className={itemClass}>
          <Avatar size="sm">
            <AvatarFallback className="text-xs font-medium text-white">
              PR
            </AvatarFallback>
          </Avatar>
        </NavLink>
      </div>
    </aside>
  )
}
