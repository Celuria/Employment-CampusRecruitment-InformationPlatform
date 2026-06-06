import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia

export { useAuthStore } from './modules/auth'
export { useUserStore } from './modules/user'
