import { cn } from '@/shared/lib/utils'
import { useOrderedItems } from '@/shared/lib/use-ordered-items'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/shared/ui/tooltip'
import { SortableList } from '@/shared/ui/sortable-list'
import { NavLink } from 'react-router-dom'
import {
  getGlobeMenuItemId,
  getGlobeMenuItemLabel,
  globeMenuItems,
  globeMenuOrderKey,
} from '../model/items'

export function GlobeMenuWidget() {
  const [items, reorder] = useOrderedItems(
    globeMenuOrderKey,
    globeMenuItems,
    getGlobeMenuItemId,
  )

  return (
    <aside className="flex h-full w-14 shrink-0 flex-col items-center border-l border-white/10 py-3">
      <nav className="flex w-full flex-col items-center px-2">
        <SortableList
          items={items}
          getId={getGlobeMenuItemId}
          getLabel={getGlobeMenuItemLabel}
          onReorder={reorder}
          label="Сервисы"
          className="w-full items-center gap-1"
        >
          {(item, { isDragging }) => (
            <Tooltip open={isDragging ? false : undefined}>
              <TooltipTrigger asChild>
                <span className="inline-flex">
                  <NavLink
                    to={item.to}
                    aria-label={item.title}
                    draggable={false}
                    className={({ isActive }) =>
                      cn(
                        'flex size-10 items-center justify-center rounded-lg outline-none transition-colors hover:bg-white/5 hover:text-white focus-visible:ring-2 focus-visible:ring-ring/50',
                        isActive ? 'bg-white/10 text-white' : 'text-white/60',
                      )
                    }
                  >
                    <item.icon className="size-5" aria-hidden="true" />
                  </NavLink>
                </span>
              </TooltipTrigger>
              <TooltipContent side="left">{item.title}</TooltipContent>
            </Tooltip>
          )}
        </SortableList>
      </nav>
    </aside>
  )
}
