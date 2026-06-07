<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCareerTalkDetailApi } from '@/api/modules/careerTalk'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import type { CareerTalk } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const calendarLoading = ref(false)
const talk = ref<CareerTalk | null>(null)

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

async function fetchDetail() {
  loading.value = true
  try {
    talk.value = await getCareerTalkDetailApi(Number(route.params.id))
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
  if (!talk.value || talk.value.inCalendar) return
  calendarLoading.value = true
  try {
    await addCalendarEventApi({ eventType: 'career_talk', refId: talk.value.id })
    talk.value.inCalendar = true
    ElMessage.success('已加入日历')
  } catch (e: any) {
    if (e?.response?.data?.code === 40901) {
      if (talk.value) talk.value.inCalendar = true
      ElMessage.info('该活动已在您的日历中')
    }
  } finally {
    calendarLoading.value = false
  }
}

function formatDateTime(dateStr?: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric', month: 'long', day: 'numeric',
    weekday: 'short', hour: '2-digit', minute: '2-digit',
  })
}

onMounted(fetchDetail)
</script>

<template>
  <div>
    <div v-if="loading" class="space-y-4">
      <div class="h-8 w-1/3 animate-pulse rounded-lg bg-ink-100" />
      <div class="h-56 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <template v-else-if="talk">
      <nav class="mb-6 flex items-center gap-2 text-sm text-ink-500">
        <RouterLink to="/career-talks" class="transition-colors hover:text-brand-600">宣讲会</RouterLink>
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <span class="truncate text-ink-800">{{ talk.title }}</span>
      </nav>

      <div class="grid gap-6 lg:grid-cols-3">
        <div class="space-y-5 lg:col-span-2">
          <div class="card-shadow rounded-2xl bg-white p-8">
            <!-- 标题区 -->
            <div class="mb-6 flex items-start gap-5">
              <div class="company-logo-placeholder flex h-20 w-20 shrink-0 items-center justify-center rounded-xl border border-brand-100">
                <span class="text-lg font-bold text-brand-700">{{ talk.company.slice(0, 2) }}</span>
              </div>
              <div>
                <h1 class="mb-1 text-xl font-bold text-ink-900">{{ talk.title }}</h1>
                <p class="text-sm text-ink-500">{{ talk.company }} · {{ talk.industry }} · {{ talk.companySize }}</p>
                <div class="mt-2 flex gap-2">
                  <span class="rounded-md px-2.5 py-1 text-xs font-medium" :class="formatTagClass[talk.format]">
                    {{ formatLabel[talk.format] }}
                  </span>
                  <span v-if="talk.status === 'ended'" class="rounded-md bg-ink-100 px-2.5 py-1 text-xs font-medium text-ink-400">
                    已结束
                  </span>
                </div>
              </div>
            </div>

            <!-- 详情网格 -->
            <div class="mb-6 grid gap-4 sm:grid-cols-2">
              <div class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">举办时间</p>
                  <p class="text-sm font-medium text-ink-800">{{ formatDateTime(talk.startTime) }}</p>
                </div>
              </div>
              <div class="flex items-start gap-3 rounded-xl bg-ink-50 p-4">
                <svg class="mt-0.5 h-5 w-5 shrink-0 text-brand-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                </svg>
                <div>
                  <p class="text-xs text-ink-500">举办地点</p>
                  <p class="text-sm font-medium text-ink-800">{{ talk.location }}</p>
                </div>
              </div>
            </div>

            <!-- 面向岗位 -->
            <div v-if="talk.positions?.length" class="mb-4">
              <p class="mb-2 text-sm font-semibold text-ink-800">面向岗位</p>
              <div class="flex flex-wrap gap-2">
                <span v-for="pos in talk.positions" :key="pos" class="rounded-lg bg-brand-50 px-3 py-1.5 text-sm font-medium text-brand-700">
                  {{ pos }}
                </span>
              </div>
            </div>

            <!-- 描述 -->
            <div v-if="(talk as any).description">
              <p class="mb-2 text-sm font-semibold text-ink-800">活动详情</p>
              <p class="whitespace-pre-wrap text-sm leading-relaxed text-ink-700">{{ (talk as any).description }}</p>
            </div>
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="space-y-5">
          <div class="card-shadow rounded-2xl bg-white p-5 space-y-3">
            <template v-if="talk.status !== 'ended'">
              <button
                v-if="!talk.inCalendar"
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
            </template>
            <span v-else class="block w-full rounded-xl bg-ink-100 py-3 text-center text-sm text-ink-400">
              活动已结束
            </span>
            <a
              v-if="talk.sourceUrl"
              :href="talk.sourceUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="flex w-full items-center justify-center gap-2 rounded-xl border border-ink-200 py-3 text-sm font-medium text-ink-700 transition-colors hover:bg-ink-50"
            >
              查看原文
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
            <RouterLink to="/career-talks" class="flex w-full items-center justify-center rounded-xl py-2 text-sm text-ink-500 hover:text-ink-700">
              ← 返回宣讲会列表
            </RouterLink>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="flex flex-col items-center justify-center py-20 text-center">
      <p class="mb-4 text-5xl font-bold text-ink-200">404</p>
      <p class="mb-4 text-ink-500">宣讲会不存在或已下架</p>
      <RouterLink to="/career-talks" class="btn-primary rounded-xl px-6 py-2.5 text-sm font-medium text-white">
        返回宣讲会列表
      </RouterLink>
    </div>
  </div>
</template>