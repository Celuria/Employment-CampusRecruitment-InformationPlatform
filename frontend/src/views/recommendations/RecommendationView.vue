<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import CareerTalkCard from '@/components/business/CareerTalkCard.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import { getRecommendationsApi } from '@/api/modules/recommendation'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import { usePagination } from '@/composables/usePagination'
import { toCareerTalk, toJobFair } from '@/utils/recommendation'
import type { RecommendationItem } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const { page, pageSize, total, totalPages, setTotal } = usePagination(20)

const loading = ref(true)
const list = ref<RecommendationItem[]>([])
const fallback = ref(false)
const calendarLoading = ref<Record<number, boolean>>({})

const profileIncomplete = computed(() => {
  const info = authStore.userInfo
  return !info || !info.college || !info.major || !info.targetPositions?.length
})

async function fetchRecommendations() {
  loading.value = true
  try {
    const res = await getRecommendationsApi({ page: page.value, pageSize: pageSize.value })
    list.value = res.list
    fallback.value = res.fallback
    setTotal(Number(res.total))
  } catch {
    list.value = []
    fallback.value = false
    setTotal(0)
  } finally {
    loading.value = false
  }
}

async function handleAddToCalendar(item: RecommendationItem) {
  if (item.inCalendar) {
    ElMessage.info('该活动已在您的日历中')
    return
  }
  calendarLoading.value[item.refId] = true
  try {
    await addCalendarEventApi({ eventType: item.eventType, refId: item.refId })
    item.inCalendar = true
    ElMessage.success('已加入日历')
  } catch (e: any) {
    if (e?.response?.data?.code === 40901) {
      item.inCalendar = true
      ElMessage.info('该活动已在您的日历中')
    }
  } finally {
    calendarLoading.value[item.refId] = false
  }
}

function handleViewDetail(item: RecommendationItem) {
  const path = item.eventType === 'career_talk'
    ? `/career-talks/${item.refId}`
    : `/job-fairs/${item.refId}`
  router.push(path)
}

