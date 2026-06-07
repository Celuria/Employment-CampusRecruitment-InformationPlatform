<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getJobFairDetailApi } from '@/api/modules/jobFair'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import type { JobFair } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const calendarLoading = ref(false)
const jobFair = ref<JobFair | null>(null)

async function fetchDetail() {
  loading.value = true
  try {
    jobFair.value = await getJobFairDetailApi(Number(route.params.id))
  } catch {
    // 错误由拦截器处理
  } finally {
    loading.value = false
  }
}

async function handleAddToCalendar() {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: route.fullPath } })
    return
  }
  if (!jobFair.value || jobFair.value.inCalendar) return
  calendarLoading.value = true
  try {
    await addCalendarEventApi({ eventType: 'job_fair', refId: jobFair.value.id })
    jobFair.value.inCalendar = true
    ElMessage.success('已加入日历')
  } catch (e: any) {
    if (e?.response?.data?.code === 40901) {
      if (jobFair.value) jobFair.value.inCalendar = true
      ElMessage.info('该活动已在您的日历中')
    }
  } finally {
    calendarLoading.value = false
  }
}

function formatDate(dateStr?: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric', month: 'long', day: 'numeric', weekday: 'short',
  })
}

function isDeadlineSoon(deadline?: string) {
  if (!deadline) return false
  const diff = new Date(deadline).getTime() - Date.now()
  return diff > 0 && diff < 3 * 24 * 3600 * 1000
}

const deadlinePassed = computed(() => {
  if (!jobFair.value?.deadline) return false
  return new Date(jobFair.value.deadline).getTime() < Date.now()
})

onMounted(fetchDetail)
</script>

<template>
  <div>
    <!-- 加载骨架 -->
    <div v-if="loading" class="space-y-4">
      <div class="h-8 w-1/3 animate-pulse rounded-lg bg-ink-100" />
      <div class="h-48 animate-pulse rounded-2xl bg-ink-100" />
      <div class="h-32 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <template v-else-if="jobFair">
      <!-- 面包屑 -->
      <nav class="mb-6 flex items-center gap-2 text-sm text-ink-500">
        <RouterLink to="/job-fairs" class="transition-colors hover:text-brand-600">双选会</RouterLink>
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="truncate text-ink-800">{{ jobFair.title }}</span>
      </nav>

      <div class="grid gap-6 lg:grid-cols-3">
        <!-- 主信息 -->
        <div class="space-y-5 lg:col-span-2">
          <div class="card-shadow rounded-2xl bg-white p-8">
            <h1 class="mb-4 text-2xl font-bold text-ink-900">{{ jobFair.title }}</h1>

            <div class="mb-6 grid gap-4 sm:grid-cols-2">
              <div class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">举办日期</p>
                  <p class="text-sm font-medium text-ink-800">{{ formatDate(jobFair.startDate) }}</p>
                  <p v-if="jobFair.endDate" class="text-xs text-ink-500">至 {{ formatDate(jobFair.endDate) }}</p>
                </div>
              </div>
              <div class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">举办地点</p>
                  <p class="text-sm font-medium text-ink-800">{{ jobFair.location }}</p>
                </div>
              </div>
              <div v-if="jobFair.companyCount" class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">参与企业</p>
                  <p class="text-sm font-medium text-ink-800">{{ jobFair.companyCount }} 家</p>
                </div>
              </div>
              <div v-if="jobFair.targetAudience" class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">面向对象</p>
                  <p class="text-sm font-medium text-ink-800">{{ jobFair.targetAudience }}</p>
                </div>
              </div>
            </div>

            <!-- 描述 -->
            <div v-if="(jobFair as any).description" class="prose prose-sm max-w-none text-ink-700">
              <h3 class="mb-2 text-base font-semibold text-ink-900">活动详情</h3>
              <p class="whitespace-pre-wrap leading-relaxed">{{ (jobFair as any).description }}</p>
            </div>
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="space-y-5">
          <!-- 报名截止 -->
          <div
            v-if="jobFair.deadline"
            class="rounded-2xl p-5"
            :class="deadlinePassed ? 'bg-ink-100' : isDeadlineSoon(jobFair.deadline) ? 'border border-red-200 bg-red-50' : 'border border-warm-200 bg-warm-50'"
          >
            <h4 class="mb-1 text-sm font-bold" :class="deadlinePassed ? 'text-ink-500' : isDeadlineSoon(jobFair.deadline) ? 'text-red-700' : 'text-ink-800'">
              报名截止时间
            </h4>
            <p class="text-sm" :class="deadlinePassed ? 'text-ink-400' : isDeadlineSoon(jobFair.deadline) ? 'font-semibold text-red-600' : 'text-ink-700'">
              {{ formatDate(jobFair.deadline) }}
            </p>
            <p v-if="deadlinePassed" class="mt-1 text-xs text-ink-400">报名已截止</p>
            <p v-else-if="isDeadlineSoon(jobFair.deadline)" class="mt-1 text-xs font-medium text-red-500">⚡ 即将截止，请尽快报名</p>
          </div>

          <!-- 操作区 -->
          <div class="card-shadow rounded-2xl bg-white p-5 space-y-3">
            <button
              v-if="!jobFair.inCalendar"
              type="button"
              class="btn-primary w-full rounded-xl py-3 text-sm font-medium text-white disabled:opacity-60"
              :disabled="calendarLoading"
              @click="handleAddToCalendar"
            >
              {{ calendarLoading ? '添加中…' : '加入我的日历' }}
            </button>
            <button v-else type="button" disabled class="w-full rounded-xl border border-brand-300 bg-brand-50 py-3 text-sm font-medium text-brand-700">
              ✓ 已加入日历
            </button>
            <a
              v-if="jobFair.detailUrl"
              :href="jobFair.detailUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="flex w-full items-center justify-center gap-2 rounded-xl border border-ink-200 py-3 text-sm font-medium text-ink-700 transition-colors hover:bg-ink-50"
            >
              前往报名页面
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
            <RouterLink to="/job-fairs" class="flex w-full items-center justify-center rounded-xl py-2 text-sm text-ink-500 hover:text-ink-700">
              ← 返回双选会列表
            </RouterLink>
          </div>
        </div>
      </div>
    </template>

    <!-- 404 -->
    <div v-else class="flex flex-col items-center justify-center py-20 text-center">
      <p class="mb-4 text-5xl font-bold text-ink-200">404</p>
      <p class="mb-4 text-ink-500">双选会不存在或已下架</p>
      <RouterLink to="/job-fairs" class="btn-primary rounded-xl px-6 py-2.5 text-sm font-medium text-white">
        返回双选会列表
      </RouterLink>
    </div>
  </div>
</template>