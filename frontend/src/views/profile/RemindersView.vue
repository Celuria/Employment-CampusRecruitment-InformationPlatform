<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getReminderLogsApi } from '@/api'
import AppPagination from '@/components/common/AppPagination.vue'
import type { ReminderLog } from '@/types'

const list = ref<ReminderLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(false)

const totalPages = ref(1)

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

function statusLabel(status: string): string {
  switch (status) {
    case 'pending': return '待提醒'
    case 'sent': return '已提醒'
    case 'failed': return '失败'
    case 'cancelled': return '已取消'
    default: return status
  }
}

function remindLabel(rb: string): string {
  switch (rb) {
    case '1h': return '提前 1 小时'
    case '1d': return '提前 1 天'
    case '3d': return '提前 3 天'
    default: return rb
  }
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
      <div v-if="list.length === 0" class="flex flex-col items-center py-16 text-ink-400">
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
              <th class="px-5 py-3">提醒方式</th>
              <th class="px-5 py-3">计划提醒时间</th>
              <th class="px-5 py-3">状态</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-ink-100">
            <tr
              v-for="item in list"
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
                {{ item.remindBefore ? remindLabel(item.remindBefore) : '-' }}
              </td>
              <td class="px-5 py-3 text-ink-500">
                {{ formatTime(item.scheduledTime) }}
              </td>
              <td class="px-5 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="{
                    'bg-yellow-100 text-yellow-700': item.status === 'pending',
                    'bg-green-100 text-green-700': item.status === 'sent',
                    'bg-red-100 text-red-600': item.status === 'failed',
                    'bg-ink-100 text-ink-500': item.status === 'cancelled',
                  }"
                >
                  {{ statusLabel(item.status) }}
                </span>
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
