import { get, post } from '@/api/request'
import type { PageResult } from '@/types/api'
import type { SyncLogVO, SyncTaskVO } from '@/types/admin'

export function triggerSyncApi(sourceType = 'all', force = false) {
  return post<SyncTaskVO>('/admin/sync', { sourceType, force })
}

export function getSyncLogsApi(page = 1, pageSize = 10) {
  return get<PageResult<SyncLogVO>>('/admin/sync/logs', { params: { page, pageSize } })
}
