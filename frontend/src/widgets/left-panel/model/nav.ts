import type { LucideIcon } from 'lucide-react'
import { GlobeIcon, HomeIcon } from 'lucide-react'

export type NavItem = {
  title: string
  to: string
  icon: LucideIcon
}

export const navItems: NavItem[] = [
  {
    title: 'Дом',
    to: '/',
    icon: HomeIcon,
  },
  {
    title: 'Работа с проектами',
    to: '/globe',
    icon: GlobeIcon,
  },
]
