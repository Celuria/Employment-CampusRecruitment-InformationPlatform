<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
}>()
</script>

<template>
  <div class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-16 text-center">
    <div class="mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-brand-50">
      <svg class="h-7 w-7 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"
        />
      </svg>
    </div>
    <h3 v-if="title" class="mb-2 text-lg font-semibold text-ink-800">{{ title }}</h3>
    <p class="max-w-sm text-sm text-ink-500">
      <slot>{{ description || '页面框架已就绪，业务功能待接入' }}</slot>
    </p>
  </div>
</template>
