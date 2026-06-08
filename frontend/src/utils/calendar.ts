import type { CalendarEvent } from '@/types'

export function formatDateParam(date: Date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

/** 计算月历可见网格的起止日期（含前后补位周） */
export function getVisibleGridRange(month: Date) {
  const year = month.getFullYear()
  const m = month.getMonth()
  const firstDay = new Date(year, m, 1)
  const lastDay = new Date(year, m + 1, 0)

  const gridStart = new Date(firstDay)
  gridStart.setDate(firstDay.getDate() - firstDay.getDay())

  const gridEnd = new Date(lastDay)
  gridEnd.setDate(lastDay.getDate() + (6 - lastDay.getDay()))

  return {
    startDate: formatDateParam(gridStart),
    endDate: formatDateParam(gridEnd),
    gridStart,
    gridEnd,
  }
}

function dayStart(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate())
}

function dayEnd(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate(), 23, 59, 59, 999)
}

/** 判断事件是否覆盖某一天（支持跨天） */
export function eventCoversDay(event: CalendarEvent, date: Date) {
  const start = new Date(event.startTime)
  const end = event.endTime ? new Date(event.endTime) : start
  return start <= dayEnd(date) && end >= dayStart(date)
}

/** 获取某一天的所有事件，按开始时间升序 */
export function getEventsForDay(date: Date, events: CalendarEvent[]) {
  return events
    .filter((e) => eventCoversDay(e, date))
    .sort((a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime())
}

/** 月历格内事件时间简写 */
export function formatEventChipTime(startTime: string) {
  const date = new Date(startTime)
  if (Number.isNaN(date.getTime())) return ''
  const h = date.getHours()
  const min = date.getMinutes()
  if (h === 0 && min === 0) return ''
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

/** 构建月历网格日期列表 */
export function buildCalendarDays(month: Date, events: CalendarEvent[]) {
  const year = month.getFullYear()
  const m = month.getMonth()
  const firstDay = new Date(year, m, 1)
  const lastDay = new Date(year, m + 1, 0)
  const startDow = firstDay.getDay()
  const days: { date: Date; events: CalendarEvent[] }[] = []

  for (let i = startDow - 1; i >= 0; i--) {
    const date = new Date(year, m, -i)
    days.push({ date, events: getEventsForDay(date, events) })
  }
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = new Date(year, m, d)
    days.push({ date, events: getEventsForDay(date, events) })
  }
  while (days.length % 7 !== 0) {
    const last = days[days.length - 1].date
    const date = new Date(last.getFullYear(), last.getMonth(), last.getDate() + 1)
    days.push({ date, events: getEventsForDay(date, events) })
  }
  return days
}
