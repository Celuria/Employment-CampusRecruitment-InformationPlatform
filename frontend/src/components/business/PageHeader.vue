<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  count?: number
  countLabel?: string
}>()
</script>

<template>
  <div class="mb-8">
    <div class="mb-2 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-ink-900">{{ title }}</h1>
        <p v-if="description" class="mt-1 text-sm text-ink-500">{{ description }}</p>
      </div>
      <div v-if="count !== undefined" class="flex items-center gap-2 text-sm text-ink-500">
        <span>共找到</span>
        <span class="text-lg font-bold text-brand-600">{{ count }}</span>
        <span>{{ countLabel }}</span>
      </div>
    </div>
  </div>
</template>
