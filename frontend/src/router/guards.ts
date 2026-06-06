import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores'

export function setupRouterGuards(router: Router) {
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()
    const title = to.meta.title
    document.title = title ? `${title} - 就业中心校招信息平台` : '就业中心校招信息平台'

    const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)

    if (requiresAuth && !authStore.isLoggedIn) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    if (to.meta.guestOnly && authStore.isLoggedIn) {
      next({ name: 'Home' })
      return
    }

    next()
  })
}
