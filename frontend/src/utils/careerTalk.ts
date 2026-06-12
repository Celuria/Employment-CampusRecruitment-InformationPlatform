import { INDUSTRY_OPTIONS } from '@/constants'
import { formatEventLocation, getCampusLabel, normalizeCampusValue } from '@/utils/location'
import type { CareerTalk, JobFair } from '@/types'

const industryLabelMap = Object.fromEntries(
  INDUSTRY_OPTIONS.filter((o) => o.value !== 'all').map((o) => [o.value, o.label]),
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
  const normalized = normalizeCampusValue(campus)
  const campusShort = normalized ? getCampusLabel(normalized).replace('校区', '') : ''

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
  const campus = normalizeCampusValue(talk.campus)

  return {
    ...talk,
    campus,
    industry: talk.industry || industryLabelMap[talk.industryCode ?? ''] || talk.industryCode || '',
    status: talk.status ?? (isEnded ? 'ended' : 'upcoming'),
    positions: talk.positions ?? [],
    inCalendar: talk.inCalendar ?? false,
    location: formatEventLocation(campus, talk.venue, talk.location),
  }
}

export function normalizeJobFair(fair: JobFair): JobFair {
  const campus = normalizeCampusValue(fair.campus)
  return {
    ...fair,
    campus,
    location: formatEventLocation(campus, fair.venue, fair.location),
    inCalendar: fair.inCalendar ?? false,
  }
}
