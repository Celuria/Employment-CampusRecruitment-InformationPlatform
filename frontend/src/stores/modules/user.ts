import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPreferencesApi, updatePreferencesApi } from '@/api/modules/user'
import type { UserPreference } from '@/types'

export const useUserStore = defineStore('user', () => {
  const preferences = ref<UserPreference | null>(null)
  const loading = ref(false)
  const profileCompleted = ref(false)

  async function fetchPreferences() {
    loading.value = true
    try {
      preferences.value = await getPreferencesApi()
      return preferences.value
    } finally {
      loading.value = false
    }
  }

  async function updatePreferences(data: UserPreference) {
    loading.value = true
    try {
      preferences.value = await updatePreferencesApi(data)
      return preferences.value
    } finally {
      loading.value = false
    }
  }

  function setProfileCompleted(completed: boolean) {
    profileCompleted.value = completed
  }

  return {
    preferences,
    loading,
    profileCompleted,
    fetchPreferences,
    updatePreferences,
    setProfileCompleted,
  }
})
