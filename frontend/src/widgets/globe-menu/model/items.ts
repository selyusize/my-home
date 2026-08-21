import type { ComponentType, SVGProps } from 'react'
import { BitrixIcon } from '@/entities/bitrix'
import { GitHubIcon, GitLabIcon } from '@/entities/git'
import { FolderGit2Icon } from 'lucide-react'

export type GlobeMenuItem = {
  title: string
  to: string
  icon: ComponentType<SVGProps<SVGSVGElement>>
}

export const globeMenuOrderKey = 'globe-menu'

export const getGlobeMenuItemId = (item: GlobeMenuItem) => item.to

export const getGlobeMenuItemLabel = (item: GlobeMenuItem) => item.title

export const globeMenuItems: GlobeMenuItem[] = [
  {
    title: 'Bitrix24',
    to: '/globe/bitrix',
    icon: BitrixIcon,
  },
  {
    title: 'GitLab',
    to: '/globe/gitlab',
    icon: GitLabIcon,
  },
  {
    title: 'GitHub',
    to: '/globe/github',
    icon: GitHubIcon,
  },
  {
    title: 'Проводник',
    to: '/globe/explorer',
    icon: FolderGit2Icon,
  },
]
