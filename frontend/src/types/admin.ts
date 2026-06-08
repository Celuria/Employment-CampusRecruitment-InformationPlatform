import type { EventFormat } from './careerTalk'
import type { UserRole } from './auth'

export type PublishStatus = 'draft' | 'published' | 'archived'
export type UserStatus = 'active' | 'locked' | 'disabled'

export interface AdminCareerTalkVO {
  id: number
  title: string
  company: string
  industryCode?: string
  companySize?: string
  startTime: string
  endTime?: string
  location: string
  campus?: string
  format: EventFormat
  positions: string[]
  targetMajors?: string[]
  registrationUrl?: string
  sourceUrl?: string
  logoUrl?: string
  description?: string
  publishStatus: PublishStatus
  sourceType: string
  createdBy?: number
  updatedBy?: number
  createdAt: string
  updatedAt: string
}

export interface AdminCareerTalkQuery {
  keyword?: string
  campus?: string
  industry?: string
  publishStatus?: PublishStatus
  sourceType?: string
  page?: number
  pageSize?: number
}

export interface AdminCareerTalkForm {
  title: string
  company: string
  industryCode?: string
  companySize?: string
  startTime: string
  endTime?: string
  location: string
  campus?: string
  format: EventFormat
  positions: string[]
  targetMajors: string[]
  registrationUrl?: string
  sourceUrl?: string
  logoUrl?: string
  description?: string
  publishStatus?: PublishStatus
}

export interface AdminJobFairVO {
  id: number
  title: string
  startDate: string
  endDate?: string
  startTime?: string
  location: string
  campus?: string
  companyCount?: number
  targetAudience?: string
  targetMajors?: string[]
  deadline?: string
  detailUrl?: string
  sourceUrl?: string
  description?: string
  publishStatus: PublishStatus
  sourceType: string
  createdBy?: number
  updatedBy?: number
  createdAt: string
  updatedAt: string
}

export interface AdminJobFairQuery {
  keyword?: string
  campus?: string
  publishStatus?: PublishStatus
  sourceType?: string
  page?: number
  pageSize?: number
}

export interface AdminJobFairForm {
  title: string
  startDate: string
  endDate?: string
  startTime?: string
  location: string
  campus?: string
  companyCount?: number
  targetAudience?: string
  targetMajors: string[]
  deadline?: string
  detailUrl?: string
  sourceUrl?: string
  description?: string
  publishStatus?: PublishStatus
}

export interface AdminUserVO {
  id: number
  username: string
  name: string
  email: string
  college?: string
  major?: string
  role: UserRole
  status: UserStatus
  lastLoginAt?: string
  createdAt: string
}

export interface AdminUserQuery {
  keyword?: string
  role?: UserRole
  status?: UserStatus
  page?: number
  pageSize?: number
}

export interface AdminUserCreateForm {
  username: string
  password: string
  name: string
  email: string
  role: UserRole
}

export interface AdminUserUpdateForm {
  name?: string
  email?: string
  college?: string
  major?: string
  role?: UserRole
}

export interface SyncTaskVO {
  taskId: string
  status: string
  startedAt: string
  message: string
}

export interface SyncLogVO {
  id: number
  taskId: string
  sourceType: string
  status: string
  addedCount: number
  updatedCount: number
  failedCount: number
  startedAt: string
  finishedAt?: string
  operatorId: number
  errorMessage?: string
}

export interface AuditLogVO {
  id: number
  operatorId: number
  operatorName: string
  action: string
  resourceType: string
  resourceId: number
  detail: string
  ip: string
  createdAt: string
}

export interface AuditLogQuery {
  operatorId?: number
  action?: string
  resourceType?: string
  startDate?: string
  endDate?: string
  page?: number
  pageSize?: number
}
