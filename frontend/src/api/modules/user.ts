import { get, put } from '../request'
import type { UserInfo, UserPreference, UpdateProfileForm } from '@/types'

/** 获取当前用户资料 */
export function getProfileApi() {
  return get<UserInfo>('/users/me')
}

/** 更新当前用户资料 */
export function updateProfileApi(data: UpdateProfileForm) {
  return put<UserInfo>('/users/me', data)
}

/** 获取用户偏好设置 */
export function getPreferencesApi() {
  return get<UserPreference>('/users/me/preferences')
}

/** 更新用户偏好设置 */
export function updatePreferencesApi(data: UserPreference) {
  return put<UserPreference>('/users/me/preferences', data)
}
