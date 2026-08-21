import { GitHubIcon } from '@/entities/git'
import { useGitSession } from '../model/use-git-session'
import { ServiceControl } from './service-control'

export function GitHubSession() {
  const session = useGitSession('github')
  return <ServiceControl label="GitHub" icon={GitHubIcon} {...session} />
}
