import { normalizeCareerTalk } from '@/utils/careerTalk'
import type { CareerTalk, JobFair, RecommendationItem } from '@/types'

export function toCareerTalk(item: RecommendationItem): CareerTalk {
  return normalizeCareerTalk({
    id: item.refId,
    title: item.title,
    company: item.company ?? '',
    industry: item.industry,
    industryCode: item.industryCode,
    companySize: item.companySize,
    startTime: item.startTime,
    location: item.location,
    format: item.format ?? 'offline',
    positions: item.positions ?? [],
    status: item.status,
    inCalendar: item.inCalendar,
  })
}

export function toJobFair(item: RecommendationItem): JobFair {
  return {
    id: item.refId,
    title: item.title,
    startDate: item.startDate ?? item.startTime,
    endDate: item.endDate,
    location: item.location,
    companyCount: item.companyCount,
    targetAudience: item.targetAudience,
    deadline: item.deadline,
    detailUrl: item.detailUrl,
    inCalendar: item.inCalendar,
  }
}
