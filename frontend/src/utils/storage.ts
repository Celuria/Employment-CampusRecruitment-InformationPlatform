const TOKEN_KEY = 'campus_recruit_token'
const REMEMBER_KEY = 'campus_recruit_remember'

export const storage = {
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken(token: string, remember = false): void {
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(REMEMBER_KEY, String(remember))
  },

  removeToken(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REMEMBER_KEY)
  },

  isRemembered(): boolean {
    return localStorage.getItem(REMEMBER_KEY) === 'true'
  },
}
