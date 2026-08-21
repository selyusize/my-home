import { GitLabIcon } from '@/entities/git'
import { useGitSession } from '../model/use-git-session'
import { ServiceControl } from './service-control'

export function GitLabSession() {
  const session = useGitSession('gitlab')
  return <ServiceControl label="GitLab" icon={GitLabIcon} {...session} />
}
