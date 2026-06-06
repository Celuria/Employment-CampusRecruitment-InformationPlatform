<script setup lang="ts">
defineProps<{
  page: number
  totalPages: number
}>()

defineEmits<{
  change: [page: number]
}>()
</script>

<template>
  <div class="mt-8 flex items-center justify-center gap-2">
    <button
      type="button"
      class="rounded-lg border border-ink-200 px-3 py-2 text-sm text-ink-400 transition-colors hover:bg-ink-50 disabled:cursor-not-allowed"
      :disabled="page <= 1"
      @click="$emit('change', page - 1)"
    >
      ‹
    </button>
    <button
      v-for="p in Math.min(totalPages, 5)"
      :key="p"
      type="button"
      class="rounded-lg px-4 py-2 text-sm font-medium transition-colors"
      :class="
        p === page
          ? 'bg-brand-500 text-white'
          : 'border border-ink-200 text-ink-600 hover:bg-ink-50'
      "
      @click="$emit('change', p)"
    >
      {{ p }}
    </button>
    <button
      type="button"
      class="rounded-lg border border-ink-200 px-3 py-2 text-sm text-ink-600 transition-colors hover:bg-ink-50 disabled:cursor-not-allowed"
      :disabled="page >= totalPages"
      @click="$emit('change', page + 1)"
    >
      ›
    </button>
  </div>
</template>
