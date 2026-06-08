<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import { getJobFairListApi } from '@/api/modules/jobFair'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import { usePagination } from '@/composables/usePagination'
import { CAMPUS_OPTIONS } from '@/constants'
import type { JobFair, JobFairQuery } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const { page, pageSize, total, totalPages, setTotal } = usePagination()

const loading = ref(false)
const list = ref<JobFair[]>([])
const calendarLoading = ref<Record<number, boolean>>({})

const filters = reactive({
  keyword: '',
  startDate: '',
  endDate: '',
  campus: 'all',
})

async function fetchList() {
  loading.value = true
  try {
    const params: JobFairQuery = {
      page: page.value,
      pageSize: pageSize.value,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.startDate) params.startDate = filters.startDate
    if (filters.endDate) params.endDate = filters.endDate
    if (filters.campus !== 'all') params.campus = filters.campus

    const res = await getJobFairListApi(params)
    list.value = res.list
    setTotal(Number(res.total))
  } catch {
    // 错误由拦截器处理
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchList()
}

function handleReset() {
  filters.keyword = ''
  filters.startDate = ''
  filters.endDate = ''
  filters.campus = 'all'
  page.value = 1
  fetchList()
}

function setCampus(value: string) {
  filters.campus = value
  handleSearch()
}

function handleDateChange() {
  if (filters.startDate && filters.endDate && filters.startDate > filters.endDate) {
    ElMessage.warning('开始日期不能晚于结束日期')
    return
  }
  handleSearch()
}

function handlePageChange(newPage: number) {
  page.value = newPage
  fetchList()
}

async function handleAddToCalendar(item: JobFair) {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: '/job-fairs' } })
    return
  }
  if (item.inCalendar) {
    ElMessage.info('该活动已在您的日历中')
    return
  }
  calendarLoading.value[item.id] = true
  try {
    await addCalendarEventApi({ eventType: 'job_fair', refId: item.id })
    item.inCalendar = true
    ElMessage.success('已加入日历')
  } catch (e: any) {
    if (e?.response?.data?.code === 40901) {
      item.inCalendar = true
      ElMessage.info('该活动已在您的日历中')
    }
  } finally {
    calendarLoading.value[item.id] = false
  }
}

