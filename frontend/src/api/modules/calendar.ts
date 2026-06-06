import { get, post, put, del } from '../request'
import type { CalendarEvent } from '@/types'

/** 我的日历事件列表 */
export function getCalendarEventsApi() {
  return get<CalendarEvent[]>('/calendar/events')
}

/** 添加到日历 */
export function addCalendarEventApi(data: { eventType: string; refId: number }) {
  return post<CalendarEvent>('/calendar/events', data)
}

/** 修改日历事件 */
export function updateCalendarEventApi(id: number, data: Partial<CalendarEvent>) {
  return put<CalendarEvent>(`/calendar/events/${id}`, data)
}

/** 删除日历事件 */
export function deleteCalendarEventApi(id: number) {
  return del<void>(`/calendar/events/${id}`)
}
