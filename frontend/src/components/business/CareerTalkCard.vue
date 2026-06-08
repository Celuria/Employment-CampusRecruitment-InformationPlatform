<script setup lang="ts">
import { formatDateTime } from '@/utils/format'
import type { CareerTalk } from '@/types'

defineProps<{
  item: CareerTalk
}>()

defineEmits<{
  addToCalendar: [id: number]
  toggleFavorite: [id: number]
  viewDetail: [id: number]
}>()

const formatTagClass: Record<string, string> = {
  online: 'tag-online',
  offline: 'tag-offline',
  hybrid: 'tag-hybrid',
}

const formatLabel: Record<string, string> = {
  online: '纯线上',
  offline: '线下专场',
  hybrid: '线上+线下',
}
</script>

<template>
  <article
    class="card-shadow group flex cursor-pointer gap-5 rounded-2xl bg-white p-6 transition-all hover:shadow-md"
    :class="{
      'border-2 border-brand-200': item.inCalendar,
      'opacity-75': item.status === 'ended',
    }"
    @click="$emit('viewDetail', item.id)"
  >
    <div class="shrink-0">
      <div
        class="company-logo-placeholder flex h-20 w-20 items-center justify-center rounded-xl border border-brand-100"
        :class="{ 'border-ink-200 grayscale': item.status === 'ended' }"
      >
        <span class="text-sm font-bold text-brand-700">{{ item.company.slice(0, 2) }}</span>
      </div>
    </div>
    <div class="min-w-0 flex-1">
      <div class="mb-2 flex items-start justify-between">
        <div class="min-w-0">
          <h3
            class="mb-1 truncate text-lg font-bold text-ink-900 transition-colors group-hover:text-brand-600"
            :class="item.status === 'ended' ? 'text-ink-700' : ''"
          >
            {{ item.title }}
          </h3>
          <p class="text-sm text-ink-500">
            {{ item.company }} · {{ item.industry }} · {{ item.companySize }}
          </p>
        </div>
        <div class="ml-4 flex shrink-0 gap-2" @click.stop>
          <template v-if="item.status === 'ended'">
            <span class="rounded-lg bg-ink-100 px-3 py-2 text-xs font-medium text-ink-400">
              已结束
            </span>
          </template>
          <template v-else>
            <button
              type="button"
              class="rounded-lg border p-2 transition-all"
              :class="
                item.favorited
                  ? 'border-red-200 bg-red-50'
                  : 'border-ink-200 hover:border-brand-300 hover:bg-brand-50'
              "
              :title="item.favorited ? '已收藏' : '收藏'"
              @click="$emit('toggleFavorite', item.id)"
            >
              <svg
                class="h-4 w-4"
                :class="item.favorited ? 'text-red-500' : 'text-ink-400'"
                :fill="item.favorited ? 'currentColor' : 'none'"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                />
              </svg>
            </button>
            <button
              v-if="!item.inCalendar"
              type="button"
              class="flex items-center gap-1.5 rounded-lg bg-brand-500 px-4 py-2 text-xs font-medium text-white transition-colors hover:bg-brand-600"
              @click="$emit('addToCalendar', item.id)"
            >
              加入日历
            </button>
            <button
              v-else
              type="button"
              class="flex items-center gap-1.5 rounded-lg border border-brand-300 bg-brand-100 px-4 py-2 text-xs font-medium text-brand-700"
              disabled
            >
              已加入日历
            </button>
          </template>
        </div>
      </div>
      <div class="mb-3 flex items-center gap-4 text-sm text-ink-600">
        <span>{{ formatDateTime(item.startTime) }}</span>
        <span>{{ item.location }}</span>
      </div>
      <div class="flex items-center gap-2">
        <span
          class="rounded-md px-2.5 py-1 text-xs font-medium"
          :class="formatTagClass[item.format]"
        >
          {{ formatLabel[item.format] }}
        </span>
        <span
          v-for="pos in item.positions"
          :key="pos"
          class="rounded-md bg-ink-100 px-2.5 py-1 text-xs font-medium text-ink-600"
        >
          {{ pos }}
        </span>
      </div>
    </div>
  </article>
</template>
