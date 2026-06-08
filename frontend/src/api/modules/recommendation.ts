import { get } from '../request'
import type { RecommendationResult } from '@/types'

/** 个性化推荐列表 */
export function getRecommendationsApi(params?: { page?: number; pageSize?: number }) {
  return get<RecommendationResult>('/recommendations', { params })
}
