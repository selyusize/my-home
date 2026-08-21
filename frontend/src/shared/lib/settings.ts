import {
  GetOrder,
  SetOrder,
} from '@bindings/github.com/selyusize/my-home/internal/settings/settingsservice'

export const settingsApi = {
  getOrder: async (key: string) => (await GetOrder(key)) ?? [],
  setOrder: (key: string, order: string[]) => SetOrder(key, order),
}
