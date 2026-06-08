<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { getCareerTalkUpcomingApi, getCareerTalkHotCompaniesApi } from '@/api/modules/careerTalk'
import { formatUpcomingLabel, getUpcomingColor, normalizeCareerTalk } from '@/utils/careerTalk'
import type { CareerTalk, HotCompany } from '@/types'

const emit = defineEmits<{
  filterCompany: [company: string]
}>()

const upcomingLoading = ref(true)
const hotLoading = ref(true)
const upcomingEvents = ref<CareerTalk[]>([])
const hotCompanies = ref<HotCompany[]>([])

async function fetchSidebar() {
  upcomingLoading.value = true
  hotLoading.value = true
  try {
    const [upcoming, hot] = await Promise.all([
      getCareerTalkUpcomingApi(),
      getCareerTalkHotCompaniesApi(6),
    ])
    upcomingEvents.value = upcoming.map(normalizeCareerTalk)
    hotCompanies.value = hot
  } catch {
    upcomingEvents.value = []
    hotCompanies.value = []
  } finally {
    upcomingLoading.value = false
    hotLoading.value = false
  }
}

function handleCompanyClick(company: string) {
  emit('filterCompany', company)
}

onMounted(fetchSidebar)
</script>

<template>
  <aside class="w-80 shrink-0 space-y-6">
    <!-- 即将开始 -->
    <div class="sidebar-card rounded-2xl bg-white p-5">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-sm font-bold text-ink-900">即将开始</h3>
        <RouterLink to="/calendar" class="text-xs font-medium text-brand-600 hover:text-brand-700">
          查看全部
        </RouterLink>
      </div>

      <div v-if="upcomingLoading" class="space-y-3">
        <div v-for="i in 2" :key="i" class="h-14 animate-pulse rounded-xl bg-ink-100" />
      </div>

      <p v-else-if="upcomingEvents.length === 0" class="py-4 text-center text-xs text-ink-400">
        24 小时内暂无即将开始的宣讲会
      </p>

      <div v-else class="space-y-3">
        <RouterLink
          v-for="(event, index) in upcomingEvents"
          :key="event.id"
          :to="`/career-talks/${event.id}`"
          class="group flex cursor-pointer gap-3 rounded-xl p-3 transition-colors hover:bg-ink-50"
        >
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg"
            :class="getUpcomingColor(index).split(' ')[0]"
          >
            <span class="text-xs font-bold" :class="getUpcomingColor(index).split(' ')[1]">
              {{ event.company.slice(0, 2) }}
            </span>
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-ink-900 transition-colors group-hover:text-brand-600">
              {{ event.title }}
            </p>
            <p class="mt-0.5 text-xs text-ink-500">
              {{ formatUpcomingLabel(event.startTime, event.campus) }}
            </p>
          </div>
        </RouterLink>
      </div>
    </div>

    <!-- 热门公司 -->
    <div class="sidebar-card rounded-2xl bg-white p-5">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-sm font-bold text-ink-900">热门公司</h3>
      </div>

      <div v-if="hotLoading" class="grid grid-cols-2 gap-2">
        <div v-for="i in 4" :key="i" class="h-12 animate-pulse rounded-xl bg-ink-100" />
      </div>

      <p v-else-if="hotCompanies.length === 0" class="py-4 text-center text-xs text-ink-400">
        暂无热门公司数据
      </p>

      <div v-else class="grid grid-cols-2 gap-2">
        <button
          v-for="item in hotCompanies"
          :key="item.company"
          type="button"
          class="flex cursor-pointer items-center gap-2 rounded-xl border border-ink-100 p-2.5 transition-all hover:border-brand-200 hover:bg-brand-50"
          @click="handleCompanyClick(item.company)"
        >
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-100">
            <span class="text-[10px] font-bold text-brand-700">{{ item.company.slice(0, 2) }}</span>
          </div>
          <span class="truncate text-xs font-medium text-ink-700">{{ item.company }}</span>
        </button>
      </div>
    </div>

    <!-- 我的日历快捷入口 -->
    <div class="rounded-2xl bg-gradient-to-br from-brand-500 to-brand-700 p-5 text-white">
      <div class="mb-3 flex items-center gap-2">
        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
          />
        </svg>
        <h3 class="text-sm font-bold">我的日历</h3>
      </div>
      <p class="mb-4 text-xs text-brand-100">本周还有 3 场宣讲会待参加</p>
      <RouterLink
        to="/calendar"
        class="inline-flex items-center gap-1.5 rounded-lg bg-white/20 px-4 py-2 text-xs font-medium transition-colors hover:bg-white/30"
      >
        查看日程
      </RouterLink>
    </div>

    <!-- 邮件提醒 -->
    <div class="rounded-2xl border border-warm-200 bg-warm-50 p-5">
      <div class="mb-2 flex items-center gap-2">
        <svg class="h-5 w-5 text-warm-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
          />
        </svg>
        <h3 class="text-sm font-bold text-ink-900">邮件提醒</h3>
      </div>
      <p class="mb-3 text-xs text-ink-500">订阅感兴趣的公司和岗位，宣讲会开始前自动邮件提醒</p>
      <RouterLink
        to="/profile/preferences"
        class="block w-full rounded-xl bg-ink-900 px-4 py-2.5 text-center text-xs font-medium text-white transition-colors hover:bg-ink-800"
      >
        管理订阅偏好
      </RouterLink>
    </div>
  </aside>
</template>
