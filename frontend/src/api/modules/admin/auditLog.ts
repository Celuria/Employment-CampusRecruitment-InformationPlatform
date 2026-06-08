import { get } from '@/api/request'
import type { PageResult } from '@/types/api'
import type { AuditLogQuery, AuditLogVO } from '@/types/admin'

export function getAuditLogsApi(params?: AuditLogQuery) {
  return get<PageResult<AuditLogVO>>('/admin/audit-logs', { params })
}
