/** 转为 el-date-picker datetime 的 value-format：YYYY-MM-DDTHH:mm:ss */
export function toDatePickerDateTime(value?: string | null): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

/** 转为 el-date-picker date 的 value-format：YYYY-MM-DD */
export function toDatePickerDate(value?: string | null): string {
  if (!value) return ''
  if (/^\d{4}-\d{2}-\d{2}$/.test(value)) return value
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

/** 日期选择器值转为 API 提交格式 */
export function toApiDateTime(value?: string): string | undefined {
  if (!value?.trim()) return undefined
  return value.replace('T', ' ')
}

/** 格式化日期时间（占位，后续接入 dayjs） */
export function formatDateTime(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** 邮箱脱敏 */
export function maskEmail(email: string): string {
  const [name, domain] = email.split('@')
  if (!name || !domain) return email
  const visible = name.slice(0, Math.min(2, name.length))
  return `${visible}***@${domain}`
}
