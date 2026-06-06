import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi, logoutApi } from '@/api'
import { getProfileApi } from '@/api/modules/user'
import { storage } from '@/utils/storage'
import type { LoginForm, UserInfo } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getToken())
  const userInfo = ref<UserInfo | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function login(form: LoginForm) {
    loading.value = true
    try {
      const result = await loginApi(form)
      token.value = result.token
      storage.setToken(result.token, form.remember)
      await fetchUserInfo()
    } finally {
      loading.value = false
    }
  }

  async function fetchUserInfo() {
    if (!token.value) return
    userInfo.value = await getProfileApi()
  }

  async function logout() {
    try {
      await logoutApi()
    } catch {
      // 忽略退出接口失败，仍清理本地状态
    } finally {
      token.value = null
      userInfo.value = null
      storage.removeToken()
    }
  }

  function initAuth() {
    if (token.value) {
      fetchUserInfo().catch(() => {
        token.value = null
        storage.removeToken()
      })
    }
  }

  return {
    token,
    userInfo,
    loading,
    isLoggedIn,
    login,
    logout,
    fetchUserInfo,
    initAuth,
  }
})
