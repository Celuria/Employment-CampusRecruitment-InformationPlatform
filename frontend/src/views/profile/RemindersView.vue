<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { getReminderLogsApi } from '@/api'
import AppPagination from '@/components/common/AppPagination.vue'
import type { ReminderLog } from '@/types'

const list = ref<ReminderLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(false)

const totalPages = ref(1)

// 客户端过滤：未到时间的 pending 不展示在提醒记录页
const visibleList = computed(() => list.value.filter((item) => item.status !== 'pending'))

async function fetchData() {
  loading.value = true
  try {
    const result = await getReminderLogsApi({ page: page.value, pageSize })
    list.value = result.list
    total.value = result.total
    totalPages.value = Math.max(1, Math.ceil(result.total / result.pageSize))
  } catch {
    list.value = []
    total.value = 0
    totalPages.value = 1
  } finally {
    loading.value = false
  }
}

function onPageChange(p: number) {
  page.value = p
}

function formatTime(iso: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/**
 * 根据 scheduledTime（提醒发出时间）+ remindBefore（提前量）反推活动实际开始时间
 */
function getEventStartTime(item: ReminderLog): string {
  if (!item.scheduledTime || !item.remindBefore) return '-'
  const scheduled = new Date(item.scheduledTime)
  switch (item.remindBefore) {
    case '1h':
      scheduled.setHours(scheduled.getHours() + 1)
      break
    case '1d':
      scheduled.setDate(scheduled.getDate() + 1)
      break
    case '3d':
      scheduled.setDate(scheduled.getDate() + 3)
      break
    default:
      scheduled.setDate(scheduled.getDate() + 1)
  }
  return formatTime(scheduled.toISOString())
}

watch(page, () => fetchData())
onMounted(fetchData)
</script>

<template>
  <div>
    <h2 class="mb-6 text-xl font-semibold text-ink-800">提醒记录</h2>

    <div v-if="loading" class="flex items-center justify-center py-16">
      <span class="text-sm text-ink-400">加载中...</span>
    </div>

    <template v-else>
      <div v-if="visibleList.length === 0" class="flex flex-col items-center py-16 text-ink-400">
        <svg class="mb-4 h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
          />
        </svg>
        <p class="text-sm">暂无提醒记录</p>
        <p class="mt-1 text-xs">将宣讲会或双选会加入日历后，系统会在到期前自动提醒</p>
      </div>

      <div v-else class="overflow-hidden rounded-xl border border-ink-200">
        <table class="w-full text-left text-sm">
          <thead class="bg-ink-50 text-xs font-medium uppercase text-ink-500">
            <tr>
              <th class="px-5 py-3">活动名称</th>
              <th class="px-5 py-3">活动开始时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-ink-100">
            <tr
              v-for="item in visibleList"
              :key="item.id"
              class="transition-colors hover:bg-ink-50"
            >
              <td class="px-5 py-3">
                <span class="font-medium text-ink-800">{{ item.eventTitle }}</span>
                <span
                  v-if="item.eventType"
                  class="ml-2 rounded bg-brand-50 px-1.5 py-0.5 text-xs text-brand-600"
                >
                  {{ item.eventType === 'career_talk' ? '宣讲会' : '双选会' }}
                </span>
              </td>
              <td class="px-5 py-3 text-ink-600">
                {{ getEventStartTime(item) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <AppPagination
        v-if="totalPages > 1"
        :page="page"
        :total-pages="totalPages"
        @change="onPageChange"
      />
    </template>
  </div>
</template>
