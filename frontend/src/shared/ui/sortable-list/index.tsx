import { useRef, useState, type CSSProperties, type ReactNode } from 'react'
import {
  closestCenter,
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
  type DragStartEvent,
  type DraggableAttributes,
} from '@dnd-kit/core'
import { restrictToHorizontalAxis, restrictToVerticalAxis } from '@dnd-kit/modifiers'
import {
  arrayMove,
  horizontalListSortingStrategy,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { cn } from '@/shared/lib/utils'
import { usePrefersReducedMotion } from '@/shared/lib/use-prefers-reduced-motion'

export type SortableItemState = {
  isDragging: boolean
}

export type SortableListProps<T> = {
  items: T[]
  getId: (item: T) => string
  getLabel?: (item: T) => string
  onReorder: (items: T[]) => void
  orientation?: 'vertical' | 'horizontal'
  className?: string
  itemClassName?: string
  label: string
  children: (item: T, state: SortableItemState) => ReactNode
}

const freezeLayoutChanges = () => false

export function SortableList<T>({
  items,
  getId,
  getLabel,
  onReorder,
  orientation = 'vertical',
  className,
  itemClassName,
  label,
  children,
}: SortableListProps<T>) {
  const skipClickRef = useRef(false)
  const [activeId, setActiveId] = useState<string | null>(null)
  const [liveMessage, setLiveMessage] = useState('')
  const reduceMotion = usePrefersReducedMotion()
  const ids = items.map(getId)
  const vertical = orientation === 'vertical'
  const axisModifier = vertical ? restrictToVerticalAxis : restrictToHorizontalAxis
  const activeItem = items.find((item) => getId(item) === activeId)
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  )

  const handleDragStart = (event: DragStartEvent) => {
    skipClickRef.current = false
    const id = String(event.active.id)
    setActiveId(id)
    const item = items.find((entry) => getId(entry) === id)
    const name = item && getLabel ? getLabel(item) : id
    setLiveMessage(`Перемещение: ${name}`)
  }

  const handleDragEnd = (event: DragEndEvent) => {
    skipClickRef.current = true
    const { active, over } = event

    if (over && active.id !== over.id) {
      const oldIndex = items.findIndex((item) => getId(item) === String(active.id))
      const newIndex = items.findIndex((item) => getId(item) === String(over.id))
      if (oldIndex >= 0 && newIndex >= 0) {
        onReorder(arrayMove(items, oldIndex, newIndex))
        setLiveMessage('Порядок обновлён')
      }
    }

    setActiveId(null)
    window.setTimeout(() => {
      skipClickRef.current = false
    }, 0)
  }

  const handleDragCancel = () => {
    skipClickRef.current = true
    setActiveId(null)
    setLiveMessage('Перемещение отменено')
    window.setTimeout(() => {
      skipClickRef.current = false
    }, 0)
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      modifiers={[axisModifier]}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
      onDragCancel={handleDragCancel}
    >
      <SortableContext
        items={ids}
        strategy={vertical ? verticalListSortingStrategy : horizontalListSortingStrategy}
      >
        <ul
          aria-label={label}
          data-dragging={activeId ? '' : undefined}
          className={cn(
            'no-drag flex select-none',
            vertical ? 'flex-col' : 'flex-row',
            className,
          )}
        >
          {items.map((item) => (
            <SortableItem
              key={getId(item)}
              id={getId(item)}
              item={item}
              reduceMotion={reduceMotion}
              skipClickRef={skipClickRef}
              className={itemClassName}
            >
              {children}
            </SortableItem>
          ))}
        </ul>
      </SortableContext>
      <DragOverlay dropAnimation={null} modifiers={[axisModifier]}>
        {activeItem ? (
          <div className="pointer-events-none">{children(activeItem, { isDragging: true })}</div>
        ) : null}
      </DragOverlay>
      <div className="sr-only" aria-live="polite">
        {liveMessage}
      </div>
    </DndContext>
  )
}

function SortableItem<T>({
  id,
  item,
  reduceMotion,
  skipClickRef,
  className,
  children,
}: {
  id: string
  item: T
  reduceMotion: boolean
  skipClickRef: { current: boolean }
  className?: string
  children: (item: T, state: SortableItemState) => ReactNode
}) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({
      id,
      animateLayoutChanges: freezeLayoutChanges,
      transition: reduceMotion ? null : { duration: 180, easing: 'cubic-bezier(0.25, 1, 0.5, 1)' },
    })
  const style: CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition: isDragging ? 'none' : transition,
  }

  return (
    <li
      ref={setNodeRef}
      style={style}
      className={cn(
        'relative list-none touch-none',
        isDragging && 'opacity-0',
        className,
      )}
      {...withLinkSafeAttributes(attributes)}
      {...listeners}
      onClickCapture={(event) => {
        if (!skipClickRef.current) {
          return
        }
        event.preventDefault()
        event.stopPropagation()
      }}
    >
      {children(item, { isDragging })}
    </li>
  )
}

function withLinkSafeAttributes(attributes: DraggableAttributes) {
  return {
    ...attributes,
    role: undefined,
    tabIndex: 0,
    'aria-roledescription': 'sortable',
  }
}
