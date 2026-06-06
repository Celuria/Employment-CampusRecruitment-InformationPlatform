/** 更新个人资料请求体 */
export interface UpdateProfileForm {
  name: string
  college: string
  major: string
  grade?: string
  targetPositions: string[]
  phone?: string
  email: string
}
