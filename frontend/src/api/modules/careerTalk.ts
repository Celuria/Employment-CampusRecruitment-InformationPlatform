import { get } from '../request'
import type { CareerTalk, CareerTalkQuery, HotCompany, PageResult } from '@/types'

/** 宣讲会列表 */
export function getCareerTalkListApi(params: CareerTalkQuery) {
  return get<PageResult<CareerTalk>>('/career-talks', { params })
}

/** 宣讲会详情 */
export function getCareerTalkDetailApi(id: number) {
  return get<CareerTalk>(`/career-talks/${id}`)
}

/** 24 小时内即将开始的宣讲会（侧边栏） */
export function getCareerTalkUpcomingApi() {
  return get<CareerTalk[]>('/career-talks/upcoming')
}

/** 热门公司 Top N（按公司规模降序） */
export function getCareerTalkHotCompaniesApi(limit = 6) {
  return get<HotCompany[]>('/career-talks/hot-companies', { params: { limit } })
}