function handlePageChange(newPage: number) {
  page.value = newPage
  fetchRecommendations()
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
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

onMounted(fetchRecommendations)
</script>

<template>
  <div>
    <PageHeader
      title="个性化推荐"
      description="根据您的学院、专业与偏好，为您匹配最合适的校招活动"
      :count="total"
      count-label="条推荐"
    />

    <!-- 资料不完整引导 -->
    <div v-if="profileIncomplete" class="mb-6 flex items-start gap-4 rounded-2xl border border-brand-200 bg-brand-50 p-5">
      <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <div class="flex-1">
        <p class="text-sm font-medium text-brand-800">完善个人资料，获得更精准的推荐</p>
        <p class="mt-1 text-xs text-brand-600">请填写学院、专业和意向岗位，系统将为您匹配更多相关校招活动。</p>
      </div>
      <RouterLink
        to="/profile/info"
        class="shrink-0 rounded-xl bg-brand-500 px-4 py-2 text-xs font-medium text-white transition-colors hover:bg-brand-600"
      >
        去完善
      </RouterLink>
    </div>

    <!-- 无精准匹配提示 -->
    <div v-if="!loading && fallback && list.length > 0" class="mb-6 flex items-start gap-3 rounded-2xl border border-ink-200 bg-ink-50 p-4">
      <svg class="mt-0.5 h-5 w-5 shrink-0 text-ink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-ink-600">
        暂未找到与您的资料与偏好精准匹配的活动，以下为您展示全部宣讲会与双选会。
      </p>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="space-y-5">
      <div v-for="i in 4" :key="i" class="h-40 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <!-- 空状态 -->
    <div v-else-if="list.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20">
      <svg class="mb-4 h-12 w-12 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
      </svg>
      <p class="mb-2 text-base font-medium text-ink-600">暂无校招活动</p>
      <p class="mb-4 max-w-xs text-center text-sm text-ink-400">
        {{ profileIncomplete ? '请先完善个人资料和偏好设置，以便系统为您生成推荐。' : '当前数据库中暂无已发布的宣讲会或双选会。' }}
      </p>
      <RouterLink v-if="profileIncomplete" to="/profile/preferences" class="btn-primary rounded-xl px-5 py-2.5 text-sm font-medium text-white">
        设置偏好
      </RouterLink>
    </div>

    <!-- 推荐列表 -->
    <div v-else class="space-y-5">
      <template v-for="item in list" :key="`${item.eventType}-${item.refId}`">
        <!-- 宣讲会卡片 -->
        <div v-if="item.eventType === 'career_talk'" class="space-y-2">
          <CareerTalkCard
            :item="toCareerTalk(item)"
            @add-to-calendar="() => handleAddToCalendar(item)"
            @toggle-favorite="() => {}"
            @view-detail="() => handleViewDetail(item)"
          />
          <div v-if="item.matchReasons?.length" class="flex flex-wrap gap-2 px-1">
            <span
              v-for="reason in item.matchReasons"
              :key="reason"
              class="flex items-center gap-1 rounded-lg bg-brand-50 px-2.5 py-1 text-xs text-brand-700"
            >
              <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              {{ reason }}
            </span>
          </div>
        </div>

        <!-- 双选会卡片 -->
        <div v-else class="space-y-2">
          <article
            class="card-shadow cursor-pointer rounded-2xl bg-white p-6 transition-all"
            :class="{ 'border-2 border-brand-200': item.inCalendar }"
            @click="handleViewDetail(item)"
          >
            <div class="flex items-start gap-6">
              <div class="shrink-0 rounded-2xl bg-brand-50 p-4 text-center">
                <div class="text-2xl font-bold leading-none text-brand-700">
                  {{ new Date(toJobFair(item).startDate).getDate() }}
                </div>
                <div class="mt-1 text-xs font-medium text-brand-500">
                  {{ new Date(toJobFair(item).startDate).toLocaleDateString('zh-CN', { month: 'short' }) }}
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <div class="mb-2 flex items-start justify-between gap-4">
                  <div class="flex min-w-0 items-center gap-2">
                    <span class="shrink-0 rounded-md bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700">双选会</span>
                    <h3 class="truncate text-lg font-bold text-ink-900">{{ item.title }}</h3>
                  </div>
                  <div class="flex shrink-0 gap-2" @click.stop>
                    <button
                      v-if="!item.inCalendar"
                      type="button"
                      class="flex items-center gap-1.5 rounded-lg bg-brand-500 px-4 py-2 text-xs font-medium text-white transition-colors hover:bg-brand-600 disabled:opacity-60"
                      :disabled="calendarLoading[item.refId]"
                      @click="handleAddToCalendar(item)"
                    >
                      {{ calendarLoading[item.refId] ? '添加中…' : '加入日历' }}
                    </button>
                    <button v-else type="button" disabled class="flex items-center gap-1.5 rounded-lg border border-brand-300 bg-brand-100 px-4 py-2 text-xs font-medium text-brand-700">
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
                    {{ formatDate(item.startDate!) }} — {{ formatDate(item.endDate) }}
                  </span>
                  <span v-else-if="item.startDate">{{ formatDate(item.startDate) }}</span>
                </div>

                <div class="flex flex-wrap items-center gap-2">
                  <span v-if="item.targetAudience" class="rounded-md bg-ink-100 px-2.5 py-1 text-xs text-ink-600">
                    面向：{{ item.targetAudience }}
                  </span>
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
                </div>
              </div>
            </div>
          </article>

          <div v-if="item.matchReasons?.length" class="flex flex-wrap gap-2 px-1">
            <span
              v-for="reason in item.matchReasons"
              :key="reason"
              class="flex items-center gap-1 rounded-lg bg-brand-50 px-2.5 py-1 text-xs text-brand-700"
            >
              <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              {{ reason }}
            </span>
          </div>
        </div>
      </template>
    </div>

    <AppPagination
      v-if="!loading && list.length > 0"
      :page="page"
      :total-pages="totalPages"
      @change="handlePageChange"
    />
  </div>
</template>
