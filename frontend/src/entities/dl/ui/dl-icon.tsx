import type { ImgHTMLAttributes } from 'react'
import { cn } from '@/shared/lib/utils'
import dlMark from './dl.png'

export function DlIcon({ className, alt = '', ...props }: ImgHTMLAttributes<HTMLImageElement>) {
  return (
    <img
      src={dlMark}
      alt={alt}
      aria-hidden={alt ? undefined : true}
      className={cn('size-4 rounded-sm object-cover', className)}
      {...props}
    />
  )
}
