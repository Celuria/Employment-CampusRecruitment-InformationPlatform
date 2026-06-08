import { get, post, put, del, patch } from '@/api/request'
import type { PageResult } from '@/types/api'
import type { AdminJobFairForm, AdminJobFairQuery, AdminJobFairVO, PublishStatus } from '@/types/admin'

export function getAdminJobFairsApi(params?: AdminJobFairQuery) {
  return get<PageResult<AdminJobFairVO>>('/admin/job-fairs', { params })
}

export function createAdminJobFairApi(data: AdminJobFairForm) {
  return post<AdminJobFairVO>('/admin/job-fairs', data)
}

export function updateAdminJobFairApi(id: number, data: Partial<AdminJobFairForm>) {
  return put<AdminJobFairVO>(`/admin/job-fairs/${id}`, data)
}

export function deleteAdminJobFairApi(id: number) {
  return del<null>(`/admin/job-fairs/${id}`)
}

export function batchAdminJobFairStatusApi(ids: number[], publishStatus: PublishStatus) {
  return patch<null>('/admin/job-fairs/batch-status', { ids, publishStatus })
}