function handleViewDetail(id: number) {
  router.push(`/job-fairs/${id}`)
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

function isDeadlineSoon(deadline?: string) {
  if (!deadline) return false
  const diff = new Date(deadline).getTime() - Date.now()
  return diff > 0 && diff < 3 * 24 * 3600 * 1000
}

function isDeadlinePassed(deadline?: string) {
  if (!deadline) return false
  return new Date(deadline).getTime() < Date.now()
}

onMounted(fetchList)
</script>

<template>
  <div>
    <PageHeader
      title="双选会"
      description="大型校园招聘会信息，把握与企业双向选择的机会"
      :count="total"
      count-label="场双选会"
    />

    <!-- 筛选栏 -->
    <div class="card-shadow mb-8 rounded-2xl bg-white p-6">
      <div class="mb-5 flex gap-3">
        <div class="relative flex-1">
          <svg class="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-ink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="filters.keyword"
            type="text"
            placeholder="搜索双选会名称..."
            class="w-full rounded-xl border border-ink-200 bg-ink-50 py-3 pl-12 pr-4 text-sm transition-all focus:border-brand-400 focus:outline-none focus:ring-2 focus:ring-brand-300"
            @keyup.enter="handleSearch"
          />
        </div>
        <button type="button" class="btn-primary rounded-xl px-8 py-3 text-sm font-medium text-white" @click="handleSearch">
          搜索
        </button>
        <button type="button" class="rounded-xl border border-ink-200 px-5 py-3 text-sm font-medium text-ink-600 hover:bg-ink-50" @click="handleReset">
          重置
        </button>
      </div>
      <div class="flex flex-wrap items-center gap-6 text-sm">
        <div class="flex items-center gap-2">
          <span class="font-medium text-ink-500">开始日期：</span>
          <input
            v-model="filters.startDate"
            type="date"
            class="rounded-lg border border-ink-200 bg-ink-50 px-3 py-1.5 text-xs text-ink-600 focus:outline-none focus:ring-2 focus:ring-brand-300"
            @change="handleDateChange"
          />
          <span class="text-ink-400">至</span>
          <input
            v-model="filters.endDate"
            type="date"
            class="rounded-lg border border-ink-200 bg-ink-50 px-3 py-1.5 text-xs text-ink-600 focus:outline-none focus:ring-2 focus:ring-brand-300"
            @change="handleDateChange"
          />
        </div>
        <div class="flex items-center gap-2">
          <span class="font-medium text-ink-500">校区：</span>
          <div class="flex gap-1.5">
            <button
              v-for="opt in CAMPUS_OPTIONS"
              :key="opt.value"
              type="button"
              class="filter-pill rounded-lg px-3 py-1.5 text-xs font-medium"
              :class="filters.campus === opt.value ? 'active' : 'bg-ink-50 text-ink-600'"
              @click="setCampus(opt.value)"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 列表 -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-40 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <div v-else-if="list.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20">
      <svg class="mb-4 h-12 w-12 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <p class="text-sm text-ink-500">暂无双选会信息</p>
    </div>

    <div v-else class="space-y-5">
      <article
        v-for="item in list"
        :key="item.id"
        class="card-shadow cursor-pointer rounded-2xl bg-white p-6 transition-all"
        :class="{ 'border-2 border-brand-200': item.inCalendar }"
        @click="handleViewDetail(item.id)"
      >
        <div class="flex items-start gap-6">
          <!-- 日期徽章 -->
          <div class="shrink-0 rounded-2xl bg-brand-50 p-4 text-center">
            <div class="text-2xl font-bold text-brand-700 leading-none">
              {{ new Date(item.startDate).getDate() }}
            </div>
            <div class="mt-1 text-xs font-medium text-brand-500">
              {{ new Date(item.startDate).toLocaleDateString('zh-CN', { month: 'short' }) }}
            </div>
          </div>

          <!-- 主体信息 -->
          <div class="min-w-0 flex-1">
            <div class="mb-2 flex items-start justify-between gap-4">
              <h3 class="text-lg font-bold text-ink-900 line-clamp-1">{{ item.title }}</h3>
              <div class="flex shrink-0 gap-2" @click.stop>
                <button
                  v-if="!item.inCalendar"
                  type="button"
                  class="flex items-center gap-1.5 rounded-lg bg-brand-500 px-4 py-2 text-xs font-medium text-white transition-colors hover:bg-brand-600 disabled:opacity-60"
                  :disabled="calendarLoading[item.id]"
                  @click="handleAddToCalendar(item)"
                >
                  <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                  {{ calendarLoading[item.id] ? '添加中…' : '加入日历' }}
                </button>
                <button v-else type="button" disabled class="flex items-center gap-1.5 rounded-lg border border-brand-300 bg-brand-100 px-4 py-2 text-xs font-medium text-brand-700">
                  <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  已加入日历
                </button>
              </div>
            </div>

            <div class="mb-3 flex flex-wrap items-center gap-4 text-sm text-ink-600">
              <span class="flex items-center gap-1">
                <svg class="h-4 w-4 text-ink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                </svg>
                {{ item.location }}
              </span>
              <span v-if="item.companyCount" class="flex items-center gap-1">
                <svg class="h-4 w-4 text-ink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5" />
                </svg>
                {{ item.companyCount }} 家企业
              </span>
              <span v-if="item.endDate" class="text-ink-500">
                {{ formatDate(item.startDate) }} — {{ formatDate(item.endDate) }}
              </span>
              <span v-else>{{ formatDate(item.startDate) }}</span>
            </div>

            <div class="flex flex-wrap items-center gap-2">
              <span v-if="item.targetAudience" class="rounded-md bg-ink-100 px-2.5 py-1 text-xs text-ink-600">
                面向：{{ item.targetAudience }}
              </span>
              <!-- 报名截止 -->
              <template v-if="item.deadline">
                <span
                  v-if="isDeadlinePassed(item.deadline)"
                  class="rounded-md bg-ink-100 px-2.5 py-1 text-xs text-ink-400"
                >
                  报名已截止
                </span>
                <span
                  v-else-if="isDeadlineSoon(item.deadline)"
                  class="rounded-md bg-red-50 px-2.5 py-1 text-xs font-medium text-red-600"
                >
                  ⚡ 截止：{{ formatDate(item.deadline) }}
                </span>
                <span v-else class="rounded-md bg-warm-50 px-2.5 py-1 text-xs text-warm-300">
                  报名截止：{{ formatDate(item.deadline) }}
                </span>
              </template>
              <a
                v-if="item.detailUrl"
                :href="item.detailUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="ml-auto flex items-center gap-1 rounded-lg border border-ink-200 px-3 py-1 text-xs text-ink-600 transition-colors hover:border-brand-300 hover:text-brand-600"
                @click.stop
              >
                查看详情
                <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </a>
            </div>
          </div>
        </div>
      </article>
    </div>

    <AppPagination :page="page" :total-pages="totalPages" @change="handlePageChange" />
  </div>
</template>