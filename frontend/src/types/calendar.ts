export type EventType = 'career_talk' | 'job_fair'

export type RemindBefore = '1h' | '1d' | '3d'

export interface CalendarEvent {
  id: number
  eventType: EventType
  refId: number
  title: string
  startTime: string
  location: string
  customNote?: string
  remindBefore: RemindBefore
  reminderStatus?: 'pending' | 'sent' | 'failed'
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
  title: string
  company: string
  startTime: string
  location: string
  matchReasons: string[]
  inCalendar?: boolean
}

export interface ReminderLog {
  id: number
  eventTitle: string
  scheduledTime: string
  sentTime?: string
  status: 'success' | 'failed'
  retryCount?: number
}
