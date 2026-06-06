import { get } from '../request'
import type { CareerTalk, CareerTalkQuery, PageResult } from '@/types'

/** 宣讲会列表 */
export function getCareerTalkListApi(params: CareerTalkQuery) {
  return get<PageResult<CareerTalk>>('/career-talks', { params })
}

/** 宣讲会详情 */
export function getCareerTalkDetailApi(id: number) {
  return get<CareerTalk>(`/career-talks/${id}`)
}
