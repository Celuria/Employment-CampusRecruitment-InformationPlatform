<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import { getRecommendationsApi } from '@/api/modules/recommendation'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import type { RecommendationItem } from '@/types'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const list = ref<RecommendationItem[]>([])
const calendarLoading = ref<Record<number, boolean>>({})

const profileIncomplete = computed(() => {
  const info = authStore.userInfo
  return !info || !info.college || !info.major || !info.targetPositions?.length
})

async function fetchRecommendations() {
  loading.value = true
  try {
    const res = await getRecommendationsApi()
    list.value = Array.isArray(res) ? res : (res as any).list || []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function handleAddToCalendar(item: RecommendationItem) {
  if (item.inCalendar) {
    ElMessage.info('该活动已在您的日历中')
    return
  }
  calendarLoading.value[item.id] = true
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
    calendarLoading.value[item.id] = false
  }
}

function handleViewDetail(item: RecommendationItem) {
  const path = item.eventType === 'career_talk' ? `/career-talks/${item.refId}` : `/job-fairs/${item.refId}`
  router.push(path)
}

function formatDateTime(str: string) {
  return new Date(str).toLocaleString('zh-CN', {
    year: 'numeric', month: 'short', day: 'numeric',
    weekday: 'short', hour: '2-digit', minute: '2-digit',
  })
}

const EVENT_TYPE_LABEL: Record<string, string> = {
  career_talk: '宣讲会',
  job_fair: '双选会',
}

onMounted(fetchRecommendations)
</script>

<template>
  <div>
    <PageHeader
      title="个性化推荐"
      description="根据您的学院、专业与偏好，为您匹配最合适的校招活动"
      :count="list.length"
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

    <!-- 加载中 -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 4" :key="i" class="h-36 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <!-- 空状态 -->
    <div v-else-if="list.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20">
      <svg class="mb-4 h-12 w-12 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
      </svg>
      <p class="mb-2 text-base font-medium text-ink-600">暂无推荐</p>
      <p class="mb-4 max-w-xs text-center text-sm text-ink-400">
        {{ profileIncomplete ? '请先完善个人资料和偏好设置，以便系统为您生成推荐。' : '当前暂无匹配的校招活动，请稍后再来查看。' }}
      </p>
      <RouterLink v-if="profileIncomplete" to="/profile/preferences" class="btn-primary rounded-xl px-5 py-2.5 text-sm font-medium text-white">
        设置偏好
      </RouterLink>
    </div>

    <!-- 推荐列表 -->
    <div v-else class="space-y-5">
      <article
        v-for="item in list"
        :key="item.id"
        class="card-shadow group cursor-pointer rounded-2xl bg-white p-6 transition-all"
        :class="{ 'border-2 border-brand-200': item.inCalendar }"
        @click="handleViewDetail(item)"
      >
        <div class="flex items-start gap-5">
          <!-- 公司徽标 -->
          <div class="company-logo-placeholder flex h-16 w-16 shrink-0 items-center justify-center rounded-xl border border-brand-100">
            <span class="text-sm font-bold text-brand-700">{{ item.company?.slice(0, 2) || '活动' }}</span>
          </div>

          <div class="min-w-0 flex-1">
            <div class="mb-1 flex items-start justify-between gap-3">
              <div>
                <span class="mr-2 rounded-md px-2 py-0.5 text-xs font-medium"
                  :class="item.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700' : 'bg-blue-100 text-blue-700'"
                >
                  {{ EVENT_TYPE_LABEL[item.eventType] }}
                </span>
                <h3 class="inline text-base font-bold text-ink-900 transition-colors group-hover:text-brand-600">
                  {{ item.title }}
                </h3>
              </div>
              <div class="flex shrink-0 gap-2" @click.stop>
                <button
                  type="button"
                  class="rounded-lg border border-ink-200 px-3 py-1.5 text-xs text-ink-600 transition-colors hover:border-brand-300 hover:text-brand-600"
                  @click="handleViewDetail(item)"
                >
                  查看详情
                </button>
                <button
                  v-if="!item.inCalendar"
                  type="button"
                  class="flex items-center gap-1 rounded-lg bg-brand-500 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-brand-600 disabled:opacity-60"
                  :disabled="calendarLoading[item.id]"
                  @click="handleAddToCalendar(item)"
                >
                  {{ calendarLoading[item.id] ? '添加中…' : '加入日历' }}
                </button>
                <button v-else type="button" disabled class="rounded-lg border border-brand-300 bg-brand-50 px-3 py-1.5 text-xs font-medium text-brand-700">
                  ✓ 已加入
                </button>
              </div>
            </div>

            <p class="mb-2 text-sm text-ink-500">
              {{ item.company }} · {{ formatDateTime(item.startTime) }}
            </p>
            <p v-if="item.location" class="mb-2 text-sm text-ink-500">📍 {{ item.location }}</p>

            <!-- 匹配原因 -->
            <div v-if="item.matchReasons?.length" class="flex flex-wrap gap-2">
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
        </div>
      </article>
    </div>
  </div>
</template>