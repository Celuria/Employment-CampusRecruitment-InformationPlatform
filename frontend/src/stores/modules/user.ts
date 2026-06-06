import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPreferencesApi } from '@/api/modules/user'
import type { UserPreference } from '@/types'

export const useUserStore = defineStore('user', () => {
  const preferences = ref<UserPreference | null>(null)
  const profileCompleted = ref(false)

  async function fetchPreferences() {
    preferences.value = await getPreferencesApi()
  }

  function setProfileCompleted(completed: boolean) {
    profileCompleted.value = completed
  }

  return {
    preferences,
    profileCompleted,
    fetchPreferences,
    setProfileCompleted,
  }
})
