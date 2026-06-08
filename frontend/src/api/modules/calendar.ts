import { get, post, put, del } from '../request'
import type { CalendarEvent, CalendarEventQuery } from '@/types'

/** 我的日历事件列表 */
export function getCalendarEventsApi(params?: CalendarEventQuery) {
  return get<CalendarEvent[]>('/calendar/events', { params })
}

/** 添加到日历 */
export function addCalendarEventApi(data: { eventType: string; refId: number }) {
  return post<CalendarEvent>('/calendar/events', data)
}

/** 修改日历事件 */
export function updateCalendarEventApi(id: number, data: { customNote?: string; remindBefore?: string[] }) {
  return put<CalendarEvent>(`/calendar/events/${id}`, data)
}

/** 删除日历事件 */
export function deleteCalendarEventApi(id: number) {
  return del<void>(`/calendar/events/${id}`)
}
