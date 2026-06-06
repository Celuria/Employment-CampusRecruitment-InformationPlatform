import type { RemindBefore } from '@/types'

/** 常见偏好城市 */
export const CITY_SUGGESTIONS = [
  '北京',
  '上海',
  '深圳',
  '杭州',
  '广州',
  '成都',
  '南京',
  '武汉',
  '西安',
  '苏州',
  '重庆',
  '天津',
]

/** 常见偏好公司 */
export const COMPANY_SUGGESTIONS = [
  '字节跳动',
  '腾讯',
  '阿里巴巴',
  '华为',
  '美团',
  '京东',
  '百度',
  '网易',
  '小米',
  '滴滴',
  '快手',
  '拼多多',
]

/** 日历提醒提前量选项 */
export const REMIND_BEFORE_OPTIONS: { label: string; value: RemindBefore }[] = [
  { label: '提前 1 小时', value: '1h' },
  { label: '提前 1 天', value: '1d' },
  { label: '提前 3 天', value: '3d' },
]
