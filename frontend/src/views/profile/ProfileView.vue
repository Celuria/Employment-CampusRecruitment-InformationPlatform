<script setup lang="ts">
import { RouterLink } from 'vue-router'

const tabs = [
  { name: 'ProfileInfo', label: '基本资料', path: '/profile/info' },
  { name: 'ProfilePreferences', label: '偏好设置', path: '/profile/preferences' },
  { name: 'ProfileReminders', label: '提醒记录', path: '/profile/reminders' },
]
</script>

<template>
  <div>
    <h1 class="mb-6 text-2xl font-bold text-ink-900">个人中心</h1>
    <div class="flex gap-8">
      <nav class="w-48 shrink-0 space-y-1">
        <RouterLink
          v-for="tab in tabs"
          :key="tab.name"
          :to="tab.path"
          class="block rounded-lg px-4 py-2.5 text-sm font-medium transition-colors"
          active-class="bg-brand-50 text-brand-600"
          inactive-class="text-ink-600 hover:bg-ink-100"
        >
          {{ tab.label }}
        </RouterLink>
      </nav>
      <div class="min-w-0 flex-1">
        <RouterView />
      </div>
    </div>
  </div>
</template>
