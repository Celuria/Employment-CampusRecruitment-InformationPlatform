import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi, logoutApi, registerApi } from '@/api'
import { getProfileApi, updateProfileApi } from '@/api/modules/user'
import { storage } from '@/utils/storage'
import type { LoginForm, RegisterForm, UpdateProfileForm, UserInfo } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(storage.getToken())
  const userInfo = ref<UserInfo | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

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

  async function register(form: RegisterForm) {
    loading.value = true
    try {
      const { confirmPassword: _, ...payload } = form
      await registerApi(payload)
    } finally {
      loading.value = false
    }
  }

  async function updateProfile(form: UpdateProfileForm) {
    loading.value = true
    try {
      userInfo.value = await updateProfileApi(form)
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
      if (token.value) {
        await logoutApi()
      }
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
    isAdmin,
    login,
    register,
    updateProfile,
    logout,
    fetchUserInfo,
    initAuth,
  }
})
