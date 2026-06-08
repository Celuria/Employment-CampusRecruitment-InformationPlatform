import type { Router } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'

export function setupRouterGuards(router: Router) {
  router.beforeEach(async (to, _from, next) => {
    const authStore = useAuthStore()
    const title = to.meta.title
    document.title = title ? `${title} - 就业中心校招信息平台` : '就业中心校招信息平台'

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
    const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin)

    if (requiresAuth && !authStore.isLoggedIn) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    if (requiresAdmin) {
      if (!authStore.userInfo) {
        try {
          await authStore.fetchUserInfo()
        } catch {
          next({ name: 'Login', query: { redirect: to.fullPath } })
          return
        }
      }
      if (!authStore.isAdmin) {
        ElMessage.warning('需要管理员权限')
        next({ name: 'Home' })
        return
      }
    }

    if (to.meta.guestOnly && authStore.isLoggedIn) {
      next({ name: 'Home' })
      return
    }

    next()
  })
}
