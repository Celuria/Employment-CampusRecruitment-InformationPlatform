import { CAMPUS_OPTIONS } from '@/constants'

const campusLabelMap = Object.fromEntries(
  CAMPUS_OPTIONS.filter((o) => o.value !== 'all').map((o) => [o.value, o.label]),
)

/** 校区中文名 */
export function getCampusLabel(campus?: string) {
  if (!campus) return ''
  return campusLabelMap[campus] ?? campus
}

/** 组合校区与楼栋为完整展示地点 */
export function formatEventLocation(campus?: string, venue?: string, fallbackLocation?: string) {
  const v = venue?.trim()
  if (campus) {
    if (campus === 'online') {
      return v ? `线上 · ${v}` : '线上'
    }
    const label = getCampusLabel(campus)
    if (v) return `${label} · ${v}`
    return label || fallbackLocation || ''
  }
  return fallbackLocation || v || ''
}

/** 从完整地点或 venue 字段解析楼栋名（编辑表单回填） */
export function parseVenueFromRecord(row: { campus?: string; venue?: string; location: string }) {
  if (row.venue?.trim()) return row.venue.trim()
  const loc = row.location?.trim() ?? ''
  if (!loc) return ''
  const parts = loc.split(/[·・]/).map((s) => s.trim())
  if (parts.length >= 2) return parts.slice(1).join(' · ')
  if (row.campus === 'online') return loc.replace(/^线上\s*·?\s*/, '')
  return loc
}

/** 旧校区值映射（兼容历史数据） */
export function normalizeCampusValue(campus?: string) {
  if (!campus) return campus
  const legacy: Record<string, string> = { main: 'nanhu', shahe: 'mafangshan' }
  return legacy[campus] ?? campus
}
