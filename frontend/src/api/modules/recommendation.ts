import { get } from '../request'
import type { RecommendationItem } from '@/types'

/** 个性化推荐列表 */
export function getRecommendationsApi() {
  return get<RecommendationItem[]>('/recommendations')
}
