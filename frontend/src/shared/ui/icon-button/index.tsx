import type { ComponentProps, ReactNode } from 'react'
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/button'
import { Spinner } from '@/shared/ui/spinner'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/shared/ui/tooltip'

type IconButtonProps = Omit<ComponentProps<typeof Button>, 'size' | 'children'> & {
  label: string
  pending?: boolean
  badge?: ReactNode
  children: ReactNode
}

export function IconButton({
  label,
  pending = false,
  badge,
  disabled,
  variant = 'ghost',
  className,
  children,
  ...props
}: IconButtonProps) {
  const isDisabled = Boolean(disabled || pending)

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className="inline-flex">
          <Button
            type="button"
            variant={variant}
            size="icon"
            aria-label={label}
            aria-busy={pending}
            disabled={isDisabled}
            className={cn(
              'relative cursor-pointer text-white hover:bg-white/10 hover:text-white [&_svg:not([class*="size-"])]:size-4',
              className,
            )}
            {...props}
          >
            {pending ? <Spinner /> : children}
            {badge}
          </Button>
        </span>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  )
}
