/** 密码规则：长度 ≥ 8，含字母与数字 */
export function validatePassword(password: string): boolean {
  return password.length >= 8 && /[A-Za-z]/.test(password) && /\d/.test(password)
}

/** 账号规则：4-32 位字母、数字、下划线 */
export function validateUsername(username: string): boolean {
  return /^[A-Za-z0-9_]{4,32}$/.test(username)
}

export const passwordRuleMessage = '密码长度不少于 8 位，且需包含字母和数字'
export const usernameRuleMessage = '账号为 4-32 位字母、数字或下划线'
