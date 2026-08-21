import * as BitrixService from '@bindings/github.com/selyusize/my-home/internal/bitrix/bitrixservice'
import type { Profile, TimeMan } from '@bindings/github.com/selyusize/my-home/pkg/bitrix/models'

export const bitrixApi = {
  isAuthenticated: BitrixService.IsAuthenticated,
  profile: BitrixService.Profile,
  timeMan: BitrixService.TimeMan,
  setWebhook: BitrixService.SetWebhook,
  logout: BitrixService.Logout,
  timeManOpen: BitrixService.TimeManOpen,
  timeManPause: BitrixService.TimeManPause,
  timeManClose: BitrixService.TimeManClose,
}

export type BitrixProfile = Profile
export type BitrixTimeMan = TimeMan
