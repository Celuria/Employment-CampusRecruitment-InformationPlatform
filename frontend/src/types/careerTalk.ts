export type EventFormat = 'online' | 'offline' | 'hybrid'

export interface CareerTalk {
  id: number
  title: string
  company: string
  industry: string
  companySize?: string
  startTime: string
  location: string
  format: EventFormat
  positions: string[]
  sourceUrl?: string
  publishedAt?: string
  logoUrl?: string
  status?: 'upcoming' | 'ended'
  inCalendar?: boolean
  favorited?: boolean
}

export interface CareerTalkQuery {
  keyword?: string
  dateRange?: string
  campus?: string
  industry?: string
  sortBy?: string
  page?: number
  pageSize?: number
}

export interface JobFair {
  id: number
  title: string
  startDate: string
  endDate?: string
  location: string
  companyCount?: number
  targetAudience?: string
  deadline?: string
  detailUrl?: string
  inCalendar?: boolean
}

export interface JobFairQuery {
  keyword?: string
  startDate?: string
  endDate?: string
  page?: number
  pageSize?: number
}
