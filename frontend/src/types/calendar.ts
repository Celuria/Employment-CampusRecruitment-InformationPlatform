export type EventType = 'career_talk' | 'job_fair'

export type RemindBefore = '1h' | '1d' | '3d'

export interface CalendarEvent {
  id: number
  eventType: EventType
  refId: number
  title: string
  startTime: string
  endTime?: string
  location: string
  customNote?: string
  remindBefore: RemindBefore[]
  reminderStatus?: 'pending' | 'sent' | 'failed' | 'cancelled'
  sourceUpdated?: boolean
  createdAt?: string
}

export interface CalendarEventQuery {
  startDate?: string
  endDate?: string
  eventType?: EventType
  view?: 'list' | 'month'
}

export interface UserPreference {
  targetPositions: string[]
  preferredCities: string[]
  preferredCompanies: string[]
  focusCompanies: string[]
  remindBefore: RemindBefore[]
}

export interface RecommendationItem {
  id: number
  eventType: EventType
  refId: number
  title: string
  company?: string
  startTime: string
  startDate?: string
  endDate?: string
  location: string
  industry?: string
  industryCode?: string
  companySize?: string
  format?: 'online' | 'offline' | 'hybrid'
  positions?: string[]
  status?: 'upcoming' | 'ended'
  companyCount?: number
  targetAudience?: string
  deadline?: string
  detailUrl?: string
  matchScore?: number
  matchReasons: string[]
  inCalendar?: boolean
}

export interface RecommendationResult {
  list: RecommendationItem[]
  total: number
  page: number
  pageSize: number
  fallback: boolean
}

export interface ReminderLog {
  id: number
  calendarEventId?: number
  eventTitle: string
  eventType?: string
  remindBefore?: string
  scheduledTime: string
  sentTime?: string
  status: 'pending' | 'sent' | 'failed' | 'cancelled'
  retryCount?: number
  failReason?: string
  createdAt?: string
}
