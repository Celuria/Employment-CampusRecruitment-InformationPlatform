/** 年级选项 */
export const GRADE_OPTIONS = [
  { label: '2026 级', value: '2026' },
  { label: '2025 级', value: '2025' },
  { label: '2024 级', value: '2024' },
  { label: '2023 级', value: '2023' },
  { label: '2022 级', value: '2022' },
  { label: '2021 级', value: '2021' },
] as const

/** 学院与专业字典 */
export const COLLEGE_MAJOR_MAP: Record<string, string[]> = {
  计算机学院: ['软件工程', '计算机科学与技术', '网络工程', '信息安全', '数据科学'],
  经济管理学院: ['工商管理', '会计学', '金融学', '市场营销', '人力资源管理'],
  机械工程学院: ['机械设计制造及其自动化', '工业设计', '车辆工程', '智能制造'],
  电子信息学院: ['电子信息工程', '通信工程', '微电子科学与工程', '自动化'],
  外国语学院: ['英语', '日语', '翻译', '商务英语'],
  数学学院: ['数学与应用数学', '信息与计算科学', '统计学'],
}

export const COLLEGE_OPTIONS = Object.keys(COLLEGE_MAJOR_MAP)

/** 常见意向岗位（可自定义输入） */
export const POSITION_SUGGESTIONS = [
  'Java 开发',
  '后端开发',
  '前端开发',
  '产品经理',
  '数据分析',
  '算法工程师',
  '测试工程师',
  '运维工程师',
  'UI 设计',
  '运营',
]
