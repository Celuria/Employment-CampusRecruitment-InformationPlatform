import { get } from '../request'
import type { ReminderLog, PageResult } from '@/types'

/** 提醒记录列表（分页） */
export function getReminderLogsApi(params?: { page?: number; pageSize?: number }) {
  return get<PageResult<ReminderLog>>('/reminders/logs', { params })
}
