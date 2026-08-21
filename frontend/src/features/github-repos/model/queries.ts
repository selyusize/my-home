import { useInfiniteQuery, useQuery } from '@tanstack/react-query'
import { listBranches, listCommits, listRepos } from '../api/github'

const REPO_PAGE_SIZE = 50
const BRANCH_PAGE_SIZE = 30
const COMMIT_PAGE_SIZE = 10

export const githubRepoKeys = {
  all: ['github-repos'] as const,
  list: () => [...githubRepoKeys.all, 'list'] as const,
  branches: (owner: string, name: string) =>
    [...githubRepoKeys.all, 'branches', owner, name] as const,
  commits: (owner: string, name: string, ref: string) =>
    [...githubRepoKeys.all, 'commits', owner, name, ref] as const,
}

export function useGitHubRepos() {
  return useInfiniteQuery({
    queryKey: githubRepoKeys.list(),
    queryFn: ({ pageParam }) =>
      listRepos({ page: pageParam, perPage: REPO_PAGE_SIZE }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) => lastPage.pageInfo.nextPage || undefined,
    retry: false,
    staleTime: 30_000,
  })
}

export function useGitHubBranches(owner: string, name: string, enabled: boolean) {
  return useQuery({
    queryKey: githubRepoKeys.branches(owner, name),
    queryFn: () =>
      listBranches(owner, name, { page: 1, perPage: BRANCH_PAGE_SIZE }),
    enabled,
    retry: false,
    staleTime: 30_000,
  })
}

export function useGitHubCommits(
  owner: string,
  name: string,
  ref: string,
  enabled: boolean,
) {
  return useQuery({
    queryKey: githubRepoKeys.commits(owner, name, ref),
    queryFn: () =>
      listCommits(owner, name, ref, { page: 1, perPage: COMMIT_PAGE_SIZE }),
    enabled,
    retry: false,
    staleTime: 15_000,
  })
}
