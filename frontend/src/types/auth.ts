export interface LoginForm {
  username: string
  password: string
  remember?: boolean
}

export interface RegisterForm {
  username: string
  password: string
  confirmPassword: string
  email: string
}

export interface AuthToken {
  token: string
  expiresIn?: number
}

export interface UserInfo {
  id: number
  username: string
  name: string
  email: string
  college: string
  major: string
  grade?: string
  avatar?: string
  targetPositions: string[]
  phone?: string
}
