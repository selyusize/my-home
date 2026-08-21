import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo } from 'react'
import { toast } from 'sonner'
import { getErrorMessage } from '@/shared/lib/error-message'
import { settingsApi } from '@/shared/lib/settings'

export function applyItemOrder<T>(
  items: T[],
  order: string[],
  getId: (item: T) => string,
) {
  const byId = new Map(items.map((item) => [getId(item), item]))
  const used = new Set<string>()
  const next: T[] = []

  for (const id of order) {
    const item = byId.get(id)
    if (!item) {
      continue
    }
    next.push(item)
    used.add(id)
  }

  for (const item of items) {
    const id = getId(item)
    if (!used.has(id)) {
      next.push(item)
    }
  }

  return next
}

export const menuOrderKeys = {
  all: ['menu-order'] as const,
  of: (key: string) => [...menuOrderKeys.all, key] as const,
}

export function useOrderedItems<T>(
  key: string,
  items: T[],
  getId: (item: T) => string,
) {
  const queryClient = useQueryClient()
  const queryKey = menuOrderKeys.of(key)
  const {
    data: persisted = [],
    isPending,
    isError,
  } = useQuery({
    queryKey,
    queryFn: () => settingsApi.getOrder(key),
    retry: false,
    staleTime: Infinity,
  })

  const ordered = useMemo(
    () => applyItemOrder(items, persisted ?? [], getId),
    [getId, items, persisted],
  )

  const save = useMutation({
    mutationFn: (order: string[]) => settingsApi.setOrder(key, order),
    onMutate: async (order) => {
      await queryClient.cancelQueries({ queryKey })
      const previous = queryClient.getQueryData<string[]>(queryKey)
      queryClient.setQueryData(queryKey, order)
      return { previous }
    },
    onError: (error, _order, context) => {
      if (context?.previous) {
        queryClient.setQueryData(queryKey, context.previous)
      }
      toast.error(getErrorMessage(error))
    },
    onSettled: () => queryClient.invalidateQueries({ queryKey }),
  })

  const reorder = (next: T[]) => {
    save.mutate(next.map(getId))
  }

  return [ordered, reorder, !isPending || isError] as const
}
