import { post } from '../request'
import type { AuthToken, LoginForm, RegisterForm } from '@/types'

/** 用户登录 */
export function loginApi(data: LoginForm) {
  return post<AuthToken>('/auth/login', data)
}

/** 用户注册 */
export function registerApi(data: Omit<RegisterForm, 'confirmPassword'>) {
  return post<void>('/auth/register', data)
}

/** 退出登录 */
export function logoutApi() {
  return post<void>('/auth/logout')
}
