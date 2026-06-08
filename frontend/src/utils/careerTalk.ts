import { CAMPUS_OPTIONS, INDUSTRY_OPTIONS } from '@/constants'
import type { CareerTalk } from '@/types'

const industryLabelMap = Object.fromEntries(
  INDUSTRY_OPTIONS.filter((o) => o.value !== 'all').map((o) => [o.value, o.label]),
)

const campusLabelMap = Object.fromEntries(
  CAMPUS_OPTIONS.filter((o) => o.value !== 'all').map((o) => [o.value, o.label.replace('校区', '')]),
)

const UPCOMING_COLORS = [
  'bg-brand-100 text-brand-700',
  'bg-orange-100 text-orange-700',
  'bg-blue-100 text-blue-700',
  'bg-emerald-100 text-emerald-700',
  'bg-violet-100 text-violet-700',
]

function isSameDay(a: Date, b: Date) {
  return a.getFullYear() === b.getFullYear()
    && a.getMonth() === b.getMonth()
    && a.getDate() === b.getDate()
}

/** 侧边栏「即将开始」时间文案 */
export function formatUpcomingLabel(startTime: string, campus?: string) {
  const date = new Date(startTime)
  if (Number.isNaN(date.getTime())) return ''

  const now = new Date()
  const timeStr = date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  const campusShort = campus ? (campusLabelMap[campus] ?? campus) : ''

  let dayLabel: string
  if (isSameDay(date, now)) {
    dayLabel = '今天'
  } else {
    const tomorrow = new Date(now)
    tomorrow.setDate(tomorrow.getDate() + 1)
    dayLabel = isSameDay(date, tomorrow)
      ? '明天'
      : date.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric' })
  }

  return campusShort ? `${dayLabel} ${timeStr} · ${campusShort}` : `${dayLabel} ${timeStr}`
}

export function getUpcomingColor(index: number) {
  return UPCOMING_COLORS[index % UPCOMING_COLORS.length]
}

/** 将 API 返回的宣讲会数据规范化为前端展示格式 */
export function normalizeCareerTalk(talk: CareerTalk): CareerTalk {
  const startMs = new Date(talk.startTime).getTime()
  const isEnded = !Number.isNaN(startMs) && startMs < Date.now()

  return {
    ...talk,
    industry: talk.industry || industryLabelMap[talk.industryCode ?? ''] || talk.industryCode || '',
    status: talk.status ?? (isEnded ? 'ended' : 'upcoming'),
    positions: talk.positions ?? [],
  }
}
