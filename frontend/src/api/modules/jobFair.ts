import { get } from '../request'
import type { JobFair, JobFairQuery, PageResult } from '@/types'

/** 双选会列表 */
export function getJobFairListApi(params: JobFairQuery) {
  return get<PageResult<JobFair>>('/job-fairs', { params })
}

/** 双选会详情 */
export function getJobFairDetailApi(id: number) {
  return get<JobFair>(`/job-fairs/${id}`)
}
