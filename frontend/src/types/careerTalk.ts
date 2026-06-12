export type EventFormat = 'online' | 'offline' | 'hybrid'

export interface CareerTalk {
  id: number
  title: string
  company: string
  industry?: string
  industryCode?: string
  companySize?: string
  startTime: string
  endTime?: string
  location: string
  campus?: string
  venue?: string
  format: EventFormat
  positions: string[]
  targetMajors?: string[]
  registrationUrl?: string
  sourceUrl?: string
  publishedAt?: string
  logoUrl?: string
  description?: string
  syncedAt?: string
  contactInfo?: string
  status?: 'upcoming' | 'ended'
  inCalendar?: boolean
  favorited?: boolean
}

export interface CareerTalkQuery {
  keyword?: string
  dateRange?: string
  campus?: string
  industry?: string
  company?: string
  sortBy?: string
  page?: number
  pageSize?: number
}

export interface HotCompany {
  company: string
  companySize?: string
}

export interface JobFair {
  id: number
  title: string
  startDate: string
  endDate?: string
  location: string
  campus?: string
  venue?: string
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
  campus?: string
  sortBy?: string
  page?: number
  pageSize?: number
}
