<script setup lang="ts">
import type { CalendarEvent } from '@/types'

defineProps<{
  date: Date
  events: CalendarEvent[]
  deleteLoadingId: number | null
}>()

const emit = defineEmits<{
  close: []
  viewDetail: [event: CalendarEvent]
  edit: [event: CalendarEvent]
  delete: [event: CalendarEvent]
}>()

const REMIND_OPTIONS = [
  { label: '提前 1 小时', value: '1h' },
  { label: '提前 1 天', value: '1d' },
  { label: '提前 3 天', value: '3d' },
]

const EVENT_TYPE_LABEL: Record<string, string> = {
  career_talk: '宣讲会',
  job_fair: '双选会',
}

function formatDateLabel(date: Date) {
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long',
  })
}

function formatDateTime(str: string) {
  return new Date(str).toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <aside class="w-80 shrink-0">
    <div class="sidebar-card sticky top-6 rounded-2xl bg-white p-5">
      <div class="mb-4 flex items-start justify-between gap-2">
        <div>
          <h3 class="text-sm font-bold text-ink-900">当日日程</h3>
          <p class="mt-0.5 text-xs text-ink-500">{{ formatDateLabel(date) }}</p>
        </div>
        <button
          type="button"
          class="rounded-lg p-1.5 text-ink-400 transition-colors hover:bg-ink-50 hover:text-ink-600"
          @click="emit('close')"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <p v-if="events.length === 0" class="py-8 text-center text-xs text-ink-400">
        当天暂无日程
      </p>

      <div v-else class="space-y-3">
        <div
          v-for="event in events"
          :key="event.id"
          class="rounded-xl border border-ink-100 p-3 transition-colors hover:border-brand-200 hover:bg-brand-50/30"
        >
          <div class="mb-2 flex items-start justify-between gap-2">
            <span
              class="shrink-0 rounded-md px-2 py-0.5 text-[10px] font-medium"
              :class="event.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700' : 'bg-blue-100 text-blue-700'"
            >
              {{ EVENT_TYPE_LABEL[event.eventType] }}
            </span>
            <div class="flex shrink-0 gap-1">
              <button
                type="button"
                class="rounded px-2 py-0.5 text-[10px] text-ink-500 hover:text-brand-600"
                @click="emit('viewDetail', event)"
              >
                详情
              </button>
              <button
                type="button"
                class="rounded px-2 py-0.5 text-[10px] text-ink-500 hover:text-brand-600"
                @click="emit('edit', event)"
              >
                编辑
              </button>
              <button
                type="button"
                class="rounded px-2 py-0.5 text-[10px] text-red-500 hover:bg-red-50 disabled:opacity-50"
                :disabled="deleteLoadingId === event.id"
                @click="emit('delete', event)"
              >
                删除
              </button>
            </div>
          </div>
          <p class="mb-1 text-sm font-medium text-ink-900 line-clamp-2">{{ event.title }}</p>
          <p class="text-xs text-ink-500">{{ formatDateTime(event.startTime) }}</p>
          <p v-if="event.location" class="mt-0.5 text-xs text-ink-400">{{ event.location }}</p>
          <div v-if="event.remindBefore?.length" class="mt-2 flex flex-wrap gap-1">
            <span
              v-for="rb in event.remindBefore"
              :key="rb"
              class="rounded bg-ink-100 px-1.5 py-0.5 text-[10px] text-ink-500"
            >
              {{ REMIND_OPTIONS.find((o) => o.value === rb)?.label || rb }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
