import { get, post, patch } from '@/api/request'
import type { PageResult } from '@/types/api'
import type {
  AdminUserCreateForm,
  AdminUserQuery,
  AdminUserUpdateForm,
  AdminUserVO,
  UserStatus,
} from '@/types/admin'

export function getAdminUsersApi(params?: AdminUserQuery) {
  return get<PageResult<AdminUserVO>>('/admin/users', { params })
}

export function createAdminUserApi(data: AdminUserCreateForm) {
  return post<AdminUserVO>('/admin/users', data)
}

export function updateAdminUserApi(id: number, data: AdminUserUpdateForm) {
  return patch<AdminUserVO>(`/admin/users/${id}`, data)
}

export function updateAdminUserStatusApi(id: number, status: UserStatus) {
  return patch<null>(`/admin/users/${id}/status`, { status })
}

export function resetAdminUserPasswordApi(id: number, newPassword: string) {
  return post<null>(`/admin/users/${id}/reset-password`, { newPassword })
}
