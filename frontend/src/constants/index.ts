export const APP_NAME = '就业中心校招信息平台'

export const API_PREFIX = '/api/v1'

export const PAGE_SIZE = 10

export const DATE_FILTER_OPTIONS = [
  { label: '全部', value: 'all' },
  { label: '今天', value: 'today' },
  { label: '明天', value: 'tomorrow' },
  { label: '本周', value: 'this_week' },
  { label: '下周', value: 'next_week' },
] as const

export const INDUSTRY_OPTIONS = [
  { label: '全部', value: 'all' },
  { label: '互联网', value: 'internet' },
  { label: '金融', value: 'finance' },
  { label: '制造', value: 'manufacturing' },
  { label: '咨询', value: 'consulting' },
] as const

export const CAMPUS_OPTIONS = [
  { label: '全部校区', value: 'all' },
  { label: '本部校区', value: 'main' },
  { label: '沙河校区', value: 'shahe' },
  { label: '线上', value: 'online' },
] as const

export const SORT_OPTIONS = [
  { label: '按时间最近', value: 'time_asc' },
  { label: '按热度', value: 'popularity' },
  { label: '按公司规模', value: 'company_size' },
] as const

export const REMIND_BEFORE_OPTIONS = [
  { label: '提前 1 小时', value: '1h' },
  { label: '提前 1 天', value: '1d' },
  { label: '提前 3 天', value: '3d' },
] as const
