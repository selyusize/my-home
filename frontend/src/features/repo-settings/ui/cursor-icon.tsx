import type { SVGProps } from 'react'

export function CursorIcon({ className, ...props }: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="currentColor"
      aria-hidden
      className={className}
      {...props}
    >
      <path d="M11.925 24 15.237 11.988 24 8.69 11.988 5.378 8.69 0 5.378 12.012 0 15.31l12.012 3.312z" />
    </svg>
  )
}
