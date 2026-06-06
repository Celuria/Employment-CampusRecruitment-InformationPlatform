import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores'

/** 认证相关组合式函数 */
export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()

  const isLoggedIn = computed(() => authStore.isLoggedIn)
  const userInfo = computed(() => authStore.userInfo)

  async function requireAuth(redirect?: string) {
    if (!authStore.isLoggedIn) {
      await router.push({
        name: 'Login',
        query: redirect ? { redirect } : undefined,
      })
      return false
    }
    return true
  }

  return {
    isLoggedIn,
    userInfo,
    requireAuth,
    login: authStore.login,
    logout: authStore.logout,
  }
}
