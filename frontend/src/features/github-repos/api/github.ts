import * as GitHubService from '@bindings/github.com/selyusize/my-home/internal/git/githubservice'
import type { ListOptions, PageInfo } from '@bindings/github.com/selyusize/my-home/pkg/git/models'
import type { GitBranch, GitCommit, GitRepository } from '@/entities/git'

async function unwrapPage<T>(
  result: Promise<[T[] | null, PageInfo]>,
): Promise<{ items: T[]; pageInfo: PageInfo }> {
  const [items, pageInfo] = await result
  return { items: items ?? [], pageInfo }
}

export function listRepos(opts: ListOptions) {
  return unwrapPage<GitRepository>(GitHubService.ListRepos(opts))
}

export function listBranches(owner: string, name: string, opts: ListOptions) {
  return unwrapPage<GitBranch>(GitHubService.ListBranches(owner, name, opts))
}

export function listCommits(
  owner: string,
  name: string,
  ref: string,
  opts: ListOptions,
) {
  return unwrapPage<GitCommit>(GitHubService.ListCommits(owner, name, ref, opts))
}
