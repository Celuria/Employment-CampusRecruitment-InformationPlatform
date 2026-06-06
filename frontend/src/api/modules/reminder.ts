import { get } from '../request'
import type { ReminderLog } from '@/types'

/** 提醒记录列表 */
export function getReminderLogsApi() {
  return get<ReminderLog[]>('/reminders/logs')
}
