import { get, post, put, del, patch } from '@/api/request'
import type { PageResult } from '@/types/api'
import type { AdminCareerTalkForm, AdminCareerTalkQuery, AdminCareerTalkVO, PublishStatus } from '@/types/admin'

export function getAdminCareerTalksApi(params?: AdminCareerTalkQuery) {
  return get<PageResult<AdminCareerTalkVO>>('/admin/career-talks', { params })
}

export function createAdminCareerTalkApi(data: AdminCareerTalkForm) {
  return post<AdminCareerTalkVO>('/admin/career-talks', data)
}

export function updateAdminCareerTalkApi(id: number, data: Partial<AdminCareerTalkForm>) {
  return put<AdminCareerTalkVO>(`/admin/career-talks/${id}`, data)
}

export function deleteAdminCareerTalkApi(id: number) {
  return del<null>(`/admin/career-talks/${id}`)
}

export function batchAdminCareerTalkStatusApi(ids: number[], publishStatus: PublishStatus) {
  return patch<null>('/admin/career-talks/batch-status', { ids, publishStatus })
}
